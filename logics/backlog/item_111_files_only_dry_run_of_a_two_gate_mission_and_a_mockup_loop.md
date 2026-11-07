## item_111_files_only_dry_run_of_a_two_gate_mission_and_a_mockup_loop - Files-only dry run of a two-gate mission and a mockup loop
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 01:33:54

# AI Context
- Summary: Run one two-gate mission and one mockup loop by hand from files, resuming on a cold session, to prove the contract is enough.
- Keywords: dry run, two-gate, mockup loop, resume, cold run
- Use when: Proving the contract before core.
- Skip when: Building core.

# Problem
- Request AC4 needs proof that brief plus version files are enough to resume.
- cdx `--continue-run` pins the parent's account, so a resume on another account is a cold run. Decision D2.

# Scope
- In:
  - Two-gate mission: plan v1, dig, plan v2, go, report, ok, each turn as a separate headless run.
  - At least one resume as a cold run on a different account class.
  - Mockup loop: three options, change, pick, approve.
  - Validate every file with item_106; write the events by hand or with a script.
- Out:
  - Publishing (the dry run stops before push).

# Acceptance criteria
- AC1: Both missions reach their last gate from files only.
- AC2: The cold resume produces a valid next version that answers every point.
- AC3: Gaps found in the contract are fixed in the spec.

# AC Traceability
- request-AC4 -> AC1. Proof: mission dirs from the dry run.
- request-AC4 -> AC2. Proof: cold resume output.

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
- Rationale: Last check before core builds on the contract.

# Notes
- Waits on D2 (resume model), item_106 and item_107 (brief for the cold resume).
- Uses the scope-aware account pick (item_095).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
