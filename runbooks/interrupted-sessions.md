# Interrupted sessions: find them, relaunch them on another account

When an account hits its 5-hour or weekly limit, the sessions running on it stop on a limit
message. They can continue on another allowed account through a cdx handoff.

## Fresh stop: use cdx
- A stop in the last 24 hours: `cdx rescue --dry-run`, then `cdx rescue` (it picks a ready account of the same
  provider, skips what it already rescued, and respects `cdx set --policy`). Policy lives in cdx; `local/accounts.yaml`
  only supplies the fallback and the memory mode.
- `crow sessions` and `crow relaunch` cover what cdx does not: older stops, a chosen target, memory linking.

## Find
- `crow sessions`: one line per session whose last turn is a limit error (retros and the
  /brief, /daily, /weekly one-shots are left out). Cheap to rerun: transcripts are cached by mtime and size.
- `crow sessions show <id-prefix>`: the full entry (transcript path, workspace, first and last prompt).
- Registry: `<state>/interrupted-sessions.json`; relaunches are logged in `<state>/relaunches.jsonl`
  and drop out of the list.

## Pick a target
- `crow relaunch status`: allowed accounts, which are blocked, and the recommended target
  (`use` accounts first, then `ask-first`). `never` accounts are not shown and are always refused.

## Relaunch
1. `crow relaunch <id> <account> --dry-run`: runs only the read-only cdx checks
   (`can-resume`, `status`) and prints the exact handoff. Nothing is prepared or launched.
2. `crow relaunch <id> <account>`: asks first. A `use` account takes `y` (or `--yes`); an `ask-first`
   account needs a typed `yes` on a terminal every time, and `--yes` does not count for it.
3. After the confirmation it opens a new kitty window (detached, real HOME, no inherited Claude/cdx
   variables) that runs `cdx handoff <source> <target> --source-conversation <id>` from the session's
   workspace. `--here` runs it in the current terminal instead.
   Never call `cdx handoff` without `--source-conversation`: cdx then offers its own candidates, which
   once included an automatic retro instead of the real session.

## Gotchas
- A session that changed directory (worktrees, subfolders) can be handed off since cdx 0.22: the relaunch
  starts in the last folder the conversation worked in that still exists.
- Check `crow runs` before relaunching a background run: a later run may already have resumed it.
- Accounts with memory mode `session` get `local/memory` linked for the relaunched session, then
  unlinked and committed when it ends.
