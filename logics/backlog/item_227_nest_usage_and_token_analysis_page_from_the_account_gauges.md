## item_227_nest_usage_and_token_analysis_page_from_the_account_gauges - Nest usage and token analysis page from the account gauges
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 95
> Confidence: 85
> Progress: 5%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

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
- req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits-AC -> This backlog slice. Proof: Owner approved the Usage board first.

# Dependencies
- Header gauges 0a925df. crow-cost-audit.

# Links
- Request: req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits.

# Priority
- Priority: Medium
- Rationale: Owner asked for it; helps the token work (item_213).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
