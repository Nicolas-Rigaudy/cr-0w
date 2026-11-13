## item_204_nest_shows_every_kind_of_work - Nest shows every kind of work
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Nest reads the core work feed everywhere: hero planets and cards, mood line, inbox and desk lists. Each item shows its kind; items under a mission roll up under the goal.
- Keywords: nest, work feed, hero, mood, inbox, desk, MissionRow, rollup
- Use when: Nest work lists or counts.
- Skip when: core feed changes.

# Problem
- Nest screens are mission-centric (`MissionRow`, `order.ts`, `Desk`, `Inbox`). The hero joins live feeds ad hoc (item_199).

# Scope
- In:
  - `data/live.ts` reads `work`; replace `heroLiveOf` joins with it; one row component for any work kind with a kind label.
  - Mission rollup row (goal, then its items) when items carry a mission id.
  - Fixtures gain each kind.
- Out:
  - Visual redesign beyond the row kind label; any new look needs rendered options first.

# Acceptance criteria
- AC1: On live data, hero, mood, inbox and desk counts match `/api/state/work`.
- AC2: A session, run, parked session and MR with red CI each show with their kind and the right state.
- AC3: Items of one mission show under its goal; others show alone.
- AC4: `npm test`, `tsc` and `scripts/hero-reel.mjs` pass.

# Priority
- Priority: High
- Rationale: Nest is the daily surface.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_060_nest_deck_track_all_work
- Primary task(s): (none yet)

# Notes
- Owner decision 2026-11-13: work = live sessions, background runs, parked sessions, MRs/PRs with CI, plus missions. Missions are an automatic optional grouping (goal, end state, delegation contract), never created by hand, never a separate list.
- Read `local/owner-prefs.md` before building. Visual changes need rendered options first.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
