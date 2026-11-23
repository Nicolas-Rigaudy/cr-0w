"""User bricks (item_265): CROW_HOME/modules/<name>/module.yaml listed, skills linked, scripts run from crow."""
import json
import os
import sys

from conftest import BIN, run

sys.path.insert(0, str(BIN))
import crow_lib

ROOT = BIN.parent


def hello(home, **extra):
    brick = home / "modules" / "hello"
    (brick / "skills" / "hello").mkdir(parents=True)
    (brick / "skills" / "hello" / "SKILL.md").write_text("---\nname: hello\ndescription: Say hello.\n---\n")
    (brick / "bin").mkdir()
    script = brick / "bin" / "crow-hello"
    script.write_text('#!/bin/sh\necho "hello $*"\n')
    script.chmod(0o755)
    lines = ["summary: Says hello.", "scripts: [crow-hello]", "skills: [hello]"]
    lines += [f"{k}: {v}" for k, v in extra.items()]
    (brick / "module.yaml").write_text("\n".join(lines) + "\n")
    return brick


def test_a_user_brick_is_listed_linked_and_run_from_crow(hub):
    home, profiles = hub
    brick = hello(home)
    out = run(hub, "crow", "modules", "json")
    assert out.returncode == 0, out.stderr
    rows = {m["name"]: m for m in json.loads(out.stdout)["modules"]}
    assert rows["hello"]["on"] and rows["hello"]["source"] == "user"
    assert rows["hello"]["dir"] == str(brick)
    assert rows["forge"]["source"] == "engine"
    text = run(hub, "crow", "modules").stdout
    assert any(line.startswith("on ") and " hello " in line and "[user:" in line for line in text.splitlines())

    claude = profiles / "main" / "claude-home" / ".claude"
    claude.mkdir(parents=True)
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    link = claude / "skills" / "hello"
    assert link.is_symlink() and link.resolve() == (brick / "skills" / "hello").resolve()

    out = run(hub, "crow", "hello", "world")
    assert out.returncode == 0, out.stderr
    assert out.stdout.strip() == "hello world"


def test_an_off_user_brick_is_unlinked_and_its_script_unknown(hub):
    home, profiles = hub
    hello(home)
    claude = profiles / "main" / "claude-home" / ".claude"
    claude.mkdir(parents=True)
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  hello: false\n")
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    assert not (claude / "skills" / "hello").exists()
    out = run(hub, "crow", "hello")
    assert out.returncode == 1 and "unknown command 'hello'" in out.stderr


def test_a_user_brick_named_like_an_engine_brick_is_refused(hub):
    home, _ = hub
    (home / "modules" / "forge").mkdir(parents=True)
    (home / "modules" / "forge" / "module.yaml").write_text("summary: mine\n")
    st = crow_lib.load_modules(str(home), str(ROOT))
    err = next(e for e in st["errors"] if e["key"] == "modules.forge")
    assert "clashes with the engine brick" in err["error"] and "rename its folder" in err["error"]
    assert err["file"] == str(home / "modules" / "forge" / "module.yaml")
    assert st["modules"]["forge"]["source"] == "engine"
    out = run(hub, "crow-doctor")
    assert "clashes with the engine brick" in out.stdout + out.stderr


def test_a_user_brick_is_off_when_a_brick_it_requires_is_off(hub):
    home, _ = hub
    hello(home, requires="[design]")
    st = crow_lib.load_modules(str(home), str(ROOT))
    assert not st["modules"]["hello"]["on"]
    assert any(e["key"] == "modules.hello" and "requires design" in e["error"] for e in st["errors"])
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  design: true\n")
    assert crow_lib.load_modules(str(home), str(ROOT))["modules"]["hello"]["on"]


def test_user_bricks_live_under_local_when_crow_home_is_the_checkout(tmp_path):
    assert crow_lib.user_modules_dir(str(ROOT), str(ROOT)) == os.path.join(str(ROOT), "local", "modules")
    assert crow_lib.user_modules_dir(str(tmp_path), str(ROOT)) == str(tmp_path / "modules")


def test_no_user_bricks_changes_nothing(tmp_path):
    st = crow_lib.load_modules(str(tmp_path), str(ROOT))
    assert set(st["modules"]) == set(crow_lib.engine_manifests(str(ROOT)))
    assert st["errors"] == []
    assert crow_lib.user_script("hello", str(tmp_path), str(ROOT)) is None


