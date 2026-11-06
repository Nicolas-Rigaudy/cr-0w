## req_029_spike_hosts_model_sources_and_local_ai_benchmarks - Spike hosts model sources and local AI benchmarks
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 18:12:01

# AI Context
- Summary: Research spike. Decide where core and missions run, what each of the owner's three machines can serve, which model sources to use for what (Claude accounts, Antigravity on Gmail accounts, local Ollama), and how quota is probed. Includes local AI benchmarks on laptop, tower and homelab.
- Keywords: hosts, control host, worker, tower idle, homelab, Ollama, benchmarks, Antigravity, Gemini, Google terms, quota probe, routing, data class
- Use when: deciding where a mission or a model call runs, or which account pays.
- Skip when: per-scope routing rules already settled in req_025.

# Needs
1. **Host roles.** Compare the two references:
   - Ultron: everything on one Mac. Simple; dies with the laptop lid.
   - Jarvis: main host (engine, HQ, timers), worker host (missions), office host; Tailscale between them; engine dispatches over SSH.
   - crow: laptop first (owner decision). Measure what a homelab control host would gain (always on, phone answers while the laptop sleeps) against the cost (moving the store, credentials on another box).
2. **Tower as idle worker.** The tower only takes work when idle and never while the owner plays. Define idle (no game process, no input for N min, GPU free), how a running job yields (pause, kill, migrate), and wake-on-LAN or not.
3. **Local AI benchmarks** on all 3 machines. Models: small and mid Gemma, Qwen, Llama, a code model. Tasks: run summary from a `cdx conversation` excerpt, mail triage label, bank line categorisation, memo to action, French short prose. Metrics: tokens per second, first token latency, memory, quality against a Claude baseline on a fixed set of 20 samples, power draw where easy.
4. **Need check.** Ultron dropped local models for lack of measured need. List what crow would actually send local (private data, bulk summaries) and the volume per day, so benchmarks answer a real question.
5. **Antigravity and Google free usage.** Check cdx `--provider antigravity`; what quota each Gmail account gives; Google terms on several accounts for one person; what data Google keeps; which mission kinds fit (personal code, research, cheap text). Rule stays: no work data.
6. **Quota probe.** Ultron probes with a one-token Haiku request and reads rate-limit headers for the 5 h and weekly windows. Check what cdx exposes per provider; define room (lower window, ≥ 20 % left, interactive reserve).
7. **Router table.** Kind × scope × data class → provider, account, model, effort, host. Fallback order per row.

# Context
- Folded in 2026-11-06: req_023 (scheduled jobs that survive the lid closed and offline). The control host choice answers it.
- Owner decisions 2026-11-06: laptop first, homelab later if it pays, tower idle only; Claude for anything that matters and all work data; Antigravity on Gmail accounts after a terms check; local after benchmarks.
- cdx already has profiles with policy (`use`, `ask-first`, `never`), priority and budget; req_025 routes by scope.
- Jarvis used Codex first, Claude second, Ollama for private and bulk; Ultron used 3 Claude accounts only.

# Acceptance criteria
- AC1: A host recommendation (control host now and later, tower rules) with the measured gain and cost.
- AC2: A benchmark table for 3 machines × chosen models × 5 tasks, with a go or no-go per task.
- AC3: A written terms check for Antigravity on several Gmail accounts and a decision on how many to use.
- AC4: A probe script returning room per account for Claude and, if possible, Antigravity.
- AC5: A router table ready to paste into core config.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_023_scheduled_jobs_that_survive_lid_closed_and_offline.md`
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`
- `logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md`
- `logics/request/req_025_account_routing_by_scope.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
