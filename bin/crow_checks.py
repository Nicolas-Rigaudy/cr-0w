"""Checks engine for bin/crow-checks: management gaps as rows for the deck and the brief.

A snapshot is built from cached files only (crow-today cache, sprint cache, vault notes,
marker files). Each rule is a pure function rule(snap, cfg) -> rows, so tests replay a
day from a fixture with no network. A row is {id, check, what, why, action, due}.
"""
import datetime as dt
import glob
import json
import os
import re
from pathlib import Path

MAX_ROWS = 3
TASK = re.compile(r"^\s*- \[ \] (.*?)📅\s*(\d{4}-\d{2}-\d{2})")
DAYS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]


def day_label(d):
    return f"{DAYS[d.weekday()]} {d.day} {d.strftime('%b')}"


def workdays(cfg):
    return set(cfg.get("workdays", [0, 1, 2, 3, 4]))


def add_workdays(d, n, days):
    while n > 0:
        d += dt.timedelta(days=1)
        if d.weekday() in days:
            n -= 1
    return d


def workdays_between(a, b, days):
    """Workdays after a up to and including b."""
    n, d = 0, a
    while d < b:
        d += dt.timedelta(days=1)
        if d.weekday() in days:
            n += 1
    return n


def event_day(e):
    return dt.date.fromisoformat(e["start"][:10])


def subject_has(e, words):
    s = (e.get("subject") or "").lower()
    return any(w.lower() in s for w in words)


# Rules: each takes the snapshot and the config, returns rows.

def rule_sprint(snap, cfg):
    """Sprint ending within N workdays with a ritual missing from the calendar."""
    c = cfg.get("sprint") or {}
    sp, events = snap.get("sprint"), snap.get("events")
    if not sp or events is None:
        return []
    today, days = snap["today"], workdays(cfg)
    end = dt.date.fromisoformat(sp["end"])
    warn = c.get("warn_workdays", 2)
    if end < today or workdays_between(today, end, days) > warn:
        return []
    until = add_workdays(end, warn, days)
    rituals = c.get("rituals") or {"planning": ["planning"], "review": ["review"], "retro": ["retro"]}
    found = {k for k, words in rituals.items()
             if any(subject_has(e, words) and today <= event_day(e) <= until for e in events)}
    missing = [k for k in rituals if k not in found]
    if not missing:
        return []
    name = sp.get("name") or "Sprint"
    return [{"id": f"sprint:{name}:{sp['end']}", "check": "sprint",
             "what": f"{name} ends {day_label(end)}",
             "why": "no " + ", ".join(missing) + " in the calendar",
             "action": "propose " + missing[0], "due": sp["end"]}]


def rule_one_on_one(snap, cfg):
    """1:1 due (vault task date passed) with no 1:1 in the calendar in the next N days."""
    c = cfg.get("one_on_one") or {}
    events = snap.get("events")
    if events is None:
        return []
    today = snap["today"]
    ahead = today + dt.timedelta(days=c.get("lookahead_days", 7))
    words = c.get("match") or ["1:1", "1-1", "one-on-one"]
    out = []
    for p in snap.get("people") or []:
        due = dt.date.fromisoformat(p["due"])
        if due > today:
            continue
        name = p["name"]
        booked = any(name.lower() in (e.get("subject") or "").lower() and subject_has(e, words)
                     and today <= event_day(e) <= ahead for e in events)
        if booked:
            continue
        late = (today - due).days
        out.append({"id": f"1on1:{name}:{p['due']}", "check": "one_on_one", "what": f"1:1 with {name}",
                    "why": f"due {day_label(due)}" + (f", {late}d late" if late else "") + ", none booked",
                    "action": "propose a 1:1", "due": p["due"]})
    return out


def rule_mail(snap, cfg):
    """Team or management mail unanswered for more than N days (crow-today's important rule)."""
    mail = snap.get("mail")
    if mail is None:
        return []
    age = (cfg.get("mail") or {}).get("older_than_days", 2)
    now = snap["now"]
    out = []
    for m in mail:
        at = dt.datetime.fromisoformat(m["at"])
        days = (now - at).total_seconds() / 86400
        if days <= age:
            continue
        out.append({"id": "mail:" + m["id"], "check": "mail", "what": f"{m.get('from', '')}: {m.get('subject', '')}",
                    "why": f"unanswered for {int(days)}d", "action": "draft a reply",
                    "due": at.date().isoformat(), "ref": m["id"]})
    return out


def rule_rituals(snap, cfg):
    """A recurring ritual not done by its deadline (weekly prepped by Monday noon, daily by 11:30)."""
    out = []
    for r in snap.get("rituals") or []:
        if r["done"] or snap["now"] < dt.datetime.fromisoformat(r["deadline"]):
            continue
        out.append({"id": f"ritual:{r['id']}:{r['period']}", "check": "ritual", "what": r["what"],
                    "why": "not done by " + r["deadline"][11:16], "action": r.get("action") or "do it now",
                    "due": r["deadline"][:10]})
    return out


def rule_tasks(snap, cfg):
    """Open vault tasks past their due date in the watched notes (DevEx check-in)."""
    out = []
    for t in snap.get("tasks") or []:
        if dt.date.fromisoformat(t["due"]) > snap["today"]:
            continue
        out.append({"id": f"task:{t['check']}:{t['due']}:{t['text'][:40]}", "check": t["check"],
                    "what": t["text"], "why": f"overdue since {day_label(dt.date.fromisoformat(t['due']))}",
                    "action": t.get("action") or "do it or re-date it", "due": t["due"]})
    return out


RULES = [rule_sprint, rule_one_on_one, rule_mail, rule_rituals, rule_tasks]


