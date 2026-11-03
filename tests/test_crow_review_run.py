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
    # Simulated runs only: a run nobody simulated is still running, cdx is never called.
    monkeypatch.setattr(mod, "cdx_run", lambda rid: {"run_id": rid, "ended_at": None})
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
    monkeypatch.setattr(rr, "dispatch", lambda *a, **k: calls.append(a) or "rv000002-full")
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
    assert cdx[cdx.index("--model") + 1] == "claude-opus-5-5"


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


# ---------- worktree and apply (simulated runs, local git only) ----------

def g(cwd, *args):
    return subprocess.run(["git", "-c", "user.name=t", "-c", "user.email=t@t", *args], cwd=cwd, check=True,
                          capture_output=True, text=True).stdout.strip()


@pytest.fixture
def forge(home, tmp_path):
    """A bare origin with branch feat, cloned as the registered checkout (no network)."""
    origin = tmp_path / "origin.git"
    g(tmp_path, "init", "-q", "--bare", str(origin))
    seed = tmp_path / "seed"
    g(tmp_path, "init", "-q", "-b", "main", str(seed))
    (seed / "a.py").write_text("x = 1\n")
    g(seed, "add", "."); g(seed, "commit", "-qm", "new[PY]: a")
    g(seed, "checkout", "-qb", "feat")
    (seed / "a.py").write_text("x = 2\n")
    g(seed, "commit", "-qam", "chg[PY]: two")
    g(seed, "push", "-q", str(origin), "main", "feat")
    checkout = tmp_path / "demo-service"
    import shutil
    shutil.rmtree(checkout)
    g(tmp_path, "clone", "-q", str(origin), str(checkout))
    return checkout


def reviewed(rr, wt, mode="self", decisions=(("accepted", ""), ("chat", "Pourquoi pas un itérateur ?"), ("refused", ""))):
    doc = json.loads((FIXTURES / "review-findings.json").read_text())
    base = doc["findings"][0]
    doc.update(repo=REPO, number=230, mode=mode, branch="feat", worktree=str(wt),
               run={"id": "rv000010-full", "state": "done", "reason": "", "model": "opus", "started_at": "t", "ended_at": "e"})
    doc["findings"] = [{**base, "id": f"f{i}", "path": "a.py", "state": s, "note": n} for i, (s, n) in enumerate(decisions, 1)]
    rr.save(rr.findings_path(REPO, 230), doc)
    return doc


def test_start_with_branch_reviews_in_a_detached_worktree(rr, forge, monkeypatch):
    calls = []
    monkeypatch.setattr(rr, "mr_detail", lambda repo, n: {"files": [], "threads": []})
    monkeypatch.setattr(rr, "dispatch", lambda *a, **k: calls.append(k) or "rv000011-full")
    rr.main(["start", REPO, "231", "--branch", "feat", "--mode", "self"])
    wt = forge.parent / "demo-service.wt" / "review-231"
    assert calls[0]["cwd"] == wt
    assert (wt / "a.py").read_text() == "x = 2\n" and g(wt, "branch", "--show-current") == ""
    doc = rr.load(rr.findings_path(REPO, 231))
    assert doc["worktree"] == str(wt) and doc["branch"] == "feat"
    # The bundle sits in the run's cwd and stays out of git status (and of the apply commit).
    assert (wt / ".crow-review" / "231.md").exists() and g(wt, "status", "--porcelain") == ""
    # A local commit in the worktree (apply output not pushed) blocks a new review.
    rr.save(rr.findings_path(REPO, 231), {**doc, "run": {**doc["run"], "state": "done"}})
    (wt / "a.py").write_text("x = 3\n")
    g(wt, "commit", "-qam", "fix[PY]: three")
    with pytest.raises(SystemExit, match="local commits"):
        rr.main(["start", REPO, "231", "--branch", "feat", "--mode", "self"])


