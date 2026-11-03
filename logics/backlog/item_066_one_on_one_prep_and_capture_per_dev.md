## item_066_one_on_one_prep_and_capture_per_dev - 1:1 prep and capture per dev
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Management
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Before a 1:1 crow writes a prep block into the dev's vault person note from Jira, MRs, review comments and past signals; after it, a capture step records notes with the star and warning markers and moves the due date.
- Keywords: one-on-one, prep, capture, vault, person note, objectives, signals
- Use when: preparing or recording a 1:1.
- Skip when: team-wide rituals (067).

# Problem
- Person notes exist with a 1:1 template but their goals, signals and history stay empty. 1:1s slip without notice.

# Scope
- In:
  - Skill `one-on-one <dev>`: gathers last 1:1 date, tickets done and carried, MRs and review threads, the dev's quarterly objectives, open signals. Writes a dated prep block in the person note.
  - Capture mode after the meeting: asks up to three questions, writes history lines with `- ⭐` and `- ⚠️`, ticks the recurring task (next due date).
  - Deck: a 1:1 row (from 063) runs prep on enter.
  - Generic skill in `.claude/skills/`, names only in the vault.
- Out:
  - Sharing notes with the dev.
  - Performance reviews.

# Acceptance criteria
- AC1: Prep writes a block with at least tickets, MRs and objectives for one dev.
- AC2: Capture adds history lines with markers and moves the task due date.
- AC3: Running prep twice for the same date replaces the block, no duplicates.

# AC Traceability
- All ACs -> this item. Proof: to fill at build.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_014_crow_management_checks_and_actionable_deck
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Monthly 1:1s with the devs were due 2026-11-01 and not booked.

# Notes
- Template `90 Templates/Person (1-1).md`; objectives in `10 Work/6 Management/Team Objectives & Motivation.md`.
- Vault writes follow the vault conventions (vault Claude owns triage).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
