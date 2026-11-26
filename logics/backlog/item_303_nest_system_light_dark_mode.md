## item_303_nest_system_light_dark_mode - System light and dark mode with Settings override
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Add a mode axis to the Nest theme engine. `<html data-mode="light|dark">` follows prefers-color-scheme unless Settings forces a mode.
- Keywords: nest, theme, mode, light, dark, prefers-color-scheme, settings
- Use when: Building the mode switch or debugging which mode Nest paints.
- Skip when: Styling the gundam theme itself (item_304).

# Problem
- Nest themes have a single palette. The gundam theme needs a dark and a light palette picked by the system, with an owner override.

# Scope
- In:
  - Stored choice `nest.mode` = auto | light | dark (default auto), next to `nest.theme` in `src/ui/theme.ts` (setMode, useMode, resolved mode).
  - `index.html` first-paint script sets `data-mode` from the stored choice or `matchMedia('(prefers-color-scheme: dark)')`; live update on system change while in auto.
  - Settings: a Mode control (Auto / Light / Dark) beside the theme picker, with the same preview-then-save behaviour as themes.
  - `theme-color` meta follows the painted `--bg` after a mode change.
  - Themes opt in: a theme with no light block keeps its single palette (modern, cyber unchanged). Theme CSS writes `[data-theme=x][data-mode=light]` for the light tokens.
- Out:
  - Any gundam tokens (item_304). Light palettes for modern/cyber.

# Acceptance criteria
- AC1: With mode auto, toggling the OS dark setting repaints Nest live without a reload.
- AC2: Settings Light/Dark forces the mode; it survives a reload with no flash of the other mode.
- AC3: Modern and cyber render pixel-identical before and after (screenshot compare per runbooks/nest-screenshots.md).
- AC4: Unit tests cover the stored choice, the resolved mode and the system listener.

# AC Traceability
- request-AC2 -> This slice. Proof: AC1, AC2.
- request-AC3 -> This slice. Proof: AC3.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_072_nest_gundam_theme_with_system_light_and_dark`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Side project. Blocks items 304 to 306.

# Notes
- Dependencies: none. First in order.
- Design source: vault folder CR-0W/Research/nest-gundam, design-notes.md (Decided: follow system, override in Settings).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
