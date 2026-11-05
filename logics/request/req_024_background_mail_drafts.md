## req_024_background_mail_drafts - Background mail reply drafts reviewed in the deck
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: deck
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Crow prepares French reply drafts in background for every mail the owner flags, from a rules file and the thread; the owner reviews, redrafts or opens them from a deck mail pane. Nothing is ever sent.
- Keywords: mail, outlook, graph, draft, reply, background, sonnet, deck, today, rules
- Use when: changing how flagged mail gets reply drafts or how the deck shows them.
- Skip when: calendar, sending mail, Teams.

# Needs
- Owner feedback 2026-11-05: enter on a flagged mail opened an interactive hub session. That is annoying and costly; the draft does not need a conversation.
- Wanted: drafts prepared in background from good rules and context, cheap, then one place in the deck to see where each mail stands and finish it.

# Context
- Today: `bin/crow-today` (Graph) lists flagged mail (`crow` true when crow flagged it), `thread <id>` prints the thread, `draft <id> [--all] < text` saves a reply draft (createReply). The Graph app has no Mail.Send: sending is impossible by design.
- Deck: Today pane rows, keys o/enter/x/s/r (item_062, 6d31aef). Enter runs `crow handoff cr-0w` (fixed in 45241b4).
- Headless runs: `crow_lib.RUN_KINDS` + `run_flags` (budget, timeout, model per kind), deny-by-default tool whitelist per kind; `crow-retro` + systemd user timer is the pattern for a periodic background job.

# Owner decisions (2026-11-05)
- Trigger: every mail the owner flags gets a draft. Mail crow flagged (category crow) does not: the owner decides, enter queues one.
- Review: a deck mail pane (summary + draft; open in Outlook, redraft with an instruction, drop).
- Model: Sonnet.
- Rules: new private `local/mail-rules.md` (general rules + per-sender notes), seeded from owner-prefs and context.md, corrected once by the owner.

# Acceptance criteria
- AC1: `local/mail-rules.md` exists, seeded and owner-approved; the drafter reads it on every run.
- AC2: Within one timer tick of the owner flagging a mail, a reply draft is in Outlook Drafts, made by a headless Sonnet run whose tools are only thread/draft/rules reads; nothing is sent; cost per mail is logged and capped.
- AC3: `crow-today json` reports per flagged mail its draft state (none, queued, drafting, ready, failed) with a short summary; `redraft <id> "<instruction>"` and `drop <id>` work.
- AC4: The deck shows the draft state on Today rows and a pane with thread summary and draft text; o opens the draft in Outlook, R redrafts with a one-line instruction, x drops; enter on a crow-flagged mail queues a draft instead of opening a session.
- AC5: Tests cover the picker (owner flag vs crow flag, already drafted), the state transitions and the deck keys.

# Out of scope
- Sending mail; auto-flagging rules; Teams replies; calendar.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Risks
- Draft quality depends on the rules file and thread context: the owner corrects mail-rules.md after the first drafts.
- A redraft must replace the previous crow draft, not pile up drafts in Outlook.
- Timer runs cost money with no one watching: per-mail budget cap, one run per mail unless redraft.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-today`, `bin/crow_lib.py`, `bin/crow-retro`, `tools/crow-deck/today_data.go`
- item_062 (Today mail actions)

# Backlog
- item_091_mail_rules_file
- item_092_background_mail_drafter
- item_093_crow_today_draft_state
- item_094_deck_mail_draft_pane

# Execution order
1. item_091 (rules) → 2. item_092 (drafter) → 3. item_093 (state + redraft/drop) → 4. item_094 (deck pane; mockup options shown to the owner first, can start in parallel with 2).


# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

