## req_057_nest_on_the_laptop_replaces_the_owner_s_terminals - Nest on the laptop replaces the owner's terminals
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: High
> Theme: Nest
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: bring the laptop Nest layout up to the phone and Fold: session chat on the Desk, several chats side by side, headless sessions started from Nest, laptop Missions and Inbox.
- Keywords: nest, laptop, desk, chat, tmux panes, headless session, missions board, inbox
- Use when: building any Nest screen at 1200 px and wider, or driving sessions from Nest without a terminal window.
- Skip when: phone and Fold screens (req_054, req_055), butler chat, push, mail.

# Needs
- Owner, 2026-11-12: Nest work went to the phone and the laptop was left behind, so on the laptop the owner still works in terminals.
- Read live session output in Nest on the laptop, with full message text, as in a terminal.
- Type into any session and answer its prompts from Nest on the laptop.
- See several sessions at once, side by side like tmux panes.
- Start a new interactive session in a project from Nest, and drive it from Nest without a terminal window.
- Laptop layouts for Missions and Inbox, which today fall back to the Fold inner layouts.

# Context
- Nest switches to the laptop layout at 1200 px (LAPTOP in tools/nest/src/ui/useMedia.ts). Only Home and Desk have a laptop layout; Missions and Inbox render their Fold layouts.
- item_177 (commit 5fcb8ce) built the session and mission chat: core GET /api/sessions/{pid}/conversation and /api/missions/{id}/conversation through `crow-runs conversation` (cdx conversation pages, tail/cursor/before), the Nest Chat screen with tool folds, Earlier paging, send, permission answer buttons, Pause and Stop. It is used by the phone cover (/desk/<pid>) and the Fold inner Desk right pane. It polls every 3 s.
- The laptop Desk (DeskLaptop in Desk.tsx) still shows the old Live pane: Timeline, Ask and TellBox. It does not use the chat.
- Core actions already exist: POST /api/sessions (crow term in a project, optional account and prompt), POST /api/sessions/{pid}/send, /interrupt, /stop, /answer, /resume.
- `crow term` opens through CROW_TERMINAL=kitty|tmux|print. From core it opens a kitty window on the owner's desktop. Send and prompt answers reach a session through its terminal (tmux fallback from item_174).
- item_148 (run conversation poller with cursors in core) and item_146 (SSE stream) are the planned push path for conversation entries; until then the chat polls.
- Laptop Missions has approved boards: local/state/shots/nest-rebuild/boards/missions-board.png and missions-board-cyber.png (task_018 note 2026-11-10: kp title row, 5 columns Ready+Scheduled, Queued, Running n/2, Needs you+failed, Done today, bottom lane timeline and 7-day runs chart). No laptop Inbox board exists.
- Owner rules: visual choices are rendered and shown side by side before any build, one screen at a time with an OK before the next; build only from the board the item names as approved; never open test windows on the live desktop (headless screenshots, runbooks/nest-screenshots.md).
- Owner answers 2026-11-12: all four needs are wanted (live output, type into sessions, several sessions at once, start new sessions). Plan first, build later item by item.

# Acceptance criteria
- AC1: On the laptop the Desk shows the selected session or run as the item_177 chat (full message text, tool folds, Earlier, send, permission answers, Pause and Stop), checked against an approved laptop board.
- AC2: The owner types into any live session and answers its prompts from the laptop, and the reply shows in the chat without a page reload.
- AC3: Two or more session chats are visible side by side on the laptop, the owner picks which sessions fill the panes, and keyboard keys move focus and send.
- AC4: The owner starts a new interactive session in a registered project from Nest; it runs without opening a window on the desktop, appears in Nest, and its chat works at once. Opening the same session in a terminal stays possible.
- AC5: Missions on the laptop matches the approved missions-board boards on live data.
- AC6: Inbox on the laptop has its own layout, from a board the owner approved before the build.
- AC7: For a normal workday the owner can keep sessions in Nest on the laptop instead of terminals; a live trial day lists any gap left, and each gap becomes a backlog item.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)

# References
- tools/nest/src/screens/Chat.tsx
- tools/nest/src/screens/Desk.tsx
- tools/nest/src/screens/Missions.tsx
- tools/nest/src/screens/Inbox.tsx
- tools/nest/src/data/useConversation.ts
- tools/nest/src/ui/useMedia.ts
- tools/crow-core/internal/api/conversation.go
- tools/crow-core/internal/api/live.go
- bin/crow
- bin/crow-runs
- logics/request/req_055_a_work_day_from_nest_without_a_terminal.md
- logics/tasks/task_018_nest_app_shell_on_core_fixtures.md
- runbooks/nest-screenshots.md

# Backlog
- `item_187_session_chat_in_the_laptop_desk_live_pane`
- `item_188_start_an_interactive_session_that_lives_in_nest`
- `item_189_several_session_chats_side_by_side_on_the_laptop`
- `item_190_laptop_missions_from_the_approved_missions_board`
- `item_191_laptop_inbox_layout`
- `item_192_laptop_trial_day_without_terminals`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

