## item_091_mail_rules_file - Mail drafting rules file seeded from prefs and context
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
- Summary: Private local/mail-rules.md: tone, language, sign-off, never-say, per-sender notes; seeded from owner-prefs and context.md, approved by the owner.
- Keywords: mail, rules, tone, sign-off, sender, prefs
- Use when: building this slice of req_024_background_mail_drafts.
- Skip when: other mail work.

# Problem
- Background drafts need written rules; today they live scattered in owner-prefs, context.md and memories.

# Scope
- In:
  - `local/mail-rules.md`: general rules (French for work, short declarative, sign-off, no em-dashes), what never goes in a mail, per-sender/per-domain notes (people from context.md).
  - Generic template `templates/mail-rules.md` (no names) so the engine stays publishable.
  - Owner reviews the seed once; corrections go in the file.
- Out:
  - Drafter code (item_092).

# Acceptance criteria
- AC1: The file exists with general rules and per-sender notes, owner-approved.
- AC2: The template outside local/ has no private names (crow-precommit-check passes).

# Dependencies
- none

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

