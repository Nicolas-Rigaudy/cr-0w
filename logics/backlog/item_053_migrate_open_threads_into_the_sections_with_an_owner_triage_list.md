## item_053_migrate_open_threads_into_the_sections_with_an_owner_triage_list - Migrate open-threads into the sections with an owner triage list
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Move the current open-threads lines into the new sections; drop verified-done items to the activity log; give the owner one keep/drop list for the stale rest.
- Keywords: open threads, migration, triage, stale, status notes
- Use when: migrating open-threads.md to the sectioned model.
- Skip when: format changes (item_052).

# Problem
- About 94 open lines and 46 ticked ones; about 60 open lines predate 2026-10-25 and many are done or superseded.

# Scope
- In:
  - Ticked lines → `activity/` log of their date (or today's), then removed.
  - Open lines verifiably done (merged MR, closed Logics item, later line saying done) → dropped with a note in the triage file.
  - Status/constraint lines → status.md or decisions.md.
  - Rest sorted into Me / Crow / Waiting; stale ones (before 2026-10-25) listed in `local/state/drafts/threads-triage.md` with crow's keep/drop suggestion; owner marks it in one pass, crow applies it.
- Out:
  - Changing skills (item_052).

# Acceptance criteria
- AC1: open-threads.md has only the three sections and no ticked lines.
- AC2: Every dropped line is traceable in the triage file or an activity log.
- AC3: The owner's keep/drop marks are applied.

# AC Traceability
- request-AC3 -> AC1-AC3. Proof: migrated file + triage file.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_012_owner_todos_and_open_work_model
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Needed before the deck and brief can read the sections.

# Notes
- Execution order: 2 of 3, after item_052. Needs the owner for the triage pass.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

