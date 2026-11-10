## road_001_crow_roadmap_after_the_deck_pass - Crow roadmap after the deck pass
> Date: 2026-11-05
> Status: Proposed
> Related product: `prod_005_crow_jarvis_layer`
> Related request: `req_000_crow_watch_and_bar_modules_ui_ux_pass_session_titles_savings_stats`, `req_001_crow_account_switching_and_session_resume`, `req_004_portfolio_readmes_profile_readme_and_personal_page`, `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`, `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`, `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`, `req_011_crow_deck_rework_terminal_true_mockups`, `req_014_crow_management_checks_and_actionable_deck`, `req_016_crow_everywhere_multi_machine_and_mobile_control`, `req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch`, `req_019_release_leverage_cdx_0_25_logics_2_24`, `req_020_terminal_design_system_and_design_review_workflow`, `req_021_load_plugins_and_mcp_per_session_kind`, `req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control`, `req_025_account_routing_by_scope`, `req_026_personal_butler`, `req_027_conversations_feed_for_core_deck_nest`, `req_039_news_board_for_tech_news_and_repos`
> Reminder: Update status, milestone scope, linked refs, risks, and success signals when you edit this doc.
> Indicators reviewed: 2026-11-10 16:17:24

# AI Context
- Summary: Order for the remaining crow work: tidy, finish in-flight, safe runs, Jarvis core, Jarvis tracks, releases, design system, management deck, side projects. The cockpit milestone folded into the Jarvis tracks on 2026-11-06.
- Keywords: roadmap, milestones, versions, crow roadmap after the deck pass
- Use when: Planning or sequencing versions for Crow roadmap after the deck pass.
- Skip when: You need execution details for a single backlog item or task.

# Summary
Agreed with the owner on 2026-11-05. Milestone numbers are roadmap steps, not crow version numbers. Each item is built in a fresh session from its backlog doc.

# Milestones
## 0.1 - Tidy up
- Goal: A clean workflow state before new work. One short session.
- Scope: close item_004 and item_015 (at 100%); `flow start` the draft requests already in delivery; close or withdraw req_002, req_003 (superseded by this roadmap) and req_012 if covered; fix `crow handoff` ignoring CROW_ACCOUNT equal to the session account (bin/crow:47) and the unread `planning:` key; fix req_018; write item_087 (work prose checklist).
- Exit signal: `logics-manager status` lists no docs at 100% still open and no drafts in delivery.

## 0.2 - Finish in-flight work
- Goal: Everything already started is merged and closed.
- Scope: req_000 (001 done 2026-11-05; 002 and 005 closed, crow watch folded into the deck) and req_001 (done 2026-11-05: 009 built in the deck, account table and no-model group); req_008 via task_003 (apply 023, re-measure 029, build 025 gather scripts, then 024 batched retros); req_007 deck (017-022, 038).
- Exit signal: req_000, req_001, req_007 and req_008 are Done; brief, daily and weekly run from gather scripts.
- Rule: no new deck features until this milestone closes.

## 0.2b - Account routing by scope (exception)
- Goal: Work, personal and hub sessions each prefer their own account and overflow by rule, on top of the cdx ranking (usability, then `--priority`).
- Scope: req_025 via item_095. One session.
- Exit signal: `crow-accounts pick --scope` routes by rank and matches `cdx select` without a scope.
- Exception (owner, 2026-11-06): ahead of 0.3 because background runs and the personal butler need the overflow rules.

## 0.3 - Safe background runs
- Goal: Unattended runs start with no authorisation, load only what their kind needs, and record cost.
- Scope: req_019: 077 sandboxed runs (fold 076 isolated review workspace into it), then 074 dispatch records runs, 075 cost per run, 073 deck Logics viewer. Groom req_021 and run the 088 trials alongside.
- Exit signal: a headless run loads no repo `.mcp.json` server it does not need, and `crow runs` shows its cost.
- Status: req_019 Done 2026-11-09 (076 last: peer reviews run isolated, self and threads stay unisolated for apply).

## 0.3b - Jarvis layer core (exception)
- Goal: work runs as background missions; questions and proposals land in one inbox; status reads as short human sentences.
- Scope: req_028 slices 1-5 (mission model, inbox with park and resume, core API and MCP server, human summaries, autonomy dial, reconciler), kinds `research`, `plan`, `review`, `design` (req_037), nest Home, Missions and Mission detail on both Fold 6 screens, deck inbox pane. Slices 6-9 (deck views, nest MVP, butler chat, model router) follow in J.
- Exception (owner, 2026-11-06): ahead of 0.4 because the owner wants to steer from the interface and run more hands off. Needs 0.3 (deny-by-default runs) first.
- Exit signal: a mission asks a question, exits, and resumes from an answer given in deck or nest.
- Status: slices 1-2 groomed 2026-11-09 into task_016 (10 items, Go core in `tools/crow-core`); slices 3-5 groomed 2026-11-10 into task_017 and task_020.
- Exception (owner, 2026-11-10): the Nest answer sheet (item_168) is pulled from slice 7 so the exit signal can be met from Nest. task_020 takes Nest live first (core API, live adapter, actions, answer sheet) and interleaves the Home design items.
- Nest early start (owner, 2026-11-09): design item_125 → item_126, then item_137 nest app shell on the item_119 fixtures, in parallel with task_016.
- Before it: spikes req_029 (hosts, models), req_030 (safety), req_031 (stack), req_032 (butler engine) run alongside 0.3. Brief: prod_005; architecture: adr_001. Design docs req_037 (mission kinds and run contract, first) and req_038 (nest design rules).

