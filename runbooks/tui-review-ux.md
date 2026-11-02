# TUI review UX

Research note for the deck MR page. Each review tool below gave one pattern. The page rules at the end follow from them. Sources are the tools' own docs, read 2026-11-02.

## Current deck states and keys

- Findings: `pending`, `accepted`, `refused`, `edited`, `chat`. Per finding after a run: applied, pushed, replied, resolved.
- Review run: `starting`, `running`, `failed`, done. Apply run: `starting`, `running`, `failed`, committed (local only), pushed.
- Keys: `A` applies the decisions and commits locally. `P` pushes, then replies and resolves the threads (self and threads modes) or posts comments (peer mode). `0` and `1` pick the reply.
- Today the stage lives in three places: a progress line above the list, the card footer, and per-thread labels. The tab badge is `MRsAwaiting + MRsThreads`.

## Tools

### gh-dash
Pattern: one row per PR with fixed icon columns for state (open, draft, merged), review status (unreviewed, approved, changes requested) and CI (running, passed, failed), plus a comment count. The list answers "what is the shape of this PR" without opening it. Sections are saved queries, so "needs my review" is its own tab.
Source: https://www.gh-dash.dev/configuration/layout/pr/ and https://www.gh-dash.dev/configuration/pr-section/

### lazygit
Pattern: a bottom line with the keys valid right now, and a spinner on that same line while an action runs. Lists show a "n of m" count footer. Risky actions ask first. The status of a running action never moves into the content panes.
Source: https://github.com/jesseduffield/lazygit/blob/master/docs/Config.md

### tig
Pattern: a view stack. Enter on a row splits the screen with the detail below, Tab moves focus, `q` pops back. Every view has a title line with name, position ("commit 1 of 61") and a load timer after 3 seconds. The list stays on screen while the detail is read.
Source: https://jonas.github.io/tig/doc/manual.html

### GitHub PR review flow
Pattern: comments batch into a pending review that only the author sees. One "Review changes" step submits them with a verdict (comment, approve, request changes). Per-file "Viewed" marks drive a progress bar, and a changed file loses its mark. Pending is a visible state, and submit is one explicit step.
Source: https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/reviewing-changes-in-pull-requests/reviewing-proposed-changes-in-a-pull-request

### GitLab MR widgets
Pattern: merge blockers show as counters on the page. The open-threads counter turns orange while any thread is open and the merge check is on. Reviewers carry one of four states: awaiting review, started, approved, requested changes. A thread reply can carry a "resolve" action that runs when the review is submitted. Delayed actions are listed on the pending comment.
Source: https://docs.gitlab.com/user/project/merge_requests/ and https://docs.gitlab.com/user/project/merge_requests/reviews/

### Graphite
Pattern: the PR inbox is a set of sections ordered by whose move it is: needs your review, approved, returned to you, merging and recently merged, drafts, waiting for review. The section is the whose-move signal. Merged work stays visible for a while in its own section.
Source: https://graphite.com/docs/use-pr-inbox

### Reviewable
Pattern: sidebar counters are red when you must act, grey when someone else must act, and grey with a red stripe when you deferred. A "waiting on" list names who holds the ball and why (unreplied discussion, file changed since reviewed, requested reviewer). When the list is empty the author holds the ball, which means ready to merge. A completion rule defines done.
Source: https://docs.reviewable.io/reviews.html

### Gerrit attention set
Pattern: the set names who is expected to act. Rules add and remove people on events: a reply removes the replier, a reply on an unresolved thread adds the other side, a new patch set re-adds reviewers whose vote went outdated, submit or abandon clears the set. The dashboard opens with a "Your turn" section and a waiting-time column. Draft (work in progress) changes get no automatic additions. The state change is driven by events, not by a manual flag.
Source: https://gerrit-review.googlesource.com/Documentation/user-attention-set.html

### Phabricator Differential
Pattern: the home screen has two groups, "Action Required" (review, revise or push) and "Waiting on Others". Accepted stays accepted after the author updates, so a small fix does not send the revision back to the start. The author drives the last step.
Source: https://secure.phabricator.com/book/phabricator/article/differential/ and https://secure.phabricator.com/book/phabricator/article/differential_faq/

## Principles for the deck MR page

Pain points refer to the req_013 Needs list: P1 decision keys lost under a thin card, P2 no sign that `A` started, P3 card stale after `P`, P4 list says "fixed" after the push, P5 no "done", P6 badge shows 4 instead of 2, P7 waiting MRs look like mine.

