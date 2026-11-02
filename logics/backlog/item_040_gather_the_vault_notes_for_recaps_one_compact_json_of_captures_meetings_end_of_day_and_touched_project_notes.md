## item_040_gather_the_vault_notes_for_recaps_one_compact_json_of_captures_meetings_end_of_day_and_touched_project_notes - Gather the vault notes for recaps: one compact JSON of captures, meetings, end-of-day and touched project notes
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Gather the vault notes for recaps into one compact JSON.
- Keywords: gather, vault, notes, recaps, compact, json, captures, meetings, end, day, touched, project
- Use when: building this slice of req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap.
- Skip when: another slice of the same request.

# Problem
- Recaps need the vault lines, and a gather script keeps the model's tool turns low (same pattern as the gather scripts in req_008).

# Scope
- In:
  - a script `bin/crow-notes --since <date|Nd>` returning JSON: per day the Capture, Meetings & notes and End of day lines, the `00 Inbox` notes added, and the project notes (`10 Work`) modified in the window with their first lines
  - skips empty template lines and the habit buttons
  - unit tests on a fixture vault
- Out:
  - writing to the vault
  - summarizing (the skill does that)

# Acceptance criteria
- AC1: on the fixture vault the JSON lists captured lines with day, time and tags
- AC2: empty template placeholders are absent
- AC3: runs in under two seconds on the real vault, output under about 6k tokens for a week

# AC Traceability
- request-AC2 -> This backlog slice. Proof: the ACs above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: depends on: item_039 (daily format).

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

