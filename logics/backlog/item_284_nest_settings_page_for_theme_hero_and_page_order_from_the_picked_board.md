## item_284_nest_settings_page_for_theme_hero_and_page_order_from_the_picked_board - Nest Settings page for theme, hero and page order from the picked board
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The Nest Settings page (phone and laptop) for theme, hero, page order and which pages show, built from the layout the owner picks.
- Keywords: nest, settings, page, theme, hero, page order, mockup, board
- Use when: building the Settings screen in tools/nest.
- Skip when: the owner has not picked a layout yet (build is gated on the pick).

# Problem
The owner edits crow.yaml by hand to change the Nest look and page list. Nest needs a Settings page on phone and laptop that shows the current values and saves through core.

# Scope
- In:
  - Gate: round 1 (A, B, C in `settings.html`) rejected as too busy. Round 2 option D picked on 2026-11-25: archeotech shell layout, sidebar of panes (Look, Pages, Modules, General), one setting per row with a hint, Save bar per pane; phone = section list then pane. Board: `tools/nest/scripts/mockups/settings-v2.html` (`?v=D&p=look|pages|list`), PNGs in vault `CR-0W/Research/nest-settings/round2/settings-D-*.png`. A styling pass on D comes before the build (owner rule: layout and styling are separate picks).
  - Settings entry in the top bar (laptop) and under More (phone).
  - Theme picker with a live preview (data-theme) before save; hero list; page list with order and show/hide; Home stays on.
  - Pending changes, discard, one save; the build state from item_283.
  - No chat panel reserved on the page.
- Out:
  - Module switches (item_285), general settings beyond a placeholder section.

# Acceptance criteria
- AC1: A Settings page lists theme, hero and pages with their known values and shows the current crow.yaml values.
- AC1a: The page matches the picked board; a screenshot next to the board is shown before calling it done (one screen at a time).
- AC1b: Saving calls item_282, shows item_283's build state, and reloads on a passing build.
- AC5: Tests and crow-check pass (Nest unit tests, typecheck, `bin/crow-check`).

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: A Settings page lists theme, hero and pages with their known values and shows the current crow.yaml values.
- request-AC5 -> This backlog slice. Proof: AC5: Tests and crow-check pass.

# Decision framing
- Product framing: Needed (layout pick)
- Product signals: owner picks layout A, B or C, then a styling.
- Product follow-up: record the pick in this item before moving it to Ready.
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
- Rationale: Order 3 of 4, after items 282 and 283. Blocked on the owner's layout pick.

# Notes
- Options: A one page in three columns with a sticky save bar; B section rail (Look, Pages, Modules, General) with a crow.yaml diff and a live preview; C tabs with page tiles and a pending queue with undo per change.
- Derived from request `req_066_nest_settings_page_theme_hero_page_order_and_module_switches`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

