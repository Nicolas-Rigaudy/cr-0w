import datetime as dt
import json

from conftest import run

TODAY = dt.date.today()
OLD = TODAY - dt.timedelta(days=20)
MONDAY = TODAY - dt.timedelta(days=TODAY.weekday())

REGISTRY = """\
vault: {path: null, crow_dir: null}
projects:
  - {name: api, aliases: [backend], scope: work, path: /x}
  - {name: toy, scope: personal, path: /y}
"""

THREADS = f"""\
# Open threads

## Me
- [ ] {TODAY} [backend] fix the seed
- [ ] {OLD} [toy] polish the readme
- [x] {TODAY} [work] done already

## Crow
- [ ] {TODAY} [hub] crow item

## Waiting
- [ ] {OLD} [work] vendor answer
"""


def setup(hub, extra_env=None):
    home, _ = hub
    state = home / "local" / "state"
    (home / "local" / "projects.yaml").write_text(REGISTRY)
    (home / "local" / "context.md").write_text(
        "- Confluence daily page: https://x/pages/1/Daily (space dt, page id 431)\n"
        "- Atlassian cloudId: 1d9177cb-ebdb-4a95-a1af-54c78e01130c (Jira project CD).\n")
    (state / "open-threads.md").write_text(THREADS)
    (state / "activity").mkdir()
    (state / "activity" / f"{TODAY}.md").write_text("- 09:00 [hub] brief\n")
    (state / "activity" / f"{OLD}.md").write_text("- 09:00 [hub] old\n")
    (state / "decisions.md").write_text(f"# Decisions\n\n## {OLD} [hub] old one\nx\n\n## {TODAY} [hub] new one\ny\n")
    (state / "last-daily").write_text(f"{TODAY - dt.timedelta(days=2)}\n")
    env = {"CROW_GATHER_OFFLINE": "1", "CROW_VAULT": str(home / "novault")}
    env.update(extra_env or {})
    return env


def gather(hub, *args, env=None):
    out = run(hub, "crow-gather", *args, "--json", env=env)
    assert out.returncode == 0, out.stderr
    return json.loads(out.stdout)


def test_brief_keeps_me_and_waiting_with_stale_flag(hub):
    data = gather(hub, "brief", env=setup(hub))
    threads = data["open_threads"]
    assert threads["Me"] == [f"0d [backend] fix the seed", f"STALE 20d [toy] polish the readme"]
    assert threads["Waiting"] == ["STALE 20d [work] vendor answer"]
    assert threads["_skipped_counts"] == {"Crow": 1}
    assert data["activity"] == {"skipped": "offline"}
    assert data["checks"] == {"skipped": "offline"}
    assert data["targets"]["activity_log"].endswith(f"{TODAY}.md")


def test_daily_since_last_daily_and_work_threads_only(hub):
    data = gather(hub, "daily", env=setup(hub))
    assert data["since"] == str(TODAY - dt.timedelta(days=1))
    assert data["tracker"] == {"cloud_id": "1d9177cb-ebdb-4a95-a1af-54c78e01130c",
                               "jira_project": "CD", "daily_page_id": "431"}
    assert data["open_threads_work"]["Me"] == ["0d [backend] fix the seed"]
    assert data["open_threads_work"]["Waiting"] == ["STALE 20d [work] vendor answer"]
    assert list(data["activity_logs"]) == [str(TODAY)]


def test_weekly_week_bounds_and_decisions(hub):
    data = gather(hub, "weekly", env=setup(hub))
    assert data["monday"] == str(MONDAY)
    assert data["iso_week"] == "{}-W{:02d}".format(*TODAY.isocalendar()[:2])
    assert len(data["decisions"]) == 1 and "new one" in data["decisions"][0]
    assert data["targets"]["rollup"].endswith(data["iso_week"] + ".md")


def test_output_fits_the_budget(hub):
    home, _ = hub
    env = setup(hub, {"CROW_GATHER_BUDGET": "3000"})
    (home / "local" / "state" / "activity" / f"{TODAY}.md").write_text("- line\n" * 2000)
    out = run(hub, "crow-gather", "daily", "--json", env=env)
    assert len(out.stdout) <= 3100
    assert json.loads(out.stdout)["_cut_to_fit"] == [f"activity_logs/{TODAY}"]


def test_vault_notes_when_crow_dir_set(hub):
    home, _ = hub
    env = setup(hub)
    notes = home / "vault" / "CR-0W"
    (notes / "Decisions").mkdir(parents=True)
    (home / "local" / "projects.yaml").write_text(
        REGISTRY.replace("{path: null, crow_dir: null}", f"{{path: {home / 'vault'}, crow_dir: CR-0W}}"))
    (notes / "Open Threads.md").write_text(THREADS.replace("fix the seed", "vault seed"))
    (notes / "Status.md").write_text("vault status\n")
    (notes / "Decisions" / f"{TODAY} new one.md").write_text("y\n")
    (notes / "Decisions" / f"{OLD} old one.md").write_text("x\n")
    brief = gather(hub, "brief", env=env)
    assert brief["open_threads"]["Me"][0] == "0d [backend] vault seed"
    assert brief["status"].startswith("vault status")
    assert brief["state_dir"] == str(home / "local" / "state")
    weekly = gather(hub, "weekly", env=env)
    assert [d.splitlines()[0] for d in weekly["decisions"]] == [f"## {TODAY} new one"]
    assert weekly["targets"]["rollup"].startswith(str(notes / "Rollups"))
