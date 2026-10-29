"""crow-doctor: identity, accounts and stale worktrees, read only."""
import json
import subprocess

from conftest import run


def git(*args):
    subprocess.run(["git", *args], check=True, capture_output=True)


def make_repo(path, email):
    git("init", "-q", "-b", "main", str(path))
    git("-C", str(path), "config", "user.email", email)
    git("-C", str(path), "config", "user.name", "T")
    (path / "a.txt").write_text("x")
    git("-C", str(path), "add", "a.txt")
    git("-C", str(path), "commit", "-q", "-m", "first")


def registry(home, repo, identity="me@example.org"):
    (home / "local" / "projects.yaml").write_text(
        f"identities: {{work: {identity}}}\nprojects:\n"
        f"  - {{name: app, scope: work, path: {repo}, forge: none}}\n")


def checks(out):
    return {c["area"]: c for c in json.loads(out.stdout)["checks"]}


def test_doctor_flags_an_identity_mismatch(hub, tmp_path):
    home, _ = hub
    repo = tmp_path / "app"
    make_repo(repo, "other@example.org")
    registry(home, repo)
    out = run(hub, "crow-doctor", "--json")
    assert out.returncode == 1
    ident = checks(out)["identity"]
    assert ident["level"] == "fail" and "other@example.org" in ident["details"][0]


def test_doctor_passes_when_identity_matches(hub, tmp_path):
    home, _ = hub
    repo = tmp_path / "app"
    make_repo(repo, "me@example.org")
    registry(home, repo)
    out = run(hub, "crow-doctor", "--json")
    assert checks(out)["identity"]["level"] == "ok"
    assert checks(out)["worktrees"]["level"] == "ok"


def test_doctor_lists_a_merged_worktree(hub, tmp_path):
    home, _ = hub
    repo = tmp_path / "app"
    make_repo(repo, "me@example.org")
    git("-C", str(repo), "worktree", "add", "-q", "-b", "done", str(tmp_path / "app.wt" / "done"))
    registry(home, repo)
    wt = checks(run(hub, "crow-doctor", "--json"))["worktrees"]
    assert wt["level"] == "warn" and wt["details"] == ["app: done"]


def test_doctor_reads_the_default_account(hub, tmp_path):
    home, _ = hub
    registry(home, tmp_path / "missing")
    acc = checks(run(hub, "crow-doctor", "--json"))["accounts"]
    assert acc["level"] in ("ok", "fail")
    assert "main: use, memory permanent" in acc["details"]


def test_doctor_text_and_bad_argument(hub, tmp_path):
    home, _ = hub
    registry(home, tmp_path / "missing")
    out = run(hub, "crow-doctor")
    assert out.stdout.startswith(("ok", "warn", "FAIL")) and "identity" in out.stdout
    assert run(hub, "crow-doctor", "--nope").returncode != 0
