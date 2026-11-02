## req_010_review_mrs_inside_the_crow_deck_with_a_background_claude_run - Review MRs inside the crow deck with a background Claude run
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Review an MR from the crow deck only: a background Claude run reviews, the deck shows findings and takes the owner's decisions, a second phase applies them. No terminal Claude needed.
- Keywords: review, mrs, inside, crow, deck, background, claude, run
- Use when: building or changing the deck review flow.
- Skip when: the work is outside the deck MRs tab.

# Needs
- Today the MRs tab is read only and nothing runs behind it. The owner wants the deck to be the review interface: one place to read findings and diffs, decide, and see the result, without working in a Claude terminal.
- Flow: (1) a read-only background run reviews the MR and writes a findings JSON; (2) the deck shows threads, diff and proposed fixes and records a/x/e/c decisions in that JSON; (3) the same Claude session is resumed (--resumes) to apply accepted findings and commit locally; (4) push, reply and resolve go through one confirm overlay per action showing the exact text.
- Why two phases of one session: the human decision is slow and asynchronous, the review needs no write rights, and the JSON is the contract the deck can render. A live window driven by crow-focus --send was rejected as fragile (no structured data, wrong keystrokes land in the session).
- Escape hatch: c opens an interactive session resumed from the review for cases that need a real discussion.

# Context
- Spec: local/state/research/crow-deck-design.md. Decisions: confirm overlay per action, one OK per action (decisions-2026-11-02.md). Findings path already planned: <state>/reviews/<group>_<name>-<number>.json.
- Risk to check first on !214: does --resumes keep enough context on a large diff.
- Autonomy stays as in CLAUDE.md: the deck and the runs never push, merge or post without the owner's OK on the exact text.

# Acceptance criteria
- AC1: From the MRs tab, one key starts a background review of the selected MR and the findings appear in the deck without opening a terminal.
- AC2: The owner can accept, refuse, edit or chat each finding in the deck and the decision is saved in the findings JSON.
- AC3: An apply phase resumes the same session, fixes accepted findings and commits locally only.
- AC4: Push, reply and resolve each show the exact text in an overlay and run only on enter.
- AC5: Tests use simulated runs; nothing in the tests launches a real run or touches a forge.

# Definition of Ready (DoR)
- [ ] Problem statement is explicit and user impact is clear.
- [ ] Scope boundaries (in/out) are explicit.
- [ ] Acceptance criteria are testable.
- [ ] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `logics_manager/flow.py`
- `logics_manager/assist.py`
- `tests/python/test_logics_manager_cli.py`

# Backlog
- item_043, item_044, item_045, item_046, item_047, item_048

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

