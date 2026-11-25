## item_285_module_switches_on_the_nest_settings_page_through_crow_modules - Module switches on the Nest Settings page through crow modules
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Second slice of the Settings page: turn bricks on or off through `crow modules`, show `requires` blocks, relink skills.
- Keywords: nest, settings, modules, bricks, switches, requires, crow modules link
- Use when: adding brick on/off to the Nest Settings page.
- Skip when: installing bricks from git (out of scope of req_066).

# Problem
Turning a brick on or off needs a terminal (`crow modules`). From a phone or laptop the owner cannot change which bricks run.

# Scope
- In:
  - Core route that lists bricks (name, summary, on, required by, requires) from `crow modules json` and switches one through `crow modules` (same line edit as `crow_yaml_switch`), then runs `crow modules link`.
  - A brick that another on brick requires cannot be turned off; the row says which brick needs it. A brick whose `requires` is off cannot be turned on until it is.
  - Bricks that ship Nest themes, heroes or pages trigger the item_283 rebuild after a switch.
  - The Modules section or tab on the page from the picked layout.
- Out:
  - `crow modules add/update/remove` (git installs), editing brick config other than `on`.

# Acceptance criteria
- AC4: Module switches (second slice) turn a brick on or off through `crow modules`, show `requires` blocks, and relink skills.
- AC4a: Tests cover a blocked off (forge needs git), a blocked on (requires off) and a switch that relinks skills.
- AC4b: Go tests, crow-modules tests, Nest tests and `bin/crow-check` pass.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC4: Module switches (second slice) turn a brick on or off through `crow modules`, show `requires` blocks, and relink skills.

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
- Priority: Low
- Rationale: Order 4 of 4, after item_284. Gated on the same layout pick.

# Notes
- Engine bricks today: core, core-api, nest (required by Nest itself), git, forge (needs git), mail, vault, team, deck, afk, releases, retro, design (off by default), office (off by default).
- Derived from request `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

