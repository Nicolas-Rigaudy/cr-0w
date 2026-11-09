"""crow-room: admit rules over cdx status rows."""
import importlib.machinery
import importlib.util
import subprocess
import sys
from datetime import UTC, datetime, timedelta

from conftest import BIN

sys.path.insert(0, str(BIN))
loader = importlib.machinery.SourceFileLoader("crow_room", str(BIN / "crow-room"))
spec = importlib.util.spec_from_loader("crow_room", loader)
room = importlib.util.module_from_spec(spec)
loader.exec_module(room)

NOW = datetime(2026, 11, 7, 12, 0, tzinfo=UTC)


def row(r5=80, rw=60, **kw):
    base = {"session_name": "a", "provider": "claude", "policy": "use", "auth_status": "authenticated",
            "remaining_5h_pct": r5, "remaining_week_pct": rw, "rate_limit_reached": None,
            "runs_out_at_iso": None, "updated_at": (NOW - timedelta(minutes=3)).isoformat()}
    base.update(kw)
    return base


def test_room_is_the_smaller_window():
    assert room.verdict(row(80, 60), NOW)[:2] == (60, True)


def test_reserve_raises_the_bar():
    r, admit, why = room.verdict(row(45, 90), NOW, reserve=30)
    assert (r, admit) == (45, False) and "50%" in why
    assert room.verdict(row(55, 90), NOW, reserve=30)[1]


def test_blockers():
    assert not room.verdict(row(policy="ask-first"), NOW)[1]
    assert not room.verdict(row(auth_status="expired"), NOW)[1]
    assert not room.verdict(row(rate_limit_reached="5h"), NOW)[1]
    assert not room.verdict(row(r5=10), NOW)[1]


def test_burn_rate_projection():
    soon = (NOW + timedelta(minutes=20)).isoformat()
    assert not room.verdict(row(runs_out_at_iso=soon), NOW, minutes=60)[1]
    assert room.verdict(row(runs_out_at_iso=soon), NOW, minutes=10)[1]


def test_antigravity_stale_and_ollama():
    old = (NOW - timedelta(hours=30)).isoformat()
    assert room.verdict(row(provider="antigravity", updated_at=old), NOW)[:2] == (None, False)
    assert room.verdict(row(provider="antigravity"), NOW)[1]
    assert room.verdict(row(provider="ollama"), NOW)[:2] == (None, False)


def test_missing_reading():
    assert room.verdict(row(r5=None), NOW)[:2] == (None, False)


def test_cli_without_cdx(tmp_path):
    r = subprocess.run([sys.executable, str(BIN / "crow-room")], check=False, capture_output=True, text=True,
                       env={"PATH": str(tmp_path)})
    assert r.returncode != 0 and "cdx did not answer" in r.stderr
