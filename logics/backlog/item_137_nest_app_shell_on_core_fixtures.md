## item_137_nest_app_shell_on_core_fixtures - Nest app shell on core fixtures
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09

# AI Context
- Summary: First nest web app code: an installable React PWA shell with Home, Missions and Mission detail rendered from the item_119 fixtures, in the look approved by item_125 and planned by item_126. No live API yet.
- Keywords: nest, React, Vite, PWA, shell, fixtures, Fold 6, missions, inbox
- Use when: starting the nest front end before the core API (slice 3) exists.
- Skip when: the look is not approved yet (item_125, item_126), or the work needs live data.

# Problem
- The nest screens sit in 0.3b slices 3 to 5, behind 10 core items with no UI. The owner wants the web app started sooner (2026-11-09).
- item_119 fixtures already describe the mission record, the run entry and the inbox row. A shell can render them before the API exists and switch to the API in slice 3.

# Scope
- In:
  - App at `tools/nest` (React 19, Vite 8, TypeScript), stack per req_031 decisions: hand-made components over req_038 tokens, no component kit.
  - Installable PWA: manifest, service worker stub, icons.
  - Routes Home, Missions, Mission detail, inbox list read only.
  - A data layer with one interface and a fixture adapter that loads the item_119 JSON fixtures. The API adapter comes with slice 3.
  - Fold 6 cover and inner screens through breakpoints and segments (req_038).
  - The visual layer from item_126: the first build item it plans, or the room as a static image if its render approach needs more.
  - Bundle size measured and a budget written down.
- Out:
  - Live API, SSE, Web Push, actions that change state.
  - Butler chat, terminal route, Work and Personal screens (J tracks).

# Acceptance criteria
- AC1: `tools/nest` builds and serves locally, and installs as a PWA from Chrome on the phone over the tailnet.
- AC2: Home, Missions, Mission detail and inbox render every item_119 fixture case (working, needs you, failed, sleeping, queued) with no hardcoded sample data in the components.
- AC3: Swapping the fixture adapter for an API adapter touches only the data layer (one interface, typed from the fixtures).
- AC4: The look matches the item_125 picks and the item_126 plan; the owner approves it side by side with the concept art.
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
- Request: `req_040_nest_deck_two_themes_3d_world`, `req_031_spike_nest_stack_and_app_shape`
- Depends on: `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures`, `item_125_workshop_concept_art_round_2_original_cast_states_phone_and_night_views`, `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Owner choice 2026-11-09: design first (125, 126), then this shell in parallel with task_016. Starts nest before slice 3.

# Notes
- 2026-11-09 owner picked option 3: design track first, shell once item_126 fixes the look. Not before the art is approved (memory design-renders-need-wow).
- Order: item_125 → item_126 → item_137, alongside task_016 (item_119 first).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
