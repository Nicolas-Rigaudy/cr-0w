## item_082_crow_deck_widget_contract_user_deck_config_tabs_rebuilt_from_widgets - Crow deck: widget contract, user deck config, tabs rebuilt from widgets
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: crow-deck
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Implement the widget contract and chosen config in crow deck; default config reproduces today's screens.
- Keywords: crow, deck, widget, contract, user, config, tabs, rebuilt, widgets
- Use when: Making the deck composable.
- Skip when: Kit internals or lazytf.

# Problem
- Deck screens are hard-coded; users cannot choose what they see.

# Scope
- In:
  - Implement the widget contract and the config shape chosen in the study.
  - Layout by aspect ratio (vertical screen first).
  - Rebuild existing tabs as composed widgets rendered through the kit; default config reproduces today's deck.
  - Screen-level review round through the workflow.
- Out:
  - New widgets beyond today's features.

# Acceptance criteria
- AC1: Default config renders the approved screens; a custom config adds, removes and moves widgets.
- AC2: Existing data tests pass; no tab loses a feature.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC1: Default config renders the approved screens; a custom config adds, removes and moves widgets.

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

