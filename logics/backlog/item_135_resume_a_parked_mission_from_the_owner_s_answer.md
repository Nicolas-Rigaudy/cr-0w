## item_135_resume_a_parked_mission_from_the_owner_s_answer - Resume a parked mission from the owner's answer
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90
> Confidence: 85
> Progress: 100%
> Complexity: Medium
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-09 23:08:29

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
- AC2: The same flow on an isolated mission resumes cold and still sees the answer. Moved to item_160: core does not dispatch isolated missions yet.
- AC3: Hitting the resume cap creates an inbox item instead of a run.

# AC Traceability
- request-AC1 -> AC1. Proof: mission `20261110-review-e2e!5` (runs d552e2b8, 066c3af3). AC2 carried by item_160. AC3: resumecap unit tests.

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
- 2026-11-09 built in c440791 (ADR in fad796d): tick resumes answered missions through `env.resume`, cold after a usage limit. Resume cap: `budget.resumes` resumes in a row without a gate park the mission `queued -> waiting` on a core `resume-cap` question (continue/stop). Continue allows one more resume, stop cancels. Unit tests cover resume, cold resume, cap ask (AC3), cap stop and the gate reset.
- 2026-11-09 real end-to-end (AC1) blocked: a research mission on rk1 (Sonnet low) ran in plan mode and wrote no questions file or version file; tick failed it "no plan, report or mockup v1". Cause is in dispatch, not resume: read-only kinds run in plan mode and the brief never names the mission folder. Captured as req_046. AC1 and AC2 wait on it.
- 2026-11-09 item_143 (2bcaf72) unblocked the first run: the e2e mission parked its question and the answer queued it. The continue run then failed in cdx (inlined settings over the 512-char cap, AlexAgo83/cdx-manager#131). AC1 and AC2 wait on that fix.
- 2026-11-10 real e2e passed on rk1 (scratch store, review kind, Sonnet low, mission `20261110-review-e2e!5`): run d552e2b8 ($0.14) asked q1 in a needs-input report block, `crow inbox answer ... msg "q1: red"` queued it, the tick continued it as run 066c3af3 (cdx#131 fixed). `report.v2.md` is ready, quotes the answer as `a1`, and records the probe edit outside the mission folder as refused; the repo stayed clean. Three fixes on the way: the brief shows the crow block shape and the resume note names the `a1..aN` refs (7c48f4a); a needs-input or blocked gate gets an inbox row without validating verbs (af00367).
- Closed 2026-11-10: AC1 live, AC3 by unit test, AC2 moved to item_160 (isolated missions are not dispatched by core; cold resume is unit tested in TestTickResumesColdAfterUsageLimit).

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
