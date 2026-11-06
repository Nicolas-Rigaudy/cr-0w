## req_035_forge_ops_in_nest - Forge ops in nest
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 65
> Confidence: 60
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Forge ops in nest: machines and tools page, nightly health with the repair catalogue, releases page for the owner's tools, cost and quota page.
- Keywords: machines, health, nightly, backups, repair catalogue, releases, cost, quota, eta
- Use when: building the tech or forge side of nest.
- Skip when: changing a tool's own release pipeline.

# Needs
1. **Machines and tools.** Per host: online or asleep, CPU, memory, disk, tool versions deployed versus main, credential expiry.
2. **Nightly health.** Scripted night jobs: tests on main, tool updates, backups. A failure queues a repair mission that follows the catalogue (at most 2 rounds, then asks).
3. **Releases page.** Release trains for the owner's tools (crow, lazytf, cdx and others): prepare, publish, update each host.
4. **Cost and quota.** Gauges per account, cost per mission kind, an ETA for queued missions when every account is out.

# Context
- Builds on item_075 (cost per run), req_006 (releases), req_016 (multi-machine), `crow doctor`.
- Scripts first; models only narrate.

# Acceptance criteria
- AC1: The machines page shows every tailnet host with tool versions.
- AC2: A red nightly test creates a repair mission and an FYI or an ask.
- AC3: The cost page shows spend per mission kind for the last 7 days.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
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
