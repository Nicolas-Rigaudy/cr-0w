## item_190_laptop_missions_from_the_approved_missions_board - Laptop Missions from the approved missions board
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: (unfilled: replace before this doc is used)
- Keywords: laptop, missions, approved, board
- Use when: (unfilled: replace before this doc is used)
- Skip when: (unfilled: replace before this doc is used)

# Problem
- Missions on the laptop renders the Fold layout; the laptop board (missions-board, missions-board-cyber) is approved but not built.

# Scope
- In:
  - Build the laptop Missions layout from local/state/shots/nest-rebuild/boards/missions-board.png and missions-board-cyber.png on live data: kp title row, 5 columns, lane timeline, 7-day runs chart with labelled axes.
  - A mission row opens its chat (/desk/m/<id>).
- Out:
  - New mission actions beyond item_147 and item_168.

# Acceptance criteria
- Side-by-side compare page of the build vs both boards gets the owner's OK.
- Every column and the chart read live core data; chart axes are labelled.
- Opening a mission from the laptop shows its last run chat.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: Side-by-side compare page of the build vs both boards gets the owner's OK.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_009_nest_laptop_parity`
- Architecture decision(s): (none yet)
- Request: `req_057_nest_on_the_laptop_replaces_the_owner_s_terminals`
- Primary task(s): `task_023_orchestrate_nest_laptop_parity`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

