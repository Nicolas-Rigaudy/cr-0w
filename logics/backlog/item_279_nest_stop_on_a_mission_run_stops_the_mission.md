## item_279_nest_stop_on_a_mission_run_stops_the_mission - Nest Stop on a mission run stops the mission
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
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
- 2026-11-25 (8b05af6, 0472767): Stop on a mission run (Desk and session chat) opens one sheet: "Stop and cancel the mission" calls `POST /missions/{id}/cancel`; "Restart run" calls the existing `/stop` (kills the run, the tick resumes it). Cancel sheet words a running mission. A reply with `ok: false` now throws and stays in the sheet. Go test: cancel on a running mission stops its run and later cancel/stop get 409.
- Left: run `stop-probe.mjs` (needs approval in an unattended run) and the live repro from Nest on a real running mission (AC1, AC2).
- Open: core's cancel ignores a failed run stop (the warning is dropped) and still cancels; the UI shows success.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T10:27:57.900812+00:00 | cdx | 470beddb074d46aaba6bf86687d56146 | nest stop build item_279 | running |  |  |
| 2026-11-25T09:51:21.811074+00:00 | cdx | b94fd623d5f44fa69c8e2263c7a86d83 | item_279 nest stop | succeeded |  |  |

# Tasks
- `task_027_nest_stop_on_a_mission_run_stops_the_mission`
