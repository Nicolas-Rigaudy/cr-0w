## req_017_token_levers_quota_switch - Token saving levers and quota-aware account switching
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: high
> Theme: cost
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-03 17:16:17

# AI Context
- Summary: Apply a colleague's 16 token-saving levers to crow and make crow warn before a usage limit and move work to an allowed account instead of failing.
- Keywords: token, cache, compaction, quota, usage limit, failover, account, headless, stream-json, mcp, budget
- Use when: changing how crow dispatches runs, watches quotas, or advises compaction.
- Skip when: per-repo token work inside a project session.

# Needs
- 2026-11-03: the default account hit its 5-hour limit at 15:19 UTC. Run ab76eee8 (lambda review-213) died, four other sessions on that account (lambda Sonar fixes, CD-546 MR 2 handoff, hub sessions) stopped and could not continue, and the owner moved to the shared account by hand. Nothing warned before, nothing switched.
- A colleague's table lists 16 levers with measured gains. Crow covers some; the rest go in this request.

# Context
## Why nothing happened on the limit
- `crow-usage` computes 5h remaining and burn rates, but only the deck, watch and brief show them. No hook pushes a warning into a live session.
- `crow-runs` dispatch never passes `cdx run --failover`. A run lands on one account and dies there.
- `cdx rescue <default> --to <shared>` only knows the account's last conversation (the review run). The other stopped sessions need `crow-sessions` to find them.
- `crow runs sessions` and `crow runs relaunch` already find and move stopped sessions, but nothing runs them or tells the owner they exist.
- The shared account has policy `ask-first` in `local/accounts.yaml`, and nothing asks. The only paths to the shared account were manual.
- The limit death shows as `failed` in the inbox hook, though `crow-runs` classifies it `limit` elsewhere.

## Levers already covered (no item)
- 6 No parallel subagents: hub rule "no speculative parallel runs", max 2 runs.
- 8 Compact before 200k: auto-compact cap at 150k, warning at 120k (runbooks/long-sessions.md).
- 11 State outside the conversation: /wrap, handoffs, open threads, Logics.
- 12 Deterministic code: `bin/` scripts first (crow-activity, crow-mrs, crow-review-check).
- 14 Code graph: graphify listed in runbooks/agent-tooling.md. Repo sessions own it.

## Levers to build (one backlog item per group)
| Lever | Gap | Item |
| --- | --- | --- |
| 15 Quota and account switch | no warning, no failover, no auto-resume | item_068 |
| 3 No useless MCP, 4 fixed model/effort, 7 guard rails, 16 headless traps | dispatch flags missing | item_069 |
| 2 Live session per task, 5 skill sent once, 13 send only the diff | multi-step runs relaunch `claude -p` | item_070 |
| 9 Compact only if worth it, 10 compact before cache expiry | single 150k rule, no idle rule | item_071 |
| 1 Measure first | cache reload/reread/output split not shown | item_072 |

# Acceptance criteria
- AC1: A live hub session gets a warning before its account's 5h window runs out, with the allowed fallback named.
- AC2: A background run that hits a limit continues on an allowed account or resumes at reset, without the owner.
- AC3: Every headless dispatch carries the lever 3/4/7/16 flags.
- AC4: The re-measure (runbooks/cost-audit.md method) shows lower cache reload share than the 2026-11-03 baseline.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-runs`, `bin/crow-usage`, `bin/crow-ctx`, `bin/crow-review-run`
- `local/accounts.yaml`
- `runbooks/dispatch-background-runs.md`, `runbooks/long-sessions.md`, `runbooks/cost-audit.md`
- Source table: colleague's "Levier / Pourquoi ça coûte / Comment faire / Gain constaté" (16 rows), received 2026-11-03.

# Backlog
- item_072_token_baseline_cache_split (done 2026-11-03: reload 32%, reread 46%, output 22%; local/state/research/token-baseline-2026-11.md)
- item_068_quota_warning_and_account_failover (done 2026-11-03: warning hook, auto fallback with caps, limit requeue, `crow rescue`, doctor check)
- item_069_headless_dispatch_flags (done 2026-11-03: model id, effort, budget, timeout, headless env, stuck watchdog; strict MCP waits on cdx per-run extra args)
- item_071_compaction_thresholds_and_idle_rule
- item_070_live_session_multi_step_runs (parked 2026-11-03: resumes already hit the 1h cache, see item notes)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
