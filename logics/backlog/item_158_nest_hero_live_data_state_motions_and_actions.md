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
> Indicators reviewed: 2026-11-10 14:53:51

# AI Context
- Summary: Drive the holotable hero from the nest data layer and make it explorable: spin the planet layer, hover for moons, click to focus then open the project.
- Keywords: nest, hero, holotable, planets, moons, spin, data layer, state motion, actions
- Use when: connecting the hero to crow data or changing how the holotable reacts to input.
- Skip when: rendering of the static holotable (item_162) or the fallback rule (item_159).

# Problem
- The holotable (item_162) is a static picture built from fixtures. Nothing builds its scene from crow data, reacts to changes, or lets the owner explore it.

# Scope
- In:
  - Data: a `heroScene` selector over the item_137 data layer, following the item_126 data map. Fixtures first, the API adapter when slice 3 lands. Planets = projects, moons = open missions, stalks = wait time, rim = 24h agenda (item_162 mapping).
  - Core mood: fail when the last run failed, ask when `inbox.open_count` > 0, run while a run is live, otherwise idle. CR-0W line types out the latest `coordinator.message`; asks chip shows `inbox.open_count`.
  - State motions per item_138: working 2.4 s pulse, needs you lift 4 px and one tick, failed one 180 ms glitch, sleeping and queued still.
  - Spin: horizontal drag on the table turns the planet layer around the table's vertical axis only. Inertia with friction, then a soft detent eases the nearest planet to the front. Faint table grid marks turn with the planets so the motion reads.
  - Fixed parts: the 24h rim (time, not space), the lens and the bot never move with the spin.
  - Bot reacts: eyes and head follow the drag direction with a slight lean; look at the hovered planet; on release, turn back to the longest-wait planet.
  - Readability: labels are billboarded (always upright, facing the viewer). Back planets shrink and dim, their labels hide except `YOUR MOVE` tags. Labels hide when the head or lens covers them.
  - Hover a planet: it lifts, its moons appear (open missions coloured by state, max 6 then `+n`), the other planets dim. Hover a moon: mission title. One lifted planet at a time (keeps the item_162 glow budget of 4).
  - Click a planet: first click rotates it to the front and floats a detail card (missions, what waits on you, last run). Second click or Enter opens Missions filtered to the project. Esc returns to the table. Click the core: opens the Desk chat. Failed planet card shows `run.error`.
  - Touch (Fold, phone): tap = hover (lift, moons, card), tap the card = open. Horizontal swipe spins, vertical swipe scrolls (`touch-action: pan-y`).
  - Keyboard: Left/Right step to the next planet (rotates and snaps), Enter opens, Esc back. Focusable list mirrors the planets for screen readers.
  - Reduced motion: step rotation, no inertia, final frames of state motions only.
  - Build: 2D canvas; planets keep a polar angle, the spin adds one offset. Labels, moons tooltips and cards are DOM overlays positioned from the projected coordinates (crisp text, a11y). Hit testing on projected positions.
  - Frame budget: render only during drag, inertia, settle, hover transitions or a state motion; otherwise the loop sleeps.
- Out:
  - Tilting the table or moving the camera on other axes; moving the bot.
  - Changing state from the hero (approve, retry): stays in Missions and Inbox.
  - Beams or lines between the bot and planets (item_162 out of scope).

# Acceptance criteria
- AC1: Switching fixtures between the 5 cases changes planets, moons, stalks, rim and core mood with the item_138 motions and nothing else moves.
- AC2: With no input, no state change and no live run, the hero renders 0 frames per second.
- AC3: Drag spins only the planet layer around the vertical axis and settles with a planet at the front; rim, lens and bot stay fixed; the bot reacts to drag and hover.
- AC4: Labels stay upright and readable at every rotation; back labels hide except `YOUR MOVE` tags; no label sits under the head.
- AC5: Hover shows moons from fixture missions; first click focuses with a card, second click or Enter opens Missions filtered to the project; core opens the Desk chat. Same paths by touch and keyboard.
- AC6: Reduced motion steps without inertia. The static hero and the canvas take the same scene object.

# AC Traceability
- item_126 AC4 -> AC1 to AC6. Proof: this item's notes and the owner's approval.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (render path and fallback rule settled in item_126)

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_043_nest_design_scope_and_jarvis_3d_hero_option`, `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
- Plan: `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`
- Depends on: `item_162_holographic_hero_on_the_dot_matrix_background`, `item_137_nest_app_shell_on_core_fixtures`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Turns the hero from a picture into the Home overview.

# Notes
- 2026-11-10: rescoped from orbit nodes to the holotable. Owner chose free spin on the vertical axis only, fixed bot that reacts, hover = lift + moons, click = details then project page. Built after item_162 is approved.
- Build order from item_126: item_155, item_156, item_157, item_158, item_159.
- Sources: vault `CR-0W/Research/nest-design/item126/` (`compare-item126.html`, `compare-spike.html`, `look-ascii.png`, `spike/holo-head.html`, `spike/img/`). Not `item139/jarvis-1-full.png` (owner rejected it as a reference, 2026-11-10). Data map and budgets: item_126.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
