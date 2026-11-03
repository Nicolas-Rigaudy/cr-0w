## item_073_deck_logics_viewer_keys - Deck v and V keys open the Logics viewer
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Deck Projects page: v opens the selected project's Logics viewer, V the fleet view.
- Keywords: deck, projects, logics viewer, keys
- Use when: adding the viewer keys to tools/crow-deck.
- Skip when: other req_019 slices.

# Problem
- Opening a project's Logics corpus means leaving the deck and typing `logics-manager view` in the right directory.

# Scope
- In:
  - `v` on a Projects card whose project has Logics (card `Logics` flag): run `logics-manager view --open --focus` detached in the project path from the registry.
  - `V` on the Projects page: one `logics-manager view --fleet --open` across Logics projects; reuse a running server when there is one.
  - Flash line on start and on error; `v` on a non-Logics card says so.
  - Footer and `?` help list both keys; `v` and `V` stay free on other pages.
  - Open the browser with the real HOME, not the cdx sandbox HOME.
- Out:
  - Embedding the viewer in the TUI.
  - LAN mode.

# Acceptance criteria
- AC1: `v` on a Logics card opens that project's viewer in the browser; the deck stays responsive.
- AC2: `V` opens the fleet view; a second `V` reuses the server.
- AC3: `v` on a non-Logics card or with no selection flashes a clear message.
- AC4: Go tests cover key routing and the command built for each case.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1 to AC4 above.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner picked it from the 2026-11-03 release review.

# Notes
- Check `cdx view` first: if it adds nothing over `logics-manager view` for this case, call logics-manager directly. Check how the fleet view finds projects (registry vs its own config).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

