"""crow-usage: per-day weekly % used and 5h windows from limit samples, and sample retention."""
import datetime as dt
import importlib.machinery
import importlib.util

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_usage_ld", str(BIN / "crow-usage"))
spec = importlib.util.spec_from_loader("crow_usage_ld", loader)
cu = importlib.util.module_from_spec(spec)
loader.exec_module(cu)

DAY1, DAY2 = "2026-11-24", "2026-11-25"


def at(day, h, m=0):
    y, mo, d = map(int, day.split("-"))
    return dt.datetime(y, mo, d, h, m).astimezone().timestamp()


def s(day, h, h5, wk, a="acct1", m=0):
    return {"t": at(day, h, m), "a": a, "h5": h5, "wk": wk}


def test_weekly_used_sums_drops_and_skips_the_reset():
    r = cu.limit_days([s(DAY1, 9, 100, 40), s(DAY1, 10, 80, 35), s(DAY1, 11, 70, 30),
                       s(DAY2, 9, 100, 100),  # weekly reset: jump up counts nothing
                       s(DAY2, 10, 90, 96)])
    assert r[DAY1]["acct1"]["week_used_pct"] == 10
    assert r[DAY2]["acct1"]["week_used_pct"] == 4


def test_5h_windows_split_on_jump_up():
    r = cu.limit_days([s(DAY1, 8, 100, 90), s(DAY1, 9, 40, 85), s(DAY1, 10, 0, 80),
                       s(DAY1, 13, 100, 80), s(DAY1, 14, 0, 75),
                       s(DAY1, 18, 100, 75), s(DAY1, 19, 90, 74)])
    assert r[DAY1]["acct1"]["windows_5h"] == [100, 100, 10]


def test_5h_window_split_on_gap_and_idle_windows_dropped():
    # Reset hidden in a gap: 85 after 70 is no jump up, but 6h passed since the window started.
    r = cu.limit_days([s(DAY1, 8, 80, 90), s(DAY1, 9, 70, 88), s(DAY1, 15, 85, 86),
                       s(DAY2, 1, 100, 86), s(DAY2, 2, 100, 86)])
    assert r[DAY1]["acct1"]["windows_5h"] == [30, 15]
    assert r[DAY2]["acct1"] == {"week_used_pct": 0, "windows_5h": []}


def test_accounts_kept_apart():
    r = cu.limit_days([s(DAY1, 9, 100, 50), s(DAY1, 9, 100, 60, a="acct2"),
                       s(DAY1, 10, 50, 45), s(DAY1, 10, 100, 60, a="acct2")])
    assert r[DAY1]["acct1"] == {"week_used_pct": 5, "windows_5h": [50]}
    assert r[DAY1]["acct2"] == {"week_used_pct": 0, "windows_5h": []}


def test_trim_keeps_recent_full_and_thins_older():
    now = at(DAY2, 12)
    recent = [s(DAY2, 11, 50, 50, m=m) for m in range(3)]
    older = [s(DAY1, 11, 50, 50, m=m) for m in range(10)]
    base = (now - 40 * 86400) // 3600 * 3600
    old = [{"t": base + m * 600, "a": "acct1", "h5": 1, "wk": 1} for m in range(6)]
    kept = cu.trim_samples(older + recent + old, now)
    assert recent[0] in kept and recent[1] in kept and recent[2] in kept
    assert sum(1 for x in kept if x in older) == 2  # one per 5 minutes
    assert [x for x in kept if x in old] == old[:1]  # past 35 days: one per hour, kept
