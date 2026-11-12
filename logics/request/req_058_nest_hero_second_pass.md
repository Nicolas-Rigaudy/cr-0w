## req_058_nest_hero_second_pass - Nest hero second pass: scale, sky planets, motion
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Second pass on the Nest hero after item_158: scale on laptop, readable moons, smooth head turn, projects as sky planets, more motion.
- Keywords: nest, hero, holotable, planets, moons, animation, scale
- Use when: planning or building hero work after item_158.
- Skip when: other Nest screens.

# Needs
- Owner feedback 2026-11-12 on the item_158 build (screenshot `Pictures/Screenshots/20261112_142808.png`):
  - Laptop head is small next to the holotable and the space it has. Phone size is good.
  - Moons are too small to read well.
  - Failed and green planets look poor. Projects should all be planets up in the sky, not one green thing on the floor.
  - The table does not spin with the planets. Should it spin a bit?
  - Head direction changes are brutal and jittery on fast left-right moves.
  - More animations: idle, reactions while following the pointer, planet motion.
- Plan first, build later, one item per fresh session.

# Context
- Hero code: `tools/nest/src/hero/` (scene.ts, Hero.tsx, declutter.ts, motion.ts).
- Shots and perf: runbook `runbooks/nest-screenshots.md` (hero section).
- Visual choices go through side-by-side boards (bin/crow-compare) and an owner pick.

# Acceptance criteria
- AC1: Each slice ends with an owner OK on shots or recordings at phone, inner and laptop.
- AC2: Declutter tests, idle and perf checks stay green (no overlap, 0 idle paints with reduced motion, under 20 ms per frame on laptop).

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/nest/src/hero/scene.ts`
- `tools/nest/src/hero/Hero.tsx`
- `runbooks/nest-screenshots.md`

# Backlog
- `item_194_hero_scale_laptop`
- `item_195_hero_moon_legibility`
- `item_196_hero_head_turn_smoothing`
- `item_197_hero_sky_planets`
- `item_198_hero_idle_follow_animations`
- `item_201_hero_gaze_follows_the_table_spin`

# Execution order
1. item_194 laptop scale (quick, no design pass).
2. item_196 head turn smoothing (motion code, no design pass).
3. item_195 moon legibility (after 194).
4. item_197 sky planets: concept sketch, options, owner pick, then build. Answers the table spin question.
5. item_198 idle and follow animations, table spin (after 196 and 197).
- Risk: 197 may move the moons, so 195 can fold into it if the sky layout lands first.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
