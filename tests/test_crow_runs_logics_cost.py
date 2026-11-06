"""crow-runs: Logics Executions records (item_074) and run cost (item_075)."""
import importlib.machinery
import importlib.util
import json
import os

from conftest import BIN

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)

REC = {"run_id": "r1", "project": "hub", "label": "do it", "cwd": "/x", "logics_ref": "item_074"}


def test_start_command():
    cmd = cr.record_run_cmd("item_074", "r1", "running", "do it")
    assert cmd == ["logics-manager", "flow", "record-run", "item_074", "--runner", "cdx", "--run-id", "r1",
                   "--state", "running", "--title", "do it"]
    assert "--progress" not in cmd


def test_final_states_map_and_keep_run_id():
    assert "succeeded" in cr.record_run_cmd("item_074", "r1", "done")
    for st in ("failed", "budget", "stuck", "limit"):
        cmd = cr.record_run_cmd("item_074", "r1", st)
        assert cmd[cmd.index("--state") + 1] == "failed" and cmd[cmd.index("--notes") + 1] == st
        assert cmd[cmd.index("--run-id") + 1] == "r1"
    assert "cancelled" in cr.record_run_cmd("item_074", "r1", "stopped")


def test_record_runs_command_in_repo_and_logs_event(monkeypatch, tmp_path):
    monkeypatch.setattr(cr, "EVENTS", tmp_path / "ev.jsonl")
    monkeypatch.setattr(cr, "REG", {"projects": [{"name": "hub", "scope": "personal"}]})
    monkeypatch.setattr(cr.shutil, "which", lambda n: "/bin/lm")
    seen = {}
    monkeypatch.setattr(cr.subprocess, "run", lambda cmd, **kw: seen.update(cmd=cmd, cwd=kw["cwd"]) or type("R", (), {"returncode": 0})())
    cr.record_run(REC, "running")
    assert seen["cwd"] == "/x" and seen["cmd"][3] == "item_074"
    ev = json.loads((tmp_path / "ev.jsonl").read_text())
    assert set(ev) == {"ts", "mission", "type", "actor", "data"} and ev["type"] == "run.running"


def test_no_ref_does_nothing(monkeypatch):
    monkeypatch.setattr(cr.subprocess, "run", lambda *a, **k: 1 / 0)
    cr.record_run({"run_id": "r1"}, "running")


def test_missing_tool_and_failure_only_warn(monkeypatch, tmp_path, capsys):
    monkeypatch.setattr(cr, "EVENTS", tmp_path / "ev.jsonl")
    monkeypatch.setattr(cr, "REG", {"projects": []})
    monkeypatch.setattr(cr.shutil, "which", lambda n: None)
    cr.record_run(REC, "running")
    assert "warning" in capsys.readouterr().err
    monkeypatch.setattr(cr.shutil, "which", lambda n: "/bin/lm")
    monkeypatch.setattr(cr.subprocess, "run", lambda *a, **k: type("R", (), {"returncode": 1, "stderr": "bad ref", "stdout": ""})())
    cr.record_run(REC, "done")
    assert "bad ref" in capsys.readouterr().err


def test_work_repo_skipped(monkeypatch, tmp_path, capsys):
    monkeypatch.setattr(cr, "EVENTS", tmp_path / "ev.jsonl")
    monkeypatch.setattr(cr, "REG", {"projects": [{"name": "hub", "scope": "work"}]})
    monkeypatch.setattr(cr.subprocess, "run", lambda *a, **k: 1 / 0)
    cr.record_run(REC, "running")
    assert "work repo" in capsys.readouterr().err


def test_run_cost_with_and_without():
    assert cr.run_cost({"cost_usd": 0.7084026}) == 0.7084
    assert cr.run_cost({"cost_usd": None}) is None
    assert cr.run_cost({}) is None
    assert cr.run_cost({"cost_usd": 0}) == 0


def test_fmt_cost():
    assert cr.fmt_cost(0.7084) == "$0.71"
    assert cr.fmt_cost(0.7084, 25.0) == "$0.71 of $25"
    assert cr.fmt_cost(None, 25.0) == ""


def _data(finished):
    return {"agents": [], "slots": {"busy": 0, "max": 2}, "finished": finished}


def _row(rid, cost, budget=None):
    return {"run_id": rid, "project": "p", "label": "l", "status": "done", "duration_s": 60, "ended_at": None,
            "commits": [], "cost_usd": cost, "budget_usd": budget}


def test_render_shows_cost_and_total_only_when_known():
    out = "\n".join(cr.render(_data([_row("aaaaaaaa", 0.5, 25.0), _row("bbbbbbbb", 1.25), _row("cccccccc", None)]), color=False, width=200))
    assert "$0.50 of $25" in out and "$1.25" in out
    assert "background runs cost $1.75" in out
    none = "\n".join(cr.render(_data([_row("dddddddd", None)]), color=False, width=200))
    assert "$" not in none
