## item_064_deck_action_model_every_row_actionable - Deck action model: every row actionable
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Management
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:48

# AI Context
- Summary: Every Home row carries an action spec: a default enter action plus shared keys (o open, enter session, x done, s snooze, d dismiss), with key hints in the footer.
- Keywords: deck, actions, keys, enter, open, snooze, dismiss, session
- Use when: wiring keys on any deck row.
- Skip when: new data sources.

# Problem
- Most deck rows are read only. Enter flashes 'not wired' on Home.

# Scope
- In:
  - Row type gains `Actions` (kind, label, command). Sources fill it: MR rows open the MR, review rows start the review session, check rows run the suggested action, Today mail rows follow item_062.
  - One dispatcher in main.go maps keys to the row's actions; missing actions are hidden in the footer.
  - Session actions reuse the `crow handoff` launcher with a prompt built from the row.
  - Done, snooze and dismiss write through the source's own command (crow-checks, crow-today). Home thread snooze writes the open-threads marker (item_114, owner decision 2026-11-07).
- Out:
  - Mouse.
  - Outward actions without a draft step (065 owns invites).

# Acceptance criteria
- AC1: Each Home pane row with an action opens, starts a session or updates state on enter.
- AC2: The footer lists only the keys the focused row supports.
- AC3: Unit tests cover the key to action mapping per row kind; `go test ./...` passes.

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
- Priority: High
- Rationale: Owner wants every item in the deck actionable (2026-11-03).

# Notes
- item_062 is the first instance (Today mail). Build it first or fold it in here.
- Reuse actionCmd/runBin in live.go.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
