## req_035_forge_ops_in_nest - Forge ops in nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Feature track. Operations of crow and the owner's machines in nest: machines and tools, nightly health with repair, releases, cost and quota, jobs. Mostly no-LLM scripts (Jarvis ran about 30) feeding sections of the core API.
- Keywords: machines, tools, versions, nightly health, repair catalogue, watchdog, releases, cost, quota, jobs, backups, credentials expiry
- Use when: building forge pages or ops jobs.
- Skip when: changing host roles (req_029).

# Needs
1. **Machines** (Jarvis). Per host: up, CPU, memory, disk, load, last seen; tool versions vs the reference host (claude, cdx, logics-manager, glab, gh, node, go); credential names and expiry, never values.
2. **Nightly health** (Jarvis three passes). Pass 1 checks (tests on main of personal tools, disk, backups, credentials). Pass 2 repairs known causes from the repair catalogue. Pass 3 queues a `ci-repair` mission for what remains and reports in the morning brief.
3. **Repair catalogue.** Known failure cause → fix script, max 2 rounds, then ask. Fed by the incidents pass.
4. **Watchdog.** Silent run 30 min, busy over 90 min (Ultron 90/20), frozen run resumed once, core restart.
5. **Releases** (Jarvis). Release trains for the owner's tools (crow, cdx config, deck): what is merged since the last tag, notes draft, publish on OK.
6. **Cost and quota** (Ultron). Spend per kind and per validated mission, cache share, compactions, quota gauges and ETA per account.
7. **Jobs.** Scheduled jobs with last run, next run, run now, pause, resume.
8. **Backups.** Daily backup of store, config, state and agent homes; age alert past 36 h.

# Context
- Owner choices 2026-11-06: machines and tools, nightly health, releases, cost and quota; scripts first; repair catalogue, failover, watchdog, incident log.
- Existing: `crow doctor`, `crow runs`, `bin/crow-ctx`, cdx usage data.

# Acceptance criteria
- AC1: Machines page shows the laptop with tool drift and credential expiry.
- AC2: A nightly run repairs one seeded failure from the catalogue and queues a mission for an unknown one.
- AC3: Cost page totals match cdx usage for a day.
- AC4: Backup age alert fires in a test.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/request/req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path.md`
- `logics/request/req_016_crow_everywhere_multi_machine_and_mobile_control.md`
- `logics/request/req_019_release_leverage_cdx_0_25_logics_2_24.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
