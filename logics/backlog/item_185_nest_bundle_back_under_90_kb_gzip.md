## item_185_nest_bundle_back_under_90_kb_gzip - Nest bundle back under 90 kB gzip
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 80%
> Progress: 0%
> Complexity: Small
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Nest JS is 95.9 kB gzip after item_171 (94.6 kB before it); the README budget and item_171 AC4 say under 90 kB.
- Keywords: nest, bundle, gzip, budget, vite, code split
- Use when: bringing the Nest bundle back inside its budget.
- Skip when: changing Nest screens or data.

# Problem
- The Nest JS bundle passed its 90 kB gzip budget before item_171 and item_171 added about 1.3 kB. The owner chose to ship item_171 and fix the budget here.

# Scope
- In:
  - measure what weighs most (`npx vite build` then a size breakdown, e.g. rollup-plugin-visualizer or source-map-explorer, dev only);
  - likely candidates: the holotable head data in the hero, screens not shown on first load (Desk, Missions detail) behind lazy imports;
  - update the README "Last build" line.
- Out:
  - visual or behaviour changes.

# Acceptance criteria
- AC1: `npx vite build` reports JS under 90 kB gzip for the first load.
- AC2: `npm test` passes and Home, Desk, Missions and Inbox look the same in `scripts/shots.mjs` shots.

# AC Traceability
- item_171 AC4 -> AC1. Proof: build output.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): (none)
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: Low
- Rationale: phone loads over the tailnet; owner chose ship first, trim later (2026-11-10).

# Notes
- Filed from item_171 (commit 06545b2).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
