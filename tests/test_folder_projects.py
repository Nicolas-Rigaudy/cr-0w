"""item_236: folder projects (no git) and the forge brick switched off."""
import datetime as dt
import json
import os
import subprocess
import sys

from conftest import BIN, run

sys.path.insert(0, str(BIN))
import crow_lib

FORGE_OFF = "version: 1\nmodules:\n  forge: false\n"


def checks(out):
    return {c["area"]: c for c in json.loads(out.stdout)["checks"]}


def test_layout_defaults_to_repo_with_git_and_folder_without(tmp_path):
    repo, notes = tmp_path / "repo", tmp_path / "notes"
    (repo / ".git").mkdir(parents=True)
    notes.mkdir()
    assert crow_lib.project_layout({"path": str(repo)}) == "repo"
    assert crow_lib.project_layout({"path": str(notes)}) == "folder"
    assert crow_lib.project_layout({"path": str(notes), "layout": "repo"}) == "repo"


def test_project_forge_is_none_for_a_folder(tmp_path):
    assert crow_lib.project_forge({"path": str(tmp_path), "forge": "github"}) == "none"
    assert crow_lib.project_forge({"path": str(tmp_path), "layout": "repo", "forge": "github"}) == "github"


def test_changed_files_skips_hidden_ignored_and_old(tmp_path):
    for rel in ("a.md", "sub/b.md", ".hidden/c.md", "node_modules/d.js", "drafts/e.md", "old.md"):
        p = tmp_path / rel
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text("x")
    old = (dt.datetime.now() - dt.timedelta(days=30)).timestamp()
    os.utime(tmp_path / "old.md", (old, old))
    got = {rel for rel, _ in crow_lib.changed_files(str(tmp_path), dt.date.today() - dt.timedelta(days=7), ["drafts"])}
    assert got == {"a.md", os.path.join("sub", "b.md")}
    assert len(crow_lib.changed_files(str(tmp_path), dt.date.today(), limit=1)) == 1


def test_registry_serves_layout_and_folder_work(hub, tmp_path):
    home, _ = hub
    notes = tmp_path / "notes"
    notes.mkdir()
    (notes / "today.md").write_text("x")
    (home / "local" / "projects.yaml").write_text(
        f"projects:\n  - {{name: notes, scope: personal, path: {notes}}}\n")
    out = run(hub, "crow-registry", "json")
    assert out.returncode == 0, out.stderr
    p = json.loads(out.stdout)["projects"][0]
    assert p["layout"] == "folder"
    assert p["work"] == {dt.date.today().isoformat(): 1}
    assert p["git"] == {"email": "", "dirty": 0, "ahead": 0, "stale_trees": 0}


def test_doctor_fails_a_folder_project_with_a_forge(hub, tmp_path):
    home, _ = hub
    (home / "local" / "projects.yaml").write_text(
        f"projects:\n  - {{name: notes, scope: personal, path: {tmp_path}, forge: github}}\n")
    layout = checks(run(hub, "crow-doctor", "--json"))["layout"]
    assert layout["level"] == "fail"
    assert any("notes" in d and "forge" in d for d in layout["details"])


def test_doctor_counts_repo_and_folder_projects(hub, tmp_path):
    home, _ = hub
    repo = tmp_path / "repo"
    subprocess.run(["git", "init", "-q", str(repo)], check=True)
    (home / "local" / "projects.yaml").write_text(
        f"projects:\n  - {{name: app, scope: personal, path: {repo}, forge: none}}\n"
        f"  - {{name: notes, scope: personal, path: {tmp_path}}}\n")
    layout = checks(run(hub, "crow-doctor", "--json"))["layout"]
    assert layout["level"] == "ok" and "1 repo and 1 folder" in layout["summary"]


def test_mrs_json_reports_off_and_show_exits_2_with_forge_off(hub):
    home, _ = hub
    (home / "crow.yaml").write_text(FORGE_OFF)
    out = run(hub, "crow-mrs", "json")
    assert out.returncode == 0, out.stderr
    r = json.loads(out.stdout)
    assert r["off"] is True and r["mrs"] == [] and r["forges"] == {}
    out = run(hub, "crow-mrs", "show", "app", "5")
    assert out.returncode == 2 and "forge" in out.stderr


def test_dispatch_refuses_isolate_on_a_folder_project(hub, tmp_path):
    home, _ = hub
    (home / "local" / "projects.yaml").write_text(
        f"projects:\n  - {{name: notes, scope: personal, path: {tmp_path}}}\n")
    out = run(hub, "crow-runs", "dispatch", "notes", "t1", "--isolate", "--prompt", "x")
    assert out.returncode != 0 and "folder project" in out.stderr + out.stdout
