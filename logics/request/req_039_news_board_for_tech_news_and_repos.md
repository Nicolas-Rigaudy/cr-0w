## req_039_news_board_for_tech_news_and_repos - News board for tech news and repos
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 60
> Confidence: 55
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: A Discover board in nest that suggests tech articles and repos worth a look. A scheduled scout job gathers from the TLDR newsletter, GitHub trending, Hacker News, Lobsters and release feeds; the owner adds Instagram posts through the PWA share target. Picks can be saved to the vault repo note or turned into research and plan missions.
- Keywords: news board, discover, scout, tldr, github trending, hacker news, lobsters, share target, instagram, vault repo note, ranking
- Use when: building the scout job, the Discover card or tab, or the share target.
- Skip when: work news or anything touching work data (personal scope only).

# Needs
1. **Sources.** TLDR newsletter from personal Gmail (read-only), GitHub trending (daily, languages the owner uses), Hacker News front page, Lobsters, release feeds of tools in use (cdx, Claude Code, Logics, Bubble Tea). Each source is a script; no model for gathering.
2. **Share in.** Instagram has no API. nest declares a PWA share target; a post or link shared from any Android app becomes a pick tagged "shared by owner" and is always shown.
3. **Dedupe and rank.** Scripts dedupe by URL and repo. A cheap model ranks against the owner's interests: registered projects, the vault note "Github repos and websites", past likes and dismisses. One sentence per pick on why it matters.
4. **Actions per pick.** Like, dismiss (both feed the ranking), save to the vault note "Github repos and websites", start a `research` mission, plan a trial (`plan` mission).
5. **Placement.** A Discover card on Home with the top three picks; the full board under Personal. No push notifications for picks.
6. **Schedule.** Morning and evening runs as a core job; a budget cap per run.

# Context
- Owner idea 2026-11-06: the TLDR newsletter goes unread; nice repo posts on Instagram get attention. A small board that suggests articles and repos fits better.
- Reuses: core jobs (scripts first, narrate-only model), the `research` and `plan` kinds from req_037, the vault repo note as the store of kept picks, Web Push silent for this board.
- Depends on core (phase 1) and the nest shell; built in phase 2 with the personal screen.

# Acceptance criteria
- AC1: The scout job writes ranked picks with source, link, one-line reason and score, and runs without a model when ranking fails (plain recency order).
- AC2: A link shared from Instagram on the phone appears on the board within a minute.
- AC3: "Save" appends the pick to the vault note "Github repos and websites" in its existing format.
- AC4: "Start research" opens a `research` mission with the pick as its goal.
- AC5: Dismissed picks and their sources weigh less in the next ranking.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/product/prod_005_crow_jarvis_layer.md`
- `logics/request/req_026_personal_butler.md`
- `logics/request/req_037_mission_kinds_and_run_contract.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
