## req_047_make_the_hub_rules_configurable_in_local_policy_yaml - Make the hub rules configurable in local/policy.yaml
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Move the hub rules (commit format, autonomy table) from code and prose into one config file read by the hooks.
- Keywords: hub, rules, configurable, local, policy, yaml
- Use when: picking up this hub improvement.
- Skip when: working on another repo.

# Needs
- Move the hub rules (commit format, autonomy table) from code and prose into one config file read by the hooks.

# Context
- The commit format is hardcoded in `bin/crow-commit-msg-check`.
- The autonomy table lives as prose in CLAUDE.md.
- The README says these are the author's example set until they are configurable.

# Acceptance criteria
- AC1: Commit types, scope rule and subject pattern come from `local/policy.yaml`.
- AC2: The autonomy table (per action, work vs personal) is read from the same file by the guard hooks.
- AC3: A missing file falls back to the current defaults.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-commit-msg-check`
- `bin/crow-hub-guard`
- `CLAUDE.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

