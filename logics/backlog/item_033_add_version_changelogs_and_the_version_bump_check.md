## item_033_add_version_changelogs_and_the_version_bump_check - Add VERSION, changelogs and the version-bump check
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
- Summary: VERSION file, changelogs and a check that versions and changelog agree.
- Keywords: version semver changelog VERSION bump
- Use when: Bumping or validating a version.
- Skip when: Packaging.

# Problem
- No VERSION, no changelog, no tag. Decision: semver from 0.1.0; crow proposes version, tag and notes, the owner approves each release.

# Scope
- In:
  - `VERSION` at the root is the single source; any other file that carries the version (package manifest if any) is checked against it by `bin/crow-check`.
  - `changelogs/CHANGELOGS_X_Y_Z.md` per release (what changed, in plain words), aggregated into `CHANGELOG.md` by a script; check fails on a stale aggregate (logics-manager pattern).
  - A `crow release propose` helper (or runbook section) that reads commits since the last tag, groups by the `new/chg/fix` subject type, and drafts the next version and notes for the owner to approve.
- Out:
  - Tagging and publishing (item_035).

# Acceptance criteria
- AC1: A mismatch between VERSION and another version carrier fails the gate.
- AC2: A stale or missing changelog entry for the VERSION fails the gate.
- AC3: The draft for the next release is produced from commits and needs no network.
- AC4: First entry written for 0.1.0.

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
- Rationale: The release workflow and every package read the version from here.

# Dependencies
- Feeds item_030 (check) and item_035.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
