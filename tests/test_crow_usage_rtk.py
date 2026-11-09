"""crow-usage: `rtk gain` parsing holds for rtk 0.35 and 1.x output."""
import importlib.machinery
import importlib.util
import json
import subprocess
import types

import pytest
from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_usage_rtk", str(BIN / "crow-usage"))
spec = importlib.util.spec_from_loader("crow_usage_rtk", loader)
cu = importlib.util.module_from_spec(spec)
loader.exec_module(cu)

SUMMARY = {"avg_savings_pct": 80.0, "avg_time_ms": 5, "total_commands": 10, "total_input": 1000,
           "total_output": 200, "total_saved": 800, "total_time_ms": 50}
DAY = {"avg_time_ms": 5, "commands": 4, "date": "2026-11-09", "input_tokens": 400, "output_tokens": 100,
       "saved_tokens": 300, "savings_pct": 75.0, "total_time_ms": 20}
GAIN = {
    "0.35": {"summary": SUMMARY, "daily": [DAY]},
    "1.2": {"summary": {**SUMMARY, "oversized_excluded": 3}, "daily": [DAY]},  # 1.0 adds oversized_excluded
}
FAILURES = "RTK Parse Failures\n══════\n\nTotal failures:    7\n"


@pytest.mark.parametrize("version", sorted(GAIN))
def test_gain_parsing_across_versions(monkeypatch, version):
    def run(cmd, **kw):
        out = json.dumps(GAIN[version]) if "--daily" in cmd else FAILURES
        return types.SimpleNamespace(stdout=out, returncode=0)
    monkeypatch.setattr(cu, "rtk_homes", lambda: ["/h"])
    monkeypatch.setattr(subprocess, "run", run)
    s = cu.rtk_stats("2026-11-09")
    assert (s["saved_tokens"], s["input_tokens"], s["commands"], s["saved_pct"]) == (800, 1000, 10, 80.0)
    assert (s["today_saved_tokens"], s["today_commands"], s["today_saved_pct"]) == (300, 4, 75.0)
    assert cu.rtk_failures() == 7
