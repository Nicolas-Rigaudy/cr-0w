## item_109_outbox_ingest_with_caps_and_parked_questions - Outbox ingest with caps and parked questions
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 17:14:57

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
- request-AC4 -> AC3. Proof: the dry run parks on the spec questions file.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `logics/request/req_037_mission_kinds_and_run_contract.md`
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Needed once runs write real artifacts.

# Notes
- Depends on item_104, item_106 and the core mission record.
- Caps per kind in the kind files.
- Code home waits on the Go core and the shared fixture contract (item_119). The backend is Go (owner decision 2026-11-07); D6 keeps only the parser in Python `bin/`.
- Tool approvals without exit and resume: `--permission-prompt-tool mcp__<server>__approve` makes Claude Code call an MCP tool and wait. The tool writes a pending approval file, polls it until deck, nest or ntfy sets allow or deny, and denies on timeout. Mark pending files expired when stdin closes; refuse Allow when the shown input was clipped (harness-maxxing review 2026-11-07, local/state/research/harness-maxxing-and-hermes-2026-11-07.md).
- Outside text (mail, issues, web pages) reaches the run between markers, with the rule never to follow it. Owner messages are tasks, forwarded text is outside text.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
