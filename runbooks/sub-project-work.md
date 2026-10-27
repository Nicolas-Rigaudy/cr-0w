# Work on a registered project from the hub

Why this shape: `docs/sub-project-work.md`. A hub session never runs a repo's hooks or
permission rules, so repo work goes to a repo session.

## Decide
- Mostly about one repo (code beyond a small fix, design, MR/push, review threads): hand off.
- Lookup, status, a one-line fix: stay, after `crow work <project>`.
- Mixed: hand off the repo part, keep the hub part (tickets, mails, vault, briefing).

## Hand off
```
crow handoff <project> --task "<what>" --end-state "<what the owner's OK unlocks>" [--context "<hub facts>"]
```
- Opens a kitty window on `crow open <project>` with the note as first prompt (see `open-project-session.md`).
- The note is saved in `<state>/handoffs/<date>-<project>.md`; mention it in the activity log.
- End state examples: "after my OK: commit, push, post replies, resolve threads"; "commit and push to a branch when tests pass, then stop".
- Context: facts the hub already gathered (ticket text, MR number, decisions). Leave out steps the repo's skills already hold.
- Background instead of a window: `crow dispatch` with the same note as prompt file.

## Stay
```
crow work <project>        # or --json
```
Read the listed instructions and the skills that match the task. Treat listed denies as forbidden
and listed hooks as checks to run by hand (lint, gates). Git hooks still run.

## Guard
`bin/crow-hub-guard` runs on every hub prompt and Bash call (`.claude/settings.json`).
- A prompt naming a repo with a forge (name, folder or `aliases:`) gets a one-time note per session.
- Commit, push, rebase, merge, cherry-pick, `gh pr create|merge`, `glab mr create|merge` on a path or
  remote of such a repo, or `git worktree add` under /tmp: the owner is asked.
- Markers live in `<state>/hub-guard/<session>.json`; delete them to re-arm the notes.
- `CROW_HUB_GUARD=off` for hub runs nobody can answer (AFK orchestration only touches `.wt/` worktrees, which pass).

## After a handoff
- `crow handoff` already opens the kitty window (it calls `crow term`). Never say it "only saved the note" and never run `crow term` again with the same note: that starts a second session on the same task (two sessions racing on a prod release).
- If the owner says the session crashed or nothing opened, check `pgrep -af "crow open"` / `crow runs` first, then ask before reopening.
