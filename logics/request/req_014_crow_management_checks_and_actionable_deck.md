## req_014_crow_management_checks_and_actionable_deck - Crow as management assistant: checks that turn into deck actions
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Complexity: High
> Theme: Management
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Crow watches calendar, mail, Jira and the vault for management gaps (sprint rituals not booked, 1:1s due, mail left unanswered) and turns each into an actionable deck row with a drafted fix.
- Keywords: crow, management, checks, deck, actions, sprint, rituals, one-on-one, calendar, mail, vault
- Use when: building management checks, deck row actions, meeting proposals, 1:1 or sprint ritual prep.
- Skip when: deck layout work (req_011), Jira sprint edits.

# Problem
- The deck shows state but rarely lets the owner act on it.
- Crow does not notice management gaps. Example from 2026-11-03: the sprint ended that day in Jira and no sprint planning was booked. The owner caught it alone.
- The vault already tracks 1:1 cadence per dev, team objectives and rituals, but nothing reads it to raise what is due.

# Needs
- A checks engine that detects gaps and feeds "You decide" in the deck and the brief.
- Every deck row has a default action and a few keys.
- Meeting proposals from free/busy, booked only after the owner's OK.
- 1:1 prep and capture per dev in the vault person note.
- Sprint ritual prep: planning, review, retro.

# Scope
- In: items 063 to 067.
- Out:
  - Sending mail or invites without an OK on the exact draft.
  - Changing the Jira sprint unless the owner asks (local runbook sprint-planning-and-ticket-status).
  - Mirroring Jira into the vault.
  - Mouse support.

# Acceptance criteria
- AC1: A missing sprint ritual, a due 1:1 and an unanswered team mail each show as a You decide row with a suggested action (063).
- AC2: Every Home row has a default enter action and the shared keys o, enter, x, s, d (064).
- AC3: A proposed meeting shows slot, attendees and body; nothing is created before the owner's OK (065).
- AC4: A 1:1 prep note lands in the dev's vault note before the meeting and the capture after it (066).
- AC5: Sprint planning prep lists velocity, carry-over and a proposal list without touching Jira (067).

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Execution order
1. item_063 checks engine (alone it already catches the sprint gap).
2. item_064 action model (generalises item_062; build 062 first or fold it in).
3. item_066 1:1 prep and capture.
4. item_067 sprint rituals.
5. item_065 calendar actions (needs Calendars.ReadWrite consent on the Entra app).

# Risks
- Calendars.ReadWrite sends real invites. Mitigation: draft first, create on OK only.
- Noise: too many checks bury the real ones. Mitigation: dismiss and snooze state, max rows per check.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-today`, `runbooks/graph-today.md`
- `tools/crow-deck/` (Next pane, You decide group)
- Vault: person notes, `90 Templates/Person (1-1).md`, `90 Templates/Retro.md`, `10 Work/6 Management/`
- `local/runbooks/sprint-planning-and-ticket-status.md`

# Backlog
- item_063_checks_engine_crow_checks_feeds_you_decide
- item_064_deck_action_model_every_row_actionable
- item_065_calendar_actions_propose_and_book_meetings
- item_066_one_on_one_prep_and_capture_per_dev
- item_067_sprint_rituals_planning_review_retro_prep

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
