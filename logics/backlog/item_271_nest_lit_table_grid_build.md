## item_271_nest_lit_table_grid_build - Nest lit table grid build cost
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 70%
> Confidence: 60%
> Progress: 0%
> Complexity: high
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Cut the CPU cost of building the lit table cell grid each frame (planet ray casts, hex foot fills, board splats) now that item_270 draws it on WebGL.
- Keywords: nest, hero, lit table, grid build, planets, hex fill, perf
- Use when: the lit table frame is too slow after the WebGL draw, or the bench p10 ratio is the target.
- Skip when: the change is about drawing glyphs (item_270 owns the GL renderer) or the table look.

# Problem
- After item_270 the canvas draw is one drawImage, but the frame still costs about 1.6-1.9x less than 2D, not the 2x item_270 AC2 asked for. Bench p10 ratio GL: laptop 2.65-2.87, inner 2.07-2.28, cover 1.41. The rest is the JS that fills the cell grid each frame in `src/hero/scene.ts` (castPlanets, hexOnPlane feet, stalks, moons, trails).

# Scope
- In:
  - Profile the GL path first (`hero-table.mjs prof a` with `TABLEGL=on`) and record the top self times here.
  - Reuse work between frames: board layers that only move with `spin` (hex grid, foot fills) cached per spin step; planet cells skipped early outside each planet's screen box.
  - Move per-cell work into the GL pass where it is a pure function of cell position (for example planet shading in a shader), if the profile points there.
- Out:
  - Look changes without owner before/after shots (crow-compare).
  - The default plain table: stays pixel-identical.
  - The budget-miss fallback (item_220 AC5).

# Acceptance criteria
- AC1: `TABLEGL=on node scripts/hero-table.mjs bench a`: p10 frame/calib ratio at least 2x lower than the 2D baseline of item_270 on every size (2D: laptop 4.48-4.49, inner 3.40-3.79, cover 2.40-2.86).
- AC2: Lit frame pixel diff vs before under 0.5% at laptop, inner and cover; default hash unchanged (`hero-table.mjs hash`).
- AC3: `hero-table.mjs fallback` still ok on all sizes.
- AC4: tsc clean, vitest green, Nest rebuilt.

# AC Traceability
- req_062 (lit table perf) -> This item. Proof: AC1 bench numbers recorded in Notes.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_062
- Primary task(s): (none yet)
- Follows: item_270 (AC2 partial, owner 2026-11-20). Budget: item_159, item_220.

# Priority
- Priority: Medium
- Rationale: Idle medians already meet the item_159 budget on GL (laptop 5.0 ms, fold 6.5-7.7 ms). This is headroom for loaded machines.

# Notes
- Owner 2026-11-20: close item_270 with AC2 partial and open this item for the grid build.
- Bench needs `PLAYWRIGHT_BROWSERS_PATH` set to the local Playwright cache; compare ratios under load, not ms.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
