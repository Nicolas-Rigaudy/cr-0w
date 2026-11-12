## item_195_hero_moon_legibility - Moon titles and moons readable on every size
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
- Summary: Make moons and their titles easy to read at every size.
- Keywords: nest, hero, moons, labels, legibility
- Use when: building or reviewing the Nest hero after item_158.
- Skip when: work outside the Nest hero.

# Problem
- In the item_158 sheets the moons are small and their titles hard to read.

# Scope
- In:
  - Moon disc size, title font size and contrast.
  - Declutter keeps working (no overlap, no clipping).
- Out:
  - Planet redesign (item_197); moon looks there follow this item.

# Acceptance criteria
- AC1: Moon titles read at arm's length on phone cover 368x812 (owner check).
- AC2: Declutter tests pass; no text over text in the layout shots.
- AC3: Owner OK on a before/after sheet.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: owner OK on shots.
- request-AC2 -> This backlog slice. Proof: tests and perf scripts.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_058_nest_hero_second_pass`
- Primary task(s): (none yet)
- Follows: `item_158_nest_hero_live_data_state_motions_and_actions`

# Priority
- Priority: High
- Rationale: owner feedback 2026-11-12 on the item_158 build.

# Notes
- Depends on item_194 (cell size). May be folded into item_197 if the sky layout changes moon placement.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Outcome (2026-11-12, commit 58c790a)
- Moons are shaded spheres again (emit 0.3) on the 1.25x grid; titles placed by declutter with `data-hide`.
- Moon titles shown: laptop 3/3, phone and inner 2/3 (one hides when no free spot).
