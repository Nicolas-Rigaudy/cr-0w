## item_089_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch - Deck review confirm shows the stale mode when an own MR got threads after the list fetch
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: High
> Theme: Operator workflow and runtime integration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: deck, review, confirm, shows, stale, mode, own, got, threads, after, list, fetch
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
The confirm overlay (`overlay.go:213`) shows `reviewMode(r)` from the list row. When threads were posted after the list fetch it says `self`, while `crow-review-run start` now runs threads mode (b292d1d). Use the loaded detail's unresolved threads for the label.

# Scope
- In:
  - one coherent delivery slice from the source request
- Out:
  - unrelated sibling slices that should stay in separate backlog items instead of widening this doc

# Acceptance criteria
- AC1: On an own MR whose loaded detail has open threads, the confirm overlay says threads even when the list row says 0.
- AC2: Scope boundaries and operator impact are explicit.
- AC3: The request is ready to be promoted into a backlog slice.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: `reviewMode(r, d)` counts the loaded detail's unresolved threads; the overlay and the run start both use it. Test `TestReviewConfirmAndStartUseThreadsPostedAfterTheListFetch` (tools/crow-deck/threads_test.go).
- request-AC2 -> This backlog slice. Proof: scope is the label and the `--mode` the deck passes; crow-review-run already corrects the mode on its side.
- request-AC3 -> This backlog slice. Proof: promoted to this item on 2026-11-05.

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `logics/request/req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Default until groomed.

# Notes
- Hybrid rationale: Derived from request `req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch.md`.
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
