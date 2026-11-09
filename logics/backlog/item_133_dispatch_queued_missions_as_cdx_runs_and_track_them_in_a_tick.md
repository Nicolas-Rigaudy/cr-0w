## item_133_dispatch_queued_missions_as_cdx_runs_and_track_them_in_a_tick - Dispatch queued missions as cdx runs and track them in a tick
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85
> Confidence: 80
> Progress: 0%
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-09 17:49:40

# AI Context
- Summary: `crow-core tick` admits queued missions, renders the brief, starts a detached cdx run through the existing dispatch path, and closes finished runs into mission state.
- Keywords: tick, dispatch, cdx run, detach, admit, run_start, run_end, exit cause, max parallel
- Use when: Turning a queued mission into a running one, or finishing a run.
- Skip when: Parking on questions (item_109) or resuming (item_135).

# Problem
- Slice 1 AC: a dispatched mission shows its state and runs in `crow mission list --json`. Today runs start from `crow-runs dispatch` with `RUN_KINDS`; missions must reuse that path, not fork it.

# Scope
- In:
  - `crow-core tick` (one shot; the service loop of slice 3 calls the same code every 2 min). No model call in the tick.
  - Admit: queued missions up to the background cap (2, CLAUDE.md autonomy table), oldest first.
  - Start: render `brief.md` (item_107), then start the run through the existing Python dispatch path (`bin/crow-runs dispatch`, extended with `--mission <id>` and kind flags from the kind file) so the item_095 account pick, item_077 allowlist and item_100 isolation stay in one place. Record `run_start` with run id and account.
  - Finish: for each running mission, read cdx run state; on exit record `run_end` with exit cause (done, usage limit, budget 125, timeout, error). A version file with a valid crow block moves the mission to `review`; no file moves it to `failed` with an incident. Usage limit and budget stops follow item_115 (park, resume later).
  - `kinds` replace `RUN_KINDS` for mission runs only; plain `crow dispatch` stays as is.
- Out:
  - Outbox and questions (item_109).
  - Watchdog, reconciler, summaries (later slices).
  - Timers and systemd (slice 3).

# Acceptance criteria
- AC1: A research mission created with `crow mission new`, then `crow-core tick`, shows `running` with one run id in `crow mission list --json`; after the run ends and another tick, it shows `review` (or `done` for the straight-through gate shape) with `run_end`.
- AC2: A third queued mission stays `queued` while two run.
- AC3: A run that exits without a version file leaves the mission `failed` with an `incident` event naming the exit cause.
- AC4: The run goes through `crow-runs dispatch` (one account pick and allowlist path), checked by a test with a fake cdx.

# AC Traceability
- request-AC1 -> AC1. Proof: tick test with a fake cdx and the real dispatch path.

# Decision framing
- Product framing: Not needed
- Architecture framing: Covered by adr_001 "One process" and "Resume".

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Primary task(s): `task_016_orchestrate_jarvis_core_slices_1_and_2_missions_and_inbox`

# Priority
- Priority: High
- Rationale: Closes slice 1.

# Notes
- Groomed 2026-11-09 from req_028 slice 1. Waits on item_132 and item_107. Uses item_106 (crow block parser, Done) to read the version file.
- From item_132: `crow mission new` accepts kinds of any phase. Admission must refuse to dispatch a kind whose phase is not live yet (build is phase 3).
- From item_107: call `mission.Env.Brief` then `brief.Write` before each start and resume. Pass preflight answers in `Input.Preflight` once the preflight item lands; nil renders the safest defaults.
- 2026-11-09 tool pass (owner picked): cdx 0.26.8 `cdx run --extra-args 'ARGS'` is per run. Dispatch passes each mission's allowlist per run, no per-account extra-args edits (memory: background runs deny by default).
- 2026-11-09: cdx 0.26.0 per-profile `launch.max_parallel` (0-5, 0 = no automatic missions; `cdx set <name> --max-parallel N`, shown in list/status JSON). The tick admits runs against it instead of crow's own cap.
- 2026-11-09: logics-manager 2.24.2 `flow record-run <ref> --runner cdx --run-id <id> --state <state>` writes each run into the Executions table of the mission's Logics doc. Call it on run start and end when a mission has a Logics ref.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
