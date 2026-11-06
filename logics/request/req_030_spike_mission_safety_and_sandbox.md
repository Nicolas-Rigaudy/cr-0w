## req_030_spike_mission_safety_and_sandbox - Spike mission safety and sandbox
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60
> Confidence: 50
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Pick the safety model for unattended missions by comparing the two colleague systems and current best practice: guard hook with validation markers, OS sandbox with egress rules, and an engine that holds all keys.
- Keywords: safety, sandbox, bubblewrap, guard hook, egress, secrets, engine publishes, auto-merge, policy
- Use when: deciding what a mission agent can touch and who pushes.
- Skip when: tuning one allowlist for one run kind (item_077).

# Needs
1. Compare: Ultron (PreToolUse guard hook, validated markers, Stop guard by phase), Jarvis (bubblewrap sandbox, network namespace blocking the LAN, agent holds no git creds, engine pushes and merges), Claude Code native sandbox and permission modes, item_077 deny-by-default runs.
2. Recommend one layered design: what stops a mistake, what stops a bad prompt, what stops data leaving.
3. Define the auto-merge gate for personal repos: green CI, secret scan over the commit range, a diff size limit, a deploy hook.
4. Keep work hard rules outside the agent: work push, MR, merge, mail and Teams only through approval items.

# Context
- Owner decision 2026-11-06: deeper subject, research what is doable and recommended, steal from both.
- Personal repos auto-merge on green CI (prod_005); work repos never merge.
- Builds on req_019 / item_077 (sandboxed runs) and adr_001 (engine holds the keys).

# Acceptance criteria
- AC1: A comparison table of the options against threats (destructive command, secret leak, LAN access, unwanted push).
- AC2: A recommended design with what is built in core, in cdx, and in run config.
- AC3: A prototype run proves an agent cannot push or reach the LAN while the engine can publish its branch.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
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
