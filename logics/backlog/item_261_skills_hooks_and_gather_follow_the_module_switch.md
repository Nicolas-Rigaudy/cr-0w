## item_261_skills_hooks_and_gather_follow_the_module_switch - Skills, hooks and gather follow the module switch
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Split step 2 of adr_002: an off brick installs no skills or git hooks, its Claude hooks exit at once, gather skips its keys and its `crow` subcommands refuse with a clear line.
- Keywords: modules, skills, hooks, gather, install, switch
- Use when: making the generic off behaviour work for any brick.
- Skip when: a single brick's own code paths (items 236, 262, 263).

# Problem
- Skills, hooks and gather keys are installed and run for every brick, whatever the setup.

# Scope
- In:
  - `crow-memory link` links only skills of on bricks; doctor's skill link check follows the on set.
  - Hub skills move from `.claude/skills/` to `skills/`; install links on bricks into `.claude/skills/` (git-ignored).
  - Claude hook scripts exit 0 when their brick is off.
  - Git hooks installed only in `repo` projects while `git` is on.
  - `crow-gather` collects only the gather keys of on bricks and lists the on set in its JSON; brief, daily and weekly skills read it.
  - `crow` subcommands of an off brick print `crow: module <name> is off (crow.yaml)` and exit 2.
- Out:
  - Brick-specific code paths.

# Acceptance criteria
- AC1: Turning `design` off removes its skills and makes its hook a no-op, with no other change.
- AC2: `crow-gather brief` with `afk` off has no afk key and the brief renders.
- AC3: With no crow.yaml, installed skills and hooks match today's set.

# Delivery
- `crow modules link [account...]` (crow_lib `link_hub_skills`, `link_profile_skills`): the hub `.claude/skills/` is generated and git-ignored; it holds the `hub_skills` of on bricks (`skills/`) plus `local/hub-skills/*`. Profiles get the `skills` of on bricks. Only links into `skills/` or `local/hub-skills/` are removed; real dirs and foreign links stay. `crow hub` and `crow-memory link` run it.
- New manifest field `hub_skills`: skills linked only in the hub (brief, afk, night, eod, retro). `skills` stays profile-wide (wrap, mr, review, review-threads).
- Design skills ship as the crow-design plugin: with design off, `crow modules link` runs `crow-skills-ext unlink` and `crow-skills-ext link` refuses.
- `crow-modules require <field> <name>`: exit 2 with `crow: module <name> is off (crow.yaml)`. `bin/crow` guards activity, note(s), projects/releases/release, focus, usage, secret, deck/watch, core/mission/inbox through it.
- Claude hooks (crow-hub-guard, crow-runs inbox --hook, crow-relaunch hook, crow-design-notes hook) call `crow_lib.exit_if_off`. Git hooks (commit-msg, pre-commit, pre-push) exit 0 when git is off. Crow has no git hook installer, so "installed only while git is on" became "a no-op while git is off".
- `crow-gather` drops the keys off bricks list under `gather.<routine>` (jobs never start) and adds `modules`. The brief, daily and weekly skills skip sections whose key is absent.
- Doctor: profile skill check expects the `skills` of on bricks and flags links to off ones.
- Owner setup: `local/crow.yaml` with `design: true` keeps today's set (engine default stays design off).
- Inventory: tagged 6 pieces left untagged by the routines commit (crow-routine, butler clear and routines routes in core/core-api, afk routes in afk).

# AC Traceability
- req_064_paul_s_feedback_and_a_modular_public_crow-AC3 -> This backlog slice. Proof: gather and install follow the on set.
- AC1 -> `tests/test_crow_modules.py::test_hooks_of_an_off_brick_are_no_ops` (design hook prints nothing, exit 0), `test_link_generates_hub_and_profile_skills_and_follows_the_switch` (off brick links removed); design plugin removed by `crow modules link` (`crow-skills-ext unlink`).
- AC2 -> `tests/test_crow_gather.py::test_brief_drops_the_keys_of_an_off_brick`; brief skill skips absent keys.
- AC3 -> `tests/test_crow_gather.py::test_without_crow_yaml_every_default_key_is_there`, `test_crow_doctor.py::test_doctor_passes_a_set_up_profile`; live run: the hub and existing profile skill sets were unchanged.

# Dependencies
- `item_260_module_manifests_crow_yaml_loader_and_doctor_check`.

# Links
- Request: `req_064_paul_s_feedback_and_a_modular_public_crow`
- Architecture decision(s): `adr_002_crow_bricks_module_manifest_and_project_kinds`
- Next: `item_236_first_module_split_forge_mrs_optional_and_folder_projects`

# Priority
- Priority: Medium
- Rationale: Gives every brick its off switch at once.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-23T13:50:04.532504+00:00 | hub-session | 2026-11-23-hub | Skills, hooks and gather follow the module switch | succeeded |  |  |
