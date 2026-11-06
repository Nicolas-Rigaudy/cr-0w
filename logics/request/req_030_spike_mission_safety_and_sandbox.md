## req_030_spike_mission_safety_and_sandbox - Spike mission safety and sandbox
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 55
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Pick the layered safety model for unattended missions from the two colleague designs and current Claude Code features: guard hook and validation records (Ultron), OS sandbox with egress rules and an engine that holds every key (Jarvis), native sandbox and permission modes, item_077 deny-by-default runs. Define the personal auto-merge gate.
- Keywords: safety, sandbox, bubblewrap, pasta, nftables, guard hook, Stop hook, allowlist, egress, secrets, engine publishes, auto-merge, secret scan
- Use when: deciding what a mission agent can touch and who publishes.
- Skip when: tuning one allowlist for one run kind (item_077).

# Needs
1. **Options to compare.**
   - Ultron: about 170 allow rules in run settings; a PreToolUse guard hook that blocks push, MR creation and destructive git unless a `validated/KEY` marker exists; a Stop hook that checks the phase output exists before letting a run end.
   - Jarvis: bubblewrap with a temp home, project dir bound read-write, no git or forge credentials inside; network through pasta with nftables blocking the LAN and the tailnet; the engine pushes, opens PRs and merges.
   - Claude Code native sandbox and permission modes; cdx `--permission` per account.
   - crow item_077: deny by default, whitelist per run kind, scrub-env.
2. **Threats.** Destructive command on the host, secret read and leak, LAN or tailnet access, unwanted push or post, prompt injection from a ticket, MR or mail, runaway cost.
3. **Layered design.** For each threat: what stops it in run config, in cdx, in core, and at the OS. Who holds which credential.
4. **Auto-merge gate (personal).** Green CI, secret scan over the commit range, diff size limit with fallback to asking, sensitive paths always ask (CI, sandbox, units, policy, merge gate), full-suite receipt before main, deploy hook after merge.
5. **Work rules.** Push, MR, merge, Jira and Confluence writes, mail and Teams only through core after a validation record for the exact text.
6. **Escape hatch.** Live terminal "take control" and how it is logged as a manual gesture.

# Context
- Owner decision 2026-11-06: research what is doable and recommended, steal from both.
- Memory: scrub-env broke headless runs until `--allowedTools` and socat; MCP inside sandboxes untested.
- Builds on req_019 / item_077 and adr_001 (engine holds the keys).

# Acceptance criteria
- AC1: A threat × option table with a verdict per cell.
- AC2: A recommended design split into run config, cdx, core and OS, with credential ownership.
- AC3: A prototype run proves an agent cannot push, read the owner's credentials or reach the LAN, while core publishes its branch.
- AC4: The auto-merge gate written as a checklist core can run, tested on one personal repo.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_019_release_leverage_cdx_0_25_logics_2_24.md`
- `logics/architecture/adr_001_core_missions_and_inbox_architecture.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
