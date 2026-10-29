## item_016_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home - 0.1 deck skeleton: chrome, tabs, theme engine, config, mock Home
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Bubble Tea app `crow deck` with the approved chrome, 5 tabs, theme table (lazytf, HUD, Cyber, Soft) switched with `t`, config file and `--theme`, Home drawn from mock data.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Go module in tools/, ANSI 16 colors only, aspect-ratio panes, landscape and portrait layouts.
  - Theme table: role-to-slot colors, border/title style, tab style, flourish glyphs. One accent per theme from the 16 slots, set in config (defaults: lazytf blue/magenta, HUD cyan, Cyber magenta, Soft cyan).
  - No background fills except the selected row (accent fill, dark text). No truecolor, no pixel-only effects.
  - Width breakpoints: 150+ landscape, 100 to 149 portrait stack, under 100 one pane with `tab`.
  - Chrome: header with health dots, sync dot states (green, yellow stale, red with source), unified status bar order, counter rule.
  - Empty, error, loading and stale states; `?` help, `:` palette, `t` theme overlays.
  - Glyphs per spec section 2 (sparkline verified; `‖` not `⏸`, `▶` not `▸`).
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: Resize tests like view_test.go for both orientations.
- AC2: `t` cycles 4 themes; lazytf is the default.
- AC3: Home renders at 170, 120 and 90 columns per the breakpoints, with empty and error states.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: Resize tests like view_test.go for both orientations.
- request-AC2 -> This backlog slice. Proof: AC2: `t` cycles 4 themes; lazytf is the default.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Ships on its own; order follows the build plan in the spec.
