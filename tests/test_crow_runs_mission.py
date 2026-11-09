"""Mission runs (item_133): crow-core tick dispatches through crow-runs and reads the end with `state`."""
import importlib.machinery
import importlib.util
import json
import os
import sys
import types

import pytest
from conftest import BIN

sys.path.insert(0, str(BIN))
import crow_lib

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs_mission", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs_mission", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)


@pytest.fixture
def fake(tmp_path, monkeypatch):
    """A fake cdx and accounts check: records every command, answers `cdx run` with `reply`."""
    calls, popen = [], []
    reply = {"out": json.dumps({"ok": True, "run_id": "abcdef0123456789"})}

    def run(cmd, **kw):
        calls.append({"cmd": list(cmd), "env": kw.get("env")})
        out = reply["out"] if cmd[:2] == ["cdx", "run"] else ""
        return types.SimpleNamespace(returncode=0, stdout=out, stderr="")
    monkeypatch.setattr(cr.subprocess, "run", run)
    monkeypatch.setattr(cr.subprocess, "Popen", lambda *a, **k: popen.append(a))
    monkeypatch.setattr(cr, "STATE", tmp_path / "state")
    monkeypatch.setattr(cr, "LEDGER", tmp_path / "state" / "runs.jsonl")
    monkeypatch.setattr(cr, "REG", {"projects": [{"name": "demo", "path": str(tmp_path)}]})
    monkeypatch.setattr(cr, "running_runs", list)
    monkeypatch.setattr(cr, "pick_account", lambda a: "acct")
    monkeypatch.setattr(crow_lib, "provider_of", lambda acct, status=None: "claude")
    return types.SimpleNamespace(calls=calls, popen=popen, reply=reply, tmp=tmp_path)


def args(**kw):
    base = {"project": "demo", "label": "research-x", "prompt_file": "/m/brief.md", "prompt": None,
            "kind": "review", "model": "opus", "power": "medium", "continue_run": None, "account": None,
            "isolate": False, "cwd": None, "permission": "review", "force": False, "budget": 5.0, "dry_run": False,
            "context": None, "logics": None, "resumes": None, "mission": "20261109-research-x", "timeout_min": 60,
            "json": True, "write_dir": None}
    return types.SimpleNamespace(**{**base, **kw})


def test_mission_dispatch_prints_json_records_the_mission_and_skips_notify(fake, capsys):
    cr.dispatch(args())
    assert json.loads(capsys.readouterr().out) == {"ok": True, "run_id": "abcdef0123456789", "account": "acct"}
    cdx = next(c for c in fake.calls if c["cmd"][:2] == ["cdx", "run"])["cmd"]
    assert cdx[2] == "acct" and "--prompt-file" in cdx and "/m/brief.md" in cdx
    assert cdx[cdx.index("--timeout-seconds") + 1] == "3600"
    assert cdx[cdx.index("--budget") + 1] == "5" and cdx[cdx.index("--permission") + 1] == "review"
    assert "--extra-args" in cdx  # the kind's sandbox allowlist rides along
    rec = json.loads(cr.LEDGER.read_text())
    assert rec["mission"] == "20261109-research-x" and rec["run_id"] == "abcdef0123456789"
    assert fake.popen == []  # no notify: the tick finishes mission runs


def test_mission_dispatch_failure_is_a_coded_json_error(fake, capsys):
    fake.reply["out"] = json.dumps({"ok": False, "error": {"code": "too_many_runs", "message": "2 running"}})
    with pytest.raises(SystemExit) as e:
        cr.dispatch(args())
    assert e.value.code == 1
    out = json.loads(capsys.readouterr().out)
    assert out["ok"] is False and out["error"]["code"] == "too_many_runs"
    assert not cr.LEDGER.exists()


