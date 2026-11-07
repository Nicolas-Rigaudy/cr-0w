import datetime as dt

from conftest import run
from test_crow_note import TPL, vault  # noqa: F401  (fixture)


def log(hub, vault, *args):
    return run(hub, "crow-log", *args, env={"CROW_VAULT": str(vault)})


def daily(vault):
    return vault / "40 Journal" / "Daily" / f"{dt.date.today().isoformat()}.md"


def test_creates_daily_and_section_before_end_of_day(hub, vault):
    r = log(hub, vault, "[hub] brief")
    assert r.returncode == 0, r.stderr
    text = daily(vault).read_text()
    assert "<%" not in text
    crow = text.split("## CR-0W\n")[1].split("## 🌙 End of day")[0]
    assert "] brief" in crow and crow.startswith("- ") and crow.endswith("\n\n")


def test_appends_and_keeps_rest_byte_identical(hub, vault):
    path = daily(vault)
    path.write_text("# Day\n\n## 📥 Capture\n- x\n\n## CR-0W\n- 08:00 [hub] one\n\n## 🌙 End of day\n- **Shipped:** \n")
    before = path.read_text()
    log(hub, vault, "[hub] two")
    after = path.read_text()
    head, tail = before.split("- 08:00 [hub] one\n")
    assert after.startswith(head + "- 08:00 [hub] one\n- ") and after.endswith("[hub] two\n" + tail)
    assert not list(path.parent.glob(".*.crow-tmp"))


def test_section_appended_when_no_end_of_day(hub, vault):
    path = daily(vault)
    path.write_text("# Day\n\n## Notes\n- a\n")
    log(hub, vault, "[x] y")
    text = path.read_text()
    assert text.startswith("# Day\n\n## Notes\n- a\n\n## CR-0W\n- ") and text.endswith("[x] y\n")


def test_unknown_template_call_falls_back_to_state(hub, vault):
    (vault / "90 Templates" / "Daily Note.md").write_text(TPL + "<% tp.system.prompt('x') %>\n")
    r = log(hub, vault, "[hub] brief")
    assert r.returncode == 0, r.stderr
    assert not daily(vault).exists()
    home, _ = hub
    legacy = home / "local" / "state" / "activity" / f"{dt.date.today().isoformat()}.md"
    assert "[hub] brief" in legacy.read_text()


def test_gather_reads_daily_section(hub, vault):
    log(hub, vault, "[hub] brief")
    r = run(hub, "crow-gather", "daily", "--json", env={"CROW_VAULT": str(vault), "CROW_GATHER_OFFLINE": "1"})
    assert r.returncode == 0, r.stderr
    import json
    data = json.loads(r.stdout)
    assert "[hub] brief" in data["activity_logs"][dt.date.today().isoformat()]
