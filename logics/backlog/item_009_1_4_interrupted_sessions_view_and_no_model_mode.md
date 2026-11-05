## item_009_1_4_interrupted_sessions_view_and_no_model_mode - 1.4 interrupted-sessions view and no-model mode
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub orchestration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: interrupted, sessions, view, model, mode
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Relaunching is a manual command and the owner has nothing to do while out of tokens.

# Scope
- In:
  - A crow watch view over the registry with relaunch, and a no-model list of useful work.
- Out:
  - Local model dispatch.

# Decisions
- 2026-11-05, owner, from rendered mockups (`local/state/shots/item-009/item-009-sheet.png`, generator `mock009.go.txt` there):
  - The interrupted-sessions view is the deck (Agents tab, `R` on a paused row), not crow watch.
  - Picker: the Relaunch overlay shows a table of every account (policy, 5h and week headroom bars with reset times, blocked reason below), j k picks, the recommended account is preselected. Blocked or never accounts cannot be picked.
  - Ask-first: choosing an ask-first account in the deck and pressing enter is the owner's explicit yes. The deck passes `--picked` to crow-relaunch, which accepts it in place of a typed yes for ask-first. `--yes` alone still does not count.
  - No-model mode: while every allowed account is blocked, Home Next shows a first group "No tokens until HH:MM" (earliest reset) listing work that needs no model: red pipelines, MRs waiting on the owner's review or with threads, repos with uncommitted files or stale worktrees, paused runs with reset times. Crow works rows dim. No new pane or key.

# Acceptance criteria
- AC1: the view fills the terminal like the other views and degrades on small terminals.
- AC2: selecting a session and an allowed account runs the guarded relaunch.

# Validation
- 2026-11-05: built in crow-deck (`tools/crow-deck/relaunch.go`), `go test ./...` passes, renders checked in all four themes (`local/state/shots/item-009/real/`).
  - AC1: the Relaunch overlay is a 72-column table drawn with the deck's box and gauge code, cut to the terminal width under 72; Home keeps its layouts, the no-model group is one more Next group. Shots at wide, portrait and small render.
  - AC2: `R` on a paused row opens the table with the recommended account preselected; j k skip blocked, never and the session's own account; enter runs `crow-relaunch <id> <target> --picked` (TestRelaunchPickerPicksAndPassesPicked). crow-relaunch accepts `--picked` as the yes for both policies (test_ask_first_hands_off_when_picked_in_the_deck, 21 relaunch tests pass).
  - No-model mode: the deck reads `crow-relaunch status --json` every 30 s; when every allowed account is blocked, Home Next opens with "No tokens until HH:MM" (MR rows of You decide, dirty repos, stale worktrees, paused runs) and Crow works shows as paused (TestNoTokens, TestHomeNoModelGroup). A week-capped account has no ISO reset time in the status, so the time is left out when only such accounts remain.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: AC1: the view fills the terminal like the other views and degrades on small terminals.
- request-AC6 -> This backlog slice. Proof: AC2: selecting a session and an allowed account runs the guarded relaunch.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_002_keep_working_when_an_account_runs_out`
- Architecture decision(s): (none yet)
- Request: `req_001_crow_account_switching_and_session_resume`
- Primary task(s): `task_002_deliver_account_switching_and_session_resume_in_three_steps`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.

# Notes
- Task `task_002_deliver_account_switching_and_session_resume_in_three_steps` was finished via `logics-manager flow finish task` on 2026-11-05.
