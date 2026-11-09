## item_152_turn_on_full_auto_for_the_first_mission_kinds - Turn on full auto for the first mission kinds
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Set the first kinds to full auto in their kind files once item_151 holds: research reports and briefs first; review fixes matching ledger patterns up to a local commit, the retro batch and mail triage (drafts only) when they run as missions.
- Keywords: autonomy, full auto, research, brief, review-fix, retro, mail triage, kind files
- Use when: deciding which kinds may run without asking.
- Skip when: the dial mechanics (item_151).

# Problem
- req_028 slice 5 names the first full-auto candidates. Some are not mission kinds yet (retro batch, mail triage run as scripts and skills today; review-fix is a phase 2 kind).

# Scope
- In:
  - research: auto with an FYI report (already the prod_005 default; confirm it in the kind file).
  - Brief as a `write` mission started by a schedule: auto, FYI item with the brief.
  - For each later candidate, the condition written in its kind file and left at notify until its kind runs as a mission: review-fix (every thread matches a ledger pattern, stop at a local commit), retro batch, mail triage (drafts only, never send).
  - One week of real runs per enabled kind, counted in the notes (done, FYI read, owner overrides).
- Out:
  - Moving retro or mail triage into missions (their own items), personal auto-merge (item_101).

# Acceptance criteria
- AC1: research and brief kind files say auto; a real research mission and a scheduled brief finish with only an FYI item each.
- AC2: The later candidates carry their auto condition and stay at notify.
- AC3: Notes record one week of runs and any owner override.

# AC Traceability
- request-AC3 -> AC1, AC2. Proof: kind files and run log.

# Decision framing
- Product framing: Not needed (prod_005 kinds table)
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_005_crow_jarvis_layer`
- Architecture decision(s): `adr_001_core_missions_and_inbox_architecture`
- Request: `req_028_jarvis_layer_missions_inbox_butler_chat_model_router`
- Depends on: `item_151_autonomy_dial_per_kind_and_mission_under_the_hard_rules`, `item_143_mission_runs_can_write_their_version_files` (real read-only runs)
- Primary task(s): `task_017_orchestrate_jarvis_core_slices_3_to_5_api_summaries_and_autonomy`

# Priority
- Priority: Medium
- Rationale: Pays off the dial; small once item_151 is done.

# Notes
- Ask the owner before switching each kind to auto (decision via prompt).
- Groomed 2026-11-09 from req_028 slice 5.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
