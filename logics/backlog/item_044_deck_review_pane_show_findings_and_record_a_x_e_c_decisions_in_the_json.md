## item_044_deck_review_pane_show_findings_and_record_a_x_e_c_decisions_in_the_json - Deck review pane: show findings and record a/x/e/c decisions in the JSON
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
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
- Shipped: `crow-review-run decide`, `crow-mrs show` syncs a running review, deck `r` (confirm overlay, then `crow-review-run start --mode peer|self`), `a/x/e/c` call `decide`, `e/c` open an input line (enter saves, esc drops), card shows severity, comment, proposed fix, note, run state.
- Proofs: AC1 `TestReviewCardFitsEveryThemeAndSize` (three themes, 80x24, 170x46, 118x60, typing on and off). AC2 `TestDecisionsAreWrittenToTheFindingsJSON`, `TestSavedDecisionsShowAfterRestart`, `test_decide_persists_and_validates`. AC3 `TestReviewRunStateShows`, `TestRunningReviewRefetchesSooner`.
- Simulated runs only. No real review run or deck action launched.
- Threads mode `a/x/e/c` still flashes: answering threads is the apply slice.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

