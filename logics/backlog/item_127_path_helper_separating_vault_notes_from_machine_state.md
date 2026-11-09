## item_127_path_helper_separating_vault_notes_from_machine_state - Path helper separating vault notes from machine state
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub data
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 10:46:09

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: path, helper, separating, vault, notes, machine, state
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Every script derives one state dir from vault.crow_dir, so machine and owner files cannot be split.

# Scope
- In:
  - crow_lib: machine_dir() always local/state; crow_dir() = vault.path joined with vault.crow_dir, else local/state; crow_path(kind, ...) maps kinds (threads, status, decisions, rollup, draft, research, retro, daily) to files.
  - bin/crow-state path <kind> [args] CLI for skills.
  - crow-registry exposes state_dir (machine) and crow_dir (notes); crow-runs, crow-hub-guard, crow-mrs, crow-project, crow-review-run, crow-sessions use the machine dir; crow-gather and crow-deck read notes through crow_dir.
  - Tests on a temp vault under the repo.
- Out:
  - Moving any file
  - Daily note writing

# Acceptance criteria
- With crow_dir set to a temp vault, runs.jsonl, hub-guard, reviews cache, handoffs and session caches still resolve to local/state.
- crow-gather and crow-deck read Open Threads.md and Status.md from the crow dir.
- Project tests pass.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: With crow_dir set to a temp vault, runs.jsonl, hub-guard, reviews cache, handoffs and session caches still resolve to local/state.
- request-AC2 -> This backlog slice. Proof: crow-gather and crow-deck read Open Threads.md and Status.md from the crow dir.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_006_crow_data_in_the_vault`
- Architecture decision(s): (none yet)
- Request: `req_041_crow_data_in_the_vault`
- Primary task(s): `task_015_build_crow_data_in_the_vault`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_015_build_crow_data_in_the_vault`

# Notes
- Task `task_015_build_crow_data_in_the_vault` was finished via `logics-manager flow finish task` on 2026-11-09.
