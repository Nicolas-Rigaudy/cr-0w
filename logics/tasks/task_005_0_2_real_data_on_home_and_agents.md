## task_005_0_2_real_data_on_home_and_agents - 0.2 real data on Home and Agents
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Owner: nicolas
> Indicators reviewed: 2026-11-05 16:15:07

# AI Context
- Summary: Home and Agents read real data from crow-runs, crow-sessions, crow-usage, crow-accounts and crow-registry; Agents tab with tail, answer keys and who-worked-when.
- Keywords: real, data, home, agents
- Use when: continuing the live data layer or the Agents tab of the crow deck.
- Skip when: the change is the MRs, Projects or Usage tab.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_017_0_2_real_data_on_home_and_agents`

# Acceptance criteria
- AC1: No model calls to render.
- AC2: Questions from agents can be answered from the deck.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_005_0_2_real_data_on_home_and_agents.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_005_0_2_real_data_on_home_and_agents.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: mappers on invented JSON, grouping and sort, sync states, cursor stability across refreshes, answer keys through an injected runner, relaunch confirm, fit tests for every theme, size and pane.
- `python3 -m pytest tests/test_crow_focus.py tests/test_crow_registry.py`: 6 passed.
- AC1: the deck only execs the bin/ scripts; crow-usage runs with CROW_TITLES_PER_DAY=0, so no model call renders anything.
- 2026-11-05: `python3 -m pytest -q tests` 288 passed; `go vet ./... && go test ./...` 138 passed (agentctl_test covers the confirm and the crow-focus flags).
- Renders of the Agents tab in 4 themes (wide, portrait, small, paused, typing, states) under local/state/afk/shots/deck-tabs/.
- command: `python3 -m pytest -q tests && cd tools/crow-deck && go vet ./... && go test ./...` | result: passed | date: 2026-11-05
- Finish workflow executed on 2026-11-05.
- Linked backlog/request close verification passed.

# Report
- Shipped crow-registry (registry, identities, git hygiene as JSON), crow-focus --send, the data layer (data.go), the live refresh loop (runs and sessions 2s, usage and accounts 30s, registry 60s) and the Agents tab.
- Answer path: a x e c stage text in the input line, enter sends it with `crow-focus --pid P --send TEXT` (kitty remote control). Only the command builder and dry run are tested; no live trial was done.
- 2026-11-05 (ec29940): pause (`crow-focus --interrupt`, Escape in the window) and stop (`crow-focus --stop`, SIGTERM, refused unless the pid runs claude) from the `:` palette, both behind a y confirm. attach (`g`, palette) and relaunch (`R`, behind a confirm) are wired.
- 2026-11-05 live trial of the answer path on kitty: the text was typed, but into the hub session's window. The lookup took the first window holding any ancestor of the pid. Fixed: the nearest process wins across all sockets (kitty and tmux); checked read only on the 4 live kitty windows.
- Not tried live: pause and stop on a real session (owner, from the deck).
- Home Blockers and the MR counts now read the MRs source (item_018). Today shows an honest note until calendar and mail are connected; ticket blockers are not wired.
- Finished on 2026-11-05.
- Linked backlog item(s): `item_017_0_2_real_data_on_home_and_agents`
- Related request(s): `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# AC Traceability
- request-AC1 -> This task. Proof: Implemented in ec29940 (Agents tab, answer, attach, relaunch, pause, stop; render execs bin/ scripts only); validated with pytest 288 passed and go test 138 passed; answer path tried live on kitty. Source: `ec29940`
- request-AC2 -> This task. Proof: Implemented in ec29940 (Agents tab, answer, attach, relaunch, pause, stop; render execs bin/ scripts only); validated with pytest 288 passed and go test 138 passed; answer path tried live on kitty. Source: `ec29940`
