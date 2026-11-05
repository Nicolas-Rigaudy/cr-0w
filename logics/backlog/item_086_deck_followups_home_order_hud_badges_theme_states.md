## item_086_deck_followups_home_order_hud_badges_theme_states - Deck follow-ups: Home blockers order, HUD tab badges, theme empty states
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
- Summary: Home blockers take the MRs page priority order; HUD tabs get the ●N badge; each theme gets its own empty, loading and error text.
- Keywords: deck, home, blockers, priority, hud, badge, empty state, theme
- Use when: finishing the deck rework after items 051 and 058.
- Skip when: new deck pages or layout work.

# Problem
- Home Blockers rank MRs by their own rule, so Home and the MRs page can disagree on what comes first.
- The HUD spec's `●N` tab badge in ANSI slot 1 is not built.
- Empty, loading and error states look the same in every theme.

# Scope
- In:
  - Home Blockers and You decide sort by the MRs priority order (review, red pipeline, threads), sharing the order from `mrs_move.go`. Home texts stay as they are (owner, 2026-11-05).
  - HUD tab badges: `●N` in ANSI slot 1 per the theme spec.
  - Per-theme empty, loading and error text and glyph, per runbooks/tui-themes-ux.md.
- Out:
  - The full stage rule on Home (owner chose order only).
  - Layout changes.

# Acceptance criteria
- AC1: With one MR on review and one with a red pipeline, Home and the MRs page list them in the same order (Go test).
- AC2: The HUD theme renders `●N` on a tab with N MRs on the owner's move; other themes unchanged (shot test).
- AC3: Each theme renders its own empty, loading and error state; the monochrome check still tells the four apart.

# AC Traceability
- AC1-AC3 -> this item. Proof: Go tests and renders.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_011_crow_deck_rework_terminal_true_mockups`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Owner decisions 2026-11-05 after the AFK deck builds (items 051, 058).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

