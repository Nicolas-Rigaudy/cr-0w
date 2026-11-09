## req_028_jarvis_layer_missions_inbox_butler_chat_model_router - Jarvis layer: background missions, one inbox, butler chat, model router
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85
> Confidence: 75
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 11:24:42

# AI Context
- Summary: Turn crow into a hands-off assistant. Work runs as background missions through cdx. Decisions, approvals and proposals land in one inbox shown in deck, nest and comlink. A butler chat talks to everything going on. A router picks the model and account per mission kind and data class. Status is told in short human sentences; detail stays one step away.
- Keywords: jarvis, mission, inbox, butler, chat, park and resume, autonomy, summaries, router, gemini, gemma, ollama, antigravity, nest, deck
- Use when: planning how crow runs work in the background and how the owner steers it from deck, nest or the phone.
- Skip when: building a single deck widget or a single ritual that does not touch missions or the inbox.

# Needs
1. **Less terminal.** The owner steers from one interface (deck or nest, laptop or phone), not from terminal sessions.
2. **Background missions.** A mission is the unit the owner thinks in: goal, kind (design, code, review, research, mail, admin), scope (work, personal, hub), budget, end state, allowlist. core splits it into cdx runs.
3. **One inbox.** Every question, approval and proposal from any mission lands in one queue, shown by every surface. Item types: choice (options with previews), approval (exact draft: push, MR, mail), review (diff, HTML pages, mockups), FYI.
4. **Proposals shown, not described.** A design mission comes back with pages; nest shows them side by side (crow-compare in the browser) and the pick resumes the mission.
5. **Butler chat.** One persistent chat that sees all missions, runs and the inbox, dispatches new missions, answers questions and summarises. It does no deep work itself.
6. **Hands off by default, detail on demand (owner, 2026-11-06).** Per mission kind, an autonomy level: full auto, notify, ask. Some flows go full auto (example: review fixes matching a known ledger pattern, up to a local commit); the MR pane stays the drill-down when the owner wants detail. The agents panel shows short human-language status per run, summarised from `cdx conversation`; the raw trail sits one key away.
7. **Model and account router.** Pick provider, model and account per mission kind, scope and data class, using every source the owner has: two Claude accounts, a capped shared fallback account, Google free tiers, local models.

# Decisions (owner, 2026-11-06)
- **Both surfaces on one core API.** Build the core API and the inbox first, with thin views in deck and nest at the same time; no surface owns logic.
- **Park and resume.** A run that needs a decision writes an inbox item and exits; core parks the mission; the answer resumes it with `cdx run --continue-run`. No tokens or sessions held while waiting.
- **Free and local models for cheap tasks and private data.** Gemini Flash free API key or local Gemma for triage, summaries, classification and private personal data (bank, health). Antigravity free tier as optional overflow for side-project code. Claude for everything that matters. Work data never leaves Claude accounts.
- Plan here, build later in fresh sessions from backlog docs.
- Review round (2026-11-06, `local/state/research/jarvis-roadmap-review-2026-11-06.md`): nest is the main surface, deck gets an inbox pane; Web Push first, ntfy only if needed; a Discord bot for butler free-text chat on personal data only; 8 mission kinds; 7 nest screens plus 3 overlays; the butler track waits for core; AFK runs become mission batches in nest; a news board (req_039). The phone is a Galaxy Z Fold 6 (req_038).
- Product brief round (2026-11-06): see `prod_005_crow_jarvis_layer` for every decision (gates per kind, preflight plus park and resume, laptop first, Web Push with digest and quiet hours, one scope-aware butler with focus mode, auto-merge on green CI for personal repos, scripts-first jobs, repair catalogue, failover, watchdog, incident log). Architecture in `adr_001_core_missions_and_inbox_architecture`.

# Context
- Builds on: req_022 (core, deck, nest, comlink), req_027 (conversation feed), req_025 (scope routing), req_026 (personal butler), req_016 (multi-machine, tailnet), req_015 (wishlist: voice, screensaver, reports), req_019/item_077 (deny-by-default runs), req_020 (terminal kit, decisions widget, crow-compare v2).
- cdx already has: `cdx run` with `--provider claude|codex|antigravity|ollama`, `--kind`, `--continue-run`, `--failover`, `--detach`; `cdx select`; `cdx conversation`; per-account `--policy`, `--priority`, `--budget`, `--max-parallel`. Remote Ollama is text only (no repo edits).
- Hard rules stay above any autonomy level: work push, MR creation, merge, mail and Teams always go through an approval item with the exact text.
- Hardware: the laptop (Iris Xe, 30 GB RAM, mostly used) runs only small models on CPU. Real local inference needs a GPU machine on the tailnet (tower, if it has one) serving Ollama.
- Research 2026-11-06 (secondary sources, verify before relying):
  - Gemini API free key: about 250 requests/day, Flash models; free tier data may be used for training.
  - Gemini CLI free personal tier (1,000/day): one October 2026 source says it stopped on 2026-06-18; others still list it.
  - Antigravity free tier: personal Gmail allowed, weekly agent quota, unpublished and capacity based, no BYOK.
  - Gemma in AI Studio: free quota unpublished, reportedly cut in July 2026; self-hosting avoids limits.
  - Several Gmail accounts to stack free quotas counts as quota circumvention under Google terms: plan on one Google account.

