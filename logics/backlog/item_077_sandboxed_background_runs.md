## item_077_sandboxed_background_runs - Background runs start with no authorisation and whitelist tools
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 60%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:14:57

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
- Blocked (owner decision 2026-11-07): the allowlists wait for cdx issue 120 (per-run `--extra-args`). AC1, AC3 and AC4 stay open until it ships.
- Depends on the cdx-manager design check (item_015, 2026-11-03 follow-up). Prior proof: scrub-env on plus `--allowedTools` passed the probe on 2026-10-29 once socat was installed; MCP tools untested.
- Claude Code ignores the sandbox domain allowlist in bypass mode. A sandboxed run must use acceptEdits or dontAsk, never bypass (harness-maxxing review 2026-11-07, local/state/research/harness-maxxing-and-hermes-2026-11-07.md).
- Secret paths go in `permissions.deny` as well as the sandbox: Read/Edit/Write are not shell commands, and deny rules hold in every mode. Fixed per-kind rules fit a PreToolUse hook run with `python -I`.

# Progress 2026-11-07
- Done: per-kind allowlists, shared denials, credential filter and denial parsing in `crow_lib` (`RUN_SANDBOX`, `RUN_DENY`, `CREDENTIAL_ENV`, `RUN_KEEP_ENV`, `sandbox_args`, `run_denials`); every dispatch path (crow-runs, crow-retro, crow oneshot) starts cdx without credential variables; `crow runs` and the deck show denials; runbook `dispatch-background-runs.md` "Run sandbox".
- AC2 holds. Probe (scratch folder, `cdx run rk1 --permission auto`, haiku): `GH_TOKEN`, `AWS_SECRET_ACCESS_KEY` and `FAKE_PROBE_TOKEN` set in the dispatching env were absent from the run's `env`. gh and glab keep their config dirs. Test: `tests/test_crow_run_sandbox.py::test_dispatch_hands_cdx_a_credential_free_env`.
- AC5 holds: the runbook lists each kind's allowlist and the widening steps.
- AC1 partial: denials are reported (probe with `--permission default`: both Bash calls came back as `permission_denials`, parsed by `run_denials`, shown in `crow runs` and the deck). Passing the kind's allowlist per run needs `cdx run --extra-args` (upstream cdx-manager issue 120). `cdx set --extra-args` is per account and also changes interactive launches.
- AC3 open: the probe showed a scratch folder's `.mcp.json` server loading in a headless run. `--strict-mcp-config` needs the same per-run transport.
- AC4 open: the allowlist probes run once the transport exists.
- Workspace isolation (item_076): kept as worktrees here; the OS boundary is item_100 (cdx `--isolate-workspace`, needs pasta).
- Next once cdx ships per-run args: dispatch passes `--permission default --extra-args` with `sandbox_args(kind)`, then review, task and retro probes.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
