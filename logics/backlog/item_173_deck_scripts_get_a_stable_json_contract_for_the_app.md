## item_173_deck_scripts_get_a_stable_json_contract_for_the_app - Deck scripts get a stable JSON contract for the app
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-10 19:27:27

# AI Context
- Summary: rework the `bin/` scripts deck and core both read so they suit a live app: a fixture test per JSON output, a short usage summary, a cached MR fetch, and JSON results from the action scripts.
- Keywords: bin scripts, json contract, fixtures, crow-usage summary, crow-mrs cache, crow-focus json, crow-runs
- Use when: changing the output of a script that deck or core reads.
- Skip when: core wiring (item_169).

# Problem
- The scripts stay the single source of data and actions (owner, 2026-11-10). Measured 2026-11-10:
  | Script | Time | Output |
  | --- | --- | --- |
  | `crow-sessions list --json` | 0.2 s | 27 kB |
  | `crow-runs json` | 0.6 s | 18 kB |
  | `crow-usage --json` | 1.5 s | 334 kB |
  | `crow-accounts json` | 0.2 s | 1 kB |
  | `crow-registry json` | 0.5 s | 19 kB |
  | `crow-checks json` | 0.05 s | empty (no sprint config) |
  | `crow-today json` | 0.1 s cached, 2.3 s fresh | 6 kB |
  | `crow-mrs json` | 11.5 s | 1.5 kB |
- No test pins the JSON shape, so a script change can break Nest silently. Usage is too heavy for a phone poll. MRs are too slow to fetch per refresh. `crow-focus` and `crow-runs stop` print text, so the app cannot show a clean result.

# Scope
- In:
  - Contract: each JSON output gets a `schema` version field and a fixture test (`tests/` per script, run in CI with the other script tests). Breaking a field means bumping the version.
  - `crow-usage --json --summary`: per account quota, reset time, today's cost and tokens, model. Target under 5 kB.
  - `crow-mrs json` cache like `crow-today` (state file, TTL, `-f` to refresh); a cached call under 200 ms.
  - Action scripts answer `--json`: `{ok, action, target, message}` and exit codes 0 ok, 1 refused, 2 error. Covers `crow-focus --send|--interrupt|--stop`, `crow-runs stop`, `crow-today done|snooze`.
  - `crow-runs json` documents that it reads `local/state/runs.jsonl` so core can watch the file (item_169).
  - deck keeps working: run deck's tests and open each tab once.
- Out:
  - Rewriting scripts in Go. New data sources.

# Acceptance criteria
- AC1: Each script's JSON output has a `schema` field and a fixture test that fails on a renamed or removed field.
- AC2: `crow-usage --json --summary` answers under 5 kB with the same quota numbers deck shows.
- AC3: A cached `crow-mrs json` answers under 200 ms; `-f` refetches.
- AC4: Each action script with `--json` prints one JSON result and exits 0, 1 or 2.
- AC5: deck tests pass and every tab still loads.

# AC Traceability
- request-AC3 -> AC2, AC3. Proof: light, fast sources. `crow-usage --json --summary` is 1.7 kB with the deck quota fields; cached `crow-mrs json` 0.09 s, `-f` 10.9 s (5b4a90d).
- request-AC4 -> AC1, AC5. Proof: one source for deck and core, pinned by tests. Schemas, fixtures and `tests/test_script_contract.py` in 5cbdd7d (rename/remove mutation fails every schema; action results exit 0/1/2); `python bin/crow-check` green; deck and core `go test` pass; deck tabs 1-5 render headless.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (owner 2026-11-10: scripts stay the single source, rework in its own item before item_169)

# Links
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Unblocks: `item_169_core_serves_every_deck_source_as_read_only_sections`, `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: first step of req_054; item_169 wires these outputs as they are.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
