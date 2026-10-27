## item_013_run_reports_per_provider - Parse run reports per provider (#5)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: Hub portability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: `final_text(stdout, provider)` parses the Claude result envelope and the Codex last agent message; crow-runs report and `crow oneshot` use it (#5).
- Keywords: parse, run, reports, per, provider
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- `final_text(stdout, provider)` parses the Claude result envelope and the Codex last agent message; crow-runs report and `crow oneshot` use it (#5).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: fixtures for both formats.
- AC2: unknown format falls back to the raw tail.
- AC3: tests green.

# AC Traceability
- request-AC4 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_005_codex_support_and_cleanups`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: order in req_005 (profiles helper first).

# Notes
- Generated locally by logics-manager.
