## task_021_orchestrate_a_work_day_from_nest - Orchestrate a work day from Nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: High
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: build order for req_055 (butler chat, mission chat, push, mail and calendar, rituals, always-on core, retro missions), after task_020 lands the sections, actions and approvals.
- Keywords: orchestrate, nest, butler, push, mail, rituals, retro, missions first
- Use when: picking the next req_055 item to build.
- Skip when: task_020 items (deck parity, Home design).

# Context
- Owner 2026-11-10: run the work day from Nest without a terminal; most runs will be missions.
- Starts once task_020 steps 4 to 7 and 11 are Done (sections, actions, prompts, approvals). item_181 spike can run any time.
- One fresh session per item, from its backlog doc. Owner prefs go into every handoff. Screens need a mockup and owner OK first.

# Plan
- [ ] 1. item_181 spike on an always-on core (parallel, any time). Exit: owner picked a host or a written limit.
- [ ] 2. item_146 SSE and item_148 conversation poller (from task_017). Exit: a run's turns stream to a curl client.
- [ ] 3. item_177 chat inside a mission or session (mockup first). Exit: read and write a running mission from the phone.
- [ ] 4. item_178 phone push. Exit: a parked mission buzzes the phone.
- [ ] 5. item_182 retro mission kind (fixes budget_unverifiable first). Exit: the weekday batch runs as a mission and proposals reach the inbox.
- [ ] 6. item_180 rituals as mission kinds. Exit: Brief from Nest.
- [ ] 7. item_176 butler chat (mockup first). Exit: "start a review of X" from Nest creates the mission.
- [ ] 8. item_179 mail and calendar actions (after item_172). Exit: a reply sent after Approve from the phone.
- [ ] Order rule: missions first in every list and action. If a design step waits on owner OK, build the next item.
- [ ] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_176_butler_chat_in_nest_on_the_core_mcp_tools`
- `item_177_chat_inside_any_session_or_mission_run_from_nest`
- `item_178_phone_push_when_something_needs_the_owner`
- `item_179_mail_and_calendar_actions_from_nest`
- `item_180_rituals_as_mission_kinds_started_from_nest`
- `item_181_core_reachable_when_the_laptop_sleeps`
- `item_182_retro_mission_kind_that_improves_crow`

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`
- Depends on: `task_020_orchestrate_nest_going_live_on_the_core_api`

# Definition of Done (DoD)
- [ ] Every listed item is Done with its acceptance criteria proven.
- [ ] One full work day run from Nest on the phone with no terminal opened, or each gap written down.
- [ ] Validation passes.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
