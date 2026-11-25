## task_030_nest_reports_page_for_settled_afk_reports_and_ritual_outputs - Nest Reports page for settled AFK reports and ritual outputs
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 95%
> Confidence: 85%
> Progress: 85%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Core `reports` section (settled AFK summaries, weekly rollups, retros) and the Nest Reports page from board B: six ritual tiles, outputs newest first, Markdown reader, Generate on the retro tile only.
- Keywords: nest, reports, page, settled, afk, ritual, outputs, retro, generate
- Use when: changing the Reports page, the core reports section or what counts as settled.
- Skip when: open decisions (Inbox), single run reports (Ops), scheduling the other rituals (item_234).

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

# Decisions
- Owner 2026-11-25: an AFK report is settled when `summary-<date>.json` has an empty `decisions[]` or every decision is in `decided[]`. Held reports show as a count pointing to Inbox.
- Owner 2026-11-25: keep the six tiles of board B. Brief, daily and eod show "no saved output". Weekly shows its rollups with "no routine yet".
- Owner 2026-11-25: Generate only on the retro tile (`POST /api/routines/crow-retro.timer/run`), shown only when the routine exists in `jobs`. Generate for brief, daily, weekly and eod waits for item_234 (routines). AC3 is met for retro only until then.

# Validation
- 2026-11-25, tools/crow-core: `go test ./...` passed (313 tests, 14 packages).
- 2026-11-25, tools/nest: `npx tsc --noEmit -p .` no errors; `npx vitest run` 333/333 in 43 files.
- Shots (`scripts/shots.mjs --fixtures --only=reports,reports-retro`) not taken yet: the background run had no approval for the command.

# Report
- 2026-11-25 data check, build paused on owner decisions:
  - AFK reports are not missions. They live in `local/state/afk/summary-<date>.json` (projects with `decisions[]` and `decided[]`) and `reports/<run>.md`. No inbox row exists for their decisions, so "settled = no open inbox row" has nothing to match.
  - Ritual outputs saved in the vault: weekly (`Rollups/<week> rollup.md`) and retro (`Retros/*.md`), both resolved through `crow_lib.crow_path`. Brief prints in chat only. Eod files answers into the daily note. Daily posts to Confluence (example skill only). Core has no endpoint for any of them.
  - Background start: only retro has a routine (`crow-retro.timer`), started by core `POST /api/routines/{unit}/run`. Brief, daily, weekly and eod have no routine (item_234). Eod and daily need the owner (questions, approval).
- 2026-11-25 build:
  - ab1ae7e core `reports` section: settled AFK summaries, rollups and retros read through the notes dir.
  - caf36f0 Nest Reports page: tiles, list, reader, phone reader page, retro Generate with a 409 refresh, fixtures and tests.
  - 1c0b8fc shots entries `reports` and `reports-retro`.
  - Left: laptop and phone shots into the vault `Research/nest-design/item_258/`, compare with board B, fix CSS gaps.

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
