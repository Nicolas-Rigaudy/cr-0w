## item_008_1_3_guarded_relaunch_and_block_detection - 1.3 guarded relaunch and block detection
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub orchestration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-05 15:23:56

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: guarded, relaunch, block, detection
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- A handoff without an id picked the wrong session, and nothing watches for a blocked account.

# Scope
- In:
  - crow command that checks `cdx can-resume`, requires an id and a confirmed allowed target, then runs the handoff; block detection from `cdx status --json`.
- Out:
  - The TUI.

# Acceptance criteria
- AC1: the command refuses to run without a conversation id or with a never account.
- AC2: an ask-first account needs an explicit yes every time.

# Validation
- 2026-11-05: AC1 and AC2 proven in task_002 (2026-10-20): bin/crow-relaunch refuses no id, bad ids, never or unknown targets; ask-first needs a typed yes on a terminal, --yes and piped input refused. `crow relaunch status` recommends an unblocked account. Full suite re-run: 273 tests pass.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1: the command refuses to run without a conversation id or with a never account.
- request-AC4 -> This backlog slice. Proof: AC2: an ask-first account needs an explicit yes every time.
- request-AC7 -> This backlog slice. Proof: AC2: an ask-first account needs an explicit yes every time.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_002_keep_working_when_an_account_runs_out`
- Architecture decision(s): (none yet)
- Request: `req_001_crow_account_switching_and_session_resume`
- Primary task(s): `task_002_deliver_account_switching_and_session_resume_in_three_steps`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Notes
- Task `task_002_deliver_account_switching_and_session_resume_in_three_steps` was finished via `logics-manager flow finish task` on 2026-11-05.
