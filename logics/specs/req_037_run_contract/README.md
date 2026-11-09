# Run contract spec (req_037)

Draft spec for mission kinds and the contract between core and a run. Data files only. Nothing here is wired into code yet. The backlog items of req_037 build it.

## Files

| Path | What |
|---|---|
| `schemas/kind.schema.json` | Shape of a kind file |
| `kinds/<kind>.yaml` | One file per decided kind (8) |
| `schemas/crow-block.schema.json` | The fenced `crow` block in plan, report and mockup files |
| `schemas/questions.schema.json` | `outbox/questions.json` |
| `schemas/validation-event.schema.json` | The validation line in the event log |
| `brief.template.md` | What core writes as `brief.md` |
| `examples/20261109-build-item-095/` | A full two-gate mission: brief, plan v1 and v2, answers, report v3, events |
| `examples/gates/` | One example per other gate shape, plus a `questions.json` |

All schemas use JSON Schema draft 2020-12. Kind files and `crow` blocks are YAML that parses to JSON.

## Kinds

| Kind | Phase | Gates | Autonomy | Skill | Key | Publish |
|---|---|---|---|---|---|---|
| research | 1 | `[]` | auto | missing | slug | none |
| plan | 1 | `[plan]` | ask | `corpus` (bundled) | slug or request id | commit-docs |
| review | 1 | `[report]` | notify | `review` peer mode | MR ref | post-comments |
| design | 1 | `[mockup]` | ask | missing (tool `bin/crow-compare`) | slug | commit-docs, next-mission |
| review-fix | 2 | `[report]` | auto | `review-threads` | MR ref | push, reply-resolve |
| write | 2 | `[report]` | notify | missing (fallbacks: local daily, weekly, mgmt-weekly) | schedule | page, mail-draft, vault-note |
| build | 3 | `[plan, report]` | ask | `implement-task`, then `review`, `mr` | ticket, item or slug | push, mr |
| repair | 3 | `[report]` | notify | missing | MR ref or branch | push |

Scope overrides (`by_scope`):

- review-fix, personal: no gate, auto.
- build, personal: no gate, auto merge, ask on sensitive paths.
- build, crow: `[plan, report]`, merge after go on the report, ask on sensitive paths.
- repair, personal and crow: no gate, auto, ask on sensitive paths.

Kind files never name an account. `account.route: scope` hands the pick to the scope-aware account pick (item_095). Models and caps in the files are draft defaults.

## Gate shapes and verbs

Gates are an ordered list. Five shapes are in use.

| Shape | Kinds | File at the gate | Verbs |
|---|---|---|---|
| `[]` straight through | research | `report.vN.md`, no wait | none (msg reopens it) |
| `[report]` | review, review-fix, write, repair | `report.vN.md` | ok, msg, abandon |
| `[plan]` | plan | `plan.vN.md` | go, dig, no |
| `[plan, report]` | build | `plan.vN.md`, then `report.vN.md` | go, dig, no; then ok, msg, abandon |
| `[mockup]` | design | `mockup.vN.md` | pick, change, approve, abandon |

- `go` runs the next step. `dig` asks for a new plan version with the owner's words. `no` closes the mission.
- `ok` validates the report and lets core publish. `msg` sends words back for a new version.
- `pick <option>` narrows the mockup. `change <words>` asks for a new round. `approve` ends the loop.
- A blocking question open in the block refuses `go`, `ok` and `approve` until it is answered.

## Mission dir

```
missions/<YYYYMMDD-kind-key>/
  brief.md            written by core, rewritten before each resume
  plan.vN.md          one crow block each
  report.vN.md
  mockup.vN.md
  answers.vN.md       the owner's answer to version N
  outbox/             artifacts and questions.json, capped
  events.jsonl        {ts, mission, type, actor, data}
```

The event types and the data of each are fixed by `logics/specs/core_contract` (`schemas/event.schema.json`).

One version counter per mission. Plan, report and mockup share it. `answers.vN.md` always answers the file of version N. In the example, plan v1, plan v2 and report v3 follow each other.

## Versions and responses

- core never rewrites a version file. Each answered version stays with its answers file.
- `answers.vN.md` has a front matter (mission, version, verb, surface, ts) and the owner's words as a numbered list. Points are referred to as `a1`, `a2`.
- From version 2 the file starts with a `## Responses` section. The block's `responses` maps each answer point or question id to `done`, `partly`, `declined` or `moot`, with a note. The schema requires it from version 2.
- A version may stay open for further words. A new message before `go` or `ok` is appended to the same answers file.

## Validation record

Same record as req_030 Need 5 and item_099: `{mission, version, action, target, payload_sha256, owner_event, ts}`, chained by hash. Schema: `validation-event.schema.json`.

