## item_095_scope_aware_account_pick - Scope-aware account pick with per-scope ranks and room thresholds
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
- Summary: Add per-scope ranks and room thresholds to local/accounts.yaml and make `crow-accounts pick --scope` and every caller route by them.
- Keywords: accounts, pick, scope, rank, scope_min_room, overflow
- Use when: building req_025.
- Skip when: quota display work.

# Problem
- Every allowed account competes for every session. Work and personal sessions cannot each prefer their own account and overflow differently.

# Scope
- In:
  - `scopes` (rank per scope) and `scope_min_room` keys in accounts.yaml, parsed by crow-accounts, shown in `json` and `crow doctor`.
  - `crow_lib.pick_accounts` / `launch_pick` filter and order by scope rank, then current room ordering.
  - `crow-accounts pick --scope`; scope from the registry in bin/crow (open, term, handoff), crow dispatch, relaunch, rescue; hub = `hub`.
  - docs/accounts.example.yaml and runbooks (interrupted-sessions, dispatch-background-runs) updated.
  - Owner config in local/accounts.yaml: work-priority account rank 1 work/hub, 2 personal with a room threshold; personal account rank 1 personal/hub, 2 work; shared account unchanged (work only).
- Out:
  - Personal butler features (req_026).

# Acceptance criteria
- AC1: pick --scope follows ranks, then room; equal ranks keep today's order.
- AC2: An account without a scope rank is never auto-picked for it; CROW_ACCOUNT overrides.
- AC3: scope_min_room skips with a reason in the skipped list.
- AC4: ask-first and fallback caps unchanged inside a scope.
- AC5: callers pass the registry scope; hub passes `hub`.
- AC6: no `scopes` key = today's behaviour, covered by tests.
- AC7: json, doctor and the example config show scopes.

# AC Traceability
- request-AC1 -> AC1. Proof: tests on pick_accounts with ranked fixtures.
- request-AC2 -> AC2. Proof: test with an unranked account.
- request-AC3 -> AC3. Proof: test below the threshold.
- request-AC4 -> AC4. Proof: existing fallback tests pass with scopes set.
- request-AC5 -> AC5. Proof: crow CLI test on a personal-scope project.
- request-AC6 -> AC6. Proof: unchanged existing tests.
- request-AC7 -> AC7. Proof: doctor/json tests.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): (none yet)
- Architecture decision(s): (none yet)
- Request: `req_025_account_routing_by_scope`
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: The personal account and the butler work (req_026) depend on it.

# Notes
- Start from `crow_lib.pick_accounts` and `launch_pick`; keep the sort key, prepend the scope rank.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
