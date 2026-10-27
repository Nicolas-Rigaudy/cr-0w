## item_014_codex_hub_instructions - Hub instructions and guard coverage for Codex (#1)
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
- Summary: Codex accounts get a read-first list (AGENTS.md, LOGICS.md, owner-prefs, memory), provider-aware model and permission check, interactive hub refused without a guard; guard coverage documented in a runbook (#1).
- Keywords: hub, instructions, guard, coverage, codex
- Use when: working on this slice of req_005.
- Skip when: the slice is done.

# Problem
- Codex accounts get a read-first list (AGENTS.md, LOGICS.md, owner-prefs, memory), provider-aware model and permission check, interactive hub refused without a guard; guard coverage documented in a runbook (#1).

# Scope
- In:
  - this slice only.
- Out:
  - other slices of req_005.

# Acceptance criteria
- AC1: `context_refs()` and prompt preamble for Codex.
- AC2: model/provider mismatch refused with a clear message.
- AC3: runbook table of guarded operations per provider.
- AC4: tests.

# AC Traceability
- request-AC5 -> this slice. Proof: its acceptance criteria.

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
