## item_027_free_account_choice_for_crow_open_term_and_handoff - Free-account choice for crow open, term and handoff
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub orchestration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: crow open, term and handoff skip accounts with a live session (cdx sessions --json), rank the rest with cdx select, ask for ask-first accounts or when none is free.
- Keywords: free, account, choice, crow, open, term, handoff
- Use when: changing account choice in bin/crow or crow-accounts.
- Skip when: relaunching limit-stopped sessions (crow-relaunch).

# Problem
- crow always launches on the default account and hits the 'session already running' prompt.

# Scope
- In:
  - pick among allowed accounts without a live session (`cdx sessions --json`), ranked by `cdx select`
  - ask-first accounts and an all-busy case ask the owner
  - `--on-busy fail` for scripted launches
  - pass `--task/--end-state` to cdx in crow handoff
- Out:
  - forcing a second session on a busy profile

# Acceptance criteria
- AC1: with the default account busy, crow handoff proposes the next allowed free account.
- AC2: tests cover busy, ask-first and none-free cases.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC1: with the default account busy, crow handoff proposes the next allowed free account.

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
