## req_011_crow_deck_rework_terminal_true_mockups - Crow deck rework: terminal-true mockups of every tab, themes, data wiring
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Rework the crow deck after the owner's first real use: terminal-true mockups of every tab (old crow-watch, current deck, new proposal side by side), a theme rework that holds on a translucent terminal, then the view and data fixes.
- Keywords: crow deck, TUI, mockups, themes, usage, projects, who worked when, azerty, data wiring
- Use when: mocking, rebuilding or wiring any crow deck tab or theme.
- Skip when: MRs tab review flow (req_010).

# Needs
- The owner is unhappy with the built deck: lists that hide who must act, boxes padded with empty space, truncated titles next to free columns, wrong project names, a Usage tab that lost most of the old crow-watch detail, and a Soft theme with no visible delimitation in the terminal.
- Audit with every finding and the wiring steps: `local/state/research/crow-deck-audit-2026-11-02.md`.

# Context
- Owner direction 2026-11-02: mock every tab in Claude Design (a Design artifact canvas) with real content that renders in a terminal (16 ANSI slots, monospace grid, no fills, translucent background). For each tab show three columns: old page (crow-watch, where it existed), current deck, new proposal. Themes need a rework, Soft has no delimitation at all.
- Old and current renders: `tools/crow-watch` and `tools/crow-deck` shot tests + `bin/ansi-shot`; live renders via a temporary test calling `goLive()` and every fetcher (see audit method).
- Previous canvas: "Crow Deck designs" artifact (Nsdtzva5oHxcNosKGfQGF2), spec `local/state/research/crow-deck-design.md`.
- MRs tab is rebuilt under req_010; mock it last, from that state.
- Today box: Microsoft Graph through the owner's own Entra app "crow" (registered 2026-11-02, no IT needed; ids in local/context.md). Device-code login once, refresh token in secret-tool, `bin/crow-today` writes the Today data (events, unread, flagged, drafts, shared support@ mailbox). Mail.ReadWrite lets crow write drafts; no Mail.Send, the owner sends.
- Next list depends on req_012 (to-do model); mock it with the actor column the audit asks for.

# Scope
- In: mockups of Home, Agents (+ who worked when, history), Projects, Usage, MRs; theme rework (4 themes, focused vs unfocused borders that read on a translucent terminal); AZERTY keys; view fixes; wiring of Jira, roadmap dates, lambda release, Graph (once IT answers), actions slice.
- Out: the to-do model itself (req_012), MR review flow (req_010).

# Acceptance criteria
- AC1: A Design canvas holds, per tab, old / current / proposal side by side, built from real data at a terminal size (235x62 landscape and 118x60 portrait), using only glyphs and colors a terminal draws.
- AC2: Each theme is shown on Home and one dense tab, with visible box delimitation on a translucent background; Soft gets real separators.
- AC3: The owner picks a proposal per tab before any Go code changes.
- AC4: `1-5` and `& é " ' (` both switch tabs; no key needs AltGr on AZERTY.
- AC5: Usage shows per session project, full title, model, tokens and cost, with today / week / month, and its chart fits in under 10 rows.
- AC6: Projects names come from the repo name, NEXT is a real next action, lambda shows its Dec 1 release.
- AC7: Each unwired source (Jira, Graph, roadmap dates, actions) is either live or removed from the screen, never a "not wired" box.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed: req_010 (MRs), req_012 (to-do model), IT answer for Graph, Atlassian API token from the owner.

# Execution order
1. Mockup session (fresh, Opus): canvas with old / current / proposal per tab + theme boards. Owner picks.
2. Quick fix, independent: AZERTY tab keys.
3. Backlog items per tab from the picked boards (Usage first: data already there).
4. Wiring: Jira token + `bin/crow-jira`, roadmap `target:` + lambda release, Graph after IT, actions slice.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `local/state/research/crow-deck-audit-2026-11-02.md`
- `local/state/research/crow-deck-design.md`
- `tools/crow-deck/`, `tools/crow-watch/`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
