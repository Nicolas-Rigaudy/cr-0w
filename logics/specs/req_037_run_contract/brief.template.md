# Brief: {mission_id}

<!-- Written by core from the kind file, the mission record and the answers. The run reads it first and never edits it. core rewrites it before each resume. -->

## Goal
{goal: one or two sentences, the owner's words}

## Mission
- kind: {kind}
- gates: {gates, or "none (straight through)"}
- key: {key}
- scope: {scope}
- repo: {repo placeholder} on branch {branch} (worktree {worktree}) <!-- omit when the kind has no worktree -->
- caps: {time_min} min, {budget_usd} USD, {resumes} resumes{, rounds N}
- current version: {N} (write `{phase}.v{N}.md`)

## Owner preferences
{owner prefs relevant to this kind: prose rules, commit format, language}

## Preflight
{one line per preflight answer, or "skipped: safest defaults" with the list}

## Answers so far
{for each answers.vK.md: version, verb, the owner's words verbatim. "None yet." on the first run}

## Constraints
- Allowed: {from kind autonomy, permission, worktree}
- Not allowed: {publish actions before the gate, anything outside the worktree, installs}
- Data classes allowed: {data_classes}

## Reply contract
- Write exactly one file for this turn: `{phase}.v{N}.md`, with one fenced `crow` block (schema `crow-block.schema.json`, contract 1).
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
- Put any file you want shown under `outbox/` and list it in `artifacts`. Caps: {outbox_file_mb} MB per file, {outbox_total_mb} MB in total.
- A choice you cannot make: add it to `questions` with `blocking: true`, set `verdict: needs-input`, and end the turn. Outside a plan or report, write `outbox/questions.json`.
- Never wait for an answer inside the run. End the turn. core resumes you with `answers.v{N}.md`.
- From version 2, start the file with a `## Responses` section and fill `responses`.
- A report with `verdict: ready` carries `proof`.
- List every outward action in `publish` with its exact text. core publishes only what the owner validated.
