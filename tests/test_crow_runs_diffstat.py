"""Run diff stat (item_269): notify stores `git diff --numstat <start>...HEAD` on the ledger row."""
import importlib.machinery
import importlib.util
import json
import os
import shutil
import subprocess

import pytest
from conftest import BIN

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs_diffstat", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs_diffstat", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)


def git(cwd, *a):
    return subprocess.run(["git", *a], cwd=cwd, check=True, capture_output=True, text=True).stdout.strip()


@pytest.fixture
def repo(tmp_path):
    r = tmp_path / "repo"
    r.mkdir()
    git(r, "init", "-q", "-b", "main")
    git(r, "config", "user.email", "t@example.com")
    git(r, "config", "user.name", "t")
    (r / "a.txt").write_text("one\n")
    git(r, "add", ".")
    git(r, "commit", "-qm", "base")
    return r


@pytest.fixture
def ledger(tmp_path, monkeypatch):
    monkeypatch.setattr(cr, "LEDGER", tmp_path / "runs.jsonl")
    return tmp_path / "runs.jsonl"


def test_finished_run_keeps_its_stat_after_the_worktree_is_removed(repo, ledger, tmp_path):
    wt = tmp_path / "repo.wt" / "item-1"
    git(repo, "worktree", "add", "-q", "-b", "item-1", str(wt))
    rec = {"run_id": "r1", "project": "demo", "cwd": str(wt), "start_sha": git(wt, "rev-parse", "HEAD")}
    ledger.write_text(json.dumps(rec) + "\n")
    (wt / "a.txt").write_text("one\ntwo\nthree\n")
    (wt / "b.bin").write_bytes(b"\x00\x01")
    git(wt, "add", ".")
    git(wt, "commit", "-qm", "work")
    (wt / "a.txt").write_text("uncommitted\n")  # not part of the run's commits

    stat = cr.record_diff_stat(rec)
    assert stat == {"base": rec["start_sha"], "files": 2, "insertions": 2, "deletions": 0,
                    "per_file": [{"path": "a.txt", "insertions": 2, "deletions": 0},
                                 {"path": "b.bin", "insertions": 0, "deletions": 0}],
                    "truncated": False}
    git(repo, "worktree", "remove", "--force", str(wt))
    rows = cr.ledger()
    assert len(rows) == 1 and rows[0]["project"] == "demo" and rows[0]["diffstat"] == stat


def test_per_file_list_is_capped(repo, ledger, monkeypatch):
    monkeypatch.setattr(cr, "DIFF_FILES_MAX", 2)
    base = git(repo, "rev-parse", "HEAD")
    for i in range(3):
        (repo / f"f{i}.txt").write_text("x\n")
    git(repo, "add", ".")
    git(repo, "commit", "-qm", "three")
    stat = cr.diff_stat(str(repo), base)
    assert stat["files"] == 3 and stat["insertions"] == 3 and len(stat["per_file"]) == 2 and stat["truncated"]


def test_no_git_cwd_gives_no_stat_and_no_error(tmp_path, ledger):
    plain = tmp_path / "plain"
    plain.mkdir()
    for rec in ({"run_id": "r2", "cwd": str(plain), "start_sha": "abc"},
                {"run_id": "r3", "cwd": str(tmp_path / "gone"), "start_sha": "abc"},
                {"run_id": "r4", "cwd": "", "start_sha": ""}):
        assert cr.record_diff_stat(rec) is None
    assert not ledger.exists()


def test_ledger_merges_records_and_skips_bad_lines(ledger):
    ledger.write_text('{"run_id": "a", "project": "p"}\nnot json\n{"run_id": "b"}\n'
                      '{"run_id": "a", "diffstat": {"files": 1}}\n')
    assert cr.ledger() == [{"run_id": "a", "project": "p", "diffstat": {"files": 1}}, {"run_id": "b"}]


pytestmark = pytest.mark.skipif(not shutil.which("git"), reason="needs git")
