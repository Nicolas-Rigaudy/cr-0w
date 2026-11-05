## item_094_deck_mail_draft_pane - Deck mail pane to review and finish drafts
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
- Summary: Deck mail pane: Today rows show draft state; a pane shows the thread summary and the draft; o opens the draft, R redrafts with an instruction, x drops; enter on a crow flag queues a draft.
- Keywords: deck, mail, pane, draft, review, today, keys
- Use when: building this slice of req_024_background_mail_drafts.
- Skip when: other mail work.

# Problem
- The owner wants one place to see where every mail stands and finish drafts quickly.

# Scope
- In:
  - Mockup options (bin/crow-compare) shown to the owner before building; the owner picks.
  - Today rows: draft state badge.
  - Mail pane: summary + draft text, keys o (open draft in Outlook), R (one-line instruction → redraft), x (drop), enter on a crow flag → queue.
  - Palette keeps `mail session` for the rare interactive case.
- Out:
  - Editing the draft text inside the deck (Outlook does that).

# Acceptance criteria
- AC1: The owner picked a mockup before the build.
- AC2: Keys map to the item_093 subcommands; tests cover them; go test passes.
- AC3: Live: the owner finishes one real draft from the pane.

# Dependencies
- item_093 (mockups can start with item_092)

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

