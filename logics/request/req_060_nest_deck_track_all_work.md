## req_060_nest_deck_track_all_work - Nest and deck track all work, missions as auto grouping
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 80%
> Progress: 0%
> Complexity: High
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Nest and crow-deck treat missions as the unit of work. Make every kind of work first-class (live sessions, background runs, parked sessions, MRs/PRs with CI) through one core work feed, and turn missions into an automatic grouping above it.
- Keywords: work feed, sessions, runs, dispatch, parked, MRs, CI, missions, nest, deck
- Use when: changing what Nest or the deck counts as work, or how missions relate to runs.
- Skip when: purely visual hero work (req_059).

# Needs
- The owner does not think of missions as separate. Crow must show all of their work.
- Today each surface stitches feeds itself (the hero joins dispatch, sessions and projects in `src/hero/model.ts` since item_199). Every surface needs the same answer.

# Acceptance criteria
- AC1: The core serves one work feed: every live session, background run, parked session, open MR/PR (with CI state and reviews waiting) and mission, each with kind, project, state, since and wait.
- AC2: Handoff, dispatch and Logics-started work create a mission automatically when they carry an end state; its sessions, runs and MRs link to it.
- AC3: Nest hero, mood line, inbox and desk read the work feed; each item shows its kind; items under a mission roll up to its goal; items without one show alone.
- AC4: crow-deck lists read the same feed with the same counts as Nest.

# Backlog (execution order)
1. item_202_core_work_feed_over_sessions_runs_parked_and_mrs (first)
2. item_203_missions_created_automatically_and_linked_to_work (after 202)
3. item_204_nest_shows_every_kind_of_work (after 202; mission rollup after 203)
4. item_205_deck_shows_every_kind_of_work (after 202; mission rollup after 203)
- 204 and 205 can run in parallel. item_200 (hero cap) should land before or with 204.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (this request)
- Primary task(s): (none yet)

# Notes
- Owner decision 2026-11-13: work = live sessions, background runs, parked sessions, MRs/PRs with CI, plus missions. Missions are an automatic optional grouping (goal, end state, delegation contract), never created by hand, never a separate list.
- Read `local/owner-prefs.md` before building. Visual changes need rendered options first.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
