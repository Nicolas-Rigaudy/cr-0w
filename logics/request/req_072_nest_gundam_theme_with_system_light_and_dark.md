## req_072_nest_gundam_theme_with_system_light_and_dark - Nest Gundam theme with system light and dark
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Add a "gundam" Nest theme. Dark = black and white with a mauve accent and full HUD state colours. Light = RX-78-2 (white, blue, red, yellow). Nest gains a light/dark mode that follows the system, with an override in Settings.
- Keywords: nest, gundam, theme, light, dark, prefers-color-scheme, rx-78-2, hud
- Use when: Building or reviewing the gundam theme, Nest light/dark mode, or its charts and hero colours.
- Skip when: Working on the archeotech shell pack (archeotech-dotfiles item_129) or the droid hero redesign.

# Needs
- A Gundam theme for Nest that matches the owner's archeotech mecha wallpaper.
- Nest has one look per theme today. The owner wants the theme to switch dark/light with the system, and to force a mode from Settings.

# Context
- Design is decided in the design session: vault folder CR-0W/Research/nest-gundam, `design-notes.md` (Decided section is final, never re-ask).
- Approved boards: `mix-X-dark.png` / `mix-X-light.png` (layout and plates), `states-3-dark.png` / `states-3-light.png` (state colours, radar chart). Canvas uses Y's blueprint grid (`mix-Y-*.png` background only).
- Mockup sources: `local/state/shots/nest-gundam/mix.html` (?v=X&m=dark|light) and `states.html` (?o=3&m=...). Copy values from there.
- Theme engine: `tools/nest/src/ui/theme.ts`, `tools/nest/src/themes/modern.css` and `cyber.css`, component overrides in `tools/nest/src/theme.css`, first paint in `index.html`.
- De-AI rules: hallmark skill (`vendor/plugins/crow-design/skills/hallmark`), run sloplint on the theme CSS.

# Acceptance criteria
- AC1: Settings offers a "gundam" theme. Picking it gives the approved look in both modes.
- AC2: Nest follows the system light/dark setting, and Settings can force light or dark. The choice survives a reload and paints at first frame.
- AC3: Modern and cyber look exactly as before.
- AC4: Hero and charts read the gundam colours in both modes.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Scope
- In: mode engine and Settings control, gundam theme CSS (both modes), hero recolour, HUD chart styling.
- Out: Gundam-helmet droid hero (later), archeotech shell pack (archeotech-dotfiles item_129, after Nest), a new "runs radar" Home widget (suggest separately).

# Execution order
1. item_303 mode engine (no visual change for existing themes).
2. item_304 gundam tokens and plates (needs 303).
3. item_305 hero recolour (needs 304).
4. item_306 HUD charts (needs 304).
Each item runs in its own fresh session from its backlog doc. One screen at a time: show board vs app before moving on.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/nest/src/ui/theme.ts`
- `tools/nest/src/themes/cyber.css`
- `tools/nest/src/theme.css`
- `tools/nest/index.html`
- `runbooks/nest-screenshots.md`

# Backlog
- `item_303_nest_system_light_dark_mode`
- `item_304_nest_gundam_theme_tokens_and_plates`
- `item_305_nest_gundam_hero_recolour`
- `item_306_nest_gundam_hud_charts`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
