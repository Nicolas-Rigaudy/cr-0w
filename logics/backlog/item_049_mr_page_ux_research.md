## item_049_mr_page_ux_research - MR page UX research: review walkthrough patterns
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 80%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Study how good review tools show a multi-step review flow and whose move it is, then write the principles the deck MR page should follow.
- Keywords: UX research, code review tools, TUI, stepper, whose move, gh-dash, lazygit
- Use when: before drawing the MR page mockups.
- Skip when: mockups already chosen.

# Problem
- The deck MR page grew one fix at a time; the walkthrough (review → decide → apply → push → reply/resolve → done) has no visible spine. Research first so the mockups start from proven patterns, not from the current layout.

# Scope
- In:
  - Tools to look at (docs, screenshots, source where useful): gh-dash, lazygit, tig, GitHub PR page and its "review changes" flow, GitLab MR page (approval/thread widgets), Graphite, Reviewable (its "who has the ball" disposition), Gerrit (attention set), Phabricator/Differential.
  - Patterns to extract: whose-turn signals (attention set, ball-in-court), stage/stepper display, inline vs side cards for threads, how actions confirm and report progress, done states, list badges and counts.
  - TUI constraints: 160x45 laptop, 240x62 wide, 118x95 vertical; ANSI theme colours only.
- Out: building anything; mockups (item_050).

# Acceptance criteria
- AC1: A research note `runbooks/tui-review-ux.md` (generic, no names) lists each tool looked at with the pattern taken from it and a source link.
- AC2: The note ends with 6-10 principles for the deck MR page, each tied to a pain point in req_013 Needs.
- AC3: It proposes the whose-move rule per MR (reviewer waiting on me, open threads, red pipeline, decisions pending, commit not pushed, draft, waiting for review/merge) and 2-3 badge options for item_050.

# AC Traceability
- AC1 -> tool list with pattern and source link. Proof: `runbooks/tui-review-ux.md`, section "Tools" (9 tools).
- AC2 -> 6-10 principles tied to req_013 pain points. Proof: `runbooks/tui-review-ux.md`, section "Principles for the deck MR page" (10 principles, P1-P7).
- AC3 -> whose-move rule and badge options. Proof: `runbooks/tui-review-ux.md`, sections "Whose-move rule" and "Badge options for item_050" (3 options).

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
- Rationale: The owner reviews all MRs from the deck; a confusing flow costs time on every review. Below prod launch and reviews themselves.

# Notes
- Cheap run: Sonnet, web research plus a short read of `mrs.go` for the current states. Execution order: 1 of 3.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
