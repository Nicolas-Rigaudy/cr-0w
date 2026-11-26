## item_307_in_flight_work_per_project_in_core_data - In-flight work per project in core data
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90
> Confidence: 85
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Core data lists, per registered project, the work started and not finished: Jira tickets in progress, branches with uncommitted or unpushed changes, and the last ended session on each branch.
- Keywords: in flight, unfinished, branch, dirty, unpushed, jira in progress, resume, projects section
- Use when: building the data behind the Project page In flight section.
- Skip when: the page UI (item_222), the Desk Ready list rows (item_247).

# Problem
- A Jira ticket was started on 2026-11-25 in a work repo (its feature branch, 11 modified files, no commit, no MR). Its session ended. Nest shows nothing about it (req_073).
- `bin/crow-registry` gives `dirty` and `ahead` counts only: no branch, no ticket, no session.

# Scope
- In:
  - `bin/crow-registry` git block gains `branch`, `ticket` (first `[A-Z][A-Z0-9]+-\d+` in the branch name), `changed` (files), `ahead`, `last_commit_at`, and the same for each worktree under `<repo>.wt/`. Main or master with nothing pending is not in flight.
  - Jira: reuse item_247's fetch script and cache; tickets of the owner with status category In Progress are joined to a project by branch ticket key, else by the repo hint.
  - Last session: the most recent ended Claude session whose cwd is that checkout (session cache the Desk resume list already reads), with its conversation id and account, when cdx can resume it.
  - Core exposes the result as `in_flight` per project in the projects section (or `/api/state/work` filtered by project), one row per branch or ticket: project, branch, ticket key, summary, Jira status, changed, ahead, MR link if one exists, last session, age.
  - Generic code only (engine side); ticket regex and Jira site come from config.
- Out:
  - Writing to Jira or git. Page UI.

# Acceptance criteria
- AC1: With today's state the work repo entry from the req_073 case has one row: its feature branch, its ticket key, its Jira status, changed 11, the 2026-11-25 session as resumable.
- AC2: A clean main checkout gives no row; a pushed branch with an open MR gives a row with the MR link.
- AC3: Jira down or no token: git and session rows stay, ticket fields empty, a section note says why.
- AC4: pytest for the registry parse and the join; Go test for the core field.

# AC Traceability
- req_073 AC1 -> this item. Proof: AC1, AC2.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Dependencies
- item_247 Jira fetch script and cache (build that part first).

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_073_nest_shows_what_needs_doing_per_project_unfinished_work_included
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Unfinished work is invisible in Nest today; the owner lost track of a started ticket.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
