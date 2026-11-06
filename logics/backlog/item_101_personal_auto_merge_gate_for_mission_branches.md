## item_101_personal_auto_merge_gate_for_mission_branches - Personal auto-merge gate for mission branches
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
- Summary: core merges a personal mission branch only after the req_030 checklist: pause off, rebase clean, commit subjects, secret scan over the range, sensitive paths ask, size limit, full-suite receipt, green CI, then deploy hook.
- Keywords: personal, auto, merge, gate, mission, branches
- Use when: Building or changing auto-merge for personal repos.
- Skip when: Work repos (never merge).

# Problem
- prod_005 settles auto-merge on green CI and a secret scan for personal repos; the checklist and its receipts do not exist yet.

# Scope
- In:
  - The ten-step checklist from req_030 Need 4, stopping at the first failure as an inbox item.
  - Global sensitive path list in core plus per-repo globs in the registry.
  - Full-suite receipt `{sha, command, exit, duration}` from a fresh isolated workspace with network off.
  - Secret scan over merge base..head with the prepush deny list plus a scanner.
  - Test on one personal repo end to end.
- Out:
  - Choosing the scanner and the size limit (decisions 5 and the scanner question).

# Acceptance criteria
- AC1: A branch touching a sensitive path asks instead of merging.
- AC2: A planted fake secret in the range blocks the merge.
- AC3: A push to the default branch without a receipt for the head SHA is refused.
- AC4: A diff above the limit asks instead of merging.
- AC5: One personal repo merges a clean mission branch through the gate and logs `publish`.

# AC Traceability
- request-AC4 -> auto-merge gate as a checklist, tested on one personal repo.

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
- Rationale: Unlocks the build and repair kinds; comes after the sandbox.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- Depends on the core publish path and isolated runs.
- Repo settings on the forge (auto-merge, branch protection) need owner OK before any change.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

