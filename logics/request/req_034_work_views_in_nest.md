## req_034_work_views_in_nest - Work views in nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 65
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Work views in nest and deck on top of core: Today page, sprint and team board, reports pipeline, incidents pass.
- Keywords: today, sprint, team, velocity, reports, daily, weekly, mgmt-weekly, incidents, roadmap proposals
- Use when: building the work side of nest.
- Skip when: changing the report skills themselves.

# Needs
1. **Today page.** Next tickets by priority and blockers, MRs and reviews waiting on the owner, one tap launches a ticket as a mission.
2. **Sprint and team board.** Sprint tickets by status, velocity per person, who is blocked.
3. **Reports pipeline.** Daily, weekly and management weekly drafted in the background from gather scripts; the model narrates facts only; retouch and approve publish from nest.
4. **Incidents pass.** Every failure and every manual fix the owner does is logged; a weekly pass proposes runbooks, autonomy upgrades or roadmap items.

# Context
- Reuses existing scripts and skills: brief, daily, weekly, mgmt-weekly, crow-activity, req_014 checks engine.
- Publishing work text keeps the approval rule (exact draft, explicit OK).
- Depends on core API (req_028 slice 3) and nest MVP.

# Acceptance criteria
- AC1: The Today page lists the owner's next actions and launches one as a mission.
- AC2: A daily draft appears in the inbox before the meeting, ready to approve.
- AC3: The weekly incidents pass produces at least one proposal from real logs.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
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
