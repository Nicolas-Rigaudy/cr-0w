## item_019_0_4_projects_and_usage_tabs - 0.4 Projects and Usage tabs
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: One box per project with NOW/NEXT/BLOCK/ROAD/MILE/OPEN, timeline underneath, action buttons; Usage with gauges, burn sparkline, daily spend bars, cost per session/project/model, savings, forecast.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Decide the roadmap file format and where it lives (logics roadmap docs are the candidate).
  - Verify per-session cost accounting in cdx stats and transcripts before building cost views.
  - Perso hidden in work mode (`m`).
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: Roadmap data comes from one file per project, edited with `e`.
- AC2: Cost per session matches `cdx stats` within rounding.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: Roadmap data comes from one file per project, edited with `e`.
- request-AC2 -> This backlog slice. Proof: AC2: Cost per session matches `cdx stats` within rounding.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Ships on its own; order follows the build plan in the spec.

# Tasks
- `task_007_0_4_projects_and_usage_tabs`
