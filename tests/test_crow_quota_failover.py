"""Quota failover: account picker rules, the warning hook, dispatch account choice and limit requeue."""
import importlib.machinery
import importlib.util
import json
import os
import sys
from datetime import datetime, timedelta, timezone

import pytest

from conftest import BIN

sys.path.insert(0, str(BIN))
import crow_lib  # noqa: E402

NOW = datetime(2026, 11, 3, 16, 0, tzinfo=timezone.utc)
USE = {"name": "main", "policy": "use"}
SHARED = {"name": "shared", "policy": "ask-first", "fallback": "auto", "max_5h_pct": 50, "max_week_pct": 50,
          "free_5h_min": 45}
PLAIN = {"name": "mate", "policy": "ask-first"}


def row(name, used5=10, usedw=10, reset_in_min=200):
    return {"session_name": name, "enabled": True, "auth_status": "authenticated",
            "remaining_5h_pct": 100 - used5, "remaining_week_pct": 100 - usedw, "reset_5h_at": "later",
            "reset_5h_at_iso": (NOW + timedelta(minutes=reset_in_min)).isoformat(), "reset_week_at": "next week"}


def pick(rows, accts=(USE, SHARED, PLAIN), current="main"):
    return crow_lib.pick_accounts(list(accts), {r["session_name"]: r for r in rows}, current, NOW)


def test_use_account_free_needs_no_fallback():
    st = pick([row("main", 40), row("shared"), row("mate")])
    assert st["current_blocked"] is None
    assert not next(a for a in st["accounts"] if a["name"] == "shared")["auto_ok"]


def test_auto_fallback_once_use_account_is_out():
    st = pick([row("main", 100), row("shared", 20), row("mate", 0)])
    assert st["current_blocked"].startswith("5h limit")
    assert st["auto_target"] == "shared"
    assert st["recommended"] == "shared"


def test_nearly_out_counts_as_blocked_for_new_work():
    st = pick([row("main", 96), row("shared", 20), row("mate")])
    assert st["current_blocked"] == "5h 96% used" and st["auto_target"] == "shared"


def test_shared_caps_block_auto_fallback():
    st = pick([row("main", 100), row("shared", 60), row("mate", 0)])
    assert "cap is 50%" in next(a for a in st["accounts"] if a["name"] == "shared")["blocked"]
    assert st["auto_target"] is None and st["recommended"] == "mate"
    st = pick([row("main", 100), row("shared", 20, usedw=55), row("mate")])
    assert st["auto_target"] is None


def test_shared_5h_cap_lifts_near_its_reset():
    st = pick([row("main", 100), row("shared", 70, reset_in_min=30), row("mate")])
    assert st["auto_target"] == "shared"


def test_plain_ask_first_is_never_auto():
    st = pick([row("main", 100), row("mate", 0)], accts=(USE, PLAIN))
    assert st["auto_target"] is None and st["recommended"] == "mate"


# ---------- crow-relaunch hook ----------