- `ok` (and `go` or `approve` where they publish) writes one `validation` event per entry of the block's `publish` list.
- `payload_sha256` is the SHA-256 of the canonical JSON of that entry: keys sorted, separators `,` and `:`, UTF-8, no ASCII escaping.
- `owner_event` is the ts of the owner's answer event. `prev_sha256` is the SHA-256 of the previous validation line of the mission as stored, 64 zeros for the first. A replayed or edited record breaks the chain.
- Before each action core recomputes the hash from the file on disk. A mismatch refuses that action and parks the mission.
- Any owner message after the validation writes `validation_dropped`. core publishes nothing until a new `ok`.
- Only the owner validates. A chat surface, the butler or a run never writes a validation.
- Each publish writes a `publish` event that names the validation it used. The example mission has two (push, then MR).

## Outbox

- Everything the owner should see goes under `outbox/` and is listed in `artifacts`.
- Types: html, image, video, diff, doc, data.
- Caps per kind (`caps.outbox_file_mb`, `caps.outbox_total_mb`). Most kinds: 20 MB per file, 100 MB in total. design: 200 MB in total for video. write: 20 MB and 50 MB. review, review-fix and repair: 5 MB and 20 MB.
- core ingests the outbox after the turn, in the order of `artifacts`. A file over the per-file cap, or over the total once earlier files filled it, is dropped and noted in the event log. A file not listed in `artifacts` is ignored. Dropped files stay on disk.
- A listed path must start with `outbox/`, hold no `..` and resolve to a regular file inside `outbox/` (links that lead out are refused).
- Each dropped file gives one `incident`: `artifact_over_cap`, `artifact_refused` or `artifact_missing`. The mission still reaches its gate.
- `outbox/questions.json` is for questions outside a plan or report. Writing it ends the turn. core parks the mission (waiting), logs a `question` event and the fold shows the questions in the inbox. Choice previews go through the same caps.
- A questions file whose `version` is below the run's version was already asked and is ignored. An invalid file gives a `bad_questions` incident and the turn ends on its version file.

## brief.md

Sections, in order: Goal, Mission, Owner preferences, Preflight, Answers so far, Constraints, Reply contract. See `brief.template.md` and the example brief. core rebuilds it before every resume. A cold run on another account reads the brief and the version files and needs nothing else.

Owner preferences come from `local/mission-prefs.yaml` (config key `prefs`): a list of one-line rules, each with optional `kinds`. See `prefs.example.yaml`. The brief names no account and no local path; the repo and worktree stay placeholders. `crow core mission brief <id> [--write]` prints or writes it.

## End of a turn

- The run writes one version file and ends the turn. It never waits for an answer inside the run.
- A usage limit, a budget stop (`cdx --budget` exits 125) or a timeout ends the run. core records it and resumes later.
- A resume after the owner's answer uses the same session where the tool allows it. A resume on another account is a cold run from the files. Owner decision D2 (2026-11-07): at night wait for the reset; by day restart cold when the reset is more than about 1 hour away, otherwise wait.
- A stop guard hook can refuse to stop when no valid version file was written. Claude Code Stop hooks can block with exit code 2 and see the last assistant message. `stop_hook_active` is true on a stop that follows a block; the hook allows the stop then, so a run cannot loop. The hook needs hooks to run in headless mode (sandbox item).

## Crow today

What exists, read from the code on 2026-11-06:

- `bin/crow-runs dispatch --kind task|review|retro`.
- `crow_lib.RUN_KINDS`: task (opus, medium, 25 USD, 120 min), review (opus, high, 10 USD, 60 min), retro (sonnet, low, 2 USD, 45 min), oneshot (opus, medium, 3 USD, 20 min).
- `run_flags` passes model, effort, timeout, failover and budget to cdx. Runs are logged in `<state>/runs.jsonl`.

The kind files replace `RUN_KINDS` once a loader exists.

## Sources

| Claim | Source | Checked |
|---|---|---|
| `--continue-run` pins the parent's profile, model and effort; not combinable with `--failover` or `--isolate-workspace` | cdx-manager README, https://github.com/AlexAgo83/cdx-manager | v0.26.6, 2026-11-06 |
| `--budget` stops with exit 125 and `budget_exceeded` | same README | v0.26.6, 2026-11-06 |
| `cdx run` flags: `--kind`, `--permission`, `--timeout-seconds`, `--budget`, `--label`, `--detach`, `--failover` | `cdx run --help` | 0.26.6, 2026-11-06 |
| `--output-format json --json-schema` returns `structured_output`; `--bare` skips hooks, CLAUDE.md and MCP; `--resume <session>` | https://code.claude.com/docs/en/headless | Claude Code 2.1.319, 2026-11-06 |
| Stop hook blocks with exit 2 or `{"decision":"block"}`, receives `last_assistant_message` | https://code.claude.com/docs/en/hooks | 2026-11-06 |
| Stop input has `stop_hook_active`; consecutive stop-hook continuations are capped at 8 | https://code.claude.com/docs/en/hooks | Claude Code 2.1.319, 2026-11-07 |
| JSON Schema draft 2020-12 meta-schema | https://json-schema.org/draft/2020-12/schema | published 2022-06-16, checked 2026-11-06 |

## Open decisions

All six were settled by the owner on 2026-11-07. See req_037 under Decisions.
