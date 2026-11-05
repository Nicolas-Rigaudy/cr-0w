# Interrupted sessions: find them, relaunch them on another account

When an account hits its 5-hour or weekly limit, the sessions running on it stop on a limit
message. They can continue on another allowed account through a cdx handoff.

## Before the stop: the quota warning
- A UserPromptSubmit hook (`crow-relaunch hook`, `cdx status --cached`) prints one QUOTA line when the current
  account passes 80% and again at 95% of its 5h or week window, with the reset time and the fallback.
  Each level shows once per account and window (`<state>/quota-warned.json`). At 95%: tell the owner and suggest /wrap.
- `crow doctor` warns when every `use` account is blocked or past 90% of its 5h window.

## Fresh stop: `crow rescue`
- `crow rescue <account> --dry-run` lists every session the account lost in its last 5h window: each one cut off by
  the limit and not relaunched yet, plus interactive sessions still open on it. `crow rescue <account>` relaunches all
  of them on one OK (`--to <target>` to choose).
- The target is the picker's choice (`crow runs relaunch status`): a `use` account first, then an account with
  `fallback: auto` in `local/accounts.yaml` once every `use` account is out. That one needs no OK while under its caps
  (`max_5h_pct`, `max_week_pct`; past the 5h cap it stays usable when its window resets within `free_5h_min`).
  Any other `ask-first` target needs a typed `yes`.
- `cdx rescue` moves only the latest conversation of an account and reads policies from `cdx set --policy`.
  Without cdx policies it treats every account as `use`, so prefer `crow rescue`.
- `crow runs sessions` and `crow runs relaunch` cover older stops and a single session.

## Find
- `crow runs sessions`: one line per session whose last turn is a limit error (retros and the
  /brief, /daily, /weekly one-shots are left out). Cheap to rerun: transcripts are cached by mtime and size.
- `crow runs sessions show <id-prefix>`: the full entry (transcript path, workspace, first and last prompt).
- Registry: `<state>/interrupted-sessions.json`; relaunches are logged in `<state>/relaunches.jsonl`
  and drop out of the list.

## Pick a target
- `crow runs relaunch status`: allowed accounts, which are blocked, and the recommended target
  (`use` accounts first, then `ask-first`). `never` accounts are not shown and are always refused.

## Relaunch
1. `crow runs relaunch <id> <account> --dry-run`: runs only the read-only cdx checks
   (`can-resume`, `status`) and prints the exact handoff. Nothing is prepared or launched.
2. `crow runs relaunch <id> <account>`: asks first. A `use` account takes `y` (or `--yes`); an `ask-first`
   account needs a typed `yes` on a terminal every time, and `--yes` does not count for it.
   From crow-deck: `R` on a paused Agents row opens a table of every account (policy, 5h and week
   headroom, why one is blocked), the recommended one preselected. Picking one and pressing enter is the
   owner's yes, so the deck passes `--picked`, which counts for both policies. Only an interactive
   picker may pass it. When every account is blocked, Home Next opens with "No tokens until HH:MM":
   work that needs no model (red pipelines, reviews, threads, dirty repos, stale worktrees, paused runs).
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
