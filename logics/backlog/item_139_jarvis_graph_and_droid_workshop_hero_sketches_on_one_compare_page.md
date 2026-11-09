## item_139_jarvis_graph_and_droid_workshop_hero_sketches_on_one_compare_page - Jarvis graph and droid workshop hero sketches on one compare page
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 14:01:24

# AI Context
- Summary: Draw the two Home hero directions, the Jarvis-style 3D graph and the droid workshop, as annotated 2D sketches on the same real crow data, put them on one compare page and get the owner's pick.
- Keywords: jarvis, graph, droid, workshop, hero, sketches, compare, page
- Use when: choosing the nest Home hero.
- Skip when: state language or screen layouts (item_138), or building the hero in code.

# Problem
- The workshop hero (req_040) needs AI art for every droid and state, and that round (item_125) is paused as oversized.
- The owner named a second direction: a Jarvis-style 3D graph built in code (refs: Jarvis 3D space panels, cyberpunk game UI, Capri-style AI apparitions, the Obsidian 3D graph).
- Per the design-renders-need-wow lesson, a direction is judged on how it maps real data, before any render.

# Scope
- In:
  - One shared data set: the real registry projects, a handful of missions from the item_119 fixtures, all 5 states, one active run, one inbox item.
  - Jarvis sketch, following the req_043 mapping: CR-0W core, project clusters, mission nodes, state as motion, runs as pulses on links, inbox and mission detail as glass panels, vault notes as a faint background graph. Annotate the 2D fallback (panels become sheets).
  - Workshop sketch on the same data, reusing the req_040 cast picks and stations, no new art.
  - Each hero in a phone (cover) view and an inner-screen view, every element labelled with the data it shows.
  - A short note per hero: build cost, art or quota needs, bundle and battery risk.
  - Both on one compare page (bin/crow-compare); ask the owner with a prompt, options: Jarvis, workshop, a named mix.
- Out:
  - AI renders, image model prompts, 3D code, the modern theme.

# Acceptance criteria
- AC1: An annotated 2D sketch of each hero (Jarvis graph, droid workshop) on the same real crow data: the same projects, missions and the 5 states, in a phone view and an inner-screen view.
- AC2: Both sketches go on one compare page (bin/crow-compare), and the owner picks a hero or a mix. The pick is recorded in this item and in req_043.

# AC Traceability
- request-AC1 -> AC1. Proof: the two sketch files.
- request-AC2 -> AC2. Proof: the compare page and the recorded pick.

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.

- 2026-11-09 owner approved the step 2 read-back (vault `CR-0W/Research/nest-design/Nest Vision Readback.md`) with edits: blue base everywhere (red and amber only for needs you and failed), FiraCode plus a condensed uppercase header face. Core form: sketch shard glyph, holo bust, full light figure, and a pixel or text face (eyes and mouth only, lip-syncs in voice mode, DedSec mask style). Open: hero visible dimmed on Missions and Inbox or Home only.
# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`
- Related: `req_040_nest_deck_two_themes_3d_world` (workshop cast picks, vault `Nest Prompts v5.md`), `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: The hero choice gates item_126 and item_140. Runs in parallel with item_138.

# Notes
- References: vault `CR-0W/Research/nest-design/Nest Hero References.md` and `Nest Moodboard.html`. The owner prunes the moodboard before the sketches start; use only the kept images.
- Reference layers: Iron Man HUD (panels on a sphere), Capri (the core as an apparition), Star Wars holoprojector (projection from a base, flicker for failed, dim for sleeping), Cyberpunk 2077 (glitch as a state), owner's Obsidian graph (real density).
- Hybrid rationale: Derived from request `req_043_nest_design_scope_and_jarvis_3d_hero_option` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`.
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
