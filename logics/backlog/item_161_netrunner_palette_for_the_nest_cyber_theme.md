## item_161_netrunner_palette_for_the_nest_cyber_theme - Netrunner palette for the Nest cyber theme
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: netrunner, palette, nest, cyber, theme
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- The cyber theme uses a saturated navy and purple base, a large teal radial, glow on every number and scanlines at 0.18.

# Scope
- In:
  - Cyber tokens in src/theme.css move to the B-netrunner values (base #05070c, panels rgba(9,18,24,.92), line #1d3a46, row #0b1820, fg #e2f4ff, dim #8fb3c3, w #2ef2ff, n #ffb02e, f #ff2a4a, d #3dff9e, b #3aa8ff, accent #c05cff).
  - Hard-coded #4fc3ff and #2ef2ff in cyber rules become tokens.
  - Glow only on the logo, the hero and the live chip. No scanlines.
  - Source of the values: the B-netrunner block in vault CR-0W/Research/nest-home-design-pass/palette-variants.mjs.
- Out:
  - Modern theme
  - Hero box (next item)

# Acceptance criteria
- AC1: The laptop cyber Home render matches B-netrunner.png and the owner approves it.
- AC2: The modern theme render is unchanged.
- AC3: No hard-coded cyan or blue hex is left in the cyber rules.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: The laptop cyber Home render matches B-netrunner.png and the owner approves it.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_007_nest_home_design_pass`
- Architecture decision(s): (none yet)
- Request: `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
- Primary task(s): `task_019_orchestrate_the_nest_home_design_pass`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

