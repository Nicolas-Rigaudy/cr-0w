# Mission end-to-end check in a scratch store

Run a real, cheap mission through crow-core (dispatch, park, answer, resume) without touching the live mission store. Cost: about $0.02 to $0.15 per run on Sonnet low.

## Setup
1. Pick a scratch folder `S` (the session scratchpad, never the repo).
2. Write `$S/core.yaml` pointing missions, store and kind overrides into `$S` (`bin/crow core config` prints the keys to override).
3. Copy the kind into `$S/kinds/<kind>.yaml` with a cheap profile: Sonnet, effort low, 15 min, budget $1, no preflight. The `review` kind is read only, so it also exercises the mission-folder write rights.
4. `export CROW_CORE_CONFIG=$S/core.yaml CROW_ACCOUNT=<an allowed claude account>` before every command. The dispatch keeps that account unless its 5h use is past `NEARLY_OUT_PCT` (crow_lib); then the scope pick moves the run to another rank-1 account and prints `crow: dispatching on ...`.

## Mission goal
Ask for three things so one run covers each path:
- a blocking question with two choices, in the crow block's `questions` (ids `q1`, `q2`...);
- after the answer, `report.v2.md` quoting the answer;
- one try at creating a probe file in the cwd, reporting whether it was refused.

## Loop
1. `bin/crow core mission new --kind review --project cr-0w --key 'e2e!N' --goal "..." --end-state "..."`. The review kind needs a key matching its MR pattern (`e2e!N` passes); use a new N per attempt.
2. `bin/crow core tick`: starts run 1. Check the brief names the absolute mission folder and shows the block shape.
3. Wait for the run (`bin/crow-runs list`), then `bin/crow core tick`.
   - Questions in `outbox/questions.json` park the mission on a question row.
   - Questions in a `needs-input` report block send it to review with a `needs-input` row (verbs `msg`/`abandon`, no validating verb).
4. `bin/crow inbox list`, then answer:
   - question row: `bin/crow inbox answer <mission>/v1/q1 <choice>`;
   - needs-input row: `bin/crow inbox answer <mission>/v1 msg "q1: <choice>"`.
   Then `bin/crow core tick` starts the `--continue-run` resume.
5. Tick again after it ends: expect `report.v2.md` with a `responses` entry `ref: a1` quoting the answer, the probe refused, `git status` clean, and a `review` row (`ok`/`msg`/`abandon`).
6. `bin/crow core rebuild` refolds the store if a row looks stale.

## When a run fails
- `json: cannot unmarshal string into CrowBlock...`: the run wrote a bare string where the contract wants an object. Check the block shape in the brief (`internal/brief/brief.go` `blockShape`).
- `vN responses miss aK`: the resume note names refs `a1..aN` per answer section (`internal/mission/resume.go`); check the run used them.
- Read the cdx launch log: `<cdx profile home>/log/cdx-run-<run_id>-attempt-1.launch.log`.
- Compare the run's recorded `launch.extra_args` with what crow-runs passed (`bin/crow-runs` ledger in `local/state/runs.jsonl`).
- A failed mission stays failed; start a new mission (new key) for the rerun.

## Not covered
- Isolated (worktree) missions: core does not dispatch them yet (item_160), so the cold resume path is only unit tested.
