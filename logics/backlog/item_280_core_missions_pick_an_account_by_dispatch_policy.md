## item_280_core_missions_pick_an_account_by_dispatch_policy - Core missions pick an account by dispatch policy
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 75%
> Confidence: 70%
> Progress: 0%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The core tick started plan mission nest-desk-bugs on an account outside the dispatch policy. Mission account choice must follow the same policy as crow dispatch.
- Keywords: core, tick, mission, account, policy
- Use when: choosing the account for a mission run.
- Skip when: interactive sessions.

# Problem
- Tick 2026-11-25 10:26: run 33329c47 started on an account outside the dispatch policy (`account_class: personal`).
- The tick picks by data class (work/personal), not by the dispatch policy that `crow dispatch` uses.

# Scope
- In:
  - The tick reads the same allowed-account list as `crow dispatch` (cdx `--policy`, `local/accounts.yaml` fallback).
  - No allowed account with quota: the mission stays queued with the reason shown.
- Out:
  - Quota failover rules.

# Acceptance criteria
- AC1: With the current config every mission starts on the dispatch account.
- AC2: An account outside the policy is never picked (go test).
- AC3: A held mission shows why in Nest and `crow mission show`.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`.

# Priority
- Priority: Medium
- Rationale: quota and identity rule.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
