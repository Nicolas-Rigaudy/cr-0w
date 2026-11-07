"""crow-checks: each rule over a fixture snapshot, the 2026-11-03 replay, dismiss and snooze.

The replay fixture holds checks.yaml and state/ outside a local/ folder (gitignored); the
tests copy them into a throwaway hub's local/.
"""
import datetime as dt
import importlib.machinery
import importlib.util
import json
import shutil
import sys

from conftest import BIN, FIXTURES, run

sys.path.insert(0, str(BIN))
import crow_checks as cc

DAY = FIXTURES / "checks" / "2026-11-03"
NOON = "2026-11-03T12:00:00+01:00"


def snap(now=NOON, **kw):
    t = dt.datetime.fromisoformat(now)
    s = {"now": t, "today": t.date(), "sprint": None, "events": [], "people": [], "mail": [], "rituals": [], "tasks": []}
    s.update(kw)
    return s


def ev(start, subject):
    return {"start": start, "end": start, "subject": subject, "all_day": False}


# sprint

def test_sprint_ending_without_planning_proposes_planning():
    s = snap(sprint={"name": "Sprint 33", "end": "2026-11-03"}, events=[ev("2026-11-03T14:00", "Sprint review")])
    rows = cc.rule_sprint(s, {"sprint": {"rituals": {"planning": ["planning"], "review": ["review"]}}})
    assert [(r["check"], r["action"], r["why"]) for r in rows] == [("sprint", "propose planning", "no planning in the calendar")]


def test_sprint_all_rituals_booked_or_far_gives_nothing():
    cfg = {"sprint": {"rituals": {"planning": ["planning"]}}}
    booked = snap(sprint={"name": "S", "end": "2026-11-04"}, events=[ev("2026-11-05T10:00", "Sprint Planning")])
    assert cc.rule_sprint(booked, cfg) == []
    far = snap(sprint={"name": "S", "end": "2026-11-11"})
    assert cc.rule_sprint(far, cfg) == []
    assert cc.rule_sprint(snap(sprint={"name": "S", "end": "2026-11-04"}, events=None), cfg) == []


def test_sprint_cadence_finds_the_current_end():
    cfg = {"sprint": {"anchor_end": "2026-01-14", "length_days": 14}}
    assert cc.sprint_of(cfg, "/nonexistent", dt.date(2026, 11, 3))["end"] == "2026-11-04"
    assert cc.sprint_of(cfg, "/nonexistent", dt.date(2026, 11, 5))["end"] == "2026-11-18"
    assert cc.sprint_of({}, "/nonexistent", dt.date(2026, 11, 5)) is None


# 1:1

def test_one_on_one_past_due_and_unbooked_yields_a_row():
    s = snap(people=[{"name": "Alex", "due": "2026-11-01"}, {"name": "Sam", "due": "2026-11-01"},
                     {"name": "Kim", "due": "2026-11-20"}],
             events=[ev("2026-11-06T10:00", "1:1 Sam")])
    rows = cc.rule_one_on_one(s, {})
    assert [r["what"] for r in rows] == ["1:1 with Alex"]
    assert rows[0]["action"] == "propose a 1:1"


def test_one_on_one_outside_the_window_or_cancelled_counts_as_missing():
    # crow-today drops cancelled events, so a cancelled 1:1 is simply absent
    s = snap(people=[{"name": "Alex", "due": "2026-11-01"}], events=[ev("2026-11-12T10:00", "1:1 Alex")])
    assert len(cc.rule_one_on_one(s, {"one_on_one": {"lookahead_days": 7}})) == 1


def test_people_read_from_vault_notes():
    people = cc.people_of(str(DAY / "vault"), {"dir": "People", "task": "1:1"})
    assert people == [{"name": "Alex", "due": "2026-11-01"}, {"name": "Kim", "due": "2026-11-20"},
                      {"name": "Sam", "due": "2026-11-01"}]


# mail

