## item_165_telemetry_strip_on_the_laptop_home - Telemetry strip on the laptop Home
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 18:01:00

# AI Context
- Summary: a thin telemetry strip on the laptop Home (tokens, cost, model, uptime) fed by the core accounts and runs sections.
- Keywords: nest, home, telemetry strip, tokens, cost, model, uptime
- Use when: building the Home telemetry strip.
- Skip when: new core data sources.

# Problem
- Home does not show session telemetry the way the darwin reference does.

# Scope
- In:
  - A thin strip with tokens, cost, model and uptime, read through `useSnapshot()` from the core `usage`, `accounts` and `sessions` sections (item_169; fixtures in dev and tests).
  - // section labels if they fit the boards.
- Out:
  - New data sources in core (item_169 serves them)

# Acceptance criteria
- AC1: No hard-coded sample values in the component.
- AC2: The owner approves the render.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1: No hard-coded sample values in the component.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_007_nest_home_design_pass`
- Architecture decision(s): (none yet)
- Request: `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
- Depends on: `item_167_nest_reads_live_data_from_the_core_api` (live data; can be built on fixtures if it slips), `item_169_core_serves_every_deck_source_as_read_only_sections` (usage and sessions sections)
- Primary task(s): `task_019_orchestrate_the_nest_home_design_pass`, run in `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: Low
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
