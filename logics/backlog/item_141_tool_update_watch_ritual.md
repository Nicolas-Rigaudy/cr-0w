## item_141_tool_update_watch_ritual - Tool update watch ritual: changelogs since last check, mapped to crow uses, proposals in the vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
> Complexity: Small
> Theme: Hub upkeep
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:28:44

# AI Context
- Summary: `bin/crow-toolwatch` with a weekly systemd user timer, built like crow-retro: a no-model pre-filter compares installed and latest versions against the last ones read, then one small-model run turns the new changelog entries into a vault note of proposals.
- Keywords: crow-toolwatch, timer, changelog, versions, proposals, vault note
- Use when: building req_045.
- Skip when: upgrading packages (the owner's dotfiles repo).

# Problem
- New features in cdx and logics-manager land weekly. crow keeps workarounds that a release already made unnecessary until the owner asks by hand.

# Scope
- In:
  - `local/tools.yaml`: per tool, the version command, the changelog source (npm package CHANGELOG.md, GitHub releases through `gh`, `claude --version` plus its changelog) and `last_read`. Start with cdx, logics-manager, Claude Code, Codex CLI, rtk.
  - `bin/crow-toolwatch check`: no model. Prints per tool the installed version, the latest version and the unread entries. It exits 0 with nothing to do when nothing is new.
  - `bin/crow-toolwatch run`: when check finds unread entries, one `crow dispatch` run on the smallest fitting model. It gets the entries and a grep of crow's uses of each tool (bin, runbooks, tools/crow-core, open Logics items). It writes `Tool updates YYYY-MM-DD.md` in the vault crow folder: what changed, what it touches in crow (file or item), proposals, and an "installed behind latest" line. It then bumps `last_read`.
  - `install-timer` like crow-retro: weekly, Monday morning, `Persistent=true`.
  - The brief lists the note while it has open proposals.
  - runbooks/tool-update-watch.md and an index line.
- Out:
  - Installing or upgrading anything (the owner's dotfiles repo).
  - Editing crow code from the run. Proposals only.

# Acceptance criteria
- AC1: `crow-toolwatch check` with no new release prints the versions and starts no model run (test with a fake tools.yaml and fake version commands).
- AC2: with a new fake release, `run` builds a prompt holding only the unread entries and the matching crow uses, and bumps `last_read` only after the note is written.
- AC3: the timer is installed with `crow-toolwatch install-timer` and shows in `systemctl --user list-timers`.
- AC4: the brief shows the latest note while it has open proposals.
- AC5: the runbook exists and is indexed.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: `tests/test_crow_toolwatch.py::test_nothing_new_starts_no_model` (fake tools.yaml, `echo` version commands, model stub that fails if called); live `crow-toolwatch run` after the first note prints "nothing new".
- request-AC2 -> This backlog slice. Proof: `test_run_prompt_holds_unread_and_uses_then_bumps` (prompt holds only unread entries and crow uses, config unchanged during the model call, bumped after the note) and `test_failed_run_keeps_last_read`.
- request-AC3 -> This backlog slice. Proof: `crow-toolwatch install-timer`; `systemctl --user list-timers` shows crow-toolwatch.timer, next Mon 2026-11-16 09:39.
- request-AC4 -> This backlog slice. Proof: `crow-gather` source `toolwatch` (latest note while it has `- [ ]` lines) plus brief SKILL step 7 and output line; live gather lists the 2026-11-09 note with 9 open proposals.
- request-AC5 -> This backlog slice. Proof: `runbooks/tool-update-watch.md`, indexed in `runbooks/README.md`.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_045_periodic_tool_update_watch`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Saves a manual pass each week. Cheap to build on the crow-retro pattern.

# Notes
- 2026-11-09: first manual pass done by hand (cdx 0.26.0..0.26.11, logics-manager 2.24.x). Seed `last_read` with cdx 0.26.11 and logics-manager 2.24.3.

- 2026-11-09: built in the hub session. Sources: GitHub releases (paged past prereleases) or a `## <version>` changelog file, both via `gh api`. Notes land in `Tool updates/` under the vault crow folder (new `toolwatch` note kind). Run = `cdx run --kind oneshot --power low`, Sonnet, 1.5 USD budget. First note read claude 2.1.321 and rtk 0.45.0..1.2.0 (17 older rtk releases left out by the 12-entry cap).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
