## item_043_review_run_read_only_dispatch_that_writes_a_findings_json_per_mr - Review run: read-only dispatch that writes a findings JSON per MR
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
- Summary: Dispatch a read-only review run for one MR and write its findings JSON.
- Keywords: review, run, read, only, dispatch, writes, findings, json, per
- Use when: building the deck review flow (req_010).
- Skip when: the work is outside the deck review flow.

# Problem
- Dispatch a read-only review run for one MR and write its findings JSON. Part of req_010. Do first. Check on !214 that --resumes later keeps the context.

# Scope
- In:
  - bin/crow-runs dispatch with a read-only permission and the review skill; the run writes <state>/reviews/<group>_<name>-<number>.json (findings: path, line, severity, comment, proposed fix, pattern hit).
- Out:
  - Define the JSON schema in the repo and document it in runbooks. Add --model opus by default now that Opus is the default.

# Acceptance criteria
- AC1: A dispatch for an MR id writes a valid findings JSON.
- AC2: The run has no write permission on the repo.
- AC3: The schema is documented and covered by a test with a fixture.

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

