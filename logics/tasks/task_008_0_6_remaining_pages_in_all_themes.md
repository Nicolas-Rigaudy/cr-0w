## task_008_0_6_remaining_pages_in_all_themes - 0.6 remaining pages in all themes
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 75%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: History behind `/` in Agents, the disabled post confirm overlay, per-theme overlay shots, and a fit check of every page, theme, state and size in the crow deck. Read only.
- Keywords: history, overlays, themes, fit
- Use when: continuing the deck after the remaining pages, or deciding the narrow-width decoration rule and the Soft density mode.
- Skip when: the change is outside the deck.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_021_0_6_remaining_pages_in_all_themes`

# Acceptance criteria
- AC1: Every tab renders cleanly in all 4 themes, portrait and landscape.

# Plan
- [x] History view with search, chips, date range, preview, live switch.
- [x] Post confirm overlay (disabled) and overlay shots per theme.
- [x] Fit tests for every page, theme, state and size.
- [ ] Narrow-width decoration rule per theme (paused: not in the spec).
- [ ] Soft lower-density mode (paused: not in the spec).

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: History (slash opens, esc leaves, search, chips, date range, mode hides perso, preview without detail, live switch, states fit), post overlay is disabled, overlays fit.
- Renders of every page and overlay in 4 themes under local/state/afk/shots/deck-final/ (landscape, portrait, single panes, states).

# Report
- Shipped History from live crow-usage sessions (no new source). Live rows lack end state, touched files, links and last messages, so the preview says they are not recorded yet. Mock rows carry invented detail.
- The key `m` served both mode and the model chip, so the model chip is `M`. A status chip `s` was added. Both are open to confirm.
- History action keys (enter, o, r, y, h) and the post overlay stay disabled. Nothing writes to a forge or Jira.
- AC1 holds for the shipped pages: the fit tests cover 4 themes at the three breakpoints.
- Open: the narrow-width decoration rule per theme and the Soft lower-density mode are visual choices the spec does not make.

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
