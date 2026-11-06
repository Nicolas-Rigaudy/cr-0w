## item_106_parse_and_validate_the_crow_block_and_questions_file - Parse and validate the crow block and questions file
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
- Summary: Extract the one fenced crow block from a plan, report or mockup file, validate it and `outbox/questions.json` against the spec schemas, and report clear errors.
- Keywords: crow block, parser, questions.json, validation, contract
- Use when: Reading any run output that follows the contract.
- Skip when: Acting on the block (publish, versions).

# Problem
- Free text in a run's last message is hard to parse (Ultron lesson 2).
- The spec fixes the block and the questions file; nothing reads them yet.

# Scope
- In:
  - Find exactly one fenced `crow` block; zero or two is an error.
  - Parse YAML or JSON, validate against `crow-block.schema.json` with the item_104 validator.
  - Validate `outbox/questions.json` against `questions.schema.json`.
  - Check the file name matches `phase` and `version`.
  - Tests with the spec examples and broken copies (no proof on a ready report, no responses from v2, two blocks).
- Out:
  - Repairing a broken block by asking the run again (stop guard item).
  - Publishing.

# Acceptance criteria
- AC1: Every spec example parses and validates.
- AC2: A ready report without proof, a v2 without responses, a file with two blocks and a name that does not match the version each fail with a clear message.
- AC3: `outbox/questions.json` from the spec validates; a question without `blocking` fails.

# AC Traceability
- request-AC2 -> AC1. Proof: spec examples validate through the parser in the test suite.

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
- Rationale: The contract is useless until something reads it.

# Notes
- Depends on item_104 (validator).
- Spec: `logics/specs/req_037_run_contract/`.
- Usable by AFK run reports before core exists.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
