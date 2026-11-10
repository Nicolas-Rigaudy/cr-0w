## item_169_core_serves_every_deck_source_as_read_only_sections - Core serves every deck source as read-only sections
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 100%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: new external sections in crow-core (`sessions`, `dispatch`, `mrs`, `today`, `usage`, `checks`, `projects`) built by exec'ing the same `bin/` scripts deck uses, refreshed in the background, served with ETags; Nest's data layer reads them.
- Keywords: core, external sections, crow-sessions, crow-runs, crow-mrs, crow-today, crow-usage, background refresh, stale, etag
- Use when: adding a deck data source to core or changing how external sections refresh.
- Skip when: mission store sections (missions, runs, inbox) or actions (item_170).

# Problem
- core serves only the mission store plus accounts and jobs. Nest shows zeros while deck shows live sessions, 121 dispatch runs, MRs, mail and quota (req_054).
- `external()` runs its script under the server mutex. A 90 s `crow-mrs json` would block every request.

# Scope
- In:
  - Sections, each the script's JSON as `data`, unchanged:
    | Section | Command | TTL |
    | --- | --- | --- |
    | `sessions` | `crow-sessions list --json` | 5 s |
    | `dispatch` | `crow-runs json` | 10 s |
    | `usage` | `crow-usage --json` with `CROW_TITLES_PER_DAY=0` | 60 s |
    | `today` | `crow-today json` (its own 5 min cache) | 60 s |
    | `mrs` | `crow-mrs json` | 120 s |
    | `checks` | `crow-checks json` | 300 s |
    | `projects` | `crow-registry json` | 300 s |
  - TTLs above are starting values in config, tuned during the build against the measured script time.
  - Refresh off the request path: one goroutine per section, at most one run at a time, timeout per script (deck's values). A request returns the last good body at once.
  - Section envelope gains `stale: true` and `error` when the last refresh failed; `updated_at` stays the last good time. First load with no data returns `data: []` and `stale: true`.
  - A section is refreshed only while it was requested in the last 2 min, so an idle phone costs no script runs.
  - `accounts` and `jobs` move to the same refresher.
  - Response fixtures for each section (`go test ./internal/api -update`), as for item_144.
  - Nest data layer: `DataSource` gains one generic `section(name)` reader over the item_167 ETag cache; `snapshot` carries the new sections; `agenda()` maps `today.events`. No screen work beyond keeping the build green.
- Out:
  - Screens (item_171, item_172), actions (item_170).
  - Changing the scripts' JSON. If a field is missing for Nest, note it and add it in the script, never in core.
  - deck reading core.

# Acceptance criteria
- AC1: `curl` over the socket returns each new section with data equal to the script's own `json` output at that time.
- AC2: With `crow-mrs` made to sleep 100 s (test double), `/api/state/missions` and `/api/state/sessions` still answer in under 100 ms.
- AC3: A failing script returns the last good data with `stale: true`, the error text, and the old `updated_at`; it clears on the next good run.
- AC4: An unchanged section answers 304 to its ETag.
- AC5: No section is refreshed when nothing requested it for 2 min (test with a fake clock).
- AC6: Nest dev build reads every new section through `useSnapshot()`; Agenda shows today's events from `today`.

# AC Traceability
- request-AC1 -> AC1, AC6. Proof: same data as deck, readable from Nest.
- request-AC3 -> AC2, AC3. Proof: slow and failing scripts isolated.
- request-AC4 -> AC1. Proof: core execs the deck scripts unchanged.

# Decision framing
- Product framing: Not needed (prod_005 surfaces)
- Architecture framing: Not needed (adr_001: one API for every surface; req_022: core exposes every deck source)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Depends on: `item_144_core_http_service_with_read_endpoints_on_the_tailnet`, `item_167_nest_reads_live_data_from_the_core_api`
- Unblocks: `item_170_core_actions_on_live_sessions_and_dispatch_runs`, `item_171_nest_home_shows_live_sessions_runs_alerts_and_quota`, `item_172_nest_work_and_system_screens_on_the_deck_sources`, `item_165_telemetry_strip_on_the_laptop_home`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: Nest shows nothing useful without it (owner, 2026-11-10). Next on task_020's critical path.

# Notes
- Build in the hub repo (crow-core is in cr-0w), fresh session.
- Script data reaches only the owner's devices: core checks the tailnet identity (item_144).
- Built 2026-11-10 in 50ad41a (core) and 1011e0e (Nest). `internal/api/refresh.go` holds one refresher per external section: last good body served at once, one run at a time, ticked every second, idle after 2 min. `accounts` (crow-accounts) and `jobs` moved onto it. TTLs are `serve.ttl` in core.yaml (defaults in `config.DefaultTTL`); script timeouts in `serve_cmd.go`. Documented in `runbooks/core-service.md`.
- Envelope: `stale` and `error` appear together, `updated_at` stays the last good change. First load is `data: []`, stale, `error: not loaded yet`. A stale body gets its own ETag, so 304 ends when a run fails.
- Nest: `DataSource.section(name)` over a shared ETag cache (concurrent reads share one request), `Snapshot.sections`, `agenda()` from `today.events` with all-day markers left out. Fixture source returns empty sections until screens bring state fixtures (171/172).
- Proof: AC1 test core on a scratch socket, every section equal to its script except fields that move between runs (`generated_at`, `now`, live token counts; two direct script runs drift the same way). AC2-AC5 `go test -race ./internal/api` (TestSlowScript, TestFailingScript, TestIdleSection with a fake clock); live 304 on `projects`, `/api/state` in 9 ms. AC6 vitest feed test (snapshot carries all 9 sections, agenda = `agendaOf(today)`), typecheck, build. Live core rebuilt and restarted: `today` serves 3 events.
- Left for the owner: open Nest on the phone and check the agenda panel shows today's meetings.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