def test_apply_resumes_the_review_with_the_decisions(rr, forge, monkeypatch):
    wt = forge.parent / "demo-service.wt" / "review-230"
    reviewed(rr, wt)
    calls = []

    def fake(project, label, prompt, model, power, **k):
        calls.append({"label": label, "prompt": prompt.read_text(), **k})
        return "ap000001-full"
    doc = rr.apply(REPO, 230, run_dispatch=fake)
    c = calls[0]
    assert c["continue_run"] == "rv000010-full" and c["cwd"] == str(wt) and c["permission"] is None
    assert "- f1 [accepted]" in c["prompt"] and "- f2 [chat]" in c["prompt"] and "Pourquoi pas" in c["prompt"]
    assert "- f3" not in c["prompt"] and "Leave every other finding alone (f3)" in c["prompt"]
    assert "Never push" in c["prompt"] and "type[SCOPE]" in c["prompt"]
    assert doc["apply"]["state"] == "running" and doc["apply"]["asked"] == {"f1": ["accepted", ""], "f2": ["chat", "Pourquoi pas un itérateur ?"]}
    assert rr.validate_file(rr.load(rr.findings_path(REPO, 230))) == []
    with pytest.raises(SystemExit, match="already running"):
        rr.apply(REPO, 230, run_dispatch=fake)


def test_peer_apply_only_answers_chat_notes(rr, tmp_path):
    reviewed(rr, tmp_path, mode="peer", decisions=(("accepted", ""), ("edited", "Texte revu")))
    with pytest.raises(SystemExit, match="no chat"):
        rr.apply(REPO, 230, run_dispatch=lambda *a, **k: pytest.fail("must not dispatch"))


def finish_apply(rr, forge, text_of, exit_code=0):
    """Run apply with a fake dispatch, let text_of(wt) act as the run, then sync the result."""
    wt = rr.review_worktree(str(forge), 230, "feat")
    reviewed(rr, wt)
    rr.apply(REPO, 230, run_dispatch=lambda *a, **k: "ap000002-full")
    text = text_of(wt)
    return rr.sync(REPO, 230, lookup=lambda rid: {"run_id": rid, "ended_at": "e", "exit_code": exit_code},
                   message=lambda row: text)


def test_apply_result_lands_on_the_findings(rr, forge):
    def run(wt):
        (wt / "a.py").write_text("x = 2  # paginated\n")
        g(wt, "commit", "-qam", "fix[PY]: paginate the list call")
        sha = g(wt, "rev-parse", "HEAD")
        return "Done.\n```json\n" + json.dumps({"context": "ok", "commit": sha, "subject": "fix[PY]: paginate the list call", "results": [
            {"id": "f1", "status": "fixed", "detail": "a.py paginates"},
            {"id": "f2", "status": "answered", "detail": "Un itérateur masque le coût des appels."}]}) + "\n```"
    doc = finish_apply(rr, forge, run)
    ap = doc["apply"]
    assert ap["state"] == "done" and ap["subject"] == "fix[PY]: paginate the list call" and ap["warnings"] == []
    f1, f2, f3 = doc["findings"]
    assert f1["apply"]["status"] == "fixed" and f1["apply"]["for"] == ["accepted", ""]
    assert f2["apply"]["detail"].startswith("Un itérateur") and f2["state"] == "chat"
    assert "apply" not in f3
    assert rr.validate_file(doc) == []
    # Nothing left: f1 and f2 are handled for the same decision. A new decision on f2 is picked up.
    assert rr.todo(doc) == []
    rr.decide(REPO, 230, "f2", "accepted")
    assert [f["id"] for f in rr.todo(rr.load(rr.findings_path(REPO, 230)))] == ["f2"]
    calls = []
    rr.apply(REPO, 230, run_dispatch=lambda *a, **k: calls.append(k) or "ap000003-full")
    assert calls[0]["continue_run"] == "ap000002-full"


