## item_108_mission_versions_responses_and_validation_record - Mission versions, responses and validation record
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 75%
> Progress: 100%
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 11:19:08
> Owner: crow

# AI Context
- Summary: core keeps every version file with its answers file, checks the Responses section, and writes one hash-chained validation record per publish action on ok.
- Keywords: versions, answers, responses, validation record, payload hash, hash chain
- Use when: Handling an owner answer at a gate.
- Skip when: The publish actions themselves (item_099).

# Problem
- Each answered version must stay with the owner's words.
- `ok` must cover exactly the payload shown; a later message must cancel it.

# Scope
- In:
  - One version counter per mission shared by plan, report and mockup files.
  - Write `answers.vN.md` (front matter plus numbered points) from each answer; append further words before the verb.
  - Check the `responses` list covers every answer point and question of the previous version.
  - On ok, go or approve where they publish: one `validation` event per publish entry with payload hash, owner event and chain hash.
  - Any owner message after it writes `validation_dropped`.
- Out:
  - Calling the forge or posting (item_099).
  - The UI that shows the history.

# Acceptance criteria
- AC1: The spec example mission replays to the same events, including both validation hashes.
- AC2: A changed publish entry after ok fails the hash check; a message after ok drops the validation.
- AC3: An edited or replayed validation line breaks the chain check.
- AC4: A version 2 that skips an answer point is refused with the missing point named.

# AC Traceability
- request-AC4 -> AC1. Proof: Implemented in 0e7b72f; validated with go test ./... and bin/crow-check. TestAnswerReplaysTheExampleMission replays the example mission: all 22 log lines and answers.v1-v3 match byte for byte, both validation hashes included. Source: `0e7b72f`
- AC2 -> `internal/mission` TestVerifyPassesThenRefusesAChangedEntry (a publish entry changed after ok gives a `validation_mismatch` incident and a refusal) and TestMessageAfterOkDropsTheValidation (a msg after ok writes `validation_dropped`; a new ok chains after the first two).
- AC3 -> `internal/mission` TestChainBreaksOnEditOrReplay: an edited line, a replayed line and a forged actor each fail `CheckChain` with the line named.
- AC4 -> `internal/mission` TestCoverageNamesTheMissingPoint and TestTickRefusesAVersionThatSkipsAPoint: a v2 that skips a2 fails the mission with "v2 responses miss a2".

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: No publish path is safe without it.

# Notes
- Record shape shared with item_099 and req_030 Need 5.
- Depends on item_106 and the core event log (adr_001).
- 2026-11-09 built (0e7b72f): `internal/mission/answer.go`. `env.Answer(ref, Reply)` writes `answers.vN.md`, the owner `answer` event, a `validation_dropped` per live validation, then one `validation` per publish entry on ok/go/approve. `Verify` checks the chain and the payload hashes against the file on disk. `CheckChain` checks the owner actor, `prev_sha256` continuity and the owner event.
- Event lines are now canonical JSON (sorted keys, `contract.Canonical`), like the spec example.
- The `answer` event omits the surface; it lives in the answers file, as in the example.
- A Verify mismatch keeps the mission in `review` (no legal move to a parked state from there).
- Nothing moves review -> running after dig/msg yet. For item_134: `crow inbox answer` calls `env.Answer`, then that move.
- Code home waits on the Go core and the shared fixture contract (item_119). The backend is Go (owner decision 2026-11-07); D6 keeps only the parser in Python `bin/`.
- Spec: `logics/specs/req_037_run_contract/`.
- 2026-11-09 groomed into task_016. Code home: Go in `tools/crow-core`; event log from item_131. Called by `crow inbox answer` (item_134).
- Task `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox` was finished via `logics-manager flow finish task` on 2026-11-10.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`
