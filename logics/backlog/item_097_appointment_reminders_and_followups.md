## item_097_appointment_reminders_and_followups - Appointment reminders and follow-ups in the personal vault
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Track appointments and follow-ups (health and admin, global level) in the personal vault and surface reminders in the personal brief.
- Keywords: butler, appointments, reminders, health, follow-up, calendar, vault
- Use when: building this butler slice.
- Skip when: work rituals on the work account.

# Problem
- Appointments and recurring check-ups (doctor, dentist, renewals) are remembered by hand. The owner wants global reminders, not detailed health tracking.

# Scope
- In:
  - A vault note (or folder) listing appointments and recurring follow-ups: who, when, what to bring, next due date, recurrence.
  - Capture from chat (`crow note`-style) and from calendar events that look like appointments.
  - Reminders: due soon in the personal brief; an optional calendar event created after an explicit OK.
- Out:
  - Sensor or wearable data, medical details, trends.
  - Booking appointments on third-party sites.

# Acceptance criteria
- AC1: Telling crow about an appointment or a check-up adds it to the vault list with its next due date.
- AC2: Items due within a configurable window appear in the personal brief.
- AC3: A recurring follow-up rolls its next due date forward when marked done.
- AC4: Creating a calendar event asks first and shows the exact event.
- AC5: No item stores medical detail beyond what the owner typed.

# AC Traceability
- AC1 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC2 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC3 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC4 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.
- AC5 -> request-AC1/AC2. Proof: test or a recorded run in the personal vault.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_026_personal_butler`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Small, low-risk, and the owner named reminders as the health need.

# Notes
- Format in the vault follows the vault Claude's conventions: delegate the note layout to the vault session if it owns one.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
