## item_249_nest_feels_smooth_on_every_screen - Nest feels smooth on every screen
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 80%
> Progress: 40%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Measure and fix the lag on the laptop Desk and other screens, not only the Home hologram.
- Keywords: performance, lag, webgl, rail head, polling, re-render, nest
- Use when: Nest feels slow, scrolling or clicks lag.
- Skip when: the Home hologram look itself (item_176).

# Problem
- Owner 2026-11-18: the hero is laggy as known, and the Desk sessions and the rest also feel slow and not smooth.
- Lead suspect: since item_176 slice 2 the butler rail draws a WebGL head (src/hero/RailHead.tsx) on every non-Home laptop screen. It renders every 33-50 ms (AMBIENT_MS/IDLE_MS) for as long as the tab is visible, even when idle and even when the rail is folded. That is 20-30 GL frames per second on the main thread next to React.
- Second suspect: data polls (api.ts poll 5000, useConversation, useButler) re-render whole screens on each tick, even when the data did not change.
- Nothing is measured yet.

# Scope
- In:
  - Measure first: Chrome performance trace on the Desk and Home, 10 s idle and 10 s scrolling, with and without the rail head (a `?nohead` flag or a stubbed head). Note frame time, long tasks, scripting vs rendering. Save numbers in the item notes.
  - Rail head: stop the loop when idle (render on change only, short ambient bursts), stop it when the rail is folded or off screen (IntersectionObserver), cap it at a low rate off Home, or show a static frame. Pick from the trace.
  - Polling: skip setState when the payload is unchanged (compare a hash or the raw text), memoize heavy derived lists (work groups, ready rows) and rows (React.memo).
  - Check CSS cost (backdrop-filter, large shadows, animations on many rows) if the trace shows paint time.
- Out:
  - Hologram redesign. Phone and Fold beyond what the shared fixes give.

# Acceptance criteria
- AC1: Before/after trace numbers for Desk and Home are in the notes.
- AC2: An idle Desk spends under 5% main-thread time (no constant GL loop off Home when nothing changes).
- AC3: A poll with unchanged data causes no screen re-render (React profiler or a render counter test).
- AC4: Scrolling the Desk lists and clicking a session feel instant to the owner; vitest and shots pass.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: bounded delivery slice.
- request-AC2 -> This backlog slice. Proof: promotable backlog item.
- request-AC3 -> This backlog slice. Proof: delivery chain includes a task-ready backlog item.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Owner feedback on the live app, 2026-11-18; lag hits every screen.

# Notes
- Generated locally by logics-manager.

## Measurements 2026-11-25
- Probe: `npm run perf -- 10` in tools/nest (scripts/nest-perf.mjs). Headless Chromium on the Intel Iris Xe GPU (ANGLE GL), 1440x900 dpr 2, fixtures, focused window, 10 s per run. "stub" refuses webgl2 so the head never starts (Home then draws its bake hero, so Home stub is not a fair baseline). Main busy = renderer main thread TaskDuration over wall time. All-process CPU = renderer + GPU + browser processes, where GL and raster work land.
- Scrolling: frame p50/p95 stayed at 16.7/16.8 ms and no long task showed on any screen, before or after. Headless scrolling does not reproduce the lag the owner feels.

Before (main at b96fbbf):

| Screen | Mode | Head | Main busy | Script | All-process CPU | GL draws |
| --- | --- | --- | --- | --- | --- | --- |
| Home | idle | hero | 16.8% | 15.6% | 28.2% | 289/s |
| Home | scroll | hero | 19.1% | 17.3% | 35.0% | 321/s |
| Desk | idle | on | 2.2% | 0.8% | 11.5% | 164/s |
| Desk | idle | stub | 0.2% | 0.0% | 4.0% | 0/s |
| Desk | scroll | on | 3.9% | 1.0% | 17.7% | 162/s |
| Desk | scroll | stub | 2.8% | 0.4% | 11.6% | 0/s |
| Ops | idle | on | 10.5% | 0.8% | 39.0% | 163/s |
| Ops | idle | stub | 8.6% | 0.0% | 31.5% | 0/s |
| Ops | scroll | on | 10.5% | 0.8% | 36.7% | 163/s |
| Inbox | scroll | on | 4.1% | 1.0% | 17.8% | 163/s |
| Inbox | scroll | stub | 3.1% | 0.4% | 12.6% | 0/s |

- Largest cause off Home: the Ops live hatch animated `background-position`, a full repaint every frame (31.5% CPU on an idle Ops with no head).
- Rail head: about +2% main thread and +7.5% all-process CPU on any laptop screen (9 GL passes per frame, about 18 fps while idle). It never stops while the window is focused: the idle head breathes, drifts and sways by design.
- Home: the hero costs 17-19% main thread, mostly script. Hologram look and cost belong to item_176.

After the hatch fix (c4c4fb0, transform on a pseudo-element; `npm run perf:hatch` finds 0 visibly different pixels, only edge anti-aliasing up to 13 levels at 1440 and none at 368):

| Screen | Mode | Head | Main busy | Script | All-process CPU | GL draws |
| --- | --- | --- | --- | --- | --- | --- |
| Ops | idle | on | 4.0% | 0.8% | 17.3% | 164/s |
| Ops | idle | stub | 1.3% | 0.0% | 13.9% | 0/s |
| Ops | scroll | on | 7.5% | 1.0% | 26.9% | 163/s |
| Ops | scroll | stub | 6.9% | 0.4% | 23.4% | 0/s |

- Desk, Inbox and Home are unchanged by this fix (same numbers as before).
- AC2 status: an idle Desk is at 0.2-2.2% main thread (under 5%), but the GL loop keeps running. Stopping it at idle freezes the breathing and sway: a look change, paused for the owner.
- AC3 (polls with unchanged data) not started: fixtures do not poll, it needs a render counter test over the api source.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-25T09:12:13.195208+00:00 | afk | afk-nest-lag | Measure, fix the Ops hatch repaint | succeeded |  | Probe and before/after in Notes; head idle loop paused on a look decision; AC3 not started |

