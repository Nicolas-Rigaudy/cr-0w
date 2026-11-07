## item_103_take_control_of_a_mission_as_a_logged_manual_gesture - Take control of a mission as a logged manual gesture
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Safety
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:49

# AI Context
- Summary: Take control stops the headless run, opens an interactive session in the same worktree, logs `manual_gesture` start and end with heads before and after, and returns to auto with a fresh isolated run.
- Keywords: take, control, mission, logged, manual, gesture
- Use when: Building the escape hatch in core and the Live tab.
- Skip when: Read-only attach to a live run (prod_005 Live tab).

# Problem
- Isolated runs are headless; the owner needs a way in that leaves a trace and keeps publishing through core.

# Scope
- In:
  - Mission sub-state `manual` with start and end events.
  - Interactive session through the existing open or handoff path, under the owner's own permission mode plus the mission guard.
  - Commits in the window tagged as owner work.
  - Return to auto: fresh isolated run with the brief and the owner's notes.
- Out:
  - Live terminal attach to isolated runs (not possible headless).

# Acceptance criteria
- AC1: Take control writes `manual_gesture` start and end events with worktree heads.
- AC2: Any pending validation record is dropped on take control.
- AC3: Publishing after take control still goes through core.
- AC4: Return to auto starts a fresh run carrying the owner's notes.
- AC5: A commit made in the manual window goes through the full merge gate; sensitive paths still ask.

# AC Traceability
- request-AC2 -> escape hatch in the layered design.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Needed once missions run isolated; manual worktree work covers it until then.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- Owner decision 2026-11-07 (req_030 decision 6): the merge gate still runs on take-control commits.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
