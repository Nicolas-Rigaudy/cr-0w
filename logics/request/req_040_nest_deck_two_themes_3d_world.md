## req_040_nest_deck_two_themes_3d_world - Two themes and a 3D mission world for nest and deck
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70%
> Confidence: 60%
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

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
