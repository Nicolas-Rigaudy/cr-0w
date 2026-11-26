## item_304_nest_gundam_theme_tokens_and_plates - Gundam theme tokens and plate styling
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
- Summary: Ship `src/themes/gundam.css` (dark + light token blocks) and the gundam component overrides: chamfered plates, part-code tags, stencil serial, panel lines, reticle corners, grid canvas.
- Keywords: nest, gundam, theme, css, tokens, rx-78-2, chamfer, hud
- Use when: Building or adjusting the gundam look.
- Skip when: Mode engine (item_303), hero (item_305), charts (item_306).

# Problem
- The approved gundam design exists only as mockups.

# Scope
- In:
  - Token blocks for `[data-theme=gundam]` dark (default) and `[data-mode=light]`, full token set as cyber.css. Values from vault folder CR-0W/Research/nest-gundam design-notes.md and `local/state/shots/nest-gundam/{mix,states}.html`:
    - Dark: canvas #0d0d0f, panels #17171a / #202024, lines #2e2e34 / #3c3c44, ink #ecebe6 / #a3a2a9 / #6e6d75, accent mauve #b99cf2 (ink #17121f, soft #2a2336).
    - Light: canvas #e3e2dd, panels #f6f5f1 / #ecebe6, lines #cfcec8 / #b7b6b0, ink #16171c / #4f5059 / #7a7b83, header RX blue #1d3a8f, actions RX blue, tab underline and V-fin RX yellow #f5b800, alert red #c8102e. No mauve in light.
    - States (Full HUD), dark: w #8fb4ff, n #f2c230, f #ff6470, d #86d3a0, q #8a8991, s #5d5c64, t #b99cf2, o #e8a06a, b #8fb4ff. Light: w #1d3a8f, n #8a6400 (text) / #f5b800 (fills, dark ink), f #c8102e, d #2c6e45, q #7a7b83, s #b7b6b0, t #5b3fa6, o #b4520f, b #1d3a8f.
  - White accents in dark: big numbers, reticle corners, focus outlines, callout rules stay white/grey.
  - Fonts: Barlow Condensed (display, headings, tags, buttons), Barlow (body), B612 Mono (numbers and serials only). Add to index.html.
  - Shapes: chamfered panels and buttons via `--clip-pn` / `--clip-btn` with a 1px line shell; filled chamfered status tags with part codes; faint stencil serial and one panel line on hero panels; reticle corners on the ask input and the hero number; subtle 24px blueprint grid via `--bg-image`.
  - Header: V-fin logo mark colours per mode, active tab underline (mauve dark, yellow light).
- Out:
  - Coloured side stripes on rows or cards (rejected). Tracked small caps, dot-joined meta, glow, card in card.
  - Hero drawing changes (item_305), chart internals (item_306).

# Acceptance criteria
- AC1: Every Nest screen under gundam matches the boards in both modes; owner OKs Home board vs app before other screens (one screen at a time).
- AC2: sloplint on gundam.css passes, and the hallmark gates listed in Out hold (no side stripes, no glow, max 3 fonts).
- AC3: `npm run audit:scroll` passes under gundam in both modes.
- AC4: Text contrast at least 4.5:1 for body and state text in both modes (yellow text uses the ochre token in light).

# AC Traceability
- request-AC1 -> This slice. Proof: AC1, AC2.

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
- Rationale: The theme itself.

# Notes
- Dependencies: item_303.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
