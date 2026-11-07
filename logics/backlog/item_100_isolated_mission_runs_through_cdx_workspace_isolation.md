## item_100_isolated_mission_runs_through_cdx_workspace_isolation - Isolated mission runs through cdx workspace isolation
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Safety
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:40:09
> Owner: crow

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

# Proof (bin/crow-isolation-probe, 2026-11-07, pasta + bwrap, all checks held)
- AC1: no readable file under the forge, ssh, cloud, cdx and claude dirs; no token env var (fake tokens planted in the launcher env).
- AC2: no git remote visible; `git commit` and push to origin fail; an HTTPS push fails on missing credentials. The LAN gateway, host loopback and tailnet listeners are refused inside the sandbox and open from the host. `api.anthropic.com:443` is reachable.
- AC3: the host commits the sandbox diff with `crow_lib.commit_run_diff` as `new[TXT]: probe run (run probe000)`, which passes `crow-commit-msg-check`.
- AC4: with a PATH without pasta, cdx refuses the run and `crow_lib.isolation_preflight` returns the pacman fix. No shared-network fallback.
- Unit tests: `python3 -m pytest tests/test_crow_isolation.py` (11 passed).

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): `adr_001`
- Request: `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- Commits: `1d7bdd1` probe, `7be7c82` dispatch and commit, `4cefc45` runbook
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: The only layer that holds against an injected agent.

# Notes
- Proposed by the req_030 spike (2026-11-06); Draft until the owner grooms it.
- passt installed by the owner 2026-11-07: unblocked. Rerun the safety probe with pasta as the first step of this item.
- Aligns with item_077: that item is the run-config layer, this one the OS layer.
- Decisions 2026-11-07 in req_030: public minus LAN, core commits the diff, fresh run from files for resume.
- Built 2026-11-07: `crow dispatch --isolate` for task (read-write) and review (read-only). The commit lives in bin/ (`crow-runs notify` calls `crow_lib.commit_run_diff`) until the core (item_119) owns mission commits.
- A live isolated Claude run needs a setup token. The full-login profiles still store one, but cdx skips it when the profile holds a full login. Owner chose a cdx fix (2026-11-07): https://github.com/AlexAgo83/cdx-manager/issues/124. Until it ships, a `<account>-runs` setup-token profile is the workaround.
- Pasta chosen over Claude Code `settings.sandbox`: Claude enforces that one itself and a bypass run skips its domain list.
- Prior art for the OS layer: Claude Code `settings.sandbox` (bubblewrap + socat), writes limited to the run folder, home and secrets hidden, per-run domain list. Compare with pasta before building; same bypass caveat as item_077.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
