## item_032_add_ci_workflow_linux_and_macos_matrix_running_the_gate - Add CI workflow: Linux and macOS matrix running the gate
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
- Summary: GitHub Actions CI: Linux and macOS matrix running bin/crow-check.
- Keywords: CI github actions matrix linux macos
- Use when: Editing or debugging the CI workflow.
- Skip when: Release publishing.

# Problem
- crow has no .github/workflows. The owner has no Mac, so a macOS runner is the only macOS test.
- Model: cdx-manager `ci.yml` (OS x language matrix) and logics-manager `ci.yml`.

# Scope
- In:
  - `.github/workflows/ci.yml`: push to main and pull_request; ubuntu-latest and macos-latest; install deps; run `bin/crow-check`.
  - WSL is Ubuntu, so the Linux job is its proof; no Windows job.
  - A weekly dependency/audit job like logics `audit.yml` if crow has any dependency to audit.
  - Pin action versions.
- Out:
  - Release publishing (item_035).

# Acceptance criteria
- AC1: CI is green on a PR for both operating systems.
- AC2: A deliberately broken commit turns CI red (shown once, then reverted).
- AC3: The job runs without `local/` and without secrets.

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
- Rationale: Gate for every merge and for the release workflow's await-ci step.

# Dependencies
- Needs item_030. Macos green depends on item_031 for the kitty paths.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
