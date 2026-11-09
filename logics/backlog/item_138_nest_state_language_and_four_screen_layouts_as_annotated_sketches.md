## item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches - Nest state language and four screen layouts as annotated sketches
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
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
- 2026-11-09 owner: reuse the archeotech cyberdeck pack work, especially motion. Read-back section "Cyberdeck carry-over": cyber = netrunner entrance (glitch, scan, land, decode, bracket lock; small surfaces stagger only), modern is not megacorp (owner: Arasaka-inspired, still cyber); modern shapes and motion need their own refs. Same colours in both.
- 2026-11-09 owner ref picks (vault `refs/step2/`): palette blue frame and data with purple for CR-0W and selection (P2, P3, P4, P7). Hub patterns to borrow: coordinator thread ending in a question (H4, H11), live run graph (H1), status words with parked and blocked states and one resume action (H8, H9, H12), mission detail with a properties rail and a run ledger (H10, H13).
- 2026-11-09 owner: modern theme feels like the archeotech-shell base pack (rounded panels, translucent surfaces, mauve active pill, quiet fades; vault `refs/step2/base-pack/`). Default palette is Catppuccin Macchiato with the mauve accent (`themes/archeotech-macchiato`), replacing the deck shades. That is the modern palette.
- 2026-11-09 owner: cyber is neon, not pastel (vault `refs/step2/palette-cyber/s-Z2.png`). Background is a radial haze, teal core to deep purple edges, never black (`#1d6470` > `#174a5c` > `#182650` > `#120b2c`). Frames `#4fc3ff`, live `#2ef2ff`, CR-0W and selection `#c05cff`, needs you `#ffb800`, failed `#ff2e63`. Themes no longer share colours; roles stay the same.
- 2026-11-09 owner: cyber background texture over the haze, fading to the edges: dot matrix 16 px, hex grid, or polar rings and spokes around the core (vault `refs/step2/palette-cyber/cyber-bg-textures.jpg`, D1, D4, D5). Try them per screen or state.
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
- Fold 6: cover 2376 x 968 px, inner 2160 x 1856 px. CSS viewports measured 2026-11-09, see "Fold pass" below.
- Hybrid rationale: Derived from request `req_043_nest_design_scope_and_jarvis_3d_hero_option` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_043_nest_design_scope_and_jarvis_3d_hero_option.md`.
- Generated locally by logics-manager.

# Owner review (2026-11-09, laptop pass 1)
- Sheets and compare page: vault `CR-0W/Research/nest-design/item138/` (generator `src138/gen.py`, page `compare-item138.html`).
- Verdict: good first draft, base kept for both themes. Some rework needed.
- Home density: dense.
- Hero: dimmed at 50% behind Missions and Inbox.
- Fonts: undecided. Compare pairs on the same screen. Cyber needs themed fonts throughout, body text included.
- Missing content to add in pass 2:
  - calendar access with upcoming meeting reminders
  - mail access
  - a stronger Missions dashboard
  - the day's missions timeline
  - a dedicated MR screen
  - reports, and validation of questions and designs
- Fold cover and inner screens wait for pass 2.

# Owner review (2026-11-09, laptop pass 2)
- Canvas: Claude Design artifact MyPAqSwJFddMnSR4grm2g6 ("Nest laptop pass 2"), co-edited by the owner. Rows: Home, Missions, MRs, Agenda and mail, Reports and validation, font pairs. Modern left, cyber right.
- Home: side rails. Hero slot free in the centre, dense rails left and right. Day strip dropped (owner removed it on the canvas).
- Home content: one card per tab, in tab order (Missions, Inbox, MRs, Agenda, Reports), each with its key count, top items and an open link. Blocking items sit in their tab's card. The meeting reminder sits in the Agenda card.
- Cyber Home: owner made the rails see-through and set the hero rings; kept.
- Fonts picked. Modern M4: Michroma display, Fira Code headers and body. Cyber C2: Orbitron display, Oxanium headers and body. Applied on every board.
- Owner: Missions, MRs, Agenda and mail, Reports approved. Go Fold.
- Home (canvas v19): mix of the side-rails version and the one-card-per-tab version. Needs you on top with the primary action inline and a tab key per row; Missions card with big KPI numbers on top; Agenda card with a big countdown header for the next meeting. Every tab card keeps its number, key count and open link. Owner asked for the big reminder header and the numbered KPIs.
- Owner: Nest must be a place to work from, not only a status view (launch, queue and steer runs without a terminal). Three launch layouts compared (butler desk, workbench, pipeline board).
- Desk and Missions split (canvas v23), approved in both themes. Tabs grow to 7: Home, Desk, Missions, Inbox, MRs, Agenda, Reports.
  - Desk (tab 2): butler desk. Needs you and Ready to run left, CR-0W chat centre (dispatch, queue, explain), slots, queue, live run and sessions right.
  - Desk mission view: opening a mission toggles the centre between Chat, Transcript, Terminal and Graph. Terminal is take control (run pauses, pty session, Hand back resumes). Right rail: properties and run ledger.
  - Missions (tab 3): pipeline board (Ready, Queued, Running, Needs you, Done today) with the day timeline and runs per day. Big KPI numbers on top of both tabs.
  - The old Missions dashboard and the three option boards are retired from the canvas.
  - Generators and a copy of the v23 boards: vault `CR-0W/Research/nest-design/item138/pass2-src/` (scripts point at a scratchpad path, update `P` before reuse).
- Next: Fold cover and inner, derived from the laptop boards.

# Fold pass (2026-11-09 decisions)
- Scope: all 7 tabs on both the cover and the inner screen.
- Themes: modern first; cyber restyles the boards after the owner approves them.
- Canvas: new Fold rows on the same Claude Design artifact (MyPAqSwJFddMnSR4grm2g6).
- Viewports measured on the Fold 6 in the browser (DPR 2.63): cover 368 x 812 (screen 369 x 906), inner 707 x 691 (screen 708 x 823). Boards use the browser viewport. Extra height as a PWA goes to the lists.

# Fold pass result (2026-11-10, approved)
- Canvas rows on artifact MyPAqSwJFddMnSR4grm2g6: `fold-<tab>-cover` and `fold-<tab>-inner`, plus a `-cyber` variant of each, for home, desk, missions, inbox, mrs, agenda and reports (28 boards).
- Covers (368 x 812): glance tabs. One tab per screen with its key count, top items and one primary action; the tab bar switches tabs. The first dense pass was rejected as packed, the second as leaving unused space; v2 approved.
- Inners (707 x 691): the laptop boards reflowed to two columns, with the Home hero slot and the Desk mission view toggle kept. Approved after a design-skill pass (better-interface: accessibility, layout, typography, colours, UI).
- Cyber: generated from the modern boards (chamfered panels, Orbitron and Oxanium, scanlines, dimmed haze, neon roles). Approved.
- Generators and a copy of every board: vault `CR-0W/Research/nest-design/item138/fold-src/` (fold.py, glance2.py, polish.py, cyber.py, fix_web.py, bodies/, boards/). Scripts point at a scratchpad path; update the paths before reuse.

# Web pass (2026-11-10, approved)
- Same design-skill pass on the 16 laptop boards (1440 x 900), both themes.
- Contrast fixes. Modern secondary text `#8087a2` and `#6e738d` become `#939ab7` (5.66 on panel, 4.74 on row). Cyber haze gets a dark overlay `rgba(10,12,40,.55)`, panels go to `.85` opacity, secondary text becomes `#8aa0c8` (6.34 or more on panels).
- Dimmed rows use a dim colour (`#a5adcb` modern, `#9fb6e0` cyber) instead of opacity. Cancel buttons get a danger outline. Kanban cards get a 4 px hit extension. Counters and times use tabular numbers. A 2 px focus ring (`#8aadf4` modern, `#2ef2ff` cyber) on every focusable element.
- Rendering, fit and keyboard checks were not run (no browser in the session). Wording was not reviewed.

# Build notes for item_137
- Boards use `<span class="btn">`: the app uses `<button>` (and `<a href>` for navigation).
- Icon-only controls (`.ib`) need an `aria-label`.
- Purple `#c05cff` and red `#ff2e63` text pass on cyber panels (4.96 and 4.65) but fail on the bare haze centre (3.86 and 3.61). Keep that text on a panel.
- Tags stay at 11 px to match the Fold boards.
- Drag on the Missions board needs a button or menu path to move a card.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
