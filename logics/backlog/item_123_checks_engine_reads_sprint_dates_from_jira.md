## item_123_checks_engine_reads_sprint_dates_from_jira - Checks engine reads sprint dates from Jira
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The checks engine reads the active sprint dates from Jira with the crow-jira token, cadence as fallback.
- Keywords: checks, engine, reads, sprint, dates, jira
- Use when: Sprint-related checks rules or the Jira token.
- Skip when: Checks rules not tied to the sprint.

# Problem
- bin/crow-checks (item_063) only knows the configured cadence (2 weeks ending Wednesday). A scoped read token crow-jira is already in the keyring (local/runbooks/service-tokens.md). Owner decision 2026-11-07: read sprint dates from Jira.

# Scope
- In:
  - Fetch the active and next sprint (start, end, name) for the configured board via api.atlassian.com/ex/jira/<cloudId>/rest/agile or rest/api/3, basic auth email:token.
  - Cache the result under the state dir with a short TTL.
  - Fall back to the cadence when the token or the network is missing, and say so in the brief.
  - Board id and cloudId in local/checks.yaml; docs/checks.example.yaml shows the keys.
- Out:
  - Writing to Jira.
  - Other Jira data (tickets, velocity).

# Acceptance criteria
- AC1: With a valid token, sprint rules use Jira dates.
- AC2: Without a token or offline, rules use the cadence and the brief shows the fallback.
- AC3: Tests use a recorded fixture, no live call.
- AC4: The token never appears in logs or the cache.

# AC Traceability
- AC1 -> this item. Proof: to record at closeout.
- AC2 -> this item. Proof: to record at closeout.
- AC3 -> this item. Proof: to record at closeout.
- AC4 -> this item. Proof: to record at closeout.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Real dates make the sprint checks trustworthy.

# Notes
- Owner decision 2026-11-07 on return from the AFK run (local/state/afk/decisions-2026-11-07-afk.md).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

