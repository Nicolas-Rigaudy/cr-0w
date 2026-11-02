# Review findings (deck review runs)

A read-only background run reviews one MR and its findings land in a JSON file the deck renders.

## Start, sync, check
- `bin/crow-review-run start <group/name> <number> [--branch B] [--mode peer|self]`: writes the context bundle, dispatches the run, writes the file with `run.state: running`. Default model `opus`, power `medium`. Refuses while a run is going or when a finding already has a decision (`--force` starts over).
- With `--branch` (the deck always passes the MR's source branch) the run works in `<repo path>.wt/review-<number>`, detached on `origin/<branch>` after a fetch. A second start moves that worktree to the branch tip; it refuses while the worktree holds changes or local commits not on the remote (apply output not pushed yet). `--force` drops them.
- `bin/crow-review-run sync <group/name> <number>`: when cdx reports the run finished, parses the final message and writes the findings (`done`) or the reason (`failed`). A no-op otherwise. `show` syncs then prints the file.
- `bin/crow-review-run check <file>`: validates a file against the schema below.
- `bin/crow-review-run decide <group/name> <number> <fN> accepted|refused|edited|chat [--note TEXT]`: records the owner's decision on one finding. `edited` and `chat` need `--note`.
- `crow-mrs show` syncs a running review, so the deck sees it end without a second call.

- `bin/crow-review-run apply <group/name> <number>`: resumes the review conversation to act on the decisions (below).

## Deck keys (MRs tab)
- `r` asks, then starts a review of the selected MR (`peer` on a review, `self` on an own MR). Refused while one runs.
- `a` accept, `x` refuse, `e` edit the comment (prefilled), `c` chat note. `enter` saves, `esc` drops. Each decision calls `decide`, so it survives a restart.
- `A` asks (counts of fixes and chat notes), then starts the apply run. Refused before the review is done, while an apply runs, or when nothing is left to apply.
- A running review or apply refetches every 20 s. A failed one shows its reason; `r` runs the review again.
- The mock deck keeps decisions in memory and starts no run.

## Read only
- The run is `crow-runs dispatch --permission review`, which is `cdx run --permission review`: Claude `--permission-mode plan`, Codex `-s read-only`. It cannot edit or write. It may share a folder with another run.
- Plan mode cannot write the findings file, so the run prints one JSON object as its final message and crow writes the file.
- Plan mode may deny Bash. crow puts everything the review needs in a bundle: the review ledgers (`local/reviews/common.md` and `<repo-dir>.md`), the MR threads and the full diff from `crow-mrs show`. In a worktree it is `<worktree>/.crow-review/<number>.md` (`/.crow-review/` added to the repo's `info/exclude`), else `<git-common-dir>/crow-review/<number>.md`. Both are inside the run's cwd and outside `git status`.
- The prompt is saved next to the findings file as `<group>_<name>-<number>.prompt.md`.

## Apply phase
`crow-review-run apply` dispatches `crow-runs dispatch --continue-run <parent> --cwd <worktree>`, which is `cdx run --continue-run`: the same Claude conversation (`claude --resume <conversation_id>`), so the run sees the bundle, its reasoning and its findings. cdx pins the parent's account, provider, model and effort (crow passes no `--model`/`--power`; cdx refuses a different one) and refuses `--failover`/`--isolate-workspace`. The permission is the account's write permission (`auto` on Claude).
- The run resumes in the review's worktree. A Claude transcript lives under its cwd, so the review and the apply share that cwd; a review run in the main checkout (no `--branch`) cannot be applied.
- What it acts on: chat notes in any mode; accepted and edited findings on a `self` review (a `peer` review posts its accepted findings instead, item_046). A finding already handled for the same decision and note is skipped, so a second `apply` after new decisions only sends those. It resumes the last finished apply run, else the review run.
- Prompt (`<file>.apply.md`): fix accepted findings, fix edited ones as the note says, answer chat notes in the note's language with no code change, leave the rest alone, one local commit `type[SCOPE]: message`, never push or post. Final message: `{"context": "ok|lost", "commit", "subject", "results": [{"id", "status": "fixed|answered|skipped|failed", "detail"}]}`.
- Sync folds it like the review: results land on each finding as `apply`, the commit on `apply`. Checks read from the worktree add `apply.warnings`: subject not `type[SCOPE]` (crow-commit-msg-check), files outside the accepted findings' paths (tests allowed), commit already on a remote branch, a finding with no result.
- Lost context: `context: lost` from the run, or cdx refusing the resume (`continue_run_unavailable`, `continue_run_not_found`), marks `apply.state: failed` with the reason; the deck shows it and `r` runs the review again.

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
   "state": "pending|accepted|refused|edited|chat", "note": "owner's edit or chat note",
   "apply": {"status": "fixed|answered|skipped|failed", "detail": "what changed or the answer",
             "for": ["decision", "note it acted on"], "run": "apply run id"}}],
 "branch": "MR source branch", "worktree": "<repo path>.wt/review-<number>",
 "apply": {"id": "cdx run id", "state": "running|done|failed", "reason": "", "parent": "resumed run id",
           "asked": {"f1": ["accepted", ""]}, "commit": "sha", "subject": "type[SCOPE]: ...",
           "warnings": [], "started_at": "ISO", "ended_at": "ISO"}
}
```
- Required: `schema`, `repo`, `number`, `project`, `mode`, `run.state` (`run.reason` when failed), and per finding `id` (unique), `title`, `body`, `severity`, `state`. `apply` and the finding's `apply` are optional; when present they need a valid `state`/`status` (`apply.reason` when failed).
- `line` is the new-file line in the diff, 0 when none.
- The run prints only `summary` and `findings` without `id`, `state` or `note`; sync adds them (`state: pending`).
- An older file holding a bare list of findings is still read by `crow-mrs show` (no run).
- Fixture: `tests/fixtures/review-findings.json`. Tests: `tests/test_crow_review_run.py` (simulated runs, fake `cdx`, a local bare origin, no forge call), `tools/crow-deck/apply_test.go` and `review_test.go` (`runBin` recorded, nothing runs).
