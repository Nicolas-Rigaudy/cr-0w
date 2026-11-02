"""crow-review-run: findings schema, parsing a run's output, sync and start with simulated runs only."""
import importlib.machinery
import importlib.util
import json
import os
import subprocess
import sys

import pytest

from conftest import BIN, FIXTURES

REPO = "grp/demo-service"


@pytest.fixture
def home(tmp_path):
    h = tmp_path / "hub"
    (h / "local" / "reviews").mkdir(parents=True)
    (h / "local" / "reviews" / "common.md").write_text("| Unpaginated AWS list calls | grep list( |\n")
    checkout = tmp_path / "demo-service"
    checkout.mkdir()
    subprocess.run(["git", "init", "-q", str(checkout)], check=True)
    (h / "local" / "projects.yaml").write_text(
        f"projects:\n  - {{name: demo-service, path: {checkout}, remote: {REPO}, forge: gitlab}}\n")
    return h


@pytest.fixture
def rr(home, monkeypatch):
    monkeypatch.setenv("CROW_HOME", str(home))
    loader = importlib.machinery.SourceFileLoader("crow_review_run", str(BIN / "crow-review-run"))
    spec = importlib.util.spec_from_loader("crow_review_run", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


def running(rr, number=220):
    doc = {"schema": 1, "repo": REPO, "number": number, "project": "demo-service", "mode": "peer",
           "run": {"id": "rv000001", "state": "running", "reason": "", "model": "opus", "started_at": "t", "ended_at": ""},
           "summary": "", "findings": []}
    rr.save(rr.findings_path(REPO, number), doc)
    return doc


def fixture_result():
    return json.loads((FIXTURES / "review-run-claude.jsonl").read_text())["result"]


# ---------- schema ----------

def test_fixture_file_is_valid(home):
    r = subprocess.run([sys.executable, str(BIN / "crow-review-run"), "check", str(FIXTURES / "review-findings.json")],
                       capture_output=True, text=True, env={**os.environ, "CROW_HOME": str(home)})
    assert r.returncode == 0 and r.stdout.strip() == "ok", r.stdout + r.stderr


def test_schema_errors_are_reported(rr):
    doc = json.loads((FIXTURES / "review-findings.json").read_text())
    doc["findings"][0]["severity"] = "major"
    doc["findings"][1]["id"] = "f1"
    doc["findings"][1]["state"] = "maybe"
    doc["run"] = {"state": "failed"}
    errs = rr.validate_file(doc)
    assert any("findings[0].severity" in e for e in errs)
    assert any("findings[1].id" in e for e in errs)
    assert any("findings[1].state" in e for e in errs)
    assert any("run.reason" in e for e in errs)


# ---------- parsing ----------

def test_parse_output_takes_the_json_block(rr):
    doc = rr.parse_output(fixture_result())
    assert doc["findings"][0]["path"] == "src/handlers/list.py"
    assert rr.validate_output(doc) == []


def test_parse_output_reads_raw_json_and_prose_around_braces(rr):
    raw = '{"summary": "s", "findings": []}'
    assert rr.parse_output(raw) == {"summary": "s", "findings": []}
    assert rr.parse_output("Here you go: " + raw + " done") == {"summary": "s", "findings": []}
    assert rr.parse_output("no json at all") is None


# ---------- sync ----------

def test_sync_folds_a_finished_run(rr):
    running(rr)
    doc = rr.sync(REPO, 220, lookup=lambda rid: {"run_id": rid, "ended_at": "e", "exit_code": 0},
                  message=lambda row: fixture_result())
    assert doc["run"]["state"] == "done"
    assert doc["findings"][0]["id"] == "f1" and doc["findings"][0]["state"] == "pending"
    assert doc["summary"].startswith("Patterns: 4 checked")
    on_disk = rr.load(rr.findings_path(REPO, 220))
    assert rr.validate_file(on_disk) == []


def test_sync_leaves_a_running_run_alone(rr):
    before = running(rr)
    doc = rr.sync(REPO, 220, lookup=lambda rid: {"run_id": rid, "ended_at": None}, message=lambda row: "")
    assert doc == before


@pytest.mark.parametrize("row,text,reason", [
    ({"exit_code": 1}, "boom\nlast line", "run exited 1: last line"),
    ({"exit_code": 0}, "I reviewed it, all good.", "no findings JSON"),
    ({"exit_code": 0}, '{"findings": [{"title": "t"}]}', "invalid findings"),
])
def test_sync_marks_failures_with_a_reason(rr, row, text, reason):
    running(rr)
    doc = rr.sync(REPO, 220, lookup=lambda rid: {"run_id": rid, "ended_at": "e", **row}, message=lambda r: text)
    assert doc["run"]["state"] == "failed" and reason in doc["run"]["reason"]
    assert rr.validate_file(doc) == []


def test_sync_unknown_run_fails(rr):
    running(rr)
    doc = rr.sync(REPO, 220, lookup=lambda rid: None, message=lambda r: "")
    assert doc["run"]["state"] == "failed" and "does not know" in doc["run"]["reason"]


def test_show_end_to_end_with_a_simulated_cdx(home, tmp_path, rr):
    """The CLI reads a simulated finished run from a fake cdx and writes the findings."""
    running(rr)
    fake = tmp_path / "fakebin"
    fake.mkdir()
    stdout = FIXTURES / "review-run-claude.jsonl"
    (fake / "cdx").write_text(f"""#!/bin/sh
case "$1" in
  runs) echo '{{"runs": [{{"run_id": "rv000001-full", "ended_at": "e", "exit_code": 0, "provider": "claude",
         "artifacts": {{"stdout_path": "{stdout}"}}}}]}}' ;;
  run-report) echo '{{"report": {{}}}}' ;;
  *) exit 2 ;;
esac
""")
    (fake / "cdx").chmod(0o755)
    env = {**os.environ, "CROW_HOME": str(home), "PATH": f"{fake}:/usr/bin:/bin"}
    r = subprocess.run([sys.executable, str(BIN / "crow-review-run"), "show", REPO, "!220"],
                       capture_output=True, text=True, env=env)
    assert r.returncode == 0, r.stderr
    doc = json.loads(r.stdout)
    assert doc["run"]["state"] == "done" and doc["findings"][0]["pattern"] == "Unpaginated AWS list calls"


# ---------- start ----------

def test_start_writes_bundle_prompt_and_running_file(rr, monkeypatch):
    calls = []
    monkeypatch.setattr(rr, "mr_detail", lambda repo, n: {
        "files": [{"path": "a.py", "add": 1, "del": 0, "diff": "@@ -0,0 +1 @@\n+x = 1"}],
        "threads": [{"author": "sam", "path": "a.py", "line": 1, "body": "why?", "resolved": False}]})
    monkeypatch.setattr(rr, "dispatch", lambda *a: calls.append(a) or "rv000002-full")
    rr.main(["start", REPO, "221"])
    project, label, prompt, model, power = calls[0]
    assert (project, label, model) == ("demo-service", "review-221", "opus")
    text = prompt.read_text()
    assert "plan mode" in text and "ONE JSON object" in text
    bundle = next(l for l in text.splitlines() if "Read the bundle" in l).split("bundle ", 1)[1].split(":")[0]
    assert "/.git/crow-review/" in bundle
    b = open(bundle).read()
    assert "Unpaginated AWS list calls" in b and "+x = 1" in b and "[open] sam" in b
    doc = rr.load(rr.findings_path(REPO, 221))
    assert doc["run"]["state"] == "running" and doc["run"]["id"] == "rv000002-full"
    assert rr.validate_file(doc) == []


def test_start_refuses_to_drop_decisions(rr, monkeypatch):
    doc = json.loads((FIXTURES / "review-findings.json").read_text())
    rr.save(rr.findings_path(REPO, 220), {**doc, "repo": REPO})
    monkeypatch.setattr(rr, "dispatch", lambda *a: pytest.fail("must not dispatch"))
    with pytest.raises(SystemExit, match="already have a decision"):
        rr.main(["start", REPO, "220"])


def test_dispatch_asks_cdx_for_review_permission(tmp_path, monkeypatch):
    """crow-runs dispatch --permission review reaches cdx as --permission review (Claude plan mode)."""
    (tmp_path / "local").mkdir()
    (tmp_path / "local" / "projects.yaml").write_text(f"projects:\n  - {{name: demo, path: {tmp_path}}}\n")
    (tmp_path / "bin").mkdir()
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    monkeypatch.setenv("CROW_ACCOUNT", "main")
    loader = importlib.machinery.SourceFileLoader("crow_runs_rv", str(BIN / "crow-runs"))
    spec = importlib.util.spec_from_loader("crow_runs_rv", loader)
    cr = importlib.util.module_from_spec(spec)
    loader.exec_module(cr)
    seen = []

    class R:
        returncode = 0
        stdout = '{"ok": true, "run_id": "rv000003-full"}'

    def fake_run(argv, **k):
        seen.append(argv)
        return R()
    monkeypatch.setattr(cr.subprocess, "run", fake_run)
    monkeypatch.setattr(cr.subprocess, "Popen", lambda *a, **k: None)
    monkeypatch.setattr(cr, "running_runs", lambda: [{"cwd": str(tmp_path)}])  # same folder: allowed for review
    monkeypatch.setattr(cr.crow_lib, "provider_of", lambda acct: "claude")
    monkeypatch.setattr(cr, "sh", lambda *a, **k: "")
    monkeypatch.setattr(sys, "argv", ["crow-runs", "dispatch", "demo", "review-1", "--prompt", "p",
                                      "--model", "opus", "--permission", "review"])
    cr.main()
    cdx = next(a for a in seen if a[:2] == ["cdx", "run"])
    assert cdx[cdx.index("--permission") + 1] == "review"
    assert cdx[cdx.index("--model") + 1] == "opus"


# ---------- decide ----------

def test_decide_persists_and_validates(rr):
    doc = json.loads((FIXTURES / "review-findings.json").read_text())
    rr.save(rr.findings_path(REPO, 220), {**doc, "repo": REPO})
    rr.main(["decide", REPO, "220", "f1", "edited", "--note", "Il manque la pagination."])
    rr.main(["decide", REPO, "220", "f2", "refused"])
    on_disk = rr.load(rr.findings_path(REPO, 220))
    assert [(f["state"], f["note"]) for f in on_disk["findings"]] == [("edited", "Il manque la pagination."), ("refused", "")]
    assert rr.validate_file(on_disk) == []
    with pytest.raises(SystemExit, match="needs --note"):
        rr.main(["decide", REPO, "220", "f1", "chat"])
    with pytest.raises(SystemExit, match="no finding f9"):
        rr.main(["decide", REPO, "220", "f9", "accepted"])
