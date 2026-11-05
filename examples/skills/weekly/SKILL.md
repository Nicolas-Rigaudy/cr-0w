---
name: weekly
description: Weekly rollup. Turns the week's activity into a summary file and a short outline for a demo or report. Use when the user says "weekly" or "rollup".
context: fork
model: sonnet
effort: medium
background: false
allowed-tools: Bash(bin/crow-gather:*)
---

# Weekly (example)

Copy this folder to `local/hub-skills/weekly/` and adapt it, then link it:
`ln -s ../../local/hub-skills/weekly .claude/skills/weekly`. See `runbooks/own-routines.md`.

Runs in its own context. Write the rollup file, then return a short summary and the outline.

## Pre-gathered
One JSON from `bin/crow-gather weekly`: `monday`, `iso_week`, `activity` (all scopes), `notes`,
`activity_logs`, `decisions`, `carry_over` (open threads) and `targets`. Don't read the files again.
!`"${CLAUDE_PROJECT_DIR}/bin/crow-gather" weekly --json 2>&1 || true`

## Gather
Optional: your tracker through its MCP server, from `monday`.

## Rollup to `targets.rollup`
```
# Week YYYY-Www
## Shipped
## Reviews
## Decisions
## Carry-over     (threads still open)
## Outline        (one achievement per slide: outcome as title, 2-3 bullets of impact)
```
Then delete this week's `activity/*.md` files (git history keeps them).
