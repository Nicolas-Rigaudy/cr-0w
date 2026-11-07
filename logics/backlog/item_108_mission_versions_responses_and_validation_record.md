## item_108_mission_versions_responses_and_validation_record - Mission versions, responses and validation record
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:50

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
- request-AC4 -> AC1. Proof: replay of the example mission from files in the test suite.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: No publish path is safe without it.

# Notes
- Record shape shared with item_099 and req_030 Need 5.
- Depends on item_106 and the core event log (adr_001).
- Code home waits on the Go core and the shared fixture contract (item_119). The backend is Go (owner decision 2026-11-07); D6 keeps only the parser in Python `bin/`.
- Spec: `logics/specs/req_037_run_contract/`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
