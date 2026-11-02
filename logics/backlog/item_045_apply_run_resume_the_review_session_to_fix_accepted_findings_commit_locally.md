## item_045_apply_run_resume_the_review_session_to_fix_accepted_findings_commit_locally - Apply run: resume the review session to fix accepted findings, commit locally
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
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
- Out:
  - Write rights only in this phase and only on a worktree under <repo>.wt/.

# Acceptance criteria
- AC1: Only accepted findings are changed.
- AC2: Commit subject follows type[SCOPE].
- AC3: Nothing is pushed.
- AC4: If the resume loses context the deck says so.

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

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

