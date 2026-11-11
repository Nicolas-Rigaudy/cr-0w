## item_168_nest_answer_sheet_for_inbox_items_and_mission_actions - Nest answer sheet for inbox items and mission actions
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 95%
> Confidence: 90%
> Progress: 90%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 16:17:25

# AI Context
- Summary: Nest answers an inbox item (its verbs, optional text, at the item's version) and stops, cancels or hides from Inbox and Mission detail, through the item_147 endpoints. Pulled forward from req_028 slice 7 so roadmap 0.3b can exit from Nest.
- Keywords: nest, answer sheet, inbox, verbs, stop, cancel, hide, csrf, 409 stale, optimistic
- Use when: building or changing a state-changing control in Nest.
- Skip when: dispatching a new mission from Nest (New mission sheet, slice 7/10), mockup tab picks, Web Push action buttons, publish.

# Problem
- Roadmap 0.3b exits when "a mission asks a question, exits, and resumes from an answer given in deck or nest". Answering from deck is slice 6 and from Nest slice 7, both in milestone J, so 0.3b had no surface to answer from except curl.
- Owner decision 2026-11-10: pull a small Nest answer sheet into 0.3b. The rest of slice 7 stays in J.

# Scope
- In:
  - Mockup first: an answer sheet on Inbox and Mission detail for the cover (bottom sheet) and inner/laptop (side panel), both themes, from the approved boards and tokens. Owner OK before building (owner prefs: show options before building a UI).
  - Verbs from the item (`verbs`), validated by core; optional text; the item `version` sent with the answer. `fyi` items get hide only.
  - Stop and cancel on Mission detail with a confirm step; hide on inbox rows.
  - POSTs carry the CSRF header (`X-Crow: 1`) and go same-origin to `/api` (item_144, item_147).
  - After a POST: refetch the touched sections (no optimistic state). 409 stale: show the newer version and keep the typed text. 403/5xx: inline error, text kept.
  - `DataSource` grows an `actions` part (answer, stop, cancel, hide); the fixture adapter implements it in memory so screens and tests run without core.
  - Vitest for the action calls and the 409 path; a Playwright shot of the sheet per board size for review (runbooks/nest-screenshots.md).
- Out:
  - Attachments in answers, New mission sheet, mockup tab flip/zoom/pick, Web Push, publish or merge actions (adr_001 publish section).

# Acceptance criteria
- AC1: The owner approves the sheet mockup before build, and the built sheet side by side with it.
- AC2: From the phone, answering a test mission's `choice` item moves the mission out of `waiting`; the tick resumes it (item_135) and Nest shows it running without a reload. This is the roadmap 0.3b exit signal.
- AC3: Answering an item whose version is stale shows the newer question and keeps the typed text; nothing changes in core.
- AC4: Stop, cancel and hide each produce the same events as the matching `crow mission` / `crow inbox` command (actor owner, surface nest).
- AC5: No publish, merge or post control exists in Nest.

# AC Traceability
- request-AC1 -> AC2. Proof: answer from Nest on the phone, mission resumes.
- request-AC3 -> AC5. Proof: no publish action in Nest or the API.

# Decision framing
- Product framing: Not needed (prod_005)
- Architecture framing: Not needed (adr_001 API and publish sections, req_031 auth)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Depends on: `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`, `item_167_nest_reads_live_data_from_the_core_api`, `item_135_resume_a_parked_mission_from_the_owner_s_answer` (Done)
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: Closes the roadmap 0.3b exit signal from Nest (owner decision 2026-11-10).

# Notes
- Groomed 2026-11-10. Roadmap exception recorded in road_001 0.3b.
- 2026-11-11 build (commit d6de8cc): shape A answers inline in the Inbox reading pane, rebuilt from the fold-inbox boards; cover verbs docked above the tab bar. Owner decision: no Mission detail screen (no board draws one). Stop is the Desk S1 STOP; Cancel is the ⋯ on Desk queue rows; Home mission links go to Desk or Inbox. Compare pages in the vault: nest-design/item_168/{inbox,desk-queue}-compare.
- AC status: AC1 done (owner OK on both compare pages). AC3 covered by the 409 path and tests. AC5 holds. AC2 (phone answer resumes a mission) and AC4 (event parity checked against core) still to verify live.
- 2026-11-11: Desk STOP on a run that belongs to a running mission now calls `missions/{id}/stop` (`missionOfRun` in data/state.ts), so core writes the owner `action` event with surface nest like `crow mission stop`. Ad-hoc runs keep the dispatch stop. AC4 code path aligned; confirm the event in the AC2 phone test.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
