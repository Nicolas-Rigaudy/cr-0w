## task_015_build_crow_data_in_the_vault - Build crow data in the vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-09 10:46:09

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: build, crow, data, vault
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Context
- Orchestrate the scaffolded request chain and keep sibling implementation slices linked.

# Plan
- [x] 1. Path helper
- [x] 2. Activity lines in the daily note
- [x] 3. Move owner notes and set crow_dir
- [x] 4. Vault side hand-off (built by the vault session 2026-11-09)
- [x] ADR 009 checkpoint: update affected Logics docs during each meaningful wave and leave the repo commit-ready.
- [x] Keep commit creation under operator control; do not force one commit per micro-step.
- [x] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_127_path_helper_separating_vault_notes_from_machine_state`
- `item_128_activity_lines_in_the_daily_note`
- `item_129_move_owner_notes_into_the_vault_and_switch_on_crow_dir`
- `item_130_vault_side_dashboard_templates_and_folder_conventions`

# Definition of Done (DoD)
- [x] Generated request, product, backlog, and task docs are present.
- [x] Context-pack handoff is available when requested.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# AC Traceability
- request-AC1 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`
- request-AC2 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`
- request-AC3 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`
- request-AC1 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`
- request-AC4 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`
- request-AC3 -> This task. Proof: Implemented in cb902fc, c284dfd, 8b8939e, 0ed4e98; validated with python3 -m pytest -q tests (444 passed) and go test ./... Source: `0ed4e98`

# Validation
- (no validation recorded yet)
- Vault check passed on 2026-11-09: Daily Note template has ## CR-0W, Weekly Review embeds the rollup and excludes CR-0W, CR-0W.md and CR-0W.base exist, vault CLAUDE.md has the CR-0W section, .stignore has .*.crow-tmp
- Finish workflow executed on 2026-11-09.
- Linked backlog/request close verification passed.

# Report
- Not started.
- Finished on 2026-11-09.
- Linked backlog item(s): `item_127_path_helper_separating_vault_notes_from_machine_state`, `item_128_activity_lines_in_the_daily_note`, `item_129_move_owner_notes_into_the_vault_and_switch_on_crow_dir`, `item_130_vault_side_dashboard_templates_and_folder_conventions`
- Related request(s): `req_041_crow_data_in_the_vault`

# Links
- Request: `req_041_crow_data_in_the_vault`
- Product brief(s): `prod_006_crow_data_in_the_vault`
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
