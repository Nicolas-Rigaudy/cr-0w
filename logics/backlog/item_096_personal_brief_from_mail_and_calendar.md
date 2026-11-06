## item_096_personal_brief_from_mail_and_calendar - Personal brief from personal mail and calendar on the personal account
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
- Summary: A personal brief on the personal account: personal mail triage, today's agenda and what needs an answer, written to the personal vault daily note.
- Keywords: butler, personal brief, gmail, calendar, vault, privacy
- Use when: building this butler slice.
- Skip when: work rituals on the work account.

# Problem
- Personal mail and calendar get no triage. The work brief runs on the work account and must never read personal data.

# Scope
- In:
  - A `personal-brief` skill built on the pattern of `.claude/skills/brief` and `runbooks/own-routines.md`, run on the personal account only (scope `personal`, req_025).
  - Sources: personal Gmail (unread, needs-reply, bills and admin mail), personal calendar (today, next 7 days), appointment follow-ups due (item_097), budget alerts (item_098) once they exist.
  - Output: a section in the personal vault daily note; a one-screen summary in the session.
  - Privacy rules written into the skill: what may show in the deck, the work brief and phone notifications (no amounts, no mail bodies outside the personal vault).
  - Runbook for connecting personal connectors on the personal account and checking them (`claude mcp get`, `/mcp`).
- Out:
  - Sending mail, accepting invites (drafts only, send with explicit OK).
  - Work data of any kind.

# Acceptance criteria
- AC1: The skill refuses to run on an account outside the personal scope and says why.
- AC2: A run lists unread personal mail needing action, today's agenda and the next 7 days, and writes them to the vault daily note.
- AC3: Mail replies are drafts only; nothing is sent without an explicit OK.
- AC4: The privacy rules are in the skill and the work brief shows at most a count of personal items.
- AC5: A runbook covers connecting and checking personal connectors on the personal account.

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
- Priority: High
- Rationale: First visible butler value; the other slices feed into it.

# Notes
- Depends on item_095 (scope routing) and on the owner connecting personal Gmail/Calendar on the personal account.
- Reuse the gather-script approach of the work brief (scripts before exploration).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
