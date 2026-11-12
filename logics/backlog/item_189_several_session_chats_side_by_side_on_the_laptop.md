## item_189_several_session_chats_side_by_side_on_the_laptop - Several session chats side by side on the laptop
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: several, session, chats, side, laptop
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- In terminals the owner keeps several sessions visible at once; Nest shows one at a time.

# Scope
- In:
  - Render layout options (2, 3, 4 panes, how panes are picked, where the slot list goes) and get the owner's pick before the build.
  - A laptop view that tiles chats of chosen sessions; the pane choice survives a reload.
  - Keyboard: move focus between panes, focus the input, send, Esc; shown in a key hint.
- Out:
  - Resizable splits by drag unless the chosen board asks for them.

# Acceptance criteria
- The owner approved a rendered multi-pane board before the build.
- Two or more chats show at once, each reading and sending on its own session.
- Keyboard focus moves between panes and sends without the mouse.
- Polling cost stays bounded with several panes (one read per pane per interval, hidden page stops reading).

# AC Traceability
- request-AC3 -> This backlog slice. Proof: The owner approved a rendered multi-pane board before the build.

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

