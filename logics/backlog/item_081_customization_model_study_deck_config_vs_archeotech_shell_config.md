## item_081_customization_model_study_deck_config_vs_archeotech_shell_config - Customization model study: deck config vs archeotech shell-config
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: design
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Compare archeotech shell-config, lazytf layout and prior art; owner picks one composition model.
- Keywords: customization, model, study, deck, config, archeotech, shell
- Use when: Deciding how users compose decks or shells.
- Skip when: Implementing the chosen config.

# Problem
- The owner wants crow deck and archeotech to be composable and asked whether they should share one model, keep archeotech's, or both move to something new.

# Scope
- In:
  - Read archeotech shell-config.json and its widget loader, lazytf's layout config, and crow deck's tab model.
  - Compare prior art: zellij layouts, waybar modules, tmux, lazygit config, Grafana dashboard JSON, Home Assistant dashboards.
  - Widget contract draft: id, config schema, min size, data source, refresh, actions.
  - Options with example configs for both apps, presented as one decision on the review page; record the pick.
- Out:
  - Changing archeotech code (a follow-up in its own session if the pick requires it).

# Acceptance criteria
- AC1: Comparison doc in local/state/research/ with 2-3 options and a recommendation.
- AC2: Owner pick recorded in the ledger, with the impact on archeotech stated.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC1: Comparison doc in local/state/research/ with 2-3 options and a recommendation.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)
- Request: `req_020_terminal_design_system_and_design_review_workflow`
- Primary task(s): `task_010_orchestrate_terminal_design_system`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

