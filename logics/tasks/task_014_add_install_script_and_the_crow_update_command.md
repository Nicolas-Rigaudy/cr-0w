## task_014_add_install_script_and_the_crow_update_command - Add install script and the crow update command
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 30%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-07 12:03:15

# AI Context
- Summary: Split engine path from data dir, add install.sh, crow update, crow version, restore and the doctor data-repo check.
- Keywords: add, install, script, crow, update, command
- Use when: Changing how crow is installed, updated or where it reads owner data.
- Skip when: Package recipes (item_036) or the release workflow (item_035).

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_034_add_install_script_and_the_crow_update_command`

# Acceptance criteria
- AC1: Fresh install on a clean container puts crow on PATH and `crow version` works.
- AC2: `crow update` on each install method: script/git updates, package says to use the package manager, `--check` changes nothing.
- AC3: A data dir with content survives an update byte for byte (tested).
- AC4: Checksum mismatch aborts the install.
- AC5: On a clean container, `install.sh --restore <url>` brings back the owner's data dir and prints the token and account steps left; the runbook covers them.
- AC6: `crow doctor` flags a data repo with no remote, a dirty tree or unpushed commits.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_014_add_install_script_and_the_crow_update_command.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_014_add_install_script_and_the_crow_update_command.md` after implementation.

# Validation
- (no validation recorded yet)

# Report
- Slice 1 (split): `CROW_ENGINE` is the script location; `CROW_HOME` is the data root holding `local/`: the variable, else the checkout when it has `local/`, else `$XDG_CONFIG_HOME/crow`. Python scripts resolve it through `crow_lib.data_home()`; engine paths (`bin/`, `tools/`) come from the script location. `crow version` prints version, engine and data paths and needs no account. Migration note in `runbooks/install-and-update.md`; nothing moves by itself.
- Decision taken during the run: the data root keeps a `local/` subdir (`~/.config/crow/local/...`), so the checkout layout and the installed layout have the same shape and every `HOME/local/...` lookup stays as is. Flattening it is a later choice.
- Known gap: a hub session still opens in `CROW_HOME`; on an installed engine that dir has no CLAUDE.md. Follow-up for the docs/packaging items.

# Links
- Request: `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
