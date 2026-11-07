## item_115_usage_limit_mid_mission_policy_by_context - Usage limit mid-mission policy by context
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: When a mission run hits a usage limit, core decides by context: wait for the reset at night, restart cold on the next account by day when the reset is far.
- Keywords: usage limit, failover, afk, night, cold restart, mission
- Use when: Building the end-of-run rule for limits in core.
- Skip when: Quota warnings (item_068) and dispatch failover flags (item_028).

# Problem
- A run that hits a usage limit ends. `--continue-run` pins the parent's account and cannot combine with `--failover`. The next step needs a policy.

# Scope
- In:
  - Context flag on the mission: AFK/night batch or daytime.
  - AFK/night: wait for the reset, then continue. A wait of about 4 hours at night is fine.
  - Daytime: restart cold on the next allowed account from `brief.md` and the version files when the reset is more than about 1 hour away. Otherwise wait.
  - Log the choice and its reason as a mission event.
- Out:
  - Account ranking (item_095) and quota probing (crow-room).

# Acceptance criteria
- AC1: A night mission cut by the limit waits and continues on the same account after the reset.
- AC2: A daytime mission with a reset 3 hours away restarts cold on the next allowed account.
- AC3: A daytime mission with a reset 20 minutes away waits.
- AC4: Every choice is logged with the reset time it used.

# AC Traceability
- request-AC1 -> AC1 to AC4. Proof: the acceptance criteria above, checked when built.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_037_mission_kinds_and_run_contract`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Decides how long missions stall and how much cache is lost.

# Notes
- 2026-11-07 owner decision 12 (req_037 D2, option C refined by context).
- The 1 hour and 4 hour values are owner numbers; tune from real resets.
- Needs reset times from `cdx status --json` (req_029 Need 6) and the brief rebuild of item_107.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

