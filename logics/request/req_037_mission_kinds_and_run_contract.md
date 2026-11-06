## req_037_mission_kinds_and_run_contract - Mission kinds and run contract
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 70
> Confidence: 60
> Complexity: High
> Theme: Hub strategy
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Define every mission kind (gates, autonomy, account, effort, outputs) and the file contract between a run and core: brief, plan and report with a fenced `crow` block, outbox, versions, answer verbs. Schema first, before the first mission runs.
- Keywords: mission kind, gate, autonomy, run contract, crow block, outbox, plan, report, versions, answer verbs, preflight
- Use when: adding a mission kind, writing a dispatch prompt, parsing what a run asks.
- Skip when: choosing hosts, sandbox or UI stack (req_029 to req_031).

# Needs
1. **Kind catalogue.** For each kind: purpose, gates, default autonomy, allowed scopes, account and model, effort, time cap, outputs, publish action. Starting set from prod_005: `ship`, `review`, `review-fix`, `rebase`, `design`, `research`, `document`, `report`, `ci-repair`, `personal-code`, `crow`, `notes`, `admin`.
2. **One skill per kind.** Ultron maps each kind to a skill and a key format (ticket key, MR iid, free slug). crow already has skills for most kinds (`review`, `review-threads`, `mr`, `daily`, `weekly`, `mgmt-weekly`). The kind names its skill; the brief adds the run contract.
3. **Preflight.** Before dispatch, the kind's preflight asks only what the run cannot decide itself (Jarvis). If the owner skips it, the run states its safest assumption in the plan and goes on.
4. **Gate shapes.**
   - Two-gate (Ultron): run writes `plan.md`, ends turn; owner answers `go`, `dig` or `no`; run writes `report.md`, ends turn; owner answers `ok` or `msg`; core publishes.
   - Mockup loop (Jarvis): run writes pages to `outbox/`; owner answers `approve` (commits pages to `docs/design/` and opens the build mission), `change` with a note (up to 5 rounds), or `abandon`. Approval only from the nest page.
   - Straight through: no gate; FYI at the end; personal repos auto-merge on green CI and a secret scan.
5. **The `crow` block.** One fenced block per plan or report, parsed by core. Fields: `kind`, `phase` (plan, report, mockup), `summary` (one sentence), `confidence` (0 to 100), `verdict` (ready, blocked, needs-input), `questions` (id, text, choices, default), `artifacts` (paths), `publish` (what core would do: push, MR, post, with exact titles and text). Markdown around it stays for humans.
6. **Versions and responses.** Each answered version is kept with the owner's words. The next version starts with a "Responses" section saying what became of each point (Ultron "Réponses"). nest shows the history.
7. **Validation record.** `ok` writes a validation event tied to the report version. Any later message drops it. core checks it before every publish.
8. **Outbox.** `outbox/questions.json` for questions outside a plan or report; artifacts as files (HTML, PNG, MP4, diff). Size caps per file and per mission.
9. **Brief.** `brief.md` carries goal, kind, scope, end state, owner prefs, prior answers, constraints (time cap, no publish), and the reply contract. Built by core, never by hand.
10. **End-of-run rules.** A run ends its turn after writing; it never waits. A run that hits a usage limit ends; core resumes on the next account.

# Context
- Owner decisions 2026-11-06: gates per kind; preflight plus park and resume; personal repos auto-merge; work hard rules unchanged.
- Ultron lesson 2: define the plan and report schema first. Free text in a hook is hard to parse.
- crow already has: skills per task type, `crow dispatch` prompts, `local/owner-prefs.md`, review ledgers.
- adr_001 sets the mission dir and the event log; this request fills in the formats.

# Acceptance criteria
- AC1: A kind table with every column of need 1 filled for the starting kinds, reviewed by the owner.
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
- `logics/request/req_028_jarvis_layer_missions_inbox_butler_chat_model_router.md`
- `logics/product/prod_005_crow_jarvis_layer.md`
- `logics/architecture/adr_001_core_missions_and_inbox_architecture.md`

# Backlog
- none

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
