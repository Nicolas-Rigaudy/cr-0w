## task_004_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home - 0.1 deck skeleton: chrome, tabs, theme engine, config, mock Home
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Go and Bubble Tea app tools/crow-deck: chrome, 5 tabs, four themes switched with `t`, config file, mock Home at three width breakpoints.
- Keywords: deck, skeleton, chrome, tabs, theme, engine, config, mock, home
- Use when: checking what the deck skeleton ships and where its tests live.
- Skip when: working on the later deck slices (real data, list pages).

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_016_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home`

# Acceptance criteria
- AC1: Resize tests like view_test.go for both orientations.
- AC2: `t` cycles 4 themes; lazytf is the default.
- AC3: Home renders at 170, 120 and 90 columns per the breakpoints, with empty and error states.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_004_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_004_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...`: 16 tests pass (AC1 TestResizeBothOrientations, AC2 TestThemeCycle, AC3 TestHomeBreakpoints and TestHomeStates, plus fit, no-truecolor, config and single-pane tab walk).
- Renders of every tab, theme, size and state checked by eye against the merged boards (local/state/afk/shots/deck-skel/).
- Finish workflow executed on 2026-10-29.
- Linked backlog/request close verification passed.

# Report
- Shipped tools/crow-deck (theme table, frame kinds, band-based Home layout, overlays, config, `--theme`, `--state` mock states) and `crow deck`.
- Pages 2 to 5 are placeholders; push/post confirm overlay and narrow-decoration drop order are deferred to later slices. Visual choices the boards leave open are recorded in the AFK card.
- Finished on 2026-10-29.
- Linked backlog item(s): `item_016_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home`
- Related request(s): `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
