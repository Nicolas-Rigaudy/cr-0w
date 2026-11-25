## item_227_nest_usage_and_token_analysis_page_from_the_account_gauges - Nest usage and token analysis page from the account gauges
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 85
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 16:59:18

# AI Context
- Summary: Clicking an account gauge in the Nest header opens a Usage page: per account 5h and week windows, burn rate, reset times, cost per run, per project and per model, and the token savers (cache hit, hub start context) from crow-cost-audit.
- Keywords: nest, usage, tokens, accounts, gauges, cost, burn rate, crow-cost-audit
- Use when: building the Usage page.
- Skip when: account picking logic (item_027, item_095).

# Problem
- Deck usage features (item_003, item_019, item_061, item_075) are not in Nest. The header gauges (0a925df) show levels but open nothing.

# Scope
- In:
  - Route `/usage` and `/usage/<account>`; header gauge click opens it.
  - Per account: 5h and week windows with reset time, burn rate, projected end of window, floor lines (req_064 quota floors).
  - Analysis: tokens and cost per day, per project, per model, per run kind; top runs; cache hit; hub start context. Data from `/api/state/usage` and `bin/crow-cost-audit --json` (add the JSON output if missing).
  - Charts follow the dataviz rules; mockups first.
- Out:
  - Changing quotas or accounts from the page.

# Acceptance criteria
- AC1: Owner approved the Usage board first.
- AC2: Clicking a gauge opens that account's page; numbers match `cdx status` and crow-cost-audit for the same window.
- AC3: Charts readable on 368x812.

# Owner picks (2026-11-25)
- Layout C: account list with 5h/week sparklines on the left (plus the 7-day tiles), the selected account on the right: facts row, 5h window chart and week chart side by side, then an ANALYSIS pane with one chart at a time.
- Board: tools/nest/scripts/mockups/usage.html (?v=C), shots in vault Research/nest-design/item_227 (usage-C.png, usage-phone-A.png for the phone).
- Owner: "I really like the graphs". Keep every chart with labelled axes, legend, floor line, projection dashed.
- Borrow from the deck usage tab (item_003, item_019, item_061, item_075):
  - period switch D / W / M on the analysis;
  - analysis tabs PER DAY, PER PROJECT, PER MODEL, PER KIND, PER SESSION (project, full title, model, tokens, cost), TOP RUNS;
  - burn sparkline per account in the list;
  - cache savings and a cost forecast to the end of the week as tiles.
- Core in this build: crow-cost-audit --json gains a per-day series and a split per account (cdx profile of each transcript). Quota history stays in item_297; until then the window charts show used-to-now as one point plus the projection.
- Phone: account picker, 5h and week charts stacked; analysis one chart per tab.

## Owner pick: units per day (2026-11-25)
- PER DAY gets a $ / tokens / limits toggle. Limits = % of weekly used per day by account, plus 5h windows as filled cells ("3 full + 10%").
- Core derives these from limit-samples.jsonl (history from 2026-11-24).

# AC Traceability
- `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits` AC1 -> This backlog slice. Proof: built on nest-usage-build, merged in e156284 and served; owner OK on the shots 2026-11-25.
- `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits` AC2 -> This backlog slice. Proof: built on nest-usage-build, merged in e156284 and served; owner OK on the shots 2026-11-25.
- `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits` AC3 -> This backlog slice. Proof: built on nest-usage-build, merged in e156284 and served; owner OK on the shots 2026-11-25.
- request-AC1 -> This backlog slice. Proof deferred to slice closeout.
- request-AC2 -> This backlog slice. Proof deferred to slice closeout.
- request-AC3 -> This backlog slice. Proof deferred to slice closeout.

# Dependencies
- Header gauges 0a925df. crow-cost-audit.

# Links
- Request: `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`

- Primary task(s): `task_031_nest_usage_page`
# Priority
- Priority: Medium
- Rationale: Owner asked for it; helps the token work (item_213).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T15:54:10.050516+00:00 | cdx | 35d8d0b0d43740a6afe7c98e1b21f792 | Nest Usage run tabs (item_227) | succeeded |  |  |
| 2026-11-25T15:49:23.337384+00:00 | cdx | 4ee772c5d9014260b172db1e3f7df7cf | Nest Usage phone fixes (item_227) | succeeded |  |  |
| 2026-11-25T15:46:39.147233+00:00 | cdx | c1e749d74f3d4639b4cc1c04191f88b2 | Nest Usage build 4 (item_227) | succeeded |  |  |
| 2026-11-25T15:37:10.353519+00:00 | cdx | 031b2ae96a1e44fbb281fc82ee61ba00 | Nest Usage build 3 (item_227) | cancelled |  | stopped |
| 2026-11-25T15:36:47.325046+00:00 | cdx | 5c42e25079a34599984c3847a22ace6e | Nest Usage build 3 (item_227) | failed |  | limit |
| 2026-11-25T15:08:57.054125+00:00 | cdx | b2b1795a99ea40dfa50fa668fee47111 | Nest Usage build 2 (item_227) | succeeded |  |  |


# Tasks
- `task_031_nest_usage_page`

# Notes
- Task `task_031_nest_usage_page` was finished via `logics-manager flow finish task` on 2026-11-25.
