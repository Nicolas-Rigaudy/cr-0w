## req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path - Clean crow releases: versioning, test pipeline, install and update path
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: release
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Give crow versioned releases people can install and update cleanly, with a test pipeline, modelled on how cdx-manager and logics release.
- Keywords: crow, release, version, changelog, CI, install, update, package, checksums
- Use when: planning or changing how crow is versioned, tested, packaged, installed or updated.
- Skip when: working on crow features or on the private `local/` config.

# Needs
- Anyone can install crow from a tagged release and update it with one command.
- Every release is gated by a test pipeline and has a changelog.
- The engine (everything outside `local/`) ships; the owner's `local/` config never does.

# Context
Today crow has no tags, no VERSION file, no CI workflow, no installer and no changelog. 50 commits on main, 21 test files, scripts in `bin/` (Python and shell), skills, runbooks, Logics docs. Users would clone and hope.

What cdx-manager does (to study and copy where it fits):
- `VERSION` file, `changelogs/`, `CONTRIBUTING.md`, `SECURITY.md`.
- `ci.yml`: OS x Python matrix, pytest with coverage, `npm pack --dry-run`, build check.
- `publish-pypi.yml` (PyPI trusted publishing, OIDC, no stored token) and `publish-npm.yml`.
- `upload-release-checksums.yml` plus `scripts/update_release_checksums.py`, `verify_release_checksums.py`, `verify_release_ci.py`: checksums are produced and verified in CI, not by hand (its issue #4 was a forgotten manual upload).
- `install.sh` / `install.ps1` and a `cdx update` command.

To study in logics: done 2026-10-29. Logics (`logics-manager`, v2.24.0) releases like this:
- Versions kept in step in `package.json`, `pyproject.toml` and a root `VERSION`.
- One curated note per release in `changelogs/CHANGELOGS_X_Y_Z.md`, folded into `CHANGELOG.md` by `npm run build:changelog`. CI fails on a stale aggregate.
- Release is a pushed `vX.Y.Z` tag. `release.yml` first waits for the green CI run on that exact commit (no duplicate test run), then creates the GitHub Release with the curated notes and publishes to npm, PyPI and the VS Code Marketplace in separate workflows, each retryable alone.
- A repo-owned contract, `logics/release/contract.json`, with `logics-manager release status|plan|validate|evidence`, so the release steps and their evidence are checkable. Worth reusing for crow since crow already runs on logics.
- CI gates beyond tests: artifact content check, npm audit, coverage floor, function-length and source-line budgets, README badge check.
- cdx-manager differs by also uploading and verifying checksums for release assets.

Ideas from the life-distributable skim (a personal agent shipped as a repo):
- One-sentence install: the user tells their agent to clone the repo and follow `SETUP.md`.
- `ADDENDA.md`: an append-only feed of offers; the owner's agent reads it weekly and raises one card per entry, applies nothing without a yes. A release note channel for agent-run installs.
- `make check` as the single gate (format, lint, tests, build).
- Opt-in, counts-only heartbeat telemetry, and agents offering fixes back as issues. Probably out of scope for crow.

# Decisions (owner, 2026-10-29)
- Distribution: GitHub Releases plus packages for anyone, not only the owner. AUR first (the owner uses it), then other package managers (Homebrew, a PyPI/pipx or npm package, .deb/.rpm if cheap). One release artifact feeds them all.
- Update path: `crow update`. Detects how crow was installed. Package-managed installs tell the user to use their package manager; script/git installs fetch the latest tag, show the changelog, never touch `local/`.
- Platforms: Linux and macOS in CI (a macOS runner tests what the owner cannot). Windows through WSL (it is Linux); native Windows only if it turns out cheap. The kitty coupling is removed from the core: terminal actions (focus, open, relaunch) go behind a small backend, kitty is one backend, others fall back to printing the command.
- Versioning: semver from 0.1.0. Crow proposes version, tag and changelog; the owner approves each release before it is tagged and published.

# Open questions
- Exact set of packages for 0.1.0 (suggest tag + GitHub Release + AUR, then Homebrew tap).
- Gate: one `bin/crow-check` reused locally and in CI (tests, commit-msg check, engine-vs-config leak check, lint).
- Which crow features are Linux-only (kitty, xdg-open, sandbox) and how the docs say so.
- Publishing a package or tag is outward-facing: each release needs the owner's OK on the exact version and notes.

# Execution order
1. item_030 crow-check (gate), item_031 terminal backend (independent, can run in parallel with 030).
2. item_032 CI (needs 030).
3. item_033 VERSION and changelogs, then item_034 engine/data split + installer + `crow update` (riskiest: do the split first and alone).
4. item_035 release workflow (needs 030, 032, 033).
5. item_036 packages (needs 034, 035), item_037 docs (runbook written along the way).
First tagged release 0.1.0 only after 1 to 5; the owner approves the version, tag and notes.

Execution is deliberate: plan here, execute each item in a fresh clear-context session started from its backlog doc (cheaper, no stale context).

# Design constraint found while planning
- `local/` (owner config, state, memory) lives inside the repo checkout and `CROW_HOME` points at the checkout. A package install is read-only, so the engine path and the data dir must be split (item_034). Default data dir `$XDG_CONFIG_HOME/crow`, current layout still honoured.

# Acceptance criteria
- AC1: A written comparison of the cdx-manager and logics release flows, with what crow adopts and what it skips.
- AC2: The open questions above are answered and recorded as decisions. The kitty-specific code paths are listed and put behind a terminal backend.
- AC3: A backlog with slices: gate script, CI workflow, versioning and changelog, install, update path, release workflow with checksums, docs (README install section, CONTRIBUTING, SECURITY).
- AC4: The first tagged release is cut by the pipeline, not by hand, and an install-then-update round trip is tested on a clean machine or container.
- AC5: Nothing under `local/` or any private name leaks into a release artifact (checked by the pipeline).

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- cdx-manager `.github/workflows/` and `scripts/` (upstream repo)
- https://github.com/mattabate/life-distributable (`ADDENDA.md`, `Makefile`, `SETUP.md`)
- `runbooks/readme-showcase.md`

# Backlog
- `logics/backlog/item_030_add_bin_crow_check_as_the_single_gate_shared_by_local_use_and_ci.md`
- `logics/backlog/item_031_put_terminal_actions_behind_a_backend_so_kitty_is_optional.md`
- `logics/backlog/item_032_add_ci_workflow_linux_and_macos_matrix_running_the_gate.md`
- `logics/backlog/item_033_add_version_changelogs_and_the_version_bump_check.md`
- `logics/backlog/item_034_add_install_script_and_the_crow_update_command.md`
- `logics/backlog/item_035_add_the_release_workflow_tag_await_ci_github_release_checksums.md`
- `logics/backlog/item_036_package_crow_aur_first_then_homebrew_tap.md`
- `logics/backlog/item_037_add_release_docs_readme_install_section_contributing_security_release_runbook.md`
