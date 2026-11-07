## item_119_go_core_and_python_sidecar_contract_with_shared_fixtures - Go core and Python sidecar contract with shared fixtures
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:22:08

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
- request-AC1 -> AC1 to AC4. Proof: the contract doc and the shared fixtures.
- request-AC2 -> AC2. Proof: the hello-world slice uses the Go core and one sidecar message.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_031_spike_nest_stack_and_app_shape`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Gates the core work and the code home of items 107 to 109.

# Notes
- 2026-11-07 owner decision 22 (req_031 decision 1).
- Frees the code homes of items 107, 108 and 109 to be Go.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