1. One stage rail per MR, always in the same place. Seven fixed steps: review, decide, apply, commit, push, resolve, done. The current step is bright, past steps are green, future steps are dim. Next to it sit the owner of the move and the key. Source pattern: GitHub progress bar, tig title line. Ties to P2, P3, P5.
2. The footer line carries the key that is valid now, with the running state beside it. A running apply shows a spinner and "background, pushes nothing" on that line. Source pattern: lazygit bottom line. Ties to P1, P2.
3. The list is grouped by whose move, not by MR number. Group order: your move, waiting on others, done. Source pattern: Graphite sections, Differential groups. Ties to P7.
4. Each row ends with one reason in plain words ("2 decisions pending", "red pipeline", "committed, not pushed", "waiting for review"). Source pattern: gh-dash icon columns, Reviewable "waiting on". Ties to P4, P7.
5. List states follow the walkthrough, one word each: deciding, applying, committed, pushed, resolved, done. "Fixed" is removed. Source pattern: GitLab reviewer states. Ties to P4.
6. A state change is event driven. `A` and `P` write the new state to the card and the list at once, then a refresh confirms it. A stale label is never shown for a state the owner just changed. Source pattern: Gerrit attention set rules. Ties to P3, P4.
7. Waiting rows are dim, your-move rows use the accent colour and a bold reason. Colour carries the meaning (green done, yellow your move, red blocked, dim waiting) and text still carries it without colour. No fills. Source pattern: Reviewable red and grey counters. Ties to P7.
8. Done is a state with a name. When all threads are resolved and the commit is pushed, the card shows one "done" line with the commit, the thread count and the next owner (ask a reviewer, or merge). The row moves to the done group. Source pattern: Reviewable completion rule. Ties to P5.
9. The decision keys and the proposed fix stay pinned at the card bottom in every layout. The thread list scrolls above them. Source pattern: tig split view. Ties to P1.
10. The tab badge counts MRs where it is your move, one per MR. Open threads are shown inside the row, not summed into the badge. Source pattern: Graphite inbox, Differential "Action Required". Ties to P6.

## Whose-move rule

The rule gives each MR one owner of the next move and one reason. Evaluate top to bottom, first match wins. Events decide, as in the attention set. A new push or a reviewer reply re-evaluates the rule.

| Order | Condition | Owner | Reason text | Key |
| --- | --- | --- | --- | --- |
| 1 | Merged or closed, or all threads resolved and pushed | nobody | done | none |
| 2 | Apply run starting or running | agent | applying | none |
| 3 | Pipeline red on my MR | me | red pipeline | open logs |
| 4 | Commit exists, not pushed | me | committed, not pushed | `P` |
| 5 | Findings pending a decision | me | N decisions pending | decide, then `A` |
| 6 | Someone requested my review, or a reviewer replied after my last push | me | reviewer waiting on me | review |
| 7 | Open threads on my MR with no fix yet | me | N open threads | decide |
| 8 | Draft | me | draft | none |
| 9 | Pushed and replies or approval pending | them | waiting for review | none |
| 10 | Approved, merge pending | them | waiting for merge | none |

Notes:
- Rows 3 to 7 are the "your move" group. Rows 8 to 10 are the "waiting" group. A draft reads quiet and sorts last in its group.
- A red pipeline outranks everything except an apply run in flight.
- For an MR where I am the reviewer, rows 5 to 7 do not apply. Waiting on the author is row 9 with the reason "waiting for author".
- Each reason is one short phrase and appears on the list row and on the stage rail.

## Badge options for item_050

Counts read from the rule above.

- Option A: your-move MR count. Number of MRs in rows 3 to 7. Matches the owner's expectation of 2 when two MRs need action. Dim when 0.
- Option B: your-move over total, for example `2/5`. Shows the MR count and the size of the queue. Costs three more characters on the tab.
- Option C: your-move count with a red marker when any of those MRs is blocked (red pipeline, or a reviewer waiting). Same number as A, with urgency on a second channel.

## TUI constraints

- Target sizes: 160x45 laptop, 240x62 wide, 118x95 vertical. Lay out by aspect ratio: list left and card right on the landscape sizes, list above card on the vertical size.
- The stage rail needs one line. Seven short labels fit in 60 columns. Under 60 columns it collapses to the current step and the next key.
- ANSI 16 colours only, from the terminal theme. Transparent background, no fills, no truecolor. State uses bold, dim and the green, yellow, red slots.
- Glyphs must exist in the fallback font set. Check with `runbooks/tui-render.md`.
- Render options as PNG and compare side by side before building (`bin/ansi-shot`, `bin/crow-compare`).
