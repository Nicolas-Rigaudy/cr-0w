## item_090_redesign_the_archeotech_crow_modules_to_the_deck_direction - Redesign the archeotech crow modules to the deck direction
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 70%
> Confidence: 70%
> Progress: 0%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Rebuild the archeotech bar modules (crow-agents, crow-usage) so they match the crow deck: same four questions, same data, shared design tokens. Includes per-session cost, the leftover of item_005.
- Keywords: archeotech, quickshell, bar panel, crow-agents, crow-usage, deck, design system
- Use when: changing integrations/archeotech/crow-*.
- Skip when: deck tab work (tools/crow-deck).

# Problem
- The bar modules were designed in October for crow watch. crow watch has folded into the deck, whose look and data model moved on (UX2 boards, design system in road_001 0.5). The panels now show a different board from the deck.
- The agents panel has no per-session cost.

# Scope
- In:
  - mockup options for both modules in the deck direction, owner picks before building.
  - rebuild crow-agents and crow-usage from the same data the deck Agents and Usage tabs read, with the design-system token roles.
  - per-session cost on agent rows.
  - qmllint and headless renders in the normal and grimdark themes.
- Out:
  - deck changes; the archeotech theme mapping itself (item_084).

# Acceptance criteria
- AC1: the owner picked one mockup per module before the build.
- AC2: each module answers the same questions, in the same order, as its deck tab, from the same JSON.
- AC3: agent rows show per-session cost when known.
- AC4: qmllint silent; renders in normal and grimdark shown to the owner.

# AC Traceability
- AC1 -> This backlog slice. Proof: chosen mockups linked in the item.
- AC2 -> This backlog slice. Proof: side-by-side render of deck tab and panel.
- AC3 -> This backlog slice. Proof: panel render with cost per row.
- AC4 -> This backlog slice. Proof: qmllint output and renders.

# Dependencies
- After road_001 0.5 (design system core: widget contract and token roles), so the panels use the same tokens as the deck.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_000_crow_watch_and_bar_modules_ui_ux_pass_session_titles_savings_stats`
- Split from: `item_005_0_2_agents_board_improvements`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: glance surface; waits for the design system.

# Notes
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

