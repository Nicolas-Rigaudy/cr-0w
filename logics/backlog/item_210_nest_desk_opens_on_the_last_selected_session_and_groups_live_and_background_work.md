## item_210_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work - Nest Desk opens on the last selected session and groups live and background work
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90
> Confidence: 80
> Progress: 10%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 10:46:33

# AI Context
- Summary: Returning to Desk jumps to the first background run. Desk forgets the pick on every tab change and its default puts runs first. Remember the last pick per device; else open the top interactive session; never default to a background run.
- Keywords: nest, desk, selection, last pick, sessions, runs
- Use when: changing how Desk picks its open session.
- Skip when: Desk layout (shipped as item_246), work list rows (item_204).

# Problem
- Seen 2026-11-16 and again 2026-11-25.
- `App.tsx:45` mounts `Desk` only on the desk route, so a tab change drops `sel` (`Desk.tsx:97`, plain `useState`, seeded only from `?pid=`).
- Default pick at `Desk.tsx:108`: `runs[0] ?? working session ?? agents[0]`. A background run comes first.
- The layout half of the first version shipped as item_246 (layout B).

# Scope
- In:
  - Remember the picked pid per device in localStorage (`nest.desk.pick`), guarded like `data/useScope.ts`. Write on every pick and on `?pid=` entry.
  - Default pick, in order: the remembered pid if still live, a run included (owner 2026-11-25: restore what the owner picked); else the first `kind: session` item of the work feed that maps to a live agent; without the feed, waiting, then working, then idle sessions, newest first; else the "Nothing open" panel.
  - Desk never picks a background run on its own.
  - The pure choice lives in `data/deskPick.ts` with unit tests; `useDesk` calls it.
- Out:
  - Layout (item_246). Phone cover routing (`/desk/<pid>`).

# Acceptance criteria
- AC1: Pick a session, switch to Ops and back: the same session is open.
- AC2: Pick a run, reload: the same run is open.
- AC3: No remembered pick (or it ended): the session that needs the owner first, else the working one; never a background run.
- AC4: Only background runs live: "Nothing open" panel.
- AC5: vitest covers `deskPick` (remembered live, remembered gone, needs before working, runs skipped, storage blocked); tsc and vitest pass.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`. Builds on item_246. Same remount as item_277.

# Priority
- Priority: High
- Rationale: the wrong default hides the live session.

# Notes
- Raised 2026-11-16; rescoped 2026-11-25 from plan mission 20261125-plan-nest-desk-bugs (owner: reuse this item for bug 1).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T09:46:18.243154+00:00 | cdx | dcfbfb5d257a441a9de8bb0a683ea0e1 | item_210 desk pick | running |  |  |


# Tasks
- `task_027_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work`
