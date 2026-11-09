## item_135_resume_a_parked_mission_from_the_owner_s_answer - Resume a parked mission from the owner's answer
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85
> Confidence: 75
> Progress: 0%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: The tick resumes an answered mission: rebuilt brief, then `cdx run --continue-run` for non-isolated runs or a cold run from files for isolated ones (item_116), with a resume cap.
- Keywords: resume, park, continue-run, cold run, answers, resume cap
- Use when: A mission was answered and must run again.
- Skip when: Usage-limit timing (item_115 policy, already decided).

# Problem
- Slice 2 AC: a test mission asks a choice, exits, and resumes with the answer in its prompt. No session is held while waiting (req_028 decision "Park and resume").

# Scope
- In:
  - In the tick, an answered mission (`queued` with a prior run) rebuilds `brief.md` (item_107) with every answers file and the next version number.
  - Non-isolated last run on the same account: `cdx run --continue-run <run>` with a short prompt pointing at the new answers file.
  - Isolated last run, or an account change: cold run from the brief and version files (item_116 path; the item_115 policy picks wait or move).
  - Resume cap from the kind budget: after N resumes in a row without reaching a gate the mission asks the owner (inbox FYI) instead of running again.
  - `run_start` records `resume_of` and the mode (continue or cold).
- Out:
  - Native isolated resume in cdx (item_116 tracks the upstream issue).
  - Moving a mission to the tower (item_120).

# Acceptance criteria
- AC1: End to end with a real cheap run (research or plan kind, Sonnet low): the mission writes `questions.json` with a choice, parks, `crow inbox answer` gives the choice, the next tick resumes it and the new version file quotes the answer in its Responses.
- AC2: The same flow on an isolated mission resumes cold and still sees the answer.
- AC3: Hitting the resume cap creates an inbox item instead of a run.

# AC Traceability
- request-AC1 -> AC1, AC2. Proof: end-to-end run log and the mission dir.

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001 "Resume" and req_037 D2.

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: Closes slice 2 and the 0.3b exit signal for the CLI.

# Notes
- Groomed 2026-11-09 from req_028 slice 2. Waits on item_134, item_133 and item_116 (cold path).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
