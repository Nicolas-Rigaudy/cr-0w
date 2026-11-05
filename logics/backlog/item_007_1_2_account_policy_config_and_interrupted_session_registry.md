## item_007_1_2_account_policy_config_and_interrupted_session_registry - 1.2 account policy config and interrupted-session registry
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub orchestration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-05 15:23:56

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: account, policy, config, interrupted, session, registry
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Account rules exist only as prose and the registry was built by a one-off script.

# Scope
- In:
  - Machine-readable accounts file in local/ and a bin/ scanner that writes the registry as JSON plus a readable list.
- Out:
  - Any UI, any handoff execution.

# Acceptance criteria
- AC1: accounts file with policies is read by one helper used by every crow command.
- AC2: the scanner excludes retros and automatic jobs and reuses a cache keyed by file mtime and size.

# Validation
- 2026-11-05: AC1 and AC2 proven in task_002 (2026-10-20): bin/crow-accounts is the single reader of local/accounts.yaml for crow, crow-runs, crow-retro and crow-memory; bin/crow-sessions excludes retros and one-shots and caches by mtime_ns and size. Full suite re-run: 273 tests pass.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: accounts file with policies is read by one helper used by every crow command.
- request-AC2 -> This backlog slice. Proof: AC2: the scanner excludes retros and automatic jobs and reuses a cache keyed by file mtime and size.
- request-AC7 -> This backlog slice. Proof: AC2: the scanner excludes retros and automatic jobs and reuses a cache keyed by file mtime and size.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_002_keep_working_when_an_account_runs_out`
- Architecture decision(s): (none yet)
- Request: `req_001_crow_account_switching_and_session_resume`
- Primary task(s): `task_002_deliver_account_switching_and_session_resume_in_three_steps`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Notes
- Task `task_002_deliver_account_switching_and_session_resume_in_three_steps` was finished via `logics-manager flow finish task` on 2026-11-05.
