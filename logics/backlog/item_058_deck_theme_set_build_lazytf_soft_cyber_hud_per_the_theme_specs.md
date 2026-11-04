## item_058_deck_theme_set_build_lazytf_soft_cyber_hud_per_the_theme_specs - Deck theme set build: lazytf, soft, cyber, hud per the theme specs
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 90%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build the four themes (lazytf, soft, cyber, hud) per runbooks/tui-themes-ux.md, each with its own drawing language, not just colors.
- Keywords: deck, theme, set, build, lazytf, soft, cyber, hud, per, specs
- Use when: theme work in tools/crow-deck/theme.go.
- Skip when: layout work.

# Problem
- Themes today change little beyond colors; Soft has no delimitation on a translucent terminal.

# Scope
- In:
  - All four themes kept (owner pick 2026-11-02).
  - lazytf: round boxes, blue focus, reverse selection, thin gauges.
  - soft: top rule + left spine per panel, chunky gauges, sentence case.
  - cyber: square, double magenta on focus, // TITLE //, [////....] gauges, caps.
  - hud: corner brackets, ◢ TITLE ◣, numbered panels, segment gauges; accent from the terminal theme (see Final spec).
  - Fix the deviations the mock layer left: HUD gauge tick fixed at 80%, table header rules missing in soft/cyber/hud, group headers sentence case in cyber/hud.
- Out:
  - Tab layouts.

# Acceptance criteria
- AC1: Each theme renders Home and MRs with visible delimitation on the translucent background (req_011 AC2).
- AC2: Only 16 ANSI slots, no fills.
- AC3: t cycles the four themes; shot tests per theme.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.
- AC1 -> every theme draws its own frame (round, top rule + spine, square/double, corner brackets). Proof: tools/crow-deck/themes_test.go TestThemesDifferInMonochrome; shots local/state/afk/shots/themes-2026-11-04/sheets/home-land-mockup-vs-build.png, mrs-land-mockup-vs-build.png.
- AC2 -> foreground slots 0-15 only, no background codes. Proof: TestThemesUseANSI16Only, TestThemesShareAccent.
- AC3 -> t cycles lazytf, soft, cyber, hud. Proof: view_test.go theme cycle test, TestShots writes each theme per page and size.
- Final spec -> gauges differ in monochrome (━─, ●○, [///...], ▕▮▯│▏). Proof: TestThemesDifferInMonochrome, TestHUDGaugeTick, sheets/mono-home-mockup-vs-build.png.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Owner found the themes half baked.

# Notes
- Depends on item_055 for the HUD accent. Execution order: after item_057.
- Mock layer: local/state/shots/deck-rework/themes.py, renders png/theme-*.png.

# Final spec (owner approved 2026-11-03)
- Superseded accent note: no per-theme accent. Colors come from the user's terminal theme (accent slot + ANSI 16). Themes must differ in monochrome by shape: lazytf ━─ lines, Soft ●○ dots, Cyber [///...] slashes, HUD ▮▯ segments with ▕ start and │ 80% tick (▰▱ not in FiraCode).
- Reference: local/state/shots/deck-rework/themes.py, mono.py (strip colors to check), png/theme-*-*.png, png/mono-grid.png.
- Comparison page: local/state/shots/deck-rework/r2-compare.html

# Open choices (owner review)
- HUD gauge fill keeps the semantic slots (ok/warn/bad), the mock painted slot 3 because amber was the HUD accent then.
- Soft panels sit 2 columns apart with no blank row, as in the mockup (the runbook says 1 blank row).
- HUD 80% tick is on every gauge, cost bars included (no policy cap data yet).
- Tab badges keep the current style; the HUD ●N badge in ANSI slot 1 is not built.
- Theme-specific empty, loading and error text and the ✕ glyph are not built.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-04T13:10:21.492582+00:00 | afk | themes-2026-11-04 | Four themes by shape, tests, shots | succeeded |  | Shots in local/state/afk/shots/themes-2026-11-04/sheets. Paused choices wait for owner review. |

