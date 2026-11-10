## item_174_nest_answers_a_session_s_permission_prompt_with_buttons - Nest answers a session's permission prompt with buttons
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 75%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: when a live Claude session stops on a permission prompt, core knows the exact prompt (tool, command, options) and Nest shows it with one button per option; the tap answers the prompt in that session.
- Keywords: permission prompt, notification hook, waiting session, allow, deny, nest buttons, crow-focus
- Use when: answering a session's tool permission from the phone.
- Skip when: free text answers (item_170 send), work approvals (item_175).

# Problem
- deck shows a session as "asks you" (`crow-runs json` agent state `waiting`) but not what it asks. A permission prompt needs a key, not text. The owner wants to answer it from the phone early (2026-11-10), before the conversation feed (req_027).

# Scope
- In:
  - Spike first (half a day): find the reliable source of the pending prompt. Candidates: the Claude Code `Notification` hook (permission prompt event) writing `local/state/prompts/<session>.json`; reading the pane with kitty or tmux. Pick one, note why in this doc.
  - `crow-runs json` (or `crow-sessions`) adds `prompt: {tool, detail, options[]}` to a waiting agent; contract per item_173.
  - Core: `POST /api/sessions/<pid>/answer {option}` runs the script that presses the option's key (`crow-focus --pid N --key K --json`). 409 when the session no longer waits on that prompt.
  - Nest: the waiting session row shows the prompt text and one button per option. Deny is plain; allow on a work repo push or MR is refused here (item_175 handles those).
  - An event per answer, as item_170.
- Out:
  - Prompts in sessions outside kitty or tmux. Changing permission rules.

# Acceptance criteria
- AC1: A session waiting on a Bash permission shows the exact command in Nest within 10 s.
- AC2: Tapping an option answers the prompt; the session carries on; the row leaves "asks you".
- AC3: An answer to a prompt that already changed gets 409 and presses nothing.
- AC4: A prompt to push or open an MR on a work repo shows only "Open in approvals" (item_175), never allow.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: pilot a waiting session from the phone.

# Decision framing
- Product framing: Not needed (prod_005 human in the loop)
- Architecture framing: Consider. The prompt source choice is recorded here after the spike.

# Links
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Depends on: `item_170_core_actions_on_live_sessions_and_dispatch_runs`, `item_173_deck_scripts_get_a_stable_json_contract_for_the_app`
- Related: `req_027_conversations_feed_for_core_deck_nest`, `item_175_work_approvals_given_from_nest_on_the_exact_text`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: High
- Rationale: owner 2026-11-10: answer prompts from the phone early.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
