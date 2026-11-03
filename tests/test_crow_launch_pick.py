"""Account choice for crow open, term and handoff: live sessions and empty quotas are skipped."""
import json
import os
import sys
from datetime import datetime, timezone

import pytest
import yaml

from conftest import BIN, run

sys.path.insert(0, str(BIN))
import crow_lib  # noqa: E402

NOW = datetime(2026, 11, 3, 16, 0, tzinfo=timezone.utc)
ACCTS = [{"name": "main", "policy": "use"}, {"name": "spare", "policy": "use"},
         {"name": "big", "policy": "use"}, {"name": "shared", "policy": "ask-first"},
         {"name": "mate", "policy": "never"}]


def row(name, left5=80, leftw=80):
    return {"session_name": name, "enabled": True, "auth_status": "authenticated", "remaining_5h_pct": left5,
            "remaining_week_pct": leftw, "reset_5h_at": "Nov 3 18:50"}


def pick(rows, busy=(), accts=ACCTS):
    st = crow_lib.pick_accounts(accts, {r["session_name"]: r for r in rows}, None, NOW)
    return crow_lib.launch_pick(st, set(busy), "main")


def test_default_stays_while_it_has_room():
    assert pick([row("main", 30), row("big", 90)])["account"] == "main"


def test_default_at_zero_moves_to_the_account_with_the_most_room():
    p = pick([row("main", 0), row("spare", 40), row("big", 70), row("shared", 100)])
    assert (p["account"], p["ask"], p["room"]) == ("big", False, 70)
    assert {"name": "main", "why": "5h limit (resets Nov 3 18:50)"} in p["skipped"]


def test_room_is_the_smaller_window():
    assert pick([row("main", 0), row("spare", 40, 90), row("big", 90, 10)])["account"] == "spare"


def test_default_nearly_out_moves_too():
    assert pick([row("main", 3), row("spare", 40)])["account"] == "spare"


def test_busy_default_moves_to_the_next_free_account():
    p = pick([row("main"), row("spare", 40), row("big", 20)], busy={"main"})
    assert p["account"] == "spare" and {"name": "main", "why": "live session"} in p["skipped"]


def test_only_ask_first_left_asks():
    p = pick([row("main", 0), row("spare", 50), row("big", 0, 50), row("shared", 60)], busy={"spare"})
    assert (p["account"], p["ask"]) == ("shared", True)


def test_none_free_returns_no_account():
    p = pick([row("main", 0), row("spare"), row("big", 50, 0), row("shared", 0)], busy={"spare"})
    assert p["account"] is None
    assert [s["name"] for s in p["skipped"]] == ["main", "spare", "big", "shared"]


# ---------- the CLI against a fake cdx, setsid and kitty ----------

FAKE_CDX = r"""#!/usr/bin/env python3
import json, os, sys
argv = sys.argv[1:]
with open(os.environ["FAKE_LOG"], "a") as f:
    f.write(json.dumps({"argv": argv, "account": os.environ.get("CROW_ACCOUNT")}) + "\n")
state = json.load(open(os.environ["FAKE_STATE"]))
if argv[:1] == ["status"]:
    print(json.dumps({"rows": state["rows"]}))
elif argv[:1] == ["sessions"]:
    print(json.dumps({"sessions": [{"session": s, "active": True} for s in state["busy"]]}))
elif argv[:1] == ["configs"]:
    print(json.dumps({"sessions": []}))
"""
# setsid -f kitty --title T <cmd...>: run <cmd...> in the foreground so the test sees the launch.
FAKE_SETSID = """#!/bin/sh
shift; shift; shift; shift
exec "$@"
"""


