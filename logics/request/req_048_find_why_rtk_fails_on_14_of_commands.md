## req_048_find_why_rtk_fails_on_14_of_commands - Find why rtk fails on 14% of commands
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: rtk failed on 14% of the commands in one week of the daily-efficiency research; find the causes and fix or exclude them.
- Keywords: find, rtk, fails, commands
- Use when: picking up this hub improvement.
- Skip when: working on another repo.

# Needs
- rtk failed on 14% of the commands in one week of the daily-efficiency research; find the causes and fix or exclude them.

# Context
- Seen 22/10/2026 in the daily-efficiency research.
- Failures cost tokens twice when the raw command reruns.

# Acceptance criteria
- AC1: The failing command kinds are listed with their cause.
- AC2: Each cause is fixed, excluded from the rewrite hook, or reported upstream.
- AC3: The failure rate is measured again after the change.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `rtk gain --history`
- `rtk discover`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

