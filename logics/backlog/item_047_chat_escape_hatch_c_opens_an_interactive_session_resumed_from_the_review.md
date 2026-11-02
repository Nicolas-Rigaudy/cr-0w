## item_047_chat_escape_hatch_c_opens_an_interactive_session_resumed_from_the_review - Chat escape hatch: c opens an interactive session resumed from the review
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
- Summary: c opens an interactive session resumed from the review.
- Keywords: chat, escape, hatch, opens, interactive, session, resumed, review
- Use when: building the deck review flow (req_010).
- Skip when: the work is outside the deck review flow.

# Problem
- c opens an interactive session resumed from the review. Part of req_010. Last, optional.

# Scope
- In:
  - From a finding, c opens a window (crow term) resuming the review session so the owner can discuss it with Claude, then returns to the deck.
- Out:
  - Escape hatch only; the normal path stays inside the deck.

# Acceptance criteria
- AC1: c opens a session with the review context.
- AC2: Nothing runs without the confirm overlay showing the command.

# AC Traceability
- req_010 AC1-AC5 -> this slice. Proof: its own acceptance criteria.

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
- Rationale: Default until groomed.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

