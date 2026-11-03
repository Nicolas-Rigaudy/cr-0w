## req_019_release_leverage_cdx_0_25_logics_2_24 - Use what cdx 0.25 and logics-manager 2.24 added
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-03 17:24:35

# AI Context
- Summary: Four small slices that put cdx 0.24-0.25 and logics-manager 2.24 features to use in crow.
- Keywords: cdx, logics-manager, viewer, record-run, cost_usd, isolate-workspace, deck
- Use when: building any of items 073 to 076.
- Skip when: working on headless flags or compaction (req_017).

# Needs
- Open a project's Logics corpus from the crow deck in one key.
- See crow background runs from the Logics doc they work on.
- See what each background run cost.
- Run reviews in a workspace cdx isolates, not a hand-built checkout.

# Context
- Release review done 2026-11-03 while closing items 069 and 071.
- logics-manager `view` takes --focus, --fleet, --open, --port 0, --read, --lan. 2.24.2 adds `flow record-run` (Executions table per doc) and `--progress auto`.
- cdx 0.25.0 adds --budget and cost_usd per run (crow uses --budget since item_069). 0.24.5-0.24.6 add --isolate-workspace; 0.25.5-0.25.8 give isolated runs git and RTK access.
- Out of scope: cdx remote SSH workers (0.25.3, 0.25.9). Per-run strict MCP goes to a cdx-manager handoff and item_015.
- Owner decisions 2026-11-03: deck `v` opens the selected project's viewer, `V` the fleet view; plan all four slices.

# Acceptance criteria
- AC1: From the deck Projects page, one key opens the selected project's Logics viewer and one opens the fleet view.
- AC2: A dispatch linked to a Logics doc shows up in that doc's Executions table with its final state.
- AC3: `crow runs` and the deck Agents page show cost_usd for finished Claude runs.
- AC4: Review runs use cdx --isolate-workspace when it covers what the review needs, or the item records why not.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/crow-deck/main.go`, `tools/crow-deck/projects.go`
- `bin/crow-runs`, `bin/crow-review-run`, `bin/crow_lib.py`
- `runbooks/dispatch-background-runs.md`

# Backlog
- `item_073_deck_logics_viewer_keys`
- `item_074_dispatch_records_logics_runs`
- `item_075_runs_cost_per_run`
- `item_076_review_runs_isolated_workspace`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
