## item_036_package_crow_aur_first_then_homebrew_tap - Package crow: AUR first, then Homebrew tap
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
- Summary: Packages built from the release tarball: AUR first, then a Homebrew tap.
- Keywords: AUR PKGBUILD homebrew tap package
- Use when: Writing or updating a package recipe.
- Skip when: Release workflow.

# Problem
- Decision: GitHub Releases plus packages for anyone; AUR first (owner uses Arch), then other package managers.

# Scope
- In:
  - `packaging/aur/PKGBUILD` (source from the release tarball, sha256 from SHA256SUMS, `crow` on PATH, data dir per item_034, optional deps declared: kitty, tmux, cdx, logics-manager).
  - `packaging/homebrew/crow.rb` for a tap repo (works on macOS and Linux).
  - Updating a package for a new release is a documented step, scripted where possible (bump version and sha).
  - Further targets (pipx/PyPI, npm, deb/rpm) listed with a cheap/expensive note; pick only if wanted.
- Out:
  - Submitting to official distro repositories.
  - Windows package managers.

# Acceptance criteria
- AC1: PKGBUILD builds in a clean Arch container (`makepkg`) and `crow version` runs.
- AC2: The Homebrew formula passes `brew audit` and installs on the macOS runner in CI.
- AC3: Both recipes say `crow update` defers to the package manager.

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
- Rationale: Packages make crow installable by anyone, not only by the owner.

# Dependencies
- Needs items 034, 035. Publishing to AUR or a tap is outward-facing: the owner approves the exact package contents first.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
