## req_071_nest_mrs_review_flow_and_visible_session_starts - Nest MRs review flow and visible session starts
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Nest MRs runs the existing review flow end to end, and every session start either opens on Desk or says why it failed.
- Keywords: nest, mrs, review, flow, visible, session, starts
- Use when: building the Nest review flow or fixing session starts from Nest.
- Skip when: the deck (TUI) review page; it already works.

# Needs
- 2026-11-26, owner on a no-terminal day: "launched the review for mr 245 from nest and the run finished but nothing got shown to me". The result then came back as a context-free question in an unrelated hub session.
- Owner: "when clicking repo session on mr, it told me it launched a session which i never saw".
- The MRs tab must do the review from start to finish: start, findings, decisions, apply, post. No hub session in between.

# Context
- Nest Review dispatches a plain `/review <url>` run (`tools/nest/src/data/mrs.ts` reviewPrompt, `Mrs.tsx` confirmsOf). item_225 left the review flow out ("Skip when: the review run itself").
- The flow exists for the deck: `bin/crow-review-run` start/sync/show/decide/apply/actions/act, findings in `local/state/reviews/<group>_<name>-<n>.json`, `crow-mrs show` returns findings and run state. Spec: `runbooks/review-findings.md`.
- Core has no endpoint over crow-review-run. `/api/state/mrs` carries the list only, no findings.
- Session start: core runs `tmux new-session -d ... crow open --on-busy fail <project>` and answers 200 at once. On 11-26 10:41 the session exited after 1.4 s; two lambda starts at 10:39 left no session either. Nest showed "started".

# Acceptance criteria
- AC1: Review on an MR in Nest starts `crow-review-run start` (peer, self or threads mode, with the source branch). No plain dispatch is left.
- AC2: The MR detail shows the run state and, when done, each finding with file, line, comment and proposed fix, and records accept, refuse, edit or chat decisions.
- AC3: Apply and post run from the MR detail, one gated action at a time with the exact text shown first. Nothing posts without the owner's OK.
- AC4: A review run never surfaces as a question in a hub session.
- AC5: A session start that dies within a few seconds returns an error with its last output; Nest shows it. A start that lives opens that session on Desk.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-review-run`, `bin/crow-mrs`, `runbooks/review-findings.md`
- `tools/nest/src/screens/Mrs.tsx`, `tools/nest/src/data/mrs.ts`
- `tools/crow-core` sessions start handler (`POST /api/sessions`)

# Backlog
- item_298_session_start_fails_loudly_and_opens_on_desk (first, independent)
- item_299_core_review_endpoints_over_crow_review_run
- item_300_mrs_review_tab_mockup_options (parallel with 299, needs the owner's pick)
- item_301_mrs_review_tab_build (after 299 and 300)

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |

