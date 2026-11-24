## item_273_ops_lane_sheet_run_changes_and_commits_tabs_from_the_worktree - Ops lane sheet: Run, Changes and Commits tabs from the worktree
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-24 17:24:23

# AI Context
- Summary: The Ops lane sheet gets tabs Run / Changes / Commits. The header shows the worktree path and diff total; Changes lists branch, ahead/behind main, uncommitted count and per-file +/−.
- Keywords: ops, lane, sheet, run, changes, commits, tabs, worktree
- Use when: Building the lane sheet changes view or its core route.
- Skip when: Lanes with no git checkout (Changes and Commits tabs hide).

# Problem
- The lane sheet (LaneSheet in tools/nest/src/screens/Ops.tsx) shows props and the run report only. The owner cannot see which worktree a run used or what it changed.
- Design approved 24/11: vault `CR-0W/Research/nest-design/item_190/lane/compare.html`, option C; decisions in `item_190/design-notes.md`.

# Scope
- In:
  - Core: a changes route per lane, read on demand when the sheet opens. Runs git in the lane cwd (ops lanes carry Cwd): branch, ahead/behind main, numstat per file, uncommitted count, commits on top of main (sha, subject, time). Bounded output and a short timeout. Finished runs use the diff stat item_269 stores on the ledger row (their worktree may be gone); only live lanes read git.
  - Nest: header line with worktree path and diff total (files, +/−). Tabs Run (props and report, default), Changes (branch card, one line per file with +/− and a bar, paths cut at the start), Commits. Buttons OPEN, COPY PATH, CLOSE. Laptop side sheet and phone bottom sheet.
- Out:
  - Full diff text. Actions on the worktree (remove, rebase, MR).

# Acceptance criteria
- AC1: For a run in a worktree, Changes and Commits match `git diff --numstat main...` and `git log main..` in that cwd.
- AC2: A lane with no repo cwd shows the Run tab only.
- AC3: Sheet matches lane/compare.html option C on laptop and 368x812, both themes.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: implemented in 31b13ed (core GET /api/changes) and 7744c01 (Nest lane sheet). TestChanges compares the route with git diff --numstat main... and git log main.. in a worktree (AC1). The sheet hides the tabs when no git data exists (AC2). The owner approved built-273.png against lane-C on 2026-11-24 (AC3).

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)
- Request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`, `req_063_nest_pages_and_inbox_as_the_home_for_everything_that_waits`
- Primary task(s): `task_023_orchestrate_nest_laptop_parity`
- Depends on: `item_190_laptop_missions_from_the_approved_missions_board`, `item_269_run_report_shows_the_diff_stat`

# Priority
- Priority: Medium
- Rationale: Finishes the Ops page (item_190). Build order: item_272, item_273, item_274; each in a fresh session.

# Notes
- Read the vault design notes and the compare page before building; build from the chosen option only.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
