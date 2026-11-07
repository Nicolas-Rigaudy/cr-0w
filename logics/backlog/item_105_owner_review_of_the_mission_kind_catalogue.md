## item_105_owner_review_of_the_mission_kind_catalogue - Owner review of the mission kind catalogue
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Low
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:49

# AI Context
- Summary: The owner reviews the 8 draft kind files (gates, autonomy, caps, models, publish, preflight) and settles where kind files live.
- Keywords: kind catalogue, owner review, caps, gates, kind file home
- Use when: Before any code reads kinds for real.
- Skip when: Schema or parser work.

# Problem
- Request AC1 needs the owner to review the kind files.
- req_028 slice 1 says kinds live in a private `missions.yaml`; the spec puts them in tracked kind files. Decision D1.

# Scope
- In:
  - Owner pass over the 8 files in the spec: gates, autonomy, caps, models, publish actions, preflight questions, scope overrides.
  - Apply D1: move or keep the files, add a private override path if chosen, update req_028 slice 1 text.
  - Record the outcome in req_037.
- Out:
  - New kinds.
  - Account names in kind files (scope routing through item_095).

# Acceptance criteria
- AC1: The owner has reviewed all 8 kind files and every change is applied.
- AC2: D1 is decided (tracked files plus private overrides) and req_028 slice 1 matches it. Done 2026-11-07.

# AC Traceability
- request-AC1 -> AC1. Proof: owner review recorded in req_037.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Every kind default waits on it.

# Notes
- Waits on the owner review of the 8 files. D1 is decided (2026-11-07): tracked kind files plus private overrides in `local/` (item_117); req_028 slice 1 is updated.
- Spec: `logics/specs/req_037_run_contract/`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
