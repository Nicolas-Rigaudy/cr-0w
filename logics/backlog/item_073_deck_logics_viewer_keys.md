## item_073_deck_logics_viewer_keys - Deck v and V keys open the Logics viewer
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
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

# Proof (2026-11-09)
- AC1: `v` starts `logics-manager view --open` detached (setsid) in the card's registry path, with the real HOME and `BROWSER=bin/crow-open`; the deck flashes, then reports within 2 s.
- AC2: `V` runs `logics-manager view --fleet --open`. logics-manager reuses its running server itself (`~/.cache/logics-manager/viewers.json`): a second start exits 0 with "Reusing the viewer already running", and the deck says "opened on the running server". Checked live with `BROWSER=true`.
- AC3: `v` flashes "v: <name> has no Logics corpus" or "v: no project selected"; a failed start flashes its last output line.
- AC4: `tools/crow-deck/logics_test.go` covers routing, the command per case, reuse, failure, other pages and the env. `go vet ./... && go test ./...` passed.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_019_release_leverage_cdx_0_25_logics_2_24`
- Commits: `347ce95`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: owner picked it from the 2026-11-03 release review.

# Notes
- Check `cdx view` first: if it adds nothing over `logics-manager view` for this case, call logics-manager directly. Check how the fleet view finds projects (registry vs its own config).

- `cdx view` has no `--fleet`; the deck calls logics-manager directly. The fleet view reads its own `fleetRoots` (viewer preferences), not the crow registry.
- Start output goes to `local/state/logics-viewer.log`.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

