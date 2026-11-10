## item_159_nest_hero_budgets_and_fallback_detection_measured_on_laptop_and_fold - Nest hero budgets and fallback detection measured on laptop and Fold
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 13:21:37

# AI Context
- Summary: Pick the static or GPU hero by the item_126 fallback rule and measure the budgets on the laptop and the Fold 6.
- Keywords: nest, hero, budgets, fallback rule, battery, Save-Data, software renderer, Fold
- Use when: tuning when the hero falls back, or checking its cost.
- Skip when: the look of either hero.

# Problem
- The spike showed SwiftShader at 814 ms per frame. The hero must fall back on its own when the GPU path is too costly, and the budgets need real numbers from the phone.

# Scope
- In:
  - A `heroMode()` check: no WebGL2, `prefers-reduced-motion`, Save-Data, battery under 20% and not charging, Fold cover, software renderer (`WEBGL_debug_renderer_info` matches SwiftShader or llvmpipe), manual toggle in settings. Re-evaluated on change events.
  - Hidden tab stops the loop; visible resumes it.
  - Bench mode (`?bench=N`) logging ms per frame, following `runbooks/headless-webgl-render.md`.
  - Measures: hero chunk size gz, frame time on the laptop iGPU at 1440x900 and on Fold inner, frames per second idle and during a live run.
- Out:
  - Changing the budgets themselves: an owner decision recorded back in item_126.

# Acceptance criteria
- AC1: Each fallback trigger, forced one at a time, shows the static hero.
- AC2: A hidden tab renders 0 frames.
- AC3: Measures written in this item's notes and checked against item_126: chunk under 30 kB gz, under 10 ms on the laptop, under 16 ms on Fold inner, 0 fps idle.
- AC4: A miss on any budget is reported to the owner with the numbers, not silently accepted.

# AC Traceability
- item_126 AC4 -> AC1 to AC4. Proof: this item's notes and the owner's approval.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (render path and fallback rule settled in item_126)

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_043_nest_design_scope_and_jarvis_3d_hero_option`
- Plan: `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`
- Depends on: `item_158_nest_hero_live_data_state_motions_and_actions`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Last step before the hero is called done.

# Notes
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`). Not `item139/jarvis-1-full.png` (owner rejected it as a reference, 2026-11-10). Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
