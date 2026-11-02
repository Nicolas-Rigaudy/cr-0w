## item_042_document_the_capture_convention_and_wire_the_recap_skills_to_the_vault_notes - Document the capture convention and wire the recap skills to the vault notes
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
- Summary: Document the capture convention and wire the recap skills and the brief to the vault notes.
- Keywords: document, capture, convention, wire, recap, skills, vault, notes
- Use when: building this slice of req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap.
- Skip when: another slice of the same request.

# Problem
- The system only works if every recap reads it and the convention is written once.

# Scope
- In:
  - a runbook `runbooks/vault-capture.md` (generic) and a private note in local/runbooks if names are needed, indexed in the README
  - add the `crow-notes` JSON to the gather list of mgmt-weekly, daily, weekly and brief
  - brief asks at most three questions on topics with no trace in 3 days (reuses the eod question logic)
  - handoff note to the vault session to add the convention to the vault CLAUDE.md
  - update the owner-prefs line about recaps
- Out:
  - changing the vault CLAUDE.md from the hub

# Acceptance criteria
- AC1: each of the four skills names the vault notes as a source and states in its output that it read them
- AC2: the runbook lists tags, sections and who writes where
- AC3: the vault session received the handoff text

# AC Traceability
- request-AC3, request-AC5, request-AC6 -> This backlog slice. Proof: the ACs above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_009_capture_the_owner_s_implicit_work_notes_through_the_day_feed_every_recap`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: depends on: item_040, item_041.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

