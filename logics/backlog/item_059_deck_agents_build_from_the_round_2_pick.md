## item_059_deck_agents_build_from_the_round_2_pick - Deck Agents build from the round 2 pick
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build the Agents tab from the board the owner picks in item_055.
- Keywords: deck, agents, build, round, pick
- Use when: building the Agents tab.
- Skip when: never; item_055 is done.

# Problem
- The current Agents tab fails the req_011 audit; the owner asked for a merged design (see item_055).

# Scope
- In:
  - Layout as the picked item_055 board at both sizes, real data only.
  - Graphs wherever the data allows one (owner direction for Usage, applies to all).
  - Day lanes per agent from activity buckets, plus the run details (tail, question, detail pane).
- Out:
  - Other tabs.

# Acceptance criteria
- AC1: Live render matches the picked board at 235x62 and 118x60.
- AC2: No 'not wired' box (req_011 AC7).
- AC3: go test passes with shot tests at both sizes.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.
- AC1 -> live renders at 235x62 and 118x60. Proof: `LIVESHOTS=<dir> go test -run TestLiveShots ./tools/crow-deck` (agents-live-*.ans); PNGs in local/state/shots/deck-agents/png/ compared with prop-agents-R2-land/port (commit ae415f9).
- AC2 -> every box shows real data or a one-line empty state; footer keys are all wired (enter, j k, a x e c, r, [ ], w, d, tab, /). Proof: `tools/crow-deck/agents.go`, `TestLiveAgentsFitEverySize`.
- AC3 -> `go test ./...` passes; `TestShots` writes Agents board (235x62) and portrait (118x60) mock shots; `TestLanesJoinAgentsAndMergeRuns`, `TestLaneGitParsesNumstat`; crow-usage `test_activity_tokens_per_bucket_and_edited_files`.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Follows the round 2 pick.

# Notes
- Depends on item_055. Execution order: after item_057, Usage first (data already there).
- Build notes 2026-11-03: crow-usage (cache v11) adds per-bucket tokens and edited files per session. Lanes join usage sessions to live agents by session id; same-titled runs share one lane. Detail reads the lane checkout with git (numstat, diff, last commit, 30 s TTL). Not built (no data source yet): the plan checklist in Detail, `s sort` and `/ filter` on the lanes title (/ stays History), `h handoff`. The cursor opens on the first lane that waits for an answer.

# Final spec (owner approved 2026-11-03)
- Board: local/state/shots/deck-rework/png/prop-agents-R2-land.png and -port.png (script agents_r2.py).
- Order: Right now first (top left), Accounts today box beside it (under it in portrait), then day lanes with A's columns and their legend on its own row, then tail + agent-asks card and Detail.
- Account tiers (req_011 round 4): add `tier: mine|team` per account in the accounts config; mine = full rows (5h + week gauges, %, reset, empty-at / week-spent warning), team = one compact line (small 5h + week gauges, %). Not hidden. Mock helpers in local/state/shots/deck-rework/real.py (TIER, row drawers).
- Comparison page: local/state/shots/deck-rework/r2-compare.html

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
