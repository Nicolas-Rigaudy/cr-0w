## item_279_nest_stop_on_a_mission_run_stops_the_mission - Nest Stop on a mission run stops the mission
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 80%
> Confidence: 70%
> Progress: 10%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 11:02:00

# AI Context
- Summary: Stop in Nest on a mission run did nothing (2026-11-25, no event logged), and mission Stop only kills the run before the tick resumes it. Make Stop end the mission and find why the press never reached core.
- Keywords: nest, stop, cancel, mission, run, tick, resume
- Use when: stop and cancel controls for missions and runs.
- Skip when: butler tools (done in 650b35b).

# Problem
- The owner pressed Stop on the running mission review lambda!242; the mission log has no stop action and the run kept going.
- The Stop sheet needs a second tap ("Tap again to stop"); a single tap may look like nothing happened. Unconfirmed.
- `mission.Env.Stop` kills the live run only. The note says "The mission moves on at the next tick": the tick resumes it up to `budget.resumes`.
- Desk offers cancel only on queued missions (`Desk.tsx:217-230`).

# Scope
- In:
  - Reproduce the press from Desk and from the session chat menu (`Chat.tsx:195`) on a running mission; log every action request in crow-core (method, path, status).
  - Stop on a mission run cancels the mission (one confirm, plain words: "Stop and cancel the mission"). Keep a separate "Restart run" action (owner, 2026-11-25): it kills the run and lets the tick resume it, as Stop does today.
  - Cancel offered on running missions, not only queued ones.
  - A failed action shows its error in place; never silent.
- Out:
  - Butler tools (done).

# Acceptance criteria
- AC1: One Stop on a running mission ends its run and the mission shows cancelled; the tick does not restart it.
- AC2: Every Nest action reaches core's log with its status.
- AC3: A refused or failed action shows its message on screen.
- AC4: Tests for the sheet and the handler; tsc, vitest and go test pass.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`.

# Priority
- Priority: High
- Rationale: an unstoppable run costs budget and sends the owner back to the terminal.

# Notes
- 2026-11-25: stopped by hand with `crow-runs stop` then `crow mission cancel`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T09:46:16.709841+00:00 | cdx | b94fd623d5f44fa69c8e2263c7a86d83 | item_279 nest stop | running |  |  |


# Tasks
- `task_027_nest_stop_on_a_mission_run_stops_the_mission`
