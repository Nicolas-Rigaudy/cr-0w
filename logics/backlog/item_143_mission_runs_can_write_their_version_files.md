## item_143_mission_runs_can_write_their_version_files - Mission runs can write their version files
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Read-only mission kinds run in the default permission mode with a per-run allowlist that grants Edit/Write only inside the mission folder, and the brief names that folder.
- Keywords: crow-core, crow-runs, mission, dispatch, plan mode, sandbox allowlist, write-dir, brief, outbox
- Use when: changing how tick dispatches mission runs, the run sandbox rules, or the brief's file instructions.
- Skip when: working on the inbox, answers or the resume cap.

# Problem
A mission run writes its version file (plan, report, mockup) and outbox/questions.json in its own mission folder, whatever the kind's permission.
Writes stay limited to that folder for read-only kinds; the project stays read-only.

# Design
- crow-runs gains `--write-dir DIR` (absolute, must exist). The run then gets a per-run settings file (under `<state>/run-sandbox/<run>.json`, not the shared `<kind>.json`) whose allow list is the kind's plus `Edit(/DIR/**)` and `Write(/DIR/**)` (Claude absolute-path rule syntax), and `--add-dir DIR`. The RUN_DENY list stays as is.
- A review-kind run with `--write-dir` runs in permission `default` instead of `review` (plan mode). Headless default mode refuses any tool outside the allowlist, so the project stays read-only: the review allowlist has no Edit/Write and no writing Bash.
- `dispatchFor` (tools/crow-core/internal/mission/tick.go) keeps kind `review` for read-only and review permissions, drops `Permission: "review"`, and passes `WriteDir` = the mission folder. Task kinds also get `WriteDir` so they can reach the folder from the project cwd.
- brief.go prints the absolute mission folder and the exact paths to write: `<dir>/<plan|report|mockup>.vN.md`, `<dir>/outbox/`, `<dir>/outbox/questions.json`. The "Not allowed" line names the folder.
- Isolated runs (crow-runs forces review for the read-only mount, line ~816) are out of scope: missions do not dispatch isolated today.

# Scope
- In:
  - crow-runs `--write-dir` + per-run settings file + `--add-dir`; crow_lib helper for the extra rules
  - Dispatch struct field + runner flag in crow-core; dispatchFor change
  - brief.go folder and paths; golden/fixture updates
  - Tests: crow_lib unit test for the rules, crow-runs dry-run test, Go tests for dispatchFor and the brief
  - Real e2e: rerun the item_135 check (scratch config pattern in item_135 notes) on rk1
- Out:
  - Isolated (bwrap) mission runs: item_100
  - The mission guard hook: item_102

# Acceptance criteria
- AC1: A read-only mission run writes report.v1.md and outbox/questions.json in its mission folder and nothing in the project.
- AC2: The brief names the absolute mission folder and the files the run must write.
- AC3: The item_135 end-to-end check passes on a real account: question parked, answer given, tick resumes the run, report.v2.md quotes the answer.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: A read-only mission run writes report.v1.md and outbox/questions.json in its mission folder and nothing in the project.
- request-AC2 -> This backlog slice. Proof: AC2: The brief names the absolute mission folder and the files the run must write.
- request-AC3 -> This backlog slice. Proof: AC3: The item_135 end-to-end check passes on a real account: question parked, answer given, tick resumes the run, report.v2.md quotes the answer.

# Decision framing
- Product framing: Not needed
- Product signals: (none detected)
- Product follow-up: No product brief follow-up is expected based on current signals.
- Architecture framing: Not needed
- Architecture signals: (none detected)
- Architecture follow-up: No architecture decision follow-up is expected based on current signals.

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `logics/request/req_046_mission_runs_can_write_their_version_files.md`
- Primary task(s): `logics/tasks/task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox.md` (step 10 prerequisite)

# Priority
- Priority: High
- Rationale: Blocks task_016 step 10 (item_135 e2e) and every read-only mission kind.

# Notes
- Risk: confirm on cdx that `--extra-args` carries `--add-dir` and the per-run `--settings` path within CDX_EXTRA_ARGS_MAX (512).
- Risk: confirm headless Claude refuses (not prompts) an Edit outside the allowed path in default mode; test it in the e2e.
- Hybrid rationale: Derived from request `req_046_mission_runs_can_write_their_version_files` and kept bounded to one coherent delivery slice.
- Source file: `logics/request/req_046_mission_runs_can_write_their_version_files.md`.
- Generated locally by logics-manager.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

