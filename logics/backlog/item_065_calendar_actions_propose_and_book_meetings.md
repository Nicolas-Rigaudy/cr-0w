## item_065_calendar_actions_propose_and_book_meetings - Calendar actions: propose and book meetings
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Management
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-today gains propose and book: it finds a free slot for the attendees, shows the exact invite, and creates it only after the owner's OK.
- Keywords: calendar, graph, invite, free busy, propose, book, sprint planning, one-on-one
- Use when: a check or row needs a meeting booked.
- Skip when: reading the calendar (already in crow-today).

# Problem
- A missing ritual or 1:1 needs a meeting. Today the owner books it by hand in Outlook.

# Scope
- In:
  - `bin/crow-today propose --who a,b --len 60 --within 3d --title ...` returns slots from getSchedule (findMeetingTimes as fallback), working hours only.
  - `bin/crow-today book <draft.json>` creates the event; the draft holds subject, slot, attendees, room, body (French).
  - Deck: a check row's enter shows the draft in a modal; y books, e opens it as an Outlook draft, esc cancels.
  - Entra app consent for Calendars.ReadWrite; runbook graph-today updated.
- Out:
  - Booking without the OK.
  - Rescheduling others' meetings.

# Acceptance criteria
- AC1: propose returns at least one free slot for the team within the window, or says none.
- AC2: Nothing is created before y on the shown draft.
- AC3: The created event matches the draft (attendees, slot, body).

# AC Traceability
- All ACs -> this item. Proof: to fill at build.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_014_crow_management_checks_and_actionable_deck
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Closes the loop on ritual and 1:1 checks; carries outward risk so it comes last.

# Notes
- Invites reach other people: the CLAUDE.md approval rule applies, the modal is the approval.
- Room mailboxes accept bookings (the planning room accepted the 2026-11-05 one).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
