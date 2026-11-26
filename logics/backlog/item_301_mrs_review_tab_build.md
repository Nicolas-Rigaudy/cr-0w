## item_301_mrs_review_tab_build - MRs review tab build
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build the Nest MR Review tab from the picked board on the core review endpoints.
- Keywords: mrs, review, tab, build
- Use when: after item_299 is done and item_300 has a pick.
- Skip when: unrelated Nest pages.

# Problem
- The Nest MRs page must run the review from start to finish (req_071).

# Scope
- In:
  - Review button calls the start endpoint (peer on a review, threads on an own MR with open threads, self otherwise), never a plain dispatch.
  - Review tab per the picked board; polls while running; decisions saved at once; apply; post actions one at a time behind a confirm sheet with the exact text.
  - Fixtures and vitest; headless shots laptop and phone.
- Out:
  - Posting without the confirm sheet.

# Acceptance criteria
- AC1: reviewPrompt dispatch removed; Review starts crow-review-run.
- AC2: A finished review shows its findings in the MR detail, no hub question.
- AC3: Decisions survive a reload; apply and post work one action at a time.
- AC4: tsc, vitest, shots pass.

# AC Traceability
- req_071 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_071_nest_mrs_review_flow_and_visible_session_starts
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Reviews rank second in the owner's priorities; no-terminal days depend on Nest.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

