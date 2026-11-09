## item_150_deck_agents_panel_shows_the_run_summary_first_and_the_trail_on_a_key - Deck agents panel shows the run summary first and the trail on a key
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
- Summary: The deck agents panel lists each live mission run with its summary sentence first; a key opens the full trail from `/runs/<id>/conversation`, scrolling back with before_cursor.
- Keywords: crow-deck, agents panel, summary, trail, conversation, before_cursor
- Use when: changing how deck shows live runs.
- Skip when: nest rendering (item_137 and slice 7).

# Problem
- Slice 4 AC: the agents panel lists each live run with a short sentence; a key opens the full trail. Today the panel shows lanes from usage and session data, not mission runs.

# Scope
- In:
  - Mission runs from `/state/runs` in the agents panel next to the existing lanes, summary first with its age.
  - Trail view on a key: entries as a timeline (text, one-line tool summaries), older pages on scroll, live tail through the stream event.
  - Fixtures and shot tests for both views.
- Out:
  - nest views, interactive session conversations (req_027 AC3).

# Acceptance criteria
- AC1: Each live mission run shows one summary sentence in the agents panel.
- AC2: The trail key opens the run's conversation; scrolling up loads older entries.
- AC3: Shot tests cover summary, no summary yet, and core down.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: slice 4 exit capture.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_027_conversations_feed_for_core_deck_nest`
- Depends on: `item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor`, `item_145_deck_inbox_pane_reads_the_core_api` (client)
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: Medium
- Rationale: Closes slice 4.

# Notes
- Show a mockup of the panel change before building (memory design-choices-ask-first).
- Groomed 2026-11-09 from req_028 slice 4.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
