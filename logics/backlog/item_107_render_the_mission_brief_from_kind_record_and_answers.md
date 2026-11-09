## item_107_render_the_mission_brief_from_kind_record_and_answers - Render the mission brief from kind, record and answers
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:45:03

# AI Context
- Summary: core writes `brief.md` from the kind file, the mission record, the owner prefs and every answers file, following the spec template.
- Keywords: brief, template, owner prefs, reply contract, resume
- Use when: Starting or resuming a mission run.
- Skip when: Writing the kind files themselves.

# Problem
- A run must restart cold on any account from files alone; the brief is what it reads first.
- Owner prefs and prior answers must reach every run without hand copying.

# Scope
- In:
  - Render `brief.template.md` sections in order: Goal, Mission, Owner preferences, Preflight, Answers so far, Constraints, Reply contract.
  - Rebuild the brief before every resume with the latest answers and the next version number.
  - Take owner prefs from the prefs file and filter them by kind.
- Out:
  - Preflight UI (preflight item).
  - Account choice (item_095).

# Acceptance criteria
- AC1: The rendered brief for the spec example mission matches the example brief in sections and facts.
- AC2: After answers v1 and v2 the brief lists both with the owner's words and asks for version 3.
- AC3: The brief names no account and no private path.

# AC Traceability
- request-AC3 -> AC1. Proof: rendered brief compared to the spec example.
- AC1 -> `internal/brief` TestFirstBriefIsTheExample: the render of the example mission equals `examples/20261109-build-item-095/brief.md` byte for byte, sections in template order.
- AC2 -> `internal/brief` TestResumeBriefCarriesTheAnswers: with answers v1 and v2 the brief lists both verbs and the owner's words, asks for `report.v3.md` and names `answers.v3.md`.
- AC3 -> `internal/brief` TestBriefNamesNoAccountOrPath: a record with an account and a worktree path renders neither, and no absolute path.

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
- Rationale: Needed for the first mission run, after the parser.

# Notes
- Depends on item_104, item_106 and the core mission record (req_028, adr_001).
- Code home waits on the Go core and the shared fixture contract (item_119). The backend is Go (owner decision 2026-11-07); D6 keeps only the parser in Python `bin/`.
- Spec: `logics/specs/req_037_run_contract/`.
- 2026-11-09 groomed into task_016. Code home: Go in `tools/crow-core` (D6 resolved by req_031). Mission record comes from item_132.

- 2026-11-09 built in `internal/brief` (5113cc4). `crow core mission brief <id> [--write]` prints or writes `brief.md`; the write is atomic.
- Owner prefs come from `prefs` in core config (default `local/mission-prefs.yaml`), a YAML list of `rule` plus optional `kinds`. `prefs.example.yaml` in the spec shows the shape.
- The key comes from the created event, since the id lowercases it. Repo and worktree stay placeholders (`<repo>`) so the brief names no local path.
- The example brief.md was regenerated from the renderer: Preflight section, data classes from the kind, full reply contract.
- item_133 renders the brief at admission and before each resume.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
