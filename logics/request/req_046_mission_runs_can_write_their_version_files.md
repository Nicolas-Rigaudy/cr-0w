## req_046_mission_runs_can_write_their_version_files - Mission runs can write their version files
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Core dispatches read-only mission kinds in Claude plan mode and never tells the run where its mission folder is, so a run cannot write outbox/questions.json or its version file.
- Keywords: crow-core, mission, dispatch, plan mode, permission, version file, outbox
- Use when: changing how tick dispatches mission runs, their permissions or the brief.
- Skip when: working on the inbox or the resume cap logic.

# Needs
- A mission run writes its version file (plan, report, mockup) and outbox/questions.json in its own mission folder, whatever the kind's permission.
- Writes stay limited to that folder for read-only kinds; the project stays read-only.

# Context
- Found in the item_135 end-to-end check on 2026-11-09: a research mission (permission read-only) ran on rk1, ended exit 0 and wrote nothing. The run said plan mode allowed only plan-file edits and ExitPlanMode was unavailable. Tick then failed the mission: "no plan, report or mockup v1".
- `dispatchFor` (tools/crow-core/internal/mission/tick.go) maps read-only and review kinds to crow-runs kind "review", which runs in plan mode.
- The run's cwd is the project dir and brief.md never names the absolute mission folder.
- crow_lib.py `sandbox_args(kind)` already carries per-kind tool allowlists through a `--settings` file; crow-runs forces review permission for isolated sandboxes. Build on these.
- Related: item_100 (isolated mission runs), item_102 (mission guard hook), item_133 (dispatch), item_135 (resume, whose real end-to-end check waits on this).

# Acceptance criteria
- AC1: A read-only mission run writes report.v1.md and outbox/questions.json in its mission folder and nothing in the project.
- AC2: The brief names the absolute mission folder and the files the run must write.
- AC3: The item_135 end-to-end check passes on a real account: question parked, answer given, tick resumes the run, report.v2.md quotes the answer.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/crow-core/internal/mission/tick.go`
- `bin/crow_lib.py`
- `bin/crow-runs`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

