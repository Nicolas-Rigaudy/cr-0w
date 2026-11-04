## item_047_chat_escape_hatch_c_opens_an_interactive_session_resumed_from_the_review - Chat escape hatch: c opens an interactive session resumed from the review
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
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
  - Owner decision 2026-11-02: `c` stays a note answered by the apply run (item_045). The interactive escape hatch is the existing `h` handoff, no new key.
  - `h` on an MR with a findings file adds its path to the handoff task, so the repo session starts from the review, the decisions and the notes.
  - Fix handoffTask: it checks `r.Role == "reviewer"` while roles are `REVIEWER`, so reviews get the author task.
- Out:
  - Resuming the review conversation in an interactive window (dropped).

# Acceptance criteria
- AC1: `h` on a reviewed MR shows the findings file path in the task, and the overlay shows it before enter.
- AC2: A reviewer MR gets the review task text, an own MR the threads text.
- AC3: Nothing runs without the confirm overlay showing the command.

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
| 2026-11-04T12:18:30.401961+00:00 | afk | afk-deckh |  | succeeded |  | AC1-3 met: go test + pytest green; findings path in task and overlay, REVIEWER case fixed, enter-only run |

