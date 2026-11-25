## task_030_nest_reports_page_for_settled_afk_reports_and_ritual_outputs - Nest Reports page for settled AFK reports and ritual outputs
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 60%
> Progress: 5%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, reports, page, settled, afk, ritual, outputs
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_258_nest_reports_page_for_settled_afk_reports_and_ritual_outputs`

# Acceptance criteria
- AC1: Owner approved the Reports board before the build.
- AC2: An AFK report with an open decision is in Inbox and not here; once answered it shows here.
- AC3: Each Generate button starts its ritual and the output appears when done.
- AC4: Works on 368x812 and on the laptop.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_030_nest_reports_page_for_settled_afk_reports_and_ritual_outputs.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_030_nest_reports_page_for_settled_afk_reports_and_ritual_outputs.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- 2026-11-25 data check, build paused on owner decisions:
  - AFK reports are not missions. They live in `local/state/afk/summary-<date>.json` (projects with `decisions[]` and `decided[]`) and `reports/<run>.md`. No inbox row exists for their decisions, so "settled = no open inbox row" has nothing to match.
  - Ritual outputs saved in the vault: weekly (`Rollups/<week> rollup.md`) and retro (`Retros/*.md`), both resolved through `crow_lib.crow_path`. Brief prints in chat only. Eod files answers into the daily note. Daily posts to Confluence (example skill only). Core has no endpoint for any of them.
  - Background start: only retro has a routine (`crow-retro.timer`), started by core `POST /api/routines/{unit}/run`. Brief, daily, weekly and eod have no routine (item_234). Eod and daily need the owner (questions, approval).

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
