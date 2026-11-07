## item_121_antigravity_weekly_quota_measurement_for_one_overflow_account - Antigravity weekly quota measurement for one overflow account
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:22:08

# AI Context
- Summary: Measure the free weekly quota of one Antigravity account, then enable it as an overflow row for public personal code and cheap text.
- Keywords: antigravity, quota, overflow, router, one account
- Use when: Before adding the Antigravity row to the router.
- Skip when: Several accounts and work data.

# Problem
- The free weekly quota has no published number; it only shows after an interactive launch.
- Owner decision 2026-11-07: one account as overflow after measuring the quota. Public personal code and cheap text only.

# Scope
- In:
  - Install `agy` and log in with the one account.
  - Run a fixed small task set and read the quota over a week.
  - Record the measured size in the req_029 spec.
  - Add the overflow row to the router table with the data class limits.
- Out:
  - Several accounts (reads as circumvention).
  - Work data and private personal data.

# Acceptance criteria
- AC1: The weekly size is written in the spec with the dates and method.
- AC2: The router row exists and refuses work and private data classes.
- AC3: A public personal code mission can overflow to it in a test.

# AC Traceability
- request-AC3 -> AC1. Proof: the measured weekly size in the spec.
- request-AC5 -> AC2. Proof: the router row refuses work and private data classes.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_029_spike_hosts_model_sources_and_local_ai_benchmarks`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Optional overflow capacity.

# Notes
- 2026-11-07 owner decision 31 (req_029 D5 option A).
- Needs `agy` login on the machine (spec Needs env).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
