## item_077_sandboxed_background_runs - Background runs start with no authorisation and whitelist tools
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Move crow background runs from auto mode (allowed unless blocked) to deny by default with a per-kind tool allowlist, scrubbed credentials and only declared MCP servers.
- Keywords: sandbox, allowedTools, scrub-env, strict-mcp-config, isolate-workspace, permission
- Use when: changing what a background run may do.
- Skip when: interactive sessions.

# Problem
- Task and one-shot runs use cdx `--permission auto`: everything is allowed unless the auto-mode safety check blocks it. Runs see the owner's credentials (forge and cloud tokens) and load the repo's own MCP servers. Reviews use plan mode, which blocks writes but not reads of secrets or network calls.
- The cdx maintainer advises the reverse: start with no authorisation, then whitelist what a run needs.

# Scope
- In:
  - Per-kind allowlists in crow_lib next to RUN_KINDS (review: read, search, git read commands, forge MR read; task: edits inside its worktree, tests, git commit; never push; oneshot and retro: minimal).
  - `--scrub-env` behaviour for runs (credentials stripped; check which forge tokens survive).
  - `--strict-mcp-config` plus an explicit MCP list per kind (none by default).
  - Workspace isolation: fold in item_076 (`--isolate-workspace`) or keep worktrees, decided in the cdx design check.
  - Denials: cdx `permission_denials` warnings shown in `crow runs` and the deck so a missing tool is visible and added on purpose.
  - Transport: a per-run parameter in cdx (owner fork handoff) if the design check keeps it small; otherwise a background-only cdx profile per dispatch account with `cdx set --scrub-env on --extra-args ...`.
- Out:
  - Interactive sessions and their permission modes.
  - Network sandboxing beyond what Claude's sandbox and scrub-env give.

# Acceptance criteria
- AC1: Every crow background run starts in the default permission mode with only its kind's allowlist; a tool outside it is denied and reported.
- AC2: A run cannot read the owner's cloud or forge tokens from its environment (probe run proves it), except tokens the kind needs and the owner approved.
- AC3: Repo MCP servers do not load unless the kind lists them.
- AC4: The review, task and retro probes still pass with the allowlists (no needless denials).
- AC5: The runbook lists each kind's allowlist and how to widen one.

# AC Traceability
- request-AC4 -> This backlog slice (absorbs or supersedes item_076). Proof: AC1 to AC5 above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: background runs act unattended with full credentials today.

# Notes
- Depends on the cdx-manager design check (item_015, 2026-11-03 follow-up). Prior proof: scrub-env on plus `--allowedTools` passed the probe on 2026-10-29 once socat was installed; MCP tools untested.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