@pytest.fixture
def fake(hub, tmp_path):
    home, _ = hub
    (home / "local/accounts.yaml").write_text(yaml.safe_dump({"default": "main", "accounts": [
        {"name": "main", "policy": "use"}, {"name": "spare", "policy": "use"},
        {"name": "shared", "policy": "ask-first"}]}))
    repo = tmp_path / "repo"
    repo.mkdir()
    (home / "local/projects.yaml").write_text(yaml.safe_dump({"projects": [
        {"name": "demo", "scope": "personal", "path": str(repo), "forge": "none"}]}))
    d = tmp_path / "fakebin"
    d.mkdir()
    for name, body in (("cdx", FAKE_CDX), ("setsid", FAKE_SETSID)):
        (d / name).write_text(body)
        (d / name).chmod(0o755)
    log, state = tmp_path / "cdx.log", tmp_path / "state.json"
    env = {"PATH": f"{d}:{os.environ['PATH']}", "FAKE_LOG": str(log), "FAKE_STATE": str(state),
           "CDX_SESSION_NAME": "", "CROW_ACCOUNT": ""}

    def set_state(rows, busy=()):
        state.write_text(json.dumps({"rows": rows, "busy": list(busy)}))

    def launches():
        calls = [json.loads(l) for l in log.read_text().splitlines()] if log.exists() else []
        return [c for c in calls if c["argv"][:1] not in (["status"], ["sessions"], ["configs"])]
    return env, set_state, launches


def test_term_with_the_default_at_zero_launches_on_the_free_account(hub, fake):
    env, set_state, launches = fake
    set_state([row("main", 0), row("spare", 60), row("shared", 100)])
    r = run(hub, "crow", "term", "demo", "hello", env=env)
    assert r.returncode == 0, r.stderr
    assert "using 'spare'" in r.stderr and "main: 5h limit" in r.stderr
    (call,) = launches()
    assert call["argv"][0] == "spare" and call["argv"][-2:] == ["--prompt", "hello"]
    assert "--on-busy" in call["argv"]


def test_open_refuses_when_no_account_has_room(hub, fake):
    env, set_state, launches = fake
    set_state([row("main", 0), row("spare", 60, 0), row("shared", 0)])
    r = run(hub, "crow", "open", "--on-busy", "fail", "demo", env=env)
    assert r.returncode == 1 and "no allowed account is free" in r.stderr
    assert launches() == []


def test_handoff_refuses_on_ask_first_only_and_writes_no_note(hub, fake):
    env, set_state, launches = fake
    set_state([row("main", 0), row("spare", 60), row("shared", 80)], busy={"spare"})
    r = run(hub, "crow", "handoff", "demo", "--task", "fix it", "--end-state", "show the diff", env=env)
    assert r.returncode == 1 and "CROW_ACCOUNT=shared" in r.stderr
    assert launches() == [] and not (hub[0] / "local/state/handoffs").exists()


def test_handoff_passes_task_and_end_state_to_cdx(hub, fake):
    env, set_state, launches = fake
    set_state([row("main"), row("spare")], busy={"main"})
    r = run(hub, "crow", "handoff", "demo", "--task", "fix it", "--end-state", "show the diff", env=env)
    assert r.returncode == 0, r.stderr
    (call,) = launches()
    argv = call["argv"]
    assert argv[0] == "spare"
    task = argv[argv.index("--task") + 1]
    assert task.startswith("fix it\n\nHand-off from the crow hub") and "End state:" not in task
    assert argv[argv.index("--end-state") + 1] == "show the diff"


def test_an_explicit_account_wins(hub, fake):
    env, set_state, launches = fake
    set_state([row("main", 0), row("spare", 60)])
    r = run(hub, "crow", "open", "demo", env={**env, "CROW_ACCOUNT": "main"})
    assert r.returncode == 0, r.stderr
    assert launches()[0]["argv"][0] == "main"


def test_the_inherited_account_of_the_calling_session_is_not_a_choice(hub, fake):
    env, set_state, launches = fake
    set_state([row("main"), row("spare", 60)], busy={"main"})
    r = run(hub, "crow", "term", "demo", env={**env, "CROW_ACCOUNT": "main", "CDX_SESSION_NAME": "main"})
    assert r.returncode == 0, r.stderr
    assert launches()[0]["argv"][0] == "spare"
