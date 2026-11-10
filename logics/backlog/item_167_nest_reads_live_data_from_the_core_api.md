## item_167_nest_reads_live_data_from_the_core_api - Nest reads live data from the core API
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 16:17:25

# AI Context
- Summary: an `apiSource()` in `tools/nest/src/data/` implements the `DataSource` interface over the item_144 read endpoints, replaces `fixtureSource()` in `main.tsx`, and polls with ETags while the page is visible.
- Keywords: nest, api adapter, data source, etag, polling, core down, live data
- Use when: wiring Nest screens to core, or changing how Nest fetches.
- Skip when: actions from Nest (item_168), SSE (item_146), Web Push (slice 7).

# Problem
- Nest renders fixtures only (item_137). The `DataSource` interface was built so an API adapter replaces the fixture adapter in one file (task_018 AC3), but no item scoped that adapter.
- Without it Nest cannot show real missions, inbox or accounts, so the app is not usable.

# Scope
- In:
  - `src/data/api.ts`: `apiSource(base = '/api')` with `missions()`, `mission(id)`, `inbox()`, `agenda()`, `now()`. Same-origin fetch (item_144 serves Nest and the API on one origin); no CORS.
  - Mapping from the item_144 response shapes to the existing `types.ts`. Nest types already are the core contract records, so the mapping stays thin; any field the fixtures had and the API lacks is listed in the item notes, not invented.
  - `agenda()` returns `[]` until a core agenda section exists; the Agenda placeholder stays.
  - `now()` is the wall clock.
  - Polling: refetch each section every 5 s with `If-None-Match` while `document.visibilityState` is visible; stop when hidden, refetch all on visible (req_031 decision 2 fallback; SSE replaces the timer in item_146). `useSnapshot()` re-renders only on a changed body.
  - "core down" state: a fetch error keeps the last data and shows the section age in the header; first load with no data shows an empty state, never fixtures.
  - Adapter choice in `main.tsx` only: API in the production build, fixtures in `npm run dev` and tests unless `VITE_NEST_SOURCE=api`. `npm run dev` proxies `/api` to the core loopback port for local work.
  - Vitest for the adapter against the item_144 response fixtures (304 path, error path, mapping).
  - The service worker never caches `/api/*` (stub today; keep it that way when caching lands).
- Out:
  - Any POST (item_168), SSE (item_146), push, offline cache.
  - Hero data wiring beyond what the interface already gives (item_158 reads the same `useSnapshot()`).

# Acceptance criteria
- AC1: On the phone over the tailnet, Home, Missions, Mission detail and Inbox show the same missions and inbox items as `crow mission list --json` and `crow inbox list --json`.
- AC2: A mission dispatched from the CLI appears in Nest within 10 s with the page open, without a reload.
- AC3: With core stopped, Nest keeps the last data and shows "core down" with the section age; it recovers by itself when core is back.
- AC4: The swap touches `src/data/` and `main.tsx` only; no screen component changes (task_018 AC3).
- AC5: Adapter tests pass against the item_144 response fixtures; bundle stays under the 90 kB gzip budget.

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: phone reads live state over the tailnet (req_028 slice 3, read half in Nest).

# Decision framing
- Product framing: Not needed (prod_005)
- Architecture framing: Not needed (adr_001 API section, req_031 decision 2)

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`, `req_031_spike_nest_stack_and_app_shape`
- Depends on: `item_144_core_http_service_with_read_endpoints_on_the_tailnet`, `item_137_nest_app_shell_on_core_fixtures`
- Unblocks: `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`, `item_165_telemetry_strip_on_the_laptop_home` (live accounts and runs)
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: The step that makes Nest show real data. Owner asked for the fastest path to a working web app (2026-11-10).

# Notes
- Groomed 2026-11-10. The README line "the live API adapter replaces it in src/main.tsx only" had no item until now.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
