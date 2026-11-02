## item_059_deck_agents_build_from_the_round_2_pick - Deck Agents build from the round 2 pick
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
- Summary: Build the Agents tab from the board the owner picks in item_055.
- Keywords: deck, agents, build, round, pick
- Use when: building the Agents tab.
- Skip when: before item_055 is picked.

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

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

