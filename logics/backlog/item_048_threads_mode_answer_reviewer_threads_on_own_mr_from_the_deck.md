## item_048_threads_mode_answer_reviewer_threads_on_own_mr_from_the_deck - Threads mode: answer reviewers' threads on my own MR from the deck
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: On the owner's own MR, the deck takes each reviewer thread from proposal to reply and resolve: a background run proposes a fix and a reply per thread, the owner decides with a/x/e/c, the apply run commits, P pushes and posts.
- Keywords: threads, own MR, review-threads, reply, resolve, deck, crow-review-run
- Use when: building the threads side of the deck MRs tab.
- Skip when: peer review findings (items 043-046 already ship that).

# Problem
- The owner wants the deck to be the single interface for MR review. Peer review and self-review work end to end (items 043-046). Threads left by reviewers on the owner's own MR do not: the deck lists them, but in threads mode a/x/e/c only flash a message, and fixing still means `h`, a repo session and the review-threads skill (2026-11-02, !214 with 4 open threads).
- Keep it simple: same keys, same flow, same overlay as findings. No new screen.

# Scope
- In:
  - `crow-review-run start --mode threads` (default when the MR is the owner's and has open threads): read-only run that follows the review-threads skill (ledger match, fix proposal, French reply) and writes one entry per open thread into the same findings JSON (`Pattern`, `Fix`, `Test`, `Reply` fields the deck already renders).
  - Deck threads mode: `r` starts that run; a/x/e/c record decisions per thread like findings (e edits the reply, c adds a note for the apply run).
  - Apply run (item_045 path): fixes accepted threads in the worktree, commits locally, drafts the final reply text (commit hash included).
  - `actions`/`act` (item_046 path): one push, then per thread one reply and one resolve, each with its exact text, one OK per action. Refused threads get no reply unless the owner wrote one with e.
  - Ledger feed: accepted patterns go to local/reviews/ like the skill does.
- Out:
  - Chat escape hatch (item_047), the rest of item_020 (history, other actions).

# Acceptance criteria
- AC1: On an own MR with open threads, `r` in the deck starts a threads run; each thread card shows the proposed fix, test and French reply when it ends, without a terminal.
- AC2: a/x/e/c on a thread card save the decision in the findings JSON (edit changes the reply text).
- AC3: Apply commits the accepted fixes locally only, in the review worktree.
- AC4: `P` shows the push, then each reply and each resolve with its exact text; enter runs only that one; a thread whose head moved or text changed is refused.
- AC5: Tests use simulated runs and faked forge calls; nothing touches GitLab.
- AC6: Dry run on !214 (lambda, 4 threads) before closing: findings render, nothing posted without OK.

# AC Traceability
- req_010 AC1-AC4 (threads side) -> this slice. Proof: its own acceptance criteria.
- req_010 AC5 -> AC5.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed (reuses the findings JSON contract, runbooks/review-findings.md)

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: req_010_review_mrs_inside_the_crow_deck_with_a_background_claude_run
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Owner wants the deck to be the MR review interface; reviews unblock devs.

# Notes
- Builds on: bin/crow-review-run (start/apply/actions/act), tools/crow-deck/mrs.go threadsMode, skills/review-threads (logic to reuse, not duplicate), runbooks/review-findings.md (update the schema section for thread entries).
- Execution order: run mode + JSON schema, deck keys, apply, actions/act, !214 dry run.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Outcome (2026-11-02)
- AC1: `reviewMode` is `threads` on an own MR with open threads; `crow-review-run start --mode threads` keys the open threads T1..Tn and the deck shows one card per thread (reviewer comment, verdict, reply or "no reply", fix, test). Proof: tests/test_crow_review_run.py threads start/fold, tools/crow-deck/threads_test.go TestThreadsReviewShowsEachAnswer.
- AC2: a/x/e/c call `decide`; `e` prefills and rewrites the reply. Proof: TestThreadsDecisionsAndApply.
- AC3: apply acts on accepted or edited threads that carry a fix, plus chat notes; one local commit. Proof: test_threads_apply_fixes_only_threads_with_a_fix.
- AC4: `actions` lists the push first, then `reply-fN` and `resolve-fN`; a fixed thread waits for the push and resolves with no reply; `agree: no` stays open; `act` keeps the hash check. Proof: test_threads_actions_wait_for_the_push_then_reply_and_resolve, TestThreadsReplyAndResolveOverlays. Resolving feeds the ledgers once per thread (test_resolving_a_thread_feeds_the_ledgers).
- AC5: simulated runs, faked forge calls. 168 pytest and the Go suite pass.
- AC6: dry run on !214: read-only run a12e1ffa gave 4 findings (4 fixes, 0 replies, 4 ledger repeats); the deck renders them as `Threads 4`; nothing pushed, replied or resolved.
- Found on the way: `cdx runs --since` leaves running runs out, so a sync during a run marked it failed. `cdx_run` now also checks `cdx runs --running`.

