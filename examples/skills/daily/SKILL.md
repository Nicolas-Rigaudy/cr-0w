---
name: daily
description: Draft your stand-up update from activity since the last one (commits, MRs, reviews, open threads). Use when the user says "daily" or "prep my daily".
context: fork
model: sonnet
effort: low
background: false
allowed-tools: Bash(bin/crow-gather:*)
---

# Daily (example)

Copy this folder to `local/hub-skills/daily/` and adapt it to your team's format, then link it:
`ln -s ../../local/hub-skills/daily .claude/skills/daily`. See `runbooks/own-routines.md`.

Runs in its own context and only drafts. Return the draft, never post it.

## Pre-gathered
One JSON from `bin/crow-gather daily`: `since` (day after `<state>/last-daily`), `activity` (work scope),
`notes`, `activity_logs`, `open_threads_work`, `tracker` ids and `targets`. Don't read the files again.
!`"${CLAUDE_PROJECT_DIR}/bin/crow-gather" daily --json 2>&1 || true`

## Gather
Optional: your tracker (Jira, Linear, GitHub issues) through its MCP server, from `since`.

## Draft
A short flat list. One line per topic, a few words each. No ticket keys, no commit hashes.
```
shipped: <what, in plain words>
in review: <what>
in progress: <what>
next: <what>

blockers
- <who or what is blocking, and since when>
```
Blockers are the most useful part. Real ones only: waiting on someone outside the team, access, a decision,
broken CI, something stuck for 2+ days. Leave the section out when there are none.

## Post
Show the draft and wait for an explicit OK or edits. Posting (chat, wiki, tracker) is your routine's last step:
put the target and the exact steps here. Then write today's date to `targets.last_daily`.
