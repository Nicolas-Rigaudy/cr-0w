## task_013_usage_limit_mid_mission_policy_by_context - Usage limit mid-mission policy by context
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-07 11:53:38

# AI Context
- Summary: dispatch --context afk|day; requeue picks wait or move through crow_lib.limit_policy and logs run.limit_policy with the reset time.
- Keywords: usage, limit, mid, mission, policy, context
- Use when: Changing what happens to a run cut by a usage limit.
- Skip when: Account ranking or quota probing.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_115_usage_limit_mid_mission_policy_by_context`

# Acceptance criteria
- AC1: A night mission cut by the limit waits and continues on the same account after the reset.
- AC2: A daytime mission with a reset 3 hours away restarts cold on the next allowed account.
- AC3: A daytime mission with a reset 20 minutes away waits.
- AC4: Every choice is logged with the reset time it used.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_013_usage_limit_mid_mission_policy_by_context.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_013_usage_limit_mid_mission_policy_by_context.md` after implementation.

# AC Traceability
- request-AC1 -> This task. Proof: tests/test_crow_quota_failover.py (test_afk_run_waits_on_same_account_despite_free_target); pytest 21 passed.
- request-AC2 -> This task. Proof: tests/test_crow_quota_failover.py (test_day_run_with_far_reset_moves_cold); pytest 21 passed.
- request-AC3 -> This task. Proof: tests/test_crow_quota_failover.py (test_day_run_with_near_reset_waits); pytest 21 passed.
- request-AC4 -> This task. Proof: tests/test_crow_quota_failover.py (run.limit_policy event asserts reset_5h_at and reset_in_min); pytest 21 passed.

# Validation
- (no validation recorded yet)
- pytest tests/test_crow_quota_failover.py tests/test_crow_retro.py: 30 passed. crow-check ruff is red on lint already in crow_lib.py before this work.
- Finish workflow executed on 2026-11-07.
- Linked backlog/request close verification passed.

# Report
- Not started.
- Finished on 2026-11-07.
- Linked backlog item(s): `item_115_usage_limit_mid_mission_policy_by_context`
- Related request(s): `req_037_mission_kinds_and_run_contract`

# Links
- Request: `req_037_mission_kinds_and_run_contract`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
