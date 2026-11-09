# Mission end-to-end check in a scratch store

Run a real, cheap mission through crow-core (dispatch, park, answer, resume) without touching the live mission store.

## Setup
1. Pick a scratch folder `S` (the session scratchpad, never the repo).
2. Write `$S/core.yaml` pointing missions, store and kind overrides into `$S` (`bin/crow core config` prints the keys to override).
3. Write `$S/kinds/research.yaml` with a cheap profile: Sonnet or Haiku, effort low, 15 min, budget $1.
4. Export `CROW_CORE_CONFIG=$S/core.yaml` and `CROW_ACCOUNT=<an allowed claude account>` in every command.

## Mission goal
Ask for three things so one run covers each path:
- a blocking question with two choices (written to `outbox/questions.json`, ids `q1`, `q2`...);
- after the answer, `report.v2.md` quoting the answer;
- one try at creating a probe file in the cwd, reporting whether it was refused.

## Loop
1. `bin/crow core mission new --kind research --project cr-0w --goal "..." --end-state "..."`
2. `bin/crow core tick`: starts run 1. Check the brief names the absolute mission folder.
3. Wait for the run (`bin/crow-runs list`), then `bin/crow core tick`: the mission parks on the question.
4. `bin/crow inbox answer <mission>/v1/q1 <choice>`, then `bin/crow core tick`: starts the resume.
5. Tick again after it ends: expect `report.v2.md` quoting the answer, the probe refused, `git status` clean.

## When a run fails
- Read the cdx launch log: `<cdx profile home>/log/cdx-run-<run_id>-attempt-1.launch.log`.
- Compare the run's recorded `launch.extra_args` with what crow-runs passed (`bin/crow-runs` ledger in `local/state/runs.jsonl`).
- A failed mission stays failed; start a new mission for the rerun.

## Known blockers
- cdx 0.26.11 inlines `--settings` JSON on `--continue-run` and then rejects it past 512 chars (upstream issue #131 on cdx-manager). Resumes of sandboxed Claude runs fail until it ships.
