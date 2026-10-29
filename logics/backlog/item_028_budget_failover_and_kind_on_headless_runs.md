## item_028_budget_failover_and_kind_on_headless_runs - Budget, failover and kind on headless runs
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:12:33

# AI Context
- Summary: Every crow headless run passes --failover and --fallback-model; retro runs pass --kind retro; a budget is optional per run, no default (owner decision 2026-10-29).
- Keywords: budget, failover, kind, headless, runs
- Use when: changing how crow-runs, the AFK queue or crow-retro call cdx run.
- Skip when: interactive launches (cdx ignores budget there).

# Problem
- A usage limit stops runs and retros; nothing resumes them on another account.

# Scope
- In:
  - crow-runs dispatch and AFK: `--failover`, `--fallback-model`
  - optional `--budget` per dispatch or AFK job, passed through to cdx, no default
  - retro runs: `--kind retro`
- Out:
  - a default spend cap
  - interactive sessions (cdx applies budget to headless runs only)

# Acceptance criteria
- AC1: every crow headless run carries failover and a fallback model.
- AC2: a per-run budget can be set and the runbook documents it.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: AC1: every crow headless run carries failover and a fallback model.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)
- Request: `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`
- Primary task(s): `task_003_deliver_the_crow_cost_and_reliability_pass`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.
