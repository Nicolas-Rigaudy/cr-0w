## item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide - Core action endpoints for dispatch answer stop cancel and hide
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 16:17:33

# AI Context
- Summary: core action endpoints: dispatch a mission, answer an inbox item (verb, text, attachments, mission version), stop, cancel and hide. Same rules as the `crow mission` and `crow inbox` CLIs, which call shared Go code.
- Keywords: core api, post, dispatch, answer, cancel, stop, hide, csrf, stale version
- Use when: adding or changing a state-changing core endpoint.
- Skip when: publish preview and confirm (req_030 publish path), MCP tools (slice 8).

# Problem
- Slice 3 needs actions second. The phone must answer a question without a terminal (req_028 AC1).

# Scope
- In:
  - Routes under `/api` (item_144 amendment): `POST /missions` (dispatch: kind, project, goal, end state, autonomy override), `POST /missions/<id>/answer`, `/stop`, `/cancel`, `POST /inbox/<id>/hide`.
  - Answer refused with 409 when the version is stale (req_031 push section). Verbs validated against the question's choices.
  - CSRF header and owner identity from item_144; every action writes its event with `actor: owner` and the surface name.
  - CLI and API call one Go function per action; the CLI keeps working without the service.
  - Request and reply fixtures.
- Out:
  - Publish, merge, posts (adr_001 publish section, req_030).
  - Resume behaviour itself (item_135): an answer queues the resume; the tick does it.

# Acceptance criteria
- AC1: `curl -X POST` an answer from the phone; the mission leaves `waiting` and its answers file holds the text.
- AC2: An answer for an older version gets 409 and changes nothing.
- AC3: A POST without the CSRF header gets 403.
- AC4: `crow inbox answer` and the endpoint produce the same events (one test over both).

# AC Traceability
- request-AC1 -> AC1. Proof: answer from the phone over the tailnet.
- request-AC3 -> AC4. Proof: no publish action exists on the API.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (adr_001 API section)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Depends on: `item_144_core_http_service_with_read_endpoints_on_the_tailnet`. Unblocks `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`. AC1 end to end with a resumed run also needs `item_135_resume_a_parked_mission_from_the_owner_s_answer` (Done); AC1 itself only needs the mission to leave `waiting`.
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api` (moved from task_017 on 2026-11-10)

# Priority
- Priority: High
- Rationale: Second half of slice 3.

# Notes
- Groomed 2026-11-09 from req_028 slice 3.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
