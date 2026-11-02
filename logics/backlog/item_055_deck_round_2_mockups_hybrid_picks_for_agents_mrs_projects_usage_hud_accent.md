## item_055_deck_round_2_mockups_hybrid_picks_for_agents_mrs_projects_usage_hud_accent - Deck round 2 mockups: hybrid picks for Agents, MRs, Projects, Usage, HUD accent
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
- Summary: Second mockup round for the tabs where the owner picked a hybrid: Agents, MRs, Projects, Usage, plus the HUD accent.
- Keywords: deck, round, mockups, hybrid, picks, agents, mrs, projects, usage, hud, accent
- Use when: before any Go change on Agents, MRs, Projects or Usage.
- Skip when: Home (A picked as is) and the AZERTY fix.

# Problem
- The owner did not pick one proposal as drawn for four tabs. Each needs a merged board before the build (req_011 AC3).

# Scope
- In:
  - Agents: B's day lanes graph (who worked when) as the main view, with A's details (table, tail, question card, detail pane) and the run details from Home A's Right now box.
  - MRs (req_013 item_050): B focused walkthrough plus a way to see more code around the highlighted line (expand context key, or a wider diff pane on demand). Show both variants if unsure.
  - Projects: A's content (next actions, who column, all details) with more air: fewer columns per row, grouped spacing, detail pane. Add a roadmap timeline that also shows when the owner worked on which project (from activity logs and git).
  - Usage: A's ledger data, B's runway and daily spend, C's by model / by project / saved. A graph wherever the data allows one. Drop B's per-session lanes (Agents covers them).
  - HUD theme redrawn with a blue or magenta accent (the owner's default family), not amber.
  - Landscape 235x62 and portrait 118x60, real data, same toolkit (mock.py, real.py, themes.py), published as new boards on the canvas.
- Out:
  - Home (picked), Go code.

# Acceptance criteria
- AC1: One merged board per tab (Agents, MRs, Projects, Usage) at both sizes on the canvas, plus HUD with the new accent.
- AC2: The owner's pick per tab is recorded in req_011 with the date.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.

# Decision framing
- Product framing: Owner decision
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Gate for items 059, 060, 061 and item_051.

# Notes
- Depends on nothing. Execution order: 1.
- Design canvas Nsdtzva5oHxcNosKGfQGF2, row 'Deck rework' (boards Rw_*); mock scripts in local/state/shots/deck-rework/
- Theme specs: runbooks/tui-themes-ux.md. Tab research: runbooks/tui-*-ux.md.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

