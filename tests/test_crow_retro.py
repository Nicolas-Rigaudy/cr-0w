import importlib.machinery
import importlib.util
import json
import os
import time

BIN = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "bin", "crow-retro")


def load(tmp_path):
    os.environ.setdefault("CROW_ACCOUNT", "test")  # the import must not need local/accounts.yaml
    loader = importlib.machinery.SourceFileLoader("crow_retro", BIN)
    spec = importlib.util.spec_from_loader("crow_retro", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    mod.STATE = str(tmp_path)
    mod.QUEUE = str(tmp_path / "queue")
    return mod


def put_failed(tmp_path, name, **entry):
    d = tmp_path / "failed"
    d.mkdir(exist_ok=True)
    (d / name).write_text(json.dumps(entry))
    return d / name


def test_limit_and_instant_provider_refusals_are_transient(tmp_path):
    mod = load(tmp_path)
    assert mod.transient("You've hit your session limit")
    assert mod.transient('{"duration_api_ms":0,"total_cost_usd":0,"usage":{}}')
    assert not mod.transient("digest failed: no such file")


def test_requeue_waits_then_stops_after_max_attempts(tmp_path):
    mod = load(tmp_path)
    f = put_failed(tmp_path, "a.json", session_id="a", transient=True, attempts=1)
    assert mod.requeue_failed() == 0  # too fresh
    os.utime(f, (time.time() - 3600,) * 2)
    assert mod.requeue_failed() == 1
    assert (tmp_path / "queue" / "a.json").exists()
    g = put_failed(tmp_path, "b.json", session_id="b", transient=True, attempts=mod.MAX_ATTEMPTS)
    os.utime(g, (time.time() - 3600,) * 2)
    assert mod.requeue_failed(force=True) == 0
    h = put_failed(tmp_path, "c.json", session_id="c", transient=False, attempts=1)
    assert mod.requeue_failed(force=True) == 0 and h.exists()


# ---------- daily batch pre-filter (item_024) ----------

import datetime as dt  # noqa: E402
import shutil  # noqa: E402
import types  # noqa: E402

FIX = os.path.join(os.path.dirname(os.path.abspath(__file__)), "fixtures", "retro")


def fixture(name):
    return os.path.join(FIX, name + ".jsonl")


def test_prefilter_signals_on_fixture_transcripts(tmp_path):
    mod = load(tmp_path)
    s = mod.signals(fixture("corrected"))
    # first prompt (the hand-off) and system reminders do not count; the interrupt does
    assert (s["prompts"], s["corrections"], s["interrupts"]) == (3, 1, 1)
    assert mod.classify(s) == (True, "1 corrections, 1 interrupts")
    assert mod.classify(mod.signals(fixture("clean"))) == (False, "no correction signal")
    assert mod.classify(mod.signals(fixture("short"))) == (False, "2 prompt(s), under 3")
    assert mod.classify(mod.signals(fixture("oneshot")))[1] == "one-shot /brief"
    assert mod.classify(mod.signals(fixture("headless")))[1].startswith("headless")
    r = mod.signals(fixture("rejected"))
    assert (r["rejected"], r["corrections"], r["prompts"]) == (1, 0, 3)  # sidechain text ignored
    assert mod.classify(r)[0]
    e = mod.signals(fixture("errors"))
    assert e["repeated_errors"] == 1 and mod.classify(e) == (True, "1 repeated errors")


def test_prefilter_is_deterministic(tmp_path):
    mod = load(tmp_path)
    for name in ("corrected", "clean", "rejected", "errors"):
        assert mod.signals(fixture(name)) == mod.signals(fixture(name))


def _profile(tmp_path, names, mtime):
    d = tmp_path / "profiles" / "acct" / "claude-home" / ".claude" / "projects" / "-work-proj"
    d.mkdir(parents=True)
    paths = []
    for n in names:
        p = d / f"{n}.jsonl"
        shutil.copy(fixture(n), p)
        os.utime(p, (mtime, mtime))
        paths.append(("acct", str(p)))
    return paths


def test_batch_runs_once_a_day_and_reports_skips(tmp_path, monkeypatch):
    mod = load(tmp_path)
    mod.LOCK, mod.LOG = str(tmp_path / ".lock"), str(tmp_path / "log")
    mod.BATCH_STATE = str(tmp_path / "batch.json")
    now = dt.datetime(2026, 11, 9, 19, 0).timestamp()  # a Monday
    files = _profile(tmp_path, ["corrected", "clean", "headless", "errors"], now - 2 * 86400)  # Saturday
    monkeypatch.setattr(mod.crow_lib, "claude_transcripts", lambda profiles=None: iter(files))
    mod.digest = lambda row: str(tmp_path / f"{row['session']}.md")
    mod.save_batch_state({"last_end": now - 3 * 86400})  # last run Friday
    calls = []

    def analyse(manifest, n):
        calls.append((manifest, n))
        open(mod.report_path("2026-11-09"), "w").write("# Retro batch: 2026-11-09\n## corrected\n")
        return {"ok": True, "run_id": "r1"}

    args = types.SimpleNamespace(since=None, force=False, dry_run=False)
    assert mod.cmd_run(args, now=now, analyse=analyse) == 0
    assert len(calls) == 1 and calls[0][1] == 2  # weekend sessions covered, two kept
    manifest = open(calls[0][0]).read()
    assert "session=corrected" in manifest and "session=errors" in manifest and "clean" not in manifest
    report = open(mod.report_path("2026-11-09")).read()
    assert "## Skipped (pre-filter)" in report
    assert "| clean | proj | no correction signal |" in report
    assert "| headless | proj | headless run" in report
    # second run the same day: no model call
    assert mod.cmd_run(args, now=now + 60, analyse=analyse) == 0 and len(calls) == 1
    # forced rerun: already-analysed sessions are skipped, so no model call either
    args.force = True
    assert mod.cmd_run(args, now=now + 120, analyse=analyse) == 0 and len(calls) == 1


def test_failed_batch_keeps_the_window(tmp_path, monkeypatch):
    mod = load(tmp_path)
    mod.LOCK, mod.LOG = str(tmp_path / ".lock"), str(tmp_path / "log")
    mod.BATCH_STATE = str(tmp_path / "batch.json")
    now = dt.datetime(2026, 11, 10, 19, 0).timestamp()
    files = _profile(tmp_path, ["corrected"], now - 3600)
    monkeypatch.setattr(mod.crow_lib, "claude_transcripts", lambda profiles=None: iter(files))
    mod.digest = lambda row: "d.md"
    mod.save_batch_state({"last_end": now - 86400})
    args = types.SimpleNamespace(since=None, force=False, dry_run=False)
    assert mod.cmd_run(args, now=now, analyse=lambda m, n: {"ok": False, "error": "usage limit"}) == 1
    st = mod.load_batch_state()
    assert st["last_end"] == now - 86400 and "last_day" not in st


def test_batch_budget_scales_and_caps(tmp_path, monkeypatch):
    mod = load(tmp_path)
    monkeypatch.delenv("CROW_RETRO_BUDGET", raising=False)
    assert mod.batch_budget(2) == 2.5 and mod.batch_budget(20) == 8.0
    monkeypatch.setenv("CROW_RETRO_BUDGET", "1.5")
    assert mod.batch_budget(5) == 1.5
