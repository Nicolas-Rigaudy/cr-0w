## item_157_nest_hero_ascii_graph_layer_around_the_core - Nest hero ASCII graph layer around the core
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
- Summary: Add the Jarvis graph to the hero shader: orbit rings, project nodes, core links, the 5 state looks and depth labels, drawn as ASCII glyph clusters in the same cell pass as the core, with HTML labels on top.
- Keywords: nest, hero, ASCII graph, orbit ring, project nodes, links, WebGL2
- Use when: building or changing the graph around the core in the GPU hero.
- Skip when: live data and motion (item_158).

# Problem
- item_156 renders the core alone. The data map in item_126 needs projects, missions and the queue around it in the option B ASCII look.

# Scope
- In:
  - Scene input as one typed object: projects (angle, depth, worst state), missions per project, queue entries, usage arcs, background graph points.
  - Inner ring (`registry.projects`), outer dashed ring (`queue[]`), core links, faint background graph, emitter base arcs per account.
  - Node looks per state: working (solid cyan with pulse ring), needs you (amber diamond), failed (broken ring, link cut), sleeping (dotted, 45% dim), queued (hollow on the outer ring).
  - HTML labels positioned from the shader camera: `project.name` plus worst state, size and brightness from depth.
  - Hit testing per node and for the core, returning the element id for item_158.
- Out:
  - Data wiring, state motions and actions (item_158).

# Acceptance criteria
- AC1: A fixture scene with all 5 states renders like `item126/look-ascii.png` (hero pass 7); the owner approves side by side.
- AC2: Labels stay pinned to their nodes on resize and on every DPR from 1 to 2.63.
- AC3: Hit testing returns the right node or core id for a tap on laptop and Fold inner.
- AC4: Frame time with the graph stays under the item_126 budget on the laptop iGPU.

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
- Depends on: `item_156_nest_hero_core_head_as_a_webgl2_ascii_shader`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Completes the hero picture before it goes live.

# Notes
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`). Not `item139/jarvis-1-full.png` (owner rejected it as a reference, 2026-11-10). Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
