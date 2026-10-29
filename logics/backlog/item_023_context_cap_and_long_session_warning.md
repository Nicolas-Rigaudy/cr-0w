## item_023_context_cap_and_long_session_warning - Context cap and long-session warning
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:12:32

# AI Context
- Summary: Opus turns over 150k context were 52% of 14-day spend; set an auto-compact cap in the cdx profiles and flag long sessions in crow-watch and the brief.
- Keywords: context, cap, long, session, warning
- Use when: changing compaction settings or the long-session flag in crow-usage.
- Skip when: work on retros, names or accounts.

# Problem
- Opus turns above 150k context are half the spend; nothing warns before a session gets there.

# Scope
- In:
  - auto-compact window at 150k in every allowed cdx profile (owner decision 2026-10-29)
  - early warning at 120k: crow-usage flags the session, crow-watch and the brief show it
  - scoping help at the warning: propose `/wrap`, then split the rest into smaller sessions (crow handoff with one narrow end state each)
  - runbook: when to wrap, how to scope a session to one end state
- Out:
  - per-turn model switching

# Acceptance criteria
- AC1: the 150k cap is set in every allowed profile and documented.
- AC2: crow-watch and the brief name any live session over 120k.
- AC3: past the warning, the hub proposes a wrap and a split into scoped handoffs.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: AC1: the chosen cap is set in every allowed profile and documented.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)
- Request: `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`
- Primary task(s): `task_003_deliver_the_crow_cost_and_reliability_pass`

# Priority
- Priority: High
- Rationale: Set by scaffold input or defaulted for grooming.
