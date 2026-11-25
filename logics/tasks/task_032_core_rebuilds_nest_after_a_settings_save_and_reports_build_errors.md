## task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors - Core rebuilds Nest after a settings save and reports build errors
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
> Indicators reviewed: 2026-11-25 16:44:21

# AI Context
- Summary: core runs `crow nest build` after a settings save or a module switch that needs it, serves the new dist only when the build is good, and Nest offers a reload chip.
- Keywords: core, rebuilds, nest, after, settings, save, reports, build, errors
- Use when: changing the Nest rebuild after a save or a switch, or how Nest picks up a new build.
- Skip when: writing crow.yaml (item_282), the Settings page (item_284) or the switches UI (item_285).

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`

# Acceptance criteria
- AC3: After a save that needs a build, core runs `crow nest build` and Nest reloads on success or shows the build error with the brick name.
- AC3a: A test with a fake build command covers ok, failed (brick name in the error) and a save during a build.
- AC3b: A theme-only save does not start a build (theme already switches live through data-theme).
- AC3c: `go test ./...` in tools/crow-core and `bin/crow-check` pass.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors.md` after implementation.

# Validation
- `npx tsc --noEmit -p .` in tools/nest: passed, no errors.
- `npx vitest run` in tools/nest: passed (44 files, 342 tests).
- `go test ./...` in tools/crow-core: passed (316 tests, 14 packages).
- Known: tests/test_module_inventory.py fails in worktrees without the generated .claude/skills folder; unrelated, not fixed.
- Finish workflow executed on 2026-11-25.
- Linked backlog/request close verification passed.

# Report
- 2abcb55: core starts `crow nest build` after a save with needs_build, one at a time with one queued; builds into a staging dir and swaps it in only on success; GET /api/settings/nest/build answers state, elapsed, log, failing brick; /api/modules carries nest_build.
- e847c6f: Nest flags a new build from nest_build and shows a "new build ready · reload" chip; tap reloads, the next route change loads the new build. No reload mid-edit (owner decision 2026-11-25).
- eaf4c22: a module switch that changes the on bricks with Nest code starts the same build (item_285 path).
- Known: tests/test_module_inventory.py fails in worktrees without the generated .claude/skills folder; not related, not fixed here.
- Finished on 2026-11-25.
- Linked backlog item(s): `item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`
- Related request(s): `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`

# Links
- Request: `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# AC Traceability
- AC3 -> 2abcb55, e847c6f. Proof: core runs the build after a save that needs it; a failed build names the brick in the build state; Nest shows the reload chip when nest_build moves (src/data/modules.test.ts, src/ui/Shell.test.tsx, src/router.test.ts).
- AC3a -> 2abcb55. Proof: TestNestBuildOK, TestNestBuildFailedKeepsDist, TestNestBuildQueuesOneMore use a fake build command; TestNestBuildOnModuleSwitch (eaf4c22) covers switches.
- AC3b -> 2abcb55. Proof: TestNestBuildOK saves a theme only and asserts no build ran.
- AC3c -> Proof: `go test ./...` in tools/crow-core passed (316 tests); `bin/crow-check` green except tests/test_module_inventory.py (worktree without .claude/skills).
- request-AC3 -> This task. Proof deferred to slice closeout.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
