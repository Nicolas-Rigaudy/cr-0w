## item_126_plan_the_workshop_nest_build_from_the_approved_concept_art - Plan the Nest Home hero build (Jarvis graph with the droid core)
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10

# AI Context
- Summary: Turn the item_139 hero pick into a build plan: the Jarvis graph (projects orbiting the CR-0W droid core, state as motion) as a raw WebGL2 ASCII shader with bloom, a static fallback, budgets, and how it plugs into the item_137 shell and the item_138 layouts. Planning only.
- Keywords: nest, hero, Jarvis graph, droid core, WebGL2, ASCII shader, bloom, static fallback, budget, crow core
- Use when: planning how the Home hero gets built after item_138 and item_139.
- Skip when: choosing the look (done in item_138 and item_139) or building the shell (item_137).

# Problem
- item_139 picked the Jarvis graph as the Home hero, with the droid from passes 5-7 as its core (visor band hull, slit-eye LED face, pincer arms, dorsal fin, no ring). item_138 fixed where the hero sits: free centre slot on Home, dimmed at 50% behind Missions and Inbox, on laptop and Fold inner.
- Nothing says yet how the graph and the core get built, what drives them, or what they cost on the phone.

# Scope
- In:
  - Data map: every visual element mapped to a crow source. Project nodes (registry), mission edges and their 5 states (crow core store), needs you (inbox rows), core moods idle, run, ask and fail (store and live runs), quota (usage).
  - State as motion: one motion per state from the item_138 state sheet, reduced-motion variant included.
  - Render approach: one raw WebGL2 ASCII shader (no three.js) for the graph and the core, the core head ported from the `src139/pass7/hero7.html` renderer. The static fallback and the rule that picks it.
  - Placement per item_138: Home centre slot, dimmed hero behind Missions and Inbox, both themes (modern and cyber).
  - Actions: a project node opens Missions filtered to it, the core opens the Desk chat, a needs-you node opens its Inbox row.
  - Budgets: shader module lazy-loaded with a size cap, frame time and frame rate caps, pause when hidden, battery rule.
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
- Build items, in order: `item_155`, `item_156`, `item_157`, `item_158`, `item_159`

# Priority
- Priority: Medium
- Rationale: R&D. The shell (item_137) can start with a static hero; this plan fills the slot.


# Data map (AC1)
Source: vault `item139/jarvis-1-full.png` (17 annotations). Idle = 0 fps; only a state change or a live run moves anything.

| # | Visual element | Crow source | Motion / action |
| --- | --- | --- | --- |
| 1 | Faint background graph | `vault.graph` (cached) | 0 fps |
| 2 | Inner orbit ring | `registry.projects` | slow drift only while a run is live |
| 3 | Outer dashed ring | `queue[]` | static |
| 4 | Core to project link | live, waiting or failed mission of the project | bright when one exists, cut when failed |
| 5 | Project node and label | `project.name`, worst `mission.state` | size and brightness from depth; tap opens Missions filtered to it |
| 6 | WORKING node | mission state `working`, `run.step`, `run.elapsed` | 2.4 s pulse ring, orbs on the link per `run.step`, timer |
| 7 | NEEDS YOU node | inbox row of the mission | amber diamond lifts 4 px, ticks once, holds; tap opens its Inbox row |
| 8 | FAILED node | mission state `failed`, `run.error` | broken ring, one 180 ms RGB-split glitch, then still; tap shows `run.error` |
| 9 | SLEEPING node | mission state `sleeping`, `schedule.next` | dotted, 45% dim, 0 fps, wake time shown |
| 10 | QUEUED node | `queue.pos` | hollow node on the outer ring, no motion until a slot frees |
| 11 | Emitter base arcs | `usage.accounts[].5h` | one arc per account; dims when every account is capped |
| 12 | CR-0W core face | core mood: idle, run (live run), ask (`inbox.open_count` > 0), fail (last run failed) | blink every 4-7 s; snap 160 ms on mood change; tap opens the Desk chat |
| 13 | Asks chip on the core | `inbox.open_count` | count only |
| 14 | CR-0W line over the scene | latest `coordinator.message` | typed out once |
| 15 | Glass panel | `ask.question`, `mission.project`, age, one primary action | 220 ms scan-in on tap |

Colours: cyan `#2ef2ff` default, amber `#ffb800` for asks only, red `#ff2e63` for failures only, slate `#6f86b8` for sleeping.

# Render and fallback (AC2)
- 2026-11-10 owner picked B (ASCII hero) on `item126/compare-item126.html` and the GPU ASCII shader as the render path.
- 2026-11-10 spike `item126/spike/holo-head.html` ported the pass 7 head to WebGL2 (cell pass to MRT, glyph atlas pass, two-level bloom). Compare `item126/compare-spike.html`: owner picked the plain port (no projector ring, beam or scanlines), glitch on fail only, and approved the port fidelity.
- Static fallback: SVG graph plus a pre-rendered ASCII head text per mood, 0 fps. Picked when any is true: no WebGL2, `prefers-reduced-motion`, Save-Data, battery under 20% and not charging, Fold cover screen, manual toggle, software renderer (`UNMASKED_RENDERER` matches SwiftShader or llvmpipe). A hidden tab pauses the shader without switching.

# Budgets (AC3)
- Bench (spike, 1440x900 at DPR 2, 360x129 cells): Iris Xe ANGLE GL 6.7 ms, ANGLE Vulkan 8.5 ms, SwiftShader 814 ms. Pass 7 on the CPU: 100-290 ms.
- Bundle: hero module lazy-loaded after first paint, under 30 kB gz including shaders; glyph atlas built at runtime from the app font. No three.js.
- Frame time: under 10 ms per frame on the laptop iGPU at 1440x900; under 16 ms on Fold inner (measured in item_159).
- Frame rate: 0 fps when nothing changes, 12 fps while a run is live, full rate only during a snap or glitch.
- Hidden tab: loop stops on `visibilitychange`, resumes on the next visible frame.
- Battery and Save-Data: static fallback per the rule above.
- Placement: dimmed to 50% behind Missions and Inbox; full brightness when the panels collapse.

# Notes
- Hero sources: vault `CR-0W/Research/nest-design/item139/` (compare page, `jarvis-1-full`, `jarvis-2-in-home`, `pass7/hero-{idle,run,ask,fail}.png`, `pass7/heads-final.png`).
- Layout sources: canvas artifact MyPAqSwJFddMnSR4grm2g6 and vault `CR-0W/Research/nest-design/item138/fold-src/`.
- Nest stack itself: req_031.
- Spike sources: vault `CR-0W/Research/nest-design/item126/spike/` (`holo-head.html`, renders in `img/`). Headless render and bench: `runbooks/headless-webgl-render.md`.
- 2026-11-10 planning closed: build items item_155 to item_159 in that order.
- 2026-11-10 realigned by item_140: the workshop is dropped as hero (item_139 pick), item_125 withdrawn. This item now plans the Jarvis graph with the droid core.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
