## task_011_load_and_validate_mission_kind_files - Load and validate mission kind files
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-11-07 01:59:44

# AI Context
- Summary: Stdlib JSON Schema subset validator plus a repo test that keeps the req_037 spec files valid.
- Keywords: load, validate, mission, kind, files
- Use when: Editing the run contract spec files or schemas.
- Skip when: Reading kinds at runtime (the loader waits on the owner review).

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_104_load_and_validate_mission_kind_files`

# Acceptance criteria
- AC1: The validator accepts the 8 kind files and rejects a file with a bad phase, duplicate gates or an unknown field.
- AC2: An invalid file fails with the file path and the JSON path of the bad field.
- AC3: A test in the repo suite validates every file under the spec dir against its schema and runs in `bin/crow-check`.
- AC4: A schema using a keyword outside the subset fails the test.

# AC Traceability
- request-AC1 -> This task. Proof: the 8 kind files validate in tests/test_run_contract_spec.py.
- request-AC2 -> This task. Proof: the spec examples (crow blocks, questions, validation events) validate in the same test file.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_011_load_and_validate_mission_kind_files.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_011_load_and_validate_mission_kind_files.md` after implementation.

# Validation
- `python3 -m pytest -q tests/test_run_contract_spec.py`: 26 passed.
- `bin/crow-check`: green (366 passed; shellcheck and ruff not installed).
- Finish workflow executed on 2026-11-07.
- Linked backlog/request close verification passed.

# Report
- `tests/schema_lite.py` is the validator: it raises on any keyword outside the subset and reports `$.json.path: message`.
- `tests/test_run_contract_spec.py` validates the 8 kind files, every crow block in the examples, the questions example and the validation events, plus negative cases (bad phase, duplicate gates, unknown field, missing field, unknown keyword). It runs in the pytest group of `bin/crow-check`.
- AC1 to AC4 covered by those tests.
- Finished on 2026-11-07.
- Linked backlog item(s): `item_104_load_and_validate_mission_kind_files`
- Related request(s): `req_037_mission_kinds_and_run_contract`

# Links
- Request: `req_037_mission_kinds_and_run_contract`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
