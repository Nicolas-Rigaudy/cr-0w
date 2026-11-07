## item_110_stop_guard_for_runs_that_end_without_a_valid_version_file - Stop guard for runs that end without a valid version file
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A Stop hook refuses to end a mission run that has not written a valid version file, and tells it what is missing.
- Keywords: stop hook, stop guard, contract, version file
- Use when: Every mission run.
- Skip when: Policy on tool calls (item_102).

# Problem
- A run can end without writing the file core expects; the mission then stalls.
- Claude Code Stop hooks can block a stop with exit 2 and see the last message (hooks docs, checked 2026-11-06).

# Scope
- In:
  - A Stop hook that runs the item_106 parser on the expected `phase.vN.md`.
  - Block with the parser errors as the reason; allow the stop once the file is valid.
  - Allow the stop when `stop_hook_active` is true and the file is still invalid, so the run cannot loop; core logs a contract failure. Claude Code also caps consecutive stop-hook continuations at 8.
- Out:
  - Runs where hooks do not load (`--bare`).

# Acceptance criteria
- AC1: A run that ends without the file is blocked once with the reason, then writes it.
- AC2: On a second stop with `stop_hook_active` true and no valid file, the stop is allowed and core logs a contract failure.

# AC Traceability
- request-AC4 -> AC1. Proof: a scripted run in the dry run.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Saves stalled missions; not needed for the first manual dry run.

# Notes
- Depends on item_106, and on how isolated runs load hooks (item_100, item_077).
- Stop input carries `stop_hook_active` and `last_assistant_message`; continuations are capped at 8 (https://code.claude.com/docs/en/hooks, Claude Code 2.1.319, checked 2026-11-07).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
