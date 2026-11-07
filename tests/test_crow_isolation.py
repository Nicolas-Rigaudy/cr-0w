"""Isolated runs (item_100): preflight refusals, the host-side commit of the run diff, dispatch wiring."""
import importlib.machinery
import importlib.util
import os
import subprocess
import sys
import types

import pytest

from conftest import BIN

sys.path.insert(0, str(BIN))
import crow_lib  # noqa: E402

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs_isolation", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs_isolation", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)


def check_msg(msg):
    return subprocess.run([str(BIN / "crow-commit-msg-check"), "-"], input=msg, text=True,
                          capture_output=True).returncode == 0


def test_only_task_and_review_isolate():
    assert crow_lib.isolation_mode("task") == "rw"
    assert crow_lib.isolation_mode("review") == "ro"
    assert crow_lib.isolation_mode("retro") is None and crow_lib.isolation_mode("oneshot") is None
    assert "cannot run isolated" in crow_lib.isolation_preflight("claude", "setup-token", kind="retro")


def test_preflight_refuses_missing_pasta_with_the_fix(tmp_path):
    for tool in ("bwrap", "python3"):
        (tmp_path / tool).symlink_to("/bin/true")
    why = crow_lib.isolation_preflight("claude", "setup-token", path=str(tmp_path))
    assert "pasta" in why and "nft" in why and "sudo pacman -S passt nftables" in why
    assert "no shared-network fallback" in why


def test_preflight_refuses_a_full_login_claude_profile(tmp_path):
    for tool in crow_lib.ISOLATION_TOOLS:
        (tmp_path / tool).symlink_to("/bin/true")
    why = crow_lib.isolation_preflight("claude", "full-login", path=str(tmp_path))
    assert "setup token" in why and "cdx login <account>-runs" in why
    assert crow_lib.isolation_preflight("claude", "setup-token", path=str(tmp_path)) is None
    assert crow_lib.isolation_preflight("codex", None, path=str(tmp_path)) is None


def test_isolation_account_prefers_the_runs_profile():
    rows = [{"session_name": "alpha"}, {"session_name": "alpha-runs", "auth_type": "setup-token"}, {"session_name": "beta"}]
    assert crow_lib.isolation_account("alpha", rows) == "alpha-runs"
    assert crow_lib.isolation_account("beta", rows) == "beta"
    assert crow_lib.account_auth("alpha-runs", rows) == "setup-token"
    assert crow_lib.account_auth("nobody", rows) is None


def test_scope_from_extension_and_shebang(tmp_path):
    (tmp_path / "tool").write_text("#!/usr/bin/env python3\n")
    (tmp_path / "run").write_text("#!/bin/bash\n")
    (tmp_path / "notes").write_text("plain\n")
    assert crow_lib.file_kind("tool", tmp_path) == "PY"
    assert crow_lib.file_kind("run", tmp_path) == "SH"
    assert crow_lib.file_kind("notes", tmp_path) == "TXT"
    assert crow_lib.file_kind("a/b.yml") == "YAML"
    assert crow_lib.commit_scope(["a.py", "b.py", "c.md"]) == "PY/MD"


def test_run_commit_message_passes_the_check():
    new = crow_lib.run_commit_message("add the parser", "0123456789ab", {"x.py": "A"})
    chg = crow_lib.run_commit_message("  fix\nthe   flaky test ", "0123456789ab", {"x.py": "M", "y.md": "A"})
    assert new == "new[PY]: add the parser (run 01234567)"
    assert chg == "chg[MD/PY]: fix the flaky test (run 01234567)"
    assert check_msg(new) and check_msg(chg)


def git(cwd, *a):
    return subprocess.run(["git", "-C", str(cwd), *a], check=True, capture_output=True, text=True).stdout.strip()


@pytest.fixture
def repo(tmp_path):
    git(tmp_path, "init", "-q", "-b", "main")
    git(tmp_path, "config", "user.email", "t@example.com")
    git(tmp_path, "config", "user.name", "t")
    (tmp_path / "a.py").write_text("x = 1\n")
    git(tmp_path, "add", ".")
    git(tmp_path, "commit", "-qm", "new[PY]: base")
    return tmp_path


def test_commit_run_diff_commits_what_the_run_left(repo):
    (repo / "a.py").write_text("x = 2\n")
    (repo / "README.md").write_text("doc\n")
    res = crow_lib.commit_run_diff(str(repo), "update a", "deadbeef00")
    assert res["sha"] == git(repo, "rev-parse", "HEAD")
    assert git(repo, "log", "-1", "--format=%s") == "chg[MD/PY]: update a (run deadbeef)"
    assert res["files"] == ["README.md", "a.py"]
    assert crow_lib.commit_run_diff(str(repo), "again", "x") == {"sha": None, "reason": "no changes"}


def test_commit_run_diff_defaults_a_blank_label(repo):
    (repo / "c.py").write_text("z = 1\n")
    res = crow_lib.commit_run_diff(str(repo), " ", "id")
    assert res["sha"] and git(repo, "log", "-1", "--format=%s") == "new[PY]: isolated run (run id)"


def test_commit_run_diff_refuses_a_bad_message_and_unstages(repo, monkeypatch):
    monkeypatch.setattr(crow_lib, "run_commit_message", lambda *a: "update stuff")
    (repo / "b.py").write_text("y = 1\n")
    head = git(repo, "rev-parse", "HEAD")
    res = crow_lib.commit_run_diff(str(repo), "l", "id")
    assert res["sha"] is None and res["reason"].startswith("message refused")
    assert git(repo, "rev-parse", "HEAD") == head and git(repo, "diff", "--cached", "--name-only") == ""


def args(**kw):
    base = dict(continue_run=None, cwd=None, kind="task", account=None, isolate=True)
    return types.SimpleNamespace(**{**base, **kw})


def test_dispatch_refusals(tmp_path, monkeypatch):
    monkeypatch.setattr(crow_lib, "provider_of", lambda acct, status=None: "claude")
    assert "--continue-run" in cr.isolation_refusal(args(continue_run="abc"), "a", "/x", [])
    assert "--cwd" in cr.isolation_refusal(args(), "a", "/x", [])
    (tmp_path / ".git").mkdir()
    assert "worktree marker" in cr.isolation_refusal(args(cwd=str(tmp_path)), "a", str(tmp_path), [])
    wt = tmp_path / "wt"
    wt.mkdir()
    (wt / ".git").write_text("gitdir: /x\n")
    rows = [{"session_name": "a", "auth_type": "full-login"}]
    assert "setup token" in cr.isolation_refusal(args(cwd=str(wt)), "a", str(wt), rows)
    assert "cannot run isolated" in cr.isolation_refusal(args(cwd=str(wt), kind="oneshot"), "a", str(wt),
                                                         [{"session_name": "a", "auth_type": "setup-token"}])


def test_notify_commits_only_writable_isolated_runs(monkeypatch):
    seen = []
    monkeypatch.setattr(crow_lib, "commit_run_diff", lambda cwd, label, rid: seen.append(rid) or {"sha": "s"})
    monkeypatch.setattr(cr, "log_event", lambda *a: None)
    cr.commit_isolated({"run_id": "r1", "cwd": "/x", "label": "l"})
    assert seen == ["r1"]
