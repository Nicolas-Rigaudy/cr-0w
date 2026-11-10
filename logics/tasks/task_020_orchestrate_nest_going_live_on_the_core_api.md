## task_020_orchestrate_nest_going_live_on_the_core_api - Orchestrate Nest going live on the core API
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-11-10 17:03:11

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
- [ ] Nit (owner 2026-11-10): installed app name reads `Nest`; set manifest `name` and `short_name` to `CR-0W`. Do it with step 2 when core starts serving Nest.
- [ ] 2. item_144 core HTTP service with read endpoints, serving Nest at `/` and the API at `/api`, systemd user unit, `tailscale serve` pointed at core. Exit: the phone loads Nest from core and `curl` reads `/api/state/inbox`.
- [ ] 3. item_158 hero interactions (spin, hover lift + moons, click to details). Design, fixtures are enough.
- [ ] 4. item_167 Nest reads live data. Exit: the phone shows the same missions and inbox as the CLI.
- [ ] 5. item_163 design skills audit of the laptop Home, now on real data.
- [ ] 6. item_147 core action endpoints. Exit: a `curl` answer moves a mission out of `waiting`.
- [ ] 7. item_165 telemetry strip on live accounts and runs.
- [ ] 8. item_168 Nest answer sheet (mockup first, owner OK). Exit: roadmap 0.3b signal, answered from Nest on the phone.
- [ ] 9. item_166 boot sequence.
- [ ] Order rule: core steps (2, 4, 6, 8) are the critical path. If a design step waits on owner OK, start the next core step instead of idling.
- [ ] After step 9: task_017 continues (item_145 deck inbox, 146 SSE, 148-152 summaries and autonomy), alternating with task_018 laptop Missions.
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

# Definition of Done (DoD)
- [ ] Every listed item is Done with its acceptance criteria proven.
- [ ] A test mission asks a question, exits, and resumes from an answer given in Nest on the phone (road_001 0.3b exit signal).
- [ ] task_018 AC1 (PWA installed on the phone) proven.
- [ ] Validation passes.

# AC Traceability
- request-AC1 -> `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`. Proof deferred to closeout: answer from Nest resumes the mission.
- request-AC3 -> `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`. Proof deferred to closeout: no publish action on the API or in Nest.

# Validation
- 2026-11-10 task_018 AC1 proven: WebAPK installed on the Fold from the tailnet URL, opens standalone with Tailscale on.

# Report
- 2026-11-10 planned from the owner ask (fastest working web app). Next: step 1 in the hub, then step 2 (item_144) in a fresh session.

# Links
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_053_nest_home_design_pass_netrunner_palette_holographic_hero_and_jarvis_chrome`
- Related: `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`, `task_018_nest_app_shell_on_core_fixtures`, `task_019_orchestrate_the_nest_home_design_pass`, `road_001_crow_roadmap_after_the_deck_pass`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Product brief(s): `prod_005_crow_jarvis_layer`, `prod_007_nest_home_design_pass`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
