## item_274_ops_on_the_phone_now_card_first_totals_table_fits_368 - Ops on the phone: Now card first, totals table fits 368
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: On the phone, Ops opens with a Now card (live runs with OPEN, the free slot and what starts next), then queue, timeline and a 4-column totals table that fits 368.
- Keywords: ops, phone, now, card, first, totals, table, fits, 368
- Use when: Building the phone layout of Ops.
- Skip when: Laptop layout (unchanged).

# Problem
- At 368 the Ops page stacks timeline, queue, totals, and the totals table runs past its card (failed and cost columns cut). What runs now is hard to read on the phone.
- Design approved 24/11: vault `CR-0W/Research/nest-design/item_190/phone/compare.html`, option B; decisions in `item_190/design-notes.md`.

# Scope
- In:
  - Phone only (below the laptop breakpoint): order summary line, Now card, queue, timeline, totals.
  - Now card: each live lane with title, run id, age, cost and OPEN to its chat; hub session live; a free slot row naming the next queued item.
  - Totals on the phone: project, active, ! (failed or stuck), cost; tapping a row shows lanes, wait, stuck.
  - Cyber theme check on the new pieces.
- Out:
  - Laptop layout. Week tab extras.

# Acceptance criteria
- AC1: At 368x812 no Ops card overflows horizontally (scrollWidth equals clientWidth), day and week, both themes.
- AC2: Now card lists every live lane with a working OPEN, and the free slot names the queue head.
- AC3: Layout matches phone/compare.html option B.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: owner-approved Ops design, built and checked on live data.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)
- Request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`, `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Primary task(s): `task_023_orchestrate_nest_laptop_parity`
- Depends on: `item_190_laptop_missions_from_the_approved_missions_board`

# Priority
- Priority: Medium
- Rationale: Finishes the Ops page (item_190). Build order: item_272, item_273, item_274; each in a fresh session.

# Notes
- Read the vault design notes and the compare page before building; build from the chosen option only.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-24T17:16:23.495444+00:00 | afk | afk-ops-phone-2 | Ops phone checks | running |  | AC1 and AC2 pass at 368 (scripts/ops-phone-shots.mjs, day and week, both themes). AC3 not compared to board B (vault not readable in this run). Not pushed: git fetch and push were refused. Branch afk/ops-phone. |
| 2026-11-24T17:14:02.549530+00:00 | afk | afk-ops-phone | Ops phone Now card and totals | cancelled |  | Built on branch afk/ops-phone (5bb1e9d), not on main. Unverified: node, npm ci and the node_modules symlink were refused, so no typecheck or screenshots; AC1-AC3 unchecked. |

