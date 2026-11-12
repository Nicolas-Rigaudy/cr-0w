## prod_009_nest_laptop_parity - Nest laptop parity
> Date: 2026-11-12
> Status: Proposed
> Related request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`
> Related backlog: `item_187_session_chat_in_the_laptop_desk_live_pane`, `item_188_start_an_interactive_session_that_lives_in_nest`, `item_189_several_session_chats_side_by_side_on_the_laptop`, `item_190_laptop_missions_from_the_approved_missions_board`, `item_191_laptop_inbox_layout`, `item_192_laptop_trial_day_without_terminals`
> Related task: `task_023_orchestrate_nest_laptop_parity`
> Related architecture: (none yet)
> Reminder: Update status, linked refs, scope, decisions, success signals, and open questions when you edit this doc.

# Overview
Nest is built for the Fold first. On the laptop the owner still runs sessions in terminals because the laptop Desk lacks the chat, cannot show several sessions at once and cannot start a session that stays inside Nest. This pass brings the laptop layout to the point where terminals become optional.

# Goals
- Every session readable and drivable from the laptop Nest.
- Several sessions on screen at once, like tmux panes.
- New sessions start inside Nest, with no desktop window.
- Missions and Inbox use the laptop width.

# Non-goals
- A full terminal emulator in the browser (raw TTY, colors, TUIs).
- Phone changes; the phone keeps the item_177 screens.
- Butler chat (item_176), push (item_178), mail and calendar (item_179).

# Scope and guardrails
- In: scaffolded request, product, backlog, orchestration task, validation, and handoff context.
- Out: unrelated workflow docs and implementation of generated tasks.

# Key product decisions
- Use structured input as the source of truth for generated docs.
- Keep generated write paths local and repo-bounded.

# Success signals
- Generated docs pass lint and audit without broad manual rewrites.
- Context-pack output can be handed to an implementation agent directly.

# References
- Product back-reference: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`
- Task back-reference: `task_023_orchestrate_nest_laptop_parity`
