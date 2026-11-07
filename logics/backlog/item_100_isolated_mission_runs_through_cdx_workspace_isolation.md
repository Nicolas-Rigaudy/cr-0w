## item_100_isolated_mission_runs_through_cdx_workspace_isolation - Isolated mission runs through cdx workspace isolation
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Safety
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:33:25

# AI Context
- Summary: Mission runs start through cdx --isolate-workspace (bwrap, volatile home, read-only git, env allowlist, pasta + nftables egress) and core commits the worktree diff afterwards.
- Keywords: isolated, mission, runs, through, cdx, workspace, isolation
- Use when: Turning on the OS boundary for background and mission runs.
- Skip when: Run-config allowlists (item_077) or the publish path.

# Problem
- Runs see the owner home and network today; the probe shows a shared network namespace reaches the LAN and host loopback services.
- cdx already ships the boundary but needs pasta, a setup token and no `--continue-run`.

# Scope
- In:
  - Add `--isolate-workspace` to the kinds that allow it in `RUN_KINDS`, read-only for review.
  - Preflight: pasta, nft and a setup token present, else refuse the isolated kind with a clear message.
  - core commits the worktree diff after the run with a checked message.
  - Resume as a fresh isolated run from brief and answers (item_116 adds the cdx issue for native resume).
  - A probe test: no credential files or env, no push, no LAN or tailnet, public inference reachable.
  - Network: public internet minus LAN (owner decision 2026-11-07).
- Out:
  - Installing pasta (owner, system package).
  - Inference-only egress for local runs (cdx change, later).

# Acceptance criteria
- AC1: An isolated run cannot read forge config, ssh or cloud dirs, nor see token env vars.
- AC2: An isolated run cannot push and cannot reach private, CGNAT or loopback ranges.
- AC3: core commits the run's diff and the commit message passes the commit check.
- AC4: A machine without pasta gets a refusal with the fix, never a shared-network fallback.

# AC Traceability
- request-AC3 -> sandboxed agent cannot push, read credentials or reach the LAN.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: The only layer that holds against an injected agent.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- passt installed by the owner 2026-11-07: unblocked. Rerun the safety probe with pasta as the first step of this item.
- Aligns with item_077: that item is the run-config layer, this one the OS layer.
- Decisions 2026-11-07 in req_030: public minus LAN, core commits the diff, fresh run from files for resume.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
