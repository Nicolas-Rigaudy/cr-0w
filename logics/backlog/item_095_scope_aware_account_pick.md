## item_095_scope_aware_account_pick - Scope-aware account pick with per-scope ranks and room thresholds
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 95%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Align crow's pick with the cdx ranking (usability, then cdx `--priority`, then room), then add per-scope ranks and room thresholds from local/accounts.yaml so `crow-accounts pick --scope` and every caller route by them.
- Keywords: accounts, pick, scope, rank, scope_min_room, overflow
- Use when: building req_025.
- Skip when: quota display work.

# Problem
- Every allowed account competes for every session. Work and personal sessions cannot each prefer their own account and overflow differently.

# Scope
- In:
  - Read `priority` from `cdx configs --json` next to `policy`; base order = usability class, cdx priority, room (same as `cdx select`).
  - `scopes` (rank per scope) and `scope_min_room` keys in accounts.yaml, parsed by crow-accounts, shown in `json` and `crow doctor`.
  - `crow_lib.pick_accounts` / `launch_pick` filter by scope, order by usability class, then scope rank, then the cdx base order; return the deciding factor.
  - `crow-accounts pick --scope`; scope from the registry in bin/crow (open, term, handoff), crow dispatch (adds `cdx run --failover`), relaunch, rescue; hub = `hub`.
  - docs/accounts.example.yaml and runbooks (interrupted-sessions, dispatch-background-runs) updated.
  - Owner config in local/accounts.yaml: work-priority account rank 1 work/hub, 2 personal with a room threshold; personal account rank 1 personal/hub, 2 work; shared account unchanged (work only).
- Out:
  - Personal butler features (req_026).

# Acceptance criteria
- AC0: no --scope = cdx order (usability, priority, room); matches `cdx select` on a fixture.
- AC1: pick --scope follows usability, then ranks, then the cdx order; a rank never beats usability.
- AC2: An account without a scope rank is never auto-picked for it; CROW_ACCOUNT overrides.
- AC3: scope_min_room skips with a reason in the skipped list.
- AC4: ask-first and fallback caps unchanged inside a scope.
- AC5: callers pass the registry scope; hub passes `hub`.
- AC6: no `scopes` key = today's behaviour, covered by tests.
- AC7: json, doctor and the example config show cdx priority and scopes.
- AC8: dispatch passes --failover; pick prints the deciding factor.

# AC Traceability
- request-AC0 -> AC0. Proof: pick_accounts test against a recorded `cdx select --json` fixture.
- request-AC1 -> AC1. Proof: tests on pick_accounts with ranked fixtures.
- request-AC2 -> AC2. Proof: test with an unranked account.
- request-AC3 -> AC3. Proof: test below the threshold.
- request-AC4 -> AC4. Proof: existing fallback tests pass with scopes set.
- request-AC5 -> AC5. Proof: crow CLI test on a personal-scope project.
- request-AC6 -> AC6. Proof: unchanged existing tests.
- request-AC7 -> AC7. Proof: doctor/json tests.
- request-AC8 -> AC8. Proof: dispatch command-line test; pick output test.

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
- Start from `crow_lib.pick_accounts` and `launch_pick`. New key: (policy != use, usability class, scope rank, not auto_ok, -cdx priority, -room).
- cdx priorities already set 2026-11-06: work-priority account 80, personal-priority 60, shared 20 (names in local/context.md). They only rank inside a usability class.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
