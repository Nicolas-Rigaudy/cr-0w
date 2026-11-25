## item_225_nest_mrs_page - Nest MRs page
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85
> Confidence: 80
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 11:59:06

# AI Context
- Summary: MRs tab on live `/api/state/mrs`: mine, to review, CI state, threads; each row opens review, threads or the MR on the forge. Optional brick (req_064 modules).
- Keywords: nest, mrs, prs, review, ci, forge, threads
- Use when: building the MRs tab.
- Skip when: the review run itself (deck items 043-048).

# Problem
- The MRs tab shows a placeholder. Core already serves `/api/state/mrs` (crow-mrs) for the deck.

# Scope
- In:
  - Sections: To review (unblock devs first, per priorities), Mine (CI, approvals, threads open), Merged today.
  - Row actions: review (dispatch review run), answer threads (review-threads session), open on forge, open repo session.
  - Mockups first (phone, laptop), reuse the deck MR board layout ideas (item_049/050).
  - Review patterns panel: top patterns from the review ledgers (local/reviews/common.md and the project ledger), with counts and the last MR that hit each one.
  - Hidden when the forge module is off (req_064).
- Out:
  - Posting comments from Nest (stays a gated action in the review flow).

# Acceptance criteria
- AC1: Owner approved the MRs board first.
- AC2: Rows match `crow-mrs` output on live data for both forges.
- AC3: Each row action works from phone and laptop.
- AC4: The patterns panel lists the ledger patterns for the MR's project, newest hits first.

# AC Traceability
- req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits-AC -> This backlog slice. Proof: Owner approved the MRs board first.
- request-AC1 -> This backlog slice. Proof: owner OK to land 2026-11-25 after the shots in vault Research/nest-mrs/build-2026-11-25.
- request-AC2 -> This backlog slice. Proof: /mrs replaces the placeholder; every row action goes through a confirm sheet (vitest).
- request-AC3 -> This backlog slice. Proof: MRs is a page, not a waiting item; review asks keep reaching Inbox.

# Dependencies
- Module flag from the module design item (can build before, flag added after).

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`

- Primary task(s): `task_029_nest_mrs_page`
# Priority
- Priority: Medium
- Rationale: Reviews rank second in the owner's priorities.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T10:33:00.906115+00:00 | cdx | 9b5f90d9a4a84f2ba941e9d7ebe77ec9 | nest mrs build (option C) | succeeded |  |  |

# Decision framing
- Owner 2026-11-25: none of A/B/C. Redo the mockups: closer to the deck, less clutter. MR list as clear as A or B, details like A but bigger, patterns small and expanded on demand. Do not reserve space for the chat panel.
- Round 2 (2026-11-25): options A (list + big detail), B (who-acts strip + wide detail), C (deck panes + tabbed detail) in `tools/nest/scripts/mockups/mrs-v2.html`, rendered by `tools/nest/scripts/mrs-v2-mockups.mjs`. Waiting on the owner's pick (AC1).
- 2026-11-25: owner picked round 2 option C (deck panes + tabbed detail). Approved board: `mrs-v2.html` option C, PNGs in the vault `CR-0W/Research/nest-mrs/round2/mrs-C-{laptop,laptop-patterns,phone}.png`. AC1 met.

# Tasks
- `task_029_nest_mrs_page`

# Notes
- Task `task_029_nest_mrs_page` was finished via `logics-manager flow finish task` on 2026-11-25.
