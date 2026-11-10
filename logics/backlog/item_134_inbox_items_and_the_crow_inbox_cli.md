## item_134_inbox_items_and_the_crow_inbox_cli - Inbox items and the crow inbox CLI
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85
> Confidence: 80
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-10 11:19:08

# AI Context
- Summary: One inbox built from mission events (version files at a gate and `questions.json`), with `crow inbox list|show|answer` taking the gate verbs of the run contract.
- Keywords: inbox, choice, approval, review, fyi, verbs, answer, crow inbox
- Use when: Showing or answering what missions ask.
- Skip when: Rendering the inbox in deck or nest (slices 6-7).

# Problem
- req_028 need 3: every question, approval and proposal lands in one queue that every surface reads. Slice 2 needs it from the CLI first.

# Scope
- In:
  - Inbox rows built in the store from events: a version file at a gate (`review` state) and each parked question (item_109). Item types: choice, approval, review, FYI (req_028 need 3), mapped from the gate shape and the block.
  - `crow inbox list [--json]` oldest first with mission, kind, type, one-line ask; `show <item>` prints the block summary, questions, options and artifact paths.
  - `crow inbox answer <item> <verb> [words]` with the verbs of req_037 "Gate shapes and verbs" (go, dig, no, ok, msg, abandon, pick, change, approve) plus a plain answer for `questions.json`. Verbs not allowed for the gate are refused with the allowed list. An open blocking question refuses go, ok and approve.
  - The answer goes through item_108 (answers file, validation record) and writes an `answer` event; the item closes and the mission moves to `queued` for resume (item_135) or to its end state (`no`, `abandon`, `approve`, `ok` with nothing to publish).
  - Surface recorded on the answer (`cli` now; deck and nest later).
- Out:
  - Publishing after `ok` (item_099).
  - deck pane and nest screens.
  - Web Push.

# Acceptance criteria
- AC1: A parked fixture mission shows one inbox item; `answer` with a valid verb closes it and writes `answers.vN.md` plus an `answer` event.
- AC2: A wrong verb for the gate and a `go` with an open blocking question are both refused with the reason.
- AC3: `crow inbox list --json` matches the item_119 inbox fixture shape.

# AC Traceability
- request-AC1 -> AC1. Proof: `TestInboxAnswersAParkedQuestion` (tools/crow-core/internal/mission/inbox_test.go): one row, answer writes `answers.v1.md` and an `answer` event, the row closes, the mission queues.
- request-AC3 -> AC2. Proof: `TestInboxRefusesWrongVerbAndOpenBlockingQuestion` and `TestInboxQuestionBlocksGoOnTheSameMission`: wrong verb refused with the allowed list, `go` refused while blocking q1 is open.
- AC3. Proof: `TestInboxMatchesTheFixtureShape`: rows from the research fixture log have the keys of `fixtures/records/inbox.json` and the same id, type, source, verbs, state and timestamps.

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001 and the req_037 spec.

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: Slice 2 of the 0.3b core.

# Notes
- Groomed 2026-11-09 from req_028 slice 2. Waits on item_109 and item_108.
- 2026-11-09 done. `internal/mission/inbox.go` fills the ask, blocking flag and artifacts from the files on disk; `inbox_cmd.go` and `bin/crow` route `crow inbox`. Answers go through `env.reply` with the row id and surface on the event.
- New legal moves: waiting -> queued and review -> queued. The tick holds a queued mission that already ran until item_135 lands the resume.
- A question row needs the verb `answer` and words. A validating verb (go, ok, approve) is refused while a blocking question row of the same mission is open.
- Task `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox` was finished via `logics-manager flow finish task` on 2026-11-10.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Tasks
- `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`