def load(name, path):
    loader = importlib.machinery.SourceFileLoader(name, str(path))
    spec = importlib.util.spec_from_loader(name, loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


@pytest.fixture
def rl(tmp_path, monkeypatch):
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    mod = load("crow_relaunch", BIN / "crow-relaunch")
    monkeypatch.setattr(mod, "WARNED", str(tmp_path / "warned.json"))
    return mod


def test_hook_warns_once_per_level(rl, monkeypatch, capsys):
    state = {"rows": [row("main", 82), row("shared", 20)]}
    monkeypatch.setattr(crow_lib, "account_state", lambda **k: pick(state["rows"]))
    rl.hook([])
    out = capsys.readouterr().out
    assert "QUOTA 5h" in out and "82%" in out and "/wrap" not in out
    rl.hook([])
    assert capsys.readouterr().out == ""
    state["rows"] = [row("main", 96), row("shared", 20)]
    rl.hook([])
    out = capsys.readouterr().out
    assert "/wrap" in out and "Fallback: shared" in out and "no OK needed" in out
    rl.hook([])
    assert capsys.readouterr().out == ""


def test_hook_quiet_below_80_and_never_fails(rl, monkeypatch, capsys):
    monkeypatch.setattr(crow_lib, "account_state", lambda **k: pick([row("main", 50)]))
    assert rl.hook([]) == 0 and capsys.readouterr().out == ""

    def boom(**k):
        raise RuntimeError("cdx down")
    monkeypatch.setattr(crow_lib, "account_state", boom)
    assert rl.hook([]) == 0


# ---------- crow-runs dispatch and requeue ----------

@pytest.fixture
def cr(tmp_path, monkeypatch):
    (tmp_path / "local").mkdir()
    (tmp_path / "local" / "projects.yaml").write_text(
        "projects:\n  - {name: p, path: %s}\n" % (tmp_path / "repo"))
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    monkeypatch.setenv("CROW_ACCOUNT", "main")
    return load("crow_runs", BIN / "crow-runs")


class A:
    account = None
    continue_run = None


def test_dispatch_stays_on_default_when_free(cr, monkeypatch):
    monkeypatch.setattr(crow_lib, "account_state", lambda cur: pick([row("main", 40), row("shared")]))
    assert cr.pick_account(A()) == "main"


def test_dispatch_moves_to_auto_fallback(cr, monkeypatch):
    monkeypatch.setattr(crow_lib, "account_state", lambda cur: pick([row("main", 100), row("shared", 20)]))
    assert cr.pick_account(A()) == "shared"


def test_dispatch_refuses_ask_first_without_ok(cr, monkeypatch):
    monkeypatch.setattr(crow_lib, "account_state",
                        lambda cur: pick([row("main", 100), row("mate")], accts=(USE, PLAIN)))
    with pytest.raises(SystemExit) as e:
        cr.pick_account(A())
    assert "--account mate" in str(e.value)
    a = A(); a.account = "mate"
    assert cr.pick_account(a) == "mate"


def test_continue_run_stays_on_parent_account(cr, monkeypatch):
    monkeypatch.setattr(cr, "ledger", lambda: [{"run_id": "abcd1234", "account": "shared"}])
    a = A(); a.continue_run = "abcd"
    assert cr.pick_account(a) == "shared"


def requeue_env(cr, monkeypatch, tmp_path, rows, ledger_extra=()):
    prompt = tmp_path / "task.md"
    prompt.write_text("do the thing")
    rec = {"run_id": "r1", "project": "p", "label": "x", "cwd": str(tmp_path / "repo"), "account": "main",
           "prompt_file": str(prompt), "power": "medium"}
    monkeypatch.setattr(cr, "ledger", lambda: [rec, *ledger_extra])
    monkeypatch.setattr(cr, "STATE", tmp_path / "state")
    monkeypatch.setattr(cr, "EVENTS", tmp_path / "state" / "runs-events.jsonl")
    monkeypatch.setattr(crow_lib, "account_state", lambda cur: pick(rows))
    calls, sleeps = [], []

    class R:
        returncode = 0
        stdout = stderr = ""
    monkeypatch.setattr(cr.subprocess, "run", lambda cmd, **k: calls.append(cmd) or R())
    monkeypatch.setattr(cr.time, "sleep", sleeps.append)
    monkeypatch.setattr(cr, "datetime", type("D", (), {"now": staticmethod(lambda tz=None: NOW)}))
    return rec, calls, sleeps


def test_requeue_moves_to_auto_fallback_with_resume_note(cr, monkeypatch, tmp_path):
    rec, calls, sleeps = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100), row("shared", 20)])
    assert cr.requeue(rec, deadline=1e12) == "resumed on shared"
    cmd = calls[0]
    assert cmd[cmd.index("--account") + 1] == "shared" and cmd[cmd.index("--resumes") + 1] == "r1"
    assert "RESUME NOTE" in open(cmd[cmd.index("--prompt-file") + 1]).read()
    assert sleeps == []


def test_requeue_waits_for_reset_then_continues(cr, monkeypatch, tmp_path):
    rec, calls, sleeps = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100, reset_in_min=30), row("shared", 80)])
    assert cr.requeue(rec, deadline=1e12) == "continued on main after the reset"
    assert sleeps and 30 * 60 < sleeps[0] <= 30 * 60 + cr.RESET_GRACE_S
    cmd = calls[0]
    assert cmd[cmd.index("--continue-run") + 1] == "r1" and cmd[cmd.index("--account") + 1] == "main"


def test_requeue_stops_after_max_chain(cr, monkeypatch, tmp_path):
    chain = [{"run_id": f"c{i}", "resumes": f"c{i - 1}" if i else "r1"} for i in range(3)]
    rec, calls, _ = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100), row("shared", 20)], chain)
    assert "not requeued" in cr.requeue(chain[-1], deadline=1e12) and calls == []


