## item_124_render_the_checks_action_options_on_deck_home - Render the checks action options on deck Home
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Render deck Home with the three ways to show checks actions so the owner can pick one.
- Keywords: render, checks, action, options, deck, home
- Use when: Choosing how checks actions fit in deck Home's Why column.
- Skip when: Checks engine logic.

# Problem
- Deck Home's Why column is 13 characters, so check actions like 'propose planning' get cut. The owner wants renders before choosing.

# Scope
- In:
  - Three PNG renders of deck Home with the same data (the 2026-11-03 fixture) and framing: (a) wider Why column, (b) one-word actions, (c) action shown on the selected row only.
  - Rendered in the terminal ANSI theme, as the deck looks today.
  - A bin/crow-compare page with the three renders side by side.
- Out:
  - Shipping the chosen option (follow-up once picked).

# Acceptance criteria
- AC1: Three PNGs, one per option, same size and data.
- AC2: A compare page opens them side by side.
- AC3: No option is picked by the agent.

# AC Traceability
- AC1 -> this item. Proof: to record at closeout.
- AC2 -> this item. Proof: to record at closeout.
- AC3 -> this item. Proof: to record at closeout.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Unblocks the checks rows on Home.

# Notes
- Owner decision 2026-11-07 on return from the AFK run (local/state/afk/decisions-2026-11-07-afk.md).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

