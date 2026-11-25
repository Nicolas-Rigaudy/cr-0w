## task_027_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work - Nest Desk opens on the last selected session and groups live and background work
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-25 10:46:33

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, desk, opens, last, selected, session, groups, live, background, work
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_210_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work`

# Acceptance criteria
- AC1: Pick a session, switch to Ops and back: the same session is open.
- AC2: Pick a run, reload: the same run is open.
- AC3: No remembered pick (or it ended): the session that needs the owner first, else the working one; never a background run.
- AC4: Only background runs live: "Nothing open" panel.
- AC5: vitest covers `deskPick` (remembered live, remembered gone, needs before working, runs skipped, storage blocked); tsc and vitest pass.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_027_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_027_nest_desk_opens_on_the_last_selected_session_and_groups_live_and_background_work.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- Not started.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
