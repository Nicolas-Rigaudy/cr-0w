## item_132_mission_record_states_and_the_crow_mission_cli - Mission record, states and the crow mission CLI
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:20:11

# AI Context
- Summary: The mission record of adr_001, its state machine, and `crow mission new|list|show|cancel` (Python CLI over `crow-core`).
- Keywords: mission, record, state machine, kind, scope, autonomy, gates, crow mission
- Use when: Creating, listing or cancelling missions.
- Skip when: Starting runs (item_133) or answering questions (item_134).

# Problem
- Missions are the unit the owner thinks in (req_028 need 2). They need a record, legal state moves and a CLI before any surface exists.

# Scope
- In:
  - Record fields from adr_001 "Mission record": id `YYYYMMDD-kind-short`, kind, batch, scope, data_class, project, goal, end_state, autonomy, gates, state, runs, version, validation, file_scope, budget, branch, worktree, mr.
  - Defaults from the kind file (gates, autonomy, budget, caps, account, effort) read from the kind files dir with private overrides when item_117 lands. The Go side reads the same YAML the Python loader of item_104 validates; the item_119 fixtures guard the shape.
  - Scope and project checked against `local/projects.yaml`; the account comes from the item_095 scope pick at dispatch, not at creation.
  - States `preflight → queued → running → waiting → review → publishing → done | failed | cancelled`. Illegal moves refused with the current state named. Every move is a `state` event.
  - `crow-core mission new|list|show|cancel` with `--json`; `crow mission ...` in `bin/crow` passes through. `new` takes goal, kind, project, end state; `preflight` is skipped until item_113 (mission goes straight to `queued`).
  - `cancel` on a running mission stops its cdx run (`cdx` stop) before moving to `cancelled`.
- Out:
  - Preflight questions (item_113).
  - Batches (req_028 need 10).
  - Starting runs (item_133).

# Acceptance criteria
- AC1: `crow mission new` with a phase 1 kind creates the mission dir, a `created` event and a `queued` mission with the kind defaults filled.
- AC2: `crow mission list --json` and `show --json` print records matching the item_119 mission fixture shape.
- AC3: A table test covers every legal and illegal state move.
- AC4: An unknown kind or project is refused with the valid names listed.

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: CLI test creating and listing a mission.
- request-AC3 -> AC1. Proof: autonomy and gates filled from the kind file.
- AC1 -> `internal/mission` TestNewFillsKindDefaults: a research mission gets its dir, `created` (owner) then `state` preflight->queued (core), gates [], autonomy auto, budget 60 min/$5/2 resumes written into `created`. TestNewAppliesScopeOverridesAndFlags covers by_scope (build, work), aliases, flags and a duplicate id. Checked end to end with `crow mission new|list|show|cancel` on a scratch CROW_HOME.
- AC2 -> `internal/mission` TestListAndShowMatchTheFixtureShape: list output has the key set of `fixtures/records/missions.json` and decodes strictly into `contract.Mission`; show takes an id or a unique prefix.
- AC3 -> `contract` TestStateMoves walks every from x to pair against the legal table; `internal/project` TestIllegalMoveNamesTheCurrentState and the "illegal move" case of TestApplyRefusesBrokenSequences check the fold refuses them.
- AC4 -> `internal/mission` TestUnknownKindOrProjectListsValidNames (kind names from the kind dirs, project names from the registry, scope outside the kind's list).

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001.

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: Slice 1 of the 0.3b core.

# Notes
- Groomed 2026-11-09 from req_028 slice 1. Waits on item_131.
- 2026-11-09 built in `internal/mission` and `crow-core mission`; `crow mission` passes through. The registry path is `projects` in config (default `local/projects.yaml`).
- Resolution order: kind file (override dir first), then its `by_scope` entry, then the flags. Scope defaults to the project's registry scope; data class defaults to the scope (`crow` gives `personal`).
- The key defaults to the first five goal words without filler, and must match the kind's key pattern. Kinds with a ticket key (build) need `--key`. An existing id is refused.
- Legal moves live in `contract.Moves`. `done -> queued` (dig) and `failed -> queued` (retry) are the only ways back.
- Creation accepts any kind whatever its phase. item_133 must gate dispatch on the kind phase.
- `cancel` on a running mission runs `$CROW_ENGINE/bin/crow-runs stop` for each live run. A stop failure is a warning; the mission is cancelled anyway.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
