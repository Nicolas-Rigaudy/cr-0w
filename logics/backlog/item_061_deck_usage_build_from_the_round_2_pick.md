## item_061_deck_usage_build_from_the_round_2_pick - Deck Usage build from the round 2 pick
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Build the Usage tab from the board the owner picks in item_055.
- Keywords: deck, usage, build, round, pick
- Use when: building the Usage tab.
- Skip when: never; item_055 is done.

# Problem
- The current Usage tab fails the req_011 audit; the owner asked for a merged design (see item_055).

# Scope
- In:
  - Layout as the picked item_055 board at both sizes, real data only.
  - Graphs wherever the data allows one (owner direction for Usage, applies to all).
  - Per session project, full title, model, tokens, cost; today / week / month; chart under 10 rows (req_011 AC5).
- Out:
  - Other tabs.

# Acceptance criteria
- AC1: Live render matches the picked board at 235x62 and 118x60.
- AC2: No 'not wired' box (req_011 AC7).
- AC3: go test passes with shot tests at both sizes.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.

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

# Final spec (owner approved 2026-11-03)
- Board: local/state/shots/deck-rework/png/prop-usage-R2-land.png and -port.png (scripts usage_r2.py, usage_r2_ledger.py).
- Runway (mine accounts full lanes, Team group compact), daily spend + Saved, ledger beside By project and By model. No per-session lanes.
- Account tiers (req_011 round 4): add `tier: mine|team` per account in the accounts config; mine = full rows (5h + week gauges, %, reset, empty-at / week-spent warning), team = one compact line (small 5h + week gauges, %). Not hidden. Mock helpers in local/state/shots/deck-rework/real.py (TIER, row drawers).
- Comparison page: local/state/shots/deck-rework/r2-compare.html

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

