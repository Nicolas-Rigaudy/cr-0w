## item_281_nest_chat_acts_like_a_terminal_session - Nest chat acts like a terminal session
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 70%
> Progress: 0%
> Complexity: High
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The butler chat runs on Sonnet with only a few core MCP tools and a short persona. It cannot read files, run crow scripts or see what the hub session sees, so the owner goes back to the terminal. Bring the chat close to a hub session within the hub rules.
- Keywords: nest, butler, chat, parity, hub session, tools, context
- Use when: what the Nest chat can know and do.
- Skip when: chat layout.

# Problem
- The butler is `claude -p --model sonnet --tools '' --allowedTools mcp__core__*`, persona only. No CLAUDE.md, owner prefs, open threads, runbooks or skills.
- Tools: brief, usage, mission_get/new/cancel, run_stop, session_new, dispatch_run. No file read, no scripts, no Jira or forge.
- It treated a message from the day before as a new order (persona fixed 650b35b).

# Decision (owner, 2026-11-25)
- A, with the butler kept. The chat header has an always visible switch between two modes:
  - Butler: the current Sonnet butler for quick status and simple actions.
  - Hub: a real headless hub session (opus, project settings, hub hooks, owner prefs, open threads, runbooks, skills) that Nest shows and drives.
- Each mode keeps its own conversation. The switch opens in the last mode used, remembered per device.
- B (widening the butler) is dropped.
- Split into build items before any work. Write the cost per day of Hub mode in the plan.

# Acceptance criteria
- AC1: The chat header switches between Butler and Hub in one tap, and Nest reopens on the last mode used on that device.
- AC2: In Hub mode the chat answers "what's on today", "status of !N", "stop X", "plan Y" as the hub session would, from the phone and the laptop.
- AC3: Every outward action still needs the same approvals as in the hub.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`. Related item_176 (butler chat).

# Priority
- Priority: High
- Rationale: the no-terminal day failed on this.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
