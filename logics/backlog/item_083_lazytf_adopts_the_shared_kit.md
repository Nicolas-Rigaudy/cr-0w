## item_083_lazytf_adopts_the_shared_kit - lazytf adopts the shared kit
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: lazytf
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: lazytf swaps its chrome for kit components with no visual change, done in a lazytf session.
- Keywords: lazytf, adopts, shared, kit
- Use when: Adopting the kit in lazytf.
- Skip when: Hub-only work.

# Problem
- lazytf keeps its own copy of chrome that the kit now owns.

# Scope
- In:
  - Run in a lazytf session (crow handoff lazytf).
  - Replace titlebar, statusbar, modal and frame rendering with kit components; keep lazytf behaviour.
- Out:
  - lazytf layout customization (possible later with the widget contract).

# Acceptance criteria
- AC1: Before/after renders show no unintended visual change.
- AC2: lazytf tests pass.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: AC1: Before/after renders show no unintended visual change.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_004_terminal_design_system`
- Architecture decision(s): (none yet)
- Request: `req_020_terminal_design_system_and_design_review_workflow`
- Primary task(s): `task_010_orchestrate_terminal_design_system`

# Priority
- Priority: Low
- Rationale: Set by scaffold input or defaulted for grooming.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

