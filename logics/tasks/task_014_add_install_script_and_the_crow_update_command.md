## task_014_add_install_script_and_the_crow_update_command - Add install script and the crow update command
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
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
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

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
- [x] Use `python3 -m logics_manager flow progress task task_014_add_install_script_and_the_crow_update_command.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_014_add_install_script_and_the_crow_update_command.md` after implementation.

# Validation
- AC1: `tests/test_install_update.py::test_install_puts_crow_on_path_with_a_data_dir`; clean-container smoke `tests/smoke/install-in-docker.sh` (python:3.13-alpine) green, log in the run's shots dir (`install-smoke.log`).
- AC2: `test_script_update_check_then_update_keeps_data`, `test_script_update_to_a_given_version`, `test_package_install_points_to_the_package_manager`, `test_git_update_check_changes_nothing_then_updates`, `test_git_update_refuses_a_dirty_tree`; script update also run in the container smoke.
- AC3: `test_script_update_check_then_update_keeps_data` and `test_git_update_check_changes_nothing_then_updates` compare a digest of the data tree (binary file included) before and after.
- AC4: `test_checksum_mismatch_aborts` (nothing under the prefix, no bin link); container smoke checks the same.
- AC5: `test_restore_clones_the_data_repo_and_lists_steps` and the container restore step; steps covered in `runbooks/install-and-update.md`.
- AC6: `tests/test_crow_doctor.py::test_doctor_warns_on_a_data_dir_without_git`, `test_doctor_flags_the_data_repo_gaps` (no remote, uncommitted, clean and pushed, 1 commit not pushed).
- Full suite: 433 passed. `bin/crow-check`: pytest, shellcheck, leak, commits, version pass; ruff red on findings present before this task (for example FURB167 in `bin/crow_lib.py:90`), none in the new files.
- Finish workflow executed on 2026-11-07.
- Linked backlog/request close verification passed.

# Report
- Slice 1 (split): `CROW_ENGINE` is the script location; `CROW_HOME` is the data root holding `local/`: the variable, else the checkout when it has `local/`, else `$XDG_CONFIG_HOME/crow`. Python scripts resolve it through `crow_lib.data_home()`; engine paths (`bin/`, `tools/`) come from the script location. `crow version` prints version, engine and data paths and needs no account. Migration note in `runbooks/install-and-update.md`; nothing moves by itself.
- Decision taken during the run: the data root keeps a `local/` subdir (`~/.config/crow/local/...`), so the checkout layout and the installed layout have the same shape and every `HOME/local/...` lookup stays as is. Flattening it is a later choice.
- Slice 2 (install, update): `install.sh` (POSIX sh) and `bin/crow-update`, commit `6ee695e`. Release contract item_035 must produce: `crow-<ver>.tar.gz` with top dir `crow-<ver>/` and no `local/`, plus `SHA256SUMS`, under `releases/download/v<ver>/`.
- Slice 3 (doctor): `crow doctor` area `data` warns when `<CROW_HOME>/local` is not a git repo, has no remote, uncommitted changes, no upstream or unpushed commits.
- Known gaps: the shell hooks (commit-msg, pre-push, pre-commit checks) still read `local/` next to the checkout; script installs are replaced whole, so local edits in an installed engine are lost on update (expected for that method).
- Known gap: a hub session still opens in `CROW_HOME`; on an installed engine that dir has no CLAUDE.md. Follow-up for the docs/packaging items.
- Finished on 2026-11-07.
- Linked backlog item(s): `item_034_add_install_script_and_the_crow_update_command`
- Related request(s): `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`

# Links
- Request: `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
