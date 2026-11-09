## task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox - Orchestrate Jarvis core slices 1 and 2: missions and inbox
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 85%
> Confidence: 80%
> Progress: 20%
> Complexity: High
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-09 17:49:36

# AI Context
- Summary: Build req_028 slices 1 (mission model) and 2 (inbox and park/resume) as a Go core in `tools/crow-core` driven from the Python CLI, one fresh session per item, in the order below.
- Keywords: orchestrate, jarvis, core, mission, inbox, park, resume, go
- Use when: Picking the next 0.3b core item.
- Skip when: Working inside a single item already started; slices 3 and later (core API, summaries, autonomy, views).

# Context
- road_001 0.3b. 0.3 (req_019 safe background runs) closed 2026-11-09, so dispatch, allowlist and isolation exist.
- Contract and formats are settled in req_037 (spec `logics/specs/req_037_run_contract/`), architecture in adr_001, backend in req_031 decision 1 (Go core, Python CLI, Python sidecars).
- Done inputs: item_077 allowlist, item_095 scope account pick, item_100 isolation, item_104 kind validator, item_105 kind catalogue, item_106 crow block parser, item_115 usage-limit policy.
- No HTTP, no service, no UI here: `crow-core` is a one-shot binary; `crow-core tick` is run by hand or by a test until slice 3 wraps it in the service.

# Plan
- [x] 1. item_119 Go/Python contract and shared fixtures (also creates the `tools/crow-core` module for its Go tests).
- [x] 2. item_131 core skeleton: config, event log, SQLite store, rebuild.
- [x] 3. item_132 mission record, states, `crow mission new|list|show|cancel`.
- [x] 4. item_107 brief rendering (Go).
- [ ] 5. item_133 tick: admit, dispatch through `crow-runs`, finish runs. Slice 1 exit: a dispatched mission shows its state and runs in `crow mission list --json`.
- [ ] 6. item_109 outbox ingest and parking on `questions.json`.
- [ ] 7. item_108 versions, answers files, validation record.
- [ ] 8. item_134 inbox rows and `crow inbox list|show|answer`.
- [ ] 9. item_116 cold resume of isolated missions from files.
- [ ] 10. item_135 resume on answer. Slice 2 exit: a test mission asks a choice, exits, and resumes with the answer.
- [ ] Steps 6 and 7 can run in parallel sessions after step 5; everything else is sequential.
- [ ] ADR 009 checkpoint: update affected Logics docs during each meaningful wave and leave the repo commit-ready.
- [ ] Keep commit creation under operator control; do not force one commit per micro-step.
- [ ] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_119_go_core_and_python_sidecar_contract_with_shared_fixtures`
- `item_131_go_core_skeleton_with_the_mission_event_log_and_store`
- `item_132_mission_record_states_and_the_crow_mission_cli`
- `item_107_render_the_mission_brief_from_kind_record_and_answers`
- `item_133_dispatch_queued_missions_as_cdx_runs_and_track_them_in_a_tick`
- `item_109_outbox_ingest_with_caps_and_parked_questions`
- `item_108_mission_versions_responses_and_validation_record`
- `item_134_inbox_items_and_the_crow_inbox_cli`
- `item_116_resume_of_isolated_missions_from_files_and_cdx_issue_for_native_resume`
- `item_135_resume_a_parked_mission_from_the_owner_s_answer`

# Definition of Done (DoD)
- [ ] Every listed item is Done with its acceptance criteria proven.
- [ ] Slice 1 and slice 2 acceptance criteria of req_028 pass end to end.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# AC Traceability
- request-AC1 -> `item_133_dispatch_queued_missions_as_cdx_runs_and_track_them_in_a_tick`. Proof deferred to slice closeout.
- request-AC1 -> `item_135_resume_a_parked_mission_from_the_owner_s_answer`. Proof deferred to slice closeout.
- request-AC3 -> `item_132_mission_record_states_and_the_crow_mission_cli`. Proof deferred to slice closeout.
- request-AC3 -> `item_134_inbox_items_and_the_crow_inbox_cli`. Proof deferred to slice closeout.

# Validation
- (no validation recorded yet)

# Report
- 2026-11-09 step 1 (item_119) done in bba2e52. Next: step 2, item_131 Go core skeleton, in a fresh session. crow-check ruff red on pre-existing re.I aliases in bin/.
- 2026-11-09 step 2 (item_131) done: `tools/crow-core` event log, store, rebuild, `crow core`. Next: step 3, item_132 mission record and `crow mission`. Ruff is red on about 90 pre-existing findings in bin/ and tests/ (PLW1510, F811, SIM115...), tracked as a separate chore.
- 2026-11-09 step 3 (item_132) done: mission record, legal moves, `crow mission new|list|show|cancel` (a29413a, 3560404). Next: step 4, item_107 brief rendering.
- 2026-11-09 step 4 (item_107) done: brief rendering, prefs file, `crow core mission brief` (5113cc4). Next: step 5, item_133 tick. Open spec gap: build `by_scope.crow` gates are [plan] while the example log shows plan, report.

# Links
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Notes
- item_119: CI runs go test -count=1 because the Go test cache ignores spec files outside the module.
- item_131: crow-core append checks the event applies before writing the log; a bad event must never poison replay.
- item_131: an inbox row's created is the ts of the event that opens it (gate ready or question); fixture inbox.json aligned to 14:52:01.