def test_apply_warns_on_bad_subject_extra_files_and_pushed_commit(rr, forge):
    def run(wt):
        (wt / "a.py").write_text("x = 9\n")
        (wt / "other.py").write_text("y = 1\n")
        g(wt, "add", ".")
        g(wt, "commit", "-qm", "Fix things")
        g(wt, "push", "-q", "origin", "HEAD:refs/heads/feat")
        g(wt, "fetch", "-q", "origin")
        return json.dumps({"context": "ok", "commit": g(wt, "rev-parse", "HEAD"), "results": [{"id": "f1", "status": "fixed", "detail": "d"}]})
    w = finish_apply(rr, forge, run)["apply"]["warnings"]
    assert any("type[SCOPE]" in x for x in w)
    assert any("other.py" in x and "a.py" not in x for x in w)
    assert any("remote branch" in x for x in w)
    assert any("no result for f2" in x for x in w)


@pytest.mark.parametrize("text,code,reason", [
    ('{"context": "lost"}', 0, "lost the review context"),
    ("I fixed it.", 0, "no result JSON"),
    ("boom", 1, "apply run exited 1"),
])
def test_apply_failures_say_why(rr, forge, text, code, reason):
    doc = finish_apply(rr, forge, lambda wt: text, exit_code=code)
    assert doc["apply"]["state"] == "failed" and reason in doc["apply"]["reason"]
    assert all("apply" not in f for f in doc["findings"])
    assert rr.validate_file(doc) == []


def test_apply_records_a_resume_cdx_refuses(rr, tmp_path):
    reviewed(rr, tmp_path)

    def refused(*a, **k):
        sys.exit("crow: dispatch failed: continue_run_unavailable: The previous run's provider session expired")
    with pytest.raises(SystemExit, match="cannot resume the review session"):
        rr.apply(REPO, 230, run_dispatch=refused)
    ap = rr.load(rr.findings_path(REPO, 230))["apply"]
    assert ap["state"] == "failed" and "r runs the review again" in ap["reason"]


def test_dispatch_continue_run_pins_nothing_and_checks_the_worktree(tmp_path, monkeypatch):
    """--continue-run reaches cdx without --model/--power; --cwd must sit under <path>.wt/."""
    proj = tmp_path / "demo"
    (proj / ".." / "demo.wt" / "review-1").resolve().mkdir(parents=True)
    (tmp_path / "local").mkdir()
    (tmp_path / "local" / "projects.yaml").write_text(f"projects:\n  - {{name: demo, path: {proj}}}\n")
    proj.mkdir()
    (tmp_path / "bin").mkdir()
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    monkeypatch.setenv("CROW_ACCOUNT", "main")
    loader = importlib.machinery.SourceFileLoader("crow_runs_ap", str(BIN / "crow-runs"))
    spec = importlib.util.spec_from_loader("crow_runs_ap", loader)
    cr = importlib.util.module_from_spec(spec)
    loader.exec_module(cr)
    seen = []

    class R:
        returncode = 0
        stdout = '{"ok": true, "run_id": "ap000009-full"}'
    monkeypatch.setattr(cr.subprocess, "run", lambda argv, **k: seen.append(argv) or R())
    monkeypatch.setattr(cr.subprocess, "Popen", lambda *a, **k: None)
    monkeypatch.setattr(cr, "running_runs", lambda: [])
    monkeypatch.setattr(cr.crow_lib, "provider_of", lambda acct: "claude")
    monkeypatch.setattr(cr, "sh", lambda *a, **k: "")
    wt = tmp_path / "demo.wt" / "review-1"
    monkeypatch.setattr(sys, "argv", ["crow-runs", "dispatch", "demo", "apply-1", "--prompt", "p",
                                      "--continue-run", "rv000010-full", "--cwd", str(wt)])
    cr.main()
    cdx = next(a for a in seen if a[:2] == ["cdx", "run"])
    assert cdx[cdx.index("--continue-run") + 1] == "rv000010-full"
    assert "--model" not in cdx and "--power" not in cdx
    assert cdx[cdx.index("--cwd") + 1] == str(wt.resolve())
    monkeypatch.setattr(sys, "argv", ["crow-runs", "dispatch", "demo", "apply-1", "--prompt", "p", "--cwd", str(tmp_path)])
    with pytest.raises(SystemExit, match="worktree under"):
        cr.main()


# ---------- push and post ----------

