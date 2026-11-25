## req_070_desk_and_ops_handle_missions_runs_and_queues - Desk and Ops handle missions, runs and queues
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 75%
> Complexity: Large
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 20:55:58

# AI Context
- Summary: Redesign the split between Desk and Ops so the owner handles missions, runs, sessions, queues and AFK launches from one model, and sees what needs them at a glance.
- Keywords: nest, desk, ops, missions, runs, sessions, queue, afk, launcher, badge
- Use when: planning Desk or Ops structure, badges, the detail view or launching work.
- Skip when: single bug fixes on Desk (req_069, mission 20261125-plan-nest-desk-bugs).

# Needs
- Owner 2026-11-25: mission rows and the Ops badge lead nowhere. Asked for research on how Desk and Ops should handle missions, runs, sessions, their details and chats, the queues, AFK and background launches, and what needs the owner.
- Direction picked: Desk-first with a shared detail view (option 1 with a bit of a mission view). Research and sketches first, no build.

# Context
- Model and open questions: vault `CR-0W/Research/nest-desk-ops/desk-ops-model.md`.
- Summary: one list of units (mission, run, session, queued job, AFK job) in one state each (needs you, moving, queued, done/failed). Desk = needs you + moving + one detail view, butler in the side panel. Ops = capacity, queues, AFK batch, history. Ops keeps usage light (one quota bar per account); no detailed cost or token breakdowns, those stay in Usage. One launcher (now / background / AFK).
- Today Ops only shows day and week stats; mission detail was lost when Ops replaced Missions.

# Acceptance criteria
- AC1: 2-3 annotated 2D sketch options built from real Nest data; the owner picks one.
- AC2: The picked option is split into build items with file scope and order, after item_210 and the desk-bugs drafts land.
- AC3: Every badge has a target row on its page in the picked option.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Backlog
- `item_293_desk_and_ops_sketch_options_from_real_nest_data`
- `item_292_mission_rows_and_ops_badge_open_the_mission` (first step, builds before the redesign)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
