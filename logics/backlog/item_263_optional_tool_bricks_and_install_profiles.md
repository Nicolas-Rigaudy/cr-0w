## item_263_optional_tool_bricks_and_install_profiles - Optional tool bricks and install profiles
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-23 17:14:51

# AI Context
- Summary: Split step 6 of adr_002: afk, retro, releases, design and the shell integration turn off cleanly, and `install.sh` asks which bricks to turn on and writes the starter crow.yaml.
- Keywords: modules, afk, retro, releases, design, install, profile
- Use when: finishing the module split and the public install path.
- Skip when: earlier splits.

# Problem
- A new install gets every brick and edits crow.yaml by hand.

# Scope
- In:
  - afk, retro, releases, design and integration-shell off: no skills, no hooks, no subcommands, no sections.
  - `install.sh --modules a,b,c` and an interactive prompt; default starter set from adr_002.
  - `crow modules on <brick>` and `off <brick>` edit crow.yaml and rerun the link step.
  - Install docs list each brick with one line.
- Out:
  - Separate packages per brick.

# Acceptance criteria
- AC1: A docker install with `--modules core` runs brief, doctor and Nest with no error.
- AC2: `crow modules off retro` then `on retro` returns the same installed set.
- AC3: With no crow.yaml, nothing changes (owner's setup).

# Delivery
- `crow modules on|off <brick>` (crow-modules `switch`): edits crow.yaml (created when missing; other keys, config and comments kept), then relinks. Refuses core, unknown bricks, an on brick whose requires are off, and an off brick an on brick requires.
- Tool brick scripts refuse when off (exit 2, `crow: module X is off (crow.yaml)`): Python ones through `crow_lib.refuse_if_off` (afk/night page, digest, cost audit, toolwatch, retro, changelog, release propose, releases, compare, crop, ansi-shot); bash ones (afk/night preflight, keep-awake) through `crow-modules require`. The deck `r` (retro) key flashes the off message.
- `install.sh --modules a,b,c` (or `CROW_MODULES`, or a /dev/tty prompt): a new data dir gets `local/crow.yaml` with core plus the listed bricks on, every other one off; default starter set `core,git,core-api,nest,deck`. Unknown names and broken requires abort before writing. An existing crow.yaml is kept; an existing data dir without one keeps every default. `fetch` reads `file://` without curl.
- `docs/install-bricks.md` lists each brick in one line; runbook install-and-update updated.

# AC Traceability
- AC1 -> docker (python:3.13-slim + git + PyYAML, no network): `install.sh --modules core` then `crow doctor` (modules ok, 1 of 15 on; cdx stubbed), `crow-gather brief --json` ok, `crow core serve` prints `module core-api is off` (Nest off cleanly). `tests/test_install_update.py::test_install_modules_core_turns_every_other_brick_off`.
- AC2 -> `tests/test_crow_modules.py::test_off_then_on_returns_the_same_installed_set`; docker run: off/on retro relinks then removes the retro hub skill.
- AC3 -> `tests/test_install_update.py::test_existing_data_dir_without_crow_yaml_keeps_every_default`; owner checkout (no crow.yaml change) still runs crow-retro and crow-releases.
- req_064_paul_s_feedback_and_a_modular_public_crow-AC1 -> This backlog slice. Proof: a second user installs with only the bricks they need, no local patches.

# Dependencies
- `item_262_mail_team_and_vault_bricks_optional`.

# Links
- Request: `req_064_paul_s_feedback_and_a_modular_public_crow`
- Architecture decision(s): `adr_002_crow_bricks_module_manifest_and_project_kinds`

# Priority
- Priority: Medium
- Rationale: Closes the public install path.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-23T16:39:53+00:00 | hub-session | 2026-11-23-hub | Optional tool bricks and install profiles | succeeded |  | 2945ac5, 32c5f0c |
