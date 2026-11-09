## req_043_nest_design_scope_and_jarvis_3d_hero_option - Nest design scope and Jarvis 3D hero option
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 13:53:36

# AI Context
- Summary: Narrow nest design to what the app needs (state language, 4 screens, one Home hero), and compare two hero directions: the droid workshop (req_040, paused) and a Jarvis-style 3D graph built in code.
- Keywords: nest, design scope, hero, Jarvis, 3D graph, holographic panels, state language, workshop
- Use when: deciding what nest design work remains before item_126 and item_137.
- Skip when: working on the core API (task_016) or tokens (req_038).

# Needs
- The droid art round (item_125) grew far beyond the need and is paused (owner, 2026-11-09). Cast picks stay in req_040.
- Nest needs designed:
  - a state language for working, needs you, failed, sleeping and queued that reads at a glance;
  - layouts for Home, Missions, Mission detail and Inbox on the Fold 6 cover and inner screens;
  - one Home hero, the "wow" piece.
- Tokens are done (req_038). The cast art only served the hero.
- Explore a second hero direction and compare it with the workshop before choosing.

# Context
- Jarvis direction, references named by the owner: Jarvis 3D space panels, cyberpunk game UI, Capri-style AI apparitions, the Obsidian 3D graph of notes.
- Data mapping (draft):
  - CR-0W = a central core, a light apparition that speaks.
  - Projects (registry) = clusters orbiting the core.
  - Missions = nodes in their project cluster.
  - State = motion, not colour: working pulses, needs you drifts toward the camera, failed glitches, sleeping dims, queued waits on the outer ring.
  - Runs = light pulses along the links.
  - Inbox and mission detail = glass panels that unfold in front of the camera on tap.
  - Vault notes = a faint background graph.
- Build: React Three Fiber with bloom. No character art and no image quota. Risks are bundle size and phone battery, with a 2D fallback where the panels become sheets.
- Workshop direction: req_040, cast picks and the 1c prompt in the vault `Nest Prompts v5.md`. It needs AI art for every droid and state.
- Method (memory design-renders-need-wow): map real crow data to visuals first, as an annotated 2D sketch. No renders before the owner's OK.

# Acceptance criteria
- AC1: An annotated 2D sketch of each hero (Jarvis graph, droid workshop) on the same real crow data: the same projects, missions and the 5 states, in a phone view and an inner-screen view.
- AC2: Both sketches go on one compare page (bin/crow-compare), and the owner picks a hero or a mix.
- AC3: The state language and the 4 screen layouts are sketched independently of the hero choice.
- AC4: item_126 (build plan) and item_137 (shell) are updated to the chosen hero; item_125 is closed or reshaped to match.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Scope
- In: design scope, the two hero sketches, the state language, the screen layouts.
- Out: AI renders, 3D code, and the modern theme (later, per req_040).

# Companion docs
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`

# References
- `logics/request/req_040_nest_deck_two_themes_3d_world.md`
- `logics/request/req_038_nest_design_rules.md`
- `logics/backlog/item_125_workshop_concept_art_round_2_original_cast_states_phone_and_night_views.md`
- `logics/backlog/item_137_nest_app_shell_on_core_fixtures.md`

# Backlog
- `item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches`
- `item_139_jarvis_graph_and_droid_workshop_hero_sketches_on_one_compare_page`
- `item_140_realign_nest_build_items_to_the_chosen_hero`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
