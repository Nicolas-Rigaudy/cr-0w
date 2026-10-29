## item_030_add_bin_crow_check_as_the_single_gate_shared_by_local_use_and_ci - Add bin/crow-check as the single gate, shared by local use and CI
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Medium
> Theme: release
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: One script, bin/crow-check, runs every release gate; developers and CI call the same thing.
- Keywords: check gate lint tests leak commit-msg CI
- Use when: Changing or running the release gate; wiring CI.
- Skip when: Feature work.

# Problem
- Today the checks are scattered: pytest (21 test files in tests/), `bin/crow-commit-msg-check`, `bin/crow-precommit-check`, `bin/crow-prepush-check`, `bin/crow-review-check`.
- No single command says 'this tree is releasable'.

# Scope
- In:
  - `bin/crow-check`: runs pytest, shell lint (shellcheck if present, skipped with a note otherwise), python lint (ruff), the engine-vs-config leak check (no private name or personal path outside `local/`), commit-subject check over the commits since the last tag, version consistency (item_033).
  - Exit non-zero on the first red group but print every group's result; `--fast` skips slow groups; `--json` for CI summaries.
  - A fixture-based test per check.
- Out:
  - Release publishing (item_035).
  - Per-project repo checks (those stay in the repos).

# Acceptance criteria
- AC1: `bin/crow-check` exits 0 on a clean tree and non-zero when any one group fails (each group proven by a failing fixture).
- AC2: The leak check fails on a planted private name outside `local/` and ignores `local/`.
- AC3: Runs on Linux and macOS without root and without `local/` present (CI has none).
- AC4: Documented in CONTRIBUTING (item_037).

# AC Traceability
- request-AC3 -> This backlog slice. Proof: bounded slice with testable criteria above.
- request-AC4 -> The slice's criteria feed the first tagged release and the install/update round trip.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: First slice: CI (item_032) and the release workflow (item_035) call it.

# Dependencies
- None. Blocks item_032, item_035.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
