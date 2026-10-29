## item_034_add_install_script_and_the_crow_update_command - Add install script and the crow update command
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
- Summary: Install script and a crow update command that respects how crow was installed.
- Keywords: install update crow update installer CROW_HOME local
- Use when: Changing install or update behaviour.
- Skip when: Package recipes (item_036).

# Problem
- Decision: `crow update` is the update command. It detects the install method: a package-managed install tells the user to use their package manager; a git/script install fetches the latest tag, shows the changelog and updates.
- Design constraint found in the code: `local/` (owner config, state, memory) lives inside the repo checkout and `CROW_HOME` points at the checkout (see `bin/crow`, `bin/crow-precommit-check`). A package installs the engine read-only, so config cannot live there.

# Scope
- In:
  - Split engine path from data path: the engine is found from the script location, the owner data dir is `CROW_HOME` (default `$XDG_CONFIG_HOME/crow`, falling back to the checkout for the owner's current layout). Migration note for the current layout, no move forced.
  - `install.sh` (Linux, macOS, WSL): installs a tagged release to a user prefix, verifies the release checksum, creates the data dir from `local.example` files, never overwrites existing data.
  - `crow update [--check] [--to VERSION]`: detects install method (package / script / git checkout), shows the changelog between versions, updates only the engine, refuses when the tree is dirty, never touches the data dir.
  - `crow version`.
- Out:
  - Windows native installer (WSL is the Windows path).
  - Auto-update in the background.

# Acceptance criteria
- AC1: Fresh install on a clean container puts crow on PATH and `crow version` works.
- AC2: `crow update` on each install method: script/git updates, package says to use the package manager, `--check` changes nothing.
- AC3: A data dir with content survives an update byte for byte (tested).
- AC4: Checksum mismatch aborts the install.

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
- Rationale: The core user-facing part of the plan.

# Dependencies
- Needs item_033 (VERSION) and item_035 (release assets with checksums). The engine/data split is the riskiest step: do it first and alone.

# Notes
- Planned 2026-10-29. Execution happens in a fresh, clear-context session started from this doc and the request; do not execute from the planning session.
