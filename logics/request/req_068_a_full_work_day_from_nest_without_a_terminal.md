## req_068_a_full_work_day_from_nest_without_a_terminal - A full work day from Nest without a terminal
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 80%
> Complexity: High
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: The owner tried a day with Nest only (2026-11-25) and went back to the terminal after 30 minutes. The butler chat relaunched a stale request, could not cancel it, and Stop in Nest did nothing. This request covers what Nest needs so a work day runs without a terminal.
- Keywords: nest, butler, chat, stop, cancel, mission, account, terminal parity
- Use when: making Nest chat and controls enough for a full day.
- Skip when: Nest visuals, layout or theme work.

# Needs
- Stop and cancel do what they say, from the chat and from the buttons.
- The chat knows what the hub session knows and can do what it does (read files, run crow scripts, check state) within the hub rules.
- The chat never acts on an old message as if it were new.
- Missions run only on accounts the dispatch policy allows.

# Context
- 2026-11-25 10:23: the butler created mission review lambda!242 from a message sent the day before. The review was already done in a terminal session.
- The butler had no cancel tool and told the owner to cancel in Nest. Stop in Nest did nothing: no stop event in the mission log.
- Mission Stop only kills the live run; the next tick resumes it (resume cap 2), so even a working Stop looks like nothing happened.
- The plan mission nest-desk-bugs ran on an account outside the dispatch policy.
- Fixed the same day (650b35b): butler tools mission_cancel and run_stop; persona acts only on the newest message.

# Acceptance criteria
- AC1: The owner works a full day from Nest (laptop and phone) without opening a terminal.
- AC2: Any mission or run can be stopped from the chat or a button, and stays stopped.

# Backlog
- item_279 Nest Stop on a mission run stops the mission
- item_280 Core missions pick an account by dispatch policy
- item_281 Nest chat acts like a terminal session
- Related: item_210, item_277, item_278 (desk bugs from the same day)
