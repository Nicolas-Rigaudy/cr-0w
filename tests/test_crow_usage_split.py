"""crow-usage cost split: parts sum to the session price, and the summary windows by day and kind."""
import datetime as dt
import importlib.machinery
import importlib.util

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_usage_split", str(BIN / "crow-usage"))
spec = importlib.util.spec_from_loader("crow_usage_split", loader)
cu = importlib.util.module_from_spec(spec)
loader.exec_module(cu)

PRICES = {"m": {"input": 4.0, "output": 20.0, "cache_read": 0.2}}


def fake_price(usage, model):
    p = PRICES[model]
    w, h = usage["cache_creation_tokens"], usage["cache_1h_tokens"]
    return (usage["input_tokens"] * p["input"] + usage["output_tokens"] * p["output"]
            + usage["cache_read_tokens"] * p["cache_read"] + (w - h) * p["input"] * 1.25 + h * p["input"] * 2) / 1e6


def test_parts_sum_to_price(monkeypatch):
    monkeypatch.setattr(cu, "price", fake_price)
    u = {"input_tokens": 10, "output_tokens": 2000, "cache_read_tokens": 400_000,
         "cache_creation_tokens": 50_000, "cache_1h_tokens": 30_000}
    split = cu.cost_split(u, "m")
    assert abs(sum(split.values()) - fake_price(u, "m")) < 1e-9
    assert split["reload"] == (20_000 * 4 * 1.25 + 30_000 * 4 * 2) / 1e6
    assert split["reread"] == 400_000 * 0.2 / 1e6


def test_summary_window_and_kinds():
    today = dt.date(2026, 11, 3)
    x = lambda r: {"cost_split": {"reload": r, "reread": 1.0, "output": 1.0, "input": 0.0}}
    sessions = [
        {"kind": "run", "by_day": {"2026-11-03": x(2.0), "2026-10-20": x(100.0)}},  # 10-20 is outside 14 days
        {"kind": "interactive", "by_day": {"2026-10-21": x(0.0)}},
    ]
    s = cu.split_summary(sessions, 14, today)
    assert s["since"] == "2026-10-21"
    assert s["by_kind"]["run"]["usd"] == {"reload": 2.0, "reread": 1.0, "output": 1.0, "input": 0.0}
    assert s["by_kind"]["all"]["cost_usd"] == 6.0
    assert s["by_kind"]["all"]["pct"]["reload"] == 33
    assert cu.split_line(s).startswith("split 14d $6: reload 33% · reread 33% · output 33% · input 0%")
