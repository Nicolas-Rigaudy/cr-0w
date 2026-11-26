## item_305_nest_gundam_hero_recolour - Hero recolour under the Gundam theme
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
- Summary: The ASCII droid head and planets read gundam tokens in both modes and repaint on a mode change. Recolour only.
- Keywords: nest, hero, droid, recolour, gundam, mode
- Use when: Hero colours look wrong under gundam or after a mode switch.
- Skip when: Redesigning the droid as a Gundam helmet (later, not planned).

# Problem
- The hero reads computed CSS colours once (`src/hero/Hero.tsx`, `RailHead.tsx`, `boot/boot.ts`). A live mode change would leave the old palette.

# Scope
- In:
  - Hero, RailHead and boot read state colours from tokens and re-read on theme or mode change.
  - Gundam values: dark = white/grey shell with HUD state colours; light = RX palette on the light canvas, readable on #f6f5f1.
- Out:
  - New droid art, helmet shape, layout changes.

# Acceptance criteria
- AC1: Switching mode in Settings repaints the hero within one frame, no reload.
- AC2: Hero states (needs, working, failed, sleeping, queued, done) show the gundam state colours in both modes; owner OKs a screenshot pair.
- AC3: Modern and cyber hero unchanged.

# AC Traceability
- request-AC4 -> This slice. Proof: AC1, AC2.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_072_nest_gundam_theme_with_system_light_and_dark`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Polish after the theme lands.

# Notes
- Dependencies: item_304.
- Later idea (not in scope): Gundam-helmet style droid, see vault folder CR-0W/Research/nest-gundam design-notes.md Later.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
