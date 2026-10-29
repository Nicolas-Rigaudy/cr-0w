## item_017_0_2_real_data_on_home_and_agents - 0.2 real data on Home and Agents
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Home and Agents fed by cdx and crow-runs: right now, needs you, running, background runs, finished, who-worked-when timeline, live tail with agent questions.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Poll crow-runs json and cdx sessions; reuse existing crow-watch data layer.
  - Answer questions from the tab (`a x e c`), `g` jump, send text, `R` relaunch; pause, stop, attach from the `:` palette.
  - Agents as one box with group headers (AgentsL_merge).
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: No model calls to render.
- AC2: Questions from agents can be answered from the deck.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: No model calls to render.
- request-AC2 -> This backlog slice. Proof: AC2: Questions from agents can be answered from the deck.

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
