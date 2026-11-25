## item_293_desk_and_ops_sketch_options_from_real_nest_data - Desk and Ops sketch options from real Nest data
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-25 14:43:51

# AI Context
- Summary: Map real core state to the Desk/Ops model, then draw 2-3 annotated 2D sketch options for the owner to pick. No renders, no build.
- Keywords: nest, desk, ops, sketch, design, missions, queue, afk
- Use when: running the design pass of req_070.
- Skip when: building (after the pick).

# Problem
- Desk and Ops do not give the owner one place to see and act on missions, runs, sessions and queues. See req_070 and the vault note `CR-0W/Research/nest-desk-ops/desk-ops-model.md`.

# Scope
- In:
  - Pull a real snapshot from core (`/api/state`: work, missions, dispatch, runs, sessions, jobs, accounts, usage) and list which unit and state each row maps to. Note gaps (fields core lacks).
  - 2-3 annotated 2D sketches at 1440x900 (Desk and Ops each), plus one 390x844 Desk. Options differ on the open questions in the note (queued group placement, detail view tabs vs scroll, AFK batch screen vs launcher mode).
  - Side-by-side compare with the flip/zoom viewer (`bin/crow-compare`), zoom crops checked for contrast.
  - Ops shows one quota bar per account and slot holders only: no cost or token breakdowns (owner 2026-11-25).
  - Decisions go in the vault design notes; the owner picks with AskUserQuestion.
- Out:
  - Image-model renders, any code in tools/nest.

# Acceptance criteria
- AC1: A data map table (unit, state, source field, gap) in the vault note.
- AC2: 2-3 sketch options in the vault folder, terminal ANSI theme look, real data in every row.
- AC3: The owner picks; the pick and the reasons are written in the note; build items drafted.

# Links
- Request: `req_070_desk_and_ops_handle_missions_runs_and_queues`
- Related: item_292 (first step), item_288 (butler out of the list), item_290 (failure cause), draft C of mission 20261125-plan-nest-desk-bugs (one chat on Desk).

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Implement item_293 (Desk and Ops sketch options from real Nest data). Read logics/backlog/item_293_desk_and_ops_sketch_options_from_real_nest_data.md and the vault note CR-0W/Research/nest-desk-ops/desk-ops-model.md. Map a live core snapshot to the model, then draw 2-3 annotated 2D sketch options into the same vault folder and show them side by side with bin/crow-compare. Ask the pick with AskUserQuestion. No renders, no code.
```

# Decision framing
- Product framing: Needed (owner picks a layout)
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
