## item_148_run_conversation_poller_with_cursors_in_core - Run conversation poller with cursors in core
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: core keeps one `cdx conversation --json` cursor per live mission run, stores new entries as a cached trail, and serves `GET /runs/<id>/conversation?tail|cursor|before`. Source for summaries (item_149) and the trail view (item_150).
- Keywords: cdx conversation, cursor, next_cursor, before_cursor, reset_required, trail, poller, back-off
- Use when: reading what a mission run said, from core or a surface.
- Skip when: interactive sessions outside missions (rest of req_027), sending prompts.

# Problem
- Slice 4 summaries are "cursor based, cached" from `cdx conversation`. req_027 makes it the one read path; no crow code parses Claude .jsonl for that content.
- Mission runs are headless; their trail is the only live view (req_031 terminal section).

# Scope
- In:
  - A sidecar step run by the tick timer, not inside the tick (adr_001: the tick calls no model; the poller calls no model either, it may run in the tick if it stays fast). Resolve by exact conversation id from the run record or `cdx sessions --json`; flag recency fallback.
  - Cursor per run in the store; bounded `--limit`; back-off for idle runs; `reset_required` restarts from tail.
  - Entries kept per run (role, text, tool_name, tool_summary, at, id) with a size cap.
  - `GET /runs/<id>/conversation` with tail, cursor, before; fixture.
  - A `conversation` SSE event when new entries land (item_146 stream).
- Out:
  - Interactive session conversations and the req_027 deck/nest panes beyond missions; waiting-on-owner detection (req_027 AC5).

# Acceptance criteria
- AC1: For a live mission run, the endpoint returns the same entries as `cdx conversation <account> --conversation <id> --json`.
- AC2: A second poll with no new entries calls cdx with the stored cursor and adds nothing.
- AC3: A test covers `reset_required` and the recency flag.

# AC Traceability
- request-AC2 -> AC1. Proof: the trail behind each summary.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (req_027, adr_001)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_027_conversations_feed_for_core_deck_nest`
- Depends on: `item_133_dispatch_queued_missions_as_cdx_runs_and_track_them_in_a_tick` (Done), `item_144_core_http_service_with_read_endpoints_on_the_tailnet` for the endpoint only
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: Medium
- Rationale: First step of slice 4.

# Notes
- Check `cdx conversation` on cdx 0.26.11 against a headless `cdx run` before building: req_027 was written for interactive sessions.
- Groomed 2026-11-09 from req_028 slice 4.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
