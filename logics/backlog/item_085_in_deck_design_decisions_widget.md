## item_085_in_deck_design_decisions_widget - In-deck design decisions widget
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: crow-deck
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Deck widget to flip options via kitty graphics and record picks in the ledger.
- Keywords: deck, design, decisions, widget
- Use when: Moving design decisions into the deck.
- Skip when: Before items 078 and 082 land.

# Problem
- The owner wants to converge daily work into crow deck and stop typing in terminals; picks still go through chat.

# Scope
- In:
  - Deck widget listing open decisions from the ledger; options shown inline with the kitty graphics protocol; flip with arrows, pick with a key, optional note; writes the ledger.
  - Fallback to opening the compare page when graphics are unavailable.
- Out:
  - Non-design decisions (a later generalization).

# Acceptance criteria
- AC1: A full round can be decided from the deck alone and the ledger updates.

# AC Traceability
- request-AC7 -> This backlog slice. Proof: AC1: A full round can be decided from the deck alone and the ledger updates.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)
- Request: `req_020_terminal_design_system_and_design_review_workflow`
- Primary task(s): `task_010_orchestrate_terminal_design_system`

# Priority
- Priority: Low
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

