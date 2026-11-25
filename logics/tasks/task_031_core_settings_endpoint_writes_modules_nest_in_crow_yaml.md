## task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml - Core settings endpoint writes modules.nest in crow.yaml
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
> Indicators reviewed: 2026-11-25 16:12:19

# AI Context
- Summary: crow modules get/set nest and core GET/POST /api/settings/nest write modules.nest in crow.yaml.
- Keywords: core, settings, endpoint, writes, modules, nest, crow, yaml
- Use when: changing how Nest theme, hero or pages are saved.
- Skip when: working on the build after a save (item_283) or the page (item_284).

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_282_core_settings_endpoint_writes_modules_nest_in_crow_yaml`

# Acceptance criteria
- AC2: Saving writes `modules.nest.*` in crow.yaml through core, keeps the rest of the file and its comments, and refuses unknown values.
- AC2a: A test edits a crow.yaml with comments and other bricks; the diff touches only the `modules.nest` lines.
- AC2b: An unknown theme, hero or page returns 400 with the bad value named; crow.yaml is unchanged.
- AC2c: `go test ./...` in tools/crow-core, the crow-modules tests and `bin/crow-check` pass.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_031_core_settings_endpoint_writes_modules_nest_in_crow_yaml.md` after implementation.

# Validation
- `go test ./...` in tools/crow-core: passed (311 tests).
- `python3 -m pytest tests/test_nest_settings.py tests/test_crow_modules.py tests/test_user_modules.py`: passed (47 tests).
- `logics-manager review check`: 0 hits.
- command: `go test ./... (tools/crow-core); python3 -m pytest tests/test_nest_settings.py tests/test_crow_modules.py tests/test_user_modules.py` | result: passed | date: 2026-11-25
- Finish workflow executed on 2026-11-25.
- Linked backlog/request close verification passed.

# Report
- b0527b0: `crow modules get|set nest` writes modules.nest in crow.yaml (line edits, read-back, rollback).
- 936deb0: core GET/POST /api/settings/nest runs the script; 400 on unknown names, 409 on a crow.yaml it cannot edit safely.
- Finished on 2026-11-25.
- Linked backlog item(s): `item_282_core_settings_endpoint_writes_modules_nest_in_crow_yaml`
- Related request(s): `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`

# Links
- Request: `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# AC Traceability
- request-AC1 -> This task. Proof: AC2 implemented in b0527b0 (crow modules get/set nest) and 936deb0 (GET/POST /api/settings/nest); validated with go test ./... in tools/crow-core and pytest tests/test_nest_settings.py tests/test_crow_modules.py. Source: `936deb0`
- request-AC2 -> This task. Proof: AC2 implemented in b0527b0 (crow modules get/set nest) and 936deb0 (GET/POST /api/settings/nest); validated with go test ./... in tools/crow-core and pytest tests/test_nest_settings.py tests/test_crow_modules.py. Source: `936deb0`
- request-AC3 -> This task. Proof: AC2 implemented in b0527b0 (crow modules get/set nest) and 936deb0 (GET/POST /api/settings/nest); validated with go test ./... in tools/crow-core and pytest tests/test_nest_settings.py tests/test_crow_modules.py. Source: `936deb0`
- request-AC4 -> This task. Proof: AC2 implemented in b0527b0 (crow modules get/set nest) and 936deb0 (GET/POST /api/settings/nest); validated with go test ./... in tools/crow-core and pytest tests/test_nest_settings.py tests/test_crow_modules.py. Source: `936deb0`
- request-AC5 -> This task. Proof: AC2 implemented in b0527b0 (crow modules get/set nest) and 936deb0 (GET/POST /api/settings/nest); validated with go test ./... in tools/crow-core and pytest tests/test_nest_settings.py tests/test_crow_modules.py. Source: `936deb0`
