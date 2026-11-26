## item_299_core_review_endpoints_over_crow_review_run - Core review endpoints over crow-review-run
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
- Summary: Core endpoints that wrap bin/crow-review-run so Nest can run the review flow.
- Keywords: core, review, endpoints, over, crow, run
- Use when: Nest needs review start, findings, decisions, apply or post.
- Skip when: the deck, which calls the script directly.

# Problem
- crow-review-run holds the whole review flow, but only the deck can call it. Core exposes no endpoint, so Nest dispatches a plain run.

# Scope
- In:
  - GET review of one MR (crow-mrs show: findings, run, apply, mode).
  - POST start (mode, branch), decide (finding, decision, note), apply, actions list, act (id, hash). Each one shells out to the script, as other core actions do, and logs to actions.jsonl.
  - Go tests with a fake script.
- Out:
  - Any change to the findings schema or the script behaviour.

# Acceptance criteria
- AC1: Each endpoint maps 1:1 to a crow-review-run subcommand and returns its JSON or its error text.
- AC2: act needs the hash from actions (no blind post).
- AC3: Go tests cover start, decide, act refusal without hash.

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
- Priority: Medium
- Rationale: Reviews rank second in the owner's priorities; no-terminal days depend on Nest.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

