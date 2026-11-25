## task_031_nest_usage_page - Nest usage page
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90
> Confidence: 80
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Build the Nest Usage page (layout C board) on top of a per-day, per-account usage series from crow-cost-audit and limit samples.
- Keywords: nest, usage, page, crow-cost-audit, limit-samples, crow-core
- Use when: building or extending the Nest Usage page or its data.
- Skip when: quota history charts (item_297), account picking (item_027).

# Context
- Spec: the "Owner picks (2026-11-25)" and "units per day" sections of item_227.
- Board: tools/nest/scripts/mockups/usage.html?v=C and ?f=phone.

# Plan
- [ ] 1. Core: crow-cost-audit --json gains a per-day series (cost and tokens per model) and a split per account, with tests.
- [ ] 2. Core: per-day weekly % used and 5h windows per account from limit-samples.jsonl (weekly reset, 5h split), with tests.
- [ ] 3. crow-core: read-only endpoint serving the usage series to Nest, with Go tests.
- [ ] 4. Nest: /usage and /usage/<account> routes, header gauge click opens it.
- [ ] 5. Nest: layout C (account list, facts, 5h and week charts) matching the board, laptop and phone.
- [ ] 6. Nest: ANALYSIS pane (D/W/M, PER DAY with $ / tokens / limits, PER PROJECT, PER MODEL, PER KIND, PER SESSION, TOP RUNS), cache savings and week forecast tiles.
- [ ] 7. Shots into vault Research/nest-design/item_227/live-*.png, board compare.
- [ ] GATE: tools/nest `npx tsc --noEmit -p .` and `npx vitest run`; tools/crow-core `go test ./...` when touched.

# Backlog
- `item_227_nest_usage_and_token_analysis_page_from_the_account_gauges`

# Definition of Done (DoD)
- [ ] Code is implemented and reviewed.
- [ ] Validation passes.
- [ ] Linked docs are synchronized.

# AC Traceability
- backlog-AC1 -> This task. Proof: board C picked by the owner 2026-11-25.
- backlog-AC2 -> This task. Proof: pending.
- backlog-AC3 -> This task. Proof: pending.

# Validation
- (no validation recorded yet)

# Report
- Not started.

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
