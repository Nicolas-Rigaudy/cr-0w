## item_295_nest_desk_rebuild_on_deck_panes - Nest Desk rebuild on deck panes
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
- Summary: Rebuild the Nest Desk page as deck panes (layout C) from the owner picks of the item_293 pass.
- Keywords: nest, desk, ops, redesign, item_293
- Use when: Building or reviewing the Desk page layout.
- Skip when: Ops page or core data changes (item_296, item_297).

# Problem
- Desk shows work as cards and count tiles that hide what needs the owner. The owner picked layout C: deck panes with counts in the titles.
- Design source: vault CR-0W/Research/nest-desk-ops/desk-ops-model.md (Owner picks 2026-11-25) and the approved sketches in that folder; mockup tools/nest/scripts/mockups/desk-ops.html.

# Scope
- In:
  - Panes [1] Needs you, [2] Moving, [3] Queued, [4] Done, plus a detail pane.
  - Queued folds to one line of titles and expands on key 3; the full queue stays on Ops.
  - Detail tabs: Overview, Transcript, Timeline, Files.
  - Launcher with three modes NOW / BACKGROUND / AFK BATCH; AFK BATCH adds the job to the Ops batch.
  - No count tiles; each pane title carries its count.
  - Butler chat folded to a button; shots taken with it folded.
  - Phone layout: panes stacked in the same order.
- Out:
  - Ops page changes (item_296).
  - New core fields (item_297); missing data shows as an empty state.

# Acceptance criteria
- AC1: Desk renders the four panes and the detail pane from the live /api data, with counts in pane titles.
- AC2: Key 1-4 focuses a pane; key 3 expands the folded Queued pane.
- AC3: The detail pane shows the four tabs for the selected row.
- AC4: The launcher offers NOW / BACKGROUND / AFK BATCH and each mode reaches its existing endpoint.
- AC5: Laptop and phone shots (chat folded) match the C sketches and are shown to the owner before merge.
- AC6: Nest tests and lint pass.

# AC Traceability
- req_070-AC -> This backlog slice. Proof: owner picks recorded 2026-11-25.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_070
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

