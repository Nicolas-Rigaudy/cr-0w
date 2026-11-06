## req_025_account_routing_by_scope - Route crow sessions to accounts by scope (work, personal, hub)
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 85%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: `crow-accounts pick` ranks every `use` account the same and ignores cdx priority. Adopt the cdx ranking (usability class, then `--priority`, then room) as crow's base order, and add a per-scope layer (work, personal, hub) on top, so work goes to the work account first, personal to the personal account first, and each overflows by its own rules.
- Keywords: accounts, pick, scope, priority, cdx priority, usability, overflow, fallback, quota, routing
- Use when: changing how crow open/term/handoff/dispatch choose an account.
- Skip when: quota display only (crow usage, deck Pulse).

# Needs
- The owner runs two own accounts: one for work first, one for personal first. Shared accounts keep their current rules (ask-first, auto fallback with caps).
- Work sessions: the work-priority account, then overflow to the personal account, then the shared fallback under its caps.
- Personal sessions: the personal account only by default. The work account takes personal work only while it has clear room left (threshold in config).
- Hub sessions (crow itself, hybrid): either own account, whichever has more room, work-priority first on a tie.
- The scope comes from the project registry (`scope` in projects.yaml); the hub is `hub`. `--scope` overrides it.

# Context
- Today: `crow_lib.pick_accounts` sorts by policy, auto_ok, then room. `launch_pick` keeps the default while it has room. No scope notion, no cdx priority.
- cdx 0.26 (docs/cli.md "How cdx picks a session"): `cdx select`, `cdx next`, `cdx run --provider`, run failover and `cdx ready` share one ranking: filters (enabled, authenticated, provider, effort floor, readiness, policy `use`), then usability class (usable now > blocked with known reset > reset known > unknown), then `--priority` 0..100 (highest first), then credits, availability, reset, effort, name. Priority never crosses a usability class. `cdx select --json` names the `deciding_factor`.
- cdx priority is one number per account, so it can only hold the default order. Set 2026-11-06: work-priority account 80, personal-priority 60, shared 20 (names in local/context.md) (`cdx set <name> --priority`). That is the work/hub order and what cdx-native picks (`cdx next`, `cdx run --failover`, `cdx rescue`) use.
- crow already reads cdx `policy` from `cdx configs --json` with accounts.yaml as fallback; `priority` sits in the same payload (absent = default).
- `local/accounts.yaml` already holds policy, memory, fallback caps per account.
- Proposed config (local/accounts.yaml): per account `scopes: {work: 1, personal: 2, hub: 1}` (rank, lower first; absent = not used for that scope) and an optional `scope_min_room: {personal: 40}` (only take that scope while at least N% of 5h and week is left).
- Callers: `bin/crow` pick_account (open, term, handoff), `crow dispatch`, `crow runs relaunch`, `crow rescue`, afk preflight.

- Layering: cdx owns the global default (policy + priority); crow owns only what cdx cannot express (per-scope rank, scope room threshold, shared-account caps). No crow-side duplicate of the cdx priority number.

# Acceptance criteria
- AC0: Without `--scope`, crow's order matches the cdx ranking: usability class, then cdx `--priority`, then room. A test compares crow's pick to `cdx select --json` on the same status fixture.
- AC1: `crow-accounts pick --scope work|personal|hub` returns the first account by scope rank that has room; equal ranks fall back to the cdx order (AC0). A scope rank never beats usability: an account out of quota never wins over a usable one.
- AC2: An account without the scope in `scopes` is never picked for it automatically; CROW_ACCOUNT still forces any allowed account.
- AC3: `scope_min_room` skips the account for that scope below the threshold, with the reason in the skipped list.
- AC4: Ask-first and `fallback: auto` caps behave as today inside each scope.
- AC5: crow open/term/handoff/dispatch derive the scope from the registry; the hub uses `hub`.
- AC6: An accounts.yaml without `scopes` behaves exactly as today (tests prove it).
- AC7: `crow-accounts json` and `crow doctor` show each account's cdx priority and scopes; docs/accounts.example.yaml documents the keys.
- AC8: Headless dispatch passes `--failover` to `cdx run` so a run that hits a limit moves to the next `use` account by cdx order, and the pick result names the deciding factor (scope rank, usability, priority, room) like `cdx select`.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Out of scope
- Mid-session account switching; changes inside cdx-manager (scope-aware priority belongs upstream only if crow's layer proves it; then via `crow handoff cdx-manager`).

# Risks
- Personal sessions starve when the personal account is out and the work threshold is high: pick must say why and suggest CROW_ACCOUNT.

# Companion docs
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)

# References
- `bin/crow-accounts`
- cdx-manager docs/cli.md (How cdx picks a session, Programmatic callers), docs/sessions-and-notifications.md (Persistent launch settings)
- `bin/crow_lib.py` (pick_accounts, launch_pick)
- `bin/crow` (pick_account)
- `tests/test_crow_accounts.py`

# Backlog
- `item_095_scope_aware_account_pick`

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
