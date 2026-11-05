## task_006_0_3_mrs_tab_both_forges_read_only - 0.3 MRs tab: both forges, read only
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.

# AI Context
- Summary: MRs tab of the crow deck: every MR or PR the owner authors, reviews or is assigned, on GitLab and GitHub, with the diff, review threads and a findings card. Read only.
- Keywords: mrs, tab, both, forges, read, only
- Use when: continuing the MRs tab, the crow-mrs script or the review findings file.
- Skip when: the change is the Projects or Usage tab.

# Definition of Done (DoD)
- [x] The backlog scope is implemented.
- [x] Acceptance criteria are covered.
- [x] Validation passes.
- [x] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# Backlog
- `item_018_0_3_mrs_tab_both_forges_read_only`

# Acceptance criteria
- AC1: Read only until 0.5: nothing is posted or pushed from this step.

# Plan
- [x] Use `python3 -m logics_manager flow progress task task_006_0_3_mrs_tab_both_forges_read_only.md --progress <n>%` during multi-wave work.
- [x] Run `python3 -m logics_manager flow finish task task_006_0_3_mrs_tab_both_forges_read_only.md` after implementation.

# Validation
- `cd tools/crow-deck && go vet ./... && go test ./...` pass: diff parser, ordering and counts, Home blockers from MRs, mode filter by scope, one detail fetch per selection with retry after the TTL, cursor kept across refreshes, a failing forge named next to the other forge's MRs, mock cards, read-only action keys, fit tests for every theme and size.
- `python3 -m pytest tests/test_crow_mrs.py`: 5 passed (GitLab and GitHub payloads invented).
- Live check: `bin/crow-mrs json` reports both forges ok; `show` returns files and threads. Output not recorded here.
- Renders of the MRs tab in 4 themes (review, own MR with threads, flipped card, portrait, single panes, states) under local/state/afk/shots/deck-tabs/.
- Finish workflow executed on 2026-11-05.
- Linked backlog/request close verification passed.

# Report
- Shipped `bin/crow-mrs` (GET only on both forges: list, diff, threads, saved findings), the MRs data layer, the MRs tab and Home Blockers and MR counts fed by the MRs source.
- AC1 holds: the tab only reads. a x e c, P and M say what arrives later and do nothing.
- Gap: no writer for the findings file yet. The review skill must save `<state>/reviews/<group>_<name>-<number>.json` (list of title, path, line, body, pattern, seen, comment, state) before live mode shows findings, patterns, fixes and replies. Until then they exist in mock data only. Record this as the next step of item_018.
- Branches ahead of the default branch with no open MR are listed for the main checkout's current branch only.
- 2026-11-05 closeout: the findings gap is closed. req_011 items 043 to 048 (Done) added the writer: the review skill saves `local/state/reviews/*.json` and `bin/crow-mrs` reads them for live findings. AC1 still holds: the tab only reads.
- Proof: `cd tools/crow-deck && go vet ./... && go test ./...` pass on branch chg/deck-0-2-finish.
- Finished on 2026-11-05.
- Linked backlog item(s): `item_018_0_3_mrs_tab_both_forges_read_only`
- Related request(s): `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`

# Links
- Request: `req_007_crow_deck_interface_redesign_tabbed_tui_4_themes_command_map`
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
