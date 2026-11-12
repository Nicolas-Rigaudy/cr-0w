## item_198_hero_idle_follow_animations - Hero idle and follow animations, table spin
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Add life to the hero: idle animations, reactions while following the pointer, and possibly a slow table spin.
- Keywords: nest, hero, animation, idle, follow, table, spin
- Use when: building or reviewing the Nest hero after item_158.
- Skip when: work outside the Nest hero.

# Problem
- Spinning moves the planets but not the table, which looks detached.
- The hero has little motion at rest and few reactions when it follows the pointer or when planets move.

# Scope
- In:
  - Owner decision: does the table spin with the planets (rings, ticks) or stay fixed.
  - Idle loop for the head (breathing, blink, glance) within the paint budget.
  - Follow reactions: head and planets respond as the pointer moves.
  - All motion off with reduced motion; timers only while visible.
- Out:
  - Planet looks (item_197).

# Acceptance criteria
- AC1: Owner picks the table spin option and the idle set from a short motion mockup.
- AC2: Idle draw rate stays bounded (`scripts/hero-idle.mjs`), 0 with reduced motion.
- AC3: Draw time stays under 20 ms per frame on laptop.
- AC4: Owner OK on a screen recording per size.

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
- After item_196 (shared follow code) and item_197 (planets in the sky change what moves).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