def committed(rr, forge):
    """A self review whose apply run committed in the worktree (simulated run)."""
    def run(wt):
        (wt / "a.py").write_text("x = 2  # paginated\n")
        g(wt, "commit", "-qam", "fix[PY]: paginate the list call")
        return "```json\n" + json.dumps({"context": "ok", "commit": g(wt, "rev-parse", "HEAD"), "subject": "fix[PY]: paginate the list call",
                                          "results": [{"id": "f1", "status": "fixed", "detail": "ok"}]}) + "\n```"
    return finish_apply(rr, forge, run)


def test_self_push_shows_the_commit_and_runs_once(rr, forge):
    doc = committed(rr, forge)
    sha = doc["apply"]["commit"]
    (a,) = rr.actions(doc)
    assert a["id"] == "push" and a["text"] == "fix[PY]: paginate the list call"
    assert a["cmd"] == ["git", "-C", doc["worktree"], "push", "origin", f"{sha}:refs/heads/feat"]
    with pytest.raises(SystemExit, match="changed since it was shown"):
        rr.act(REPO, 230, "push", "0" * 12)
    assert g(forge.parent / "origin.git", "rev-parse", "feat") != sha
    rr.act(REPO, 230, "push", a["hash"])  # the bare origin is a local folder: no network
    assert g(forge.parent / "origin.git", "rev-parse", "feat") == sha
    doc = rr.load(rr.findings_path(REPO, 230))
    assert doc["apply"]["pushed"]["commit"] == sha and rr.actions(doc) == []
    with pytest.raises(SystemExit, match="no longer pending"):
        rr.act(REPO, 230, "push", a["hash"])


def test_peer_post_one_comment_per_accepted_finding(rr, tmp_path):
    doc = reviewed(rr, tmp_path, mode="peer", decisions=(("accepted", ""), ("edited", "Texte revu"), ("refused", ""), ("chat", "pourquoi")))
    doc["findings"][1].pop("line")
    rr.save(rr.findings_path(REPO, 230), doc)
    acts = rr.actions(doc)
    assert [a["id"] for a in acts] == ["post-f1", "post-f2"]
    f1, f2 = acts
    assert f1["text"] == doc["findings"][0]["comment"] and f2["text"] == "Texte revu"
    assert f1["cmd"][4].endswith("/discussions") and "position[head_sha]={head_sha}" in f1["cmd"]
    assert f2["cmd"][4].endswith("/notes") and f2["cmd"][-1] == "body=Texte revu"
    calls = []

    def call(cmd):
        calls.append(cmd)
        if cmd[:2] == ["glab", "api"] and "-X" not in cmd:
            return json.dumps({"diff_refs": {"base_sha": "b" * 40, "start_sha": "s" * 40, "head_sha": "h" * 40}})
        return json.dumps({"id": "d123"})
    rr.act(REPO, 230, "post-f1", f1["hash"], call=call)
    sent = calls[-1]
    assert "position[head_sha]=" + "h" * 40 in sent and "body=" + f1["text"] in sent
    assert len([c for c in calls if "-X" in c]) == 1
    doc = rr.load(rr.findings_path(REPO, 230))
    assert doc["findings"][0]["posted"]["ref"] == "d123"
    assert [a["id"] for a in rr.actions(doc)] == ["post-f2"]


def test_post_refuses_when_the_mr_moved_and_records_nothing_on_failure(rr, forge):
    wt = rr.review_worktree(str(forge), 230, "feat")
    reviewed(rr, wt, mode="peer", decisions=(("accepted", ""),))
    (a,) = rr.actions(rr.load(rr.findings_path(REPO, 230)))
    moved = lambda cmd: json.dumps({"diff_refs": {"base_sha": "b", "start_sha": "s", "head_sha": "f" * 40}})
    with pytest.raises(SystemExit, match="MR moved since the review"):
        rr.act(REPO, 230, a["id"], a["hash"], call=moved)
    tip = g(wt, "rev-parse", "HEAD")

    def broken(cmd):
        if "-X" in cmd:
            raise RuntimeError("403 Forbidden")
        return json.dumps({"diff_refs": {"base_sha": "b", "start_sha": "s", "head_sha": tip}})
    with pytest.raises(SystemExit, match="post-f1 failed: 403 Forbidden"):
        rr.act(REPO, 230, a["id"], a["hash"], call=broken)
    assert "posted" not in rr.load(rr.findings_path(REPO, 230))["findings"][0]


