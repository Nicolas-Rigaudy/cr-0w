## item_026_deterministic_session_names - Deterministic session names
> From version: 1.0.0
> Schema version: 1.0
> Status: Obsolete
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:14:20

# AI Context
- Summary: Names are `<project> · <kind> · <slug>`, computed at launch or from the transcript; the Haiku titler in crow-usage goes away.
- Keywords: deterministic, session, names
- Use when: changing session titles in crow-usage, crow-runs or crow-watch.
- Skip when: renaming inside Claude Code (/rename keeps priority).

# Problem
- Titles come from four sources and a small model; they change between refreshes.

# Scope
- In:
  - scheme `<project> · <kind> · <slug>` computed at launch by crow and from transcripts for older sessions
  - remove CROW_TITLE_CMD and the titler worker
  - /rename still wins
  - wishlist entry: per-launch name passthrough in cdx
- Out:
  - renaming sessions inside Claude Code

# Acceptance criteria
- AC1: same session gives the same name on every run of crow-usage.
- AC2: no `claude -p` call left in crow-usage.

# AC Traceability
- request-AC3 -> This backlog slice. Proof: AC1: same session gives the same name on every run of crow-usage.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)
- Request: `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`
- Primary task(s): `task_003_deliver_the_crow_cost_and_reliability_pass`
- Superseded by: `item_029_trim_hub_start_context_and_fix_per_model_usage`

# Priority
- Priority: Medium
- Rationale: Set by scaffold input or defaulted for grooming.
