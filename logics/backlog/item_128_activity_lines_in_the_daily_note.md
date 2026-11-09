## item_128_activity_lines_in_the_daily_note - Activity lines in the daily note
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub data
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 10:46:09

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: activity, lines, daily, note
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Activity goes to local/state/activity/<date>.md, which the owner never sees.

# Scope
- In:
  - bin/crow-log "[tag] text": lock in local/state, append '- HH:MM [tag] text' at the end of the '## CR-0W' section of today's daily note (heading added before 'End of day' if missing), atomic replace.
  - Create a missing daily note from the vault Daily Note template, substituting the known date calls; any other template call falls back to local/state/activity.
  - crow-gather activity_logs and crow-deck doneToday read the daily note section; fallback to old activity files for past days.
  - Skills and runbooks that append activity (brief, eod, wrap, retro, afk) call crow-log.
- Out:
  - Editing the vault Daily Note template (vault Claude)

# Acceptance criteria
- crow-log on a temp vault creates or appends the section, keeps other note content byte-identical.
- gather and deck show today's lines from the daily note.
- Tests pass.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: crow-log on a temp vault creates or appends the section, keeps other note content byte-identical.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_006_crow_data_in_the_vault`
- Architecture decision(s): (none yet)
- Request: `req_041_crow_data_in_the_vault`
- Primary task(s): `task_015_build_crow_data_in_the_vault`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_015_build_crow_data_in_the_vault`

# Notes
- Task `task_015_build_crow_data_in_the_vault` was finished via `logics-manager flow finish task` on 2026-11-09.