def test_github_comments_use_the_pr_head(rr):
    f = {"path": "b.go", "line": 4}
    cmd = rr.post_cmd("github", "me/demo-tool", 12, f, "Renomme.")
    assert cmd[:5] == ["gh", "api", "-X", "POST", "repos/me/demo-tool/pulls/12/comments"] and "commit_id={head_sha}" in cmd
    assert rr.post_cmd("github", "me/demo-tool", 12, {}, "Ok.") == ["gh", "pr", "comment", "12", "-R", "me/demo-tool", "--body", "Ok."]


# ---------- threads mode (reviewers' threads on the owner's MR) ----------

THREADS = [
    {"id": "disc-1", "author": "rev", "at": "t", "path": "a.py", "line": 1, "body": "Pagine cet appel.", "resolved": False},
    {"id": "disc-0", "author": "rev", "at": "t", "path": "b.py", "line": 3, "body": "Déjà réglé.", "resolved": True},
    {"id": "disc-2", "author": "rev", "at": "t", "path": "", "line": 0, "body": "Pourquoi ce nom ?", "resolved": False},
]


def test_threads_start_lists_open_threads_and_keeps_their_ids(rr, monkeypatch):
    calls = []
    monkeypatch.setattr(rr, "mr_detail", lambda repo, n: {"files": [], "threads": THREADS})
    monkeypatch.setattr(rr, "dispatch", lambda *a, **k: calls.append(a) or "rv000020-full")
    rr.main(["start", REPO, "240", "--mode", "threads"])
    doc = rr.load(rr.findings_path(REPO, 240))
    assert doc["mode"] == "threads" and [t["key"] for t in doc["threads"]] == ["T1", "T2"]
    assert [t["id"] for t in doc["threads"]] == ["disc-1", "disc-2"]
    prompt = calls[0][2].read_text()
    assert "review-threads skill" in prompt and "NO reply" in prompt and '"thread": "T1"' in prompt
    bundle = next(rr.findings_path(REPO, 240).parent.parent.parent.parent.glob("**/crow-review/240.md"), None) \
        or next((p for p in rr.HOME.parent.glob("**/240.md")), None)
    text = bundle.read_text()
    assert "### T1: rev on a.py:1" in text and "### T2: rev on the MR:0" in text and "Déjà réglé" not in text.split("## Existing")[0]


def test_self_start_switches_to_threads_when_the_mr_has_open_threads(rr, monkeypatch):
    calls = []
    monkeypatch.setattr(rr, "mr_detail", lambda repo, n: {"files": [], "threads": THREADS})
    monkeypatch.setattr(rr, "dispatch", lambda *a, **k: calls.append(a) or "rv000021-full")
    rr.main(["start", REPO, "243", "--mode", "self"])
    doc = rr.load(rr.findings_path(REPO, 243))
    assert doc["mode"] == "threads" and [t["id"] for t in doc["threads"]] == ["disc-1", "disc-2"]
    assert "NO reply" in calls[0][2].read_text()


def test_threads_start_refuses_an_mr_without_open_threads(rr, monkeypatch):
    monkeypatch.setattr(rr, "mr_detail", lambda repo, n: {"files": [], "threads": [THREADS[1]]})
    with pytest.raises(SystemExit, match="no open thread"):
        rr.main(["start", REPO, "241", "--mode", "threads"])


