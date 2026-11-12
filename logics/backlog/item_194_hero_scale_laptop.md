## item_194_hero_scale_laptop - Hero scales to the space it has on laptop
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: One scale rule for phone cover, phone inner and laptop: head and scene take the same share of the hero everywhere, tuned on today's phone look.
- Keywords: nest, hero, scale, laptop, head, holotable
- Use when: building or reviewing the Nest hero after item_158.
- Skip when: work outside the Nest hero.

# Problem
- On laptop the head is small next to the holotable and the empty column (owner screenshot 2026-11-12 14:28). On phone it is a good size.
- Cause: `drawScene` sets the cell size to `HS * clamp(w/560, 0.7, 1)`. The factor stops at 1, so a 650 px wide hero draws the head at phone size while the table ellipse stretches to the width.
- The pass also made the phone head a bit smaller than before, which helped the planets fit there.

# Scope
- In:
  - One rule for all sizes: cell size from the hero width and height (the smaller fit), no per-device cap.
  - Calibrate on the phone cover look the owner likes today; inner and laptop follow the same proportions.
  - Head, table and planets keep one cell size (item_158 rule).
- Out:
  - Planet layout and looks (item_197).

# Acceptance criteria
- AC1: Head height over hero height is the same within 10% on cover 368x812, inner 707x691 and laptop 1440x900.
- AC2: The phone cover keeps today's look (owner check); inner and laptop match it in proportion.
- AC3: Draw time per frame stays under 20 ms on laptop (`scripts/hero-perf.mjs`).
- AC4: Owner OK on one sheet with the three sizes side by side, before and after.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: owner OK on shots.
- request-AC2 -> This backlog slice. Proof: tests and perf scripts.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_058_nest_hero_second_pass`
- Primary task(s): (none yet)
- Follows: `item_158_nest_hero_live_data_state_motions_and_actions`

# Priority
- Priority: High
- Rationale: owner feedback 2026-11-12 on the item_158 build.

# Notes
- Quick slice, no design pass needed: before/after sheet via bin/crow-compare.
- Do it before item_195, since bigger cells already help moon reading on laptop.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
