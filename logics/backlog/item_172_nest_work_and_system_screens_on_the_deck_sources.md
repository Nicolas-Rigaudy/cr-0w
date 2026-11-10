## item_172_nest_work_and_system_screens_on_the_deck_sources - Nest Work and System screens on the deck sources
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: first versions of the prod_005 Work screen (MRs and reviews waiting, today's mail and meetings) and System screen (accounts and quota, usage, jobs, dispatch history) over the item_169 sections.
- Keywords: nest work screen, system screen, mrs, reviews, today, mail, usage, accounts, jobs, dispatch history
- Use when: building Work or System in Nest.
- Skip when: sprint board and reports (req_034 later items), machines and releases (req_035).

# Problem
- deck's MR, Today, Usage and history tabs have no Nest view. The phone cannot see what waits on the owner outside missions (req_054).

# Scope
- In:
  - Mockup first, both screens, owner OK before code.
  - Work: MRs and reviews waiting from `mrs` (who acts, pipeline, threads), flagged mail and today's meetings from `today`, with done and snooze (item_170). Read-only link out to the forge.
  - System: accounts with quota and ETA from `accounts` and `usage`, jobs from `jobs`, dispatch run history from `dispatch` (label, project, status, duration, cost, commits).
  - Navigation entries per prod_005 phone bar (Work, More > System).
- Out:
  - Sprint board and Jira (req_034), machines, releases, backups (req_035), review actions.

# Acceptance criteria
- AC1: Work lists the same MRs waiting on the owner and the same flagged mail as deck's MR and Today tabs.
- AC2: System lists the same accounts, quota and the last dispatch runs as deck's Usage and history.
- AC3: Done and snooze from Work act on the mail (item_170) and the row leaves the list.
- AC4: Both Fold 6 screens pass req_038 checks; bundle budget holds (lazy route chunks allowed).

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: deck's other tabs visible from the phone.

# Decision framing
- Product framing: Not needed (prod_005 screens)
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`, `req_034_work_views_in_nest`, `req_035_forge_ops_in_nest`
- Depends on: `item_169_core_serves_every_deck_source_as_read_only_sections`, `item_170_core_actions_on_live_sessions_and_dispatch_runs`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: Medium
- Rationale: after Home; Home covers the daily glance.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
