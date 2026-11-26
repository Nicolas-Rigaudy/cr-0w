## item_298_session_start_fails_loudly_and_opens_on_desk - Session start fails loudly and opens on Desk
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 90%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-26 14:36:14

# AI Context
- Summary: A Nest session start reports a session that died, and a live one opens on Desk.
- Keywords: session, start, fails, loudly, opens, desk
- Use when: a Nest start says started but no session shows.
- Skip when: unrelated Nest pages.

# Problem
- Core answers 200 to `POST /api/sessions` as soon as tmux exists. `crow open --on-busy fail` can exit at once (busy account or project, cdx refusal) and the tmux session goes with it. Nest says "started" and nothing shows. Seen 11-26 10:39 (lambda x2) and 10:41 (cr-0w, exited after 1.4 s).

# Scope
- In:
  - Core keeps the pane after exit (remain-on-exit), waits up to ~3 s, and on exit answers an error with the exit code and the last lines of the pane, then kills the session.
  - Nest start actions (Desk launcher, MRs Repo session, threads, branch) show that error in place.
  - On success Nest opens the new session on Desk (navigate to /desk/<pid> once the pid is known).
- Out:
  - Changing the on-busy policy itself.

# Acceptance criteria
- AC1: Reproduce first: find why the 11-26 starts exited (run `crow open --on-busy fail lambda` the way core does, capture output). Record the cause in Notes.
- AC2: A start that exits within the wait returns non-2xx with the reason; a Go test covers it with a fake command that exits.
- AC3: Nest shows the reason next to the action; vitest covers it.
- AC4: A start that lives lands on that session in Desk.

# AC Traceability
- req_071 -> this slice. Proof: its acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed
- 26/11 owner: busy account start goes through a new cdx --on-busy allow mode (handed to cdx-manager). Until it lands, the Desk New session note says a busy account fails.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_071_nest_mrs_review_flow_and_visible_session_starts
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Reviews rank second in the owner's priorities; no-terminal days depend on Nest.

# Notes
- 2026-11-26 14:22 local: a fourth case. Nest start of the hub project (tmux pane named project-account-time) exited after 1.4 s; cdx `launch_history.jsonl` has no launch on that account at that time and no transcript exists, so it failed before claude started. The previous session on that account had ended 85 s earlier (13:20:39Z). Suspect `--on-busy fail` (cdx `session_busy`) or a quota refusal; the pane output was lost. owner-prefs say several live sessions on one account are expected and a launch must never be refused for a busy account: check whether `--on-busy fail` from core contradicts that once the cause is known.
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_034_session_start_fails_loudly_and_opens_on_desk`
