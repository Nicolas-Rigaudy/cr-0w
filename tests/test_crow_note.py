import datetime as dt
import json

from conftest import run


def note(hub, vault, *args):
    return run(hub, "crow-note", *args, env={"CROW_VAULT": str(vault)})


def notes(hub, vault, *args):
    return run(hub, "crow-notes", *args, env={"CROW_VAULT": str(vault)})


def today():
    return dt.date.today().isoformat()


def test_creates_daily_from_template_and_appends(hub, vault):
    r = note(hub, vault, "call vendor", "#naval")
    assert r.returncode == 0, r.stderr
    text = (vault / "40 Journal" / "Daily" / f"{today()}.md").read_text()
    assert f"created: {today()}" in text
    assert "<%" not in text
    cap = text.split("## 📥 Capture")[1].split("## ✅")[0]
    assert "call vendor #naval" in cap and cap.count("- ") == 1


def test_second_note_appends_after_first(hub, vault):
    note(hub, vault, "one")
    note(hub, vault, "two")
    text = (vault / "40 Journal" / "Daily" / f"{today()}.md").read_text()
    cap = text.split("## 📥 Capture")[1].split("## ✅")[0]
    assert cap.index("one") < cap.index("two")
    assert "## 🌙 End of day" in text and "- **Shipped:** " in text


def test_other_section_untouched(hub, vault):
    note(hub, vault, "x")
    text = (vault / "40 Journal" / "Daily" / f"{today()}.md").read_text()
    assert text.split("## 📆 Meetings & notes")[1].split("##")[0].strip() == "-"


def test_section_option(hub, vault):
    note(hub, vault, "met the CEO", "--section", "meetings")
    text = (vault / "40 Journal" / "Daily" / f"{today()}.md").read_text()
    assert "met the CEO" in text.split("## 📆 Meetings & notes")[1]


def test_missing_template_fails(hub, vault):
    (vault / "90 Templates" / "Daily Note.md").unlink()
    r = note(hub, vault, "x")
    assert r.returncode != 0 and "template" in (r.stderr + r.stdout)


def test_notes_json(hub, vault):
    note(hub, vault, "alpha", "#team")
    note(hub, vault, "done shipping", "--section", "eod")
    (vault / "00 Inbox" / "idea.md").write_text("---\ntags: [x]\n---\n# Idea\nbuy a probe\n")
    (vault / "10 Work" / "p").mkdir()
    (vault / "10 Work" / "p" / "plan.md").write_text("# Plan\nship in Q1\n")
    r = notes(hub, vault, "--since", "2d")
    assert r.returncode == 0, r.stderr
    data = json.loads(r.stdout)
    day = data["days"][today()]
    assert any("alpha #team" in l for l in day["capture"])
    assert any("done shipping" in l for l in day["eod"])
    assert all("Shipped:**" not in l for l in day["eod"])
    assert data["inbox"][0]["head"] == ["buy a probe"]
    assert data["work"][0]["note"] == "10 Work/p/plan.md"


def test_notes_skips_old_days(hub, vault):
    old = dt.date.today() - dt.timedelta(days=30)
    (vault / "40 Journal" / "Daily" / f"{old}.md").write_text("## 📥 Capture\n- 09:00 old thing\n")
    data = json.loads(notes(hub, vault, "--since", "7d").stdout)
    assert data["days"] == {}
