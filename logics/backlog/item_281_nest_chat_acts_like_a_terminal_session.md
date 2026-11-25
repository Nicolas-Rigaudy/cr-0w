## item_281_nest_chat_acts_like_a_terminal_session - Nest chat acts like a terminal session
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60%
> Confidence: 50%
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

# Options to settle in the plan
- A: the chat drives a real headless hub session (opus, project settings, hub hooks) and Nest shows it; the butler persona becomes its voice.
- B: keep the butler and widen it: read-only file and crow-script tools, hub context (owner prefs, open threads, context.md) in its prompt, model choice per message.
- Pick with the owner before any build. Write cost per day for each.

# Acceptance criteria
- AC1: Owner picks A or B from a short written comparison.
- AC2: The chat answers "what's on today", "status of !N", "stop X", "plan Y" as the hub session would, from the phone and the laptop.
- AC3: Every outward action still needs the same approvals as in the hub.

# Links
- Request: `req_068_a_full_work_day_from_nest_without_a_terminal`. Related item_176 (butler chat).

# Priority
- Priority: High
- Rationale: the no-terminal day failed on this.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
