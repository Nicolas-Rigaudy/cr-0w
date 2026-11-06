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
- caps: 120 min, 25 USD, 2 resumes
- current version: 1 (write `plan.v1.md`)

## Owner preferences
- Prose: short and declarative. No em-dashes.
- Commit subjects: `type[SCOPE]: message`.
- Never push or merge before the report gate passes.

## Answers so far
None yet.

## Constraints
- Allowed: read and edit inside the worktree, run the repo checks, commit on the mission branch.
- Not allowed: push, merge, post, install packages, touch files outside the worktree.
- Data classes allowed: code, docs.

## Reply contract
- Write exactly one file for this turn: `plan.v1.md`, with one fenced `crow` block (schema: `crow-block.schema.json`, contract 1).
- Put any file you want shown under `outbox/` and list it in `artifacts`.
- A choice you cannot make: add it to `questions` with `blocking: true`, set `verdict: needs-input`, and end the turn.
- Never wait for an answer inside the run. End the turn. core resumes you with `answers.v1.md`.
- From version 2, start the file with a `## Responses` section and fill `responses`.
