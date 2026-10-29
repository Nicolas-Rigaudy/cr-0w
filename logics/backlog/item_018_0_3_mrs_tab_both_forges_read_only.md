## item_018_0_3_mrs_tab_both_forges_read_only - 0.3 MRs tab: both forges, read only
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 85%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: MRs and PRs where the owner is assigned or a requested reviewer, on GitLab and GitHub. Reviewer and author views, thread list, thread card with fix diff, End of review box.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Forge fetch by script, cached.
  - Finding/thread cards with accept, refuse, edit, chat; diff view.
  - Structured review findings from the review skill (gap: today free text).
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: Read only until 0.5: nothing is posted or pushed from this step.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: Read only until 0.5: nothing is posted or pushed from this step.

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
- `task_006_0_3_mrs_tab_both_forges_read_only`
