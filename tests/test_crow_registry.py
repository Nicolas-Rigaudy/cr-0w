"""crow-registry: registry JSON with the first alias as short name and cheap git state."""
import importlib.machinery
import importlib.util
import json
import subprocess

from conftest import BIN, run


def load_registry():
    loader = importlib.machinery.SourceFileLoader("crow_registry", str(BIN / "crow-registry"))
    spec = importlib.util.spec_from_loader("crow_registry", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


def test_registry_json_lists_projects_with_git_state(hub, tmp_path):
    home, _ = hub
    repo = tmp_path / "repo"
    subprocess.run(["git", "init", "-q", str(repo)], check=True)
    subprocess.run(["git", "-C", str(repo), "config", "user.email", "me@example.org"], check=True)
    (repo / "a.txt").write_text("x")
    (home / "local" / "projects.yaml").write_text(
        f"identities: {{work: me@example.org}}\nprojects:\n"
        f"  - {{name: long-repo-name, aliases: [short], scope: work, path: {repo}, forge: gitlab}}\n"
        f"  - {{name: other, scope: personal, path: {tmp_path}/missing}}\n")
    out = run(hub, "crow-registry", "json")
    assert out.returncode == 0, out.stderr
    reg = json.loads(out.stdout)
    assert reg["identities"] == {"work": "me@example.org"}
    assert reg["state_dir"] == str(home / "local" / "state")
    first, second = reg["projects"]
    assert first["short"] == "short" and first["git"]["email"] == "me@example.org" and first["git"]["dirty"] == 1
    assert second["short"] == "other" and second["git"] == {"email": "", "dirty": 0, "ahead": 0, "stale_trees": 0}


def test_registry_rejects_unknown_argument(hub):
    assert run(hub, "crow-registry", "nope").returncode != 0


ROAD_DOC = """## road_001 - Demo to 1.0
> Date: 2026-01-01

# Milestones
## 0.1 - First flow (delivered)
- Goal: a loop.
- Exit signal: Met: tests pass.

## 0.2 - Correctness
- Goal: it is right.
- Exit signal: Met: checks pass.

## 0.3 - Any layout
- Goal: strangers.
- Exit signal: Three sample repos plan with config only.

## 1.0 - Public release
- Exit signal: Tag v1.0.0 gives a release.

# Sequencing
## not a milestone
- Exit signal: ignored.
"""


def test_road_info_counts_milestones_and_names_the_current_one(tmp_path):
    mod = load_registry()
    (tmp_path / "road.md").write_text(ROAD_DOC)
    road = mod.road_info(tmp_path, "road.md")
    assert road == {"done": 2, "total": 4, "current": "0.3 - Any layout",
                    "exit": "Three sample repos plan with config only."}


def test_road_info_accepts_an_absolute_path_and_all_done(tmp_path):
    mod = load_registry()
    doc = tmp_path / "r.md"
    doc.write_text("# Milestones\n## 1.0 - Done (delivered)\n- Exit signal: Met.\n")
    assert mod.road_info("/nowhere", str(doc)) == {"done": 1, "total": 1, "current": "", "exit": ""}


def test_road_info_is_none_without_a_readable_roadmap(tmp_path):
    mod = load_registry()
    assert mod.road_info(tmp_path, "") is None
    assert mod.road_info(tmp_path, "missing.md") is None
    (tmp_path / "x.md").write_text("no milestones here")
    assert mod.road_info(tmp_path, "x.md") is None
