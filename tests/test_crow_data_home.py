"""Engine path vs data dir: CROW_HOME, else the checkout with local/, else $XDG_CONFIG_HOME/crow."""
import os
import shutil
import subprocess
import sys
from pathlib import Path

REPO = Path(__file__).resolve().parent.parent
sys.path.insert(0, str(REPO / "bin"))
import crow_lib


def test_crow_home_wins(monkeypatch, tmp_path):
    monkeypatch.setattr(crow_lib, "ENGINE", str(tmp_path / "engine"))
    assert crow_lib.data_home({"CROW_HOME": "/data"}) == "/data"


def test_checkout_with_local_is_the_data_home(monkeypatch, tmp_path):
    (tmp_path / "local").mkdir()
    monkeypatch.setattr(crow_lib, "ENGINE", str(tmp_path))
    assert crow_lib.data_home({"HOME": "/h"}) == str(tmp_path)


def test_install_without_local_uses_xdg(monkeypatch, tmp_path):
    monkeypatch.setattr(crow_lib, "ENGINE", str(tmp_path))
    assert crow_lib.data_home({"HOME": "/h"}) == "/h/.config/crow"
    assert crow_lib.data_home({"HOME": "/h", "XDG_CONFIG_HOME": "/x"}) == "/x/crow"


def _engine(tmp_path):
    eng = tmp_path / "engine"
    (eng / "bin").mkdir(parents=True)
    shutil.copy(REPO / "bin" / "crow", eng / "bin" / "crow")
    (eng / "VERSION").write_text("9.9.9\n")
    return eng


def _version(eng, **env):
    e = {k: v for k, v in os.environ.items() if k not in ("CROW_HOME", "XDG_CONFIG_HOME")}
    e.update(env)
    return subprocess.run([str(eng / "bin" / "crow"), "version"], check=False, capture_output=True, text=True, env=e)


def test_version_on_a_fresh_install_needs_no_data(tmp_path):
    eng = _engine(tmp_path)
    out = _version(eng, HOME=str(tmp_path / "h"))
    assert out.returncode == 0, out.stderr
    assert out.stdout.splitlines() == ["crow 9.9.9", f"engine: {eng}", f"data:   {tmp_path}/h/.config/crow"]


def test_version_in_a_checkout_with_local(tmp_path):
    eng = _engine(tmp_path)
    (eng / "local").mkdir()
    assert f"data:   {eng}" in _version(eng, HOME=str(tmp_path)).stdout
    assert "data:   /d" in _version(eng, HOME=str(tmp_path), CROW_HOME="/d").stdout
