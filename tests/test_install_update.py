"""install.sh and crow update against a fake release served over file://."""
import hashlib
import io
import json
import os
import subprocess
import tarfile
from pathlib import Path

import pytest

REPO = Path(__file__).resolve().parent.parent


def engine_files():
    out = subprocess.run(["git", "ls-files", "-z", "--cached", "--others", "--exclude-standard"],
                         cwd=REPO, capture_output=True, text=True, check=True).stdout
    return [n for n in out.split("\0") if n and not n.startswith("local/") and (REPO / n).is_file()]


def tarball(version, extra=None):
    """Release tarball crow-<version>/ built from the working tree, VERSION set to version."""
    buf = io.BytesIO()
    over = {"VERSION": f"{version}\n", **(extra or {})}
    with tarfile.open(fileobj=buf, mode="w:gz") as tar:
        for name in engine_files():
            if name in over:
                continue
            tar.add(REPO / name, arcname=f"crow-{version}/{name}")
        for name, text in over.items():
            data = text.encode()
            info = tarfile.TarInfo(f"crow-{version}/{name}")
            info.size, info.mode = len(data), 0o644
            tar.addfile(info, io.BytesIO(data))
    return buf.getvalue()


@pytest.fixture
def release(tmp_path):
    """publish(version, corrupt=False) adds a release; it becomes the latest."""
    rel, api = tmp_path / "rel", tmp_path / "api" / "releases"
    (api / "tags").mkdir(parents=True)

    def publish(version, corrupt=False, body=""):
        d = rel / "download" / f"v{version}"
        d.mkdir(parents=True)
        blob = tarball(version, {f"changelogs/CHANGELOGS_{version.replace('.', '_')}.md": f"# {version}\n- notes of {version}\n"})
        (d / f"crow-{version}.tar.gz").write_bytes(blob)
        digest = hashlib.sha256(blob + (b"x" if corrupt else b"")).hexdigest()
        (d / "SHA256SUMS").write_text(f"{digest}  crow-{version}.tar.gz\n")
        doc = json.dumps({"tag_name": f"v{version}", "body": body or f"notes of {version}"})
        (api / "latest").write_text(doc)
        (api / "tags" / f"v{version}").write_text(doc)

    env = {"CROW_RELEASE_URL": f"file://{rel}", "CROW_API_URL": f"file://{tmp_path / 'api'}"}
    return publish, env


def user_env(tmp_path, extra):
    home = tmp_path / "home"
    home.mkdir(exist_ok=True)
    e = {k: v for k, v in os.environ.items() if not k.startswith(("CROW_", "XDG_"))}
    e.update(HOME=str(home), PATH=f"{home}/.local/bin:{os.environ['PATH']}", **extra)
    return home, e


def install(env, *args):
    return subprocess.run(["sh", str(REPO / "install.sh"), *args], check=False, capture_output=True, text=True, env=env)


def crow(env, *args):
    return subprocess.run(["crow", *args], check=False, capture_output=True, text=True, env=env)


def tree_digest(root):
    h = hashlib.sha256()
    for p in sorted(Path(root).rglob("*")):
        h.update(str(p.relative_to(root)).encode())
        if p.is_file():
            h.update(p.read_bytes())
    return h.hexdigest()


def test_install_puts_crow_on_path_with_a_data_dir(tmp_path, release):
    publish, renv = release
    publish("0.1.0")
    home, env = user_env(tmp_path, renv)
    out = install(env)
    assert out.returncode == 0, out.stderr
    v = crow(env, "version")
    assert v.returncode == 0, v.stderr
    assert v.stdout.splitlines()[0] == "crow 0.1.0"
    assert f"data:   {home}/.config/crow" in v.stdout
    assert (home / ".config/crow/local/accounts.yaml").is_file()
    assert not (home / ".local/share/crow/current/local").exists()


def test_reinstall_keeps_existing_data(tmp_path, release):
    publish, renv = release
    publish("0.1.0")
    home, env = user_env(tmp_path, renv)
    data = home / ".config/crow/local"
    data.mkdir(parents=True)
    (data / "accounts.yaml").write_text("default: mine\n")
    assert install(env).returncode == 0
    assert (data / "accounts.yaml").read_text() == "default: mine\n"


def test_checksum_mismatch_aborts(tmp_path, release):
    publish, renv = release
    publish("0.1.0", corrupt=True)
    home, env = user_env(tmp_path, renv)
    out = install(env)
    assert out.returncode != 0
    assert "checksum mismatch" in out.stderr
    assert not (home / ".local/share/crow").exists()
    assert not (home / ".local/bin/crow").exists()


def test_script_update_check_then_update_keeps_data(tmp_path, release):
    publish, renv = release
    publish("0.1.0")
    home, env = user_env(tmp_path, renv)
    assert install(env).returncode == 0
    data = home / ".config/crow/local"
    (data / "state").mkdir(exist_ok=True)
    (data / "state" / "notes.md").write_bytes(b"\x00binary and text\n")
    before = tree_digest(data)
    publish("0.2.0", body="- notes of 0.2.0")
    current = home / ".local/share/crow/current"

    chk = crow(env, "update", "--check")
    assert chk.returncode == 0, chk.stderr
    assert "0.1.0 -> 0.2.0" in chk.stdout and "notes of 0.2.0" in chk.stdout
    assert os.readlink(current) == "versions/0.1.0"

    up = crow(env, "update")
    assert up.returncode == 0, up.stderr
    assert "- notes of 0.2.0" in up.stdout and "notes of 0.1.0" not in up.stdout
    assert os.readlink(current) == "versions/0.2.0"
    assert crow(env, "version").stdout.splitlines()[0] == "crow 0.2.0"
    assert tree_digest(data) == before
    assert "up to date" in crow(env, "update").stdout


