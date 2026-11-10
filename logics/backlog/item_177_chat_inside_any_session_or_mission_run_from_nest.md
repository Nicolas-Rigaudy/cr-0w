## item_177_chat_inside_any_session_or_mission_run_from_nest - Chat inside any session or mission run from Nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: open any running mission (first) or session and read its conversation as it goes, then write into it; built on the item_148 conversation poller and item_170 send.
- Keywords: run conversation, mission chat, session chat, cursor, send text, nest
- Use when: building the conversation view of a run.
- Skip when: butler chat (item_176).

# Problem
- Nest shows a run's state, not what it says. To follow or steer a mission, the owner opens the terminal. req_027 planned the feed; item_148 plans the poller.

# Scope
- In:
  - Core: `GET /api/runs/<id>/conversation?cursor=` from item_148, the same for live sessions by session id; new turns pushed over SSE.
  - Nest: run detail shows the last turns (assistant text, tool calls folded to one line), the item_149 one-line summary on top, a text box that sends through item_170 (mission answer for a parked mission, send for a live session).
  - Missions listed first; sessions and dispatch runs after (owner 2026-11-10: most runs will be missions).
- Out:
  - Live terminal attach. Editing past turns.

# Acceptance criteria
- AC1: A running mission's new turns show in Nest within 5 s.
- AC2: Text typed in Nest lands in that mission or session and its reply shows.
- AC3: A long conversation loads the last 50 turns first and pages back.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: read and write inside a run from the phone.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (req_027, item_148)

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`, `req_027_conversations_feed_for_core_deck_nest`
- Depends on: `item_148_run_conversation_poller_with_cursors_in_core`, `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- Related: `item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: High
- Rationale: steering missions needs their words, not only their state.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
