## req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map - crow deck: interface redesign (tabbed TUI, 4 themes, command map)
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: High
> Theme: Hub observability
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Build crow deck, the tabbed TUI hub (Home, Agents, MRs, Projects, Usage) with 4 themes and a regrouped crow command map, from the approved designs.
- Keywords: crow, deck, tui, themes, mrs, projects, usage, commands
- Use when: working on the crow deck UI, its data layer or the crow CLI naming.
- Skip when: the change is to the cost audit or batched retros (separate audit request).

# Needs
- One place to see agents, MRs, projects and usage, and to act on them, without opening five tools.
- Approved designs exist: 14 page boards (landscape and portrait) and 4 themes (lazytf, HUD, Cyber, Soft). The Notes tab is parked.
- The CLI has grown command by command; group it by what the owner does.

# Context
- Spec: local/state/research/crow-deck-design.md. Boards: local/state/research/crow-deck-boards/. Vault mirror: CR-0W/Research.
- Go + Bubble Tea, ANSI 16 colors, lazytf look as default theme.
- Failover stays off while only one account is `use` in local/accounts.yaml.
- Data in the boards is invented (MR numbers, second-account usage, costs, history rows).
- Prerequisites: Jira API token (owner creates it), M365 later, per-session cost accounting to verify, structured review findings, roadmap file format.

# Acceptance criteria
- AC1: `crow deck` opens the 5 tabs in landscape and portrait with real data, `crow watch` still works as an alias.
- AC2: The 4 themes switch with `t` on every tab.
- AC3: Nothing is pushed, posted or sent without an explicit OK shown in the deck.
- AC4: `crow --help` shows the grouped command map; old names keep working for one release.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed: roadmap format, cost accounting check, Jira token open.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `local/state/research/crow-deck-design.md`
- `runbooks/tui-render.md`

# Backlog
- `item_016_0_1_deck_skeleton_chrome_tabs_theme_engine_config_mock_home`
- `item_017_0_2_real_data_on_home_and_agents`
- `item_018_0_3_mrs_tab_both_forges_read_only`
- `item_019_0_4_projects_and_usage_tabs`
- `item_020_0_5_actions_and_history`
- `item_021_0_6_remaining_pages_in_all_themes`
- `item_022_0_7_crow_command_map_and_doctor`
