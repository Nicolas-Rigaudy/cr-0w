## item_076_review_runs_isolated_workspace - Review runs use cdx isolated workspaces
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Review runs use cdx --isolate-workspace instead of crow's own checkout preparation, if it covers the review's needs.
- Keywords: review, isolate-workspace, worktree, cdx
- Use when: changing how review runs get their checkout.
- Skip when: other req_019 slices.

# Problem
- crow-review-run prepares the checkout itself; cdx now isolates a run's workspace with git and RTK access.

# Scope
- In:
  - Read how crow-review-run prepares the checkout today and what the review needs (MR branch, diff base, repo MCP, writing findings back).
  - Try one real review with --isolate-workspace and compare result and cost.
  - Switch if it covers everything; otherwise record the gap here and keep the current path.
- Out:
  - Task runs (they need to commit in a worktree under <repo>.wt/).

# Acceptance criteria
- AC1: A trial review run with --isolate-workspace is compared against the current path (findings, cost, time).
- AC2: Either review runs use it with tests updated, or this item lists the blocking gap.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC1 to AC2 above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner picked it from the 2026-11-03 release review.

# Notes
- Isolated workspaces must not land in /tmp if the hub guard forbids it there; check where cdx puts them.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

