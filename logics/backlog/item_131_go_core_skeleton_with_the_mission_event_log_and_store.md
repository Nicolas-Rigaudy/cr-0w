## item_131_go_core_skeleton_with_the_mission_event_log_and_store - Go core skeleton with the mission event log and store
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85
> Confidence: 80
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 16:47:13

# AI Context
- Summary: Start the Go core as `tools/crow-core`: config paths, the append-only mission event log and a SQLite store rebuilt from it. No HTTP yet; a one-shot binary the Python CLI calls.
- Keywords: core, go, event log, jsonl, sqlite, store, rebuild, config
- Use when: Writing the first core code or any code that reads or writes mission events.
- Skip when: HTTP API (req_028 slice 3), nest, deck views.

# Problem
- adr_001 makes an append-only event log the source of truth and the store a rebuildable view. Nothing exists yet.
- req_031 decided a Go core with the CLI staying Python. Slices 1 and 2 need core logic before the service and the API exist.

# Scope
- In:
  - Go module `tools/crow-core` (own `go.mod`, next to `tools/crow-deck`), binary `crow-core`, built by `bin/crow` the way the deck is built. CI already loops over Go dirs; check it picks the new one up.
  - Config: mission dir root, store path, kind files dir, private overrides dir. Defaults under `$CROW_HOME/local/state/missions/` (machine state). Paths from config so core can move host later.
  - Event log per mission `missions/<id>/events.jsonl` with `{ts, mission, type, actor, data}` (req_037 spec "Mission dir"). Append with fsync and a file lock. Types used by slices 1-2: `created`, `state`, `run_start`, `run_end`, `gate`, `question`, `answer`, `validation`, `publish`, `incident`. Event, mission, run and inbox shapes are fixed by `logics/specs/core_contract` (item_119); decode them with the `crow-core/contract` package and add fixtures there for any new type.
  - Store: SQLite through a pure-Go driver (no cgo). Tables for missions, runs, inbox. `crow-core rebuild` drops and replays every `events.jsonl`; the store records the last applied offset per mission so normal writes replay only the tail.
  - Event fixtures from item_119 drive the Go tests.
- Out:
  - HTTP, SSE, unix socket, systemd unit (slice 3).
  - FTS index, stats, backups (later slices).
  - Mission semantics (item_132).

# Acceptance criteria
- AC1: `crow-core rebuild` on the fixture missions gives the same store rows as writing those events one by one.
- AC2: Two concurrent appends to one mission never interleave a line (test with parallel writers).
- AC3: Deleting the SQLite file and running `rebuild` restores every mission, run and inbox row.
- AC4: `go vet` and `go test` for `tools/crow-core` run in CI.

# AC Traceability
- request-AC1 -> AC1, AC3. Proof: replay tests on fixture missions.

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001 (event log and store) and req_031 decision 1 (Go core).

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: Every later core item writes or reads events.

# Notes
- Groomed 2026-11-09 from req_028 slice 1. Waits on item_119 (fixtures and envelope).
- Code home `tools/crow-core` chosen at grooming to match `tools/crow-deck`. The req_020 shared Go kit repo is UI only and does not host core.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
