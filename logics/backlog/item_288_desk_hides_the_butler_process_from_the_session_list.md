## item_288_desk_hides_the_butler_process_from_the_session_list - Desk hides the butler process from the session list
> From version: 1.0.0
> Schema version: 1.0
> Status: Obsolete
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The butler runs as a headless claude process that the dispatch scan lists as a background run, so Desk offers it and its transcript opens in the main panel next to the butler panel.
- Keywords: nest, desk, butler, crow-butler, dispatch, agents, duplicate chat
- Use when: the butler conversation shows outside the butler panel.
- Skip when: the Desk Ask CR-0W / Launch box (mission draft C, see Links).

# Problem
- Owner 2026-11-25: "butler and session chat still show in double". The screenshot `owner screenshot 20261125_141736.png` shows the butler conversation in the Desk main panel (messages tagged `[nest]`, label crow) and again in the side panel.
- Cause (checked in core state 2026-11-25):
  - The butler process runs with cwd `~/.local/state/crow-butler` (`tools/crow-core/butler_cmd.go:60`).
  - The `dispatch` section lists it as `kind: run`, `cwd: ~/.local/state/crow-butler`, with the butler's last reply as `last_action`.
  - Desk builds its agents from that section, so the butler is one more pickable run, and its transcript opens in the session chat.
  - The butler itself reports "My run (crow-butler) shows as working".
- This is a second source of the double chat. The first is the Desk's own Ask CR-0W / Launch box, drafted as item C of the desk-bugs plan mission.

# Scope
- In:
  - Mark butler agents at the source: the dispatch scan (or core when it merges the section) tags a process whose cwd is the butler state dir with `role: butler`. Check how `internal/ops/tidy.go` already maps `crow-butler` to the `butler:` lane and reuse that test.
  - Nest drops `role: butler` agents from Desk lists, the work feed, deskPick defaults and the running count. Ops keeps counting the butler cost.
  - The butler's `brief · dispatch` and `brief · sessions` answers leave itself out, so it stops reporting itself as a working run.
- Out:
  - Removing the Desk Ask CR-0W / Launch box: draft C of the desk-bugs plan mission.
  - Butler internals (item_176).

# Acceptance criteria
- AC1: With the butler live, Desk lists no `crow-butler` run, and no work row opens the butler transcript.
- AC2: The running count and the hero "runs going" line ignore the butler process.
- AC3: Ops day view still shows butler cost under its lane.
- AC4: Asked "what runs now", the butler does not list itself.
- AC5: Go test on the tagging, vitest on the Nest filter; tsc, vitest and go test pass.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Overlaps: draft C "One CR-0W chat on Desk, in the side panel" in `local/state/missions/20261125-plan-nest-desk-bugs/outbox/items.md` (bug 3 of the desk-bugs mission). Build this item first or alongside; do not duplicate C.
- Related: item_210 / `data/deskPick.ts` (default pick must skip the butler).
- Superseded by: `item_278_desk_hides_the_butler_from_its_session_list`

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_288 (Desk hides the butler process from the session list) end to end. Read logics/backlog/item_288_desk_hides_the_butler_process_from_the_session_list.md first. Confirm in the live core state (curl --unix-socket local/state/core.sock /api/state, section dispatch) that the crow-butler process shows as a run, then tag it at the source, filter it in Nest (Desk, work feed, deskPick, counts) and in the butler brief. Do not touch the Desk Ask CR-0W box: that is draft C of mission 20261125-plan-nest-desk-bugs. Run go test, tsc and vitest; verify on the running Nest that the butler no longer appears on Desk.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
