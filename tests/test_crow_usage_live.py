"""crow-usage live_session_ids: Claude's sessions/<pid>.json beside the transcripts dir cdx reports."""
import importlib.machinery
import importlib.util
import json
import os

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_usage_live", str(BIN / "crow-usage"))
spec = importlib.util.spec_from_loader("crow_usage_live", loader)
cu = importlib.util.module_from_spec(spec)
loader.exec_module(cu)


def test_live_ids_come_from_running_pids_under_the_cdx_transcripts_dir(tmp_path, monkeypatch):
    claude = tmp_path / "p" / "home" / ".claude"
    (claude / "projects").mkdir(parents=True)
    (claude / "sessions").mkdir()
    (claude / "sessions" / f"{os.getpid()}.json").write_text(json.dumps({"sessionId": "live-1"}))
    (claude / "sessions" / "999999999.json").write_text(json.dumps({"sessionId": "dead-1"}))
    monkeypatch.setattr(cu.crow_lib, "transcript_dirs", lambda profiles=None: [("p", str(claude / "projects"))])
    assert cu.live_session_ids() == {"live-1"}


def test_live_ids_unknown_without_any_sessions_dir(monkeypatch):
    monkeypatch.setattr(cu.crow_lib, "transcript_dirs", lambda profiles=None: [])
    assert cu.live_session_ids() is None