def test_mail_older_than_threshold():
    s = snap(mail=[{"id": "a", "subject": "Scope", "from": "Pat", "at": "2026-10-31T16:00:00Z"},
                   {"id": "b", "subject": "Lunch", "from": "Sam", "at": "2026-11-03T07:00:00Z"}])
    rows = cc.rule_mail(s, {"mail": {"older_than_days": 2}})
    assert [(r["id"], r["action"]) for r in rows] == [("mail:a", "draft a reply")]
    assert cc.rule_mail(snap(mail=None), {}) == []


# rituals

def test_rituals_before_and_after_deadline(tmp_path):
    cfg = {"rituals": [{"id": "daily", "what": "Daily not posted", "by": "11:30", "marker": "last-daily"},
                       {"id": "weekly", "what": "Weekly not prepped", "weekday": 0, "by": "12:00", "marker": "last-weekly"}]}
    (tmp_path / "last-daily").write_text("2026-11-02\n")
    tue = dt.datetime.fromisoformat(NOON)
    rit = cc.rituals_of(cfg, str(tmp_path), None, tue)
    assert [r["id"] for r in rit] == ["daily"]  # weekly only on Monday
    assert [r["id"] for r in cc.rule_rituals(snap(rituals=rit), cfg)] == ["ritual:daily:2026-11-03"]
    early = dt.datetime.fromisoformat("2026-11-03T10:00:00+01:00")
    assert cc.rule_rituals(snap(now=early.isoformat(), rituals=cc.rituals_of(cfg, str(tmp_path), None, early)), cfg) == []
    (tmp_path / "last-daily").write_text("2026-11-03\n")
    assert cc.rule_rituals(snap(rituals=cc.rituals_of(cfg, str(tmp_path), None, tue)), cfg) == []
    mon = dt.datetime.fromisoformat("2026-11-09T12:30:00+01:00")
    rows = cc.rule_rituals(snap(now=mon.isoformat(), rituals=cc.rituals_of(cfg, str(tmp_path), None, mon)), cfg)
    assert {r["id"] for r in rows} == {"ritual:daily:2026-11-09", "ritual:weekly:2026-W46"}


def test_ritual_vault_path_pattern(tmp_path):
    (tmp_path / "Weekly 2026-W46 prep.md").write_text("x")
    r = {"id": "weekly", "weekday": 0, "path": "Weekly {year}-W{week} *.md"}
    mon = dt.datetime.fromisoformat("2026-11-09T12:30:00+01:00")
    assert cc._done(r, "/nonexistent", str(tmp_path), mon, "2026-W46")
    assert not cc._done(dict(r, path="Other {year}*.md"), "/nonexistent", str(tmp_path), mon, "2026-W46")


# tasks (DevEx check-in)

def test_overdue_task_in_a_watched_note():
    tasks = cc.tasks_of(str(DAY / "vault"), [{"id": "devex", "note": "Management/DevEx.md", "match": "check-in"}])
    assert tasks == [{"check": "devex", "text": "Team check-in survey", "due": "2026-10-30", "action": None}]
    rows = cc.rule_tasks(snap(tasks=tasks), {})
    assert [r["check"] for r in rows] == ["devex"]
    assert cc.rule_tasks(snap(now="2026-10-29T09:00:00+01:00", tasks=tasks), {}) == []


# engine

def test_state_and_row_limit():
    s = snap(people=[{"name": n, "due": "2026-11-01"} for n in "ABCDE"])
    rows = cc.evaluate(s, {"max_rows": 2})
    assert [r["what"] for r in rows] == ["1:1 with A", "1:1 with B"]
    state = {"dismissed": ["1on1:A:2026-11-01"], "snoozed": {"1on1:B:2026-11-01": "2026-11-04"}}
    assert [r["what"] for r in cc.evaluate(s, {"max_rows": 9}, state)] == ["1:1 with C", "1:1 with D", "1:1 with E"]
    back = snap(now="2026-11-04T09:00:00+01:00", people=s["people"])
    assert "1:1 with B" in [r["what"] for r in cc.evaluate(back, {"max_rows": 9}, state)]


