## item_187_session_chat_in_the_laptop_desk_live_pane - Session chat in the laptop Desk live pane
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: session, chat, laptop, desk, live, pane
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- The laptop Desk Live pane shows the old timeline and tell box; the item_177 chat only appears on the phone and the Fold inner Desk.

# Scope
- In:
  - Render options for the laptop Live pane with the chat (where Ask CR-0W goes, Chat/Log toggle, header actions) and get the owner's pick before the build.
  - DeskLaptop uses the Chat component for the selected agent; send and permission answers work from the laptop.
  - Headless screenshots vs the approved board.
- Out:
  - Several panes (next item).
  - Moving the chat from polling to item_148/146.

# Acceptance criteria
- The owner approved a rendered laptop Desk board with the chat before the build.
- Selecting a slot on the laptop Desk shows its chat with full message text, Earlier, send and permission answers.
- A message sent from the laptop reaches the live session and its reply appears without reload (live test on a real session).
- Nest tests and tsc pass; a laptop screenshot is compared with the board.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: The owner approved a rendered laptop Desk board with the chat before the build.
- request-AC2 -> This backlog slice. Proof: Selecting a slot on the laptop Desk shows its chat with full message text, Earlier, send and permission answers.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)
- Request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`
- Primary task(s): `task_023_orchestrate_nest_laptop_parity`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

