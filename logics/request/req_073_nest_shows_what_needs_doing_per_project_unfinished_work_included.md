## req_073_nest_shows_what_needs_doing_per_project_unfinished_work_included - Nest shows what needs doing per project, unfinished work included
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Nest shows per project what needs doing, including work started and not finished (a Jira ticket in progress, a branch with uncommitted or unpushed changes, a session that ended mid-task), and the Desk loses its project filter.
- Keywords: nest, project page, in flight, unfinished, jira, branch, desk filter
- Use when: building the Project page, the in-flight data or the Desk filter removal.
- Skip when: the MRs review flow (req_071), the hero itself (item_221).

# Needs
- Owner 2026-11-26, from Nest: "priorities are MR page and projects page so it helps me see what needs doing per project etc (also project filter on desk was not necessary i think). one thing for example is i started work on a jira ticket yesterday and it's not finished. but this shows up nowhere in nest".
- The case: a Jira ticket started in a work repo. A feature branch, 11 modified files, no commit, no MR, the session ended on 2026-11-25. Nest only lists live sessions, MRs and Logics items, so nothing shows it.
- The owner wants one place per project that answers "what needs doing here", and a way back into unfinished work with one click.

# Context
- `bin/crow-registry` already gives each project `git: {email, dirty, ahead, stale_trees}`; core serves it as the projects section. No branch name, no ticket, no last session.
- item_247 plans a generic Jira fetch of the owner's open sprint tickets (key, summary, status, repo hint), cached for the core. Not built.
- item_222 plans the Project page (route p/<name>: What needs you, Live, MRs, Missions, activity), mockups first. Not built.
- Desk has a project segment filter (`tools/nest/src/screens/Desk.tsx:699`, `proj` state, `panesOf(..., proj)`).
- Priority order from the owner: MR page first (req_071), then the project page (this request).

# Acceptance criteria
- AC1: For each registered project, core data lists in-flight work: Jira tickets of the owner in progress, branches with uncommitted or unpushed changes (ticket key read from the branch name), and the last ended session on that branch when it can be resumed.
- AC2: The Project page shows an In flight section from that data; each row has an action (resume the session, start a session on the ticket, open the MR or the branch).
- AC3: With today's state, the work repo page from the case shows that ticket with its branch, 11 changed files, its Jira status and a resume action.
- AC4: The Desk has no project filter; panes show every project.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-registry`, `local/state/sprint/fetch.py` (token read), `tools/crow-core` projects section
- `tools/nest/src/screens/Desk.tsx`
- req_063 (Nest pages), req_071 (MRs review flow), item_222, item_247

# Backlog
- item_308_desk_drops_the_project_filter (independent, small)
- item_247_desk_ready_list_shows_open_sprint_tickets_for_work_projects (its Jira fetch script is reused by item_307)
- item_307_in_flight_work_per_project_in_core_data (after item_247's fetch script)
- item_222_nest_project_page_and_planet_click_routing (gains the In flight section; mockups before build)

# Execution order
1. req_071 first (owner's order): item_298 start fails loudly, item_299 core review endpoints with item_300 review tab mockups, then item_301 build.
2. item_308 (any time, small).
3. item_247 Jira fetch, then item_307 in-flight data.
4. item_222 mockups (owner pick), then the build one screen at a time.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
