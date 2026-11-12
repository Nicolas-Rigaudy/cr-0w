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
- Summary: Size the head and scene from the space the hero has, so the laptop hero fills its column like the phone does.
- Keywords: nest, hero, scale, laptop, head, holotable
- Use when: building or reviewing the Nest hero after item_158.
- Skip when: work outside the Nest hero.

# Problem
- On laptop the head is small next to the holotable and the empty column (owner screenshot 2026-11-12 14:28). On phone it is a good size.
- Cause: `drawScene` sets the cell size to `HS * clamp(w/560, 0.7, 1)`. The factor stops at 1, so a 650 px wide hero draws the head at phone size while the table ellipse stretches to the width.
- The pass also made the phone head a bit smaller than before, which helped the planets fit there.

# Scope
- In:
  - Cell size from the hero width and height, with a higher cap on laptop.
  - Head, table and planets keep one cell size (item_158 rule).
  - Keep the phone at today's size.
- Out:
  - Planet layout and looks (item_197).

# Acceptance criteria
- AC1: On laptop 1440x900 the head fills about the same share of the hero as on phone 368x812.
- AC2: Phone and inner shots unchanged or approved by the owner.
- AC3: Draw time per frame stays under 20 ms on laptop (`scripts/hero-perf.mjs`).
- AC4: Owner OK on before/after shots side by side.

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
