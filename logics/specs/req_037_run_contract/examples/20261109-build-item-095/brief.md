# Brief: 20261109-build-item-095

<!-- Written by core. The run reads it first and never edits it. -->

## Goal
Pick the account for a dispatch from the scope of the target repo (backlog item item_095).

## Mission
- kind: build
- gates: plan, report
- key: item_095
- scope: crow
- repo: `<repo>` on branch `mission/20261109-build-item-095` (worktree `<repo>.wt/20261109-build-item-095`)
- caps: 120 min, 25 USD, 3 resumes, 3 rounds
- current version: 1 (write `plan.v1.md`)

## Owner preferences
- Prose: short and declarative. No em-dashes.
- Commit subjects: `type[SCOPE]: message`.
- Never push or merge before the report gate passes.

## Preflight
Skipped: use the safest defaults.
- Base branch? the default branch.
- Which checks gate the report? the repo's full check command.

## Answers so far
None yet.

## Constraints
- Allowed: read and edit inside the worktree, run the repo checks, commit on the mission branch.
- Not allowed: push, merge, post, install packages, touch files outside the worktree.
- Data classes allowed: work, personal.

## Reply contract
- Write exactly one file for this turn: `plan.v1.md`, with one fenced `crow` block (schema `crow-block.schema.json`, contract 1).
- Block shape (YAML or JSON; drop the keys you do not use; list entries are objects, never bare strings):
```yaml
contract: 1
mission: <mission id>
kind: <kind>
phase: <plan|report|mockup>
version: <N>
summary: <one line>
confidence: <0-100>
verdict: <ready|needs-input|blocked>
questions:
  - {id: q1, text: <question>, blocking: true, choices: [{id: a, label: <choice>}, {id: b, label: <choice>}]}
responses:
  - {ref: a1, outcome: done, note: <quote of the answer>}
artifacts:
  - {path: outbox/<file>}
publish:
  - {action: <action>, target: <where>, body: <exact text>}
proof:
  - {what: command, cmd: <command>, exit: 0}
```
- Put any file you want shown under `outbox/` and list it in `artifacts`. Caps: 20 MB per file, 100 MB in total.
- A choice you cannot make: add it to `questions` with `blocking: true`, set `verdict: needs-input`, and end the turn. Outside a plan or report, write `outbox/questions.json`.
- Never wait for an answer inside the run. End the turn. core resumes you with `answers.v1.md`.
- From version 2, start the file with a `## Responses` section and fill `responses`.
- A report with `verdict: ready` carries `proof`.
- List every outward action in `publish` with its exact text. core publishes only what the owner validated.
