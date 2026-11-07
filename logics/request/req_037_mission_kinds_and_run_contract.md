## req_037_mission_kinds_and_run_contract - Mission kinds and run contract
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 92
> Confidence: 82
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Define every mission kind (gates, autonomy, account, effort, outputs) and the file contract between a run and core: brief, plan and report with a fenced `crow` block, outbox, versions, answer verbs. Schema first, before the first mission runs.
- Keywords: mission kind, gate, autonomy, run contract, crow block, outbox, plan, report, versions, answer verbs, preflight
- Use when: adding a mission kind, writing a dispatch prompt, parsing what a run asks.
- Skip when: choosing hosts, sandbox or UI stack (req_029 to req_031).

# Needs
1. **Kind catalogue.** For each kind: purpose, gates, default autonomy, allowed scopes, account and model, effort, time cap, outputs, publish action. Set decided 2026-11-06 (13 merged into 8, table in prod_005): `research`, `plan` (absorbs `notes`), `review`, `design` in phase 1; `review-fix`, `write` (absorbs `document` and `report`; reports started by a schedule) in phase 2; `build` (absorbs `ship`, `personal-code`, `crow`; gates set by scope) and `repair` (absorbs `rebase`, `ci-repair`) in phase 3 after the sandbox. `admin` became butler jobs.
   - **Kinds as data files.** One file per kind (`kinds/<kind>.yaml` or similar) read by core; a new kind needs no engine change.
2. **One skill per kind.** Ultron maps each kind to a skill and a key format (ticket key, MR iid, free slug). crow already has skills for most kinds (`review`, `review-threads`, `mr`, `daily`, `weekly`, `mgmt-weekly`). The kind names its skill; the brief adds the run contract.
3. **Preflight.** Before dispatch, the kind's preflight asks only what the run cannot decide itself (Jarvis). If the owner skips it, the run states its safest assumption in the plan and goes on.
4. **Gate shapes.**
   - Two-gate (Ultron): run writes `plan.md`, ends turn; owner answers `go`, `dig` or `no`; run writes `report.md`, ends turn; owner answers `ok` or `msg`; core publishes.
   - Mockup loop (Jarvis): run writes pages to `outbox/`; owner answers `approve` (commits pages to `docs/design/` and opens the build mission), `change` with a note (up to 5 rounds), or `abandon`. Approval only from the nest page.
   - Straight through: no gate; FYI at the end; personal repos auto-merge on green CI and a secret scan.
5. **The `crow` block.** One fenced block per plan or report, parsed by core. Fields: `kind`, `phase` (plan, report, mockup), `summary` (one sentence), `confidence` (0 to 100), `verdict` (ready, blocked, needs-input), `questions` (id, text, choices, default), `artifacts` (paths), `publish` (what core would do: push, MR, post, with exact titles and text), `proof` (tests run, commands and their exit codes, links checked; required on reports). Markdown around it stays for humans.
6. **Versions and responses.** Each answered version is kept with the owner's words. The next version starts with a "Responses" section saying what became of each point (Ultron "Réponses"). nest shows the history.
7. **Validation record.** `ok` writes a validation event tied to the report version. Any later message drops it. core checks it before every publish.
8. **Outbox.** `outbox/questions.json` for questions outside a plan or report; artifacts as files (HTML, PNG, MP4, diff). Size caps per file and per mission.
9. **Brief.** `brief.md` carries goal, kind, scope, end state, owner prefs, prior answers, constraints (time cap, no publish), and the reply contract. Built by core, never by hand.
10. **End-of-run rules.** A run ends its turn after writing; it never waits. A run that hits a usage limit ends; core resumes on the next account.

# Context
- Owner decisions 2026-11-06: gates per kind; preflight plus park and resume; personal repos auto-merge; work hard rules unchanged.
- Owner decisions 2026-11-06 (review): 8 kinds, phase 1 kinds research, plan, review, design; proof of work in reports; kinds as data files.
- AFK runs become batches of missions with preflight up front; their reports use this contract.
- Ultron lesson 2: define the plan and report schema first. Free text in a hook is hard to parse.
- crow already has: skills per task type, `crow dispatch` prompts, `local/owner-prefs.md`, review ledgers.
- adr_001 sets the mission dir and the event log; this request fills in the formats.

# Spec
Draft spec in `logics/specs/req_037_run_contract/` (groomed 2026-11-06): kind schema and the 8 kind files, `crow` block, `questions.json` and validation event schemas, one example per gate shape, a full two-gate example mission, the brief template, and a README with the gate table, verbs, mission dir, versions, validation record, outbox caps, end-of-run rules and sources.

Grooming findings:
- Five gate shapes are in use, so the spec models gates as an ordered list: `[]`, `[report]`, `[plan]`, `[plan, report]`, `[mockup]`.
- One version counter per mission so `answers.vN.md` always answers version N.
- The validation record matches req_030 Need 5 and item_099: one hash-chained record per publish action.
- cdx `--continue-run` pins the parent's account, model and effort and cannot be combined with `--failover` (cdx-manager README, v0.26.6, checked 2026-11-06). Need 10's "resumes on the next account" needs a cold run. See D2.
- Backend language is still open (req_031). Only the spec validator (item_104) is Ready: it is a repo test, useful whatever the backend. The parser waits on D6.
- Stop hooks carry `stop_hook_active`, so the stop guard (item_110) needs no retry counter of its own (hooks docs, Claude Code 2.1.319, checked 2026-11-07).

