## item_283_core_rebuilds_nest_after_a_settings_save_and_reports_build_errors - Core rebuilds Nest after a settings save and reports build errors
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

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
  - On success: Nest reloads itself (version bump in the live feed).
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
- Request: `logics/request/req_066_nest_settings_page_theme_hero_page_order_and_module_switches.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Order 2 of 4, after item_282. Draft until the owner picks a Settings layout (the build state panel differs per option).

# Notes
- Risk: build time on the owner's machine (~40 s in the mockup); the page must stay usable while it runs.
- Mockups: options A, B, C in `tools/nest/scripts/mockups/settings.html` show the save bar, a running build (B) and a failed build (C).
- Derived from request `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

