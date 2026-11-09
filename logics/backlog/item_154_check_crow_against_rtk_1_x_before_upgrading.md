## item_154_check_crow_against_rtk_1_x_before_upgrading - Check crow against rtk 1.x before upgrading
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub upkeep
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Before any rtk upgrade, review the 1.x licence and test crow's rtk uses against 1.2.0 output: `rtk read`/`rewrite` parsing, `rtk gain` family parsing, `rtk grep` flags in runbooks, and the `proxy` to `bypass` rename in docs.
- Keywords: rtk, 1.2.0, licence, rtk read, rtk rewrite, rtk gain, rtk grep, rtk bypass
- Use when: preparing the rtk 0.35 to 1.x move.
- Skip when: running the upgrade itself (not taken up in triage; a later item once this one is done).

# Problem
- rtk is at 0.35.0; 1.2.0 is out. Several releases break formats crow parses or flags crow documents, and 1.0.0 changes the default binary's licence.

# Scope
- In:
  - Licence: read the 1.x licence terms and write a short verdict in the item (fit for a personal, public-portfolio repo and for work use). Stop the rest if it does not fit.
  - `bin/crow_lib.py`: test the `rtk read` parser against 1.2.0 (tab-separated line numbers, raw bytes on invalid UTF-8) and `rtk rewrite` against 0.51.0 (no implicit shell expansion). Use a 1.2.0 binary in a scratch dir, not the installed one.
  - `bin/crow-usage`: test `rtk gain`, `saved`, `history` and `could` parsing against 1.2.0 output. Add fixtures for both formats if the parser must handle both.
  - `runbooks/dispatch-background-runs.md`: re-check `rtk grep` flags (0.47.0 removed `--file-type`; `-l`, `-m`, `-t` now native).
  - Docs: replace `rtk proxy` with `rtk bypass` in `local/runbooks/lambda-red-pipeline.md`, `local/runbooks/service-tokens.md`, `runbooks/readme-showcase.md` and the RTK.md notes, only once 1.x is adopted (the alias keeps old docs working).
- Out:
  - Installing 1.x on the machine (owner OK, separate item).

# Acceptance criteria
- AC1: The licence verdict is written in this item before any other step.
- AC2: `bin/crow-check` passes with fixtures of 1.2.0 `rtk read` and `rtk gain` output, and still with 0.35 output.
- AC3: No crow runbook uses a removed `rtk grep` flag.
- AC4: The `proxy` to `bypass` doc edits are ready (or done if 1.x is adopted).

# AC Traceability
- AC1 -> licence section in Notes.
- AC2 -> tests in `tests/` for crow_lib and crow-usage parsers.
- AC3 -> grep of `runbooks/` and `local/runbooks/`.
- AC4 -> diff of the four docs.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_045_periodic_tool_update_watch
- Primary task(s): (none yet)
- Source: vault note `CR-0W/Tool updates/Tool updates 2026-11-09.md`

# Priority
- Priority: Low
- Rationale: 0.35 still works; this unblocks a later upgrade.

# Notes
- Groomed from the toolwatch triage of 2026-11-09. The upgrade in steps was not taken up.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
