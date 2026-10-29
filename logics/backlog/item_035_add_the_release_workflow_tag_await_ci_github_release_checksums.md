## item_035_add_the_release_workflow_tag_await_ci_github_release_checksums - Add the release workflow: tag, await CI, GitHub Release, checksums
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
- Summary: Tag-triggered release workflow: wait for green CI, GitHub Release with notes and checksummed assets.
- Keywords: release workflow tag github release checksums await-ci
- Use when: Cutting a release or changing the release workflow.
- Skip when: Package recipes.

# Problem
- Model: logics-manager `release.yml` (waits for green CI on the tagged commit, then GitHub Release with curated notes; retryable separate publish workflows) and cdx-manager checksum upload and verification (its issue #4 was a forgotten manual upload).
- logics already has `logics-manager release status|plan|validate|evidence` with a repo contract `logics/release/contract.json`.

# Scope
- In:
  - `.github/workflows/release.yml` on `v*.*.*` tags: await CI on the commit, build the release tarball, write `SHA256SUMS`, create the GitHub Release with the changelog entry as notes and the assets attached.
  - A step that verifies the checksums of the uploaded assets after upload.
  - `logics/release/contract.json` created with `logics-manager release discover`, so `release plan/validate` work for crow.
  - Owner approval gate: tags are pushed by the owner only (push on crow is personal-scope but publishing is outward-facing); the helper in item_033 only proposes.
  - Nothing from `local/` can enter the tarball (the check from item_030 runs on the built artifact too).
- Out:
  - Package-manager publishing (item_036).

# Acceptance criteria
- AC1: Pushing a test tag on a scratch branch or fork produces a Release with notes, tarball and SHA256SUMS.
- AC2: A red CI on the commit blocks the release.
- AC3: The tarball contains no `local/` and no file matching the deny list (verified in CI).
- AC4: `logics-manager release validate 0.1.0` passes.

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
- Rationale: Makes releases a pipeline action, not a manual one.

# Dependencies
- Needs items 030, 032, 033.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
