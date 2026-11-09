# Tool update watch

Crow is built on a few fast-moving tools (cdx, logics-manager, Claude Code, Codex CLI, rtk).
Their release notes are read once per release, weekly, and turned into proposals for crow. Nothing is installed or edited.

## Pipeline
1. **Timer**: `crow-toolwatch.timer` (systemd user unit) runs `bin/crow-toolwatch run` Monday 09:30, `Persistent=true` (a missed Monday runs on wake).
2. **Check** (no model): for each tool in `local/tools.yaml`, the installed version (`version` command, first dotted number) and the changelog (`gh api`):
   - `github_releases: owner/repo`: published releases, drafts and prereleases skipped, paged until one at or below `last_read` shows up.
   - `github_file: owner/repo:path`: a Markdown changelog with one `## <version>` heading per release.
   - Unread = newer than `last_read` (else the installed version). Newest 12 kept per tool, commit and PR links stripped.
   - A tool not installed is reported with its latest version and its notes are skipped.
   - Nothing unread: logs "nothing new", exits 0, no model.
3. **Uses**: crow lines matching the tool's `uses` regex in `bin/`, `runbooks/`, `skills/`, `.claude/skills/`, `tools/crow-core/`, `integrations/`, `local/runbooks/`, grouped by subcommand and flags; plus open Logics items naming the tool.
4. **Analysis**: one `cdx run --kind oneshot --power low` (Sonnet, budget `CROW_TOOLWATCH_BUDGET`, default 1.5 USD) on `local/state/toolwatch/prompt-YYYY-MM-DD.md`. The answer is the note body.
5. **Note**: `Tool updates/Tool updates YYYY-MM-DD.md` in the vault crow folder: per tool, the changes that touch crow, proposals as `- [ ]` lines, what was ignored, then a deterministic "Installed vs latest" list.
6. **Bump**: `last_read` moves to the newest read version only after the note is written. A failed run keeps it, so the next run reads the same entries.

The brief (`toolwatch` source) shows the latest note while it has unchecked proposals. Tick a proposal in the note once handled or declined.

## Operate
- `bin/crow-toolwatch check [--json]`: installed, last read, latest and unread count per tool. Seconds, no model.
- `bin/crow-toolwatch run --dry-run`: writes the prompt file only.
- `bin/crow-toolwatch status`: last note, open proposals, log tail (`local/state/toolwatch/worker.log`).
- `bin/crow-toolwatch install-timer`: (re)writes and enables the units in the real home. Check with `systemctl --user list-timers crow-toolwatch.timer`.
- Add a tool: one entry in `local/tools.yaml` with `version`, `source`, optional `uses`, and `last_read` set to the installed version (else every past release is unread).
- Re-read a range: lower `last_read` and run.

## Pitfalls
- Repos that publish many prereleases (rtk) bury stable tags past the first page; paging stops at `last_read`, capped at 8 pages.
- Large release bodies (Codex) can time out on `gh api` (HTTP 504). The error shows in `check`, the tool is skipped, the others go on.
- Upgrades are the owner's call and live in the dotfiles repo. Proposals never install.
- An upgrade the owner OKs from a cdx session must run with the real HOME (`HOME=<real home> <real bin>/claude update`). Under the sandboxed HOME, `claude update` installs into the sandbox and the real install stays on the old version. Confirm with `<real bin>/claude --version`.
