## item_039_add_crow_note_one_line_timestamped_capture_into_the_daily_note - Add crow note: one-line timestamped capture into the daily note
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
- Summary: Add `crow note`: append one timestamped line to the daily Capture.
- Keywords: add, crow, note, line, timestamped, capture, daily
- Use when: building this slice of req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap.
- Skip when: another slice of the same request.

# Problem
- The owner needs a note path under five seconds.

# Scope
- In:
  - a `bin/crow-note` script and a `crow note` entry in the CLI help (group Data)
  - reads `vault.path` from local/projects.yaml, finds `40 Journal/Daily/YYYY-MM-DD.md`, creates it from `90 Templates/Daily Note.md` (Templater fields filled for the date) when missing
  - appends `- HH:MM text #tags` under `📥 Capture`, keeps the existing lines
  - tags are free; with no tag the line is still written
  - tests with a fixture vault
- Out:
  - real vault files other than the daily Capture
  - any model call

# Acceptance criteria
- AC1: running `crow note "x" #naval` twice adds two timestamped lines under Capture and changes nothing else in the file (diff)
- AC2: with no daily for the day, the file is created from the template with all sections
- AC3: no network or model call (script only)

# AC Traceability
- request-AC1 -> This backlog slice. Proof: the ACs above.

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
- Rationale: first slice, others depend on it. Depends on: none.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

