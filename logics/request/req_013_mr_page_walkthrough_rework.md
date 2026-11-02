## req_013_mr_page_walkthrough_rework - MR page walkthrough rework in crow deck
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Complexity: Medium
> Theme: Hub observability
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.
> Indicators reviewed: 2026-11-02 18:23:56

# AI Context
- Summary: Rework the deck MRs tab so a review walkthrough (review, decide, apply, push, reply, resolve) reads as one clear sequence: where the MR stands, whose move it is, what the next key does, and when it is done.
- Keywords: crow deck, MRs tab, walkthrough, review, threads, UX, whose move, badge
- Use when: changing the layout or states of the deck MRs tab.
- Skip when: crow-review-run behaviour (the backend works; this is presentation).

# Needs
- The owner reviews MRs only from the deck. On !214 (2026-11-02) the flow worked end to end but was hard to follow:
  - The decision keys and the proposed fix got lost under a thin card. Partly fixed (wider card, pinned buttons, commits 6f0f1c2, 4245910).
  - No sign that `A` had started, where it was, or that it ran in the background. Partly fixed (apply line).
  - After `P` pushed, the card still said "nothing pushed yet" (stale for one refresh), and the push state was hard to find.
  - The list said "fixed" after the push; it should move to pushed, then resolved.
  - Once everything was resolved, nothing said "done".
  - The tab badge showed 4 (open threads summed) where the owner expected 2 (the MR count).
  - MRs waiting on someone else look the same as MRs where the owner must act.
- Owner, verbatim: "needs a bit of rework and research on UI UX of this MR page cause here hard to follow the walkthrough".

# Context
- Code: `tools/crow-deck/mrs.go` (list, diff, card, foot), `mrs_data.go`, `chrome.go:145` (badge via `fixtures.mrsNeedYou`), `main.go` (detail refresh, postDoneMsg).
- Backend contract: `runbooks/review-findings.md`. Render pipeline: `runbooks/tui-render.md`, `bin/ansi-shot`, `bin/crow-compare`.
- Owner prefs: visual choices as rendered images side by side, never ASCII; keep info visible; terminal ANSI theme (lazytf look). See memories design-choices-ask-first, visual-decisions-side-by-side, deck-ux2-final.

# Acceptance criteria
- AC1: At any point of the walkthrough, one place on screen says the stage (reviewing, deciding, applying, committed, pushed, replying/resolving, done), whose move it is, and the key that moves it on.
- AC2: The MR list marks each MR as your move (highlighted) or waiting on someone else (greyed), with the reason.
- AC3: The tab badge shows a count the owner reads as intended (decided in item_050), never a sum of threads.
- AC4: A finished walkthrough shows a clear done state on the card and in the list.
- AC5: After any `A` or `P` action the card reflects the new state without a manual refresh.
- AC6: The layout was chosen by the owner from rendered options before it was built.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Scope
- In: MRs tab presentation, list states, badge, stage/progress display, refresh after actions.
- Out: crow-review-run logic, new forge actions, other deck tabs (Home Blockers may reuse the whose-move rule).

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `tools/crow-deck/mrs.go`
- `runbooks/review-findings.md`
- `local/state/shots/deck-threads/page.html` (current state, real !214)

# Backlog
- item_049_mr_page_ux_research (first)
- item_050_mr_page_walkthrough_mockups (after 049)
- item_051_mr_page_walkthrough_build (after 050, owner's pick)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

# Notes
- 2026-11-02 owner pick (canvas boards Rw_MRs_*): layout B focused walkthrough plus more code context around the highlighted line, merged board in req_011 item_055. Badge: count of MRs on the owner's move, colored dot only when there is a problem.
