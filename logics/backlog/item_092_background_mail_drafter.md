## item_092_background_mail_drafter - Background reply drafter for owner-flagged mail
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: deck
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Headless Sonnet run per owner-flagged mail: reads the thread and mail-rules.md, saves a reply draft via crow-today draft; timer-driven, capped, never sends.
- Keywords: mail, draft, background, sonnet, headless, timer, budget
- Use when: building this slice of req_024_background_mail_drafts.
- Skip when: other mail work.

# Problem
- Enter on a mail opens an interactive session; the owner wants drafts ready without a conversation.

# Scope
- In:
  - `bin/crow-maildraft`: picks flagged mails with crow=false and no crow draft yet (plus queued ids); one headless run per mail via `run_flags` with a new RUN_KIND `maildraft` (Sonnet, small budget, short timeout).
  - Tool whitelist: `bin/crow-today thread`, `bin/crow-today draft`, Read of local/mail-rules.md, owner-prefs, context.md. Nothing else.
  - State in `local/state/maildrafts.json` (id → state, draft id, summary, cost, at).
  - systemd user timer (crow-retro.timer pattern), e.g. every 10 minutes in working hours; lock file so runs never overlap.
  - Runbook entry in runbooks/dispatch-background-runs.md.
- Out:
  - Deck UI (item_094).
  - Sending (impossible: no Mail.Send).

# Acceptance criteria
- AC1: Flagging a mail in Outlook yields a reply draft in Drafts within one tick, nothing sent.
- AC2: A crow-flagged mail gets no draft unless queued.
- AC3: Each run logs its cost; the per-mail cap holds; a mail is drafted once unless redrafted.
- AC4: Tests cover the picker and the lock.

# Dependencies
- item_091

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_024_background_mail_drafts
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: owner request 2026-11-05 after the first live mail session.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

