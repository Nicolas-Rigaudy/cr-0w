## req_069_nest_butler_chat_polish_and_failure_clarity - Nest butler chat polish and failure clarity
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Owner feedback 2026-11-25 on the Nest butler chat: raw rendering, the butler shown twice, Sir on every reply, an unexplained fail state and a dead Ops badge.
- Keywords: nest, butler, chat, render, sir, fail, ops badge, mission
- Use when: planning or building any of items 287 to 292.
- Skip when: the Desk tab bugs of mission 20261125-plan-nest-desk-bugs (drafts A to C in its outbox).

# Needs
- Owner feedback 2026-11-25 (screenshot `owner screenshot 20261125_141736.png`):
  1. Butler chat and session chat render the same reply differently.
  2. The butler conversation shows twice: side panel and Desk main panel.
  3. "Sir" on every message is annoying; use it sparingly.
  4. Nest shows a fail state and the butler cannot say why.
  5. Ops has a badge of 1 that points to nothing, and the failed Plan row on Desk opens Ops.

# Context
- The fail state is the desk-bugs plan mission: core rejected its v2 ("responses miss q1, q2, q3") at 14:17. Its plan v2 and question q4 never reached the inbox, although the owner was told it waits there.
- The desk-bugs mission drafts (Desk remembers the pick, drafts survive a remount, one CR-0W chat on Desk) are not duplicated here. Item_288 covers the second source of the double chat (the butler process listed as a run); draft C covers the Desk Ask box.

# Acceptance criteria
- AC1: Each need maps to one backlog item with cause, scope, acceptance checks and a launch prompt.
- AC2: Overlaps with mission 20261125-plan-nest-desk-bugs are links, not copies.
- AC3: Execution order and file conflicts are stated.

# Execution order
1. item_291 (plan resumes accept withdrawn questions): core only; unblocks the desk-bugs mission, which the owner then relaunches.
2. item_290 (failed mission says why): core and persona.
3. item_288 (Desk hides the butler process): core tag plus Nest filter.
4. item_287 (butler chat renders like the session chat) and item_289 (Sir sparingly): independent, can run in parallel.
5. item_292 (mission rows and Ops badge): edits `screens/Desk.tsx`; after item_210 and desk-bugs drafts B and C land. Needs the owner pick on navigation first.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Backlog
- `item_287_butler_chat_renders_replies_like_the_session_chat`
- `item_288_desk_hides_the_butler_process_from_the_session_list`
- `item_289_butler_says_sir_sparingly`
- `item_290_failed_mission_says_why_in_work_feed_and_butler`
- `item_291_plan_resumes_accept_withdrawn_questions`
- `item_292_mission_rows_and_ops_badge_open_the_mission`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
