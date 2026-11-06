## item_075_runs_cost_per_run - crow runs shows cost per run
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Show cdx cost_usd per run in `crow runs` and the deck Agents page.
- Keywords: cost_usd, runs, budget, deck agents
- Use when: surfacing run cost.
- Skip when: other req_019 slices.

# Problem
- crow caps runs with --budget but never shows what a run cost.

# Scope
- In:
  - Read cost_usd from `cdx runs --json` / run-report for Claude runs.
  - `crow runs` listing: a cost column for finished runs, budget next to it when the run had one.
  - Deck Agents page: cost on the run line and in the lane detail.
  - Daily total of background run cost in `crow runs` footer.
- Out:
  - Cost of interactive sessions (crow-usage covers it).

# Acceptance criteria
- AC1: Finished Claude runs show cost to the cent; runs without cost show nothing, not 0.
- AC2: The deck Agents page shows the same figure.
- AC3: Tests cover parsing with and without cost_usd.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1 to AC3 above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner picked it from the 2026-11-03 release review.

# Notes
- Check the exact field name and where cdx writes it (runs --json vs run-report) before building.

# Delivery
- Done 2026-11-06. cdx `runs --json` carries a top-level `cost_usd` (null when unknown); `crow runs` shows `$0.71 of $25` per finished run (budget from the ledger, stored at dispatch) and a `background runs cost $N` footer; unknown cost shows nothing, never 0.
- Deck Agents page: finished run Meta and detail tail show the same figure (`tools/crow-deck/data.go`).
- Proof: tests/test_crow_runs_logics_cost.py and TestFinishedRunCost; pytest 303 passed, Go 141 passed.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
