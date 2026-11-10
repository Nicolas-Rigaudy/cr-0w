## item_182_retro_mission_kind_that_improves_crow - Retro mission kind that improves crow
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: the retro becomes a mission kind: the daily batch from `crow-retro` runs as a retro mission, each proposed crow fix lands in the Nest inbox, an approved fix runs as a build mission on cr-0w, and the owner can start a retro on any session from Nest.
- Keywords: retro, mission kind, crow-retro, improvement, inbox, build mission, budget_unverifiable
- Use when: changing how retros run or how their fixes are applied.
- Skip when: the retro skill's analysis itself.

# Problem
- Retros run as a timer batch and write vault reports; the owner reads them and applies tracked fixes by hand. The 2026-11-06 batch failed 3/3 with cdx `budget_unverifiable` ("Measured usage cannot be priced; USD budget enforcement is unavailable.") and nothing surfaced it. Owner 2026-11-10: run retros easily, as missions.

# Scope
- In:
  - `kinds/retro.yaml`: skill retro, key = session id or `batch:<date>`, scope personal (hub), autonomy notify, outputs report plus a list of proposals.
  - Local-only lessons (runbooks, memories) stay applied by the retro as today. Each tracked proposal becomes an inbox item `proposal` with the diff summary; Approve starts a build mission on cr-0w with the proposal as input; Reject records the note.
  - `crow-retro` timer creates the batch mission instead of a raw cdx run. A failed run becomes an inbox alert.
  - Fix `budget_unverifiable`: find why cdx cannot price usage (model id or account plan), set the kind's caps so the run starts; if cdx cannot price it, use a time cap and note it.
  - Nest: "Retro this" on any finished session or mission.
- Out:
  - Retros on work repos' code (only crow and the owner's habits).

# Acceptance criteria
- AC1: The weekday batch runs as a retro mission and its report reads in Nest.
- AC2: Each tracked proposal is an inbox item; Approve starts a cr-0w build mission linked to it.
- AC3: "Retro this" from Nest on a session starts a retro mission for that transcript.
- AC4: A retro run that fails shows as an alert in Nest the same day; the budget_unverifiable cause is fixed or written down.

# AC Traceability
- request-AC6 -> AC1, AC2, AC3, AC4. Proof: retros as missions feeding approved fixes.

# Decision framing
- Product framing: Not needed (owner 2026-11-10 picked the retro mission kind)
- Architecture framing: Not needed (req_037 contract, adr_001 inbox)

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`, `req_037_mission_kinds_and_run_contract`
- Depends on: `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`, `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: Medium
- Rationale: makes every day's friction cheaper to fix.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
