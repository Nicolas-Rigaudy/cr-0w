"""Headless run flags (req_017 levers 3, 4, 7, 16) and the budget / stuck run states."""
import importlib.machinery
import importlib.util
import os
import subprocess
import sys

from conftest import BIN

sys.path.insert(0, str(BIN))
import crow_lib  # noqa: E402

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs_flags", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs_flags", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)


def opt(args, flag):
    return args[args.index(flag) + 1]


def test_every_kind_pins_model_id_effort_budget_timeout():
    for kind in crow_lib.RUN_KINDS:
        args, env = crow_lib.run_flags(kind, "claude", env={})
        assert opt(args, "--model").startswith("claude-")
        assert opt(args, "--power")
        assert float(opt(args, "--budget")) > 0
        assert int(opt(args, "--timeout-seconds")) > 0
        assert env == crow_lib.HEADLESS_ENV


def test_alias_maps_to_exact_id():
    args, _ = crow_lib.run_flags("task", "claude", model="sonnet", env={})
    assert opt(args, "--model") == "claude-sonnet-5-5"


def test_continue_run_passes_no_model_or_effort():
    args, _ = crow_lib.run_flags("review", "claude", model="opus", power="high", continue_run="r1", env={})
    assert "--model" not in args and "--power" not in args and "--budget" in args
    assert "--failover" not in args  # cdx refuses --failover with --continue-run


def test_env_overrides_budget_and_timeout():
    args, _ = crow_lib.run_flags("retro", "claude", env={"CROW_RETRO_BUDGET": "0.5", "CROW_RETRO_TIMEOUT_MIN": "3"})
    assert opt(args, "--budget") == "0.5" and opt(args, "--timeout-seconds") == "180"


def test_codex_gets_no_budget_nor_claude_env():
    args, env = crow_lib.run_flags("task", "codex", model="gpt-5.5", env={})
    assert "--budget" not in args and opt(args, "--model") == "gpt-5.5" and env == {}


def test_states_budget_and_stuck(monkeypatch, tmp_path):
    monkeypatch.setattr(cr, "STUCK", tmp_path / "stuck.json")
    assert cr.run_state({"run_id": "a", "exit_code": 125, "error": {"code": "budget_exceeded"}}) == "budget"
    assert cr.run_state({"run_id": "b", "exit_code": 1, "status": "timed_out"}) == "stuck"


def test_watchdog_kills_idle_run_and_spares_young_one(monkeypatch, tmp_path):
    monkeypatch.setattr(cr, "STUCK", tmp_path / "stuck.json")
    monkeypatch.setattr(cr, "run_transcript", lambda row: None)
    monkeypatch.setattr(cr, "agent_procs", lambda: [])
    p = subprocess.Popen(["sleep", "60"])
    from datetime import datetime, timezone
    young = {"run_id": "y", "started_at": datetime.now(timezone.utc).isoformat(), "pid": p.pid}
    assert not cr.watchdog(young)
    old = {"run_id": "o", "started_at": "2026-11-03T10:00:00+00:00", "pid": p.pid}
    assert cr.watchdog(old)
    assert p.wait(timeout=5) < 0
    assert cr.run_state({"run_id": "o", "exit_code": 1}) == "stuck"


def test_every_run_carries_failover_and_retro_kind():
    for kind in crow_lib.RUN_KINDS:
        for provider in ("claude", "codex"):
            args, _ = crow_lib.run_flags(kind, provider, model="gpt-5.5" if provider == "codex" else None, env={})
            assert "--failover" in args
            assert ("--kind" in args) == (kind == "retro")
    assert opt(crow_lib.run_flags("retro", "claude", env={})[0], "--kind") == "retro"


def test_per_run_budget_overrides_kind_and_env():
    args, _ = crow_lib.run_flags("task", "claude", env={"CROW_TASK_BUDGET": "9"}, budget=1.5)
    assert opt(args, "--budget") == "1.5"


def test_stop_kills_one_running_run_and_refuses_ambiguity(monkeypatch, capsys):
    import pytest
    monkeypatch.setattr(cr, "agent_procs", lambda: [])
    p = subprocess.Popen(["sleep", "60"])
    monkeypatch.setattr(cr, "running_runs", lambda: [{"run_id": "abc123", "pid": p.pid}, {"run_id": "abd456"}])
    with pytest.raises(SystemExit):
        cr.stop("ab")
    cr.stop("abc")
    assert p.wait(timeout=5) < 0
    assert f"stopped abc123 (pid {p.pid})" in capsys.readouterr().out
    with pytest.raises(SystemExit):
        cr.stop("abd")


def test_stop_refuses_an_empty_id_and_records_the_stopped_run(monkeypatch, tmp_path):
    import pytest
    monkeypatch.setattr(cr, "STOPPED", tmp_path / "stopped.json")
    monkeypatch.setattr(cr, "agent_procs", lambda: [])
    p = subprocess.Popen(["sleep", "60"])
    monkeypatch.setattr(cr, "running_runs", lambda: [{"run_id": "abc123", "pid": p.pid}])
    with pytest.raises(SystemExit):
        cr.stop("")
    assert p.poll() is None
    cr.stop("abc")
    p.wait(timeout=5)
    assert cr.run_state({"run_id": "abc123", "exit_code": -15}) == "stopped"
