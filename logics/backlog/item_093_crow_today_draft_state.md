## item_093_crow_today_draft_state - crow-today exposes draft state, redraft and drop
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: deck
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-today json carries draft state and summary per flagged mail; new subcommands queue, redraft with an instruction, drop.
- Keywords: crow-today, draft, state, redraft, drop, queue
- Use when: building this slice of req_024_background_mail_drafts.
- Skip when: other mail work.

# Problem
- The deck needs to know where each flagged mail stands and to act on a draft.

# Scope
- In:
  - Flagged rows gain `draft`: {state: none|queued|drafting|ready|failed, summary, draftLink}.
  - `crow-today queue <id>`, `redraft <id> "<instruction>"` (deletes the previous crow draft, queues with the instruction), `drop <id>` (deletes the crow draft, state none).
  - Done/answered mail clears its state.
- Out:
  - Deck UI.

# Acceptance criteria
- AC1: json shows the right state through queue → drafting → ready.
- AC2: redraft leaves exactly one crow draft in Outlook; drop leaves none.
- AC3: Tests cover the transitions.

# Dependencies
- item_092

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_024_background_mail_drafts
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner request 2026-11-05 after the first live mail session.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