def test_stale_today_cache_skips_calendar_rules():
    assert cc.events_of({"at": "2026-11-02T09:00:00+01:00", "events": []}, dt.date(2026, 11, 3)) is None
    week = [ev("2026-11-09T10:00", "x")]
    assert cc.events_of({"at": "2026-11-03T09:00", "events": [], "week": week}, dt.date(2026, 11, 3)) == week


# CLI

def home_copy(tmp_path):
    home = tmp_path / "hub"
    (home / "local").mkdir(parents=True)
    shutil.copy(DAY / "checks.yaml", home / "local" / "checks.yaml")
    shutil.copytree(DAY / "state", home / "local" / "state")
    shutil.copytree(DAY / "vault", home / "vault")
    (home / "bin").symlink_to(BIN)
    return (home, tmp_path), {"CROW_VAULT": str(home / "vault")}


def test_replay_2026_11_03_offline(tmp_path):
    hub, env = home_copy(tmp_path)
    p = run(hub, "crow-checks", "json", "--now", NOON, env=env)
    assert p.returncode == 0, p.stderr
    out = json.loads(p.stdout)
    sprint = [r for r in out["rows"] if r["check"] == "sprint"]
    assert len(sprint) == 1 and sprint[0]["action"] == "propose planning"
    assert {r["check"] for r in out["rows"]} == {"sprint", "one_on_one", "mail", "ritual", "devex"}
    assert out["skipped"] == {}


def test_dismiss_and_snooze_cli(tmp_path):
    hub, env = home_copy(tmp_path)
    ids = lambda now: [r["id"] for r in json.loads(run(hub, "crow-checks", "json", "--now", now, env=env).stdout)["rows"]]
    assert run(hub, "crow-checks", "dismiss", "mail:m-old", env=env).returncode == 0
    assert run(hub, "crow-checks", "snooze", "1on1:Alex:2026-11-01", "2026-11-04", env=env).returncode == 0
    now = ids(NOON)
    assert "mail:m-old" not in now and "1on1:Alex:2026-11-01" not in now
    state = json.loads((hub[0] / "local" / "state" / "checks.json").read_text())
    assert state == {"dismissed": ["mail:m-old"], "snoozed": {"1on1:Alex:2026-11-01": "2026-11-04"}}


def test_snapshot_round_trips_through_replay(tmp_path):
    hub, env = home_copy(tmp_path)
    p = run(hub, "crow-checks", "snapshot", "--now", NOON, env=env)
    f = tmp_path / "snap.json"
    f.write_text(p.stdout)
    r = run(hub, "crow-checks", "replay", str(f), env=env)
    assert r.returncode == 0, r.stderr
    assert any(x["action"] == "propose planning" for x in json.loads(r.stdout)["rows"])


# crow-today calendar window

def test_crow_today_week_window(monkeypatch):
    loader = importlib.machinery.SourceFileLoader("crow_today", str(BIN / "crow-today"))
    spec = importlib.util.spec_from_loader("crow_today", loader)
    ct = importlib.util.module_from_spec(spec)
    loader.exec_module(ct)
    seen = {}

    def fake(tok, path, extra=None):
        seen["path"] = path
        e = lambda s, subj, **k: {"subject": subj, "start": {"dateTime": s}, "end": {"dateTime": s}, **k}
        return {"value": [e("2026-11-05T10:00:00", "1:1 Alex"), e("2026-11-06T10:00:00", "gone", isCancelled=True),
                          e("2026-11-20T10:00:00", "outside")]}

    monkeypatch.setattr(ct, "get", fake)
    day = dt.datetime(2026, 11, 3)
    out = ct.events(None, day, 8)
    assert [x["subject"] for x in out] == ["1:1 Alex"]
    assert "endDateTime=2026-11-11" in seen["path"]
