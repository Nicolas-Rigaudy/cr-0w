## item_272_ops_queue_missions_keep_a_position_reorder_and_run_now - Ops queue: missions keep a position, reorder and run now
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Queued missions get a stored position. The Ops queue reorders them (arrows and drag) and RUN opens a menu: next in line or start now over cap.
- Keywords: ops, queue, missions, position, reorder, run, now
- Use when: Building mission queue order, reorder routes or the RUN menu on Ops.
- Skip when: Routine or AFK queue rows (unchanged).

# Problem
- The mission tick (tools/crow-core/internal/mission/tick.go Tick) admits queued missions by Created then ID while fewer than maxRuns run. The owner cannot reorder the queue or start one now.
- Design approved 24/11: vault `CR-0W/Research/nest-design/item_190/queue/compare.html`, option A; decisions in `item_190/design-notes.md`.

# Scope
- In:
  - Core: a stored queue position on missions; tick admits by position, then Created. New missions go to the end.
  - Core routes next to the existing POST /api/missions/{id}/cancel (internal/api/actions.go): move (up, down, to index), next (to the head), start now over cap (a run past maxRuns, recorded as over cap).
  - Nest Ops queue (src/screens/Ops.tsx): each mission row gets a drag handle and up/down arrows on laptop and phone (long-press drag on phone); RUN opens a two-line menu, 'Next in line' highlighted by default, 'Start now · over cap' in the warn colour; slots chip shows 3/2 while over cap.
  - Tests for tick order and each route; fixtures for the Ops queue.
- Out:
  - Routine and AFK rows. Priority classes. Changing maxRuns.

# Acceptance criteria
- AC1: Tick starts queued missions by stored position, then age (unit test).
- AC2: Up, down, drag and 'Next in line' change the order on live data, laptop and 368x812, and the next tick follows it.
- AC3: 'Start now · over cap' starts a third run at once; the queue shows slots 3/2 until one ends.
- AC4: The menu, arrows and drag match queue/compare.html option A in both themes.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: owner-approved Ops design, built and checked on live data.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)
- Request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`, `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Primary task(s): `task_023_orchestrate_nest_laptop_parity`
- Depends on: `item_190_laptop_missions_from_the_approved_missions_board`

# Priority
- Priority: Medium
- Rationale: Finishes the Ops page (item_190). Build order: item_272, item_273, item_274; each in a fresh session.

# Notes
- Read the vault design notes and the compare page before building; build from the chosen option only.

- Proof AC1: `internal/mission/queue_test.go` TestQueueCompare, TestTickFollowsQueueOrder.
- Proof AC2: scratch core on live routes (24/11): ▼ at 1440, long-press drag and Next in line at 368 reordered Alpha/Beta/Gamma; `crow-core tick` with max 1 started the head. Fixture checks: `tools/nest/scripts/queue-shots.mjs` (36 checks, both themes).
- Proof AC3: start now on the scratch core dispatched `--over-cap` (2 of 1 busy); TestTickStartNowOverCap; slots 3/2 shown in queue-shots `built-overcap-*`.
- Proof AC4: shots in vault `item_190/queue/built/` next to queue-A.png; cyber check noted in design-notes.md.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-24T10:48:30.883207+00:00 | hub-session | 20261124-item272 | Queue position, reorder, run now | succeeded |  | Go 302 pass, pytest 848 pass, vitest 270 pass, queue-shots 36 checks; live scratch core: arrows/drag/next reorder, tick follows, start now dispatched over cap |
