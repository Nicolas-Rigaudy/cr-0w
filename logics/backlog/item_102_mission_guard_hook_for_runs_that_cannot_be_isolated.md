## item_102_mission_guard_hook_for_runs_that_cannot_be_isolated - Mission guard hook for runs that cannot be isolated
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Safety
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A PreToolUse guard denies push, MR creation, merge, tracker and wiki writes and gate answers in runs that cannot be isolated, unless core published them.
- Keywords: mission, guard, hook, runs, cannot, isolated
- Use when: Continued runs, interactive take-control sessions and sessions of providers without isolation.
- Skip when: Isolated runs, where the OS boundary already holds.

# Problem
- Some runs stay outside the sandbox: `--continue-run`, take control, interactive hub work. Permission rules miss commands by path and inside `sh -c`.

# Scope
- In:
  - One guard script shared by mission kinds, built on `crow-hub-guard` patterns.
  - Never-list: force push, push to default branch, branch delete, merge, gate answers (`ok`, `go`, approve calls).
  - Allowed only when core holds a validation record for the exact action.
  - Compound commands containing a push are denied.
- Out:
  - Treating the guard as a security boundary (it is a policy layer).

# Acceptance criteria
- AC1: The never-list is denied in every guarded run, including inside `sh -c` and compound commands.
- AC2: A run cannot write an owner validation or answer a gate.
- AC3: Denials show in `crow runs` next to item_077 permission denials.

# AC Traceability
- request-AC2 -> run-config layer for non-isolated runs.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Covers the gap the OS boundary leaves.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- Reuses the reference design's guard idea, generic here.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

