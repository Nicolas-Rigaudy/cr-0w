## task_033_nest_desk_rebuild_on_deck_panes - Nest Desk rebuild on deck panes
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-25 20:55:59

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: nest, desk, rebuild, deck, panes
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_295_nest_desk_rebuild_on_deck_panes`

# Acceptance criteria
- AC1: Desk renders the four panes and the detail pane from the live /api data, with counts in pane titles.
- AC2: Key 1-4 focuses a pane; key 3 expands the folded Queued pane.
- AC3: The detail pane shows the four tabs for the selected row.
- AC4: The launcher offers NOW / BACKGROUND / AFK BATCH and each mode reaches its existing endpoint.
- AC5: Laptop and phone shots (chat folded) match the C sketches and are shown to the owner before merge.
- AC6: Nest tests and lint pass.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_033_nest_desk_rebuild_on_deck_panes.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_033_nest_desk_rebuild_on_deck_panes.md` after implementation.

# Validation
- (no validation recorded yet)
- command: `cd tools/nest && npm test && npm run typecheck && npm run build` | result: passed | date: 2026-11-25
- Finish workflow executed on 2026-11-25.
- Linked backlog/request close verification passed.

# Report
- Not started.
- Finished on 2026-11-25.
- Linked backlog item(s): `item_295_nest_desk_rebuild_on_deck_panes`
- Related request(s): `req_070_desk_and_ops_handle_missions_runs_and_queues`

# Links
- Request: (none yet)
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# AC Traceability
- request-AC1 -> This task. Proof: builds the picked option C (desk-C.png, desk-phone-C.png); implemented in 2489716 and 772e521, shots in the vault nest-desk-ops/built folder.
- request-AC2 -> This task. Proof: this task is the Desk build item of the split (item_295); Ops and core stay in item_296 and item_297.
- request-AC3 -> This task. Proof: each Desk count (pane titles, summary, phone tabs) sits on its pane and the pane lists the rows; covered by Desk.test.tsx layout C tests, validated with cd tools/nest && npm test.
