## task_012_parse_and_validate_the_crow_block_and_questions_file - Parse and validate the crow block and questions file
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: Python parser in bin/ for the crow block and questions.json, built on the stdlib validator.
- Keywords: parse, validate, crow, block, questions, file
- Use when: Reading any run output that follows the contract.
- Skip when: Acting on the block (publish, versions).

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_106_parse_and_validate_the_crow_block_and_questions_file`

# Acceptance criteria
- AC1: Every spec example parses and validates.
- AC2: A ready report without proof, a v2 without responses, a file with two blocks and a name that does not match the version each fail with a clear message.
- AC3: `outbox/questions.json` from the spec validates; a question without `blocking` fails.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_012_parse_and_validate_the_crow_block_and_questions_file.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_012_parse_and_validate_the_crow_block_and_questions_file.md` after implementation.

# Validation
- `python3 -m pytest tests/test_crow_contract.py tests/test_run_contract_spec.py`: 49 passed.
- `bin/crow-check`: pytest 397 passed; ruff red on files this task did not touch (crow_lib, crow_term, conftest and others).
- `bin/crow-contract <spec example>` prints ok and exits 0; a file with no block exits 1.
- Finish workflow executed on 2026-11-07.
- Linked backlog/request close verification passed.

# Report
- `bin/crow_contract.py` extracts exactly one fenced crow block (YAML or JSON), validates it and `outbox/questions.json`, and checks file name against phase and version. `bin/crow-contract FILE...` is the CLI. The validator moved to `bin/crow_schema_lite.py`; `tests/schema_lite.py` re-exports it.
- AC1: every spec example parses (test_spec_examples_parse). AC2: tests for missing proof, v2 without responses, two blocks, name/version mismatch. AC3: spec questions.json validates, missing blocking fails.
- Finished on 2026-11-07.
- Linked backlog item(s): `item_106_parse_and_validate_the_crow_block_and_questions_file`
- Related request(s): `req_037_mission_kinds_and_run_contract`

# Links
- Request: `req_037_mission_kinds_and_run_contract`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
