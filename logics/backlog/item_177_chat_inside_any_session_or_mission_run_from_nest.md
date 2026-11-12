## item_177_chat_inside_any_session_or_mission_run_from_nest - Chat inside any session or mission run from Nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 80%
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
  - Phone first (owner 2026-11-11): the phone is where the owner chats with crow. Every word of crow's replies, line breaks and bullets kept, never cut with "…". No diffs or tool logs: each run of tool calls folds to one line, tap to open. Reply box at the bottom.
  - Design source (approved by owner 2026-11-11): vault `Research/nest-chat/boards/` (`chat-cover.dc.html`, `chat-cover-ask.dc.html`, `chat-inner.dc.html`), built on the item138 fold-src boards. Compare board vs app before the owner sees a build.
  - Cover: full session chat opened from the hero Open or a session row. Back bar with session and state, your bubbles right, full replies left, one dashed fold line per tool-call run, ask card with choice buttons docked above the reply box, nav kept.
  - Inner: left column unchanged. Run pane = header (session, account, time), Chat|Log toggle, Pause and red Stop as header buttons, chat, reply box. The hub Ask box leaves Desk (lives on Home).
- Out:
  - Live terminal attach. Editing past turns.

# Acceptance criteria
- AC1: A running mission's new turns show in Nest within 5 s.
- AC2: Text typed in Nest lands in that mission or session and its reply shows.
- AC3: A long conversation loads the last 50 turns first and pages back.
- AC4: On the cover screen a 1400-char reply with bullets reads whole, with its line breaks; tool calls show as one folded line per run.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: read and write inside a run from the phone.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (req_027, item_148)

# Links
- Request: `req_055_a_work_day_from_nest_without_a_terminal`, `req_027_conversations_feed_for_core_deck_nest`
- Depends on: `item_148_run_conversation_poller_with_cursors_in_core`, `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- Related: `item_186_deck_chat_pane_in_the_terminal` (same message shape, terminal side), `item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: High
- Rationale: steering missions needs their words, not only their state.
- 2026-11-11 build, uncommitted: core GET sessions/{pid}/conversation and missions/{id}/conversation through crow-runs conversation; Nest Chat cover (/desk/<pid>, /desk/m/<id>) and the inner Desk right pane with Chat/Log toggle and Pause/Stop in the header. Owner OK on both screens vs boards (vault Research/nest-design/item_177/chat-cover-compare, chat-inner-compare). Earlier keys on the before cursor: cdx has_more on a tail page speaks of later turns. Live probe OK. Gaps: no waits-since time on the permission ask, 3 s polling instead of the item_148 store, fixtures share one conversation. Next: owner reviews the diff, then commit.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