def test_script_update_to_a_given_version(tmp_path, release):
    publish, renv = release
    publish("0.1.0")
    publish("0.2.0")
    _, env = user_env(tmp_path, renv)
    assert install(env, "--version", "0.1.0").returncode == 0
    assert crow(env, "update", "--to", "0.2.0").returncode == 0
    assert crow(env, "version").stdout.splitlines()[0] == "crow 0.2.0"


def test_package_install_points_to_the_package_manager(tmp_path):
    eng = tmp_path / "eng"
    (eng / "bin").mkdir(parents=True)
    for n in ("crow", "crow-update", "crow_lib.py"):
        (eng / "bin" / n).write_bytes((REPO / "bin" / n).read_bytes())
        (eng / "bin" / n).chmod(0o755)
    (eng / "VERSION").write_text("0.1.0\n")
    (eng / ".crow-install").write_text("method=package\nmanager=pacman\n")
    _, env = user_env(tmp_path, {})
    out = subprocess.run([str(eng / "bin" / "crow"), "update"], check=False, capture_output=True, text=True, env=env)
    assert out.returncode == 0
    assert "installed by pacman" in out.stdout


def git(cwd, *args):
    return subprocess.run(["git", "-C", str(cwd), *args], capture_output=True, text=True, check=True,
                          env={**os.environ, "GIT_AUTHOR_NAME": "t", "GIT_AUTHOR_EMAIL": "t@example.com",
                               "GIT_COMMITTER_NAME": "t", "GIT_COMMITTER_EMAIL": "t@example.com"})


@pytest.fixture
def git_install(tmp_path):
    """A clone at v0.1.0 of an origin that also has v0.2.0, with local/ data in the checkout."""
    src = tmp_path / "src"
    src.mkdir()
    git(src, "init", "-q", "-b", "main")
    (src / "bin").mkdir()
    for n in ("crow", "crow-update", "crow_lib.py"):
        (src / "bin" / n).write_bytes((REPO / "bin" / n).read_bytes())
        (src / "bin" / n).chmod(0o755)
    (src / ".gitignore").write_text("local/\n")
    (src / "changelogs").mkdir()
    for v in ("0.1.0", "0.2.0"):
        (src / "VERSION").write_text(v + "\n")
        (src / "changelogs" / f"CHANGELOGS_{v.replace('.', '_')}.md").write_text(f"# {v}\n- notes of {v}\n")
        git(src, "add", "-A")
        git(src, "commit", "-q", "-m", v)
        git(src, "tag", f"v{v}")
    origin = tmp_path / "origin.git"
    subprocess.run(["git", "clone", "-q", "--bare", str(src), str(origin)], check=True)
    eng = tmp_path / "eng"
    subprocess.run(["git", "clone", "-q", str(origin), str(eng)], check=True)
    git(eng, "reset", "-q", "--hard", "v0.1.0")
    (eng / "local").mkdir()
    (eng / "local" / "projects.yaml").write_text("projects: []\n")
    _, env = user_env(tmp_path, {})
    return eng, env


def test_git_update_check_changes_nothing_then_updates(git_install):
    eng, env = git_install
    cmd = [str(eng / "bin" / "crow"), "update"]
    head = git(eng, "rev-parse", "HEAD").stdout
    chk = subprocess.run([*cmd, "--check"], check=False, capture_output=True, text=True, env=env)
    assert chk.returncode == 0, chk.stderr
    assert "0.1.0 -> 0.2.0" in chk.stdout
    assert git(eng, "rev-parse", "HEAD").stdout == head
    before = tree_digest(eng / "local")
    up = subprocess.run(cmd, check=False, capture_output=True, text=True, env=env)
    assert up.returncode == 0, up.stderr
    assert "notes of 0.2.0" in up.stdout and "notes of 0.1.0" not in up.stdout
    assert (eng / "VERSION").read_text() == "0.2.0\n"
    assert tree_digest(eng / "local") == before


def test_git_update_refuses_a_dirty_tree(git_install):
    eng, env = git_install
    (eng / "VERSION").write_text("hacked\n")
    out = subprocess.run([str(eng / "bin" / "crow"), "update"], check=False, capture_output=True, text=True, env=env)
    assert out.returncode == 1
    assert "uncommitted changes" in out.stderr


def test_restore_clones_the_data_repo_and_lists_steps(tmp_path, release):
    publish, renv = release
    publish("0.1.0")
    data_src = tmp_path / "data-src"
    data_src.mkdir()
    git(data_src, "init", "-q", "-b", "main")
    (data_src / "tokens.yaml").write_text("tokens:\n  - service: forge   # api\n    user: a@example.com\n")
    (data_src / "accounts.yaml").write_text("default: main\naccounts:\n  - {name: main, policy: use}\n")
    git(data_src, "add", "-A")
    git(data_src, "commit", "-q", "-m", "data")
    home, env = user_env(tmp_path, renv)
    out = install(env, "--restore", str(data_src))
    assert out.returncode == 0, out.stderr
    assert (home / ".config/crow/local/tokens.yaml").is_file()
    assert "crow secret set forge" in out.stdout
    assert "cdx add main && cdx login main" in out.stdout
    again = install(env, "--restore", str(data_src))
    assert again.returncode != 0 and "already has content" in again.stderr
