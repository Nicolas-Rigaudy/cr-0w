## item_045_apply_run_resume_the_review_session_to_fix_accepted_findings_commit_locally - Apply run: resume the review session to fix accepted findings, commit locally
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Resume the review session to apply the accepted findings.
- Keywords: apply, run, resume, review, session, fix, accepted, findings, commit, locally
- Use when: building the deck review flow (req_010).
- Skip when: the work is outside the deck review flow.

# Problem
- Resume the review session to apply the accepted findings. Part of req_010. Needs items 043 and 044. First test on !214.

# Scope
- In:
  - A second dispatch with --resumes on the review run applies accepted and edited findings, commits locally as type[SCOPE] and reports a decision table. No push.
  - Findings in state `chat`: the apply run reads the owner's note, answers it in the decision table and changes no code for them. The answer is stored on the finding so the deck shows it; the owner then decides with a/x/e.
- Out:
  - Write rights only in this phase and only on a worktree under <repo>.wt/.

# Acceptance criteria
- AC1: Only accepted findings are changed.
- AC2: Commit subject follows type[SCOPE].
- AC3: Nothing is pushed.
- AC4: If the resume loses context the deck says so.
- AC5: Each chat note gets an answer shown on its finding card, with no code change.
- AC6: Tested with simulated runs only (owner decision 2026-11-02).

# AC Traceability
- req_010 AC1-AC5 -> this slice. Proof: its own acceptance criteria.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: (to be linked)
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Default until groomed.

# Notes
- Generated locally by logics-manager.
- Shipped: `crow-review-run start --branch` reviews in `<repo>.wt/review-<n>` (detached on the fetched branch, bundle in `.crow-review/` excluded from git). `crow-review-run apply` resumes the last finished run with `crow-runs dispatch --continue-run --cwd` (new flags; cdx pins model and effort). Sync folds the result onto each finding (`apply.status/detail/for`) and checks the commit in the worktree. Deck `A`: confirm overlay with counts, then `crow-review-run apply`; card shows the commit (not pushed), warnings, the answer to a chat note, and flags an answer older than the current decision.
- Peer reviews: apply only answers chat notes; accepted findings are posted by item_046.
- Proofs: AC1 `apply_prompt` lists only the to-do findings (`test_apply_resumes_the_review_with_the_decisions`) and sync warns on files outside them (`test_apply_warns_on_bad_subject_extra_files_and_pushed_commit`). AC2 prompt rule plus the `crow-commit-msg-check` warning (same test). AC3 prompt rule, no push in the code path, warning when the commit is on a remote. AC4 `test_apply_failures_say_why` (context lost), `test_apply_records_a_resume_cdx_refuses`, `TestApplyFailureSaysWhy`. AC5 `test_apply_result_lands_on_the_findings`, `TestApplyResultShows`. AC6 all tests use fake dispatch, a local bare origin and recorded `runBin`; no real run, deck action or forge call.
- Open risk from req_010: whether `--continue-run` keeps enough context on a large diff (!214) is only checkable with a real run; the run answers `context: lost` when it cannot see its review.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