def threads_done(rr, wt, decisions):
    """A finished threads review: T1 has a fix, T2 is reply only, T3 the run disagrees with."""
    threads = [{"key": f"T{i}", "id": f"disc-{i}", "comment_id": 0, "author": "rev", "body": b, "path": p, "line": l}
               for i, (b, p, l) in enumerate((("Pagine.", "a.py", 1), ("Pourquoi ?", "", 0), ("Renomme.", "a.py", 1)), 1)]
    doc = {"schema": 1, "repo": REPO, "number": 242, "project": "demo-service", "mode": "threads",
           "branch": "feat", "worktree": str(wt), "threads": threads,
           "run": {"id": "rv000030-full", "state": "running", "reason": "", "model": "opus", "started_at": "t", "ended_at": ""},
           "summary": "", "findings": []}
    rr.save(rr.findings_path(REPO, 242), doc)
    out = {"summary": "Threads: 3 open", "findings": [
        {"thread": "T1", "title": "Liste non paginée", "severity": "blocking", "agree": "yes", "path": "a.py", "line": 1,
         "body": "list() sans pagination.", "comment": "", "fix": "-x = 2\n+x = 3", "test": "+assert x == 3",
         "pattern": "Unpaginated AWS list calls", "pattern_kind": "repeat"},
        {"thread": "T2", "title": "Nom", "severity": "question", "agree": "yes", "body": "Question sur le nom.",
         "comment": "Le nom suit le module voisin.", "pattern": "", "pattern_kind": ""},
        {"thread": "T3", "title": "Renommage", "severity": "nit", "agree": "no", "path": "a.py", "line": 1,
         "body": "Le nom actuel suit la convention.", "comment": "Je garde ce nom, il suit les voisins.",
         "pattern": "Rename request", "pattern_kind": "new", "how": "compare with siblings"},
        {"thread": "T9", "title": "inconnu", "severity": "nit", "body": "pas un fil ouvert"}]}
    rr.sync(REPO, 242, lookup=lambda rid: {"run_id": rid, "ended_at": "e", "exit_code": 0}, message=lambda row: json.dumps(out))
    for fid, (state, note) in zip(("f1", "f2", "f3"), decisions):
        rr.decide(REPO, 242, fid, state, note)
    return rr.load(rr.findings_path(REPO, 242))


def test_threads_fold_attaches_each_finding_to_its_thread(rr, tmp_path):
    doc = threads_done(rr, tmp_path, (("pending", ""),) * 3)
    assert doc["run"]["state"] == "done" and rr.validate_file(doc) == []
    assert [(f["id"], f["thread"]["id"]) for f in doc["findings"]] == [("f1", "disc-1"), ("f2", "disc-2"), ("f3", "disc-3")]
    assert doc["findings"][0]["test"] == "+assert x == 3" and doc["findings"][2]["agree"] == "no"


def test_threads_apply_fixes_only_threads_with_a_fix(rr, tmp_path):
    doc = threads_done(rr, tmp_path, (("edited", "Corrigé, merci."), ("accepted", ""), ("chat", "Et si on renommait quand même ?")))
    assert [f["id"] for f in rr.todo(doc)] == ["f1", "f3"]
    text = rr.apply_prompt(doc, rr.todo(doc))
    assert "- f1 [edited]" in text and "only rewrote the reply" in text and "Corrigé, merci." not in text
    assert "- f3 [chat]" in text and "Et si on renommait" in text and "- f2" not in text


