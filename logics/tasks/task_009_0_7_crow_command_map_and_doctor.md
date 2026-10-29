## task_009_0_7_crow_command_map_and_doctor - 0.7 crow command map and doctor
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 90%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: `crow --help` in five groups (Start, See, Delegate, Rituals, Data), `crow projects` for `crow releases`, `crow runs sessions` and `crow runs relaunch`, and a read-only `crow doctor`. Old names stay as aliases with a one-line hint.
- Keywords: crow, cli, doctor, aliases
- Use when: changing the crow command map or adding a doctor check.
- Skip when: the change is outside the CLI.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_022_0_7_crow_command_map_and_doctor`

# Acceptance criteria
- AC1: `crow --help` shows the 5 groups; old commands still work.

# Plan
- [x] Help text in five groups.
- [x] `crow projects`, `crow runs sessions`, `crow runs relaunch`, `crow doctor`.
- [x] Old names print a hint on stderr and still run.
- [x] README, runbook, hub CLAUDE.md, widget default and messages use the new names.
- [ ] `crow watch` as an alias of `crow deck` (paused: the deck has no actions yet).
- [ ] `crow runs stop` (no stop exists in crow-runs).

# Validation
- `python3 -m pytest tests -q`: 125 passed, including tests/test_crow_cli.py (five groups, hints on stderr only, new names silent) and tests/test_crow_doctor.py (identity mismatch fails, merged worktree warns, default account, text and bad argument).
- `cd tools/crow-watch && go vet ./... && go test ./...` pass after the relaunch hint text changed.
- Live: `crow doctor` read the real registry and accounts; `crow releases`, `crow sessions` and `crow relaunch status` printed the hint and ran.
- Renders of `crow --help` and `crow doctor` under local/state/afk/shots/deck-final/png/.

# Report
- `crow doctor` reuses crow-registry for identity and crow-accounts for the policy. It fails on an identity mismatch, a missing cdx or a default account that crow may not target, and warns on merged worktrees. It fixes nothing; `crow tidy --apply` stays the cleanup.
- `crow runs` takes `sessions` and `relaunch` as subcommands. The scripts keep their names, so the widget, tests and old calls work.
- Not done: `crow watch` stays the older board, not an alias of `crow deck`, until the deck has the actions. `crow runs stop` does not exist.
- Not edited: runbooks/dispatch-background-runs.md and bin/crow-runs carry uncommitted owner changes in the main checkout and were left alone.

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
