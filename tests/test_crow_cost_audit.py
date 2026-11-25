"""bin/crow-cost-audit: dedup, context buckets, skill classes, hub start context."""
import json

from conftest import run, write_jsonl
from test_crow_usage_projects import day, setup


def turn(mid, rid, ctx_cache, out=100, model="claude-opus-5-5", cwd="/w/hub"):
    return {"type": "assistant", "cwd": cwd, "requestId": rid, "timestamp": f"{day(0)}T12:00:00.000Z",
            "message": {"id": mid, "role": "assistant", "model": model, "content": [{"type": "text", "text": "ok"}],
                        "usage": {"input_tokens": 10, "cache_read_input_tokens": ctx_cache, "output_tokens": out}}}


def audit(hub, *args):
    r = run(hub, "crow-cost-audit", "--json", "--days", "3", "--hub", "/w/hub", *args, env={"CROW_TITLES_OFF": "1"})
    assert r.returncode == 0, r.stderr
    return json.loads(r.stdout)


def test_dedup_buckets_and_start_context(hub):
    setup(hub, "s1", [turn("a", "r1", 30_000), turn("a", "r1", 30_000),  # streamed duplicate
                      turn("b", "r2", 120_000), turn("c", "r3", 200_000, model="claude-haiku-5-5")])
    r = audit(hub)
    assert r["turns"] == 3
    assert r["buckets"]["opus-5-5|<50k"]["turns"] == 1
    assert r["buckets"]["opus-5-5|100-150k"]["turns"] == 1
    assert r["buckets"][ "haiku-5-5|>150k"]["turns"] == 1
    assert r["top_bucket"]["turns"] == 1 and r["top_bucket"]["median_context"] == 200_010
    assert r["by_skill"]["interactive"]["turns"] == 3
    assert r["start_context"] == {"sessions": 1, "median": 30_010, "min": 30_010, "max": 30_010}
    assert set(r["by_model"]) == {"opus-5-5", "haiku-5-5"}


def test_day_series_and_account_split(hub):
    _, profiles = hub
    older = {**turn("d", "r4", 5_000, model="claude-sonnet-5-5"), "timestamp": f"{day(1)}T12:00:00.000Z"}
    setup(hub, "s1", [turn("a", "r1", 30_000), older])
    write_jsonl(profiles / "acct2" / "claude-home" / ".claude" / "projects" / "-w-api" / "s2.jsonl",
                [{"type": "user", "cwd": "/w/api", "timestamp": f"{day(0)}T08:00:00.000Z",
                  "message": {"role": "user", "content": "go"}}, turn("b", "r2", 1_000, cwd="/w/api")])
    r = audit(hub)
    assert set(r["by_account"]) == {"main", "acct2"}
    assert r["by_account"]["main"]["turns"] == 2 and r["by_account"]["acct2"]["turns"] == 1
    days = {d["date"]: d for d in r["by_day"]}
    assert [d["date"] for d in r["by_day"]] == [day(1), day(0)]
    today = days[day(0)]
    assert today["turns"] == 2
    assert today["cache_read_tokens"] == 31_000 and today["input_tokens"] == 20 and today["output_tokens"] == 200
    assert today["tokens"] == 31_220
    assert set(today["by_account"]) == {"main", "acct2"}
    assert today["by_model"]["opus-5-5"]["cache_read_tokens"] == 31_000
    assert set(today["by_project"]) == {"hub", "api"}
    assert days[day(1)]["by_model"]["sonnet-5-5"]["turns"] == 1
    assert abs(sum(d["cost_usd"] for d in r["by_day"]) - r["cost_usd"]) < 0.01
