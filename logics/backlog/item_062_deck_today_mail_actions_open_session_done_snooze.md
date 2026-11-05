## item_062_deck_today_mail_actions_open_session_done_snooze - Deck Today mail actions: open, crow session, done, snooze
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Mail rows in the deck Today rail get three actions: o opens the mail in Outlook, enter starts a crow session that drafts a reply, x marks the flag done and s snoozes it to tomorrow.
- Keywords: deck, today, mail, actions, open, crow, session, done, snooze, graph, flag
- Use when: wiring actions on the Today mail rows (tools/crow-deck, bin/crow-today).
- Skip when: the calendar rows, mouse support, other tabs.

# Problem
- The Today rail lists flagged mail (crow's and the owner's) but nothing can be done from it: enter flashes "actions are not wired yet".

# Scope
- In:
  - `bin/crow-today`: flagged rows carry `id`, `webLink`, `conversationId`; subcommands `done <id>` (flag complete) and `snooze <id>` (flag due tomorrow 08:00 local, hidden from Today until then).
  - Deck: the cursor walks mail rows in the Today pane (calendar rows are skipped); `o` runs `bin/crow-open <webLink>` (real browser); `x` and `s` call the subcommands, flash the result and refetch Today.
  - `enter`: a hub session about the mail, via the existing session launcher (`crow handoff` pattern, hub project). Prompt: the mail id, subject and sender, the rule "read the thread with Graph, say what is needed, draft the reply into Outlook Drafts in French, never send". Owner prefs go in the prompt.
  - `bin/crow-today thread <id>`: prints the thread (sender, date, plain body, newest first) so the session reads it without the claude.ai connector.
  - Key hints in the footer when the Today pane has focus.
- Out:
  - Mouse support (owner declined 2026-11-03).
  - Sending mail (the app has no Mail.Send).
  - Calendar row actions.

# Acceptance criteria
- AC1: On a crow or owner flag, `o` opens the mail in Outlook on the web in the real browser.
- AC2: `x` marks the flag complete in Outlook and the row leaves Today on the next fetch; `s` sets it due tomorrow and hides it until then.
- AC3: `enter` starts a hub session whose first answer summarises the thread and leaves a reply draft in Outlook Drafts, nothing sent.
- AC4: Unit tests cover the cursor on mail rows only, the key to command mapping, and the snooze filter; `go test ./...` passes.

# AC Traceability
- AC1 -> `todayKey` "o" runs `crow-open <webLink>`; webLink comes from Graph (`bin/crow-today json`). Proof: TestTodayKeysRunTheirCommands.
- AC2 -> `crow-today done|snooze` + `awake()` filter. Proof: live snooze 2026-11-03 (mail left the picker, came back with until Wed 08:00, flag cleared after), TestSnoozedFlagsHideUntilTheyReturn.
- AC3 -> enter runs `crow handoff cr-0w` with the thread/draft task; `crow-today thread` checked live, `draft` uses createReply/createReplyAll. Proof: TestTodayKeysRunTheirCommands; live 2026-11-05 by the owner: enter on a real flagged mail left a reply draft in Outlook Drafts, nothing sent (after fix 45241b4; owner deleted the test draft).
- AC4 -> today_test.go (cursor on mail rows, key mapping, snooze filter, mock rows run nothing); `go test ./...` passes.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_011_crow_deck_rework_terminal_true_mockups
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Owner asked 2026-11-03 right after the Today flags went live.

# Notes
- Owner picks 2026-11-03: o open, enter crow session, x done / s snooze; no mouse.
- Today today: `bin/crow-today` (runbooks/graph-today.md), deck `tools/crow-deck/today_data.go`, rendering `home.go` todayLine. Home key handling: `main.go` (enter on page 0 flashes "not wired").
- Reuse: actionCmd/runBin in live.go for x/s/o; the MR handoff call in main.go (~l.447) for enter. Check `crow handoff --help` for a hub target before adding a new launcher.
- Depends on nothing; independent of 059-061.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
