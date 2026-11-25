## item_286_nest_mrs_follow_ups_review_through_crow_review_run_key_hints_rerun_ci - Nest MRs follow-ups: review through crow-review-run, key hints, Rerun CI
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, mrs, follow, ups, review, through, crow, run, key, hints, rerun
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- The Nest MRs page (item_225) landed without three parts of board C. Owner 2026-11-25 picked crow-review-run for the Review action (same runner as the deck); today Nest dispatches `/review <url>`.

# Scope
- In:
  - crow-core action that runs `crow-review-run start <repo> <n>` under the same autonomy gate as crow-deck (`runGated`, kind MR); Nest Review action calls it instead of `reviewPrompt`.
  - Key hints bar from board C (1-4 panes, j k, r, t, s, o, p, c) with the shortcuts wired.
  - Rerun CI on the detail actions, once a forge retry action exists (search crow-core/crow-mrs first); hidden when the forge has none.
- Out:
  - Posting comments from Nest.

# Acceptance criteria
- AC1: Review on an MR row starts a crow-review-run review through crow-core, gated like the deck; vitest and a Go test cover it.
- AC2: The key hints bar shows on laptop and each listed key works.
- AC3: Rerun CI retries the failed pipeline through a confirm sheet, or the button is absent when the forge offers no retry.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: bounded delivery slice.
- request-AC2 -> This backlog slice. Proof: promotable backlog item.
- request-AC3 -> This backlog slice. Proof: delivery chain includes a task-ready backlog item.

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

