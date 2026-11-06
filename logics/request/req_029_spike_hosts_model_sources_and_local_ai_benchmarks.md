## req_029_spike_hosts_model_sources_and_local_ai_benchmarks - Spike hosts model sources and local AI benchmarks
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60
> Confidence: 50
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Research spike. Decide where core runs (laptop now, homelab later?), which machines can serve local models, and how Google free usage through Antigravity in cdx fits the router.
- Keywords: hosts, homelab, tower, laptop, ollama, benchmarks, antigravity, gemini, quota, router
- Use when: choosing the control host, a local model host, or a non-Claude model source.
- Skip when: building the router rules themselves (req_028 slice 9).

# Needs
1. **Control host.** Laptop first. Study whether a homelab as control is worth it (always-on, power, backups, moving the store) and the best practice for a small always-on agent host.
2. **Tower as idle worker.** The tower is also the gaming PC. Find a way to use it only when idle (game or GPU busy pauses jobs; they wait or fall back to cloud). Wake on LAN, suspend rules.
3. **Local model benchmarks.** Benchmark the laptop, tower and homelab: tokens per second and quality on summaries, triage and classification with a few Gemma, Qwen and Llama sizes. Pick a default model per task and host.
4. **Google free usage through Antigravity.** Use the Antigravity provider in cdx with the owner's Gmail accounts. Check the terms for several accounts first, check headless use, measure the real weekly quota. Compare with a Gemini API key for cheap text tasks.

# Context
- Owner decisions 2026-11-06: laptop first; tower must never take performance while the owner plays; local models for private data and summaries; Antigravity via cdx is the Google path.
- Research 2026-11-06 (see req_028): Gemini key about 250 requests/day on Flash, training on free data; Antigravity free tier weekly and capacity based; several Gmail accounts to stack quota counts as circumvention under Google terms. Verify each point.
- Laptop: Iris Xe, 30 GB RAM, mostly used. Tower GPU and homelab specs to collect.

# Acceptance criteria
- AC1: A short report names the control host for now and the trigger to move it.
- AC2: A benchmark table per machine and model with a default pick per task.
- AC3: An idle-only worker rule for the tower that was tested during a game.
- AC4: A go or no-go on Antigravity per account count, with the terms check cited.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`
- `logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md`
- `logics/request/req_025_account_routing_by_scope.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
