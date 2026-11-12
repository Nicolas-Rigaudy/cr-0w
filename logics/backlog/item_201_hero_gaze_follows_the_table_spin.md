## item_201_hero_gaze_follows_the_table_spin - Hero gaze follows the table spin
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: small
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: While the owner drags the holo table, the head looks at the pointer (the hand) while the planets turn the other way. The gaze should follow the planets' motion instead.
- Keywords: hero, head, gaze, drag, spin, look
- Use when: changing where the head looks during a drag or a coast.
- Skip when: planet data (req_059).

# Problem
- Drag sets `look` from the pointer (`Hero.tsx` step input) and the arms trail the drag (`headCtl.ts` trail spring). The table rotates opposite to where the head looks, so the head reads as disconnected from the table.

# Scope
- In:
  - During a drag and the coast after release, drive the gaze from the spin velocity (`spin.vel` sign and size), leading in the direction the front planets move.
  - Optional: lock on the planet crossing the front of the table while it passes, then pick up the next.
  - Back to pointer or idle look once the spin settles.
- Out:
  - Arm trail and claw spread (keep as they are unless they fight the new gaze).

# Decision framing
- Product framing: Needed. Show two variants in a short reel side by side (gaze leads the spin direction vs gaze tracks the front planet) and let the owner pick.
- Architecture framing: Not needed

# Acceptance criteria
- AC1: In the drag beats of `scripts/hero-reel.mjs`, the eyes and head turn the same way the front planets move.
- AC2: During the coast the gaze eases back as the spin slows; no snap when the spin stops.
- AC3: The owner picked the variant from side-by-side reels.
- AC4: `npm test` and `tsc` pass.

# AC Traceability
- AC1, AC2 -> reel review. AC3 -> owner pick. AC4 -> tests.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_058_nest_hero_second_pass`
- Related: item_198_hero_idle_follow_animations
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Small, visible polish; independent of req_059.

# Notes
- Raised by the owner 2026-11-12.
- Code entry points: `tools/nest/src/hero/Hero.tsx` (step input around `drag: s.drag?.moved ? s.spin.vel`), `tools/nest/src/hero/headCtl.ts` (`look`, `drag`, trail spring).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

