## item_119_go_core_and_python_sidecar_contract_with_shared_fixtures - Go core and Python sidecar contract with shared fixtures
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 16:33:30

# AI Context
- Summary: Define the contract between the Go core and the Python sidecars: local JSON messages, versions and shared fixtures that both sides test against.
- Keywords: go core, python sidecar, local json, fixtures, contract
- Use when: Before writing the Go core or any AI sidecar.
- Skip when: Front end work.

# Problem
- Owner decision 2026-11-07: the nest backend is a Go core (HTTP/SSE, state, dispatch to laptop and tower, supervision, a module shared with crow-deck) plus Python sidecars for AI pieces (butler/Agent SDK, STT, embeddings) over local JSON. The CLI stays Python.
- Two languages share formats (kinds, crow block, events, questions). Without fixtures the formats drift.

# Scope
- In:
  - Message envelope over local JSON (transport chosen in the item: stdio or unix socket), with a version field and error shape.
  - A fixtures directory read by Go tests and Python tests: crow block, questions, events, kinds, sidecar requests and replies.
  - Sidecar supervision rules: start, health, restart, timeouts.
  - Which formats the Go module shares with crow-deck.
- Out:
  - Writing the core itself.
  - The Python parser in bin/ (item_106) except as a fixture consumer.

# Acceptance criteria
- AC1: The contract doc lists every message with a request and a reply fixture.
- AC2: A Go test and a Python test read the same fixtures and both pass.
- AC3: A change to a fixture without a matching change on the other side fails CI.
- AC4: The shared Go module boundary with crow-deck is written down.

# AC Traceability
- request-AC1 -> AC1 to AC4. Proof: Implemented in c2e7799, 9ae4701, 923eae4, d524c68. `logics/specs/core_contract` lists hello, health and shutdown with request and reply fixtures (AC1); `tests/test_core_contract.py` and `tools/crow-core/contract` read the same fixtures, 30 and 9 tests pass (AC2); a fixture field added on one side fails both suites, checked by hand (AC3); README section "Shared Go module boundary with crow-deck" (AC4).
- request-AC2 -> AC2. Proof: Implemented in c2e7799, 9ae4701, 923eae4, d524c68. `TestPythonSidecar` runs `bin/crow_sidecar.py` from Go over a real pipe and matches the lifecycle fixtures.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_031_spike_nest_stack_and_app_shape`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: Medium
- Rationale: Gates the core work and the code home of items 107 to 109.

# Notes
- 2026-11-07 owner decision 22 (req_031 decision 1).
- Frees the code homes of items 107, 108 and 109 to be Go.
- 2026-11-09 groomed into task_016 (first step). Fixtures also cover the mission record, the run entry and the inbox row read by items 132 and 134. Go side lives in `tools/crow-core` (item_131).
- 2026-11-09 done. Go package `crow-core/contract` holds envelope, mission, run, inbox, event, crow block, questions and kind types. A schema walk test ties Go json tags to the schemas. Python base sidecar `bin/crow_sidecar.py`. CI runs `go test -count=1` because the Go test cache does not track files outside the module. Event names are `run_start`/`run_end` everywhere; the req_037 build example log now carries the full state chain.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
