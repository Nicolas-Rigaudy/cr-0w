## item_137_nest_app_shell_on_core_fixtures - Nest app shell on core fixtures
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 80%
> Progress: 45%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 11:23:54

# AI Context
- Summary: First nest web app code: an installable React PWA shell with Home, Missions and Mission detail rendered from the item_119 fixtures, in the look approved in item_138 (laptop and Fold boards, both themes), with the Home hero slot from item_126. No live API yet.
- Keywords: nest, React, Vite, PWA, shell, fixtures, Fold 6, missions, inbox
- Use when: starting the nest front end before the core API (slice 3) exists.
- Skip when: the work needs live data.

# Problem
- The nest screens sit in 0.3b slices 3 to 5, behind 10 core items with no UI. The owner wants the web app started sooner (2026-11-09).
- item_119 fixtures already describe the mission record, the run entry and the inbox row. A shell can render them before the API exists and switch to the API in slice 3.

# Scope
- In:
  - App at `tools/nest` (React 19, Vite 8, TypeScript), stack per req_031 decisions: hand-made components over req_038 tokens, no component kit.
  - Installable PWA: manifest, service worker stub, icons.
  - Routes Home, Missions, Mission detail, inbox list read only. The tab bar shows the 7 item_138 tabs; Desk, MRs, Agenda and Reports stay placeholders until their data exists.
  - A data layer with one interface and a fixture adapter that loads the item_119 JSON fixtures. The API adapter comes with slice 3.
  - Fold 6 cover (glance tabs) and inner screens through breakpoints and segments (req_038), laid out from the item_138 boards.
  - Home hero slot: the item_126 2D SVG fallback of the Jarvis graph and droid core, or a static render (`item139/pass7/hero-idle.png`) if item_126 has not planned it yet. Dimmed at 50% behind Missions and Inbox.
  - The item_138 build notes: real `<button>` and `<a href>` controls, `aria-label` on icon-only controls, the web pass contrast tokens and focus ring, purple and red text kept on panels in cyber, a button or menu path for every drag.
  - Bundle size measured and a budget written down.
- Out:
  - Live API, SSE, Web Push, actions that change state.
  - Butler chat, terminal route, Work and Personal screens (J tracks).

# Acceptance criteria
- AC1: `tools/nest` builds and serves locally, and installs as a PWA from Chrome on the phone over the tailnet.
- AC2: Home, Missions, Mission detail and inbox render every item_119 fixture case (working, needs you, failed, sleeping, queued) with no hardcoded sample data in the components.
- AC3: Swapping the fixture adapter for an API adapter touches only the data layer (one interface, typed from the fixtures).
- AC4: The look matches the approved item_138 boards (laptop, Fold cover and inner, both themes); the owner approves it side by side with the boards on the canvas.
- AC5: Bundle size recorded with a budget in the notes.

# AC Traceability
- request-AC1 -> AC1, AC2, AC4. Proof: running shell on the phone and the owner's approval.
- request-AC2 -> AC3. Proof: data layer interface and fixture adapter.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (stack settled in req_031, API ownership in adr_001)

# Links
- Product brief(s): `prod_005`
- Architecture decision(s): `adr_001`
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_031_spike_nest_stack_and_app_shape`, `req_043_nest_design_scope_and_jarvis_3d_hero_option`
- Depends on: `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures`, `item_138_nest_state_language_and_four_screen_layouts_as_annotated_sketches`, `item_139_jarvis_graph_and_droid_workshop_hero_sketches_on_one_compare_page`. Hero build: `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art` (the shell does not wait for it).
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Owner choice 2026-11-09: design first (125, 126), then this shell in parallel with task_016. Starts nest before slice 3.

# Notes
- 2026-11-09 owner picked option 3: design track first, shell once item_126 fixes the look. Not before the art is approved (memory design-renders-need-wow).
- Order: item_138 → item_137, alongside task_016 (item_119 first). item_126 runs in parallel; its first build item fills the hero slot.
- 2026-11-10 item_126 planned the hero: item_155 (static hero) fills the Home slot, item_156 to item_159 add the WebGL2 ASCII hero.
- 2026-11-10 realigned by item_140: item_138 boards approved (canvas artifact MyPAqSwJFddMnSR4grm2g6, sources in vault `CR-0W/Research/nest-design/item138/fold-src/`), hero = Jarvis graph with the droid core (item_139), item_125 withdrawn. The shell can start now; item_126 fills the hero slot later.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_018_nest_app_shell_on_core_fixtures`
