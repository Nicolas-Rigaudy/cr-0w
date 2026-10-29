## item_037_add_release_docs_readme_install_section_contributing_security_release_runbook - Add release docs: README install section, CONTRIBUTING, SECURITY, release runbook
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
- Summary: Release documentation: README install section, CONTRIBUTING, SECURITY, release runbook, feature-by-platform table.
- Keywords: docs README CONTRIBUTING SECURITY runbook release
- Use when: Writing or checking release docs.
- Skip when: Code changes.

# Problem
- No CONTRIBUTING.md, SECURITY.md or release runbook. Model: cdx-manager and logics-manager both ship CONTRIBUTING.md, SECURITY.md, docs/release.md.

# Scope
- In:
  - README: install per method (package, script, git), update, supported platforms (Linux, macOS, Windows via WSL) and the feature-by-platform table from item_031.
  - `CONTRIBUTING.md` (run `bin/crow-check`, commit subject format `type[SCOPE]: message`, tests).
  - `SECURITY.md` (how to report, what crow touches: tokens, local state; no secrets in the repo).
  - `runbooks/release.md`: the exact steps to cut a release with approvals, indexed in `runbooks/README.md`.
  - Everything generic: no private names or personal paths (the leak check enforces it).
- Out:
  - Marketing and portfolio pages (req_004).

# Acceptance criteria
- AC1: A reader on a clean machine can install, update and cut a release from the docs alone.
- AC2: The release runbook is indexed and was followed for the first release.
- AC3: The leak check passes on all new docs.

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
- Priority: Medium
- Rationale: Last slice, but the runbook is written as each step is learned.

# Dependencies
- Runbook written alongside items 033 to 036.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
