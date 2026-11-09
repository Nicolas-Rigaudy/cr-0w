## item_131_go_core_skeleton_with_the_mission_event_log_and_store - Go core skeleton with the mission event log and store
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
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
- AC1 -> `internal/store` TestOneByOneEqualsRebuild: each fixture event appended through the log and synced alone gives the same dump as a rebuild; applied offset equals the log end after each sync.
- AC2 -> `internal/eventlog` TestParallelAppendsInOneProcess (8 goroutines x 25 lines of 70 KB) and TestParallelAppendsAcrossProcesses (4 re-exec writer processes x 20): every line decodes and the count matches.
- AC3 -> `internal/store` TestDeletedStoreRebuilds (db, -wal and -shm removed, rebuild, identical dump: 2 missions, 4 runs, 4 inbox rows). Also checked end to end through `crow core append|dump|rebuild` on a scratch CROW_HOME.
- AC4 -> `.github/workflows/ci.yml` already loops over `tools/crow-core` with `go vet ./... && go test -count=1 ./...`.

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
- 2026-11-09 built. Packages: `internal/config` (paths, optional `local/core.yaml`), `internal/eventlog` (append under flock + fsync, tail read that leaves a partial last line), `internal/project` (pure fold of events into mission, runs, inbox rows), `internal/store` (modernc SQLite, WAL, immediate tx; Sync, SyncAll, Rebuild, Append, Dump). `crow core <cmd>` builds and runs the binary like the deck.
- `crow-core append` goes through `Store.Append`: the event must apply to the mission as it stands before it reaches the log, under the store write lock. A refused event never poisons the log. Writers must use it rather than append to `events.jsonl` directly.
- Inbox rows carry only what events say. `ask`, `artifacts` and `blocking` stay empty for item_134 to fill from the version files and `questions.json`.
- A row's `created` is the ts of the event that opens it (`gate` ready or `question`). The research fixture had the following `state` ts; `fixtures/records/inbox.json` was aligned (14:52:01).
- A straight-through report row is `fyi` with verbs `[dig]`, as in the fixture. The req_037 table says "none (msg reopens it)"; reconcile there when item_134 lands.
- A `created` event without gates, autonomy or budget takes the kind file defaults (override dir first). item_132 should write the resolved values into `created` so replay never depends on kind files changing.
- `validation_dropped` is accepted but not folded; item_108 owns the record chain.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
