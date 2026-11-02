# Review findings (deck review runs)

A read-only background run reviews one MR and its findings land in a JSON file the deck renders.

## Start, sync, check
- `bin/crow-review-run start <group/name> <number> [--mode peer|self]`: writes the context bundle, dispatches the run, writes the file with `run.state: running`. Default model `opus`, power `medium`. Refuses while a run is going or when a finding already has a decision (`--force` starts over).
- `bin/crow-review-run sync <group/name> <number>`: when cdx reports the run finished, parses the final message and writes the findings (`done`) or the reason (`failed`). A no-op otherwise. `show` syncs then prints the file.
- `bin/crow-review-run check <file>`: validates a file against the schema below.

## Read only
- The run is `crow-runs dispatch --permission review`, which is `cdx run --permission review`: Claude `--permission-mode plan`, Codex `-s read-only`. It cannot edit or write. It may share a folder with another run.
- Plan mode cannot write the findings file, so the run prints one JSON object as its final message and crow writes the file.
- Plan mode may deny Bash. crow puts everything the review needs in a bundle under the repo's git dir (`<git-common-dir>/crow-review/<number>.md`): the review ledgers (`local/reviews/common.md` and `<repo-dir>.md`), the MR threads and the full diff from `crow-mrs show`. That path is inside the run's cwd and outside `git status`.
- The prompt is saved next to the findings file as `<group>_<name>-<number>.prompt.md`.

## Resume for the apply phase
`cdx run --continue-run <run.id>` resumes the same Claude conversation (`claude --resume <conversation_id>`), so the apply run sees the bundle, the reasoning and the findings. Same account and provider, the parent's model and effort (cdx pins them and refuses a different `--model`), no `--failover` or `--isolate-workspace`. The permission is not pinned: the apply run can pass `auto`.

## File: `<state>/reviews/<group>_<name>-<number>.json`
```json
{
 "schema": 1, "repo": "group/name", "number": 214, "project": "registry name", "mode": "peer|self",
 "run": {"id": "cdx run id", "state": "running|done|failed", "reason": "why it failed",
         "model": "opus", "started_at": "ISO", "ended_at": "ISO"},
 "summary": "Patterns: N checked, K hits, M n/a. One-line verdict.",
 "findings": [{
   "id": "f1", "title": "short", "severity": "blocking|suggestion|nit|question",
   "path": "file from the diff, empty for a commit or MR-level finding", "line": 42,
   "body": "the problem", "comment": "text to post", "fix": "proposed fix (diff or replacement)",
   "pattern": "ledger row title or empty", "seen": 3,
   "state": "pending|accepted|refused|edited|chat", "note": "owner's edit or chat note"}]
}
```
- Required: `schema`, `repo`, `number`, `project`, `mode`, `run.state` (`run.reason` when failed), and per finding `id` (unique), `title`, `body`, `severity`, `state`.
- `line` is the new-file line in the diff, 0 when none.
- The run prints only `summary` and `findings` without `id`, `state` or `note`; sync adds them (`state: pending`).
- An older file holding a bare list of findings is still read by `crow-mrs show` (no run).
- Fixture: `tests/fixtures/review-findings.json`. Tests: `tests/test_crow_review_run.py` (simulated runs, fake `cdx`, no forge call).
