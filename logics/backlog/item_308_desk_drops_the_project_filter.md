## item_308_desk_drops_the_project_filter - Desk drops the project filter
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95
> Confidence: 90
> Progress: 100%
> Complexity: Low
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Remove the project segment filter from the Nest Desk; panes always show every project.
- Keywords: desk, project filter, segment, panes
- Use when: removing the Desk project filter.
- Skip when: the scope chip of other pages, the Project page.

# Problem
- Owner 2026-11-26: "project filter on desk was not necessary". Per-project views move to the Project page (item_222).

# Scope
- In:
  - `tools/nest/src/screens/Desk.tsx`: drop the `proj` state, the segment row (around line 699), the `allProjects` prop and the `proj` argument to `panesOf`; `projectsOfPanes` goes if nothing else uses it.
  - Update Desk tests and fixtures.
- Out:
  - Launcher project pickers (Start, Dispatch forms) stay.

# Acceptance criteria
- AC1: No project filter on the Desk, laptop or phone.
- AC2: tsc clean, vitest green, `npm run audit:scroll` exit 0, Desk shots laptop and phone.

# AC Traceability
- req_073 AC4 -> this item. Proof: AC1.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_073_nest_shows_what_needs_doing_per_project_unfinished_work_included
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Small cleanup; owner said it is not needed.

# Outcome
- Commit 7af1674, live on Nest 26/11. Desk panes keep every project; launcher pickers unchanged. Tests 418 pass, scroll audit clean, shots in vault CR-0W/Research/nest-desk/item308/.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
