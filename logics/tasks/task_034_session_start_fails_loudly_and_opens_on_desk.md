## task_034_session_start_fails_loudly_and_opens_on_desk - Session start fails loudly and opens on Desk
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 95%
> Confidence: 90%
> Progress: 80%
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
- `go test ./...` in tools/crow-core: 327 passed.
- `pytest tests/test_crow_focus.py`: 35 passed.
- tools/nest `npx tsc --noEmit` clean, `npx vitest run`: 411 passed.

# Report
- AC1 cause (from code, no live run): `crow-accounts pick` picks an account that already has a live session, by design. `crow open` still passes `--on-busy fail` to cdx. cdx then refuses any account with a running session (`session_busy`) and exits. The tmux session closed with it and its output was lost.
- 81c961c: headless sessions keep their pane when the command exits.
- 005ab08 (AC2, AC3): crow_term waits up to 3 s (`CROW_TERM_SETTLE`). A pane that died returns `{ok:false, code:2}` with the exit status and last lines. Core answers 502 with that text. The Nest sheet shows it on several lines.
- 4cd95d4 (AC4): starts from MRs and Agenda go to `/desk?tmux=<name>`. The Desk opens that session once it shows up.
- Open: owner-prefs say a busy account must never refuse a launch. cdx has no mode that allows a second session without a prompt. Owner decision pending.

# Links
- Request: `req_071_nest_mrs_review_flow_and_visible_session_starts`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
