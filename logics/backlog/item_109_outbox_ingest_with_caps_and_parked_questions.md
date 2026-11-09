## item_109_outbox_ingest_with_caps_and_parked_questions - Outbox ingest with caps and parked questions
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 85%
> Confidence: 75%
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 19:05:08
> Owner: crow

# AI Context
- Summary: After a turn core ingests listed outbox files within the kind caps and parks the mission on `outbox/questions.json`.
- Keywords: outbox, artifacts, caps, questions, park
- Use when: Ending a run turn.
- Skip when: Rendering artifacts in a UI.

# Problem
- Runs produce images, HTML, video and diffs; size must stay bounded.
- Questions outside a plan or report need a fixed place.

# Scope
- In:
  - Keep only files listed in `artifacts`; drop files over the per-file or per-mission cap and log it.
  - Refuse paths outside `outbox/`.
  - On a valid `questions.json`, park the mission and add the questions to the inbox.
- Out:
  - Thumbnails and previews (nest).
  - Artifact storage outside the mission dir.

# Acceptance criteria
- AC1: An over-cap file is dropped and logged; the mission still reaches its gate.
- AC2: An artifact path with `..` or outside `outbox/` is refused.
- AC3: A questions file parks the mission and its questions reach the inbox.

# AC Traceability
- request-AC4 -> AC3. Proof: Implemented in 74075f4; validated with go test -count=1 ./... (outbox_test.go) and bin/crow-check. TestTickParksOnQuestions copies the spec questions file and parks the mission. Source: `74075f4`
- AC1 -> `internal/mission` TestIngestCaps (per-file and total caps, missing file) and TestTickLogsDroppedArtifacts: a 21 MB listed video gives an `artifact_over_cap` incident and the mission still reaches its gate and `done`.
- AC2 -> `internal/mission` TestIngestRefusesPathsOutsideOutbox: `outbox/../x`, `outbox/sub/../../x`, an absolute path, `notoutbox/x`, a link leading out and a directory give `artifact_refused`.
- AC3 -> `internal/mission` TestTickParksOnQuestions: the spec example questions file moves the mission to `waiting` with a `question` event, and the inbox holds `<id>/v1/q1`.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: Medium
- Rationale: Needed once runs write real artifacts.

# Notes
- Depends on item_104, item_106 and the core mission record.
- Caps per kind in the kind files.
- Code home waits on the Go core and the shared fixture contract (item_119). The backend is Go (owner decision 2026-11-07); D6 keeps only the parser in Python `bin/`.
- Tool approvals without exit and resume: `--permission-prompt-tool mcp__<server>__approve` makes Claude Code call an MCP tool and wait. The tool writes a pending approval file, polls it until deck, nest or ntfy sets allow or deny, and denies on timeout. Mark pending files expired when stdin closes; refuse Allow when the shown input was clipped (harness-maxxing review 2026-11-07, local/state/research/harness-maxxing-and-hermes-2026-11-07.md).
- Outside text (mail, issues, web pages) reaches the run between markers, with the rule never to follow it. Owner messages are tasks, forwarded text is outside text.
- 2026-11-09 groomed into task_016. Code home: Go in `tools/crow-core`; runs in the tick after a run ends (item_133). Inbox rows land through item_134.
- 2026-11-09 built (74075f4): `internal/mission/outbox.go`. `Ingest(dir, paths, caps)` checks listed files in list order and returns kept paths and drops; it changes nothing on disk (dropped files stay, only logged). The tick logs one incident per drop (`artifact_over_cap`, `artifact_refused`, `artifact_missing`) after run_end and before the gate.
- A valid `outbox/questions.json` for the run's version ends the turn: run_end, preview drops, a `question` event (actor run), state `waiting` with reason "parked: questions". A file with a lower version is stale and ignored; an invalid file gives `bad_questions` and the turn ends on its version file.
- For item_134: reuse `Ingest` with the kind caps to fill inbox artifacts and previews.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
