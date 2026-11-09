## item_155_nest_hero_static_fallback_in_the_home_slot - Nest hero static fallback in the Home slot
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
- Summary: Fill the item_137 Home hero slot with the static hero: an SVG Jarvis graph around a pre-rendered ASCII head per mood, 0 fps. It ships first and stays as the fallback for the shader.
- Keywords: nest, hero, static fallback, SVG graph, ASCII head, Home slot, item_137
- Use when: building the first hero in the nest shell, or changing the static fallback.
- Skip when: the WebGL2 shader (item_156, item_157).

# Problem
- item_137 reserves a Home hero slot and falls back to a static PNG until item_126 lands.
- The fallback rule from item_126 needs a static hero that works with no GPU, on the Fold cover and with reduced motion.

# Scope
- In:
  - A `Hero` component in `tools/nest` with a static renderer: SVG orbit rings, project nodes, links and labels in the option B ASCII look (`look-ascii.png`), the ASCII head as preformatted text in the centre.
  - Head text per mood (idle, run, ask, fail) generated once from `spike/holo-head.html` (render to text, not PNG) and checked into the app as data.
  - Node look per state (working, needs you, failed, sleeping, queued) as static marks, colours from the item_126 palette.
  - Fixture data from the item_137 data layer; no live API.
  - Placement per item_138: Home centre slot, 50% dim behind Missions and Inbox, both themes, laptop, Fold inner and Fold cover.
- Out:
  - The WebGL2 shader (item_156, item_157).
  - Motion and live data (item_158).

# Acceptance criteria
- AC1: The Home slot renders the static hero from the item_119 fixtures on laptop, Fold inner and Fold cover, both themes.
- AC2: Each mood shows its own ASCII head text, matching the spike port.
- AC3: Each of the 5 mission states has a distinct static node mark from the item_126 data map.
- AC4: Zero animation frames: no rAF loop and no CSS animation in the static hero.
- AC5: The owner approves the static hero side by side with `item139/jarvis-1-full.png` and `look-ascii.png`.

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
- Depends on: `item_137_nest_app_shell_on_core_fixtures` (the shell and its hero slot), `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: First hero build item. Unblocks a real hero in the shell without any GPU work.

# Notes
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`) and `item139/jarvis-1-full.png`. Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
