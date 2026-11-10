## item_160_core_runs_workspace_write_missions_isolated - Core runs workspace-write missions isolated and resumes them cold
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 80
> Confidence: 70
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow-core dispatches every mission unisolated; workspace-write kinds (build, design, plan, repair, review-fix) should run isolated in a worktree, write their version files, and resume cold from files.
- Keywords: crow-core, isolate, worktree, cold resume, write-dir, mission folder
- Use when: wiring isolated dispatch into the mission tick or proving the cold resume path live.
- Skip when: read-only missions (review, research); they already run and resume (item_143, item_135).

# Problem
- `dispatchFor` (tools/crow-core/internal/mission/tick.go) never asks crow-runs for `--isolate`, so a workspace-write mission edits the live checkout.
- crow-runs refuses `--write-dir` on isolated runs (item_143), so an isolated run has no way to write its version file and outbox in the mission folder.
- Isolated runs cannot continue a conversation; the cold resume path (item_116) is only unit tested. item_135 AC2 moved here.

# Scope
- In:
  - Core dispatches workspace-write kinds with `--isolate` on a worktree under `<repo>.wt/`.
  - A way for the isolated run to deliver its version file and outbox to the mission folder (bind the folder into the sandbox, or the host copies them from the worktree at run end).
  - Resume of an isolated mission is cold from brief.md plus the resume note.
  - Real e2e on a cheap kind profile (runbooks/mission-e2e-scratch.md).
- Out:
  - Merge gate for mission branches: item_101.
  - Guard hook for runs that cannot be isolated: item_102.

# Acceptance criteria
- AC1: A workspace-write mission runs isolated in a worktree and leaves the live checkout untouched.
- AC2: The isolated run's version file and outbox reach the mission folder and the tick ingests them.
- AC3: Answering a parked isolated mission resumes it cold and the next version file quotes the answer (former item_135 AC2).

# AC Traceability
- AC1 -> e2e mission on a scratch store. Proof: worktree diff and clean checkout.
- AC2 -> same mission. Proof: version file and outbox in the mission folder.
- AC3 -> same mission after an answer. Proof: run_start with `cold: true` and the quoted answer.

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001 (missions) and item_100 (isolation); the delivery path in AC2 may need a short decision.

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): (none yet)
- Related: item_100, item_116, item_135, item_143

# Priority
- Priority: Medium
- Rationale: Needed before any build-kind mission runs unattended; read-only missions work without it.

# Notes
- Split from item_135 on 2026-11-10 after its read-only e2e passed.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

