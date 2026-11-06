## item_020_0_5_actions_and_history - 0.5 actions and History
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
- Summary: Actions that change things: answer an agent, handoff to a repo session, start a review, snooze, edit roadmap. Session History behind `/` in Agents.
- Keywords: crow, deck, tui
- Use when: building this slice of the crow deck.
- Skip when: the change is outside the deck (see the other slices of req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map).

# Problem
- Part of the crow deck redesign. Spec: local/state/research/crow-deck-design.md (mirror in the vault, CR-0W/Research). Boards: local/state/research/crow-deck-boards/ (build from the boards listed at the top of the spec: `_merge` landscape boards replace `_ux2`). Decisions: local/state/research/crow-deck-critique.md.

# Scope
- In:
  - Every push, post or send stays behind an explicit OK shown in the deck.
  - History reads cdx history and transcripts.
- Out:
  - Other slices of the same request.

# Acceptance criteria
- AC1: Autonomy table of CLAUDE.md is enforced in code, not by convention.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: Autonomy table of CLAUDE.md is enforced in code, not by convention.

# Proof
- AC1 holds. `tools/crow-deck/autonomy.go` holds the Autonomy table as rules per action kind and scope; `runGated` is the only path that starts an action and refuses `never` kinds even without an overlay.
- `TestAutonomyMatchesClaudeMd` parses the CLAUDE.md table and fails on any drift.
- `TestEveryProcessGoesThroughTheGate` walks the Go AST: a new place that starts a process fails until it goes through the gate or is listed as read only or navigation.
- `TestPostMergeOnWorkIsRefused`, `TestRunGatedRefusesWithoutRunning`, `TestProposeShowsExactTextAndRunsOnce`: work merge refused, work push and mail ask, one OK runs once.
- Actions: agent answer, pause, stop, mail done and snooze, MR handoff, History resume, retro, handoff note, copy id, transcript, Projects continue, research, session, handoff, review jump, Home handoff and review. Each change shows the exact text and command in one confirm overlay.
- History: live rows read cwd, transcript and edited files from crow-usage; the preview shows touched files and the last messages from the transcript tail.
- Commits 7aec64e, 1087f89, 9c2e41c, 0bc15bb. Go tests 154 passed, bin/crow-check green.
- Shots: local/state/afk/shots/2026-11-06-deck-actions/ (README maps board to build).

# Follow-ups
- Home `x` snooze is not built: where a snooze is kept is an owner decision (deck file, a marker in open threads, or no snooze).
- Usage `P` policy edit and the MRs keys g, R and s are not wired yet.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Ships on its own; order follows the build plan in the spec.
