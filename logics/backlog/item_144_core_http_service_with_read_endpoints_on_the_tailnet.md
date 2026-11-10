## item_144_core_http_service_with_read_endpoints_on_the_tailnet - Core HTTP service with read endpoints on the tailnet
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 16:17:33

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
- Tailscale identity headers need v1.92+ (req_031 sources). Installed on archeotech: 1.104.5 (checked 2026-11-10).
- Groomed 2026-11-09 from req_028 slice 3.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
