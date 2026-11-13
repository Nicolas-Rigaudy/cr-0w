## item_203_missions_created_automatically_and_linked_to_work - Missions created automatically and linked to work
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 75%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Missions become an automatic grouping: `crow handoff`, `crow dispatch` and Logics-started work create one when they carry an end state, and the sessions, runs and MRs they spawn link to it.
- Keywords: mission, handoff, dispatch, logics, end state, delegation contract, rollup
- Use when: changing how missions start or link to work.
- Skip when: feed shape (202) or UI (204, 205).

# Problem
- Missions exist (`contract/mission.go`: goal, end_state, autonomy, budget, gates, runs[]) but almost nothing creates them; they read as a separate concept.

# Scope
- In:
  - Create or reuse a mission on handoff/dispatch with `--end-state`, keyed by project + Logics ref when given.
  - Link: sessions by handoff note / cwd + worktree, runs by run id, MRs by branch.
  - Mission state derived from its linked items (failed > needs > working > done when end state reached).
  - Work items carry the mission id (feeds 202).
- Out:
  - Manual mission creation UI. Changing autonomy/budget rules.

# Acceptance criteria
- AC1: A `crow handoff ... --end-state` creates one mission; a second handoff for the same Logics ref reuses it.
- AC2: The session, run and MR it spawns appear in the work feed with that mission id.
- AC3: Work started without an end state creates no mission and still shows in the feed.
- AC4: Tests cover create, reuse and linking.

# Priority
- Priority: Medium
- Rationale: Grouping adds value once the feed exists.

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
