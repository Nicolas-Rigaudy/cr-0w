## item_197_hero_sky_planets - Projects as planets up in the sky, with state looks
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Rework the planets so every project floats as a planet in the sky above the table, with clear looks for running, failed and idle.
- Keywords: nest, hero, planets, sky, projects, orbit, failed, design
- Use when: building or reviewing the Nest hero after item_158.
- Skip when: work outside the Nest hero.

# Problem
- The owner's idea was always to see all projects as planets up in the sky. Today they ride the table ellipse and read as one green thing on the floor.
- The failed ('crashed') planet and the green planet do not look good. Planets in general need rework.

# Scope
- In:
  - Concept first: map real project data (state, mission count, recency) to placement and look, as an annotated 2D sketch.
  - Then 2 to 3 options as boards, side by side with bin/crow-compare, owner picks.
  - Build the picked option: planet placement in the sky, per-state looks, moons orbiting their planet.
- Out:
  - Head scale (item_194), head smoothing (item_196).

# Acceptance criteria
- AC1: Every project with missions shows as its own planet above the table on phone, inner and laptop.
- AC2: Failed, needs-you, running and idle planets are told apart at a glance (owner check).
- AC3: Planets, moons and labels never overlap each other or the head (declutter tests + layout shots).
- AC4: Build matches the approved board; owner OK per size.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: owner OK on shots.
- request-AC2 -> This backlog slice. Proof: tests and perf scripts.

# Decision framing
- Product framing: Needed: owner picks the option
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_058_nest_hero_second_pass`
- Primary task(s): (none yet)
- Follows: `item_158_nest_hero_live_data_state_motions_and_actions`

# Priority
- Priority: Medium
- Rationale: owner feedback 2026-11-12 on the item_158 build.

# Notes
- Design pass rules: concept and sketch before any render (design-renders-need-wow), no renders before OK.
- Decides whether the table spins with the planets; that answer feeds item_198.
- Reference: github.com/JIN2-ASCII/-JIN2-ASCII animated pieces (no license: technique only, no copied art).
> Non-semantic edit: reference moved from Open Threads.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Outcome (2026-11-12, commit 58c790a)
- Planets sit evenly on a sky ring parallel to the table, under the head. Altitude by priority (failed highest, done just above the table); the highest sits at the front, a spin shows the rest.
- The 24h table clock and agenda ticks are gone (owner: confusing, agenda lives elsewhere).
- Rim, owner picks from `Research/nest-design/item_197/rings.html`: bold bezel (minor ticks every 10 deg, accent `#` every 45) that turns with the planets; each planet marks the rim with an arc in its colour, brighter on hover and at the front (no extra front marker); a cyan scan sweep every 12 s on the slow tick.
- A planet label that has to move now counts its own planet as an obstacle (declutter), so it never covers its sphere.
