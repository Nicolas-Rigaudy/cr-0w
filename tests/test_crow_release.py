"""bin/crow-changelog and bin/crow-release-propose."""
import json
import subprocess
import sys

from conftest import BIN


def git(root, *args):
    subprocess.run(["git", "-c", "user.email=t@example.com", "-c", "user.name=t", *args],
                   cwd=root, check=True, capture_output=True)


def run(tool, root, *args):
    return subprocess.run([sys.executable, str(BIN / tool), "--root", str(root), *args], check=False,
                          capture_output=True, text=True)


def test_changelog_aggregates_newest_first_and_check(tmp_path):
    d = tmp_path / "changelogs"
    d.mkdir()
    (d / "CHANGELOGS_0_2_0.md").write_text("- two\n")
    (d / "CHANGELOGS_0_10_0.md").write_text("- ten\n")
    (d / "CHANGELOGS_0_1_0.md").write_text("- one\n")
    assert run("crow-changelog", tmp_path, "--check").returncode == 1  # missing aggregate
    assert run("crow-changelog", tmp_path).returncode == 0
    text = (tmp_path / "CHANGELOG.md").read_text()
    assert text.index("## 0.10.0") < text.index("## 0.2.0") < text.index("## 0.1.0")
    assert run("crow-changelog", tmp_path, "--check").returncode == 0


def test_propose_after_tag_groups_commits_offline(tmp_path):
    git(tmp_path, "init", "-q", "-b", "main")
    (tmp_path / "VERSION").write_text("0.1.0\n")
    (tmp_path / "a").write_text("1")
    git(tmp_path, "add", ".")
    git(tmp_path, "commit", "-q", "-m", "new[MD]: first")
    git(tmp_path, "tag", "v0.1.0")
    (tmp_path / "a").write_text("2")
    git(tmp_path, "commit", "-qam", "fix[PY]: repair the thing")
    p = json.loads(run("crow-release-propose", tmp_path, "--json").stdout)
    assert (p["since"], p["version"], p["commits"]) == ("v0.1.0", "0.1.1", 1)
    assert "### Fixed" in p["notes"] and "repair the thing" in p["notes"]
    (tmp_path / "a").write_text("3")
    git(tmp_path, "commit", "-qam", "new[PY]: add a feature")
    assert json.loads(run("crow-release-propose", tmp_path, "--json").stdout)["version"] == "0.2.0"


def test_propose_without_tag_keeps_version_file(tmp_path):
    git(tmp_path, "init", "-q", "-b", "main")
    (tmp_path / "VERSION").write_text("0.1.0\n")
    git(tmp_path, "add", ".")
    git(tmp_path, "commit", "-q", "-m", "new[MD]: first")
    p = json.loads(run("crow-release-propose", tmp_path, "--json").stdout)
    assert p["since"] is None and p["version"] == "0.1.0"
