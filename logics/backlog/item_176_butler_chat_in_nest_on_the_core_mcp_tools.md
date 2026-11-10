## item_176_butler_chat_in_nest_on_the_core_mcp_tools - Butler chat in Nest on the core MCP tools
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 70%
> Progress: 0%
> Complexity: High
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: a chat in Nest with the butler (CR-0W persona), run by the req_032 engine: one long-lived `claude -p` stream-json process per conversation, core MCP tools only, a core bridge that streams replies to Nest. The butler starts missions first, and sessions or dispatch only when asked.
- Keywords: butler, chat, orchestrator, claude -p, stream-json, mcp, core bridge, nest
- Use when: building the butler conversation or its core bridge.
- Skip when: chat inside one run (item_177), persona text (item_118).

# Problem
- Every request to crow today starts with a terminal hub session. The owner wants to talk to the orchestrator from Nest (2026-11-10). The engine is decided (req_032, owner 2026-11-07) but nothing runs it.

# Scope
- In:
  - Core bridge: start, keep and stop one butler process per conversation; `POST /api/butler/<conv>/say`, replies streamed over SSE (item_146) with a polling fallback; transcript kept in core state; one process idles out after 30 min and resumes from the transcript.
  - Flags per req_032: `--tools "" --strict-mcp-config --allowedTools "mcp__core__*" --permission-mode dontAsk`, on an account allowed by the cdx policy.
  - Core MCP server exposes the item_169 reads and the item_147/170 actions. Missions first: `mission.new` with a kind, then session and dispatch.
  - Any action that needs a tap (confirmation, approval) becomes an inbox item. The butler never confirms for the owner.
  - Persona from item_118 ("Sir", formal and calm).
  - Nest: a Chat tab with the conversation, a text box, action chips the butler proposes (tap = the same core action).
- Out:
  - Voice. Several parallel butler conversations in the UI (one is enough at first). Butler editing files directly.

# Acceptance criteria
- AC1: "Start a review of api!412" in Nest creates a review mission; the reply links it.
- AC2: "What's running?" answers from the core sections, matching Home.
- AC3: An action that needs approval lands in the inbox; the butler says so and nothing runs before the tap.
- AC4: Closing Nest and coming back an hour later shows the same conversation; the next message continues it.
- AC5: The butler process has no tool outside `mcp__core__*` (checked from its init event).

# AC Traceability
- request-AC1 -> AC1, AC2, AC3, AC5. Proof: talk to the orchestrator from the phone under the same rules.

# Decision framing
- Product framing: Not needed (prod_005, req_032)
- Architecture framing: Not needed (req_032 decided the engine)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Request: `req_055_a_work_day_from_nest_without_a_terminal`
- Depends on: `item_169_core_serves_every_deck_source_as_read_only_sections`, `item_170_core_actions_on_live_sessions_and_dispatch_runs`, `item_146_section_change_stream_over_sse_with_a_polling_fallback`
- Related: `item_118_butler_persona_file_steward_with_a_droid_flavour`, `req_032_spike_butler_engine_and_persona`
- Primary task(s): `task_021_orchestrate_a_work_day_from_nest`

# Priority
- Priority: High
- Rationale: the main surface the owner asked for.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
