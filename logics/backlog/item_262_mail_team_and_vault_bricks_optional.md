## item_262_mail_team_and_vault_bricks_optional - Mail, team and vault bricks optional
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90
> Confidence: 90
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-23 16:40:52

# AI Context
- Summary: Split step 5 of adr_002: mail (calendar and mail counts, Today box), team (management checks) and vault (notes, daily log, open threads) each turn off with one switch; with vault off, notes go to CROW_HOME.
- Keywords: modules, mail, team, vault, notes, today, checks
- Use when: making mail, team or vault optional.
- Skip when: forge (item_236) or tool bricks (item_263).

# Problem
- A setup with no mail account, no team rituals or no vault gets errors from `crow-today`, `crow-checks` and the vault scripts.

# Scope
- In:
  - mail off: no `today` section or routes, no Today box, brief without calendar and mail.
  - team off: no `checks` section, no management rows in brief or deck.
  - vault off: `crow-note`, `crow-log` and open threads write Markdown under `CROW_HOME/notes/`; `crow-state path` points there; `eod` skill not installed.
  - Doctor quiet for off bricks.
- Out:
  - Moving `vault.path` from `projects.yaml` under `modules.vault` (later, if the owner wants it).

# Acceptance criteria
- AC1: With mail, team and vault off, brief, gather, doctor and Nest run with no error.
- AC2: With vault off, `crow note` and open threads round-trip through `CROW_HOME/notes/`.
- AC3: With all three on, nothing changes (owner's setup).

# AC Traceability
- req_064_paul_s_feedback_and_a_modular_public_crow-AC3 -> This backlog slice. Proof: brief, gather, doctor and Nest run with mail, team and vault off.
- AC1 -> `tests/test_crow_gather.py::test_all_optional_bricks_off_still_gather` (brief, daily, weekly), `tests/test_crow_notes_mode.py::test_doctor_is_quiet_for_off_bricks`; core `TestModuleOffNoRefresh` (no script run for an off section); Nest `modules.test.ts` (today and checks off, never fetched); brief skill step 8 skips mail.
- AC2 -> `tests/test_crow_notes_mode.py::test_note_and_log_round_trip_through_crow_home_notes`, `test_threads_path_points_to_crow_home_notes`, `tests/test_crow_gather.py::test_weekly_with_vault_off_reads_threads_from_crow_home_notes`; eod not linked (`test_eod_skill_is_not_linked_with_vault_off`).
- AC3 -> `tests/test_crow_notes_mode.py::test_all_on_keeps_the_vault_daily`, `tests/test_crow_gather.py::test_without_crow_yaml_every_default_key_is_there`; full `bin/crow-check` green (commits bbf208c, 16249c9, 60c4d31).

# Dependencies
- `item_261_skills_hooks_and_gather_follow_the_module_switch`, `item_237_nest_and_crow_core_as_separate_front_and_back_bricks`.

# Links
- Request: `req_064_paul_s_feedback_and_a_modular_public_crow`
- Architecture decision(s): `adr_002_crow_bricks_module_manifest_and_project_kinds`
- Next: `item_263_optional_tool_bricks_and_install_profiles`

# Priority
- Priority: Medium
- Rationale: Needed by a second user with no work mail or vault.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
