## item_151_autonomy_dial_per_kind_and_mission_under_the_hard_rules - Autonomy dial per kind and mission under the hard rules
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Autonomy (auto, notify, ask) from the kind file, overridable per mission, drives what core does at each gate and at the end: auto finishes with an FYI inbox item, notify parks nothing and sends an item, ask parks at the first decision. Hard rules sit above the dial.
- Keywords: autonomy, auto, notify, ask, gates, hard rules, fyi, inbox, policy, kinds
- Use when: changing what core does at a gate or at the end of a mission.
- Skip when: the publish path itself (req_030), the deck action table (`tools/crow-deck/autonomy.go`, which governs deck actions, not missions).

# Problem
- Slice 5. The mission record already carries `autonomy`, and the tick sends a gateless `auto` mission with a `ready` verdict straight to done (tools/crow-core/internal/mission/tick.go). There is no FYI item, `notify` behaves like `ask`, and nothing stops a dial from reaching a hard rule.

# Scope
- In:
  - Policy table in core: per gate and end state, what auto, notify and ask do. Hard rules (work push, MR, merge, Jira and Confluence writes, mail, Teams; personal sensitive paths per adr_001) always produce an approval item, whatever the dial.
  - FYI inbox item type (read, no verb needed, dismiss) for auto and notify ends.
  - `ask`: the brief tells the run to stop at its first decision; a run that passes a decision without a question is flagged in the report check.
  - Override at dispatch (`crow mission new --autonomy`, item_147 field); the effective value and its source shown in `crow mission show`.
  - Fixtures for each dial and gate pair; brief template line "Allowed" filled from the dial.
- Out:
  - Turning auto on for real kinds (item_152), personal auto-merge (item_101), mission guard hook (item_102), stop guard (item_110).

# Acceptance criteria
- AC1: A test mission with autonomy auto finishes with only an FYI item.
- AC2: A test mission with autonomy ask stops at its first decision and parks.
- AC3: A work mission set to auto that reaches a push or post still gets an approval item (test).
- AC4: A table test covers every dial against every gate.

# AC Traceability
- request-AC3 -> AC1, AC2, AC3, AC4. Proof: slice 5 exit.

# Decision framing
- Product framing: Not needed (prod_005 kinds table)
- Architecture framing: Not needed (adr_001 publish section, req_037)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_037_mission_kinds_and_run_contract`
- Depends on: `item_134_inbox_items_and_the_crow_inbox_cli` (Done). AC2 end to end needs a run that asks, which works today (parking proven 2026-11-09); resume is not needed.
- Related: `item_102`, `item_110`, `item_112`
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: High
- Rationale: Safety layer before any kind runs unattended.

# Notes
- Groomed 2026-11-09 from req_028 slice 5.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
