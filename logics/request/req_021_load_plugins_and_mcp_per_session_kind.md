## req_021_load_plugins_and_mcp_per_session_kind - Load heavy plugins and MCP servers per session kind
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub cost
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-05 14:15:09

# AI Context
- Summary: Plugins, skills and MCP servers load per session kind (hub, repo, background run) instead of all of them in every profile session.
- Keywords: plugins, mcp, skills, context, start context, profile, session kind
- Use when: deciding which tools a session loads.
- Skip when: status line or deck work.

# Needs
- Hub sessions start at a median of 48.8k tokens (item_029 baseline, 2026-11-04). Every enabled plugin and MCP server adds its skill and tool descriptions to every session on the profile, used or not.
- The owner wants a better way than all-or-nothing per profile: keep tools that are useful in some sessions (Jira through Atlassian, design skills) without paying for them everywhere.
- Status-line tools like claude-hud cost no model context: the status line is a separate process. Only their slash-command descriptions load.

# Context
- 2026-11-05: ui-ux-pro-max disabled on the main profile; claude-hud and Atlassian kept.
- Leads to check first: `--strict-mcp-config` with a per-kind MCP list (item_077 plans it for background runs), project-level `.claude/settings.json` enabling plugins only in the repos that need them, cdx `--extra-args` per launch, crow handoff/term passing the set for the target repo.

# Acceptance criteria
- AC1: Each session kind (hub, each registered repo, each background run kind) has a declared tool set in one config file.
- AC2: A hub session no longer loads tools only a repo session uses; start context re-measured with bin/crow-cost-audit.
- AC3: A tool used in one repo still loads there with no manual step.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics/backlog/item_029_trim_hub_start_context_and_fix_per_model_usage.md`
- `logics/backlog/item_077_sandboxed_background_runs.md`
- `runbooks/cost-audit.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
