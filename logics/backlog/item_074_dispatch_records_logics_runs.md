## item_074_dispatch_records_logics_runs - crow dispatch records runs in the linked Logics doc
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 18:12:20

# AI Context
- Summary: crow dispatch records each run in the linked Logics doc's Executions table via `flow record-run`.
- Keywords: dispatch, record-run, executions, logics
- Use when: linking crow runs to Logics docs.
- Skip when: other req_019 slices.

# Problem
- A Logics doc does not show the background runs that worked on it; the link lives only in the crow ledger.

# Scope
- In:
  - `crow dispatch --logics <ref>` (or detection of an item ref in the task text, decide while building) stores the ref in the ledger row.
  - On dispatch: `logics-manager flow record-run <ref> --runner cdx --run-id <id> --state running` in the target repo.
  - On the final state seen by the notify loop: record-run again with done, failed, budget, stuck or limit.
  - Never `--progress auto` by default; the owner sets progress.
- Out:
  - Writing to repos without Logics.
  - Backfilling past runs.

# Acceptance criteria
- AC1: A dispatch with a ref adds a running row to the doc's Executions table.
- AC2: The final state updates the same row, not a new one.
- AC3: A missing logics-manager or a bad ref logs a warning and never fails the dispatch.
- AC4: Tests cover the commands built for start and final states.

# AC Traceability
- request-AC2 -> This backlog slice. Proof: AC1 to AC4 above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner picked it from the 2026-11-03 release review.

# Notes
- 2026-11-06 (Jarvis review): write each dispatch record as a core-shaped event (`{ts, mission, type, actor, data}`, adr_001) so the core store imports them on day one.
- record-run writes to the target repo's working tree: in a work repo that is an uncommitted change in someone's checkout. Decide while building whether work repos are skipped or only hub docs get rows.

# Delivery
- Done 2026-11-06. `crow dispatch --logics <ref>` stores `logics_ref` in the ledger; dispatch and the notify loop call `flow record-run` (running, then succeeded/failed/cancelled with the crow state in Notes) via `record_run` in `bin/crow-runs`. The same run id updates the same row (checked with `--dry-run`).
- Work repos are skipped with a warning (the doc write would dirty someone's checkout); only personal and hub docs get rows. A missing logics-manager or a bad ref only warns.
- Each record is also appended to `runs-events.jsonl` as `{ts, mission, type, actor, data}`.
- Proof: tests/test_crow_runs_logics_cost.py (start and final commands, warnings, work-repo skip, event shape); pytest 303 passed, crow-check green.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
