## req_034_work_views_in_nest - Work views in nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Feature track. Bring the owner's work day into nest: Today page, sprint and team board, reports pipeline (daily, weekly, management weekly), Monday incidents pass, FU registry. Built on existing crow skills and gather scripts; models only narrate.
- Keywords: today, Jira, sprint, team board, velocity, daily, weekly, management weekly, recap, incidents, follow-ups, reports pipeline
- Use when: building work views or report jobs.
- Skip when: changing the report skills' content rules (edit the skill).

# Needs
1. **Today page** (Ultron). Ranked list: reviews waiting on the owner, the owner's MRs with new threads or red CI, sprint tickets by priority and blockers, open threads due, meetings. Ranking inputs: priority order from CLAUDE.md (prod launch, reviews, sprint, roadmap, R&D, side projects), due dates, who waits. One tap launches the matching mission kind.
2. **Sprint and team board.** Board by status for the team, velocity committed vs done per person (existing sprint velocity notes), who is blocked and on whom, carry-over.
3. **Reports pipeline** (Ultron recap and weekly, Jarvis digests). Gather scripts write facts; the skill drafts; the draft lands as a review item; core publishes after OK. Reports: daily (before 11:30), Friday weekly and slides, Monday management weekly in business terms, mission digest at evening.
4. **Recap "Fait" block.** Done items from git, MRs, Jira and the activity log, with links, editable before publish.
5. **Weekly candidates table.** Items proposed for the weekly with keep or drop per row.
6. **Incidents pass** (Ultron, Monday). Group the week's incidents and manual gestures by cause; propose one runbook, rule or automation each; accepted ones become Logics items.
7. **FU registry** (Ultron). Deferred follow-ups per project, surfaced when the project comes up.

# Context
- Owner choices 2026-11-06: Today page, sprint and team board, reports pipeline, incidents pass.
- Existing: `brief`, `daily`, `weekly`, `mgmt-weekly` skills, `bin/crow-activity`, review ledgers, open threads.
- Owner rules: read Jira, mail and calendar before any team draft; French for work artifacts; never send to the CEO directly.

# Acceptance criteria
- AC1: Today page ranks from real data and launches a review mission in one tap.
- AC2: The daily draft appears as a review item before 11:30 and publishes only after OK.
- AC3: A Monday incidents pass produces grouped causes with proposals from a real week of events.
- AC4: Board and velocity match the Jira board for the current sprint.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_014_crow_management_checks_and_actionable_deck.md`
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
