## task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy - Orchestrate Jarvis core slices 3 to 5: API, summaries and autonomy
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Build req_028 slices 3 (core API), 4 (human summaries) and 5 (autonomy dial) on the Go core in `tools/crow-core`, one fresh session per item, in the order below.
- Keywords: orchestrate, jarvis, core, http api, sse, summaries, conversation, autonomy
- Use when: picking the next core item after task_016.
- Skip when: working inside an item already started; slices 6 and later (deck views, nest MVP, butler, router, AFK batches).

# Context
- road_001 0.3b. task_016 steps 1-9 are done; step 10 (item_135 resume) waits on AlexAgo83/cdx-manager#131 and item_143. Nothing here needs resume except the end-to-end half of item_147 AC1, so this task can start now.
- Settled: adr_001 (API section, publish rules), req_031 decisions 1 and 2 (Go core, Python sidecars, SSE plus ETag refetch), req_037 run contract, req_027 conversation feed.
- slice 3 folds in the req_022 core item (systemd user service, `crow core status|restart`).
- Each item names its own mockup step where it changes deck screens.

# Plan
- [ ] 1. item_144 core HTTP service with read endpoints on the tailnet.
- [ ] 2. item_145 deck inbox pane reads the core API. Slice 3 read exit: deck and curl read the same inbox.
- [ ] 3. item_147 action endpoints. Slice 3 exit: an answer posted from the phone moves the mission out of `waiting`.
- [ ] 4. item_146 SSE section stream with polling fallback.
- [ ] 5. item_148 run conversation poller with cursors.
- [ ] 6. item_149 one-line run summaries (benchmark Haiku vs local first).
- [ ] 7. item_150 deck agents panel: summary first, trail on a key. Slice 4 exit.
- [ ] 8. item_151 autonomy dial under the hard rules. Slice 5 exit: an auto mission ends with only an FYI item; an ask mission stops at its first decision.
- [ ] 9. item_152 first full-auto kinds (owner OK per kind; needs item_143 for real read-only runs).
- [ ] Parallel: step 8 has no API dependency and can run beside steps 1-7. Steps 5-6 can start after step 1.
- [ ] ADR 009 checkpoint: update affected Logics docs during each meaningful wave and leave the repo commit-ready.
- [ ] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_144_core_http_service_with_read_endpoints_on_the_tailnet`
- `item_145_deck_inbox_pane_reads_the_core_api`
- `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`
- `item_146_section_change_stream_over_sse_with_a_polling_fallback`
- `item_148_run_conversation_poller_with_cursors_in_core`
- `item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor`
- `item_150_deck_agents_panel_shows_the_run_summary_first_and_the_trail_on_a_key`
- `item_151_autonomy_dial_per_kind_and_mission_under_the_hard_rules`
- `item_152_turn_on_full_auto_for_the_first_mission_kinds`

# Links
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Related: `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`, `req_022`, `req_027`, `item_137_nest_app_shell_on_core_fixtures`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Product brief(s): `prod_005_crow_jarvis_layer`

# Definition of Done (DoD)
- [ ] Every listed item is Done with its acceptance criteria proven.
- [ ] Slice 3, 4 and 5 acceptance criteria of req_028 pass end to end.
- [ ] Validation passes.

# AC Traceability
- request-AC1 -> `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`. Proof deferred to slice closeout.
- request-AC2 -> `item_150_deck_agents_panel_shows_the_run_summary_first_and_the_trail_on_a_key`. Proof deferred to slice closeout.
- request-AC3 -> `item_151_autonomy_dial_per_kind_and_mission_under_the_hard_rules`. Proof deferred to slice closeout.
- request-AC4 -> `item_149_one_line_run_summaries_from_a_cheap_model_cached_per_cursor`. Proof deferred to slice closeout.

# Validation
- (no validation recorded yet)

# Report
- 2026-11-09 groomed from req_028 slices 3-5. Next: step 1, item_144, in a fresh session.
