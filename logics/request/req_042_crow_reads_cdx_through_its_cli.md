## req_042_crow_reads_cdx_through_its_cli - Finish moving crow onto the cdx CLI (crow #6)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub portability
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: crow still globs cdx profile internals (claude-home/.claude/projects, sessions/<pid>.json, cdx-run-<id>.stdout.log) where cdx 0.23+ gives the same data through its JSON CLI. Finish the switch and close GitHub issue #6.
- Keywords: cdx, profile_dir, transcripts_dir, cdx sessions, can-resume, conversation, issue 6, portability, codex
- Use when: changing how crow finds cdx transcripts, live sessions or run logs.
- Skip when: the read has no cdx CLI equivalent (rtk history, settings.json, setup token file, crow-doctor profile checks).

# Needs
- crow reads cdx profile paths through `cdx sessions --all --json` (`profile_dir`, `transcripts_dir`) and run logs through `cdx runs --json` artifacts, so a cdx layout change does not break crow silently and Codex profiles work too.
- crow relaunches an exact conversation with `cdx can-resume/resume --conversation <id>` on both providers.
- Live interactive sessions started through cdx come from `cdx sessions --json`; the /proc scan stays only for sessions started outside cdx.

# Context
- Issue #6 (opened by the cdx maintainer, 2026-10-27) asked what crow needs from cdx. cdx 0.23.0 shipped all five asks (AlexAgo83/cdx-manager#51).
- item_015 (req_005) posted the wishlist and adopted `runs --running` and `final_message`. The other three were left open.
- 2026-11-09: cdx 0.26.11 lets `--continue-run` carry a permissions-only `--settings`; crow dropped its restart-fresh workaround in ac1d396. Installed cdx: 0.26.11.

# Acceptance criteria
- AC1: no crow script outside a single helper in crow_lib builds `<profile>/claude-home/.claude/projects` or `.../log/cdx-run-*` paths by hand.
- AC2: relaunch of a conversation id uses `cdx can-resume --conversation` and works for a Claude and a Codex conversation.
- AC3: live sessions launched through cdx are read from `cdx sessions --json`; /proc remains a fallback for others.
- AC4: tests stay green with temp profiles (the helper is fakeable); issue #6 gets a closing comment the owner OKs.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- https://github.com/Nicolas-Rigaudy/cr-0w/issues/6
- `logics/backlog/item_015_cdx_cli_wishlist.md`
- `bin/crow_lib.py`, `bin/crow-usage`, `bin/crow-runs`, `bin/crow-relaunch`, `bin/crow-memory`, `bin/crow-cost-audit`

# Backlog
- `item_136_read_cdx_profiles_through_the_cli`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
