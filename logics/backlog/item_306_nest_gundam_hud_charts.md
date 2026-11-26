## item_306_nest_gundam_hud_charts - HUD chart styling for the Gundam theme
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
- Summary: Under gundam, Nest charts take the approved radar HUD style: hairline rings or grids, tick marks, square point markers, arc or bar segments in state colours, reticle on the item that needs attention.
- Keywords: nest, charts, hud, radar, gundam, usage
- Use when: Styling HomeHud or Usage charts for gundam.
- Skip when: Adding a new chart widget.

# Problem
- Usage charts (`src/screens/usageCharts.tsx`) and the Home HUD (`src/ui/HomeHud.tsx`) use the default chart look.

# Scope
- In:
  - Chart tokens for gundam (`--c1`..`--c4`, grid, tick colour) in both modes.
  - Theme-driven marker shape (square), tick marks on axes, hairline grids, reticle on the current value. Prefer CSS and small token-driven props over a forked component.
  - Reference: `states-3-*.png` radar and `local/state/shots/nest-gundam/states.html`.
- Out:
  - A new "runs radar" Home widget. Suggest it to the owner as a separate item.

# Acceptance criteria
- AC1: Usage and Home HUD charts under gundam match the radar board style in both modes; owner OKs board vs app.
- AC2: Other themes keep their chart look.
- AC3: Chart text stays readable (4.5:1) in both modes.

# AC Traceability
- request-AC4 -> This slice. Proof: AC1.

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

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
