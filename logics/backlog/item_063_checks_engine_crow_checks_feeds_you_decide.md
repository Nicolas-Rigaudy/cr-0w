## item_063_checks_engine_crow_checks_feeds_you_decide - Checks engine: crow-checks rules feed You decide
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 75%
> Progress: 100%
> Complexity: Medium
> Theme: Management
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-06 18:12:20

# AI Context
- Summary: bin/crow-checks runs small rule functions over calendar, mail, Jira and the vault and prints JSON rows that the deck shows in Next > You decide and the brief lists.
- Keywords: checks, rules, sprint, rituals, one-on-one, mail, devex, you decide, brief
- Use when: adding or changing a management check.
- Skip when: deck rendering of rows (064).

# Problem
- Management gaps go unnoticed until the owner sees them: sprint ending with no planning booked, 1:1s past due, team mail waiting.

# Scope
- In:
  - `bin/crow-checks json [-f]`, cached like crow-today. Each rule returns rows `{id, check, what, why, action, due}`.
  - Rules v1:
    - Sprint ending within 2 workdays with no planning, review or retro in the calendar. Sprint dates from Jira (crow-jira token), events from Graph.
    - 1:1 due per dev: vault person note due date vs calendar; a cancelled 1:1 counts as missing.
    - Team or management mail unanswered for more than 2 days (reuse the crow-today important rule).
    - Monday management weekly not prepped by Monday noon; daily not posted by 11:30.
    - DevEx check-in overdue (vault management note).
  - State `local/state/checks.json`: dismissed ids and snoozed-until dates.
  - Deck source srcChecks feeding the You decide group; `brief` skill lists the same rows.
  - Rules and thresholds in `local/checks.yaml`; generic engine in `bin/`.
- Out:
  - Acting on the rows (064, 065).
  - Writing to Jira.

# Acceptance criteria
- AC1: Replaying 2026-11-03 data (sprint ends, no planning booked) yields one sprint row with the suggested action "propose planning".
- AC2: A dev note with a past 1:1 due date and no 1:1 in the next 7 days yields a 1:1 row.
- AC3: Dismissed rows stay hidden; snoozed rows come back on the date.
- AC4: Each rule has a unit test with fixture input; the deck shows the rows in You decide.

# AC Traceability
- AC1 -> `tests/test_crow_checks.py::test_replay_2026_11_03_offline` (fixture `tests/fixtures/checks/2026-11-03`, offline): one sprint row, action "propose planning".
- AC2 -> `test_one_on_one_past_due_and_unbooked_yields_a_row`.
- AC3 -> `test_state_and_row_limit`, `test_dismiss_and_snooze_cli`.
- AC4 -> one unit test per rule in `tests/test_crow_checks.py`; deck `tools/crow-deck/checks_test.go::TestHomeShowsChecksInYouDecide`.

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
- Rationale: Owner asked 2026-11-03; the sprint planning gap is the founding example.

# Notes
- 2026-11-06 (Jarvis review): build the checks engine as a core job from the start, so the deck You decide pane and the nest Work screen read one engine.
- Vault people notes carry a monthly 1:1 recurring task with a due date; read it with crow_vault.py.
- Sprints last 2 weeks and end on Wednesday.
- Keep each rule a pure function over a snapshot so tests need no network.
- 2026-11-07 build: no crow-jira tool exists, so sprint dates come from a cache file (`sprint.file`) with a cadence fallback (`anchor_end`, `length_days`). A writer for that file from the tracker is future work.
- 2026-11-07 build: crow-today now caches a `week` of events for the 1:1 window. The live Graph call is covered by a mocked test only.
- 2026-11-07 build: the owner copies `docs/checks.example.yaml` to `local/checks.yaml` to turn the checks on.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
