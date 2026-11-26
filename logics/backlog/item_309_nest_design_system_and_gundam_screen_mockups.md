## item_309_nest_design_system_and_gundam_screen_mockups - Nest design system and gundam mockups for every screen
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.


# AI Context
- Summary: Inventory every Nest UI component, define one design system (component set, tokens, UI/UX rules), and show every screen as a gundam mockup in both modes before any build.
- Keywords: nest, design system, components, coherency, ux rules, gundam, mockups
- Use when: Deciding how a Nest component looks or behaves, or before building item_304.
- Skip when: Building code. This item produces design docs and mockups only.

# Problem
- Nest has too many different components for the same job and no real design system or UI/UX rules (owner, 26/11). Each screen solves rows, panels, buttons, tags and lists its own way.
- A theme on top of that only reskins the inconsistency. item_304 pass 1 was a recolour of cyber.

# Scope
- In:
  - Component inventory: every visual pattern in `tools/nest/src` (screens and ui), grouped by job (row, panel, button, tag, chip, tab, list, sheet, field, chart, empty state, status), with where each is used and duplicates flagged.
  - Design system doc: the reduced component set, each with anatomy, variants, states (hover, active, focus, disabled), sizes and the tokens it reads. Theme-independent structure; gundam is the first skin.
  - UI/UX rules: spacing scale, type scale, hierarchy, density, actions per row, labels, empty and loading states, scroll (mouse first, everything overflowing scrolls), phone vs laptop.
  - Gundam mockups of every screen (Home, Desk, Ops, Inbox and review, MRs and MR review, Agenda, Reports, Settings, Usage, Missions and detail, chat, sheets) at laptop and phone, dark and light, built from the design system components only.
  - Owner review one screen at a time as image columns (today vs proposal).
- Out:
  - App code changes (item_304 builds from this item's approved output).
  - Hero redesign (item_305), charts (item_306).

# Acceptance criteria
- AC1: Inventory lists every component with its usage and the duplicates to merge.
- AC2: Design system doc defines the reduced component set and the UI/UX rules; owner OK.
- AC3: Every screen has a gundam mockup in both modes made only of design system components; owner OK per screen.
- AC4: Mockups pass the hallmark gates (no side stripes, no glow, max 3 fonts) and 4.5:1 contrast for body and state text.

# AC Traceability
- request-AC1 -> This slice. Proof: AC2, AC3.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_072_nest_gundam_theme_with_system_light_and_dark`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: item_304 and every later screen build depend on it.

# Notes
- Dependencies: item_303 (done). Blocks item_304.
- Inputs: vault CR-0W/Research/nest-gundam/design-notes.md (Decided), item304/home-gap-list.md, home-proposal-{rail,header}.png and home-proposal-mock.html (Home look approved 26/11; the header serial tag needs better styling).
- Boards are a style reference: each screen keeps its job and content; components may merge or change when the design system calls for it.
- Outputs go in the vault (Research/nest-gundam/design-system/), never local/state.
- Method: HTML mockups rendered headless, overlaid or side by side with real shots; small crops only to keep context low. One screen per review step.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
