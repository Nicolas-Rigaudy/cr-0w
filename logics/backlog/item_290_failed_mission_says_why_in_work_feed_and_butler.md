## item_290_failed_mission_says_why_in_work_feed_and_butler - A failed mission says why, in the work feed and in the butler answers
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A mission failed on a core validation incident, yet its work row says "done" and the butler, asked why Nest is in fail state, found nothing. Carry the incident detail to the work feed, the hero and the butler.
- Keywords: nest, hero, fail, mission, incident, work feed, butler, brief
- Use when: Nest shows a failure the owner cannot explain.
- Skip when: the cause of this specific failure (item_291).

# Problem
- Owner 2026-11-25: "why are you in fail state" and "it can't even understand why it's in fail state".
- What failed (core state and `local/state/missions/20261125-plan-nest-desk-bugs/events.jsonl`): the desk-bugs plan mission went `running -> failed` at 14:17:43 with incident `no_version_file`: "plan.v2.md: v2 responses miss q1, q2, q3". It is the only failed row in the work feed, so the hero face shows fail and Ops shows a badge of 1.
- Causes:
  - `tools/crow-core/internal/work/work.go:338-341`: a failed mission takes `Why` from the last run cause, which is `done` (exit 0). The row reads failed / done. The incident detail never reaches the feed.
  - The butler answered from `brief · dispatch`, `sessions` and `jobs` (seq 68 to 70), not `work` or `missions`, then guessed stale usage-limit sessions. The persona has no rule for "why is Nest red". It also said the mission "still shows as running".
  - The hero model has `cause` (`hero/model.ts:157,188-192`) but the butler cannot see what Nest shows.

# Scope
- In:
  - core: for a failed mission, `Why` is the last incident's `detail` (fallback cause, then last run cause), trimmed to one line. Same for the mission section so `mission_get` returns it.
  - Butler brief: a `health` (or `work`) view that lists failed and stuck items with their why, the same rows that drive the hero mood.
  - persona.md: asked about a failure, a red face or a badge, read the failed work rows first and answer with item, cause and the next step (Retry, answer, cancel). Never guess from other sections.
  - Hero line keeps naming the failed item; the work row shows the why under the title.
- Out:
  - Fixing the validation rule itself (item_291).
  - Mission detail navigation (item_292).

# Acceptance criteria
- AC1: With the desk-bugs mission failed as today, its work row shows "plan.v2.md: v2 responses miss q1, q2, q3" (or a plain rewording), not "done".
- AC2: Asked "why are you in fail state", the butler names the failed mission and its cause in its first reply.
- AC3: `mission_get` on a failed mission returns the incident detail.
- AC4: Go tests for `work.go` failed-mission why and the brief view; persona change checked with a live question on the butler.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Related: item_291 (root cause of today's failure), item_292 (where the row and badge lead), item_158 (hero mood order).

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_290 (a failed mission says why) end to end. Read logics/backlog/item_290_failed_mission_says_why_in_work_feed_and_butler.md first, then tools/crow-core/internal/work/work.go around the failed-mission case and tools/crow-core/internal/butler/persona.md. Reproduce with the desk-bugs mission events in local/state/missions/20261125-plan-nest-desk-bugs/events.jsonl (copy into a test fixture, do not edit the live mission). Carry the incident detail to the work feed and mission_get, add the butler failure view and persona rule, then ask the live butler "why are you in fail state" and check its answer.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