def evaluate(snap, cfg, state=None):
    """All rule rows minus dismissed and snoozed ones, at most max_rows per check."""
    state = state or {}
    dismissed = set(state.get("dismissed") or [])
    snoozed = state.get("snoozed") or {}
    today = snap["today"].isoformat()
    limit = cfg.get("max_rows", MAX_ROWS)
    out = []
    for rule in RULES:
        rows = [r for r in rule(snap, cfg)
                if r["id"] not in dismissed and snoozed.get(r["id"], "") <= today]
        rows.sort(key=lambda r: r["due"])
        out += rows[:limit]
    return out


# Snapshot: cached files to plain data.

def _load_json(path):
    try:
        with open(path) as f:
            return json.load(f)
    except (OSError, ValueError):
        return None


def sprint_of(cfg, home, today):
    """Current sprint {name, end}: the cache file when it covers today, else the cadence."""
    c = cfg.get("sprint") or {}
    if c.get("file"):
        sp = _load_json(os.path.join(home, c["file"]))
        if sp and sp.get("end") and sp["end"] >= today.isoformat():
            return {"name": sp.get("name") or "Sprint", "end": sp["end"]}
    if c.get("anchor_end") and c.get("length_days"):
        end = dt.date.fromisoformat(str(c["anchor_end"]))
        step = dt.timedelta(days=c["length_days"])
        while end < today:
            end += step
        while end - step >= today:
            end -= step
        return {"name": c.get("name") or "Sprint", "end": end.isoformat()}
    return None


def events_of(today_json, today):
    """Calendar events from the crow-today cache; None when the cache is not from today."""
    if not today_json or not str(today_json.get("at", "")).startswith(today.isoformat()):
        return None
    if today_json.get("week") is not None:
        return today_json["week"]
    return (today_json.get("events") or []) + (today_json.get("tomorrow") or [])


def people_of(vault, c):
    """Open 1:1 tasks per person note: [{name, due}], earliest open due date per person."""
    if not vault or not c.get("dir"):
        return []
    word = c.get("task", "1:1")
    out = {}
    for path in sorted(Path(vault, c["dir"]).glob("*.md")):
        name = path.stem.lstrip("@")
        for line in path.read_text(errors="replace").splitlines():
            m = TASK.match(line)
            if m and word in m.group(1):
                out[name] = min(out.get(name, m.group(2)), m.group(2))
    return [{"name": k, "due": v} for k, v in sorted(out.items())]


def tasks_of(vault, entries):
    out = []
    for t in entries or []:
        path = Path(vault or "", t["note"])
        if not vault or not path.is_file():
            continue
        for line in path.read_text(errors="replace").splitlines():
            m = TASK.match(line)
            if m and (not t.get("match") or t["match"].lower() in m.group(1).lower()):
                text = re.sub(r"[🔁⏫🔼🔽].*$", "", m.group(1)).strip()
                out.append({"check": t["id"], "text": text, "due": m.group(2), "action": t.get("action")})
    return out


def _period(r, now):
    if r.get("weekday") is not None:
        y, w, _ = now.isocalendar()
        return f"{y}-W{w:02d}"
    return now.date().isoformat()


def _done(r, home, vault, now, period):
    """A marker file holding a date of this period, or a vault file matching the pattern."""
    if r.get("marker"):
        try:
            txt = Path(home, r["marker"]).read_text().strip()[:10]
        except OSError:
            txt = ""
        start = now.date() - dt.timedelta(days=now.weekday()) if r.get("weekday") is not None else now.date()
        return txt >= start.isoformat()
    if r.get("path") and vault:
        y, w, _ = now.isocalendar()
        pat = r["path"].format(date=now.date().isoformat(), week=w, year=y)
        return bool(glob.glob(os.path.join(glob.escape(str(vault)), pat)))
    return True


def rituals_of(cfg, home, vault, now):
    out = []
    for r in cfg.get("rituals") or []:
        if r.get("weekday") is not None and now.weekday() != r["weekday"]:
            continue
        if r.get("weekday") is None and now.weekday() not in workdays(cfg):
            continue
        hh, mm = map(int, str(r.get("by", "12:00")).split(":"))
        deadline = now.replace(hour=hh, minute=mm, second=0, microsecond=0)
        period = _period(r, now)
        out.append({"id": r["id"], "what": r.get("what") or r["id"], "deadline": deadline.isoformat(),
                    "period": period, "action": r.get("action"), "done": _done(r, home, vault, now, period)})
    return out


def snapshot(cfg, home, vault, now):
    """Everything the rules read, from files only. A missing source is None (rule skipped)."""
    today = now.date()
    tj = _load_json(os.path.join(home, "local", "state", "today.json"))
    events = events_of(tj, today)
    return {
        "now": now, "today": today,
        "sprint": sprint_of(cfg, home, today),
        "events": events,
        "people": people_of(vault, cfg.get("one_on_one") or {}),
        "mail": (tj or {}).get("important") if events is not None else None,
        "rituals": rituals_of(cfg, home, vault, now),
        "tasks": tasks_of(vault, cfg.get("tasks")),
    }


def skipped(snap):
    """Sources the rules could not read, for the JSON output."""
    out = {}
    if snap.get("events") is None:
        out["calendar"] = "no crow-today cache from today"
    if snap.get("sprint") is None:
        out["sprint"] = "no sprint cache or cadence in the config"
    return out


def to_json(snap):
    """A snapshot as JSON-ready data (replay files use this shape)."""
    d = dict(snap)
    d["now"], d["today"] = snap["now"].isoformat(), snap["today"].isoformat()
    return d


def from_json(d):
    snap = dict(d)
    snap["now"] = dt.datetime.fromisoformat(d["now"])
    snap["today"] = snap["now"].date()
    return snap
