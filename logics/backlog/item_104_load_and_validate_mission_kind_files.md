## item_104_load_and_validate_mission_kind_files - Load and validate mission kind files
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 92%
> Confidence: 88%
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 01:59:45

# AI Context
- Summary: A stdlib validator for the JSON Schema subset the spec uses and a repo test that validates every spec file, so the spec stays valid while it is reviewed.
- Keywords: kind file, loader, validator, json schema subset, spec test
- Use when: Reading kind files or checking the spec files stay valid.
- Skip when: Reading kinds at runtime (the loader waits on D6 and the owner review).

# Problem
- The spec defines kinds as YAML files and a JSON Schema, but nothing checks them.
- `jsonschema` is not installed and the repo avoids new dependencies.
- A spec check is a test, not core code: it is useful whatever backend req_031 picks.

# Scope
- In:
  - A stdlib validator for the keywords the spec schemas use: type, enum, const, required, properties, additionalProperties, pattern, minimum, maximum, minLength, maxLength, items, minItems, uniqueItems, $ref with $defs, allOf, if and then, propertyNames. `format` is an annotation only. YAML through PyYAML, already used by `bin/`.
  - A repo test that validates the 8 kind files, the three schemas' examples and the example mission against their schemas, plus negative cases.
  - The validator fails on a keyword it does not know, so a schema edit cannot pass silently.
- Out:
  - A runtime loader that reads kinds from a configured dir (waits on D6 and item_105; goes with the first consumer).
  - Replacing `RUN_KINDS` in dispatch (after the owner review).
  - The core port if req_031 picks Go.

# Acceptance criteria
- AC1: The validator accepts the 8 kind files and rejects a file with a bad phase, duplicate gates or an unknown field.
- AC2: An invalid file fails with the file path and the JSON path of the bad field.
- AC3: A test in the repo suite validates every file under the spec dir against its schema and runs in `bin/crow-check`.
- AC4: A schema using a keyword outside the subset fails the test.

# AC Traceability
- request-AC1 -> AC1. Proof: the 8 kind files validate in the test suite.
- request-AC2 -> AC3. Proof: the spec examples validate in the test suite.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_037_mission_kinds_and_run_contract`
- Primary task(s): `task_011_load_and_validate_mission_kind_files`

# Priority
- Priority: High
- Rationale: Every other item reads kinds or checks a schema.

# Notes
- Spec: `logics/specs/req_037_run_contract/`.
- A throwaway prototype of the validator passed on the spec files during grooming.
- First item to build: it needs no owner answer.
- Python because `bin/` and logics-manager are Python. A Go backend (req_031) ports the loader, the schemas stay.
- Task `task_011_load_and_validate_mission_kind_files` was finished via `logics-manager flow finish task` on 2026-11-07.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_011_load_and_validate_mission_kind_files`
