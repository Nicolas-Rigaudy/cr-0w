## item_032_add_ci_workflow_linux_and_macos_matrix_running_the_gate - Add CI workflow: Linux and macOS matrix running the gate
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 80%
> Progress: 100%
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

# Validation
- `.github/workflows/ci.yml`: push to main and pull_request, ubuntu-latest and macos-latest, actions pinned by commit SHA, runs `bin/crow-check` then `go vet` and `go test` for both Go tools. No weekly audit job: the only dependencies are pytest, pyyaml and markdown, installed unpinned.
- AC1: run 41833856752 on main is green on both operating systems (a PR run was not opened; the push run uses the same job).
- AC2: the first five runs were red on a clean checkout, each for a real defect, fixed in order: tests needing `local/projects.yaml`, a shellcheck warning in `crow-prepush-check`, bash 3.2 empty array in `crow term`, BSD sed in `crow-commit-msg-check`, `crow-runs` needing `cdx` and `/proc`. No commit was broken on purpose.
- AC3: the checkout has no `local/` and the workflow sets no secrets.
- The CI runs found macOS and cdx-less defects the local run hid, because the local shell had CROW_HOME and cdx set. Local check without them: `env -u CROW_HOME -u CROW_ACCOUNT python3 -m pytest -q tests` with cdx off PATH.

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
