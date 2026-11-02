## item_044_deck_review_pane_show_findings_and_record_a_x_e_c_decisions_in_the_json - Deck review pane: show findings and record a/x/e/c decisions in the JSON
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
- Summary: Show the findings in the MRs tab and record decisions.
- Keywords: deck, review, pane, show, findings, record, decisions, json
- Use when: building the deck review flow (req_010).
- Skip when: the work is outside the deck review flow.

# Problem
- Show the findings in the MRs tab and record decisions. Part of req_010. Needs item_043.

# Scope
- In:
  - MRs tab reads the findings JSON, shows each finding next to its diff, and a/x/e/c write accept, refuse, edit or chat into the JSON. A key starts item_043's run and shows its state.
- Out:
  - Spec boards MrsT_* and _merge. Read only on the repo.

# Acceptance criteria
- AC1: Findings render in the three themes at 80x24 and wide.
- AC2: a/x/e/c persist and survive a restart.
- AC3: A running review shows as running, failed shows the reason.

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

