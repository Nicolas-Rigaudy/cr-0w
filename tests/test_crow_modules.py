"""crow_lib.modules and `crow modules`: adr_002 defaults, crow.yaml switches and config, errors.

The cases in fixtures/modules are shared with the Go loader (config.Modules)."""
import json
import os
import sys
from pathlib import Path

import pytest
import yaml
from conftest import BIN, run

sys.path.insert(0, str(BIN))
import crow_lib

ROOT = BIN.parent
CASES = Path(__file__).parent / "fixtures" / "modules"


def result(home):
    st = crow_lib.load_modules(str(home), str(ROOT))
    on = sorted(n for n, m in st["modules"].items() if m["on"])
    return {"on": on, "config": {n: st["modules"][n]["config"] for n in on},
            "errors": sorted(e["key"] for e in st["errors"])}


@pytest.mark.parametrize("case", sorted(p.name for p in CASES.iterdir() if p.is_dir()))
def test_shared_fixture_cases(case):
    want = json.loads((CASES / case / "want.json").read_text())
    assert result(CASES / case / "home") == want


def test_modules_returns_the_config_of_on_bricks(tmp_path):
    on = crow_lib.modules(str(tmp_path), str(ROOT))
    assert "design" not in on and "integration-shell" not in on
    assert on["nest"] == {"core_url": "http://127.0.0.1:8787"}


def test_the_example_crow_yaml_loads_clean_and_matches_the_defaults(tmp_path):
    text = (ROOT / "docs" / "crow.example.yaml").read_text()
    (tmp_path / "crow.yaml").write_text(text)
    st = crow_lib.load_modules(str(tmp_path), str(ROOT))
    assert st["errors"] == []
    assert all(m["on"] == m["default"] for m in st["modules"].values())
    assert set(st["modules"]) - {"core"} == set(yaml.safe_load(text)["modules"])


def test_errors_name_the_file_and_key():
    home = CASES / "errors" / "home"
    st = crow_lib.load_modules(str(home), str(ROOT))
    by_key = {e["key"]: e for e in st["errors"]}
    assert by_key["modules.billing"]["file"] == os.path.join(str(home), "crow.yaml")
    assert "unknown brick" in by_key["modules.billing"]["error"]
    assert "requires core-api" in by_key["modules.nest"]["error"]
    assert "core is always on" in by_key["modules.core"]["error"]


def test_crow_modules_json_lists_the_defaults_without_crow_yaml(hub):
    out = run(hub, "crow", "modules", "json")
    assert out.returncode == 0, out.stderr
    r = json.loads(out.stdout)
    assert r["exists"] is False and r["errors"] == []
    assert r["file"].endswith(os.path.join("local", "crow.yaml"))
    state = {m["name"]: m["on"] for m in r["modules"]}
    assert state["core"] and state["git"] and not state["design"] and not state["integration-shell"]


def test_crow_modules_text_marks_a_switched_brick(hub):
    home, _ = hub
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  design: true\n  billing: true\n")
    out = run(hub, "crow", "modules")
    assert out.returncode == 0
    assert "(none, defaults)" not in out.stdout
    line = next(l for l in out.stdout.splitlines() if " design " in l)
    assert line.startswith("on ") and line.endswith("(default off)")
    assert "modules.billing: unknown brick" in out.stderr


def test_crow_modules_rejects_other_arguments(hub):
    out = run(hub, "crow", "modules", "yaml")
    assert out.returncode != 0
    assert "crow-modules json" in out.stderr


def test_require_exits_2_with_the_message_when_the_brick_is_off(hub):
    home, _ = hub
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  vault: false\n")
    out = run(hub, "crow-modules", "require", "scripts", "crow-notes")
    assert out.returncode == 2 and out.stderr.strip() == "crow: module vault is off (crow.yaml)"
    assert run(hub, "crow-modules", "require", "scripts", "crow-activity").returncode == 0
    assert run(hub, "crow-modules", "require", "scripts", "not-a-piece").returncode == 0


def test_crow_dispatch_refuses_an_off_subcommand(hub):
    home, _ = hub
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  vault: false\n")
    out = run(hub, "crow", "notes")
    assert out.returncode == 2 and "crow: module vault is off (crow.yaml)" in out.stderr


