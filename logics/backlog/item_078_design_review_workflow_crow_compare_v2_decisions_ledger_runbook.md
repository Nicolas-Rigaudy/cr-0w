## item_078_design_review_workflow_crow_compare_v2_decisions_ledger_runbook - Design review workflow: crow-compare v2, decisions ledger, runbook
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: design
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-compare v2 with open/changed/decided sections, decision IDs, copy-picks, render cache, plus the runbook.
- Keywords: design, review, workflow, crow, compare, decisions, ledger, runbook
- Use when: Starting the design workflow; first item to build.
- Skip when: Kit or deck code changes.

# Problem
- Design rounds mix three channels, pages go stale, changes are hard to spot and decisions get repeated.

# Scope
- In:
  - crow-compare: fixed sections in order: Needs decision (options as columns, recommendation tag, flip viewer), Changed since last round (before/after flip per changed render, optional diff overlay), Decided (collapsed ledger).
  - Decision IDs (D<n>) and option letters; a copy-picks button that copies 'D7=B D8=A note: ...'.
  - Previous-round snapshot kept per topic so 'changed' is computed, not hand-written.
  - Render cache by content hash: only changed scenes re-render.
  - design/decisions.md ledger format (ID, date, question, pick, note, superseded-by).
  - runbooks/terminal-design-workflow.md: layers (system, components, screens), max ~3 decisions per round, real data fixtures, both orientations, monochrome check, text-first verification of ANSI before reading PNGs. Index it; mark design-canvas.md as not for terminal UIs.
- Out:
  - In-deck decision widget (later item).
  - Artifact or Claude Design output.

# Acceptance criteria
- AC1: A sample spec with two rounds shows open decisions, a computed changed section and a decided ledger; a re-run with no change shows no stale content.
- AC2: Copy-picks output parses back into the ledger with a script.
- AC3: Runbook indexed in runbooks/README.md.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: A sample spec with two rounds shows open decisions, a computed changed section and a decided ledger; a re-run with no change shows no stale content.
- request-AC2 -> This backlog slice. Proof: AC2: Copy-picks output parses back into the ledger with a script.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)
- Request: `req_020_terminal_design_system_and_design_review_workflow`
- Primary task(s): `task_010_orchestrate_terminal_design_system`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

