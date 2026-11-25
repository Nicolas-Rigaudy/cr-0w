## item_276_nest_afk_queue_runner_view - Nest shows the AFK queue runner with pause and needs-you rows
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
- Summary: Nest Ops shows the AFK runner: on or paused, each row's state (queued, running, landed, needs you, branch ready), the account and floor that holds a row back, with pause, resume and retry land actions.
- Keywords: nest, ops, afk, queue, runner, pause
- Use when: the Nest view of the AFK queue.
- Skip when: the runner logic itself (item_233, item_275).

# Problem
- From Nest the owner cannot see why the queue is not moving or act on it without a Claude session.

# Scope
- In: runner state and per-row state in the existing Ops queue panel (item_190, item_272); pause/resume; retry land on a needs-you row; the reason a row waits (floor, cap, reset time).
- Out: editing job prompts (already in Ops).

# Acceptance criteria
- AC1: Ops shows each row's state and the waiting reason.
- AC2: Pause stops new launches without stopping running runs; resume restarts them.
- AC3: Retry land reruns item_275's land step for one row.

# AC Traceability
- req_067-AC2 -> AC1. req_067-AC3 -> AC2.

# Dependencies
- item_233, item_275.

# Links
- Request: req_067_afk_queue_runs_itself.

# Priority
- Priority: High
- Rationale: Nest-only use of crow.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
