## req_026_personal_butler - Personal butler on the personal account: mail, calendar, budget, health
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
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
- Budget and banking: import statements or a read-only aggregator, categorise, monthly budget vs actual, alerts.
- Health: log and trends from exported data (sleep, activity, weight, appointments).
- Same autonomy model as work: read and draft freely, send or pay only with explicit OK.

# Context
- Depends on req_025 (personal scope routes to the personal account).
- Data placement: personal data stays in the personal vault or local/, never in the publishable engine.
- Hub skills already split generic engine (skills/, runbooks/) from private config (local/); personal rituals follow `runbooks/own-routines.md`.

# Open questions (grooming)
- Which connectors exist for the personal account (Gmail, Google Calendar, Drive) and which are read only?
- Banking: statement exports (CSV/OFX) or an aggregator API? Which banks?
- Health sources: which apps or devices export data?
- Where personal state lives: the personal vault or local/state?
- Privacy rules: what may appear in the deck, the brief, notifications.

# Acceptance criteria
- AC1: Grooming splits this into backlog slices (connectors, personal brief, budget, health) with their own AC.
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
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
