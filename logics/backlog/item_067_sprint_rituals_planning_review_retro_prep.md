## item_067_sprint_rituals_planning_review_retro_prep - Sprint rituals: planning, review and retro prep
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
- Summary: A sprint skill preps planning (velocity, carry-over, proposal list), review (done tickets in business terms) and retro (Retro template with sprint facts), without changing Jira.
- Keywords: sprint, planning, review, retro, velocity, carry-over, jira
- Use when: preparing a sprint ritual.
- Skip when: booking the meeting (065) or editing the sprint in Jira.

# Problem
- Rituals are prepped by hand or skipped. The last retro note says meetings often get skipped.

# Scope
- In:
  - Skill `sprint <planning|review|retro>`.
  - Planning: committed vs done points per person, carry-over list, a proposed ticket list for the next sprint. Proposal only, per local runbook sprint-planning-and-ticket-status.
  - Review: done tickets grouped by theme, business wording (French).
  - Retro: note from `90 Templates/Retro.md` pre-filled with sprint facts, in the vault.
  - The 063 sprint check links its row to the matching prep.
- Out:
  - Moving tickets between sprints.
  - Posting to Teams or Confluence without OK.

# Acceptance criteria
- AC1: Planning prep for the current sprint shows velocity per person and carry-over, matching Jira.
- AC2: No Jira write happens in any mode.
- AC3: Retro mode creates one dated note from the template.

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
- Rationale: Sprint planning nearly slipped on 2026-11-03.

# Notes
- Velocity history per the sprint-velocity memory (CD board).
- The S32-S35 sprint plan lives in the vault launch notes.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