def test_run_state_reads_limit_from_tail(cr, tmp_path):
    out = tmp_path / "out"
    out.write_text("You've hit your session limit · resets 6:50pm")
    assert cr.run_state({"exit_code": 1, "artifacts": {"stdout_path": str(out)}}) == "limit"
    assert cr.run_state({"exit_code": 1, "artifacts": {}}) == "failed"
    assert cr.run_state({"exit_code": 0}) == "done"


# ---------- crow rescue ----------

def test_rescue_lists_every_stopped_and_open_session(rl, monkeypatch, capsys):
    now = datetime.now(timezone.utc)
    ago = lambda h: (now - timedelta(hours=h)).isoformat().replace("+00:00", "Z")
    reg = {"sessions": [
        {"conversation_id": f"cut{i}", "account": "main", "cwd": "/w", "last_event": ago(i), "message": "limit",
         "relaunched": None, "first_prompt": f"task {i}"} for i in range(1, 4)] + [
        {"conversation_id": "stale9", "account": "main", "cwd": "/w", "last_event": ago(9), "relaunched": None},
        {"conversation_id": "moved", "account": "main", "cwd": "/w", "last_event": ago(1), "relaunched": {"target": "shared"}},
        {"conversation_id": "other", "account": "shared", "cwd": "/w", "last_event": ago(1), "relaunched": None}]}
    live = {"agents": [{"kind": "session", "account": "main", "session_id": "open1", "cwd": "/w", "state": "idle", "pid": 7},
                       {"kind": "run", "account": "main", "session_id": "run1"}]}
    monkeypatch.setattr(rl, "helper", lambda *a: json.dumps(reg if a[0] == "crow-sessions" else live))
    monkeypatch.setattr(rl, "accounts", lambda: {"main": USE, "shared": SHARED})
    monkeypatch.setattr(rl, "account_status", lambda cur: pick([row("main", 100), row("shared", 20)]))
    assert rl.rescue(["main", "--dry-run"]) == 0
    out = capsys.readouterr().out
    assert "4 session(s) on main -> shared" in out
    for sid in ("cut1", "cut2", "cut3", "open1"):
        assert sid in out
    for sid in ("stale9", "moved", "other", "run1"):
        assert sid not in out


# ---------- usage limit policy by context (item_115) ----------

def events(cr):
    f = cr.EVENTS
    return [json.loads(l) for l in f.read_text().splitlines()] if f.exists() else []


def test_limit_policy_rules():
    assert crow_lib.limit_policy("afk", 240, True)[0] == "wait"
    assert crow_lib.limit_policy("day", 180, True)[0] == "move"
    assert crow_lib.limit_policy("day", 20, True)[0] == "wait"
    assert crow_lib.limit_policy("day", 180, False)[0] == "wait"
    assert crow_lib.limit_policy("afk", None, True)[0] == "move"


def test_afk_run_waits_on_same_account_despite_free_target(cr, monkeypatch, tmp_path):
    rec, calls, sleeps = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100, reset_in_min=240), row("shared", 20)])
    rec["context"] = "afk"
    assert cr.requeue(rec, deadline=1e12) == "continued on main after the reset"
    assert sleeps and sleeps[0] > 240 * 60
    cmd = calls[0]
    assert cmd[cmd.index("--account") + 1] == "main" and cmd[cmd.index("--context") + 1] == "afk"
    ev = events(cr)[-1]
    assert ev["type"] == "run.limit_policy" and ev["data"]["choice"] == "wait"
    assert ev["data"]["reset_5h_at"] == (NOW + timedelta(minutes=240)).isoformat()


def test_day_run_with_far_reset_moves_cold(cr, monkeypatch, tmp_path):
    rec, calls, sleeps = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100, reset_in_min=180), row("shared", 20)])
    rec["context"] = "day"
    assert cr.requeue(rec, deadline=1e12) == "resumed on shared"
    assert sleeps == []
    ev = events(cr)[-1]["data"]
    assert ev["choice"] == "move" and ev["reset_in_min"] == 180


def test_day_run_with_near_reset_waits(cr, monkeypatch, tmp_path):
    rec, calls, sleeps = requeue_env(cr, monkeypatch, tmp_path, [row("main", 100, reset_in_min=20), row("shared", 20)])
    rec["context"] = "day"
    assert cr.requeue(rec, deadline=1e12) == "continued on main after the reset"
    assert events(cr)[-1]["data"]["choice"] == "wait"
