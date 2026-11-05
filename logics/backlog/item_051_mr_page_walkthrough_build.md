## item_051_mr_page_walkthrough_build - MR page build: chosen walkthrough layout and whose-move states
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 80%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-05 11:33:57

# AI Context
- Summary: Build the layout the owner picked in item_050: stage line, whose-move list states, badge, done state, refresh after actions.
- Keywords: crow deck, mrs.go, build, whose move, badge, done
- Use when: after the owner picked a layout in item_050.
- Skip when: no pick recorded yet.

# Problem
- See req_013 Needs. The current patches (6f0f1c2, 4245910) fix symptoms; this slice replaces them with the chosen design.

# Scope
- In:
  - Stage line per the pick: stage, whose move, next key.
  - List: your move highlighted, waiting greyed, reason shown; same rule feeds Home Blockers.
  - Badge per the pick (`chrome.go`, `fixtures.mrsNeedYou`).
  - Done state on card and list.
  - Immediate detail refetch after `postDoneMsg` and apply start (no wait for the next tick).
  - Tests in `tools/crow-deck` for each stage and state; renders at three sizes compared to the mockups.
- Out: crow-review-run changes.

# Acceptance criteria
- AC1: Renders of the real !214 at the four stages match the picked mockup at three sizes.
- AC2: Go tests cover stage line, list states, badge and done state.
- AC3: `runbooks/review-findings.md` deck section updated.

# AC Traceability
- request-AC1 -> stage line. Proof: tests and renders.
- request-AC2 -> list states. Proof: tests.
- request-AC3 -> badge. Proof: tests.
- request-AC4 -> done state. Proof: tests and renders.
- request-AC5 -> refetch after actions. Proof: test on postDoneMsg.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_013_mr_page_walkthrough_rework
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Follows the owner's pick.

# Notes
- Depends on item_050. Execution order: 3 of 3. Personal repo, commit on main allowed.
- Delivered 2026-11-04: `mrs_move.go` holds the stage rule (stage line, panel, list groups, badge). `mrs.go` draws the three columns. Tests in `mrs_move_test.go`. Real !214 renders at four stages and three sizes: env-gated `TestShotsReal` (SHOTS + REAL_MR).
- AC1 met, AC2 met (`mrs_move_test.go`), AC3 met (runbook section "Deck page").
- Left for the owner: Home Blockers still use `mrBlockers` (row data only, loaded before MR details), not the stage rule.

# Final spec (owner approved 2026-11-03)
- Board: MRs R3 approved 2026-11-02 (req_011 round 3), supersedes the item_050 pick: local/state/shots/deck-rework/png/prop-mrs-R3-land.png and -port.png (scripts mrs_r3.py, mrs_r2.py, mrs_lib.py).
- 3 columns: MR list | code scrolling in place (j k, J K by hunk, [ ] file, scrollbar + line position) | Your move panel (steps with current marked, finding, fix, big keycaps a x e c, Then: n next, A apply, P push). Commented line reversed with markers. Badge = count of MRs on the owner's move, dot only on a real problem.
- Comparison page: local/state/shots/deck-rework/r2-compare.html

# Owner decisions 2026-11-05
- Portrait list key stays `tab` (m is the mode switch). Keycap detail stays on its own line. Home Blockers take the MRs priority order (review, red pipeline, threads) and keep their texts: follow-up item.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
