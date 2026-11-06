## req_026_personal_butler - Personal butler on the personal account: mail, calendar, budget, health
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 75%
> Confidence: 70%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 16:13:51

# AI Context
- Summary: Grow crow from a work chief of staff into a personal butler on the personal account: personal mail, calendar, budget and banking, health, home admin. Draft: needs grooming into slices before any build.
- Keywords: butler, personal, gmail, calendar, budget, banking, health, assistant
- Use when: planning personal-life features of crow.
- Skip when: work rituals (brief, daily, weekly) on the work account.

# Needs
- Personal connectors on the personal account only (mail, calendar, drive), never on shared or work accounts.
- Personal rituals: a personal brief (mail triage, agenda, bills due, budget drift), a weekly personal review.
- Budget and banking: import the bank's CSV export, categorise, monthly budget vs actual, alerts.
- Health: appointment reminders and follow-ups, global only.
- Same autonomy model as work: read and draft freely, send or pay only with explicit OK.

# Context
- Depends on req_025 (personal scope routes to the personal account).
- Data placement: personal data stays in the personal vault or local/, never in the publishable engine.
- Hub skills already split generic engine (skills/, runbooks/) from private config (local/); personal rituals follow `runbooks/own-routines.md`.

# Decisions (owner, 2026-11-06)
- Personal state lives in the personal vault (budget reviews, appointment log, personal brief history). Raw imports may sit next to it in the vault's attachments; nothing in the publishable engine.
- Health: global, not detailed. Appointment reminders and follow-ups (who, when, what to bring, next check-up due), no sensor or wearable data.
- Banking: start from the bank's own export (CSV, rolling 12 months, "download my operations" on the web account). No third party, no stored bank credentials. The importer is source-agnostic so a DSP2 aggregator (French providers such as Powens or Bridge, consent renewed every 90 to 180 days) can replace the manual export later if downloading becomes a chore. Bank name and account specifics in local/context.md.

# Open questions (grooming)
- Which personal connectors exist on the personal account (Gmail, Calendar, Drive) once connected, and which are read only?
- Privacy rules: what may appear in the deck, the brief and phone notifications (amounts, appointment details).
- Bank export check: confirm the export button and format on the owner's own account page before the budget slice.

# Acceptance criteria
- AC1: Grooming splits this into backlog slices (connectors and personal brief, appointments and reminders, budget from bank exports) with their own AC.
- AC2: Each slice states its data location and its approval points.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics/request/req_027_conversations_feed_for_core_deck_nest.md` (rk1 conversations as the butler's memory of requests and follow-ups)
- `runbooks/own-routines.md`
- `logics/request/req_025_account_routing_by_scope.md`

# Backlog
- `item_096_personal_brief_from_mail_and_calendar`
- `item_097_appointment_reminders_and_followups`
- `item_098_monthly_budget_from_bank_exports`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
