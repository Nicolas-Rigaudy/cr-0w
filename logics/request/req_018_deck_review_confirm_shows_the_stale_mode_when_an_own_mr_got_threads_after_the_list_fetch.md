## req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch - Deck review confirm shows the stale mode when an own MR got threads after the list fetch
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: the `r` confirm overlay names the mode from the MR list row, which can predate new threads.
- Keywords: deck, review, confirm, shows, stale, mode, own, got, threads, after, list, fetch
- Use when: touching the deck's review start (overlay.go reviewMode) or crow-review-run start modes.
- Skip when: work on findings display, apply or posting.

# Needs
- The confirm overlay (`overlay.go:213`) shows `reviewMode(r)` from the list row. When threads were posted after the list fetch it says `self`, while `crow-review-run start` now runs threads mode (b292d1d). Use the loaded detail's unresolved threads for the label.

# Context
- !213, 2026-11-03: the list had threads_open 0, threads posted a minute before `r`; a self review ran instead of answering them. The run side is fixed in b292d1d; only the label is left.

# Acceptance criteria
- AC1: On an own MR whose loaded detail has open threads, the confirm overlay says threads even when the list row says 0.
- AC2: Scope boundaries and operator impact are explicit.
- AC3: The request is ready to be promoted into a backlog slice.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics_manager/flow.py`
- `logics_manager/assist.py`
- `tests/python/test_logics_manager_cli.py`

# Backlog
- `item_089_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
