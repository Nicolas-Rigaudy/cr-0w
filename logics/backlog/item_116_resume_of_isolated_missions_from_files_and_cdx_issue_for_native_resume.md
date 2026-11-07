## item_116_resume_of_isolated_missions_from_files_and_cdx_issue_for_native_resume - Resume of isolated missions from files and cdx issue for native resume
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:14:58

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
- AC1: An isolated mission answered by the owner resumes through a fresh run and the version file chain stays valid.
- AC2: Done 2026-11-07: owner-approved issue posted (cdx-manager #122).
- AC3: Fresh run remains the path when the mission changes account.

# AC Traceability
- request-AC1 -> AC1 to AC3. Proof: the acceptance criteria above, checked when built.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_030_spike_mission_safety_and_sandbox`
- Primary task(s): (none yet)

# Priority
- Priority: Low
- Rationale: The fresh run works now; native resume only saves cache.

# Notes
- 2026-11-07 owner decision 7. Owner OK'd an issue in principle. Show the exact text before posting (upstream posts need approval).
- Related: item_100 (isolated runs), item_115 (limits), item_107 (brief).
- Cost after `--resume`: `total_cost_usd` is the running total of the whole Claude thread. Cost per turn = increase since the previous total of the same thread (or the total itself when it drops: new thread). Summing raw totals overcounts quadratically; check `cdx runs` cost for resumed runs.
- If `--resume` fails because Claude purged the transcript, start a new thread from the mission files and log a note; keep the displayed history.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
