"""bin/crow-check: each gate group proven by a failing fixture, green on a clean tree."""
import json
import os
import shutil
import subprocess
import sys

import pytest
from conftest import BIN

CHECK = str(BIN / "crow-check")
HOME_ALICE = "/ho" + "me/alice"  # built at runtime so this file passes the leak check


def git(root, *args):
    subprocess.run(["git", "-c", "user.email=t@example.com", "-c", "user.name=t", *args],
                   cwd=root, check=True, capture_output=True)


@pytest.fixture
def tree(tmp_path):
    root = tmp_path / "proj"
    (root / "tests").mkdir(parents=True)
    (root / "tests" / "test_ok.py").write_text("def test_ok():\n    assert True\n")
    (root / "README.md").write_text("hello\n")
    git(root, "init", "-q", "-b", "main")
    git(root, "add", ".")
    git(root, "commit", "-q", "-m", "new[MD]: first")
    git(root, "tag", "v0.0.1")
    return root


def check(root, *args, env=None, path=None):
    e = dict(os.environ, **(env or {}))
    e.pop("CROW_CHECK_DENY", None) if not env else None
    if path is not None:
        e["PATH"] = path
    r = subprocess.run([sys.executable, CHECK, "--root", str(root), "--json", *args], check=False,
                       capture_output=True, text=True, env=e)
    return r.returncode, {g["group"]: g for g in json.loads(r.stdout)["groups"]}


def test_clean_tree_is_green(tree):
    code, g = check(tree)
    assert code == 0, g
    assert g["pytest"]["status"] == "pass"
    assert g["leak"]["status"] == "pass"
    assert g["commits"]["status"] == "pass"
    assert g["version"]["status"] == "pass"


def test_runs_without_local_dir(tree):
    assert not (tree / "local").exists()
    assert check(tree, "--fast")[0] == 0


def test_fast_skips_pytest(tree):
    (tree / "tests" / "test_ok.py").write_text("def test_bad():\n    assert False\n")
    code, g = check(tree, "--fast")
    assert code == 0 and g["pytest"]["status"] == "skip"


def test_pytest_failure_is_red(tree):
    (tree / "tests" / "test_ok.py").write_text("def test_bad():\n    assert False\n")
    code, g = check(tree)
    assert code == 1 and g["pytest"]["status"] == "fail"
    assert g["leak"]["status"] == "pass"  # every group still reports


def test_leak_planted_name_fails(tree, tmp_path):
    deny = tmp_path / "deny.txt"
    deny.write_text("# private names\nacmecorp\n")
    (tree / "notes.md").write_text("ship it to AcmeCorp\n")
    code, g = check(tree, "--fast", env={"CROW_CHECK_DENY": str(deny)})
    assert code == 1 and g["leak"]["status"] == "fail" and "notes.md:1" in g["leak"]["note"]


def test_leak_personal_path_fails(tree):
    (tree / "run.sh").write_text("cd " + HOME_ALICE + "/work\n")
    code, g = check(tree, "--fast")
    assert code == 1 and g["leak"]["status"] == "fail"


def test_leak_placeholder_paths_pass(tree):
    (tree / "doc.md").write_text("/home/me/x /home/user/y /home/$USER/z\n")
    assert check(tree, "--fast")[0] == 0


def test_leak_ignores_local(tree, tmp_path):
    deny = tmp_path / "deny.txt"
    deny.write_text("acmecorp\n")
    (tree / "local").mkdir()
    (tree / "local" / "private.md").write_text("acmecorp at " + HOME_ALICE + "\n")
    code, g = check(tree, "--fast", env={"CROW_CHECK_DENY": str(deny)})
    assert code == 0 and g["leak"]["status"] == "pass"


def test_bad_commit_subject_fails(tree):
    (tree / "a.md").write_text("a\n")
    git(tree, "add", ".")
    git(tree, "commit", "-q", "-m", "update stuff")
    code, g = check(tree, "--fast")
    assert code == 1 and g["commits"]["status"] == "fail"


def test_commits_before_tag_are_ignored(tree):
    (tree / "a.md").write_text("a\n")
    git(tree, "add", ".")
    git(tree, "commit", "-q", "-m", "chg[MD]: more")
    assert check(tree, "--fast")[1]["commits"]["status"] == "pass"


def test_version_bad_semver_fails(tree):
    (tree / "VERSION").write_text("one\n")
    code, g = check(tree, "--fast")
    assert code == 1 and g["version"]["status"] == "fail"


def write_release(root, ver="1.2.0"):
    (root / "VERSION").write_text(ver + "\n")
    (root / "changelogs").mkdir(exist_ok=True)
    (root / "changelogs" / ("CHANGELOGS_" + ver.replace(".", "_") + ".md")).write_text("- x\n")
    subprocess.run([sys.executable, str(BIN / "crow-changelog"), "--root", str(root)], check=True, capture_output=True)


def test_version_missing_changelog_entry_fails(tree):
    (tree / "VERSION").write_text("1.2.0\n")
    assert check(tree, "--fast")[1]["version"]["status"] == "fail"
    write_release(tree)
    assert check(tree, "--fast")[1]["version"]["status"] == "pass"


def test_version_stale_aggregate_fails(tree):
    write_release(tree)
    (tree / "changelogs" / "CHANGELOGS_1_2_0.md").write_text("- changed after the aggregate\n")
    g = check(tree, "--fast")[1]["version"]
    assert g["status"] == "fail" and "stale" in g["note"]


def test_version_carrier_mismatch_fails(tree):
    write_release(tree)
    (tree / "package.json").write_text('{"name": "x", "version": "1.1.0"}\n')
    g = check(tree, "--fast")[1]["version"]
    assert g["status"] == "fail" and "package.json says 1.1.0" in g["note"]
    (tree / "package.json").write_text('{"name": "x", "version": "1.2.0"}\n')
    assert check(tree, "--fast")[1]["version"]["status"] == "pass"


def test_missing_linters_skip_with_note(tree, tmp_path):
    empty = tmp_path / "nobin"
    empty.mkdir()
    (empty / "git").symlink_to(shutil.which("git"))
    (tree / "x.sh").write_text("#!/bin/bash\necho hi\n")
    code, g = check(tree, "--fast", path=str(empty))
    assert code == 0
    assert g["shellcheck"]["status"] == "skip" and "not installed" in g["shellcheck"]["note"]
    assert g["ruff"]["status"] == "skip"


def _fake(tmp_path, name, rc):
    d = tmp_path / "fakebin"
    d.mkdir(exist_ok=True)
    f = d / name
    f.write_text(f"#!/bin/sh\necho {name} complaint\nexit {rc}\n")
    f.chmod(0o755)
    return str(d) + os.pathsep + os.environ["PATH"]


@pytest.mark.parametrize("tool", ["ruff", "shellcheck"])
def test_linter_failure_is_red(tree, tmp_path, tool):
    (tree / "x.sh").write_text("#!/bin/bash\necho hi\n")
    code, g = check(tree, "--fast", path=_fake(tmp_path, tool, 1))
    assert code == 1 and g[tool]["status"] == "fail"


def test_text_output_lists_every_group(tree):
    r = subprocess.run([sys.executable, CHECK, "--root", str(tree), "--fast"], check=False, capture_output=True, text=True)
    for name in ("pytest", "shellcheck", "ruff", "leak", "commits", "version"):
        assert name in r.stdout
    assert r.stdout.strip().endswith("green")
