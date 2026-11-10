## item_175_work_approvals_given_from_nest_on_the_exact_text - Work approvals given from Nest on the exact text
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 70%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: a session that needs the owner's OK for a work push, MR, mail or Teams post files an approval request in the core inbox with the exact text; Nest shows it in full with Approve and Reject; the session waits and acts only on Approve.
- Keywords: approval, work hard rules, push, mr, mail, teams, exact text, inbox, nest
- Use when: building approvals or their request command.
- Skip when: tool permission prompts (item_174).

# Problem
- Work rules: push, MR, mail and Teams always need an OK on the exact text (CLAUDE.md, prod_005). Today the OK is only possible in the session's terminal. The owner wants to give it from Nest (2026-11-10).

# Scope
- In:
  - Reuse the core inbox (adr_001): a new item kind `approval` with `{action, target, text, diff_ref, session}`. No second queue.
  - `crow approve request --action push|mr|mail|teams --target T --text-file F` posts it over the socket and blocks until answered or timed out; prints the answer as JSON. Sessions call it before the action; the push and MR helpers and `crow-hub-guard` learn to call it.
  - Nest: the item shows the full text (mail body, MR description, branch and commit list for a push), Approve and Reject, and Reject with a note. Approve is one tap after the full text was shown; no bulk approve.
  - The answer is an event with actor, surface and the hash of the shown text. The action runs only if the text hash still matches.
  - Terminal answer stays possible (deck and the session), first answer wins.
- Out:
  - Agents holding work credentials (never). Auto-approve. Merge (never on work).

# Acceptance criteria
- AC1: A test session's push request appears in Nest with the branch, commits and remote; Approve lets the push run; Reject stops it with the note.
- AC2: If the text changes after it was shown, Approve gets 409 and nothing runs.
- AC3: Each answer leaves an event with the text hash.
- AC4: A mail draft is shown in full before Approve is enabled.

# AC Traceability
- request-AC2 -> AC1, AC2. Proof: work approvals from the phone under the hard rules.

# Decision framing
- Product framing: Not needed (prod_005 hard rules)
- Architecture framing: Consider. Approval as an inbox kind fits adr_001; confirm during the build.

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_054_nest_sees_and_pilots_everything_deck_does`
- Depends on: `item_147_core_action_endpoints_for_dispatch_answer_stop_cancel_and_hide`, `item_168_nest_answer_sheet_for_inbox_items_and_mission_actions`
- Related: `item_151_autonomy_dial_per_kind_and_mission_under_the_hard_rules`, `item_174_nest_answers_a_session_s_permission_prompt_with_buttons`
- Primary task(s): `task_020_orchestrate_nest_going_live_on_the_core_api`

# Priority
- Priority: Medium
- Rationale: owner 2026-11-10: plan now. Builds on the inbox actions and answer sheet.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
