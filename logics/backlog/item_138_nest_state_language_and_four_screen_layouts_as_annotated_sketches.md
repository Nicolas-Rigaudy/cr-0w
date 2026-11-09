## item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches - Nest state language and four screen layouts as annotated sketches
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 14:39:33

# AI Context
- Summary: Sketch the nest state language (working, needs you, failed, sleeping, queued) and the Home, Missions, Mission detail and Inbox layouts on the Fold 6 cover and inner screens, independent of the hero.
- Keywords: nest, state, language, four, screen, layouts, annotated, sketches
- Use when: designing nest states or screen layouts before item_137.
- Skip when: choosing or drawing the Home hero (item_139).

# Problem
- Nest has tokens (req_038) but no agreed way to show a mission's state at a glance, and no screen layouts.
- The 4 screens and the 5 states are needed whatever hero the owner picks, so they should not wait for item_139.

# Scope
- In:
  - State language: one treatment per state (working, needs you, failed, sleeping, queued) using motion, shape and label first, colour from req_038 tokens second. Readable in grey scale and at row size.
  - Layouts for Home (hero slot left as a labelled box), Missions, Mission detail and Inbox, each on the Fold 6 cover screen and the inner screen.
  - Content from the item_119 fixtures and the real project registry, every element labelled with the data it shows.
  - Flat annotated 2D sketches (HTML/SVG), one compare page via bin/crow-compare.
- Out:
  - The hero itself (item_139), AI renders, 3D, app code (item_137).

# Acceptance criteria
- AC1: One sheet shows the 5 states side by side at row size and at card size, each annotated with its motion and token.
- AC2: Each of the 4 screens is sketched on the cover and inner screen, with every element labelled by its fixture field.
- AC3: The owner reviews the sheets on one compare page and approves or marks changes; the result is recorded in this item.

# AC Traceability
- request-AC3 -> AC1, AC2, AC3. Proof: the compare page and the owner's recorded decision.

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.
- 2026-11-09 owner: laptop layout first, Fold cover and inner derived from it. Scope to be updated before redrawing.
- 2026-11-09 owner: design from the approved Moodboard v2 and the vault note 'Github repos and websites' (see local/runbooks/design-toolkit.md), never from a generic palette.
- 2026-11-09 first canvas pass (artifact K4pqN8J7QHo69LrSdx3n5U, state sheet + Home cover/inner) rejected by the owner; item paused until step 2 (references and moodboard read-back with confidence per area) is approved.

- 2026-11-09 owner approved the step 2 read-back (vault `CR-0W/Research/nest-design/Nest Vision Readback.md`) with edits: blue base everywhere (red and amber only for needs you and failed; purple accent next to the cyan and neon blue, no Capri magenta or green), FiraCode plus a condensed uppercase header face. Core form: sketch shard glyph, holo bust, full light figure, and a pixel or text face (eyes and mouth only, lip-syncs in voice mode, DedSec mask style). Open: hero visible dimmed on Missions and Inbox or Home only.
- 2026-11-09 owner, step 2 follow-up: themes share colours for now and differ by shapes, effects, decorations, animations (and sounds with voice); density medium, sketch several levels; item_138 sketches the hero both dimmed behind Missions/Inbox and Home only. Extra refs gathered first: talking pixel faces, blue + purple hologram UIs, condensed header fonts, agent hub UIs.
- 2026-11-09 owner font picks (sheet `refs/step2/fonts/fonts-sheet.png`): three layers, display (logo, screen titles, big counters) / header (panels, labels) / FiraCode (data, body). Cyber: display Orbitron, header Chakra Petch or Oxanium. Modern: display Michroma, Quantico or Russo One, header Rajdhani or FiraCode. Sketches try the pairs.
- 2026-11-09 owner: reuse the archeotech cyberdeck pack work, especially motion. Read-back section "Cyberdeck carry-over": cyber = netrunner entrance (glitch, scan, land, decode, bracket lock; small surfaces stagger only), modern = megacorp flat square hairlines with code tabs and a plain fade. Same colours in both.
# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`
- Depends on: `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures` (fixtures), `req_038_nest_design_rules` (tokens)
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Unblocks item_137 whatever hero wins. Runs in parallel with item_139.

# Notes
- References: vault `CR-0W/Research/nest-design/Nest Hero References.md` (Fold 6 screens, jakubkrehel `state-machine` and `variant` skills, reduced-motion rules).
- Fold 6: cover 2376 x 968 px (about 22:9, CSS ~402 x 956 unverified); inner 2160 x 1856 px (nearly square). Measure the CSS viewports on the phone before the layouts are final.
- Hybrid rationale: Derived from request `req_043_nest_design_scope_and_jarvis_3d_hero_option` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`.
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
