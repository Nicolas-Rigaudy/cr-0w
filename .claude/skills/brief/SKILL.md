---
name: brief
description: Morning briefing. One screen covering reviews waiting, sprint, open MRs/PRs, open threads, vault tasks, background runs, plus mail/calendar when connected. Ends with a suggested top 3. Use when the user says "brief", "briefing", "what's on today", or starts the day.
context: fork
model: sonnet
effort: low
background: false
allowed-tools: Bash(bin/crow-activity:*) Bash(bin/crow-releases:*) Bash(bin/crow-tidy:*) Bash(bin/crow-retro:*) Bash(cdx runs:*) Bash(cdx status:*) Bash(logics-manager fleet:*)
---

# Brief

Runs in its own context (no conversation history): everything needed is gathered below. Return only the brief.

## Pre-gathered
Forge activity since yesterday:
!`"${CLAUDE_PROJECT_DIR}/bin/crow-activity" --since yesterday 2>&1 | head -150 || true`

Release tracks:
!`"${CLAUDE_PROJECT_DIR}/bin/crow-releases" 2>&1 | head -40 || true`

Repo hygiene (safe cleanups already applied, see `runbooks/repo-tidy.md`):
!`"${CLAUDE_PROJECT_DIR}/bin/crow-tidy" --apply --brief 2>&1 | head -60 || true`

Next moves across every Logics repo (reviews are left out, crow-activity has them):
!`logics-manager fleet brief --root "$(dirname "${CLAUDE_PROJECT_DIR}")" --no-reviews --limit 5 2>&1 | head -12 || true`

Account headroom:
!`cdx status --small --cached 2>&1 | grep -v '^Resolv' | head -10 || true`

Read-only except for one activity-log line. Target: **one screen, ≤ 40 lines**. Gather cheaply, in this order.

## Gather
0. **Last AFK run first**: if `local/state/afk/log.md` has entries since yesterday evening, read it and the AFK report before anything else. Never state how the night ended without it; if the report is missing, say which jobs never finished and offer to compile it (see `runbooks/afk-runs.md`, Morning).
1. **Open threads + status**: state dir `open-threads.md` and `status.md` (state dir = `vault.crow_dir` or `local/state/`).
2. **Forge activity**: already above (Pre-gathered). Don't explore repos further. "My open MRs" lists every open MR with its unresolved thread count.
   - "My GitLab/GitHub actions" are the owner's OWN outgoing actions (comments he left on others' MRs, MRs he accepted). Never put them in "Needs you" and never write "comments from you" as a to-do: a comment he left on a colleague's MR is waiting on the colleague. Only "Awaiting my review", unresolved threads on "My open MRs", approvals pending and failed pipelines on his MRs belong in "Needs you". If all three are "none", write "Nothing waiting on you (no review, no open MR)".
   - Comments on MRs may be SonarCloud's quality-gate note posted through the owner's token (crow-activity filters them). Never read an action as a review he gave unless it is clearly his.
   - Use the actions only as context ("reviewed !210"), and only when it changes the day (e.g. a review he gave that is still unanswered after 2+ days).
3. **Tracker** (optional, e.g. Jira through the Atlassian MCP; skip if none is connected): `searchJiraIssuesUsingJql` with
   `assignee = currentUser() AND sprint in openSprints() AND statusCategory != Done ORDER BY priority DESC`
   Ask only for key, summary, status and priority fields.
4. **Vault tasks** (only if a vault is configured): `rg -n "^\s*- \[ \]" "<vault>/00 Inbox" "<vault>/10 Work" "<vault>/20 Personal"` and keep only tasks with a due date (📅 YYYY-MM-DD or `due:`) on or before today. Also count the files in `00 Inbox`. Don't read notes further.
5. **Release tracks**: already above (Pre-gathered).
6. **Background runs**: `cdx runs --since today --json`. Keep only running or failed ones. **Accounts**: from the headroom table above, one line only when the recommended account is under 30% or an account the owner uses is blocked (with its reset).
   **Logics next moves**: from the fleet list above, only moves that unblock work (blocked docs, a task at 100% in a repo in delivery); skip pure closeout chores unless nothing else is open.
7. **Retros**: report files in `local/state/retros/` dated since the last brief (one line each: project, what was applied), unchecked `- [ ]` lines in `local/state/retros/proposals.md`, and `bin/crow-retro status` for failures.
8. **Mail + calendar**: only if Microsoft 365 tools are available in this session. Today's meetings, and unread mail that needs a reply. If unavailable, print one line saying so and move on.

Skip any source that errors and note it in one line at the bottom.

## Output (exactly this shape, skip empty sections)
```
# Brief: <weekday DD/MM>

## 🔴 Needs you
- review feedback on my MRs (unresolved threads, always first: fixing it gets them merged), reviews waiting on me (they unblock devs), approvals pending, failed runs, overdue tasks

## 📅 Today
- meetings and recurring rituals (listed in `local/context.md`), deadlines

## 🏃 Sprint
- KEY status: summary   (mine, in priority order)

## 🧵 Open threads
- stale (>3 days) or due ones only

## 🚀 Release tracks
- one line per track: name → target, next milestone (or "no roadmap yet"); then one suggested next step for the most recently active personal project

## 🧹 Repos
- what crow-tidy cleaned (one line), then per project only: uncommitted work, unpushed main, branches older than 3 days (offer to merge/delete), worktrees in /tmp

## 🔁 Learned
- one line per retro since the last brief (what changed); pending proposals to approve; failed retros

## ✅ Suggested top 3
1. …  (follow the priority order in CLAUDE.md)
```

## After
Append `- HH:MM [hub] brief` to today's activity log. Nothing else.
