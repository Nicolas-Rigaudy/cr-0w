"""crow-registry: registry JSON with the first alias as short name and cheap git state."""
import json
import subprocess

from conftest import run


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
