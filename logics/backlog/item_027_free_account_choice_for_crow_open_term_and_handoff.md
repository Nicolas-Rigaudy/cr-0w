## item_027_free_account_choice_for_crow_open_term_and_handoff - Free-account choice for crow open, term and handoff
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Low
> Theme: Hub orchestration
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Owner: crow
> Indicators reviewed: 2026-11-03 18:29:24

# AI Context
- Summary: crow open, term and handoff skip accounts with a live session (cdx sessions --json) or no quota left (cdx status), rank the rest by room, ask for ask-first accounts or when none is free.
- Keywords: free, account, choice, crow, open, term, handoff
- Use when: changing account choice in bin/crow or crow-accounts.
- Skip when: relaunching limit-stopped sessions (crow-relaunch).

# Problem
- crow always launches on the default account and hits the 'session already running' prompt.
- crow also launches on a default account with no 5h quota left: on 2026-11-03 `crow term` opened two sessions on the default account at 0% 5h.

# Scope
- In:
  - pick among allowed accounts without a live session (`cdx sessions --json`), ranked by room left (`cdx status`, as `crow_lib.pick_accounts` and `cdx select` do)
  - skip allowed accounts with no 5h or weekly quota left, or past their crow cap (`cdx status` rows)
  - an explicit `CROW_ACCOUNT=<name>` still wins
  - ask-first accounts and an all-busy case ask the owner
  - `--on-busy fail` for scripted launches
  - pass `--task/--end-state` to cdx in crow handoff
- Out:
  - forcing a second session on a busy profile

# Acceptance criteria
- AC1: with the default account busy, crow handoff proposes the next allowed free account.
- AC2: tests cover busy, ask-first and none-free cases.
- AC3: with the default account at 0% 5h, crow term/open/handoff use the allowed account with the most room, or refuse when none has room.

# AC Traceability
- request-AC4 -> This backlog slice. Proof: AC1: with the default account busy, crow handoff proposes the next allowed free account.
- request-AC4 -> This backlog slice. Proof: AC3: with the default account at 0% 5h, crow term/open/handoff use the allowed account with the most room, or refuse when none has room.

# Validation
- Implemented in f85ab46: `crow-accounts pick` skips accounts with a live session or no quota left; crow open, term and handoff use it; handoff passes `--task/--end-state` to cdx.
- `python3 -m pytest -q`: 227 passed. `tests/test_crow_launch_pick.py` covers busy, ask-first, none-free, default at 0% 5h, explicit and inherited `CROW_ACCOUNT`.
- Live check: `bin/crow-accounts pick --on-busy fail` with the default at 0% 5h and the other allowed account busy refuses with both reasons.

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
