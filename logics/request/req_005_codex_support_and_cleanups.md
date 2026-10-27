## req_005_codex_support_and_cleanups - Claude and Codex through cdx: fix the six open issues
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub portability
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Crow assumes Claude and one profiles path. Make it work with Codex accounts and a relocated cdx home, using the cdx CLI where it has the data.
- Keywords: claude, codex, through, cdx, fix, six, open, issues
- Use when: changing profile lookup, live agent and run counting, run reports, or hub instructions for a non-Claude account.
- Skip when: working on UI widgets or usage pricing (separate requests).

# Needs
- Crow runs for an owner who uses Codex accounts, a non-default cdx home, or a non-Arch Linux box (GitHub issues #1 to #6).
- The engine stays generic and small: one shared helper, cdx JSON where it has the data, fixtures for every parser.

# Context
- Fork review and decisions: `local/state/drafts/codex-support-fork-review.md`. The fork is inspiration only.
- cdx 0.22 `cdx runs --json` rows carry `status`, `kind`, `pid`, `provider`, `ended_at`.
- Out of scope: price table, Codex rollout readers, cross-provider relaunch, widget badges.

# Acceptance criteria
- AC1: One helper resolves the profiles dir (`CDX_PROFILES`, then `CDX_HOME/profiles`, then `~/.cdx/profiles`); no script reads those variables itself (#4).
- AC2: Background runs are counted from `cdx runs --json`; /proc is used for interactive sessions only and matches `claude` and `codex` with `-p` and `--print` (#2).
- AC3: README lists Linux requirements; tests pass without `local/accounts.yaml` and without a `/home/<user>` assumption (#3).
- AC4: Run reports parse Claude and Codex output, covered by fixtures (#5).
- AC5: A Codex account gets the hub instructions (read-first list), a provider-aware model and prompt, and the guard coverage is documented (#1).
- AC6: The list of what crow still needs from the cdx CLI is written as a reply for the issue (#6).

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- https://github.com/AlexAgo83/cr-0w/issues (#1 to #6)
- `local/state/drafts/codex-support-fork-review.md`

# Backlog
- `item_010_profiles_dir_helper`
- `item_011_background_runs_from_cdx`
- `item_012_linux_docs_and_portable_tests`
- `item_013_run_reports_per_provider`
- `item_014_codex_hub_instructions`
- `item_015_cdx_cli_wishlist`
