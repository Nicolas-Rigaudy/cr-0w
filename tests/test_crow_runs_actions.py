"""crow-runs summarize_tool: a tool call's own description wins over its raw command."""
import importlib.machinery
import importlib.util

import pytest

from conftest import BIN


@pytest.fixture
def cr(tmp_path, monkeypatch):
    (tmp_path / "local").mkdir()
    (tmp_path / "local" / "projects.yaml").write_text("projects: {}\n")
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    monkeypatch.setenv("CROW_ACCOUNT", "main")
    loader = importlib.machinery.SourceFileLoader("crow_runs", str(BIN / "crow-runs"))
    spec = importlib.util.spec_from_loader("crow_runs", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


def test_bash_description_wins(cr):
    inp = {"command": "cd /home/me/Projects/hub && git status", "description": "Show working tree status"}
    assert cr.summarize_tool("Bash", inp) == "Bash: Show working tree status"


def test_bare_command_is_shortened(cr):
    inp = {"command": "cd /home/me/Projects/hub/local; python3 /home/me/Projects/hub/bin/crow-usage --json\nmore"}
    assert cr.summarize_tool("Bash", inp) == "Bash: python3 …/crow-usage --json"


def test_other_tools_unchanged(cr):
    assert cr.summarize_tool("Edit", {"file_path": "/a/b.go"}) == "Edit: /a/b.go"
    assert cr.summarize_tool("Agent", {"description": "Find callers", "prompt": "long"}) == "Agent: Find callers"
    assert cr.summarize_tool("TodoWrite", {}) == "TodoWrite"


def _proc(root, pid, argv, exe, env="HOME=/h\0CDX_SESSION_NAME=acct\0"):
    d = root / str(pid)
    d.mkdir()
    (d / "cmdline").write_bytes("\0".join(argv).encode() + b"\0")
    (d / "environ").write_bytes(env.encode())
    (d / "exe").symlink_to(exe)
    (d / "cwd").symlink_to("/tmp")


def test_agent_procs_match_claude_and_codex(cr, tmp_path):
    _proc(tmp_path, 11, ["claude", "-p", "x"], "/usr/bin/claude")
    _proc(tmp_path, 12, ["/opt/codex", "exec", "x"], "/opt/codex")
    _proc(tmp_path, 13, ["python3", "x"], "/usr/bin/python3")
    got = {p["pid"]: p for p in cr.agent_procs(str(tmp_path))}
    assert {k: v["provider"] for k, v in got.items()} == {11: "claude", 12: "codex"}
    assert got[12]["account"] == "acct"


def test_is_headless_covers_p_print_and_exec(cr):
    assert cr.is_headless(["claude", "-p", "x"])
    assert cr.is_headless(["claude", "--print"])
    assert cr.is_headless(["codex", "exec", "x"])
    assert not cr.is_headless(["claude", "--resume", "id"])


def test_running_runs_are_the_ones_without_ended_at(cr, monkeypatch):
    rows = {"a": {"run_id": "a", "ended_at": None, "cwd": "/x"}, "b": {"run_id": "b", "ended_at": "t"}}
    monkeypatch.setattr(cr, "finished_rows", lambda since, ttl=8: rows)
    assert [r["run_id"] for r in cr.running_runs()] == ["a"]
    assert list(cr.finished("7d")) == ["b"]
