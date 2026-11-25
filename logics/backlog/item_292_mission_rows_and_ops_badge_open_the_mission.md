## item_292_mission_rows_and_ops_badge_open_the_mission - Mission rows and the Ops badge open the mission
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Mission rows on Desk and the Ops tab badge point to an Ops page that no longer shows missions, so a failed or waiting mission leads nowhere.
- Keywords: nest, desk, ops, badge, mission, navigation, failed
- Use when: a mission row or a tab badge leads to a page without the item.
- Skip when: the why text itself (item_290).

# Problem
- Owner 2026-11-25: "ops has a 1 icon but it doesn't point to anything, and clicking the Plan logics item in Desk brings me to Ops; I don't understand what that is for".
- Cause:
  - Ops replaced the Missions page (commit 921d4cb). `router.tsx:19` redirects `/missions/<id>` to `/ops`.
  - Desk sends mission rows to `/ops?mission=<id>` (`screens/Desk.tsx:245-246`, and the Answer action at `Desk.tsx:310`). `screens/Ops.tsx` never reads the `mission` param: it shows day and week stats only.
  - The Ops badge is the count of failed missions (`ui/Shell.tsx:35-38`, `countByState(...).failed`). Ops lists no failed mission, so the badge has no target. Today it counts the desk-bugs plan mission (see item_290 and item_291).
- A mission view exists: the session chat takes a mission target (`ChatTarget` with a mission, `tkey` `chat:m<mission>` in `Chat.tsx`), with Stop through `missionSheets`.

# Scope
- In (owner picked 2026-11-25: open on Desk, first step of req_070):
  - A mission row on Desk opens the mission in the Desk main panel (the existing Chat mission target): state, why (item_290), latest version summary and actions: Answer (inbox card), Retry/Resume, Cancel.
  - The Answer action on a waiting mission opens its inbox card, not Ops.
  - The failed badge moves from Ops to Desk (failed rows live there); tapping it filters Desk to failed rows. Ops keeps no badge.
  - `/ops?mission=` and `/missions/<id>` redirect to `/desk?mission=<id>`.
- Out:
  - Ops page content.
  - Layout restyle of the Desk main panel (owner picks styling separately per design prefs).

# Acceptance criteria
- AC1: Clicking the failed desk-bugs mission row on Desk opens that mission with its cause and actions; no jump to Ops.
- AC2: Answer on a waiting mission row opens its inbox card.
- AC3: Every tab badge count matches rows visible on its target page; Ops shows no badge.
- AC4: Old `/missions/<id>` and `/ops?mission=<id>` links land on the mission.
- AC5: vitest on routing and badges; shots at 1440x900 and 390x844 compared with the boards before showing the owner.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Depends on: item_210 build (`cr-0w.wt/item-210-desk-pick`) and drafts B and C of mission 20261125-plan-nest-desk-bugs: all edit `screens/Desk.tsx`; build after them.
- Related: item_290.
- Part of: `req_070_desk_and_ops_handle_missions_runs_and_queues` (first step; the full redesign follows the item_293 sketch pick).

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_292 (mission rows and the Ops badge open the mission) end to end, after item_210 and the desk-bugs drafts B and C have landed (check git log on main first). Read logics/backlog/item_292_mission_rows_and_ops_badge_open_the_mission.md first, then tools/nest/src/screens/Desk.tsx mission row handling, router.tsx redirect, ui/Shell.tsx badges and screens/Chat.tsx mission target. Keep the detail view minimal: the full layout comes from the item_293 pick. Board-vs-app compare before showing anything.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
