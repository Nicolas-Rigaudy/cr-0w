## item_145_deck_inbox_pane_reads_the_core_api - Deck inbox pane reads the core API
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
- Summary: crow deck gets a read-only inbox list fed by core's `/state/inbox`, so deck and a curl client read the same inbox (slice 3 exit).
- Keywords: crow-deck, inbox, core api, http client, etag, read only
- Use when: wiring deck to core's HTTP API.
- Skip when: inbox widgets on the req_020 kit (slice 6), actions from deck.

# Problem
- deck reads bin/ scripts and state files. adr_001 makes surfaces thin clients over core. The slice 3 AC needs deck reading the inbox through the API.

# Scope
- In:
  - A small core client in `tools/crow-deck` (base URL from config, loopback default, ETag revalidation).
  - An inbox list (mission, question, age, verbs) on an existing screen, read only, opening the mission detail text on enter.
  - "core down" state shown in place, with the section age; no fallback to state files.
  - Mock data for the deck shot tests from the item_144 fixtures.
- Out:
  - Answering from deck (after item_147, slice 6), widget kit styling (slice 6), conversation views.

# Acceptance criteria
- AC1: With core running, deck and `curl /state/inbox` list the same items in the same order.
- AC2: With core stopped, deck shows "core down" and the last age, and does not crash.
- AC3: Deck shot tests render the inbox from fixtures.

# AC Traceability
- request-AC1 -> AC1. Proof: slice 3 exit, side by side capture of deck and curl.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (adr_001: surfaces hold no logic)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Depends on: `item_144_core_http_service_with_read_endpoints_on_the_tailnet`
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: High
- Rationale: Closes the slice 3 AC.

# Notes
- Groomed 2026-11-09 from req_028 slice 3. Placement on a deck screen follows the deck UX2 boards (memory deck-ux2-final); show a mockup before building (memory design-choices-ask-first).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