def hub_skills(home):
    d = home / ".claude" / "skills"
    return sorted(p.name for p in d.iterdir()) if d.is_dir() else []


def test_link_generates_hub_and_profile_skills_and_follows_the_switch(hub):
    home, profiles = hub
    priv = home / "local" / "hub-skills" / "daily"
    priv.mkdir(parents=True)
    (priv / "SKILL.md").write_text("x")
    claude = profiles / "main" / "claude-home" / ".claude"
    (claude / "skills" / "mine").mkdir(parents=True)  # a real dir the owner put there stays
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    assert hub_skills(home) == ["afk", "brief", "daily", "eod", "night", "retro"]
    assert (home / ".claude" / "skills" / "brief" / "SKILL.md").is_file()
    prof = sorted(p.name for p in (claude / "skills").iterdir())
    assert prof == ["mine", "mr", "review", "review-threads", "wrap"]

    (home / "crow.yaml").write_text("version: 1\nmodules:\n  afk: false\n  forge: false\n")
    assert run(hub, "crow-modules", "link", "main").returncode == 0
    assert hub_skills(home) == ["brief", "daily", "eod", "retro"]
    assert sorted(p.name for p in (claude / "skills").iterdir()) == ["mine", "wrap"]


def test_hooks_of_an_off_brick_are_no_ops(hub, tmp_path):
    home, _ = hub
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  design: false\n  git: false\n")
    out = run(hub, "crow-design-notes", "hook", input='{"transcript_path": "/nope", "cwd": "/nope"}')
    assert out.returncode == 0 and out.stdout == ""
    assert run(hub, "crow-commit-msg-check", "-", input="not the format\n").returncode == 0
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  design: true\n")
    assert run(hub, "crow-commit-msg-check", "-", input="not the format\n").returncode == 1


def on_set(hub):
    out = run(hub, "crow-modules", "json")
    return sorted(m["name"] for m in json.loads(out.stdout)["modules"] if m["on"])


def test_off_then_on_returns_the_same_installed_set(hub):
    home, _ = hub
    assert run(hub, "crow-modules", "link").returncode == 0
    before = (on_set(hub), hub_skills(home))
    off = run(hub, "crow-modules", "off", "retro")
    assert off.returncode == 0, off.stderr
    assert "retro" not in on_set(hub) and "retro" not in hub_skills(home)
    on = run(hub, "crow-modules", "on", "retro")
    assert on.returncode == 0, on.stderr
    assert (on_set(hub), hub_skills(home)) == before


def test_switch_keeps_config_and_comments(hub):
    home, _ = hub
    (home / "local" / "crow.yaml").write_text(
        "version: 1\nmodules:\n  forge:   # mrs\n    on: true\n    hosts: {gl: gitlab.example}\n  design: false  # tools\n")
    assert run(hub, "crow-modules", "off", "forge").returncode == 0
    assert run(hub, "crow-modules", "on", "design").returncode == 0
    text = (home / "local" / "crow.yaml").read_text()
    assert "    on: false\n    hosts: {gl: gitlab.example}" in text and "design: true  # tools" in text
    st = crow_lib.load_modules(str(home), str(ROOT))
    assert st["errors"] == [] and not st["modules"]["forge"]["on"] and st["modules"]["design"]["on"]


@pytest.mark.parametrize("args,msg", [(["off", "core"], "always on"), (["on", "bogus"], "unknown brick"),
                                      (["off", "core-api"], "nest requires core-api")])
def test_switch_refuses_core_unknown_and_broken_requires(hub, args, msg):
    out = run(hub, "crow-modules", *args)
    assert out.returncode != 0 and msg in out.stderr
    assert not (hub[0] / "local" / "crow.yaml").exists()


@pytest.mark.parametrize("script", ["crow-retro", "crow-releases", "crow-afk-preflight", "crow-compare"])
def test_scripts_of_an_off_brick_refuse(hub, script):
    home, _ = hub
    (home / "local" / "crow.yaml").write_text("version: 1\nmodules:\n  retro: false\n  releases: false\n  afk: false\n")
    out = run(hub, script, "--help")
    assert out.returncode == 2 and "is off (crow.yaml)" in out.stderr
