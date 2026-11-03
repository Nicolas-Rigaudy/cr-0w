## item_080_kit_components_extracted_from_crow_deck_and_lazytf - Kit components extracted from crow deck and lazytf
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: High
> Theme: design
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Move deck and lazytf chrome into kit components as pure render functions with all states, approved via review rounds.
- Keywords: kit, components, extracted, crow, deck, lazytf
- Use when: Building or reviewing kit components.
- Skip when: Deck layout or config work.

# Problem
- Deck chrome lives in a flat package mixed with data code; lazytf has equivalent pieces in internal/ui.

# Scope
- In:
  - Inventory both apps' chrome; list overlaps.
  - Components as pure render functions (data + size in, string out) with states empty, loading, error, overflow, focused: frame/box with corner counter, tab bar, table with header, gauge with cap marker, sparkline, status bar, modal, title bar.
  - Encode the approved _ux2 look and the req_011 round-2 picks.
  - Component review rounds through the workflow.
- Out:
  - Switching the apps over (later items).

# Acceptance criteria
- AC1: Each component appears in the gallery in every state and both orientations.
- AC2: Owner approves the component rounds; picks in the ledger.
- AC3: Unit tests on render output widths and states.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1: Each component appears in the gallery in every state and both orientations.

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

