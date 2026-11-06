## req_025_account_routing_by_scope - Route crow sessions to accounts by scope (work, personal, hub)
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: `crow-accounts pick` ranks every `use` account the same. Give each account a priority per scope (work, personal, hub) so work goes to the work account first, personal to the personal account first, and each overflows by its own rules.
- Keywords: accounts, pick, scope, priority, overflow, fallback, quota, routing
- Use when: changing how crow open/term/handoff/dispatch choose an account.
- Skip when: quota display only (crow usage, deck Pulse).

# Needs
- The owner runs two own accounts: one for work first, one for personal first. Shared accounts keep their current rules (ask-first, auto fallback with caps).
- Work sessions: the work-priority account, then overflow to the personal account, then the shared fallback under its caps.
- Personal sessions: the personal account only by default. The work account takes personal work only while it has clear room left (threshold in config).
- Hub sessions (crow itself, hybrid): either own account, whichever has more room, work-priority first on a tie.
- The scope comes from the project registry (`scope` in projects.yaml); the hub is `hub`. `--scope` overrides it.

# Context
- Today: `crow_lib.pick_accounts` sorts by policy, auto_ok, then room. `launch_pick` keeps the default while it has room. No scope notion.
- `local/accounts.yaml` already holds policy, memory, fallback caps per account.
- Proposed config (local/accounts.yaml): per account `scopes: {work: 1, personal: 2, hub: 1}` (rank, lower first; absent = not used for that scope) and an optional `scope_min_room: {personal: 40}` (only take that scope while at least N% of 5h and week is left).
- Callers: `bin/crow` pick_account (open, term, handoff), `crow dispatch`, `crow runs relaunch`, `crow rescue`, afk preflight.

# Acceptance criteria
- AC1: `crow-accounts pick --scope work|personal|hub` returns the first account by scope rank that has room; ties and equal ranks fall back to the current room ordering.
- AC2: An account without the scope in `scopes` is never picked for it automatically; CROW_ACCOUNT still forces any allowed account.
- AC3: `scope_min_room` skips the account for that scope below the threshold, with the reason in the skipped list.
- AC4: Ask-first and `fallback: auto` caps behave as today inside each scope.
- AC5: crow open/term/handoff/dispatch derive the scope from the registry; the hub uses `hub`.
- AC6: An accounts.yaml without `scopes` behaves exactly as today (tests prove it).
- AC7: `crow-accounts json` and `crow doctor` show each account's scopes; docs/accounts.example.yaml documents the keys.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Out of scope
- Mid-session account switching; cdx-side policy changes.

# Risks
- Personal sessions starve when the personal account is out and the work threshold is high: pick must say why and suggest CROW_ACCOUNT.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-accounts`
- `bin/crow_lib.py` (pick_accounts, launch_pick)
- `bin/crow` (pick_account)
- `tests/test_crow_accounts.py`

# Backlog
- `item_095_scope_aware_account_pick`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
