## req_067_afk_queue_runs_itself - AFK queue runs itself without a Claude session
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: The AFK queue moves forward with no Claude session: crow core launches the next eligible row within the floors (item_233), lands each finished branch on main when its checks pass (item_275), and Nest shows the runner with pause and needs-you rows (item_276).
- Keywords: afk, queue, runner, land, floors, nest, autonomous
- Use when: changing how AFK runs are launched, landed or shown.
- Skip when: writing AFK job prompts (afk skill).

# Needs
- Owner 2026-11-25: "why does it need a session prompt in nest, shouldn't it handle the afk run itself". Goal of the day: use crow from Nest only.
- Today a finished run only leaves an afk/* branch. A hub Claude session rebases, checks, fast-forwards main, pushes and starts the next queue row. Without that session the queue stops after the runs in flight.
- The night of 24/11 stalled 14h on one permission prompt in the hub session.

# Context
- item_233 already plans the launcher (queue config with floors, account pick, decision log). It leaves merging to the hub. This request adds landing and the Nest view and asks item_233 to run inside the core tick (item_250) instead of a separate loop.
- Missions (adr_001, internal/mission/tick.go) already admit queued work under a max-runs cap; reuse that admission instead of a second scheduler where it fits.
- Hub landing steps today (runbooks/afk-runs.md): rebase on main, crow-check, merge --ff-only, crow-prepush-check --check, push HEAD:main.

# Scope
- In: launch (item_233), auto-land for personal repos whose registry scope allows push to main (item_275), Nest view (item_276).
- Out: work repos (never push main; branch stays for the owner). Rendering inside runs (separate sandbox issue).

# Acceptance criteria
- AC1: A queue of 3 rows on 2 accounts runs to the end with no Claude session: each finished branch is on origin/main and the next row started within the floors.
- AC2: A branch that fails checks or does not fast-forward is not landed; its row shows "needs you" in Nest with the reason, and the queue moves on.
- AC3: The owner pauses and resumes the runner from Nest.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): adr_001

# References
- `tools/crow-core/tick_loop.go`
- `tools/crow-core/internal/mission/tick.go`
- `tools/crow-core/internal/api/queue.go`
- `runbooks/afk-runs.md`

# Backlog
- item_233_afk_quota_scheduler_report_languages_and_summary_schema (launcher; run it from the core tick)
- item_275_core_lands_finished_afk_branches
- item_276_nest_afk_queue_runner_view
- Order: item_233 launcher part, then item_275, then item_276.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
