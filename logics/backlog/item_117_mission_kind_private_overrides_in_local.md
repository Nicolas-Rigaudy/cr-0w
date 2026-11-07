## item_117_mission_kind_private_overrides_in_local - Mission kind private overrides in local
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Load tracked generic kind files and merge optional private overrides from `local/` on top. Private values never enter tracked files.
- Keywords: kind files, overrides, local, loader, publishable
- Use when: Reading kinds at runtime.
- Skip when: Spec validation (item_104).

# Problem
- req_028 slice 1 said a private `local/missions.yaml`; the spec uses tracked `kinds/<kind>.yaml`.
- Owner decision 2026-11-07 (D1): shared tracked files plus optional private overrides in `local/`.

# Scope
- In:
  - Override path in `local/` (one file per kind or a single overrides file, chosen when built).
  - Merge rules: overrides replace scalar fields and lists as a whole, and may not drop a hard rule.
  - Overrides validate against the same kind schema after the merge.
  - Doctor shows which kinds carry an override.
- Out:
  - Account names in kind files (scope routing through item_095).
  - New kinds.

# Acceptance criteria
- AC1: A kind without an override loads as tracked.
- AC2: An override changes a cap and the merged kind validates.
- AC3: An override that removes a gate required by a hard rule fails validation.
- AC4: No private value appears in a tracked file or in test fixtures.

# AC Traceability
- request-AC1 -> AC1 to AC4. Proof: the acceptance criteria above, checked when built.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_037_mission_kinds_and_run_contract`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Keeps the engine publishable while allowing private tuning.

# Notes
- 2026-11-07 owner decision 11 (req_037 D1 option A). req_028 slice 1 updated to match.
- Extends item_104 (validator) and item_105 (owner review).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

