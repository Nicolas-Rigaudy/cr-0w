## task_034_session_start_fails_loudly_and_opens_on_desk - Session start fails loudly and opens on Desk
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
> Indicators reviewed: 2026-11-26 14:36:14

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: session, start, fails, loudly, opens, desk
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_298_session_start_fails_loudly_and_opens_on_desk`

# Acceptance criteria
- AC1: Reproduce first: find why the 11-26 starts exited (run `crow open --on-busy fail lambda` the way core does, capture output). Record the cause in Notes.
- AC2: A start that exits within the wait returns non-2xx with the reason; a Go test covers it with a fake command that exits.
- AC3: Nest shows the reason next to the action; vitest covers it.
- AC4: A start that lives lands on that session in Desk.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_034_session_start_fails_loudly_and_opens_on_desk.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_034_session_start_fails_loudly_and_opens_on_desk.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- Not started.

# Links
- Request: `req_071_nest_mrs_review_flow_and_visible_session_starts`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
