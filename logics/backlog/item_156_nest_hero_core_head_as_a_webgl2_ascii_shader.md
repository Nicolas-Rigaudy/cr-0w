## item_156_nest_hero_core_head_as_a_webgl2_ascii_shader - Nest hero core head as a WebGL2 ASCII shader
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Port the spike head (`item126/spike/holo-head.html`) into a lazy-loaded nest module: raw WebGL2, cell pass, glyph atlas pass, two-level bloom, moods idle, run, ask and fail, glitch on fail only. No holo extras.
- Keywords: nest, hero, WebGL2, ASCII shader, glyph atlas, bloom, CR-0W head, moods
- Use when: building or changing the GPU core head of the Home hero.
- Skip when: the static fallback (item_155) or the graph around the core (item_157).

# Problem
- The owner approved the plain shader port of the pass 7 head (spike 2026-11-10). It lives in one HTML file with URL params and needs to become an app module with a typed API.

# Scope
- In:
  - A `heroGl` module in `tools/nest`, lazy-loaded, raw WebGL2, no three.js. Shaders as TS string modules.
  - Cell pass (head sphere with face primitives, arms, fin, octa and held sphere by SDF march) and glyph pass on a runtime-built Fira Code atlas, bloom at half and quarter resolution, composite.
  - The spike holo layer stays out: no emitter ring, beam, floor disc or scanlines.
  - API: `mount(canvas)`, `setMood(mood, since)`, `setSize(w, h, dpr)`, `pause()`, `resume()`, `dispose()`. Context loss handled by rebuilding the resources.
  - Motion from the spike: idle drift, hashed blink every 4-7 s, 160 ms easeOutBack snap on mood change, 180 ms glitch on fail. 12 fps settled, full rate only during a snap or glitch.
  - Mounted over the item_155 static hero; the static head shows until the first GPU frame.
- Out:
  - The graph layer (item_157).
  - Fallback detection and measured budgets (item_159).

# Acceptance criteria
- AC1: Renders of idle, ask and fail match `item126/spike/img/*-port.png` side by side; the owner approves.
- AC2: Setting a mood snaps in 160 ms; fail plays one 180 ms glitch, then holds still.
- AC3: The rAF loop runs at 12 fps when settled and stops on `pause()`.
- AC4: WebGL context loss and restore leave the head rendering with no reload.
- AC5: The module is its own lazy chunk; the main bundle does not grow by more than an import stub.

# AC Traceability
- item_126 AC4 -> AC1 to AC5. Proof: this item's notes and the owner's approval.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (render path and fallback rule settled in item_126)

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_043_nest_design_scope_and_jarvis_3d_hero_option`
- Plan: `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`
- Depends on: `item_155_nest_hero_static_fallback_in_the_home_slot`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: The core is the hero's anchor; the graph layer builds on its passes.

# Notes
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`) and `item139/jarvis-1-full.png`. Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
