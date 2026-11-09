## item_116_resume_of_isolated_missions_from_files_and_cdx_issue_for_native_resume - Resume of isolated missions from files and cdx issue for native resume
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:28:44

# AI Context
- Summary: Isolated missions resume as a fresh run from the brief and version files. A cdx issue asks for native resume of isolated runs as a same-account optimisation.
- Keywords: resume, isolated, continue-run, cdx issue, fresh run
- Use when: Resuming a mission after an answer when the run was isolated.
- Skip when: Resume of non-isolated runs (cdx --continue-run).

# Problem
- cdx `--continue-run` cannot be isolated, and it cannot move to another account.
- Owner decision 2026-11-07: both. The fresh run from files is the guaranteed path (needed anyway when a mission moves account). A cdx issue asks for native resume of isolated runs.

# Scope
- In:
  - Fresh isolated run from `brief.md`, answers and version files is the default resume (shared with item_100).
  - cdx issue posted 2026-11-07 after the owner OK on the exact text: https://github.com/AlexAgo83/cdx-manager/issues/122.
  - When cdx ships it, use native resume on the same account and keep the fresh run as fallback.
- Out:
  - Posting the issue before the owner OKs the text.
  - Building native resume in cdx.

# Acceptance criteria
- AC1: An isolated mission answered by the owner resumes through a fresh run and the version file chain stays valid. Done 2026-11-09: `TestResumeFallsBackColdWhenContinueIsRefused` (workspace_gone, settings_mismatch, unavailable) starts a cold run from brief.md with resume_of and version 2.
- AC2: Done 2026-11-07: owner-approved issue posted (cdx-manager #122).
- AC3: Fresh run remains the path when the mission changes account. Done 2026-11-09: `TestResumeColdWhenTheMissionMovesAccount` (the caller passes cold; item_135 sets it from the item_115 move policy).

# AC Traceability
- request-AC1 -> AC1 to AC3. Proof: the acceptance criteria above, checked when built.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_030_spike_mission_safety_and_sandbox`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: Low
- Rationale: The fresh run works now; native resume only saves cache.

# Notes
- 2026-11-07 owner decision 7. Owner OK'd an issue in principle. Show the exact text before posting (upstream posts need approval).
- Related: item_100 (isolated runs), item_115 (limits), item_107 (brief).
- Cost after `--resume`: `total_cost_usd` is the running total of the whole Claude thread. Cost per turn = increase since the previous total of the same thread (or the total itself when it drops: new thread). Summing raw totals overcounts quadratically; check `cdx runs` cost for resumed runs.
- If `--resume` fails because Claude purged the transcript, start a new thread from the mission files and log a note; keep the displayed history.
- 2026-11-09 groomed into task_016: the cold resume path used by item_135. cdx issue #122 still open upstream.
- 2026-11-09 tool pass (owner picked): cdx 0.26.10 `cdx run --continue-run RUN_ID` resumes an isolated run in its own sandbox. 0.26.11 keeps the allowlist through a permissions-only `--settings` (crow adopted it in ac1d396 for review runs). Resume missions natively with it first; fall back to the file-based resume on `continue_run_workspace_gone` or `continue_run_settings_mismatch`.
- 2026-11-09 built as the shared resume primitive `env.resume` in tools/crow-core/internal/mission/resume.go. It rewrites brief.md, then tries `--continue-run` of the last run with a resume.md note quoting the new answers. A refusal with workspace_gone, settings_mismatch, unavailable or not_found falls back to a cold run from brief.md (`--resumes`) and logs an incident `continue_refused`. Other refusals start nothing. run_start records `cold`, `resume_of` and the version the brief asks for. Tick wiring, the resume cap and the end-to-end run stay in item_135.
- Isolating core-dispatched missions is out of this item: the tick dispatch has no isolate option yet. Suggest a separate item when slice 3 needs it.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
