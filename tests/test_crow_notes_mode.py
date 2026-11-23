"""Vault brick off: notes, activity lines and open threads live under CROW_HOME/notes (item_262)."""
import datetime as dt
import json

from conftest import run

OFF = "version: 1\nmodules:\n  mail: false\n  team: false\n  vault: false\n"


def off(hub):
    home, _ = hub
    (home / "crow.yaml").write_text(OFF)
    return home


def today():
    return dt.date.today().isoformat()


def test_note_and_log_round_trip_through_crow_home_notes(hub, tmp_path):
    home = off(hub)
    env = {"CROW_VAULT": str(tmp_path / "novault")}
    r = run(hub, "crow-note", "call vendor", env=env)
    assert r.returncode == 0, r.stderr
    assert run(hub, "crow-log", "[hub] brief", env=env).returncode == 0
    daily = home / "notes" / "daily" / f"{today()}.md"
    text = daily.read_text()
    assert text.startswith(f"# {today()}")
    assert "call vendor" in text.split("## 📥 Capture")[1].split("##")[0]
    assert "[hub] brief" in text.split("## CR-0W")[1]
    assert not (home / "local" / "state" / "activity").exists()


def test_threads_path_points_to_crow_home_notes(hub):
    home = off(hub)
    r = run(hub, "crow-state", "path", "threads")
    assert r.returncode == 0, r.stderr
    assert r.stdout.strip() == str(home / "notes" / "Open Threads.md")


def test_doctor_is_quiet_for_off_bricks(hub):
    off(hub)
    r = run(hub, "crow-doctor", "--json")
    assert r.returncode in (0, 1), r.stderr
    text = r.stdout.lower()
    assert "vault" not in text and "mail" not in text and "outlook" not in text


def test_eod_skill_is_not_linked_with_vault_off(hub):
    home = off(hub)
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    assert "eod" not in [p.name for p in (home / ".claude" / "skills").iterdir()]


def test_all_on_keeps_the_vault_daily(hub, vault):
    home, _ = hub
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  mail: true\n  team: true\n  vault: true\n")
    assert run(hub, "crow-note", "x", env={"CROW_VAULT": str(vault)}).returncode == 0
    assert (vault / "40 Journal" / "Daily" / f"{today()}.md").exists()
    assert not (home / "notes").exists()