def git(cwd, *args):
    import subprocess
    subprocess.run(["git", "-C", str(cwd), *args], check=True, capture_output=True,
                   env={**os.environ, "GIT_AUTHOR_NAME": "t", "GIT_AUTHOR_EMAIL": "t@t", "GIT_COMMITTER_NAME": "t",
                        "GIT_COMMITTER_EMAIL": "t@t"})


def brick_repo(tmp_path):
    """A git repo holding a user brick, tagged v1 then v2."""
    repo = tmp_path / "src" / "greet.git"
    hello(repo.parent / "work")
    work = repo.parent / "work" / "modules" / "hello"
    git(work, "init", "-q", "-b", "main")
    git(work, "add", "-A")
    git(work, "commit", "-q", "-m", "v1", "--no-verify")
    git(work, "tag", "v1")
    (work / "bin" / "crow-hello").write_text('#!/bin/sh\necho "hello v2 $*"\n')
    git(work, "commit", "-q", "-am", "v2", "--no-verify")
    git(work, "tag", "v2")
    return work


def lock(home):
    import yaml
    p = home / "modules.lock"
    return (yaml.safe_load(p.read_text()) or {}).get("modules") if p.exists() else None


def test_add_update_remove_from_a_git_url(hub, tmp_path):
    home, _ = hub
    work = brick_repo(tmp_path)
    url = f"file://{work}"
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  design: false  # keep\n  greet: false\n")

    out = run(hub, "crow", "modules", "add", url, "--ref", "v1", "--name", "greet")
    assert out.returncode == 0, out.stderr
    sha1 = lock(home)["greet"]["commit"]
    assert lock(home)["greet"] == {"source": url, "ref": "v1", "commit": sha1}
    assert "greet: true" in (home / "crow.yaml").read_text() and "design: false  # keep" in (home / "crow.yaml").read_text()
    assert run(hub, "crow", "hello").stdout.strip() == "hello"
    assert "@ v1" in run(hub, "crow", "modules").stdout

    out = run(hub, "crow", "modules", "update", "greet")
    assert out.returncode == 0, out.stderr
    assert lock(home)["greet"]["ref"] == "v2" and lock(home)["greet"]["commit"] != sha1
    assert run(hub, "crow", "hello").stdout.strip() == "hello v2"

    out = run(hub, "crow", "modules", "remove", "greet")
    assert out.returncode == 0, out.stderr
    assert not (home / "modules" / "greet").exists()
    assert lock(home) is None
    assert "greet" not in (home / "crow.yaml").read_text()
    assert "design: false  # keep" in (home / "crow.yaml").read_text()


def test_add_defaults_to_the_newest_tag_and_refuses_engine_names(hub, tmp_path):
    home, _ = hub
    url = f"file://{brick_repo(tmp_path)}"
    out = run(hub, "crow", "modules", "add", url, "--name", "forge")
    assert out.returncode == 2 and "clashes with the engine brick" in out.stderr
    assert not (home / "modules" / "forge").exists()
    assert run(hub, "crow", "modules", "add", url, "--name", "greet").returncode == 0
    assert lock(home)["greet"]["ref"] == "v2"
    assert not (home / "crow.yaml").exists()  # on by default: crow.yaml untouched


def test_doctor_flags_a_clone_whose_head_differs_from_the_lock(hub, tmp_path):
    home, _ = hub
    url = f"file://{brick_repo(tmp_path)}"
    assert run(hub, "crow", "modules", "add", url, "--ref", "v2", "--name", "greet").returncode == 0
    assert crow_lib.lock_problems(str(home), str(ROOT)) == []
    git(home / "modules" / "greet", "checkout", "-q", "v1")
    problems = crow_lib.lock_problems(str(home), str(ROOT))
    assert len(problems) == 1 and "differs from the lock" in problems[0]
    out = run(hub, "crow-doctor")
    assert "out of step with modules.lock" in out.stdout
    (home / "modules" / "greet" / "bin" / "crow-hello").write_text("changed")
    assert any("local changes" in p for p in crow_lib.lock_problems(str(home), str(ROOT)))
