## task_007_0_4_projects_and_usage_tabs - 0.4 Projects and Usage tabs
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 80%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Projects tab (cards, roadmap progress from the Logics roadmap docs, timeline) and Usage tab (accounts, forecast, daily spend, savings, cost today) of the crow deck, read only.
- Keywords: projects, usage, tabs
- Use when: continuing the deck after the Projects and Usage tabs, or wiring `e` (edit roadmap) and the policy keys.
- Skip when: the change is outside the deck.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_019_0_4_projects_and_usage_tabs`

# Acceptance criteria
- AC1: Roadmap data comes from one file per project, edited with `e`.
- AC2: Cost per session matches `cdx stats` within rounding.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_007_0_4_projects_and_usage_tabs.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_007_0_4_projects_and_usage_tabs.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: Usage sections, single-pane walk, keys (d, p, j, r), loading and empty states, `usageViewOf` on a payload; Projects cards, actions on the selected card only, keys (j, c, z), mode filter, live cards without dates, registry road decoding, timeline window; fit tests for every theme and size on both pages.
- `python3 -m pytest tests/test_crow_registry.py`: 5 passed (roadmap docs invented: milestones counted, current milestone and exit signal, absolute path, no roadmap).
- Live check: `bin/crow-registry json` returned a `road` block for the four registry projects that name a roadmap doc. Output not recorded here.
- Renders of both tabs in 4 themes (default, selected card, zoom, work mode, portrait, single panes, states) under local/state/afk/shots/deck-tabs/.

# Report
- Shipped the Usage tab (accounts and policy gauges with resets, forecast, daily spend 7 or 14 days, savings, cost today in three columns), fed by crow-usage and crow-accounts, and the Projects tab (cards, ROAD and MILE from the roadmap docs through crow-registry, timeline, repos without a roadmap on one line). Home's Projects box now reads the same cards.
- AC1 is partly met: roadmap data is one Logics roadmap doc per project, read by `crow-registry`. The `e` key says it is not wired and edits nothing.
- AC2 is not verified: costs are API-equivalent dollars from crow-usage and the tab says so. Compare a session with `cdx stats` by hand.
- Open: the roadmap docs carry milestones but no target dates, so live cards show "no date" and the timeline shows only a note. Dated cards and the timeline exist in mock data only. A date format is a product decision.
- Open: the Forecast box has no "afk tonight" line (no source), the policy keys p l [ ] and the card keys c r s e o only flash.

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
