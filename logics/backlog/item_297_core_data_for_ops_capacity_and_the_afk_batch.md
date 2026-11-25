## item_297_core_data_for_ops_capacity_and_the_afk_batch - Core data for Ops capacity and the AFK batch
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Close the core data gaps the Desk and Ops rebuild needs: slot holders, AFK batch state, quota history, failure cause, run state.
- Keywords: nest, desk, ops, redesign, item_293
- Use when: crow-core work feeding item_295 and item_296.
- Skip when: Nest UI work.

# Problem
- The live data map of 2026-11-25 shows gaps: the butler run is not counted in dispatch.slots, there is no AFK batch state or order, no quota history for the 5h lanes, failed missions have why=done, runs[].state is null, stopped and failed are folded into why.
- Design source: vault CR-0W/Research/nest-desk-ops/desk-ops-model.md (Owner picks 2026-11-25) and the approved sketches in that folder; mockup tools/nest/scripts/mockups/desk-ops.html.

# Scope
- In:
  - dispatch.slots lists holders per slot and an outside-slots list (butler run).
  - AFK batch state in core: ordered jobs, edit, remove, reorder, start now, start at; routes for each.
  - Quota history per account over the 5h window (samples for the lane).
  - Failure cause recorded on failed missions.
  - runs[].state filled; done/stopped/failed kept apart.
- Out:
  - Nest UI (item_295, item_296).
  - A pending-dispatch queue outside the AFK batch.

# Acceptance criteria
- AC1: /api exposes slot holders and outside-slots runs; busy counts match the agents shown.
- AC2: AFK batch routes create, edit, reorder, remove and start a batch, with tests.
- AC3: /api exposes quota samples per account for the current 5h window.
- AC4: A failed mission carries its cause; runs[].state is set; stopped and failed are distinct states.
- AC5: Go tests and lint pass.

# AC Traceability
- req_070-AC -> This backlog slice. Proof: owner picks recorded 2026-11-25.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_070
- Blocks: item_296 (batch actions, capacity lanes)
- Related: item_293 (sketches), item_292
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Default until groomed.

# Notes
- Build order: item_297 core first, then item_296 and item_295 (independent).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

