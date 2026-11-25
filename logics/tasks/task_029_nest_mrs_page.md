## task_029_nest_mrs_page - Nest MRs page
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: cdx
> Indicators reviewed: 2026-11-25 11:59:06

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, mrs, page
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_225_nest_mrs_page`

# Acceptance criteria
- AC1: Owner approved the MRs board first.
- AC2: Rows match `crow-mrs` output on live data for both forges.
- AC3: Each row action works from phone and laptop.
- AC4: The patterns panel lists the ledger patterns for the MR's project, newest hits first.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_029_nest_mrs_page.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_029_nest_mrs_page.md` after implementation.

# Validation
- (no validation recorded yet)
- command: `npx tsc --noEmit -p . && npx vitest run (tools/nest)` | result: tsc clean, 40 files / 311 tests passed | date: 2026-11-25 | note: Shots checked against board C, owner OK 2026-11-25. Follow-ups: review action via crow-review-run, key hints bar, Rerun CI.
- Finish workflow executed on 2026-11-25.
- Linked backlog/request close verification passed.

# Report
- Not started.
- Finished on 2026-11-25.
- Linked backlog item(s): `item_225_nest_mrs_page`
- Related request(s): `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# AC Traceability
- request-AC1 -> This task. Proof: owner OK to land 2026-11-25 after the shots in vault Research/nest-mrs/build-2026-11-25.
- request-AC2 -> This task. Proof: /mrs replaces the placeholder; every row action goes through a confirm sheet (vitest).
- request-AC3 -> This task. Proof: MRs is a page, not a waiting item; review asks keep reaching Inbox.
