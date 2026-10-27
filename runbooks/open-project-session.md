# Open an interactive project session from the hub

The hub runs inside Claude Code with no terminal, so it cannot start a TUI itself.
Use `crow term <project>`: it opens a new kitty window running `crow open <project>`.

## Why not just `kitty crow open ...`
A process spawned from a Claude session inherits two things that break the new window:
- **Sandboxed HOME** (cdx remaps it): kitty doesn't find `~/.config/kitty`, so the window has no theme and looks all black.
- **Claude Code session markers** (`CLAUDECODE`, `CLAUDE_CODE_CHILD_SESSION`, `CLAUDE_CODE_SESSION_ID`, ...): the new Claude treats itself as a child session and shows
  "Transcript saving is off — inherited CLAUDE_CODE_CHILD_SESSION marker" (no resume, no cdx usage backfill).

`crow term` restores the real HOME (from `getent passwd`) and unsets `CLAUDE*`, `ANTHROPIC_*` and `CDX_*`.

## Starting prompt
`crow handoff <project> --task ... --end-state ...` builds this prompt from the registry (see `sub-project-work.md`).
`crow handoff` writes the note and opens the kitty window itself (it execs `crow term`). Never run `crow term` again with the same note: that starts a second session on the same task. Check `crow runs` / `ListAgents` for an existing session first.
`crow open|term <project> "<prompt>"` passes a first message (cdx ≥ 0.21.4, `--prompt`).

Writing the handoff prompt:
- Name the project skill to use ("use the mr-review skill") instead of restating its steps; the skill holds the current rules.
- State the end state the user's approval unlocks (e.g. "once I approve an MR: commit, push, post replies, resolve threads").
  A bare "show me before committing" made a session stop after the commit and wait for a push it could have done.
- Keep the task list; drop instructions the skill already covers.
- Start with "Read local/owner-prefs.md in the hub" (absolute path from `local/context.md`): repo sessions don't load hub memory.
- Name the worktree to use (`<repo>.wt/<topic>`), never /tmp.

## Checks before `crow handoff`
- Pass the registry name from `crow ls` (`<project>`), not an alias: `crow term` rejects aliases after the note is already written.
- Pass `--task` as a literal string. Never build it with `sed`/command substitution: an empty match saves a note with `Task:` blank (2026-10-27, the session only recovered because it searched the repo). `grep -n '^Task: .' <note>` right after.
