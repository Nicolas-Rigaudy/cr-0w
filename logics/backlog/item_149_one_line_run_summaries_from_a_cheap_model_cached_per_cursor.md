## item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor - One line run summaries from a cheap model cached per cursor
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
- Summary: A summary sidecar turns new trail entries of each live run into one or two plain sentences, cached per cursor, routed by data class: work runs on Claude Haiku or a local model, personal non-private on the cheapest allowed model, private on local only.
- Keywords: summary, sidecar, haiku, gemma, ollama, gemini flash, data class, cache, cursor, cost
- Use when: changing how run summaries are made or routed.
- Skip when: the model router in general (slice 9).

# Problem
- req_028 AC2: every live run shows a short human-language status. Trails are long; surfaces must not call models.
- prod_005 routing: run summaries stay on Claude or local for work; never work data to non-Claude providers (req_028 AC4).

# Scope
- In:
  - Python sidecar per the item_119 contract (`summarise` request: run, data class, entries since last summary, previous summary; reply: text, model, cost).
  - Started by the timer after the poller; only runs whose cursor moved since the last summary; max N per pass; skipped when nothing changed (no model call).
  - Model per data class from config. Decide in the item, with one benchmark on 10 real trails: Haiku through `cdx run` oneshot vs local Ollama (Gemma) on laptop. Record cost and latency in the notes; pick the cheapest that reads well.
  - `summary` event per run (text, cursor, model, cost); runs section of `/state` carries the latest summary and its age.
  - Fallback text without a model: last tool_summary plus elapsed time.
- Out:
  - Gemini or Antigravity wiring (router slice 9); nest rendering.

# Acceptance criteria
- AC1: A live run gets a summary of at most two lines within one timer pass after new entries.
- AC2: No new entries, no model call (test with a fake model).
- AC3: A work run is never sent to a non-Claude remote provider (test on the routing table).
- AC4: Benchmark numbers and the chosen model are in the notes; cost per day estimated at 2 live runs.

# AC Traceability
- request-AC2 -> AC1. Proof: runs section with summaries.
- request-AC4 -> AC3. Proof: routing test.

# Decision framing
- Product framing: Not needed (prod_005 models table)
- Architecture framing: Not needed (req_031 decision 1: Python sidecars for AI pieces)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Depends on: `item_148_run_conversation_poller_with_cursors_in_core`, `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures` (Done)
- Related: `item_075` (cost per mission)
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: Medium
- Rationale: Core of slice 4.

# Notes
- A `cdx run` per summary pays session start-up; measure it in the benchmark before choosing Haiku.
- Groomed 2026-11-09 from req_028 slice 4.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
