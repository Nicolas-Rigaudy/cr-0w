## task_020_orchestrate_nest_going_live_on_the_core_api - Orchestrate Nest going live on the core API
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 57%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-11-10 21:33:49
> Owner: crow

# AI Context
- Summary: Fastest path to a working Nest: core read API serving Nest on the tailnet, Nest on live data, core actions, a Nest answer sheet. Core items alternate with the Nest Home design items (owner decision 2026-11-10). Ends on the roadmap 0.3b exit signal answered from Nest.
- Keywords: orchestrate, nest, live, core api, tailnet, pwa, answer sheet, design pass, 0.3b
- Use when: picking the next Nest or core item.
- Skip when: working inside an item already started; slices 4-5 (summaries, autonomy) and the deck inbox pane, which stay in task_017.

# Context
- Owner ask 2026-11-10: get the web app up and working the fastest. Nest was a finished shell on fixtures; item_135 resume is Done and cdx-manager#131 is closed, so the whole park/answer/resume loop can now run.
- Takes items 144 and 147 out of task_017 (steps 1 and 3), adds item_167 (Nest API adapter, never scoped before) and item_168 (Nest answer sheet, pulled from req_028 slice 7 into 0.3b; roadmap exception).
- Interleaves the task_019 Home design items 158, 163, 165, 166 (owner decision). Each design item still needs owner OK on renders before the next design item.
- Hosting: core serves Nest and the API on one origin behind `tailscale serve` (item_144 amendment; req_031 sends no CORS headers).
- One fresh session per item, from its backlog doc. Owner prefs go into every handoff.

# Plan
- [x] 1. Nest on the phone on fixtures, here in the hub: `cd tools/nest && npm run build && npm run preview -- --port 4173`, `tailscale serve --bg 4173`, install from Chrome on the Fold. Closes task_018 AC1. Replaced by the core unit at step 2. Done 2026-11-10 (runbooks/nest-phone-install.md).
- [x] Nit (owner 2026-11-10): installed app name reads `Nest`; set manifest `name` and `short_name` to `CR-0W`. Done with step 2 (76534a4); reinstall the app on the phone to pick up the name.
- [x] 2. item_144 core HTTP service with read endpoints, serving Nest at `/` and the API at `/api`, systemd user unit, `tailscale serve` pointed at core. Exit: the phone loads Nest from core and `curl` reads `/api/state/inbox`. Done 2026-11-10: unit active with linger, tailscale serve proxies `/` to 127.0.0.1:8787, tailnet curl reads `/api/state/inbox`, vite preview stopped (runbooks/core-service.md).
- [x] 3. item_167 Nest reads live data. Exit: the phone shows the same missions and inbox as the CLI. Done 2026-11-10, moved before item_158 (owner: core work before design polish).
- [x] 4. item_173 script JSON contracts (fixtures, usage summary, MR cache, JSON from action scripts). Exit: deck still works; every output pinned by a test.
- [x] 5. item_169 core serves every deck source (sessions, dispatch runs, MRs, today, usage, checks, projects) with background refresh. Exit: the phone's Nest data layer reads the same sessions and runs as deck (req_054, owner 2026-11-10: Nest shows zeros, it must see everything deck does). Done 2026-11-10.
- [x] 6. item_147 core mission action endpoints, then item_170 session and run actions (start, resume, send, interrupt, stop, dispatch) on the same guard. Exit: a `curl` answer moves a mission out of `waiting`; text sent from the phone lands in a live session. item_147 done 2026-11-10 (4fbef04); item_170 done 2026-11-10 (99c977d core, 340aa19 Nest Desk).
- [x] 7. item_174 permission prompts answered from Nest (spike on the prompt source first). Done 2026-11-10: PermissionRequest hook source (bfa3141, d248562), core answer route (19800ca), Nest layout B (89fb516), doctor check (416f8c6); live tmux test passed allow, deny and stale.
- [x] 8. item_171 Nest Home on the deck sources (mockup first, owner OK). Exit: Home matches deck's agents panel with an empty mission store. Done 2026-11-11 (06545b2); bundle budget follow-up item_185.
- [ ] 9. item_158 hero interactions (spin, hover lift + moons, click to details). Open question for the owner: hero mood order, model.ts `moodOf` (run > fail > idle, item_162) vs item_158 spec (fail > ask > run > idle).
- [ ] 10. item_168 Nest answer sheet (mockup first, owner OK). Exit: roadmap 0.3b signal, answered from Nest on the phone.
- [ ] 11. item_175 work approvals in Nest on the exact text (inbox kind `approval`).
- [ ] 12. item_165 telemetry strip on the live usage and accounts sections.
- [ ] 13. item_172 Nest Work and System screens (mockup first, owner OK).
- [ ] 14. item_163 design skills audit of the laptop Home, now on real data.
- [ ] 15. item_166 boot sequence.
- [ ] Order rule: core steps (2 to 7, 11) are the critical path. If a design step waits on owner OK, start the next core or build step instead of idling.
- [ ] After step 15: task_017 continues (item_145 deck inbox, 146 SSE, 148-152 summaries and autonomy), alternating with task_018 laptop Missions. req_055 continues in task_021.
- [ ] ADR 009 checkpoint: update affected Logics docs during each meaningful wave and leave the repo commit-ready.
- [ ] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_144_core_http_service_with_read_endpoints_on_the_tailnet`
- `item_158_nest_hero_live_data_state_motions_and_actions`
- `item_167_nest_reads_live_data_from_the_core_api`
- `item_163_design_skills_audit_of_the_laptop_home`
- `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`
- `item_165_telemetry_strip_on_the_laptop_home`
- `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`
- `item_166_boot_sequence_on_first_load_of_the_nest`
- `item_169_core_serves_every_deck_source_as_read_only_sections`
- `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- `item_171_nest_home_shows_live_sessions_runs_alerts_and_quota`
- `item_172_nest_work_and_system_screens_on_the_deck_sources`
- `item_173_deck_scripts_get_a_stable_json_contract_for_the_app`
- `item_174_nest_answers_a_session_s_permission_prompt_with_buttons`
- `item_175_work_approvals_given_from_nest_on_the_exact_text`

