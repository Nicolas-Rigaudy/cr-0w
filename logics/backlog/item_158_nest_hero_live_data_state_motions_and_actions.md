## item_158_nest_hero_live_data_state_motions_and_actions - Nest hero live data, state motions and actions
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
- Summary: Drive the hero from the nest data layer: scene and core mood from missions, inbox, queue and usage, the item_138 state motions, and the tap actions.
- Keywords: nest, hero, data layer, state motion, actions, inbox, missions, desk
- Use when: connecting the hero to crow data or changing what a hero tap does.
- Skip when: rendering changes (item_156, item_157) or the fallback rule (item_159).

# Problem
- The shader and the static hero take a scene object. Nothing yet builds it from crow data or reacts to changes.

# Scope
- In:
  - A `heroScene` selector over the item_137 data layer, following the item_126 data map row by row. Fixtures first, the API adapter when slice 3 lands.
  - Core mood: fail when the last run failed, ask when `inbox.open_count` > 0, run while a run is live, otherwise idle.
  - State motions per item_138: working 2.4 s pulse and orbs per `run.step`, needs you lift 4 px and one tick, failed one 180 ms glitch, sleeping and queued still. Reduced motion shows final frames only.
  - Orbit drift only while a run is live; otherwise the loop sleeps.
  - CR-0W line typing out the latest `coordinator.message`; asks chip with `inbox.open_count`.
  - Actions: project node opens Missions filtered to it, core opens the Desk chat, needs-you node opens its Inbox row, failed node shows `run.error`. Keyboard path for each through a focusable list.
- Out:
  - Changing state from the hero (approve, retry): stays in Missions and Inbox.

# Acceptance criteria
- AC1: Switching fixtures between the 5 cases changes the scene and the core mood with the item_138 motions and nothing else moves.
- AC2: With no state change and no live run, the hero renders 0 frames per second.
- AC3: Each action opens the right route with the right filter, by tap and by keyboard.
- AC4: The static hero and the shader take the same scene object.

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
- Depends on: `item_157_nest_hero_ascii_graph_layer_around_the_core`, `item_137_nest_app_shell_on_core_fixtures`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Turns the hero from a picture into the Home overview.

# Notes
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`). Not `item139/jarvis-1-full.png` (owner rejected it as a reference, 2026-11-10). Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
