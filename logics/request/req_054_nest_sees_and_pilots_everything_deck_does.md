## req_054_nest_sees_and_pilots_everything_deck_does - Nest sees and pilots everything deck does
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: core serves the data deck reads from `bin/` scripts (live sessions, dispatch runs, MRs, today, usage, accounts, checks, projects) and the deck session actions, so Nest shows and pilots the owner's real work and not only missions.
- Keywords: nest, deck parity, core sections, sessions, dispatch runs, mrs, today, usage, pilot
- Use when: grooming or building anything that brings deck data or deck actions into core and Nest.
- Skip when: mission model work (task_016, task_017) or Nest visual design passes.

# Needs
- Nest live on the phone shows zeros (owner, 2026-11-10). core only serves the mission store, and no work goes through missions yet. All current work runs in Claude sessions and `crow dispatch`, which deck shows from `bin/` scripts that core never reads.
- Owner, 2026-11-10: deck and Nest are the same tool with different views. Nest must reach everything and pilot everything from the phone.
- req_022 already set this rule: core exposes sessions, runs, MRs, threads, today and usage; every UI reads it; no UI keeps its own logic. adr_001 follow-up asked for dispatch records to be readable by core on day one. task_020 built the mission side only.
- Order (owner, 2026-11-10): read every deck source first, then actions, then the Nest screens. Deck keeps exec'ing scripts in this pass; moving deck onto core is a later item.

# Context
- deck sources (tools/crow-deck `data.go`, `today_data.go`, `mrs_data.go`, `checks_data.go`): `crow-sessions list --json`, `crow-runs json`, `crow-mrs json`, `crow-today json`, `crow-usage --json` (with `CROW_TITLES_PER_DAY=0`), `crow-accounts json`, `crow-checks json`, `crow-registry json`. Timeouts there run from 10 to 90 s.
- deck session actions (`actions.go`, `live.go`): `crow-focus --pid N --send TEXT`, `--interrupt`, `--stop`; `crow-runs stop`; `crow-today done|snooze ID`.
- core already serves `accounts` and `jobs` as external sections: a fetch function run by exec, cached for a TTL (`tools/crow-core/internal/api/api.go` `external()`, `serve_cmd.go`). It runs under the server mutex, so a slow script blocks every request. The new sources need a refresh off the request path.
- prod_005 screens: Home (asks, today, running work as one-liners, alerts; header with quota per account and live run count), Work (Today, MRs and reviews waiting, Sprint), System (accounts, quota, jobs, machines).
- `crow-today json` carries the calendar, so Nest `agenda()` can stop returning `[]`.

# Acceptance criteria
- AC1: With no mission in the store, Nest on the phone shows the same live sessions, dispatch runs, MRs, today items and quota as deck.
- AC2: From Nest on the phone, the owner sends text to a live session, interrupts or stops it, and stops a dispatch run, with the same confirmation deck asks.
- AC3: A slow or failing script never blocks the other sections or the mission endpoints; Nest shows that source as stale with its age.
- AC4: No data logic is copied into Nest or core: core execs the same scripts deck uses.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control` (core exposes every deck source)
- `req_027_conversations_feed_for_core_deck_nest` (session conversation feed, later)
- `req_034_work_views_in_nest`, `req_035_forge_ops_in_nest`
- `tools/crow-deck/data.go`, `tools/crow-core/internal/api/api.go`, `tools/crow-core/serve_cmd.go`

# Backlog
- `item_169_core_serves_every_deck_source_as_read_only_sections`
- `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- `item_171_nest_home_shows_live_sessions_runs_alerts_and_quota`
- `item_172_nest_work_and_system_screens_on_the_deck_sources`

# Out of scope
- deck reading core instead of the scripts (follow-up once item_169 is stable).
- Conversation view of a session (req_027), live terminal attach, butler chat.
- MR review actions (`crow-review-run`), publish actions (adr_001 publish section).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