def test_plain_dispatch_keeps_notify_and_text(fake, capsys):
    cr.dispatch(args(mission=None, json=False, timeout_min=None, kind="task", permission=None))
    out = capsys.readouterr().out
    assert out.startswith("dispatched abcdef01") and len(fake.popen) == 1
    cdx = next(c for c in fake.calls if c["cmd"][:2] == ["cdx", "run"])["cmd"]
    assert cdx[cdx.index("--timeout-seconds") + 1] == str(crow_lib.RUN_KINDS["task"]["timeout_min"] * 60)


def test_state_of_a_live_a_done_and_a_limit_run(fake, monkeypatch):
    rows = {"aa11": {"run_id": "aa11", "ended_at": None},
            "bb22": {"run_id": "bb22", "ended_at": "t", "exit_code": 0},
            "cc33": {"run_id": "cc33", "ended_at": "t", "exit_code": 1}}
    monkeypatch.setattr(cr, "finished_rows", lambda since, ttl=8: rows)
    monkeypatch.setattr(cr, "run_state", lambda r: "done" if r["exit_code"] == 0 else "limit")
    monkeypatch.setattr(cr, "ledger", lambda: [{"run_id": "cc33", "account": "acct", "context": "day"}])
    monkeypatch.setattr(crow_lib, "account_state", lambda acct: {
        "accounts": [{"name": "acct", "reset_5h_at_iso": None}], "current_blocked": "5h", "auto_target": "spare"})
    monkeypatch.setattr(crow_lib, "minutes_to", lambda iso, now: 200)
    assert cr.run_end_state("aa11") == {"run_id": "aa11", "ended": False}
    done = cr.run_end_state("bb22")
    assert done["ended"] and done["state"] == "done" and done["exit_code"] == 0 and "policy" not in done
    lim = cr.run_end_state("cc")
    assert lim["state"] == "limit" and lim["policy"]["choice"] == "move" and lim["policy"]["target"] == "spare"


def test_write_dir_leaves_plan_mode_and_allows_only_the_mission_folder(fake, capsys):
    """item_143 AC1: a read-only mission run edits its folder only; the review allowlist keeps the project read only."""
    import shlex
    from pathlib import Path
    mdir = fake.tmp / "missions" / "20261109-research-x"
    mdir.mkdir(parents=True)
    cr.dispatch(args(write_dir=str(mdir)))
    cdx = next(c for c in fake.calls if c["cmd"][:2] == ["cdx", "run"])["cmd"]
    assert cdx[cdx.index("--permission") + 1] == "default"
    extra = shlex.split(cdx[cdx.index("--extra-args") + 1])
    assert len(cdx[cdx.index("--extra-args") + 1]) <= crow_lib.CDX_EXTRA_ARGS_MAX
    assert extra[extra.index("--add-dir") + 1] == os.path.realpath(mdir)
    allow = json.loads(Path(extra[extra.index("--settings") + 1]).read_text())["permissions"]["allow"]
    d = "/" + os.path.realpath(mdir)
    assert f"Edit({d}/**)" in allow and f"Write({d}/**)" in allow
    assert "Edit" not in allow and "Write" not in allow  # nothing else is writable
    assert json.loads(cr.LEDGER.read_text())["write_dir"] == os.path.realpath(mdir)


def test_write_dir_must_exist(fake, capsys):
    with pytest.raises(SystemExit):
        cr.dispatch(args(write_dir=str(fake.tmp / "nope")))
    assert json.loads(capsys.readouterr().out)["error"]["code"] == "write_dir_missing"


def test_write_dir_refused_on_codex(fake, capsys, monkeypatch, tmp_path):
    monkeypatch.setattr(crow_lib, "provider_of", lambda acct, status=None: "codex")
    with pytest.raises(SystemExit):
        cr.dispatch(args(write_dir=str(tmp_path)))
    assert json.loads(capsys.readouterr().out)["error"]["code"] == "write_dir_unsupported"


def test_review_kind_on_claude_shares_a_busy_cwd(fake, capsys, monkeypatch):
    monkeypatch.setattr(cr, "running_runs", lambda: [{"cwd": str(fake.tmp)}])
    cr.dispatch(args(permission=None, write_dir=str(fake.tmp)))
    assert json.loads(capsys.readouterr().out)["ok"] is True
