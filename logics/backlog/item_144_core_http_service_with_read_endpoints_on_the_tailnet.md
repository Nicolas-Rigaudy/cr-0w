## item_144_core_http_service_with_read_endpoints_on_the_tailnet - Core HTTP service with read endpoints on the tailnet
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 17:29:46
> Owner: crow

# AI Context
- Summary: core runs as a systemd user service and serves read-only HTTP on localhost and the tailnet: `/state` sections (inbox, missions, runs, accounts, jobs) with `updated_at` and ETag, mission detail, artifacts and the event feed.
- Keywords: crow-core, http, systemd, tailnet, tailscale serve, identity header, csrf, etag, state sections, read api
- Use when: building or changing core's HTTP server, its auth, or a read endpoint.
- Skip when: action endpoints (item_147), live stream (item_146), conversation feed (item_148).

# Problem
- Every surface (deck, nest, butler, bots) must read one API (adr_001). Today core is a CLI plus a one-shot tick; deck reads bin/ scripts and state files.
- Slice 3 of req_028, folding in the req_022 core item: read endpoints first, actions second.

# Scope
- In:
  - `crow-core serve` in `tools/crow-core` (Go `net/http`, req_031 decision 1), a systemd user unit installed by `crow core install` and `crow core status|restart` (req_022 names).
  - Bind 127.0.0.1 plus the tailnet through `tailscale serve` (HTTPS). Owner check on the Tailscale identity header; refuse when missing on a non-loopback request. POST routes need a custom header (CSRF), checked from this item on so item_147 inherits it.
  - `GET /state` and `GET /state/<section>` for inbox, missions, runs, accounts, jobs. Each section carries `updated_at` and an ETag keyed on the store version; `Cache-Control: no-cache` (req_031 decision on fetching).
  - `GET /missions/<id>` (record, versions, questions, answers, runs), `GET /missions/<id>/artifacts/<path>` (outbox files, size caps from item_109, path traversal refused), `GET /events?after=<id>`.
  - JSON shapes typed from the item_119 fixtures; new fixtures for each response so Go and Python tests share them.
  - The tick keeps running from its timer; serve reads the store and does not tick.
  - One origin for Nest and the API (amended 2026-11-10, task_020): every API route lives under `/api` (`/api/state/<section>`, `/api/missions/<id>`, `/api/events`; the paths in items 145-148 read relative to it). `/` serves the built Nest from a configured dir (default `tools/nest/dist`), with an `index.html` fallback for client routes and no-cache on `index.html` and `sw.js`. Reason: req_031 sends no CORS headers ever, so Nest must be same-origin with core. The identity check covers the static routes too.
  - Response fixtures for the sections Nest reads (missions, inbox, accounts, runs) are the input of item_167.
- Out:
  - Actions (item_147), SSE (item_146), Web Push, MCP server (butler slice 8), conversation endpoint (item_148).
  - Moving core to the homelab (paths stay in config).

# Acceptance criteria
- AC1: `crow core install` starts the unit; it survives a logout and login; `crow core status` shows the listen addresses.
- AC2: `curl` from the phone over the tailnet reads `/state/inbox`; the same request without the owner identity header (from loopback through a proxy test) is refused.
- AC3: A repeat GET with `If-None-Match` returns 304 until the store version changes.
- AC4: Response fixtures exist for every endpoint and pass in Go and Python tests.
- AC5: An artifact request outside the mission dir returns 404.

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: the phone reads the inbox over HTTP (part of the slice 3 exit with item_145).
- request-AC2 -> AC4. Proof: runs section fixture carries the summary field filled by item_149.

# Decision framing
- Product framing: Not needed (prod_005)
- Architecture framing: Not needed (adr_001 API section, req_031 decisions 1 and 2)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control`
- Depends on: `item_131_go_core_skeleton_with_the_mission_event_log_and_store`, `item_134_inbox_items_and_the_crow_inbox_cli` (both Done). Not on item_135: reads work while resume is blocked.
- Unblocks: `item_167_nest_reads_live_data_from_the_core_api` (Nest API adapter), item_145, item_146, item_147
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api` (moved from task_017 on 2026-11-10)

# Priority
- Priority: High
- Rationale: First step of slice 3 and of the Nest live path (task_020).

# Notes
- Built 2026-11-10 (02a08dc, ecfe6c1, 76534a4):
  - Listen `127.0.0.1:8787` (config `serve.listen`), unix socket `local/state/core.sock` (0600, trusted, no identity header), Nest from `tools/nest/dist`.
  - Owner login: config `serve.owner_login`, else `tailscale status --json` (Self user), retried at most once a minute. TCP requests without a matching `Tailscale-User-Login` get 403 JSON `{"error": ...}`.
  - Non-GET routes need `X-Crow: 1` and a same-origin `Origin`; the static route answers 405 to writes.
  - Sections share the envelope `{section, version, updated_at, data}`; `GET /api/state` is `{version, sections}`. Store sections are cached on the store version (`<missions>.<sum of applied>`); accounts (`bin/crow-accounts json`) and jobs (`systemctl --user list-timers`, crow-* only) are cached 30 s and keep the last good value on error. `?all=1` on the inbox adds closed rows.
  - Mission detail `{mission, versions[{version, file, summary, questions, answers}], runs}`. Artifacts: regular files under `outbox/` within the kind file cap, served with `Content-Security-Policy: sandbox` and nosniff; anything else 404 (dot segments get the mux redirect to the cleaned path first).
  - Events `{events, next, more}`, id `<ts UTC µs>~<mission>~<byte offset>`, `after=<id>`, `limit` up to 500.
  - Fixtures in `logics/specs/core_contract/fixtures/api/` (Go writes them with `-update`, Python checks them against the record schemas). That check found two drifts, fixed in 02a08dc: the inbox schema lacked `needs-input`, and the fold left `data_class` empty for created events that omit it.
  - Runbook: `runbooks/core-service.md`.
- Tailscale identity headers need v1.92+ (req_031 sources). Installed on archeotech: 1.104.5 (checked 2026-11-10).
- Groomed 2026-11-09 from req_028 slice 3.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
