## item_162_holographic_hero_on_the_dot_matrix_background - Holographic hero on the dot-matrix background
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 40%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: holographic, hero, dot, matrix, background
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- The .hero-box has its own radial background and frame. It sits on top of the page and hides the dot matrix.

# Scope
- In:
  - Remove the hero box background and frame. The page dot matrix runs under the hero.
  - The head reads as a projection: additive or screen blending, a faint cyan cone or vignette, optional slow scan line on the head only.
  - The prompt line and state chips float without a box.
  - Holotable (owner OK 2026-11-10, sketch vault Research/nest-home-design-pass/holotable-sketch.html): a projector lens at the bottom centre and a galaxy-map table ellipse. The CR-0W head hovers above the lens and replaces the old orbit ring.
  - Planets = projects placed freely on the table. Size = open missions. Colour = most urgent state.
  - A project that waits on you rises on a stalk. Stalk height = wait time. Tag `YOUR MOVE · <wait>`.
  - Done projects lie flat as faint green rings. Failed projects show as a broken red planet with a cause label. Working projects show run time in the label.
  - Bot intent: the head turns toward the planet with the longest wait and its eyes take that planet's state colour.
  - Rim = 24h ruler (merged from item_164): hour ticks, the day so far as a lit arc, a now notch, amber notches for agenda events.
- Out:
  - Beams, lanes or lines between the bot or lens and the planets
  - Stage lanes or columns of the old war table
  - Panel corner brackets (dropped with item_164, too crowded)

# Acceptance criteria
- AC1: The dot matrix is visible behind and around the head on laptop cyber and modern.
- AC2: Head and labels stay readable (contrast checked).
- AC3: The owner approves it side by side with the palette render.
- AC4: Every planet, stalk, ring and rim notch maps to a fixture value (projects, missions, wait times, agenda).
- AC5: The glow budget holds (at most 4 glowing elements).

# AC Traceability
- request-AC2 -> This backlog slice. Proof: AC1: The dot matrix is visible behind and around the head on laptop cyber and modern.

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

# Validation
- Build: .hero-box has no background. Both themes paint the dot matrix (--dots) on .main at laptop width. ::before is a faint --c-w vignette, ::after a soft conic beam from the bottom edge, both under the canvas (isolation, z-index -1). Cyber hero-say gets a dark halo instead of a glow. No scan line: the hero stays one static frame.
- AC1: dots visible behind and around the head on cyber and modern renders.
- AC2: label contrast on the brightest backdrop (vignette peak + beam): cyber dim 6.5, fg 12.9, needs 8.0, failed 3.9 (the failed label sits outside the peak, about 5:1 there, with a dark halo); modern all >= 4.3.
- AC3: first pass (no box, dots, cone) compare page vault Research/nest-home-design-pass/item162-compare.html. Owner asked for the holotable on top of it; holotable build pending.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

