## item_296_nest_ops_capacity_lanes_and_afk_batch_tab - Nest Ops capacity lanes and AFK batch tab
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
- Summary: Add a CAPACITY lane group to the Ops timeline and an AFK BATCH tab to the Ops queue, keeping today's Ops.
- Keywords: nest, desk, ops, redesign, item_293
- Use when: Building or reviewing the Ops page.
- Skip when: Desk layout (item_295) or core fields (item_297).

# Problem
- Ops does not show who holds the two slots, the 5h window per account, or the AFK batch. The owner picked capacity as timeline lanes (C) and the AFK batch as a queue tab (X). Today's timeline graph, queue and day totals stay.
- Design source: vault CR-0W/Research/nest-desk-ops/desk-ops-model.md (Owner picks 2026-11-25) and the approved sketches in that folder; mockup tools/nest/scripts/mockups/desk-ops.html.

# Scope
- In:
  - CAPACITY lane group in the timeline on the same axis: slot 1 and slot 2 with their holders, outside slots (butler run), 5h window per account (used to now, dashed until reset, % left, burn rate).
  - Queue pane tabs QUEUE / AFK BATCH.
  - AFK tab: job list with drag reorder, EDIT and ✕, GO · START BATCH, START AT, + ADD JOB.
  - Phone layout of both.
- Out:
  - Removing or replacing the existing lanes, queue or totals.
  - Core batch state and quota history (item_297).

# Acceptance criteria
- AC1: The existing timeline lanes, queue and day totals render unchanged.
- AC2: The CAPACITY group renders slot holders and one 5h lane per account from /api data.
- AC3: The AFK BATCH tab lists the batch jobs and its actions call core routes; without item_297 they show disabled.
- AC4: Laptop and phone shots (chat folded) match sketches C and X and are shown to the owner before merge.
- AC5: Nest tests and lint pass.

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

