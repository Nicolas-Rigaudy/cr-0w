## task_005_0_2_real_data_on_home_and_agents - 0.2 real data on Home and Agents
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 90%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: nicolas
> Indicators reviewed: 2026-10-29 18:27:41

# AI Context
- Summary: Home and Agents read real data from crow-runs, crow-sessions, crow-usage, crow-accounts and crow-registry; Agents tab with tail, answer keys and who-worked-when.
- Keywords: real, data, home, agents
- Use when: continuing the live data layer or the Agents tab of the crow deck.
- Skip when: the change is the MRs, Projects or Usage tab.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented (pause and stop from the `:` palette are deferred, see Report).
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_017_0_2_real_data_on_home_and_agents`

# Acceptance criteria
- AC1: No model calls to render.
- AC2: Questions from agents can be answered from the deck.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_005_0_2_real_data_on_home_and_agents.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_005_0_2_real_data_on_home_and_agents.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: mappers on invented JSON, grouping and sort, sync states, cursor stability across refreshes, answer keys through an injected runner, relaunch confirm, fit tests for every theme, size and pane.
- `python3 -m pytest tests/test_crow_focus.py tests/test_crow_registry.py`: 6 passed.
- AC1: the deck only execs the bin/ scripts; crow-usage runs with CROW_TITLES_PER_DAY=0, so no model call renders anything.
- Renders of the Agents tab in 4 themes (wide, portrait, small, paused, typing, states) under local/state/afk/shots/deck-tabs/.

# Report
- Shipped crow-registry (registry, identities, git hygiene as JSON), crow-focus --send, the data layer (data.go), the live refresh loop (runs and sessions 2s, usage and accounts 30s, registry 60s) and the Agents tab.
- Answer path: a x e c stage text in the input line, enter sends it with `crow-focus --pid P --send TEXT` (kitty remote control). Only the command builder and dry run are tested; no live trial was done.
- Deferred: pause and stop from the palette (destructive, cannot be verified unattended). attach (`g`, palette) and relaunch (`R`, behind a confirm) are wired.
- Home Blockers and Today show an honest note until the forges (item_018) and calendar/mail are connected.

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
