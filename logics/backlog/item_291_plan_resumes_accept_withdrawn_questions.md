## item_291_plan_resumes_accept_withdrawn_questions - Plan resumes accept withdrawn questions
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Small
> Theme: nest
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: A plan resume that withdraws or rewrites earlier questions fails core validation, because every earlier question needs a response entry and the resume note only mentions answer points.
- Keywords: crow-core, mission, resume, responses, questions, validation, plan
- Use when: a mission fails with "vN responses miss qK".
- Skip when: display of the failure (item_290).

# Problem
- The desk-bugs plan mission v2 (2026-11-25 14:17) settled q2 and q3, rewrote q1 as q4, and answered a1. Core failed the mission: "v2 responses miss q1, q2, q3". The plan and its new question q4 never reached the inbox. The owner and the butler both believe it waits on the owner.
- Cause:
  - `tools/crow-core/internal/mission/answer.go:269-283` requires a `responses` entry for every question of the previous version unless an answer point covers it.
  - The resume note (`internal/mission/resume.go:137`) only tells the run to give each answer point `ref: aN`. It never says earlier questions need a response, and no documented outcome exists for a withdrawn or replaced question.
  - A validation miss on a run that exited 0 goes straight to `failed`, with resumes left in budget (2 allowed, 1 used).

# Scope
- In:
  - Resume note: list the previous version's question ids and say each needs a `responses` entry (`ref: qK`, outcome `done`, `withdrawn` or `replaced`, note pointing to the new id).
  - Contract: accept `withdrawn` and `replaced` outcomes; validation counts them as covered.
  - On a `responses miss` incident with resumes left, core resumes the run once with the validation error in the note instead of failing.
  - Runbook `runbooks/mission-e2e-scratch.md` troubleshooting line updated.
- Out:
  - Relaunching the live desk-bugs mission: the owner relaunches it (owner prefs: a failed run the owner launched is fixed, not relaunched unasked).

# Acceptance criteria
- AC1: A v2 that withdraws q2, replaces q1 with q4 and answers a1 passes validation with the new outcomes.
- AC2: The resume note for a version with questions names each question id.
- AC3: A v2 missing a question response triggers one automatic resume with the error, then fails if still missing.
- AC4: go test covers the three cases with the desk-bugs plan.v2.md as a fixture.

# Links
- Request: `req_069_nest_butler_chat_polish_and_failure_clarity`
- Unblocks: mission 20261125-plan-nest-desk-bugs (desk bugs 1 to 3, drafts in its outbox).
- Related: item_290.

# Launch prompt
```text
Read local/owner-prefs.md and local/prefs/design.md first. Work in a worktree under <repo>.wt/, commit with type[SCOPE]: message, no AI attribution, never push to main without the owner OK. Implement item_291 (plan resumes accept withdrawn questions) end to end. Read logics/backlog/item_291_plan_resumes_accept_withdrawn_questions.md first, then tools/crow-core/internal/mission/answer.go (Coverage) and resume.go. Use local/state/missions/20261125-plan-nest-desk-bugs/plan.v1.md, plan.v2.md and answers.v1.md as test fixtures (copy them, never edit the live mission). Add the withdrawn/replaced outcomes, name previous questions in the resume note, add the one automatic resume on a coverage miss, update the runbook. Do not relaunch the live mission; report that the owner can relaunch it.
```

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
