## item_056_deck_azerty_tab_keys - Deck AZERTY tab keys
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 95%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: AZERTY tab keys: `& é " ' (` switch tabs like `1-5`; no key needs AltGr.
- Keywords: deck, azerty, tab, keys
- Use when: any time; independent quick fix.
- Skip when: layout work.

# Problem
- On the owner's AZERTY keyboard the digit row sends & é " ' (, so 1-5 need Shift.

# Scope
- In:
  - Map & é " ' ( to tabs 1-5 in tools/crow-deck key handling, keep 1-5.
  - Audit every bound key for AltGr-only characters and rebind them.
  - Footer and ? help show the AZERTY keys.
- Out:
  - Layout changes.

# Acceptance criteria
- AC1: & é " ' ( and 1-5 both switch tabs (req_011 AC4).
- AC2: No binding needs AltGr on AZERTY.
- AC3: go test passes.

# AC Traceability
- request-AC3 -> owner pick before Go code. Proof: req_011 § Owner picks.
- AC1, AC2, AC3 -> Implemented in 21c6fbe; validated with go test ./... (TestAzertyTabKeys, TestHourRangeKeysWithoutAltGr).

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Daily friction, small change.

# Notes
- Execution order: independent, any time.
- Done 2026-11-03 in 21c6fbe (branch deck-azerty): `& é " ' (` switch tabs 1-5; `< >` replace the AltGr-only `[ ]` for the Agents hour range and Usage priority (`[ ]` kept); footer, ? help and runbooks/tui-usage-ux.md updated. Validated with `go vet ./... && go test ./...` in tools/crow-deck (keys_test.go covers AC1 and the `< >` keys).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
