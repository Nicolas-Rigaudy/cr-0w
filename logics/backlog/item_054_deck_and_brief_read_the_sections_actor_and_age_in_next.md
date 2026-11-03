## item_054_deck_and_brief_read_the_sections_actor_and_age_in_next - Deck and brief read the sections: actor and age in Next
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
> Complexity: Small
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Deck Next list and project NEXT read the sectioned open threads and show actor and age; brief groups the same way.
- Keywords: crow deck, Next, openThreads, actor, age, data.go
- Use when: wiring the deck Next list to the new open-threads format.
- Skip when: the Next layout itself (req_011 mockups decide it).

# Problem
- `tools/crow-deck/data.go:openThreads` reads unchecked lines with no actor; the Next list hides who must act (req_011 audit).

# Scope
- In:
  - Parser tracks the current `## ` section as actor (me / crow / person from Waiting lines) and the date as age.
  - nextRow gains Actor and Since; fixtures and shot tests updated.
  - Layout follows the req_011 pick for Home Next.
- Out:
  - The to-do model (item_052).

# Acceptance criteria
- AC1: Next rows carry actor and age from the file.
- AC2: Unit test covers the three sections and a Waiting line with a person.
- AC3: Rendered as picked in req_011.

# AC Traceability
- request-AC4 -> AC1-AC3. Proof: AC1-AC2 `tools/crow-deck/data.go` openThreads (Actor, Since, Why) and `TestOpenThreads`; AC3 lands with the Home A build (item_057).

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_012_owner_todos_and_open_work_model
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: The deck Next list is the owner's main view of what needs them.

# Notes
- Execution order: 3 of 3, after item_053 and the req_011 Home pick.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

