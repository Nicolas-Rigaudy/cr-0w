## item_126_plan_the_workshop_nest_build_from_the_approved_concept_art - Plan the Nest Home hero build (Jarvis graph with the droid core)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 13:23:18

# AI Context
- Summary: Turn the item_139 hero pick into a build plan: the Jarvis graph (projects orbiting the CR-0W droid core, state as motion) in React Three Fiber with bloom, a 2D SVG fallback, budgets, and how it plugs into the item_137 shell and the item_138 layouts. Planning only.
- Keywords: nest, hero, Jarvis graph, droid core, React Three Fiber, bloom, SVG fallback, budget, crow core
- Use when: planning how the Home hero gets built after item_138 and item_139.
- Skip when: choosing the look (done in item_138 and item_139) or building the shell (item_137).

# Problem
- item_139 picked the Jarvis graph as the Home hero, with the droid from passes 5-7 as its core (visor band hull, slit-eye LED face, pincer arms, dorsal fin, no ring). item_138 fixed where the hero sits: free centre slot on Home, dimmed at 50% behind Missions and Inbox, on laptop and Fold inner.
- Nothing says yet how the graph and the core get built, what drives them, or what they cost on the phone.

# Scope
- In:
  - Data map: every visual element mapped to a crow source. Project nodes (registry), mission edges and their 5 states (crow core store), needs you (inbox rows), core moods idle, run, ask and fail (store and live runs), quota (usage).
  - State as motion: one motion per state from the item_138 state sheet, reduced-motion variant included.
  - Render approach: React Three Fiber with bloom for the graph; the droid core as glTF, or drawn from the `src139/pass7/hero7.html` renderer. The 2D SVG fallback and the rule that picks it (no WebGL, reduced motion, battery saver, Fold cover).
  - Placement per item_138: Home centre slot, dimmed hero behind Missions and Inbox, both themes (modern and cyber).
  - Actions: a project node opens Missions filtered to it, the core opens the Desk chat, a needs-you node opens its Inbox row.
  - Budgets: three.js chunk lazy-loaded with a size cap, frame rate cap, pause when hidden, battery rule.
  - Build backlog items with ACs, dependencies and order, the first one plugging into the item_137 hero slot. Then stop.
- Out:
  - The droid workshop room and its cast (item_125, withdrawn).
  - Building anything in this item.

# Acceptance criteria
- AC1: Data map table, one row per visual element and its crow source.
- AC2: Render approach and fallback rule chosen by the owner from at most 3 options, shown side by side.
- AC3: Budgets written down: bundle size, frame rate, battery and hidden-tab behaviour.
- AC4: Build backlog items written and linted, in execution order, the first one filling the item_137 hero slot.

# AC Traceability
- request-AC1 -> AC1, AC2, AC3, AC4. Proof: this item's notes and the new backlog items.

# Decision framing
- Product framing: Not needed
- Architecture framing: Consider

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_043_nest_design_scope_and_jarvis_3d_hero_option`
- Depends on: `item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches`, `item_139_jarvis_graph_and_droid_workshop_hero_sketches_on_one_compare_page`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: R&D. The shell (item_137) can start with a static hero; this plan fills the slot.

# Notes
- Hero sources: vault `CR-0W/Research/nest-design/item139/` (compare page, `jarvis-1-full`, `jarvis-2-in-home`, `pass7/hero-{idle,run,ask,fail}.png`, `pass7/heads-final.png`).
- Layout sources: canvas artifact MyPAqSwJFddMnSR4grm2g6 and vault `CR-0W/Research/nest-design/item138/fold-src/`.
- Nest stack itself: req_031.
- 2026-11-10 realigned by item_140: the workshop is dropped as hero (item_139 pick), item_125 withdrawn. This item now plans the Jarvis graph with the droid core.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
