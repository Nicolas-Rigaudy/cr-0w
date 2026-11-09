## item_140_realign_nest_build_items_to_the_chosen_hero - Realign nest build items to the chosen hero
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Once the owner picks a hero in item_139, update item_126 (build plan) and item_137 (shell) to it, and close or reshape item_125.
- Keywords: realign, nest, build, items, chosen, hero
- Use when: after the hero pick in item_139 is recorded.
- Skip when: before item_139 has a recorded pick.

# Problem
- item_126 and item_137 assume the workshop concept art from item_125. If the owner picks the Jarvis graph or a mix, both are wrong, and item_125 may no longer be needed.

# Scope
- In:
  - Rewrite item_126 as the build plan for the chosen hero (for Jarvis: React Three Fiber with bloom, bundle and battery budget, 2D fallback).
  - Update item_137 dependencies, AC4 and the visual layer line to the chosen hero and to the item_138 layouts.
  - Close item_125 if the workshop is dropped, or shrink it to the states the chosen mix still needs.
  - Mark req_043 done when all ACs hold.
- Out:
  - Any design or code work.

# Acceptance criteria
- AC1: item_126 and item_137 name the chosen hero and depend on item_138 and item_139 instead of the item_125 art where that no longer applies.
- AC2: item_125 is closed or reshaped, with the reason recorded.
- AC3: `logics-manager lint` passes on the touched docs.

# AC Traceability
- request-AC4 -> AC1, AC2. Proof: the updated item_125, item_126 and item_137 docs.

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`
- Depends on: `item_139_jarvis_graph_and_droid_workshop_hero_sketches_on_one_compare_page`, `item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Small doc pass; it keeps item_137 from building against a stale plan.

# Notes
- Hybrid rationale: Derived from request `req_043_nest_design_scope_and_jarvis_3d_hero_option` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`.
- Generated locally by logics-manager.

# Result (2026-11-10)
- item_126 rewritten as the build plan for the Jarvis graph hero with the droid core (R3F + bloom, 2D SVG fallback, budgets), depending on item_138 and item_139.
- item_137 depends on item_119, item_138 and item_139; AC4 and the hero slot point at the item_138 boards and the item_126 fallback; the item_138 build notes are in scope. It no longer waits for item_126.
- item_125 withdrawn, superseded by item_139: the workshop is no longer the hero.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
