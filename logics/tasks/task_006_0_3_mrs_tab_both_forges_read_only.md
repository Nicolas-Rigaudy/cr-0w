## task_006_0_3_mrs_tab_both_forges_read_only - 0.3 MRs tab: both forges, read only
> From version: 1.0.0
> Schema version: 1.0
> Status: In progress
> Understanding: 90%
> Confidence: 85%
> Progress: 85%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: MRs tab of the crow deck: every MR or PR the owner authors, reviews or is assigned, on GitLab and GitHub, with the diff, review threads and a findings card. Read only.
- Keywords: mrs, tab, both, forges, read, only
- Use when: continuing the MRs tab, the crow-mrs script or the review findings file.
- Skip when: the change is the Projects or Usage tab.

# Definition of Done (DoD)
- [ ] The backlog scope is implemented.
- [ ] Acceptance criteria are covered.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_018_0_3_mrs_tab_both_forges_read_only`

# Acceptance criteria
- AC1: Read only until 0.5: nothing is posted or pushed from this step.

# Plan
- [ ] Use `python3 -m logics_manager flow progress task task_006_0_3_mrs_tab_both_forges_read_only.md --progress <n>%` during multi-wave work.
- [ ] Run `python3 -m logics_manager flow finish task task_006_0_3_mrs_tab_both_forges_read_only.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: diff parser, ordering and counts, Home blockers from MRs, mode filter by scope, one detail fetch per selection with retry after the TTL, cursor kept across refreshes, a failing forge named next to the other forge's MRs, mock cards, read-only action keys, fit tests for every theme and size.
- `python3 -m pytest tests/test_crow_mrs.py`: 5 passed (GitLab and GitHub payloads invented).
- Live check: `bin/crow-mrs json` reports both forges ok; `show` returns files and threads. Output not recorded here.
- Renders of the MRs tab in 4 themes (review, own MR with threads, flipped card, portrait, single panes, states) under local/state/afk/shots/deck-tabs/.

# Report
- Shipped `bin/crow-mrs` (GET only on both forges: list, diff, threads, saved findings), the MRs data layer, the MRs tab and Home Blockers and MR counts fed by the MRs source.
- AC1 holds: the tab only reads. a x e c, P and M say what arrives later and do nothing.
- Gap: no writer for the findings file yet. The review skill must save `<state>/reviews/<group>_<name>-<number>.json` (list of title, path, line, body, pattern, seen, comment, state) before live mode shows findings, patterns, fixes and replies. Until then they exist in mock data only. Record this as the next step of item_018.
- Branches ahead of the default branch with no open MR are listed for the main checkout's current branch only.

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
