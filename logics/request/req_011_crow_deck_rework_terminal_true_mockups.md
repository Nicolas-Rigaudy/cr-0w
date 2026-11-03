## req_011_crow_deck_rework_terminal_true_mockups - Crow deck rework: terminal-true mockups of every tab, themes, data wiring
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 95%
> Confidence: 90%
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
- Jira: scoped read-only token in the keyring, tested 2026-11-02 (details in local/context.md). `bin/crow-jira json` reads it: sprint tickets assigned to the owner, flagged or blocked, waiting on the owner, age; feeds Blockers and project cards.
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
- [x] Dependencies and known risks are listed: req_010 (MRs), req_012 (to-do model), admin consent from Alex for the Graph app.

# Execution order
1. Mockup session (fresh, Opus): canvas with old / current / proposal per tab + theme boards. Owner picks.
2. Quick fix, independent: AZERTY tab keys.
3. Backlog items per tab from the picked boards (Usage first: data already there).
4. Wiring: `bin/crow-jira` (token ready), roadmap `target:` + lambda release, Graph after IT, actions slice.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `local/state/research/crow-deck-audit-2026-11-02.md`
- `local/state/research/crow-deck-design.md`
- `tools/crow-deck/`, `tools/crow-watch/`

# Owner picks (2026-11-02, canvas row 'Deck rework')
- Home: A (triage table). Keep B's day timeline for another view (used for the Agents day graph).
- Agents: B's day graph (who worked when) with A's details and the run details from Home A's Right now box. Round 2 board needed.
- MRs: B (focused walkthrough) plus an easy way to see more code around the highlighted line. Round 2 board needed.
- Projects: A for content (next actions, details) but too dense; needs an airier layout and a roadmap timeline showing when the owner worked on which project. Round 2 board needed.
- Usage: A's ledger data, B's runway and daily spend, C's by model / by project / saved. Graph wherever possible. B's session lanes are redundant with Agents. Round 2 board needed.
- MR badge: count of MRs where it's the owner's move, plus a colored dot only when there is a problem.
- Themes: keep all four (lazytf, soft, cyber, hud). HUD accent: blue or magenta (default family), not amber.

# Backlog
- item_055_deck_round_2_mockups_hybrid_picks_for_agents_mrs_projects_usage_hud_accent (1)
- item_056_deck_azerty_tab_keys (independent)
- item_057_deck_home_a_build_triage_table_pulse_projects_done_today (2)
- item_058_deck_theme_set_build_lazytf_soft_cyber_hud_per_the_theme_specs
- item_059_deck_agents_build_from_the_round_2_pick
- item_060_deck_projects_build_from_the_round_2_pick
- item_061_deck_usage_build_from_the_round_2_pick

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Owner picks, round 2 (2026-11-02, page local/state/shots/deck-rework/r2-compare.html)
- Agents: R2 approved. Asked where the "Accounts today" overview lives (B had it, R2 dropped it): add a compact accounts strip back.
- MRs: neither variant. Actions (a x e c) must be far more visible, with an obvious way to advance to the next step. Keeping 3 columns is fine. No mode switch wanted; a way to scroll the code in place instead. Round 3 board needed.
- Projects: R2 approved. "who" column only earns its place when it is not the owner (waiting on someone, management items); hide "you".
- Usage: R2 approved with a fix: move By model next to By project to fill the empty space.
- Themes: the accent comes from the user's terminal theme accent, and the palette from the user's terminal theme. Each theme must feel different in monochrome, by shape alone (frames, titles, gauges, selection, chrome). HUD amber is dropped; HUD uses the shared accent slot.

# Owner picks, round 3 (2026-11-02)
- MRs: R3 approved (3 columns, code scrolls in place, Your move panel with big keys and Then). Spec for item_051.
- Projects: R2 approved (who prefix only when not the owner). Spec for item_060.
- Usage: R2 approved with a fix: accounts whose policy is never (unused) are hidden from Usage, or kept as a single small line. Home Pulse keeps the quick overview of current usage per account. Spec for item_061.
- Agents: R2 approved with a fix: Right now goes back to the top as the first thing to see; the day lanes come below it. Spec for item_059.
- Themes: Soft gauges become dots (●●●○○), HUD gauges segments (▮▮▯▯; ▰▱ are not in FiraCode), so the four themes differ in monochrome. lazytf ━─ and Cyber [///...] stay.

# Owner picks, round 4 (2026-11-03)
- Agents: R2 approved, but the accounts strip is too small and hidden between the details. It must be visible at a glance (top area, next to or under Right now).
- Accounts come in two tiers. Mine (the owner's main account, the one shared with the team, a future personal one): full detail. Team (teammates' accounts, policy never for crow): compact one-liner each (5h and week gauges, nothing more). Not hidden. Applies to Agents, Usage and Home Pulse. Needs a tier field in local/accounts.yaml (e.g. `tier: mine|team`) so a new account only needs one line.
- Themes: ok for now (lazytf lines, Soft dots, Cyber slashes, HUD segments).
- Redrawn 2026-11-03: Accounts box at the top of Agents beside Right now (under it in portrait); Usage Runway has a Team group; Home Pulse has one compact team row.
