## item_146_section_change_stream_over_sse_with_a_polling_fallback - Section change stream over SSE with a polling fallback
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
- Summary: core streams "section X changed, version N" over SSE; clients refetch that section with its ETag, close the stream when hidden, and fall back to 5 s ETag polling after two stream failures (req_031 decision 2).
- Keywords: sse, server-sent events, etag, polling fallback, visibilitychange, live updates
- Use when: adding live updates to a core client or the stream itself.
- Skip when: Web Push (nest slice 7), live terminal (WebSocket, later).

# Problem
- deck and nest need fresh sections without polling cdx or the store themselves (adr_001, req_027 AC2).

# Scope
- In:
  - `GET /events/stream` (SSE): `section` events with name and version, a heartbeat, `Last-Event-ID` resume.
  - Store version bump hooks so a tick or an action announces the sections it touched.
  - deck client subscribes and refetches (item_145 client); a shared JS helper is left to nest.
  - Fixture for the event shape.
- Out:
  - Push to closed apps, conversation entries on the stream (item_148 adds its own event type).

# Acceptance criteria
- AC1: A tick that parks a mission makes deck refresh its inbox within 2 s without a manual reload.
- AC2: Killing the stream twice switches the client to 5 s polling; it returns to SSE on the next success.
- AC3: Three idle streams plus one poller keep core idle CPU within the req_031 probe numbers.

# AC Traceability
- request-AC1 -> AC1. Proof: parked question visible on deck without reload.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (req_031 decision 2)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_031_spike_nest_stack_and_app_shape`
- Depends on: `item_144_core_http_service_with_read_endpoints_on_the_tailnet`, `item_145_deck_inbox_pane_reads_the_core_api`
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: Medium
- Rationale: Polling with ETags works without it; needed before the agents panel goes live.

# Notes
- Groomed 2026-11-09 from req_028 slice 3.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
