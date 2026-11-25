## item_275_core_lands_finished_afk_branches - Core lands finished AFK branches on main
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: When an AFK run ends done, core lands its afk/* branch: rebase on main, crow-check, ff-only merge, prepush check, push. Any failure leaves the branch and marks the row needs-you with the reason.
- Keywords: afk, land, rebase, push, main, core tick
- Use when: changing how finished AFK branches reach main.
- Skip when: work repos (never auto-landed).

# Problem
- Landing is a hub Claude step today (runbooks/afk-runs.md). With no session the branch waits and the queue stalls.

# Scope
- In: a land step run by the core tick after a run ends done, for repos whose registry scope is personal. Steps: fetch, rebase the worktree branch on origin/main, bin/crow-check, merge --ff-only into main, bin/crow-prepush-check --check, push HEAD:main. Log each step to the AFK log.
- In: a dirty main worktree or a moved origin is handled by rebase and retry once, never by force.
- Out: work repos (scope work): no landing, row shows "branch ready" for the owner. Conflict resolution.

# Acceptance criteria
- AC1: A done run on a personal repo ends with its commits on origin/main and the row marked landed with the commit.
- AC2: Failing checks, a rebase conflict or a non-ff merge leave main untouched and mark the row needs-you with the failing step.
- AC3: A work-scope repo is never pushed.
- AC4: No step can open a permission prompt (runs as the core process, no Claude).

# AC Traceability
- req_067-AC1 -> AC1. req_067-AC2 -> AC2.

# Dependencies
- item_233 launcher (queue rows and states). item_250 core tick.

# Links
- Request: req_067_afk_queue_runs_itself.

# Priority
- Priority: High
- Rationale: Without it the queue needs a hub session; owner goal is Nest only.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
