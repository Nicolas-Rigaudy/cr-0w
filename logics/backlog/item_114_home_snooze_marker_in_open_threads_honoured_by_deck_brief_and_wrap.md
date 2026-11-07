## item_114_home_snooze_marker_in_open_threads_honoured_by_deck_brief_and_wrap - Home snooze marker in open threads honoured by deck brief and wrap
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Home `x` snooze writes a `snoozed until <date>` marker on the open-thread line. The deck, the brief and wrap hide the thread until that date.
- Keywords: snooze, open threads, deck home, brief, wrap
- Use when: Building the Home snooze key.
- Skip when: Mail snooze (item_062) and checks snooze (item_063).

# Problem
- Home `x` snooze is not built (item_020). Where a snooze is kept was open.
- Owner decision 2026-11-07: the marker lives in `local/state/open-threads.md` (or the vault equivalent). Deck, brief and wrap honour it.

# Scope
- In:
  - A `snoozed until YYYY-MM-DD` marker format on an open-thread line, documented in the open-threads runbook.
  - Deck Home `s`/snooze key writes the marker through the action gate.
  - Deck, brief and wrap skip a thread while its date is in the future and show it again on the date.
  - Stale detection ignores a snoozed thread.
- Out:
  - Mail flag snooze (item_062) and checks snooze (item_063).

# Acceptance criteria
- AC1: A snoozed thread disappears from deck Home, the brief and wrap until its date, then returns.
- AC2: The marker survives a wrap rewrite of open-threads.md.
- AC3: An unparsable date is shown as an error row, never hidden silently.

# AC Traceability
- request-AC1 -> AC1 to AC3. Proof: the acceptance criteria above, checked when built.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_014_crow_management_checks_and_actionable_deck`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: Small gap in the Home action model.

# Notes
- 2026-11-07 owner decision 1: open-threads marker, honoured by deck, brief and wrap.
- Depends on item_064 (action model) and item_052 (open-thread sections).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