# Decisions (open, for the owner)
- **D1. Where kind files live.** req_028 slice 1 says a private `missions.yaml`; the spec uses one tracked file per kind.
  - A. Tracked `kinds/<kind>.yaml` with generic defaults, plus optional private overrides. Shareable, no names in tracked files.
  - B. One private `missions.yaml` as req_028 says. Simple, but the catalogue is not publishable.
  - C. Tracked files only, no overrides. Simplest; private tuning goes through scope routing only.
  - Recommendation: A, and update req_028 slice 1.
- **D2. Resume after a usage limit on another account.** `--continue-run` keeps the parent's account.
  - A. Cold run on the next allowed account from `brief.md` and the version files. Works now; loses the session cache.
  - B. Wait for the same account to reset, then continue the session. Cheapest in tokens; can stall for hours.
  - C. A then B: cold run when the reset is more than an hour away, otherwise wait.
  - Recommendation: C. item_100 already resumes isolated runs cold.
- **D3. What `go` does for the plan kind.** The plan file is the content; there is no report gate.
  - A. `go` commits the docs on a mission branch; landing follows the scope rules (work needs a push OK). Safe and traceable.
  - B. `go` commits straight to the default branch on personal and crow repos. Fastest; no branch to clean up.
  - C. Add a report gate: `go` writes the docs, `ok` lands them. Owner sees the final files; one more round trip.
  - Recommendation: A.
- **D4. What design `approve` does before the build kind exists.** Build is phase 3; design is phase 1.
  - A. Commit the pages to `docs/design/` and open a plan mission instead. Keeps the chain moving.
  - B. Commit the pages only and leave a note. Simplest.
  - Recommendation: A, switch to the build mission in phase 3.
- **D5. Skills for research, design, repair and a generic write.**
  - A. New skills for research and design in phase 1, repair and write in their phases. Same pattern as other kinds.
  - B. Instructions in the brief only, no skill. Fewer files; behaviour drifts between runs.
  - Recommendation: A.
- **D6. Where the contract code lives before the backend is picked.**
  - A. Python in `bin/` now, shared by crow dispatch and the tests; core ports it if req_031 picks Go. Usable for AFK runs today.
  - B. Wait for req_031. No throwaway risk; nothing reads the contract meanwhile.
  - Recommendation: A for the parser (item_106) so AFK reports can use the contract now, B for the rest. The spec validator (item_104) is a test and does not wait on D6.

# Acceptance criteria
- AC1: A kind file with every column of need 1 filled for the 8 kinds, phase 1 kinds first, reviewed by the owner.
- AC2: A JSON schema for the `crow` block and for `outbox/questions.json`, with one valid example per gate shape.
- AC3: A brief template core renders, including owner prefs and the reply contract.
- AC4: A dry run of one two-gate mission and one mockup loop with files only (no core yet) shows the contract is enough to resume.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`

# References
- `logics/specs/req_037_run_contract/README.md`
- `logics/request/req_030_spike_mission_safety_and_sandbox.md`
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`
- `logics/product/prod_005_crow_jarvis_layer.md`
- `logics/architecture/adr_001_core_missions_and_inbox_architecture.md`

# Backlog
- `item_104_load_and_validate_mission_kind_files`
- `item_105_owner_review_of_the_mission_kind_catalogue`
- `item_106_parse_and_validate_the_crow_block_and_questions_file`
- `item_107_render_the_mission_brief_from_kind_record_and_answers`
- `item_108_mission_versions_responses_and_validation_record`
- `item_109_outbox_ingest_with_caps_and_parked_questions`
- `item_110_stop_guard_for_runs_that_end_without_a_valid_version_file`
- `item_111_files_only_dry_run_of_a_two_gate_mission_and_a_mockup_loop`
- `item_112_skill_and_reply_contract_for_each_mission_kind`
- `item_113_preflight_questions_per_mission_kind_at_dispatch`

# Execution order
| Order | Item | Waits on | Status |
| --- | --- | --- | --- |
| 1 | item_104 spec validator and spec test | nothing | Ready |
| 2 | item_105 owner review of the kind files | owner review, D1 | Draft |
| 3 | item_106 crow block and questions parser | item_104, D6 | Draft |
| 4 | item_112 skill per kind | D5, item_105 | Draft |
| 5 | item_107 brief rendering | item_106, core mission record, D6 | Draft |
| 6 | item_108 versions and validation record | item_106, core event log, D6 | Draft |
| 7 | item_109 outbox ingest | item_106, core mission record | Draft |
| 8 | item_111 files-only dry run | item_106, item_107, D2 | Draft |
| 9 | item_110 stop guard | item_106, headless hook loading (sandbox items) | Draft |
| 10 | item_113 preflight at dispatch | item_105, item_107 | Draft |

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
