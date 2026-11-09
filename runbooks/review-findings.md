# Review findings (deck review runs)

A read-only background run reviews one MR and its findings land in a JSON file the deck renders.

## Start, sync, check
- `bin/crow-review-run start <group/name> <number> [--branch B] [--mode peer|self|threads]`: writes the context bundle, dispatches the run, writes the file with `run.state: running`. Default model `opus`, power `medium`. Refuses while a run is going or when a finding already has a decision (`--force` starts over).
- With `--branch` (the deck always passes the MR's source branch) the run works in `<repo path>.wt/review-<number>`, detached on `origin/<branch>` after a fetch. A second start moves that worktree to the branch tip; it refuses while the worktree holds changes or local commits not on the remote (apply output not pushed yet). `--force` drops them.
- `bin/crow-review-run sync <group/name> <number>`: when cdx reports the run finished, parses the final message and writes the findings (`done`) or the reason (`failed`). A no-op otherwise. `show` syncs then prints the file.
- `bin/crow-review-run check <file>`: validates a file against the schema below.
- `bin/crow-review-run decide <group/name> <number> <fN> accepted|refused|edited|chat [--note TEXT]`: records the owner's decision on one finding. `edited` and `chat` need `--note`.
- `crow-mrs show` syncs a running review, so the deck sees it end without a second call.

- `bin/crow-review-run apply <group/name> <number>`: resumes the review conversation to act on the decisions (below).
- `bin/crow-review-run actions <group/name> <number>` and `act <group/name> <number> <id> --hash H`: push and post, one action per call (below).

## Deck page (MRs tab)
Three columns in landscape: the MR list, the code, the Your move panel. Portrait stacks the stage line, the code (or the list on `tab`) and a compact panel.
- One rule sets each MR's stage (`mrStageOf` in `tools/crow-deck/mrs_move.go`): steps review, decide, apply, commit, push, resolve (threads mode), done. Peer reviews use review, decide, post, done. It also sets who moves next (you, crow, others, done), the list reason and the key that moves it on.
- The list groups MRs by whose move: Your move, Waiting on others, Done. A done MR shows `✓` and `done, pushed <sha>`.
- The MRs tab badge counts the MRs on the owner's move. A dot follows the count only for a real problem: a failed pipeline, review or apply.
- The stage box above the code shows `✓ done ─ ● current ─ ○ next` and `next <keys>`.
- The code pane opens on the current finding's file, the commented line reversed with `▶` and `◀ Fx author`, the comment boxed under it, other findings marked `◆`. A strip lists the files with findings. A scrollbar and `↑ N above · ↓ M below · lines a–b` give the position.
- The panel lists the steps (current one marked `▶ ●`), the finding, the proposed fix, then the keycaps `a x e c` and `n` / `A` / `P`. Once the finding is pushed or posted, the keycaps fold into one `Decided Fx` line. A done MR shows the commit, the thread count and who acts next.
- Keys: `j k` scroll the focused pane, `J K` jump by hunk, `[ ]` by file, `n p` move between findings (`n` past the last one opens the next MR), `tab` cycles list, code, panel.
- `a` accept, `x` refuse, `e` edit the comment or reply (prefilled), `c` chat note. `enter` saves, `esc` drops. Each decision calls `decide`, so it survives a restart.
- `r` asks, then starts a review of the selected MR: `peer` on a review, `threads` on an own MR with open threads, `self` on another own MR. Refused while one runs.
- `A` asks (counts of fixes and chat notes), then starts the apply run. Refused before the review is done, while an apply runs, or when nothing is left to apply.
- `P` lists the pending actions, then shows them one overlay at a time: title, exact text, exact command, warnings. `enter` runs that one action, `s` skips it, `esc` stops. Keys wait while it runs. A failure closes the overlay with the reason.
- A post that succeeds and an apply that starts refetch the MR at once. A running review or apply refetches every 20 s. A failed one shows its reason; `r` runs the review again.
- The mock deck keeps decisions in memory and starts no run.

## Threads mode (answer reviewers on the owner's MR)
- `start --mode self` on an MR with open threads runs in threads mode: the deck picks the mode from a list that can predate the threads (!213, 2026-11-03).
- `start --mode threads` refuses an MR with no open thread. It keys the unresolved threads T1..Tn in the bundle ("## Open threads to answer") and stores them as `threads` (key, forge id, GitHub comment id, author, body, path, line).
- The prompt follows the review-threads skill: one finding per thread, `agree` yes|no|partly, `fix` and `test` when code changes, `comment` the reply in the reviewer's language and empty when the fix answers it, `pattern` matched against the ledgers (`pattern_kind` repeat|new, `how` for a new row).
- Sync drops findings whose key is not an open thread and attaches `thread` to each one.
- The deck shows these findings in place of the raw threads (title `Threads N`, `T` flips back): reviewer comment, verdict, reply or "no reply", fix, test. `e` rewrites the reply; the fix stays.
- Layout while answering: the card takes about half the width (55% of the height in portrait). Order: comment, verdict, reply, fix, test, apply result, reasoning, pattern; a short window cuts from the end. Buttons and the run/apply status stay pinned at the bottom. The diff opens at the thread's file and line (`thread.path`/`thread.line`).
- `A` shows `apply starting` at once, then `running Ns`, then `committed <sha> · local only` or `failed: reason`. The list shows each thread's apply result (fixed, answered, failed, skipped).
- `A` applies accepted or edited findings that have a fix, plus chat notes. Reply-only threads need no apply run.
- `P`: the push first (when there is a commit), then per thread a `reply` (when the reply is not empty) and a `resolve` (unless `agree: no`, which stays open for the reviewer). A thread with a fix waits until its commit is pushed. Fixed threads are resolved with no reply.
- GitLab: `glab api -X POST .../discussions/<id>/notes -f body=...`, `glab api -X PUT .../discussions/<id> -f resolved=true`. GitHub: `gh api -X POST repos/<r>/pulls/<n>/comments/<comment_id>/replies`, GraphQL `resolveReviewThread`.
- `act` records `posted` on a reply and `resolved` on a resolve. Resolving feeds the ledgers once per thread (`ledger` on the finding): a matching row in `<repo-dir>.md` then `common.md` gets its count bumped and the MR appended, otherwise a new row goes to `<repo-dir>.md`. `agree: no` feeds nothing.

## Read only
- The run is `crow-runs dispatch --permission review`, which is `cdx run --permission review`: Claude `--permission-mode plan`, Codex `-s read-only`. It cannot edit or write. It may share a folder with another run.
- Plan mode cannot write the findings file, so the run prints one JSON object as its final message and crow writes the file.
- Plan mode may deny Bash. crow puts everything the review needs in a bundle: the review ledgers (`local/reviews/common.md` and `<repo-dir>.md`), the MR threads and the full diff from `crow-mrs show`. In a worktree it is `<worktree>/.crow-review/<number>.md` (`/.crow-review/` added to the repo's `info/exclude`), else `<git-common-dir>/crow-review/<number>.md`. Both are inside the run's cwd and outside `git status`.
- The prompt is saved next to the findings file as `<group>_<name>-<number>.prompt.md`.

## Apply phase
`crow-review-run apply` dispatches `crow-runs dispatch --continue-run <parent> --cwd <worktree>`, which is `cdx run --continue-run`: the same Claude conversation (`claude --resume <conversation_id>`), so the run sees the bundle, its reasoning and its findings. cdx pins the parent's account, provider, model and effort (crow passes no `--model`/`--power`; cdx refuses a different one) and refuses `--failover`. An isolated parent stays isolated: cdx resumes it in the same sandbox with the same permission. The permission is `auto` on Claude with `--strict-mcp-config` and no allowlist (cdx refuses `--settings` on a continued run; see dispatch-background-runs.md, Run sandbox).
- Isolation (item_076): a peer review with a worktree dispatches `--isolate` (`isolated: true` in the file): someone else's code, read in cdx's OS sandbox with no credentials. Its apply only answers chat notes and stays read-only in the same sandbox. Self and threads reviews run unisolated, since their apply writes the fixes.
- The run resumes in the review's worktree. A Claude transcript lives under its cwd, so the review and the apply share that cwd; a review run in the main checkout (no `--branch`) cannot be applied.
- What it acts on: chat notes in any mode; accepted and edited findings on a `self` review (a `peer` review posts its accepted findings instead, item_046). A finding already handled for the same decision and note is skipped, so a second `apply` after new decisions only sends those. It resumes the last finished apply run, else the review run.
- Prompt (`<file>.apply.md`): fix accepted findings, fix edited ones as the note says, answer chat notes in the note's language with no code change, leave the rest alone, one local commit `type[SCOPE]: message`, never push or post. Final message: `{"context": "ok|lost", "commit", "subject", "results": [{"id", "status": "fixed|answered|skipped|failed", "detail"}]}`.
- Sync folds it like the review: results land on each finding as `apply`, the commit on `apply`. Checks read from the worktree add `apply.warnings`: subject not `type[SCOPE]` (crow-commit-msg-check), files outside the accepted findings' paths (tests allowed), commit already on a remote branch, a finding with no result.
- Lost context: `context: lost` from the run, or cdx refusing the resume (`continue_run_unavailable`, `continue_run_not_found`), marks `apply.state: failed` with the reason; the deck shows it and `r` runs the review again.

## Push and post
The only commands that leave the machine. Each one runs after the owner's enter on its own overlay; that enter is the explicit OK for that action on a work repo.
- `actions` prints a JSON list: `{"id", "kind": "push|comment", "title", "text", "warnings", "cmd", "hash"}`. `hash` covers `cmd` and `text`.
- Self review: one `push` once the apply run committed: `git -C <worktree> push origin <sha>:refs/heads/<branch>` (fast-forward only, never forced). Text is the commit subject, warnings are `apply.warnings`. Recorded as `apply.pushed`.
- Peer review: one `comment` per accepted or edited finding not posted yet. Text is `comment`, or the owner's `note` when edited (French on work repos, as the review wrote it). With path and line: GitLab `POST .../discussions` with a text position, GitHub `POST .../pulls/<n>/comments` on the RIGHT side. Without: a note (`.../notes`, `gh pr comment`). Recorded as the finding's `posted`.
- `act` recomputes the list and refuses an id that is gone or whose hash changed (a decision changed after the overlay showed it). For a diff-line comment it fills `{base_sha}`, `{start_sha}`, `{head_sha}` from the MR and refuses when the MR head differs from the worktree's reviewed tip (line numbers would be off; `r` reviews again).
- Forge CLIs and git run with the real HOME (cdx remaps it).
- Threads mode adds `reply` and `resolve` actions (above).

## File: `<state>/reviews/<group>_<name>-<number>.json`
```json
{
 "schema": 1, "repo": "group/name", "number": 214, "project": "registry name", "mode": "peer|self|threads",
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
             "for": ["decision", "note it acted on"], "run": "apply run id"},
   "posted": {"at": "ISO", "text": "what was posted", "ref": "forge note id or URL"},
   "thread": {"key": "T1", "id": "forge thread id", "comment_id": 0, "author": "", "body": "", "path": "", "line": 0},
   "agree": "yes|no|partly", "test": "test diff", "pattern_kind": "repeat|new", "how": "ledger How",
   "resolved": {"at": "ISO"}, "ledger": "bumped in common.md"}],
 "threads": [{"key": "T1", "id": "...", "comment_id": 0, "author": "", "body": "", "path": "", "line": 0}],
 "branch": "MR source branch", "worktree": "<repo path>.wt/review-<number>", "isolated": false,
 "apply": {"id": "cdx run id", "state": "running|done|failed", "reason": "", "parent": "resumed run id",
           "asked": {"f1": ["accepted", ""]}, "commit": "sha", "subject": "type[SCOPE]: ...",
           "warnings": [], "started_at": "ISO", "ended_at": "ISO",
           "pushed": {"commit": "sha", "at": "ISO"}}
}
```
- Required: `schema`, `repo`, `number`, `project`, `mode`, `run.state` (`run.reason` when failed), and per finding `id` (unique), `title`, `body`, `severity`, `state`. `apply` and the finding's `apply` are optional; when present they need a valid `state`/`status` (`apply.reason` when failed).
- `line` is the new-file line in the diff, 0 when none.
- The run prints only `summary` and `findings` without `id`, `state` or `note`; sync adds them (`state: pending`).
- An older file holding a bare list of findings is still read by `crow-mrs show` (no run).
- Fixture: `tests/fixtures/review-findings.json`. Tests: `tests/test_crow_review_run.py` (simulated runs, fake `cdx`, a local bare origin for the push, forge calls injected), `tools/crow-deck/apply_test.go`, `review_test.go`, `post_test.go` and `threads_test.go` (`runBin` recorded, nothing runs).
