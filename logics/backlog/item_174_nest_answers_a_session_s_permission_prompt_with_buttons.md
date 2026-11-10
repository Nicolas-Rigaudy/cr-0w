## item_174_nest_answers_a_session_s_permission_prompt_with_buttons - Nest answers a session's permission prompt with buttons
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
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
- request-AC2 -> AC1, AC2, AC3, AC4. Proof: implemented in bfa3141, d248562, 19800ca, 89fb516, 416f8c6; validated with pytest (617 passed), go test ./... in tools/crow-core, vitest (26 passed) and the live tmux test below. Source: `89fb516`

# Spike decision (2026-11-10)
- Source: a crow `PermissionRequest` hook (`bin/crow-prompt hook`) in every cdx profile's user settings. It writes `local/state/prompts/<pid>.json`, pid being the claude process in the hook's parent chain.
- Payload (Claude Code 2.1.322): session_id, transcript_path, cwd, permission_mode, prompt_id, tool_name, tool_input, permission_suggestions. No tool_use_id. The hook prints nothing: any stdout would decide the prompt.
- Pending check: a recorded prompt counts only while `<HOME>/.claude/sessions/<pid>.json` says status `waiting`, waitingFor `permission prompt`, same sessionId. A new prompt overwrites the file and gets a new id.
- Keys: menu is 1 Yes / 2 always allow dir / 3 auto mode / 4 No. Allow once = `1`, deny = Escape (the turn stops). Always-allow is never offered.
- Rejected: `cdx notify` (drops tool_input); reading the pane (fragile, no stable id).

# Delivery
- bfa3141 `bin/crow-prompt` (hook, show, answer), `crow_lib.pending_prompt`, `crow_term.send_key`.
- d248562 `crow-runs json` adds `prompt` to a waiting agent.
- 19800ca core `POST /api/sessions/<pid>/answer {prompt, option}`: 409 when stale or on a work allow; one event per answer.
- 89fb516 Nest layout B: Needs row "Answer" opens the live panel with the prompt above the timeline; work approval shows only deny plus "give it on the exact text in approvals".
- 416f8c6 `crow doctor` and the new-account runbook require the hook.

# Live test (2026-11-10, hidden tmux session, haiku, default mode)
- AC1: prompt file 2 s after the request; `crow-prompt show` and `crow-runs json` show `date +%s > stamp.txt`.
- AC3: a stale id is refused (exit 1); no file written.
- AC2: allow ran the command (stamp.txt written); the prompt cleared.
- Deny: Escape, `touch denied.txt` never ran, prompt cleared, session back at input.
- AC4: covered by tests (`tests/test_crow_prompt.py`, `src/data/live.test.ts`, core fixtures).

# Decision framing
- Product framing: Not needed (prod_005 human in the loop)
- Architecture framing: Consider. The prompt source is recorded under Spike decision.

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
