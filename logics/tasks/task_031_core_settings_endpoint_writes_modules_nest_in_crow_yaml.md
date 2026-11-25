## task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml - Core settings endpoint writes modules.nest in crow.yaml
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-25 16:03:26

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: core, settings, endpoint, writes, modules, nest, crow, yaml
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_282_core_settings_endpoint_writes_modules_nest_in_crow_yaml`

# Acceptance criteria
- AC2: Saving writes `modules.nest.*` in crow.yaml through core, keeps the rest of the file and its comments, and refuses unknown values.
- AC2a: A test edits a crow.yaml with comments and other bricks; the diff touches only the `modules.nest` lines.
- AC2b: An unknown theme, hero or page returns 400 with the bad value named; crow.yaml is unchanged.
- AC2c: `go test ./...` in tools/crow-core, the crow-modules tests and `bin/crow-check` pass.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- Not started.

# Links
- Request: `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
