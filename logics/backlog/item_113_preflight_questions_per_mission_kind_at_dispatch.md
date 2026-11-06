## item_113_preflight_questions_per_mission_kind_at_dispatch - Preflight questions per mission kind at dispatch
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Before dispatch, ask the kind's preflight questions; on skip, the run takes the safest defaults and lists them as assumptions.
- Keywords: preflight, safest default, assumptions, dispatch
- Use when: Creating a mission.
- Skip when: Questions a run asks mid-run (outbox item).

# Problem
- A run should ask up front only what it cannot decide itself.
- A skipped preflight must not block; the safest default applies and is stated.

# Scope
- In:
  - Read `preflight` from the kind file and ask at creation.
  - On skip, write the safest defaults into the brief; the run lists them in `assumptions`.
  - AFK batches ask all preflights of the batch at once.
- Out:
  - The nest or deck screen that asks them (UI items).

# Acceptance criteria
- AC1: A mission created with preflight answers carries them in the brief.
- AC2: A skipped preflight puts every safest default in the brief and the run's `assumptions`.

# AC Traceability
- request-AC1 -> AC2. Proof: brief and plan of a mission with skipped preflight.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Useful once several kinds run; defaults work meanwhile.

# Notes
- Depends on item_105 (reviewed questions) and item_107 (brief).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
