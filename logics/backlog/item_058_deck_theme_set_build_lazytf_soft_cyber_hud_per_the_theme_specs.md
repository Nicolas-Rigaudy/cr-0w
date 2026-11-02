## item_058_deck_theme_set_build_lazytf_soft_cyber_hud_per_the_theme_specs - Deck theme set build: lazytf, soft, cyber, hud per the theme specs
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
  - hud: corner brackets, ◢ TITLE ◣, numbered panels, ruler gauges; accent from item_055 (blue or magenta, not amber).
  - Fix the deviations the mock layer left: HUD gauge tick fixed at 80%, table header rules missing in soft/cyber/hud, group headers sentence case in cyber/hud.
- Out:
  - Tab layouts.

# Acceptance criteria
- AC1: Each theme renders Home and MRs with visible delimitation on the translucent background (req_011 AC2).
- AC2: Only 16 ANSI slots, no fills.
- AC3: t cycles the four themes; shot tests per theme.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.

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

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