# Definition of Done (DoD)
- [ ] Every listed item is Done with its acceptance criteria proven.
- [ ] A test mission asks a question, exits, and resumes from an answer given in Nest on the phone (road_001 0.3b exit signal).
- [ ] task_018 AC1 (PWA installed on the phone) proven.
- [ ] Validation passes.

# AC Traceability
- request-AC1 -> `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`. Proof deferred to closeout: answer from Nest resumes the mission.
- request-AC3 -> `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`. Proof deferred to closeout: no publish action on the API or in Nest.
- request-AC4 -> This task. Proof: item_173: implemented in 5b4a90d and 5cbdd7d; validated with python bin/crow-check (pytest, ruff, leak), tests/test_script_contract.py, go test in tools/crow-deck and tools/crow-core, deck tabs 1-5 rendered headless. Source: `5cbdd7d`
- request-AC2 -> This task. Proof: Implemented in 99c977d (core routes) and 340aa19 (Nest Desk); validated with go test ./... (168 passed), vitest (24 passed), a 403 on POST without X-Crow, and a phone test 2026-11-10 logged in actions.jsonl (start + send, actor owner, surface nest). Source: `340aa19`

# Validation
- 2026-11-10 task_018 AC1 proven: WebAPK installed on the Fold from the tailnet URL, opens standalone with Tailscale on.

# Report
- 2026-11-11 step 8 built and committed (06545b2), owner OK on Home vs D boards. Bundle 95.9 kB gzip, follow-up item_185. Next: step 9 item_158 (hero mood order question open) or item_176 butler chat on Home.
- 2026-11-10 step 7 done (item_174). Next: step 8 item_171 Home mockup, owner OK first.
- 2026-11-10 step 8 design approved (pass 3b option D, vault Research/nest-design/item_171). Butler chat moved to Home (item_176). Next: build item_171 in a fresh session.
- 2026-11-10 planned from the owner ask (fastest working web app). Next: step 1 in the hub, then step 2 (item_144) in a fresh session.

# Links
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`, `req_054_nest_sees_and_pilots_everything_deck_does`
- Related: `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`, `task_018_nest_app_shell_on_core_fixtures`, `task_019_orchestrate_the_nest_home_design_pass`, `road_001_crow_roadmap_after_the_deck_pass`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Product brief(s): `prod_005_crow_jarvis_layer`, `prod_007_nest_home_design_pass`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
