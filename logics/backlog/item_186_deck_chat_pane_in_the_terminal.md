## item_186_deck_chat_pane_in_the_terminal - Deck chat pane in the terminal
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: crow-deck
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A Chat pane in the deck that shows the owner's prompts and crow's full replies (tool calls folded) for the selected lane, with the ask card docked under it. Terminal only.
- Keywords: deck, chat, pane, transcript, reply, said, split
- Use when: building or changing how the deck shows a session's conversation.
- Skip when: the phone (Nest, item_177), or the lane list and ruler.

# Problem
- The deck lane tail shows crow's replies flattened to one line by `said()` in bin/crow-runs (text whole since f8612bd, up to 2000 chars), so bullets and paragraphs run together and a long reply is cut at "… N more lines".
- The owner wants to read whole replies without diffs or tool logs (owner, 2026-11-11).
- The deck runs in the owner's terminal only. The phone uses Nest (owner, 2026-11-11), covered by item_177.

# Design (owner picks 2026-11-11: option A, fitted to the _ux2 chrome, split layout)
- Mockups: vault `Research/Deck chat options.md`, section "Pick: A, fitted into the deck design (v2)". Images in `Research/deck-chat/` (v2-inner-split.png, v2-inner-split-ask.png for wide, v2-cover-chat.png, v2-cover-ask.png for narrow), compare page compare2.html. Build from v2 only.
- The tail pane becomes Chat (pane 2 of 4, same `pane N of 4 · tab` corner).
- Chat: `you` (yellow) and `crow` (green) headers, each message wrapped with its line breaks kept, bullets and numbered lists with a hanging indent. Strip `**` and backticks.
- Each run of tool calls between two messages folds into one grey line `· N tool calls hidden`. `t` toggles them shown.
- Newest message at the bottom, follows new output unless the owner scrolled. `j k` scroll.
- When the lane waits on the owner, the ask card docks under the chat; `enter` sends the typed reply (same path as the ask card answer).
- Wide terminal (100 cols and more): lanes left, Chat right (split). Narrow: one pane at a time, switched with `tab` as today.

# Scope
- In:
  - bin/crow-runs: add a `messages` list (role you/crow, full text with line breaks, tool count before each message). Keep `recent`/`last_action` as they are for crow-watch and Nest.
  - tools/crow-deck: Chat pane replacing tail, split at wide sizes, scroll, tool toggle, docked ask card.
  - Tests for wrapping with line breaks, the tool fold and both layouts; TestShots renders compared to the v2 mockups before showing the owner.
- Out:
  - Markdown styling beyond line breaks, bullets and lists.
  - Nest and the phone (item_177).

# Acceptance criteria
- AC1: A 1400-char reply with bullets reads whole with its line breaks in the Chat pane, scrollable with j k.
- AC2: Tool calls show as one grey line per run; `t` shows them.
- AC3: With a waiting lane, a reply typed in the docked ask card reaches the session.
- AC4: At 100+ cols lanes sit left and Chat right, as in v2-inner-split.png; narrower, `tab` switches panes.
- AC5: Renders match the v2 mockups (owner sees a side-by-side before commit).

# AC Traceability
- AC1 -> crow-runs messages + chat wrap tests.
- AC2 -> tool fold test.
- AC3 -> reply path test.
- AC4, AC5 -> TestShots renders vs v2 mockups.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control`
- Related: `item_177_chat_inside_any_session_or_mission_run_from_nest` (same message shape; when item_148 lands, read conversations from core instead of crow-runs)
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Main way the owner reads crow at the desk without the full log.

# Notes
- Follows f8612bd (tail wrap, whole messages up to 2000 chars).
- UI rebuild rule: build the pane, show it against the mockup, wait for OK.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
