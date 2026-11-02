## item_050_mr_page_walkthrough_mockups - MR page mockups: walkthrough layout options to choose from
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 80%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Draw 2-3 layout options for the MR page from item_049's principles, render them as terminal-true images at three sizes, and let the owner pick.
- Keywords: mockups, ansi-shot, crow-compare, MR page, walkthrough
- Use when: after item_049.
- Skip when: the owner already picked a layout.

# Problem
- The owner decides visual layouts from rendered images side by side, never from ASCII or prose.

# Scope
- In:
  - 2-3 options, each shown at four stages of the real !214 walkthrough: deciding, applying (running), pushed with threads to resolve, done. Plus the MR list with one your-move and one waiting MR, and the tab badge.
  - Written as ANSI (or from a throwaway View() variant) and rendered with `bin/ansi-shot --glass 1b5b70` at laptop, wide, vertical.
  - One `bin/crow-compare` page: rows = stage x size, columns = options; opened with xdg-open and the real HOME.
  - Ask the choice with AskUserQuestion (self-contained, no ASCII previews), including the badge option.
- Out: building the chosen option (item_051).

# Acceptance criteria
- AC1: The compare page shows every option at every stage listed above, at three sizes.
- AC2: The owner's pick (layout and badge) is recorded in this doc with the date.

# AC Traceability
- request-AC6 -> owner pick from rendered options. Proof: compare page and recorded choice.
- request-AC3 -> badge option picked. Proof: recorded choice.

# Decision framing
- Product framing: Owner decision
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_013_mr_page_walkthrough_rework
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Gate for item_051.

# Notes
- Pick 2026-11-02 (req_011 canvas, boards Rw_MRs_*): layout B focused walkthrough, plus a way to see more code around the highlighted line; round 2 board in item_055. Badge: count of MRs on the owner's move, colored dot only when there is a problem.
- Depends on item_049. Execution order: 2 of 3. Real data: `local/state/shots/deck-threads/214.json` and the review file for !214.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
