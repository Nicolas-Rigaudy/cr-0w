## item_302_shared_page_shell_and_deck_components_for_every_nest_page - Shared page shell and deck components for every Nest page
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: One page shell and one set of deck components shared by every Nest page, so layout fixes land everywhere at once.
- Keywords: shared, page, shell, deck, components, nest
- Use when: a layout fix or a new page in Nest.
- Skip when: page content or data.

# Problem
- Owner 2026-11-26: a layout fix on one page should apply everywhere, from a shared library like a design system.
- Today only tokens (colours, fonts, radius) are shared. Each page root sets its own padding, width and height rules (`.screen`, `.dk/.mrs.m-lp`, `.usage.u-lp`, `.agenda.m-lp`, `.reports.r-lp`, `.ib-wide`, `.ib-lp`). The Deck frame is a component in Mrs.tsx, and its CSS is copied for Usage. Phone/laptop modes are decided per page.
- 2026-11-26 stopgap: one `--page-max` variable for all page widths.

# Scope
- In:
  - A `Page` shell component (title row, filters, body) that owns width, padding, full-height fill and the laptop/phone split for every page.
  - Move Deck, list rows, tabs, fact strips, buttons and the reader into `src/ui/` with one CSS block each; pages compose them.
  - Port every page; drop the per-page layout rules.
  - Headless shots of every page at 390, 1440 and 1920 wide before and after; no visual change except the fixes.
- Out:
  - New visuals or a theme change.

# Acceptance criteria
- AC1: Every page renders through the shared shell; no page root sets its own max-width, padding or height rule.
- AC2: Deck and its CSS exist once, in src/ui.
- AC3: A width or height change in the shell shows on every page (headless check over all pages).
- AC4: tsc, vitest and the shot checks pass.

# AC Traceability
- Owner request 2026-11-26 -> AC1 to AC3.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Each layout bug is fixed page by page today and comes back on the next page.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

