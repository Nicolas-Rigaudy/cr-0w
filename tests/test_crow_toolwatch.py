import importlib.machinery
import importlib.util
import os
import types
from pathlib import Path

BIN = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), "bin", "crow-toolwatch")

CONFIG = """# comment kept
tools:
  alpha:
    version: echo 'alpha 1.2.0'
    source: {github_releases: o/alpha}
    last_read: 1.2.0
  beta:
    version: echo 'beta version 2.0.1 (build)'
    source: {github_file: 'o/beta:CHANGELOG.md'}
    uses: '\\bbeta\\b'
    last_read: 2.0.1
  gone:
    version: "false"
    source: {github_releases: o/gone}
    last_read:
"""

CHANGELOG = """# Changelog

## 2.1.0
- beta run gains --fast ([abc1234](https://x/commit/abc1234))

## 2.0.2
- fix beta check output

## 2.0.1
- old entry, already read
"""


def releases_alpha(newest="1.2.0"):
    return [
        {"tag_name": "v1.3.0-rc.1", "prerelease": True, "body": "rc"},
        {"tag_name": f"v{newest}", "published_at": "2026-11-01T00:00:00Z", "body": "alpha notes"},
        {"tag_name": "v1.1.0", "published_at": "2026-10-01T00:00:00Z", "body": "older"},
    ]


def fake_fetch(alpha_newest="1.2.0", beta=CHANGELOG):
    calls = []

    def fetch(path, raw=False):
        calls.append(path)
        if path.startswith("repos/o/alpha/"):
            return releases_alpha(alpha_newest)
        if path.startswith("repos/o/gone/"):
            return [{"tag_name": "v9.0.0", "body": "x"}]
        return beta
    fetch.calls = calls
    return fetch


def load(tmp_path):
    os.environ.setdefault("CROW_ACCOUNT", "test")
    loader = importlib.machinery.SourceFileLoader("crow_toolwatch", BIN)
    spec = importlib.util.spec_from_loader("crow_toolwatch", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    home = tmp_path / "home"
    (home / "bin").mkdir(parents=True)
    (home / "bin" / "crow-x").write_text("beta run --fast\nbeta check\n")
    (home / "logics" / "backlog").mkdir(parents=True)
    (home / "logics" / "backlog" / "item_001_wait_for_beta.md").write_text("> Status: Ready\nneeds beta 2.1\n")
    (home / "logics" / "backlog" / "item_002_done.md").write_text("> Status: Done\nbeta\n")
    (home / "local").mkdir()
    (home / "local" / "tools.yaml").write_text(CONFIG)
    mod.CROW_HOME = str(home)
    mod.CONFIG = str(home / "local" / "tools.yaml")
    mod.STATE = str(tmp_path / "state")
    mod.LOG = str(tmp_path / "state" / "worker.log")
    mod.LOCK = str(tmp_path / "state" / ".lock")
    mod.NOTES = str(tmp_path / "notes")
    return mod


def no_model(prompt_file):
    raise AssertionError("no model run expected")


def test_nothing_new_starts_no_model(tmp_path):
    """AC1: fake config and version commands, nothing unread: exit 0, no model, no note."""
    mod = load(tmp_path)
    rc = mod.cmd_run(types.SimpleNamespace(dry_run=False), fetch=fake_fetch(beta="## 2.0.1\n- read\n"), model=no_model)
    assert rc == 0
    assert not os.path.exists(mod.NOTES)
    rows = {r["tool"]: r for r in mod.check(mod.load_config(), fake_fetch())}
    assert rows["alpha"]["installed"] == "1.2.0" and rows["alpha"]["latest"] == "1.2.0"  # rc skipped
    assert rows["gone"]["installed"] is None and rows["gone"]["unread"] == []


def test_run_prompt_holds_unread_and_uses_then_bumps(tmp_path):
    """AC2: prompt = unread entries + matching crow uses; last_read bumped only after the note."""
    mod = load(tmp_path)
    seen = {}

    def model(prompt_file):
        seen["prompt"] = Path(prompt_file).read_text()
        seen["config_during_run"] = Path(mod.CONFIG).read_text()
        return True, "## beta 2.0.2..2.1.0\n- [ ] adopt --fast (bin/crow-x)"

    rc = mod.cmd_run(types.SimpleNamespace(dry_run=False), fetch=fake_fetch(alpha_newest="1.4.0"), model=model,
                     today="2026-11-09")
    assert rc == 0
    p = seen["prompt"]
    assert "beta 2.1.0" in p and "beta 2.0.2" in p and "old entry" not in p
    assert "alpha 1.4.0" in p and "alpha notes" in p and "older" not in p and "rc" not in p.split("# alpha")[1][:40]
    assert "abc1234" not in p  # commit links stripped
    assert "`beta run --fast`: bin/crow-x" in p and "item_001" in p and "item_002" not in p
    assert "last_read: 2.0.1" in seen["config_during_run"]  # not bumped before the note
    note = Path(os.path.join(mod.NOTES, "Tool updates 2026-11-09.md")).read_text()
    assert "- [ ] adopt --fast" in note and "gone: not installed (latest 9.0.0)" in note
    cfg = Path(mod.CONFIG).read_text()
    assert "# comment kept" in cfg and "last_read: 2.1.0" in cfg and "last_read: 1.4.0" in cfg
    assert mod.open_proposals(mod.latest_note()) == ["- [ ] adopt --fast (bin/crow-x)"]


def test_failed_run_keeps_last_read(tmp_path):
    mod = load(tmp_path)
    rc = mod.cmd_run(types.SimpleNamespace(dry_run=False), fetch=fake_fetch(), model=lambda f: (False, "limit"))
    assert rc == 1
    assert "last_read: 2.0.1" in Path(mod.CONFIG).read_text()
    assert not os.path.exists(mod.NOTES) or not os.listdir(mod.NOTES)
