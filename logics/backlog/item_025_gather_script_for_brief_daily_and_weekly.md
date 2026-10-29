## item_025_gather_script_for_brief_daily_and_weekly - Gather script for brief, daily and weekly
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: bin/crow-gather returns one JSON per routine so the brief, daily and weekly forks stop spending tool turns on files and commands.
- Keywords: gather, script, brief, daily, weekly
- Use when: changing what the brief, daily or weekly skills gather.
- Skip when: changing their output format or posting rules.

# Problem
- The forks spend 6-13 tool turns at 40k+ context each gathering files and command output.

# Scope
- In:
  - bin/crow-gather brief|daily|weekly --json (activity, releases, tidy, fleet, AFK log, open threads, status, vault tasks, runs, retros, since date, target paths, tracker ids from local config)
  - skills read the JSON through one `!` block
- Out:
  - MCP calls (tracker, mail) stay in the model

# Acceptance criteria
- AC1: each skill has one pre-gathered JSON block and no file-reading steps.
- AC2: JSON stays under 15k tokens on a busy day.

# AC Traceability
- request-AC2 -> This backlog slice. Proof: AC1: each skill has one pre-gathered JSON block and no file-reading steps.

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
