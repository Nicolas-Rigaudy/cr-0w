"""bin/crow-cost-audit: dedup, context buckets, skill classes, hub start context."""
import json

from test_crow_usage_projects import day, setup

from conftest import run


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
