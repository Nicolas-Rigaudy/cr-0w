## item_057_deck_home_a_build_triage_table_pulse_projects_done_today - Deck Home A build: triage table, pulse, projects, done today
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build Home direction A in the Go deck: Pulse per account, one Next list grouped You decide / Crow works / Waiting on others, Projects table, Done today, Right now and Today rail.
- Keywords: deck, home, build, triage, table, pulse, projects, done, today
- Use when: building the Home tab.
- Skip when: other tabs.

# Problem
- The current Home hides who must act and pads boxes with empty space.

# Scope
- In:
  - Layout as board Rw_Home_A (mock: local/state/shots/deck-rework/home_a.py) at landscape and portrait.
  - Next reads the sectioned open threads (req_012, item_054) for the actor column.
  - Done today reads the day's activity log, newest first.
  - Today rail: Graph data when wired, otherwise hidden (req_011 AC7).
  - Keep Home B's day timeline idea for the Agents day graph (item_055).
- Out:
  - Agents, Projects, Usage tabs.

# Acceptance criteria
- AC1: Live render matches Rw_Home_A at 235x62 and 118x60 with real data.
- AC2: No empty padded boxes; no 'not wired' box.
- AC3: go test passes and a shot test covers both sizes.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.
- AC1 -> live renders at 235x62 and 118x60, all and work modes. Proof: `LIVESHOTS=<dir> go test -run TestLiveShots ./tools/crow-deck`; PNGs in local/state/shots/deck-rework/build/ compared with prop-home-A-land/port/land-work.
- AC2 -> panes take their content height; Done and Today hide without data; the rail keeps Right now alone when Today has nothing. Proof: `tools/crow-deck/home.go` stack/floor, `TestHomeBreakpoints`, `TestHomeStates`.
- AC3 -> `go test ./...` passes; `TestShots` (mock) and `TestLiveShots` (real) cover both sizes.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Picked as is; Home is the daily entry point.

# Notes
- Depends on item_054 for the actor column. Execution order: 2.
- Owner picks 2026-11-02 recorded in req_011 § Owner picks.
- Build notes 2026-11-03: Next keeps the full row highlight used by every tab (board shows a `›` marker). Projects get up to 11 rows before Next scrolls. Today rail reads `bin/crow-today` (Graph app crow, consent granted 2026-11-02, login done 2026-11-03; runbooks/graph-today.md). Thread tags match the registry by name, short name or checkout folder, so `[Nicolas-Rigaudy]` counts as personal.

# Final spec (owner approved 2026-11-03)
- Board: local/state/shots/deck-rework/png/prop-home-A-land.png, prop-home-A-port.png, prop-home-A-land-work.png (script home_a.py).
- Pulse: Account tiers (req_011 round 4): add `tier: mine|team` per account in the accounts config; mine = full rows (5h + week gauges, %, reset, empty-at / week-spent warning), team = one compact line (small 5h + week gauges, %). Not hidden. Mock helpers in local/state/shots/deck-rework/real.py (TIER, row drawers). Team accounts share one compact row.
- Comparison page: local/state/shots/deck-rework/r2-compare.html

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