# Architecture sketch
```
 owner ── nest (chat, inbox, mission board, proposal viewer; laptop + phone PWA)
       ── deck (inbox pane in the terminal)        ── Web Push (ntfy fallback) · Discord bot (personal chat)
                          │  core API (tailnet only)
                        core: missions · inbox · scheduler · router · policy · summaries
                          │
           cdx run ── claude (work · personal · capped shared) · codex · antigravity · ollama (gemma)
                   └─ gemini API key (cheap text tasks)
```

# Slices (to groom into backlog items, in order; 1-2 groomed 2026-11-09 into task_016)
1. **Mission model.** `crow mission new|list|show|cancel`; mission record (goal, kind, scope, budget, end state, autonomy, state `queued → running → waiting → review → done|failed`); one mission maps to one or more cdx runs; kinds and defaults in tracked `kinds/<kind>.yaml` files with optional private overrides in `local/` (owner decision 2026-11-07, req_037 D1). Depends on item_077 (allowlist per kind) and req_025 (scope → account).
   - AC: a dispatched mission shows its state and runs in `crow mission list --json`.
2. **Inbox and park/resume.** File contract a run writes (question or proposal + artifacts dir); core turns it into an inbox item, parks the mission, resumes it with `--continue-run` on answer; `crow inbox list|show|answer`.
   - AC: a test mission asks a choice, exits, and resumes with the answer in its prompt.
3. **core API.** systemd user service, HTTP bound to the tailnet, read endpoints first (missions, inbox, runs, summaries), actions second (dispatch, answer, cancel). Fold in the req_022 core item.
   - AC: deck and a curl client read the same inbox.
4. **Human summaries.** Per run, a one or two line status in plain language from `cdx conversation` (cursor based, cached), by a cheap model (Gemini Flash key, local Gemma or Haiku). deck agents panel and nest show the summary first, the trail on demand.
   - AC: the agents panel lists each live run with a short sentence; a key opens the full trail.
5. **Autonomy dial.** Per kind: full auto, notify, ask. Hard rules override. First full-auto candidates: review fixes matching ledger patterns up to a local commit, retro batch, briefs, mail triage (drafts only).
   - AC: a full-auto mission finishes with only an FYI item; an ask mission stops at its first decision.
6. **deck views.** Inbox and mission board widgets on the req_020 kit; the MR pane stays the detail view reached from an item.
7. **nest MVP.** Web app over core: Home with inbox, Missions, Mission detail with the mockup tab (pages side by side, flip and zoom, pick resumes the mission); phone PWA laid out for both Fold 6 screens; Web Push with action buttons for approvals.
8. **Butler chat.** Persistent session with a crow tool surface (missions, inbox, runs, summaries, dispatch) instead of reading transcripts; context kept small; same chat reachable from deck and nest.
9. **Model router.** Rules by kind × scope × data class: work → Claude only; private personal data → local only; cheap text → Gemini key or Gemma; side-project code overflow → Antigravity. Uses `cdx select` and account policies; records cost per mission (item_075).
10. **AFK batches.** The New mission sheet in batch mode gathers jobs and asks every preflight question up front (from the `afk` skill and `crow-afk-preflight`); core runs the batch; the morning report becomes a batch digest in Missions. Retires `crow-afk-page` and the night pages.
   - AC: an AFK batch launched from the phone runs overnight and its digest opens from one push in the morning.
11. **Later.** Voice in and out, mission screensaver (req_015), multi-machine routing (req_016), report pipelines.

# Open questions (grooming)
- Butler: one scope-aware butler with a focus mode (decided); engine decided 2026-11-07 in req_032 (long-lived `claude -p` stream-json, Sonnet, routing by data scope).
- Mission store: SQLite proposed in adr_001. Backend decided 2026-11-07: Go core plus Python sidecars (req_031); the store follows the Go core.
- Hosts and local models: req_029 (decided 2026-11-07: laptop is the control host now, laptop and tower are both mission hosts).
- Gemini and Antigravity free usage: req_029 (Antigravity via cdx is the Google path; terms check first).

# Acceptance criteria
- AC1: The owner can dispatch a mission, leave, answer its question from the phone, and find the result as a review item, without opening a terminal.
- AC2: Every live run shows a short human-language status; full detail is one action away.
- AC3: Each mission kind has an autonomy level, and hard rules (work push, MR, merge, mail) still require an approval item.
- AC4: The router never sends work data to non-Claude providers and never sends private personal data off the machine.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_015_wishlist_ambient_mobile_voice_colleague_ideas.md`
- `logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md`
- `logics/request/req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control.md`
- `logics/request/req_025_account_routing_by_scope.md`
- `logics/request/req_026_personal_butler.md`
- `logics/request/req_027_conversations_feed_for_core_deck_nest.md`
- `logics/roadmap/road_001_crow_roadmap_after_the_deck_pass.md`
- Spikes: `logics/request/req_029_spike_hosts_model_sources_and_local_ai_benchmarks.md`, `logics/request/req_030_spike_mission_safety_and_sandbox.md`, `logics/request/req_031_spike_nest_stack_and_app_shape.md`, `logics/request/req_032_spike_butler_engine_and_persona.md`, `logics/request/req_033_spike_android_extras_for_nest.md`
- Tracks: `logics/request/req_034_work_views_in_nest.md`, `logics/request/req_035_forge_ops_in_nest.md`, `logics/request/req_036_fun_layer_for_crow.md`
- Design: `logics/request/req_037_mission_kinds_and_run_contract.md`, `logics/request/req_038_nest_design_rules.md`

# Backlog
- Slices 1-2 groomed 2026-11-09, run order in `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`:
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

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