def test_threads_actions_wait_for_the_push_then_reply_and_resolve(rr, tmp_path, monkeypatch):
    doc = threads_done(rr, tmp_path, (("accepted", ""), ("edited", "Le nom suit le module voisin, je le garde."), ("accepted", "")))
    # Before the fix is applied and pushed, only the reply-only threads can go.
    ids = [a["id"] for a in rr.actions(doc)]
    assert ids == ["reply-f2", "resolve-f2", "reply-f3"]  # f3: the run disagrees, the thread stays open
    reply = rr.actions(doc)[0]
    assert reply["text"] == "Le nom suit le module voisin, je le garde."
    assert reply["cmd"][:5] == ["glab", "api", "-X", "POST", "projects/grp%2Fdemo-service/merge_requests/242/discussions/disc-2/notes"]
    doc["findings"][0]["apply"] = {"status": "fixed", "detail": "ok", "for": ["accepted", ""], "run": "ap1"}
    doc["apply"] = {"id": "ap1", "state": "done", "reason": "", "commit": "c" * 40, "subject": "fix[PY]: x", "warnings": []}
    rr.save(rr.findings_path(REPO, 242), doc)
    assert [a["id"] for a in rr.actions(doc)][0] == "push" and "resolve-f1" not in [a["id"] for a in rr.actions(doc)]
    doc["apply"]["pushed"] = {"commit": "c" * 40, "at": "t"}
    rr.save(rr.findings_path(REPO, 242), doc)
    acts = {a["id"]: a for a in rr.actions(doc)}
    assert "reply-f1" not in acts  # fixed thread: resolved silently
    assert acts["resolve-f1"]["cmd"] == ["glab", "api", "-X", "PUT", "projects/grp%2Fdemo-service/merge_requests/242/discussions/disc-1", "-f", "resolved=true"]
    calls = []
    rr.act(REPO, 242, "resolve-f1", acts["resolve-f1"]["hash"], call=lambda c: calls.append(c) or "{}")
    rr.act(REPO, 242, "reply-f2", acts["reply-f2"]["hash"], call=lambda c: calls.append(c) or '{"id": 77}')
    assert len(calls) == 2
    doc = rr.load(rr.findings_path(REPO, 242))
    assert doc["findings"][0]["resolved"] and doc["findings"][1]["posted"]["ref"] == "77"
    assert [a["id"] for a in rr.actions(doc)] == ["resolve-f2", "reply-f3"]


def test_resolving_a_thread_feeds_the_ledgers(rr, tmp_path, home):
    (home / "local" / "reviews" / "common.md").write_text(
        "| Pattern | Count | Seen in | How |\n|---|---|---|---|\n| Unpaginated AWS list calls | 2 | demo-service !1 | grep list( |\n")
    doc = threads_done(rr, tmp_path, (("accepted", ""), ("accepted", ""), ("accepted", "")))
    f1, f3 = doc["findings"][0], doc["findings"][2]
    assert rr.feed_ledger(doc, f1) == "bumped in common.md"
    assert "| Unpaginated AWS list calls | 3 | demo-service !1 demo-service !242 | grep list( |" in (home / "local" / "reviews" / "common.md").read_text()
    assert rr.feed_ledger(doc, f1) == ""  # once per thread
    assert rr.feed_ledger(doc, f3) == ""  # the run disagreed with the reviewer
    f3["agree"] = "partly"
    assert rr.feed_ledger(doc, f3) == "added to demo-service.md"
    assert "| Rename request | 1 | !242 | compare with siblings |" in (home / "local" / "reviews" / "demo-service.md").read_text()


def test_github_thread_reply_and_resolve(rr):
    t = {"id": "PRRT_x", "comment_id": 991}
    assert rr.thread_cmd("github", "me/tool", 5, t, "reply", "Ok.") == \
        ["gh", "api", "-X", "POST", "repos/me/tool/pulls/5/comments/991/replies", "-f", "body=Ok."]
    cmd = rr.thread_cmd("github", "me/tool", 5, t, "resolve")
    assert "resolveReviewThread" in cmd[4] and cmd[-1] == "id=PRRT_x"


def test_cdx_run_finds_a_running_run(home, monkeypatch):
    """The --since list leaves running runs out; the lookup also asks --running."""
    monkeypatch.setenv("CROW_HOME", str(home))
    loader = importlib.machinery.SourceFileLoader("crow_review_run_raw", str(BIN / "crow-review-run"))
    mod = importlib.util.module_from_spec(importlib.util.spec_from_loader("crow_review_run_raw", loader))
    loader.exec_module(mod)
    lists = {"--since": [{"run_id": "old1"}], "--running": [{"run_id": "abc123full", "ended_at": None}]}
    monkeypatch.setattr(mod, "sh", lambda *a: json.dumps({"runs": lists[a[2]]}))
    assert mod.cdx_run("abc123")["run_id"] == "abc123full" and mod.cdx_run("zzz") is None
