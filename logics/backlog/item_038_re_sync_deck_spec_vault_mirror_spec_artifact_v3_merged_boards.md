## item_038_re_sync_deck_spec_vault_mirror_spec_artifact_v3_merged_boards - Re-sync deck spec: vault mirror, spec artifact v3, merged boards
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Bring the vault mirror, the "Crow Deck Spec" artifact and the board copies in line with the critique decisions and the `_merge` / `MrsT_*` boards.
- Keywords: sync, deck, spec, vault, mirror, artifact, merged, boards
- Use when: the deck spec or boards changed after the last sync (29/10 Opus critique).
- Skip when: the build slices (item_016 and later) are the topic.

# Problem
- The spec in the repo has the critique decisions. The vault mirror (CR-0W/Research) and the spec artifact predate them. The canvas has new boards (`_merge`, MrsT_hud/cyber/soft).

# Scope
- In:
  - Read the live canvas first, copy new boards into local/state/research/crow-deck-boards/.
  - Copy spec and boards to the vault CR-0W/Research/.
  - Republish the artifact as v3 with the new boards inline.
  - Check the spec page list names the `_merge` boards and the command map.
- Out:
  - Redrawing any board.

# Acceptance criteria
- AC1: Repo, vault and artifact hold the same spec text and the same board files (cmp clean).
- AC2: The spec names the `_merge` boards as the build reference.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: bounded delivery slice.
- request-AC2 -> This backlog slice. Proof: promotable backlog item.
- request-AC3 -> This backlog slice. Proof: delivery chain includes a task-ready backlog item.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Do before item_016 so the build session reads one consistent spec.

# Notes
- Generated locally by logics-manager.
