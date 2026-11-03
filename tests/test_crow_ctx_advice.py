"""crow-ctx compact advice (warm 65k, cold 125k) and the one-shot idle nudge (req_017 levers 9, 10)."""
import importlib.machinery
import importlib.util
import json
import os
import time
from datetime import datetime, timedelta, timezone

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_ctx_advice", str(BIN / "crow-ctx"))
spec = importlib.util.spec_from_loader("crow_ctx_advice", loader)
cc = importlib.util.module_from_spec(spec)
loader.exec_module(cc)


def sess(ctx, idle_min, one_h=1):
    end = (datetime.now(timezone.utc) - timedelta(minutes=idle_min)).isoformat()
    return {"id": "s", "context_tokens": ctx, "end": end, "cache_1h_tokens": one_h}


def test_three_fixture_sessions():
    assert cc.compact_advice(sess(80_000, 5)) == "compact"     # warm, past 65k
    assert cc.compact_advice(sess(80_000, 70)) == "reload"     # cold, under 125k: reload is cheaper
    assert cc.compact_advice(sess(130_000, 70)) == "compact"   # cold, past 125k
    assert cc.compact_advice(sess(40_000, 5)) == "none"


def test_five_minute_cache_goes_cold_fast():
    assert cc.compact_advice(sess(80_000, 10, one_h=0)) == "reload"


def transcript(tmp_path, ctx, one_h=True):
    t = tmp_path / "t.jsonl"
    cc_ = {"ephemeral_1h_input_tokens": 500, "ephemeral_5m_input_tokens": 0} if one_h else \
          {"ephemeral_1h_input_tokens": 0, "ephemeral_5m_input_tokens": 500}
    t.write_text(json.dumps({"type": "assistant", "message": {"usage": {
        "input_tokens": 10, "cache_read_input_tokens": ctx - 510, "cache_creation_input_tokens": 500,
        "cache_creation": cc_}}}) + "\n")
    return t


def test_idle_nudge_fires_once_when_still_idle(tmp_path, monkeypatch):
    sent = []
    monkeypatch.setattr(cc.subprocess, "run", lambda argv, **k: sent.append(argv))
    t = transcript(tmp_path, 90_000)
    waits = []
    msg = cc.idle_watch(str(t), "sid", "/w/proj", sleep=waits.append)
    assert msg and "proj" in msg and "/compact" in msg and len(sent) == 1
    assert 48 * 60 < waits[0] <= 50 * 60


def test_idle_nudge_skips_new_activity_small_sessions_and_short_ttl(tmp_path, monkeypatch):
    sent = []
    monkeypatch.setattr(cc.subprocess, "run", lambda argv, **k: sent.append(argv))
    t = transcript(tmp_path, 90_000)
    assert cc.idle_watch(str(t), "sid", "/w/p", sleep=lambda s: os.utime(t, (time.time() + 5, time.time() + 5))) is None
    assert cc.idle_watch(str(transcript(tmp_path, 30_000)), "sid", "/w/p", sleep=lambda s: None) is None
    assert cc.idle_watch(str(transcript(tmp_path, 90_000, one_h=False)), "sid", "/w/p", sleep=lambda s: None) is None
    assert sent == []


def test_idle_hook_skips_headless_runs(monkeypatch):
    started = []
    monkeypatch.setattr(cc.subprocess, "Popen", lambda *a, **k: started.append(a))
    monkeypatch.setenv("CLAUDE_CODE_ENTRYPOINT", "sdk-cli")
    monkeypatch.setattr(cc.sys, "stdin", __import__("io").StringIO('{"session_id": "s", "transcript_path": "/t"}'))
    cc.idle_hook()
    assert started == []
