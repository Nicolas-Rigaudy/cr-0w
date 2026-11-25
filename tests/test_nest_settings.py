"""`crow modules get|set nest` (item_282): the Nest settings writer core calls for /api/settings/nest."""
import difflib
import json
import re
import shutil
import sys
from importlib.machinery import SourceFileLoader

import pytest
from conftest import BIN, run

sys.path.insert(0, str(BIN))

ROOT = BIN.parent
GARDEN = ROOT / "tools" / "nest" / "fixtures" / "bricks" / "garden"

COMMENTED = """\
# crow.yaml of a test hub
version: 1
modules:
  forge:   # merge requests
    on: true
    hosts: {gl: gitlab.example}
  nest:
    core_url: http://127.0.0.1:8787  # the core
    theme: modern   # first paint
    pages:
      - home
      - desk
  mail: false  # no mail here

# trailing note
"""


def garden(home):
    """A user brick with a theme (moss), a hero (garden) and a page (garden)."""
    brick = home / "modules" / "garden"
    shutil.copytree(GARDEN, brick)
    (brick / "module.yaml").write_text("summary: Garden pieces for Nest.\n")
    return brick


def get(hub):
    out = run(hub, "crow-modules", "get", "nest")
    assert out.returncode == 0, out.stderr
    return json.loads(out.stdout)


def changed_lines(before, after):
    return [line for line in difflib.unified_diff(before.splitlines(), after.splitlines(), lineterm="", n=0)
            if line[:1] in "+-" and not line.startswith(("+++", "---"))]


def test_get_lists_the_defaults_and_the_names_of_on_bricks(hub):
    home, _ = hub
    garden(home)
    got = get(hub)
    assert (got["theme"], got["hero"], got["pages"]) == ("modern", "holotable", [])
    assert got["known"]["themes"] == ["modern", "cyber", "moss"]
    assert got["known"]["heroes"] == ["holotable", "none", "garden"]
    assert got["known"]["pages"][-1] == "garden" and "home" in got["known"]["pages"]
    assert got["errors"] == []
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  garden: false\n")
    assert get(hub)["known"]["themes"] == ["modern", "cyber"]


def test_set_touches_only_the_nest_lines_and_keeps_comments(hub):
    home, _ = hub
    garden(home)
    (home / "crow.yaml").write_text(COMMENTED)
    out = run(hub, "crow-modules", "set", "nest", "theme", "moss", "hero", "garden", "pages", "home,garden,home", "--json")
    assert out.returncode == 0, out.stderr
    got = json.loads(out.stdout)
    assert (got["theme"], got["hero"], got["pages"]) == ("moss", "garden", ["home", "garden"])
    assert got["changed"] == ["theme", "hero", "pages"] and got["needs_build"] is True
    after = (home / "crow.yaml").read_text()
    assert changed_lines(COMMENTED, after) == [
        "-    theme: modern   # first paint", "-    pages:", "-      - home", "-      - desk",
        "+    theme: moss  # first paint", "+    pages: [home, garden]", "+    hero: garden"]
    assert get(hub)["errors"] == []


def test_a_theme_only_save_needs_no_build(hub):
    home, _ = hub
    (home / "crow.yaml").write_text(COMMENTED)
    got = json.loads(run(hub, "crow-modules", "set", "nest", "theme", "cyber", "--json").stdout)
    assert got["changed"] == ["theme"] and got["needs_build"] is False
    same = json.loads(run(hub, "crow-modules", "set", "nest", "pages", "home,desk", "--json").stdout)
    assert same["changed"] == [] and same["needs_build"] is False


@pytest.mark.parametrize("args,bad", [(["theme", "neon"], "'neon'"), (["hero", "garden"], "'garden'"),
                                      (["pages", "home,attic"], "'attic'"), (["colour", "red"], "'colour'"),
                                      (["theme", "cyber", "hero", "nope"], "'nope'"), (["theme"], "pairs")])
def test_an_unknown_value_exits_2_names_it_and_writes_nothing(hub, args, bad):
    home, _ = hub
    (home / "crow.yaml").write_text(COMMENTED)
    out = run(hub, "crow-modules", "set", "nest", *args)
    assert out.returncode == 2 and bad in out.stderr
    assert (home / "crow.yaml").read_text() == COMMENTED


def test_set_adds_the_nest_mapping_when_missing_and_keeps_a_switch(hub):
    home, _ = hub
    assert run(hub, "crow-modules", "set", "nest", "hero", "none").returncode == 0
    assert (home / "local" / "crow.yaml").read_text() == "version: 1\nmodules:\n  nest:\n    hero: none\n"
    (home / "crow.yaml").write_text("version: 1\nmodules:\n  nest: true  # web\n")
    out = run(hub, "crow-modules", "set", "nest", "pages", "")
    assert out.returncode == 0, out.stderr
    assert (home / "crow.yaml").read_text() == "version: 1\nmodules:\n  nest:  # web\n    on: true\n    pages: []\n"
    assert "every page" in out.stdout


def test_a_flow_mapping_is_refused_unchanged(hub):
    home, _ = hub
    text = "version: 1\nmodules:\n  nest: {theme: cyber}\n"
    (home / "crow.yaml").write_text(text)
    out = run(hub, "crow-modules", "set", "nest", "theme", "modern")
    assert out.returncode == 1 and "block style" in out.stderr
    assert (home / "crow.yaml").read_text() == text


def test_only_nest_settings_can_be_set(hub):
    out = run(hub, "crow-modules", "set", "forge", "hosts", "x")
    assert out.returncode == 2 and "only nest" in out.stderr


def test_engine_names_match_nest_modules_ts():
    ts = (ROOT / "tools" / "nest" / "nest-modules.ts").read_text()
    mods = SourceFileLoader("crow_modules", str(BIN / "crow-modules")).load_module()
    for const, py in (("ENGINE_THEMES", "NEST_THEMES"), ("ENGINE_HEROES", "NEST_HEROES"), ("ENGINE_PAGES", "NEST_PAGES")):
        names = re.search(const + r" = \[([^\]]*)\]", ts).group(1)
        assert [n.strip(" '") for n in names.split(",")] == getattr(mods, py), const
