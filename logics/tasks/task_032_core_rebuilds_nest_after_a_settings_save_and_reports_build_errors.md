## task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors - Core rebuilds Nest after a settings save and reports build errors
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-25 16:12:38

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: core, rebuilds, nest, after, settings, save, reports, build, errors
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`

# Acceptance criteria
- AC3: After a save that needs a build, core runs `crow nest build` and Nest reloads on success or shows the build error with the brick name.
- AC3a: A test with a fake build command covers ok, failed (brick name in the error) and a save during a build.
- AC3b: A theme-only save does not start a build (theme already switches live through data-theme).
- AC3c: `go test ./...` in tools/crow-core and `bin/crow-check` pass.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors.md` after implementation.

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
