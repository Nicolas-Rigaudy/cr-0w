## item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors - Core rebuilds Nest after a settings save and reports build errors
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 80%
> Progress: 100%
> Complexity: Medium
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 16:44:21

# AI Context
- Summary: After a settings save that needs a build, core runs `crow nest build`, streams its state, and Nest reloads on success or shows the error with the brick name.
- Keywords: core, rebuild, crow nest build, settings, build error, brick, reload
- Use when: wiring the rebuild step behind the Nest Settings save.
- Skip when: the change only writes crow.yaml (item_282) or is page layout (item_284).

# Problem
Brick themes, heroes and pages are compiled in. A hero or page change in crow.yaml shows nothing until `crow nest build` runs, and today only the terminal can run it.

# Scope
- In:
  - Core starts `crow nest build` after a save that returned `needs_build`; one build at a time, a second save while building queues one more build.
  - A build state the page can poll or receive live: `idle | building | ok | failed`, elapsed time, last lines of the log.
  - On failure: crow.yaml keeps the new values, the last good dist stays served, the error names the brick (the `brick <name>:` prefix from `nest-modules.ts` errors).
  - On success: Nest shows a "new build ready · reload" chip (nest_build moves in /api/modules); tap reloads, else the next route change does. Never a reload mid-edit (owner, 2026-11-25).
  - A module switch that changes the on bricks with Nest code starts the same build (owner, 2026-11-25).
- Out:
  - Rebuilding the owner's live dist during tests: tests use a temp dist dir and a fake build command.
  - Rolling crow.yaml back on failure (the owner fixes the brick or changes the value again).

# Acceptance criteria
- AC3: After a save that needs a build, core runs `crow nest build` and Nest reloads on success or shows the build error with the brick name.
- AC3a: A test with a fake build command covers ok, failed (brick name in the error) and a save during a build.
- AC3b: A theme-only save does not start a build (theme already switches live through data-theme).
- AC3c: `go test ./...` in tools/crow-core and `bin/crow-check` pass.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC3: After a save that needs a build, core runs `crow nest build` and Nest reloads on success or shows the build error with the brick name.
- AC3 -> 2abcb55 (core build, brick name in the failed state), e847c6f (Nest reload chip), eaf4c22 (module switch builds). Proof: TestNestBuildOK, TestNestBuildFailedKeepsDist, TestNestBuildOnModuleSwitch; Nest tests modules.test.ts, Shell.test.tsx, router.test.ts.
- AC3a -> Proof: TestNestBuildOK, TestNestBuildFailedKeepsDist, TestNestBuildQueuesOneMore (fake build command).
- AC3b -> Proof: TestNestBuildOK asserts a theme-only save runs no build.
- AC3c -> Proof: `go test ./...` in tools/crow-core passed (316 tests); `bin/crow-check` green except tests/test_module_inventory.py (worktree without .claude/skills).

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`
- Primary task(s): `task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`

# Priority
- Priority: Medium
- Rationale: Order 2 of 4, after item_282. Draft until the owner picks a Settings layout (the build state panel differs per option).

# Notes
- Risk: build time on the owner's machine (~40 s in the mockup); the page must stay usable while it runs.
- Mockups: options A, B, C in `tools/nest/scripts/mockups/settings.html` show the save bar, a running build (B) and a failed build (C).
- Derived from request `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`.
- Task `task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors` was finished via `logics-manager flow finish task` on 2026-11-25.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T15:44:59.085881+00:00 | cdx | 3814751952664ffb83222fa44cd7a61c | Settings core closeout (item_283) | succeeded |  |  |
| 2026-11-25T15:36:54.813214+00:00 | cdx | a6f35983988146baacef59b5fbf49d4c | Nest Settings core 2 (item_283) | cancelled |  | stopped |
| 2026-11-25T15:36:31.719471+00:00 | cdx | 86a703d4b67947d79562ca0878c3c400 | Nest Settings core 2 (item_283) | failed |  | limit |

# Tasks
- `task_032_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors`