## J - Jarvis tracks (after 0.3b)
- Goal: nest becomes the daily interface, then the feature tracks land.
- Scope: J1 daily interface: req_028 slices 6-10 with req_031 and req_032 (Web Push, butler, AFK batches), kinds `review-fix` and `write`, Work and Personal screens, req_026 butler items as core jobs, req_039 news board. J2 tracks: kinds `build` and `repair` after the sandbox, System and Logics screens (req_034, req_035), persona, voice, Discord bot. J3: req_033 Android extras, req_036 fun layer.
- Folded in (2026-11-06): 0.7 work cockpit (req_022 core and nest parts). req_023 (scheduled jobs with the lid closed) moves into the hosts spike req_029. The Crow HQ design research (2026-10-29) is superseded by nest plus the deck inbox pane.
- Exit signal: a full work day runs from nest with no terminal opened except deep repo sessions.

## 0.4 - Release path
- Goal: crow is installable and updatable by someone else.
- Scope: req_006 and req_004: 037 docs (README install, CONTRIBUTING, SECURITY, release runbook), 034 install script and `crow update`, 035 release workflow (tag, CI, GitHub Release, checksums), 036 AUR then Homebrew tap.
- Exit signal: a fresh machine installs a tagged release and `crow update` moves it to the next one.
- Before 0.5 because the kit repo is also public.

## 0.5 - Design system core
- Goal: The deck is rebuilt from shared widgets with a design review loop.
- Scope: req_020 via task_010: 078 crow-compare v2 and decisions ledger, 079 shared kit repo, 080 components extracted from deck and lazytf, 082 widget contract with tabs rebuilt from widgets, 085 in-deck design picks; req_011 item 086 deck follow-ups; then item_090 redesigns the archeotech crow modules (agents, usage panels) on the same widgets and tokens.
- Exit signal: every deck tab renders from widgets declared in the user deck config.
- Before 0.6 so the management features are not built twice.

## 0.6 - Management deck
- Goal: The deck drives management work, every row has an action.
- Scope: req_014: 064 action model, 063 checks engine feeding You decide, 066 1:1 prep and capture, 067 sprint rituals, 065 calendar actions.
- Exit signal: a 1:1 and a sprint planning are prepared from the deck alone.
- Exception: 063 is mostly a script and may move ahead of 0.5 if management needs it sooner.

## 0.7 - Work cockpit (folded into J, 2026-11-06)
- core, nest and req_027 moved into 0.3b and J. The terminal cockpit layout is dropped; nest is the main surface.
- Kept: a spike of a Claude Code panel mod, run when convenient.

## B - Personal butler (parallel track)
- Goal: crow runs the owner's personal admin on the personal account: personal brief, mail and calendar, appointment reminders, budget.
- Waits for core (owner, 2026-11-06): its items are built as core jobs in J1.
- Scope: req_026, built on the personal account. Slices: personal connectors and personal brief, appointments and reminders, budget from bank exports. Later: conversation-backed memory (req_027) once core exists.
- Exit signal: a personal brief runs daily from the vault and the monthly budget review needs no manual sorting.
- Parallel: it runs on the personal account and touches no work milestone.

## 0.8 - Side projects and later
- Goal: Lower priority work once the hub is settled.
- Scope: req_020 items 081 customization study, 083 lazytf adopts the kit, 084 archeotech themes mapped to token roles (overlaps the archeotech-dotfiles cyberpunk pack); req_016 multi-machine; req_015 wishlist. item_070 stays parked.
- Exit signal: each one is either scoped through /corpus or cut.

# Sequencing
- Deliver milestones in ascending version order unless dependencies force a documented exception.
- Keep each increment independently reviewable and linked to concrete workflow docs.

# Risks
- Long-term scope can drift unless every milestone keeps a clear exit signal.
- Version labels are planning targets, not release promises.

# References
- Product brief(s): `prod_005_crow_jarvis_layer`
- Request(s): `req_000_crow_watch_and_bar_modules_ui_ux_pass_session_titles_savings_stats`, `req_001_crow_account_switching_and_session_resume`, `req_004_portfolio_readmes_profile_readme_and_personal_page`, `req_006_clean_crow_releases_versioning_test_pipeline_install_and_update_path`, `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`, `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`, `req_011_crow_deck_rework_terminal_true_mockups`, `req_014_crow_management_checks_and_actionable_deck`, `req_016_crow_everywhere_multi_machine_and_mobile_control`, `req_018_deck_review_confirm_shows_the_stale_mode_when_an_own_mr_got_threads_after_the_list_fetch`, `req_019_release_leverage_cdx_0_25_logics_2_24`, `req_020_terminal_design_system_and_design_review_workflow`, `req_021_load_plugins_and_mcp_per_session_kind`, `req_022_crow_deck_as_the_work_hub_terminal_cockpit_and_phone_control`, `req_025_account_routing_by_scope`, `req_026_personal_butler`, `req_027_conversations_feed_for_core_deck_nest`
- Backlog item(s): (none yet)
- Task(s): `task_003_deliver_the_crow_cost_and_reliability_pass`, `task_010_orchestrate_terminal_design_system`
