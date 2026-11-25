## item_282_core_settings_endpoint_writes_modules_nest_in_crow_yaml - Core settings endpoint writes modules.nest in crow.yaml
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Core GET and POST routes for Nest settings: read and write `modules.nest.{theme,hero,pages}` in crow.yaml, keep comments, refuse unknown values.
- Keywords: core, settings, endpoint, crow.yaml, modules.nest, theme, hero, pages
- Use when: building the write side of the Nest Settings page.
- Skip when: the change is the page UI (item_284) or the rebuild (item_283).

# Problem
Nest can show theme, hero and pages but cannot change them. Core is the only writer of crow.yaml, so Nest needs a core route that reads the current values with their known lists and writes a change safely.

# Scope
- In:
  - `GET /api/settings/nest`: current `theme`, `hero`, `pages` from crow.yaml and the known values (engine lists plus brick themes, heroes and pages from `crow modules json`).
  - `POST /api/settings/nest`: body `{theme?, hero?, pages?}`; validates every value against the known lists, writes only `modules.nest.*`, returns the new values and `needs_build` (true when hero or pages changed).
  - One writer: a `crow-modules` subcommand (for example `crow modules set nest theme cyber`) built on the line edits of `crow_yaml_set` / `crow_yaml_switch`, called by core like other script actions. Core does not parse or rewrite YAML itself.
  - Same guard as the other write routes (owner login, CSRF header, origin).
- Out:
  - Running the build (item_283), the page UI (item_284), brick on/off (item_285).

# Acceptance criteria
- AC2: Saving writes `modules.nest.*` in crow.yaml through core, keeps the rest of the file and its comments, and refuses unknown values.
- AC2a: A test edits a crow.yaml with comments and other bricks; the diff touches only the `modules.nest` lines.
- AC2b: An unknown theme, hero or page returns 400 with the bad value named; crow.yaml is unchanged.
- AC2c: `go test ./...` in tools/crow-core, the crow-modules tests and `bin/crow-check` pass.

# AC Traceability
- request-AC2 -> This backlog slice. Proof: AC2: Saving writes `modules.nest.*` in crow.yaml through core, keeps the rest of the file and its comments, and refuses unknown values.

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
- Rationale: First in order (282 → 283 → 284 → 285). Does not depend on the layout pick, but stays Draft until the owner picks a Settings layout so the page and its API ship together.

# Notes
- Order: 1 of 4. Reuse `yaml_key_span` and `crow_yaml_set` in `bin/crow-modules`; the write side pattern is `tools/crow-core/internal/api/actions.go`.
- Derived from request `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

