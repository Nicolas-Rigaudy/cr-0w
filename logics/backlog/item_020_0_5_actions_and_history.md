## item_020_0_5_actions_and_history - 0.5 actions and History
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
- Summary: Actions that change things: answer an agent, handoff to a repo session, start a review, snooze, edit roadmap. Session History behind `/` in Agents.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Every push, post or send stays behind an explicit OK shown in the deck.
  - History reads cdx history and transcripts.
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: Autonomy table of CLAUDE.md is enforced in code, not by convention.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: Autonomy table of CLAUDE.md is enforced in code, not by convention.

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
