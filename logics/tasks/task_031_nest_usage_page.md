## task_031_nest_usage_page - Nest usage page
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90
> Confidence: 80
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-11-25 16:59:19

# AI Context
- Summary: Build the Nest Usage page (layout C board) on top of a per-day, per-account usage series from crow-cost-audit and limit samples.
- Keywords: nest, usage, page, crow-cost-audit, limit-samples, crow-core
- Use when: building or extending the Nest Usage page or its data.
- Skip when: quota history charts (item_297), account picking (item_027).

# Context
- Spec: the "Owner picks (2026-11-25)" and "units per day" sections of item_227.
- Board: tools/nest/scripts/mockups/usage.html?v=C and ?f=phone.

# Plan
- [x] 1. Core: crow-cost-audit --json gains a per-day series (cost and tokens per model) and a split per account, with tests.
- [x] 2. Core: per-day weekly % used and 5h windows per account from limit-samples.jsonl (weekly reset, 5h split), with tests.
- [x] 3. crow-core: read-only endpoint serving the usage series to Nest, with Go tests.
- [x] 4. Nest: /usage and /usage/<account> routes, header gauge click opens it.
- [x] 5. Nest: layout C (account list, facts, 5h and week charts) matching the board, laptop and phone.
- [x] 6. Nest: ANALYSIS pane (D/W/M, PER DAY with $ / tokens / limits, PER PROJECT, PER MODEL, PER KIND, PER SESSION, TOP RUNS), cache savings and week forecast tiles.
- [x] 7. Shots into vault Research/nest-design/item_227/live-*.png, board compare.
- [x] GATE: tools/nest `npx tsc --noEmit -p .` and `npx vitest run`; tools/crow-core `go test ./...` when touched.

# Backlog
- `item_227_nest_usage_and_token_analysis_page_from_the_account_gauges`

# Definition of Done (DoD)
- [x] Code is implemented and reviewed.
- [x] Validation passes.
- [x] Linked docs are synchronized.

# AC Traceability
- backlog-AC1 -> This task. Proof: board C picked by the owner 2026-11-25.
- backlog-AC2 -> This task. Proof: pending.
- backlog-AC3 -> This task. Proof: pending.
- request-AC1 -> This task. Proof: Usage page merged in e156284 and served; owner approved the shots against board C and phone board A; validated with tsc, vitest (378 passed) and go test. Source: `e156284`
- request-AC2 -> This task. Proof: Usage page merged in e156284 and served; owner approved the shots against board C and phone board A; validated with tsc, vitest (378 passed) and go test. Source: `e156284`
- request-AC3 -> This task. Proof: Usage page merged in e156284 and served; owner approved the shots against board C and phone board A; validated with tsc, vitest (378 passed) and go test. Source: `e156284`

# Validation
- command: `python3 -m pytest -q tests/ -k "usage or cost or afk"` | result: 63 passed, 1 skipped | date: 2026-11-25 | note: steps 1-2.
- command: `go test ./...` (tools/crow-core) | result: 311 passed in 14 packages | date: 2026-11-25 | note: step 3.
- command: `npx tsc --noEmit -p . && npx vitest run (378 passed); go test ./...` | result: passed | date: 2026-11-25
- Finish workflow executed on 2026-11-25.
- Linked backlog/request close verification passed.

# Report
- Steps 1-2 done 2026-11-25: `by_account` and `by_day` (per model, account, project, kind, plus `limits`) in crow-cost-audit --json; `limit_days()` and `trim_samples()` in crow-usage.
- limit-samples.jsonl was trimmed to the last 6h on every write, so no per-day history existed. Retention is now 6h full, then one sample per account per 5 minutes up to 35 days. Limits history starts when this lands on main; earlier days show "no data".
- Step 3 done 2026-11-25: `GET /api/usage` (internal/api/usage.go) answers `{version, cost, limits, sessions, generated_at, stale, error}`. `cost` is the `crow-cost-audit --json --days 35` body as a section (5 min TTL, own ext, not in /api/state, refreshed only when the page asks). `limits` and `sessions` come from the usage section; sessions are slimmed to the fields the page reads (no paths, no activity). 404 "module off" with the usage brick off.
- Step 4 code is written but not committed: the Nest gate cannot run in this worktree (no node_modules; linking or installing it needs owner approval).
- Finished on 2026-11-25.
- Linked backlog item(s): `item_227_nest_usage_and_token_analysis_page_from_the_account_gauges`
- Related request(s): `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
