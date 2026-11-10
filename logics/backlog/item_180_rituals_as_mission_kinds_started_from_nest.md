## item_180_rituals_as_mission_kinds_started_from_nest - Rituals as mission kinds started from Nest
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
- Summary: brief, daily, eod, wrap and weekly become mission kinds (req_037 run contract) started with one tap in Nest; their output is a report in Nest and their questions come back as inbox asks.
- Keywords: rituals, brief, daily, eod, wrap, weekly, mission kind, run contract
- Use when: making a ritual skill run as a mission.
- Skip when: retro (item_182).

# Problem
- Rituals are skills the owner types in a terminal session. Most runs will be missions (owner 2026-11-10), rituals included.

# Scope
- In:
  - Kind files under logics/specs/req_037_run_contract/kinds/ for brief, daily, eod, wrap, weekly: skill, scopes, account route, caps, outputs, publish actions.
  - Interactive steps (eod questions, daily approval) park the mission as an inbox ask; the answer resumes it (item_147).
  - Publish steps (daily to Confluence, weekly slides) keep their approval (item_175).
  - wrap as a mission targets a named session (its transcript), not the mission itself.
  - Nest: a Rituals row on Home with the due ones lit (brief in the morning, daily before 11:30, eod in the evening, weekly on Friday).
- Out:
  - Changing what each ritual does. mgmt-weekly (later).

# Acceptance criteria
- AC1: Tapping Brief starts a brief mission; its report reads in Nest.
- AC2: eod asks its questions as inbox items and files the answers.
- AC3: daily's Confluence post waits on an approval on the exact text.

# AC Traceability
- request-AC4 -> AC1, AC2, AC3. Proof: rituals from Nest as missions.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (req_037 contract)

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`, `req_037_mission_kinds_and_run_contract`
- Depends on: `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`, `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`, `item_175_work_approvals_given_from_nest_on_the_exact_text`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: Medium
- Rationale: daily habit; cheap once kinds and answers work.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
