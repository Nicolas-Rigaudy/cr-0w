## item_170_core_actions_on_live_sessions_and_dispatch_runs - Core actions on live sessions and dispatch runs
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: POST routes in crow-core that run deck's existing session and run actions (send text, interrupt, stop a session; stop a dispatch run; mark a mail done or snoozed), with the item_147 CSRF and identity checks, plus the Nest confirm sheet that calls them.
- Keywords: core actions, crow-focus, send, interrupt, stop, crow-runs stop, crow-today done, csrf, confirm
- Use when: piloting live sessions or dispatch runs from Nest.
- Skip when: mission actions (item_147), MR review actions, publish.

# Problem
- Nest can only watch. The owner wants to pilot from the phone, and today's work runs in sessions and dispatch runs, not missions (req_054).
- deck already has these actions as script calls (`tools/crow-deck/actions.go`, `live.go`). core must call the same scripts, not reimplement them.

# Scope
- In:
  - Routes under `/api`, same guard as item_147 (owner identity, CSRF header, socket trusted):
    | Route | Command |
    | --- | --- |
    | `POST /sessions/<pid>/send` `{text}` | `crow-focus --pid N --send TEXT` |
    | `POST /sessions/<pid>/interrupt` | `crow-focus --pid N --interrupt` |
    | `POST /sessions/<pid>/stop` | `crow-focus --pid N --stop` |
    | `POST /dispatch/<run_id>/stop` | `crow-runs stop RUN` |
    | `POST /today/<id>/done`, `/snooze` | `crow-today done|snooze ID` |
  - The pid must be in the current `sessions` section and the run id in `dispatch`; else 404. Stale view: 409.
  - Each action writes an event `{type: action, actor: owner, surface}` to the core event log, then refreshes the touched section at once.
  - Nest: one confirm sheet (deck's propose text: what happens, to which session and account), send box on a session row. Heavy actions (stop) need a second tap (prod_005 confirmations).
  - Reply and request fixtures; Go tests with the scripts replaced by test doubles.
- Out:
  - Dispatching a new `crow dispatch` run from Nest (missions dispatch is item_147; new work goes through missions).
  - Answering a permission prompt by key (needs req_027 to know what the session waits on). Send text covers a question.
  - `crow-review-run` actions, publish, merge.

# Acceptance criteria
- AC1: From the phone, text sent to a live session appears in that session's window and is submitted.
- AC2: Interrupt and stop act on the right pid; a pid not in `sessions` gets 404 and runs nothing.
- AC3: A POST without the CSRF header gets 403 and runs nothing.
- AC4: Every action leaves one event in the log with actor and surface.
- AC5: Stop asks twice in Nest; send shows the exact text and target before it runs.

# AC Traceability
- request-AC2 -> AC1, AC2, AC5. Proof: pilot sessions and runs from the phone with deck's confirmations.

# Decision framing
- Product framing: Not needed (prod_005 human in the loop, confirmations)
- Architecture framing: Not needed (adr_001 API: actions second, CSRF)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Depends on: `item_169_core_serves_every_deck_source_as_read_only_sections`, `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide` (guard and event shape; build together or right after)
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: piloting from the phone is the point of Nest (owner, 2026-11-10).

# Notes
- Confirm sheet visuals: mockup first if it differs from item_168's answer sheet; reuse that sheet when possible.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
