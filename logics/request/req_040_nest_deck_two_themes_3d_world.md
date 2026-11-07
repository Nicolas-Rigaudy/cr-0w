## req_040_nest_deck_two_themes_3d_world - Two themes and a 3D mission world for nest and deck
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 75%
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:03:34

# AI Context
- Summary: Nest and deck keep two themes only, modern (work default) and cyber (personal default). A theme changes look, motion, 2D/3D elements and the butler voice. Nest gets a 3D world that shows missions and agents at work. One focus switch later drives the whole desktop. Design research comes first.
- Keywords: themes, modern, cyber, 3D, world, agents, focus switch, Claude Design, design system, three.js
- Use when: designing nest or deck look, the 3D view, or the work/personal theme switch.
- Skip when: picking the nest stack (req_031) or writing base design rules (req_038).

# Needs
1. **Two themes, everywhere.** Nest and deck ship exactly two themes:
   - Modern: the archeotech shell base glass look and lazytf.
   - Cyber: the archeotech cyberdeck pack, crow cyber and hud.
   The deck's current four themes (lazytf, soft, cyber, hud) fold into these two.
2. **A theme is more than colour.** It sets tokens, type, motion, sounds, the 2D and 3D elements, the agent cast and the butler voice. Modern pairs with the neutral work voice (req_032 decision 20). Cyber pairs with the Steward droid persona.
3. **Focus switch.** Work focus defaults to modern, personal focus to cyber. One focus state that every surface reads (nest, deck, then the archeotech shell, terminal and wallpaper). Phased: nest and deck first. The desktop side is a request in archeotech-dotfiles.
4. **3D mission world.** A view where agents (runs) move between places and the things waiting on the owner pile up visibly. The reference is Thao's Cockpit (rooms per area, factory line per stage). It must fit the owner's themes, not copy that look. What the world represents is open: research proposes the options.
5. **Devices.** Full 3D on desktop and tablet (Fold inner screen). On the phone, a flat 2D or isometric version of the same map with the same data and tap targets per req_038.
6. **Agent cast per theme.** Each theme has its own cast. Cyber: low-poly droids and CRT/HUD props. Modern: characters too, no abstract shapes (owner: "not a fan of abstract shapes for modern"). The modern cast is open to research.
7. **Own look, non-AI design principles.** Grounded in real design craft (composition, hierarchy, restraint, a consistent material and lighting model, hand-picked references), not the generic AI web app look (purple gradients, glass cards everywhere, emoji, stock 3D blobs). No AI-generated images or models in the shipped assets.

# Owner references (2026-11-07)
- Cyber: Cyberpunk 2077, Star Wars (droids, cantina), Dune, Warhammer, cyberdecks, Gundam, Eureka Seven, The Matrix, Tron, Ender's Game, Iron Man, Evangelion, sci-fi in general, r/unixporn rices, Caelestia shell, Aulumu, EDC gear.
- Modern: Apple design, Devialet, Samsung devices.
- First target: the laptop web app, then the phone.
- Thao's Cockpit screenshots are the only external sample so far; more references will come from the owner.

# Design research (first slice)
- References from outside AI tools: game UI and diorama games, architecture models, industrial and product design, sci-fi interfaces and film UI, the owner's own rices. A board per theme with sources.
- 3 or 4 world concepts (what the map means, how agents move, where asks pile up), each drawn in both themes.
- Agent cast sketches per theme, modern characters in at least two directions.
- Two Design System artifacts in Claude Design (modern, cyber), built from the existing tokens (req_038, terminal kit item_079/item_084, archeotech packs). Nest canvases built on them.
- Every choice shown side by side with `bin/crow-compare`, asked with AskUserQuestion. No build before the owner picks.
- Tech check: three.js or react-three-fiber with modelled glTF assets (Blender), cost on the Fold inner screen, reduced motion, a 2D fallback.

## Delivered slice (2026-11-07)
- Delivered: AC1 boards and the renders for AC2 and AC3. The files sit in the hub's private local state: a research folder and a shots folder for the run. They are not committed because they hold reference links and renders.
- Research boards: one per theme with 12 references each, split into two directions. Each board is an HTML page with links plus a PNG.
- World concepts: rooms per area, mission line, dock and tower, district map. Each has a short spec and the same iso camera in both themes, plus a phone top-down 2D version per theme.
- Cast sheets: modern M1 scale-model people and M2 toy crew; cyber C1 utility droids and C2 protocol droids. Same framing, four roles each.
- Contact sheets and `bin/crow-compare` specs for four decisions: world concept, modern cast, cyber cast, board direction.
- Tech note: three.js or react-three-fiber with glTF from Blender, the Fold GPU budget, reduced motion, the 2D fallback.
- Second pass: the modern theme on a light base next to the dark one, a wider cast framing, bloom on cyber, checked reference links, 1920x1080 laptop hero frames, and a fifth decision on the modern base (light or dark).
- Pending: the owner's picks. No build has started.

