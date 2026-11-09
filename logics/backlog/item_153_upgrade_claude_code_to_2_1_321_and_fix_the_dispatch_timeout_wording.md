## item_153_upgrade_claude_code_to_2_1_321_and_fix_the_dispatch_timeout_wording - Upgrade Claude Code to 2.1.321 and fix the dispatch timeout wording
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Small
> Theme: Hub upkeep
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Move Claude Code from 2.1.319 to 2.1.321 for the sandbox and deny-rule fixes, then make the dispatch runbook state the 10 min background-command timeout that `claude -p` now reports.
- Keywords: claude update, 2.1.321, sandbox, deny rules, background timeout, crow dispatch
- Use when: acting on the claude proposals of the tool-updates note of 2026-11-09.
- Skip when: changing the crow-usage price table (declined).

# Problem
- 2.1.319 lets a read-only Bash command run outside the sandbox without a prompt, and its deny rules miss `cp`/`mv` into a directory. Deny-by-default background runs rely on both.
- `claude -p` used to report a 30 min background-command timeout. The real limit is 10 min. The dispatch prompt and runbook may carry the old figure.

# Scope
- In:
  - Upgrade the native install to 2.1.321. This is an install: ask the owner before running it.
  - After the upgrade, run one headless dispatch smoke test and `bin/crow-check`.
  - Grep the dispatch prompt builder and `runbooks/dispatch-background-runs.md` for a background timeout figure; state 10 min where one is given.
- Out:
  - Fable 5.2 rates in `bin/crow-usage` (declined in triage).
  - Package updates beyond Claude Code (owner's dotfiles repo).

# Acceptance criteria
- AC1: `claude --version` prints 2.1.321 or later, after an explicit owner OK.
- AC2: A headless `crow dispatch` smoke run ends ✓ done on the new version.
- AC3: No crow prompt or runbook states a 30 min background-command timeout.

# AC Traceability
- AC1 -> upgrade step. Proof: `claude --version` output.
- AC2 -> smoke run. Proof: `crow runs` row.
- AC3 -> grep of `bin/` and `runbooks/` for the figure.

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
- Priority: Medium
- Rationale: Security fixes for unattended runs; small.

# Notes
- Groomed from the toolwatch triage of 2026-11-09.
- 2026-11-09 done. AC1: `claude --version` prints 2.1.321. `claude update` must run with the real HOME; under cdx it installs into the sandboxed home and leaves the real install untouched. AC2: oneshot smoke run b69f79e7 ended done in 8 s; its `claude --version` call was denied, as deny-by-default intends. AC3: no 30 min background-timeout figure in `bin/`, `runbooks/`, `local/runbooks/` or skills, so no edit was needed.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
