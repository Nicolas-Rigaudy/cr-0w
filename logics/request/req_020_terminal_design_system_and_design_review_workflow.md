## req_020_terminal_design_system_and_design_review_workflow - Terminal design system and design review workflow
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: High
> Theme: design
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Replace mixed design channels with a code-rendered review loop and a shared Go terminal kit; crow deck pilots, lazytf and archeotech follow.
- Keywords: terminal, design, system, review, workflow
- Use when: Planning or reviewing any terminal or shell UI design.
- Skip when: Web UIs or non-visual work.

# Needs
- One design review channel for terminal apps and shell work, replacing the mix of Claude Design canvas boards, hand-patched artifacts and bare imv windows.
- A review page that is easy to read, shows only what needs a decision, shows what changed since the last round, and never keeps stale content.
- Designs rendered from real code (ANSI output through the user's kitty theme and font), never imitated in HTML or a pixel canvas.
- A modular terminal design system (tokens, components, gallery) shared by crow deck and lazytf, built and reviewed layer by layer: system, then components, then screens.
- Crow deck becomes user-composable: the user picks which widgets appear, where, and with what settings.
- One look across the desktop and the terminal: archeotech themes drive kitty's ANSI 16, TUIs map roles to those slots.
- Long term: design decisions and other daily work are made from inside crow deck, with little typing in terminals.

# Context
- Owner feedback 2026-11-03: Claude Design boards had too many renders and no component work; artifacts kept stale sections and hid what changed; imv windows had no context.
- bin/ansi-shot renders ANSI to PNG; bin/crow-compare builds a side-by-side decision page with a flip viewer. Both exist and are the base.
- crow deck (tools/crow-deck) is ~16k lines of Go in one flat package; frame.go and the tab pages mix rendering with data code.
- lazytf is Go + Bubble Tea + Lipgloss with a separate internal/ui package (titlebar, sidebar, statusbar, modal) and a YAML layout config with per-project override.
- archeotech-shell is QML (Quickshell). shell-config.json composes regions as lists of {id, config, align}; themes generate kitty.conf. It shares tokens and config ideas, not Go components.
- Decided 2026-11-03: shared kit lives in a new public repo on the personal GitHub; crow deck is the pilot; picks go through a copy-picks line or AskUserQuestion until in-deck decisions exist; the deck config shape is decided after a comparison (archeotech shape, lazytf style, or a new shared model).
- Visual rules: terminal ANSI 16 only, transparent background, lazytf look, themes told apart in monochrome by shape, key info visible at a glance, vertical-screen aspect first. No ASCII mockups in chat or prompts.

# Acceptance criteria
- AC1: One documented design review workflow (runbook) covers terminal and shell design; crow-compare pages show open decisions, changes since last round and the decided ledger, regenerated from a spec each round.
- AC2: Every design decision has an ID and lives in a decisions ledger read at the start of each round; the owner never repeats a decision.
- AC3: A public Go module provides tokens (roles to ANSI slots, glyphs, borders, monochrome mode), components as pure render functions, and a gallery that renders every component in every state.
- AC4: crow deck renders through the kit and is composed from widgets declared in a user config.
- AC5: lazytf uses the kit for its shared chrome without visual regression.
- AC6: archeotech themes map to the kit's token roles, documented end to end (archeotech theme to kitty to TUI).
- AC7: A path exists to review and pick design decisions from inside crow deck.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)

# References
- bin/crow-compare
- bin/ansi-shot
- runbooks/design-canvas.md
- tools/crow-deck/
- logics/request/req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map.md
- logics/request/req_011_crow_deck_rework_terminal_true_mockups.md
- local/state/research/crow-deck-design.md
- local/owner-prefs.md
- lazytf: internal/ui/
- lazytf: internal/config/layout.go
- archeotech-shell: shell-config.json
- archeotech-shell: themes/

# Backlog
- `item_078_design_review_workflow_crow_compare_v2_decisions_ledger_runbook`
- `item_079_shared_terminal_kit_repo_tokens_and_gallery`
- `item_080_kit_components_extracted_from_crow_deck_and_lazytf`
- `item_081_customization_model_study_deck_config_vs_archeotech_shell_config`
- `item_082_crow_deck_widget_contract_user_deck_config_tabs_rebuilt_from_widgets`
- `item_083_lazytf_adopts_the_shared_kit`
- `item_084_archeotech_themes_mapped_to_kit_token_roles`
- `item_085_in_deck_design_decisions_widget`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