## Next steps
- AC4: build two Claude Design design systems (modern, cyber) on the picked board directions. Build them from one token set shared with the deck (req_038 status colours, terminal kit, archeotech Protocol register).
- AC5: draw the nest key screens on those systems from the picked world concept and cast.
- Model the picked cast and world kit in Blender as glTF, then measure on the Fold before any effect.

# Owner review of the design passes (2026-11-07)
- Both AFK passes (58d5c24d, d490afb6) are rejected as a whole: "very blocky, not inspired, too copied from Thao, no character, no personal design, no wow factor". The light vs dark base gets the same verdict. No world, cast or board pick stands.
- Cause: procedural Python primitives and a brief anchored on Thao's Cockpit.
- Next: the owner makes concept art in an image model and drops the picks in `local/state/research/nest-design/refs/` with a note on each. The next pass works only from those, without the Thao rooms or factory line. AI images serve as references only (Need 7 still holds for shipped assets).

# Acceptance criteria
- AC1: A research board per theme with cited references and a short "why it fits" per reference.
- AC2: At least 3 world concepts rendered in both themes, compared side by side. The owner picks one.
- AC3: The agent cast for each theme is picked from rendered options, modern with characters.
- AC4: Two Claude Design design systems (modern, cyber) exist and map to one token set shared with the deck.
- AC5: Nest key screens (Home, Missions, 3D world, its phone 2D version) drawn in both themes on those systems.
- AC6: The focus state format (where it lives, who writes it, how surfaces react) is specified. The matching archeotech-dotfiles request exists.
- AC7: The deck theme list is reduced to modern and cyber in the backlog, with the mapping from the four old themes.

# Out of scope
- Building nest screens (after the design picks, as backlog items).
- The archeotech shell implementation (archeotech-dotfiles).

# Context
- Owner request 2026-11-07, after seeing Thao's agent orchestrator (3D rooms and a mission factory line).
- req_038 nest design rules, req_031 nest stack (PWA, Go core), req_032 butler persona, req_036 fun layer.
- Deck UX2 boards are the approved deck look (item_055, item_058, item_086).
- archeotech-dotfiles req_001 theming packs and req_011 cyberdeck pack.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics/request/req_038_nest_design_rules.md`
- `logics/request/req_031_spike_nest_stack_and_app_shape.md`
- `runbooks/tui-themes-ux.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Purpose of the world (owner, 2026-11-07, restart from zero)
After pass 3 (Blender kits) was rejected ("shows nothing the system is doing, random models in groups with no purpose"), the design restarts from the purpose.
- **When:** glanced at all day on a second screen, on return from AFK or a meeting, on the phone while away, and shown to other people.
- **Answers at a glance:** something needs me; the work is moving or stuck; where my quota goes; what happened while I was away (a quick overview or a shortcut to a dedicated view, not the detail).
- **Mode:** mostly watch, with a few actions. An interactive board: zones and droids open their panel (an MR area opens the MR panel, a droid opens its mission), a few launch buttons, like the buttons in Thao's renders.
- **Why a world:** fun, reads better at a glance (a schematic beats a list), and a portfolio showpiece.
- **Order:** purpose, then the content inventory (what the board shows), then the rendering choice. No render before the owner agrees on the content.
- **Content (owner answers):** decisions show as "something needs me" and open a panel or the decision page on click, never the full detail in the world. Runs: simple status. Quota: a small bar, low priority. Away summary: no text, finished work waits in one place (end of a factory line, a spaceport, a spot on a grid). Few actions, mostly a viewer. No vault, Jira backlog or calendar detail.
- **Core idea (owner):** see the flows crow is built from: when things happen and what state they are in. Comparable: a teammate's Jarvis view, 2D nodes per machine and per function with glowing orbs travelling the links between services.
- **No orbs (owner):** a living world where droids walk from one workshop station to the next, keeping the simplicity of the orb view. Dispatch, AFK and MR are one mission workflow, and each kind of mission shows differently on it. Retro, quota and the rest are background items, maybe later.
- **Sketch round (concepts4/sketches.html, flat annotated):** the owner understands all four ideas (hangar, workshop, grid city, holo war table). War table: no cards in hand. Hangar: it read as an airport. It must be a sci-fi Star Wars hangar.

# World pick (owner, 2026-11-07)
- Concept images v4 were made with Gemini from the v4 prompts, one per world in both themes. They sit in the hub's private research folder (refs/v4, compare-concepts-v4).
- **Main world: the droid workshop.** Droids carry their mission crate from bench to bench along one loop: intake hatch, power cells per account, workbenches, scan arch, the owner's desk with an amber lamp (needs you), review counter, rest racks (done). The base is the cyber take 1. Modern is the atelier version with a redesigned original cast (no film droid look-alikes).
- **Kept in mind:** the hangar (biggest wow, ships), the grid city (simplest, closest to a node view) and the war table ("your move" signal). They are candidates for later views or extra rooms, not dropped.
- How the concepts got there: purpose first, then content, then flat annotated sketches, then image prompts that name a game genre and never "dashboard/UI". This replaces the rejected render passes 1 to 3.
- Next: `item_125_workshop_concept_art_round_2_original_cast_states_phone_and_night_views`, then `item_126_plan_the_workshop_nest_build_from_the_approved_concept_art`.
