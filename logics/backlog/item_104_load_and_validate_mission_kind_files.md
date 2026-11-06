## item_104_load_and_validate_mission_kind_files - Load and validate mission kind files
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A stdlib validator for the JSON Schema subset the spec uses, a loader that reads kind files from a configured path, and a test that validates every spec file.
- Keywords: kind file, loader, validator, json schema subset, spec test
- Use when: Reading kind files or checking the spec files stay valid.
- Skip when: Wiring kinds into dispatch (needs the owner review item).

# Problem
- The spec defines kinds as YAML files and a JSON Schema, but nothing checks them.
- `jsonschema` is not installed and the repo avoids new dependencies.
- core needs one loader; crow dispatch can share it while core is not built.

# Scope
- In:
  - A stdlib validator for the keywords the spec schemas use: type, enum, const, required, properties, additionalProperties, pattern, min and max, items, minItems, uniqueItems, $ref, allOf, if and then, propertyNames.
  - A loader that reads `<dir>/<kind>.yaml`, validates it and returns the kind. The dir comes from config, never a hard-coded path.
  - A test that validates the 8 kind files and every example against its schema, plus negative cases.
- Out:
  - Replacing `RUN_KINDS` in dispatch (after the owner review).
  - The core port if req_031 picks Go.

# Acceptance criteria
- AC1: The validator accepts the 8 kind files and rejects a file with a bad phase, duplicate gates or an unknown field.
- AC2: The loader reads the kind dir from config and fails with the file and field on an invalid kind.
- AC3: A test in the repo suite validates every file under the spec dir against its schema.

# AC Traceability
- request-AC1 -> AC1. Proof: the 8 kind files validate in the test suite.
- request-AC2 -> AC3. Proof: the spec examples validate in the test suite.

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
- Rationale: Every other item reads kinds or checks a schema.

# Notes
- Spec: `logics/specs/req_037_run_contract/`.
- A throwaway prototype of the validator passed on the spec files during grooming.
- Python because `bin/` and logics-manager are Python. A Go backend (req_031) ports the loader, the schemas stay.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
