## item_095_scope_aware_account_pick - Scope-aware account pick with per-scope ranks and room thresholds
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 85%
> Progress: 100%
> Complexity: Medium
> Theme: General
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-07 11:21:49

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
- request-AC0 -> AC0. Proof: tests/test_crow_scope_pick.py::test_no_scope_matches_the_recorded_cdx_select against tests/fixtures/cdx-select.json (recorded `cdx select`/`status`/`configs`), plus the priority and default-vs-priority tests.
- request-AC1 -> AC1. Proof: test_rank_orders_the_scope, test_rank_never_beats_usability, test_rank_beats_priority_and_room_and_names_the_factor, test_equal_rank_falls_back_to_the_cdx_order.
- request-AC2 -> AC2. Proof: test_account_without_a_rank_is_never_picked_for_the_scope; CLI: test_an_explicit_account_beats_the_scope (CROW_ACCOUNT overrides).
- request-AC3 -> AC3. Proof: test_min_room_skips_with_a_reason, test_min_room_lets_the_account_through_above_the_floor.
- request-AC4 -> AC4. Proof: test_ask_first_stays_ask_first_inside_a_scope, test_caps_still_block_inside_a_scope; existing fallback tests pass.
- request-AC5 -> AC5. Proof: tests/test_crow_launch_pick.py scope tests (personal project, work project and overflow, hub passes `hub`, hub keeps the default without scopes); dispatch: test_dispatch_routes_by_the_project_scope; rescue uses the stranded sessions' shared scope (crow-relaunch scope_of).
- request-AC6 -> AC6. Proof: all 294 earlier tests unchanged and green (317 with the new ones); test_no_scopes_anywhere_ignores_the_scope.
- request-AC7 -> AC7. Proof: test_json_and_doctor_show_priority_and_scopes; docs/accounts.example.yaml documents scopes, scope_min_room and cdx priority; live `crow doctor` shows them for the owner config.
- request-AC8 -> AC8. Proof: test_dispatch_passes_failover (run_flags already carried `--failover`); pick prints "decided by <factor>" (tests/test_crow_launch_pick.py asserts "decided by scope rank").

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
- Owner decision 2026-11-07: the 40% room floor of the work account on personal scope stays.
- cdx priorities already set 2026-11-06: work-priority account 80, personal-priority 60, shared 20 (names in local/context.md). They only rank inside a usability class.

# Executions
| When | Runner | Run id | Title | State | PR | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| 2026-11-06T19:47:47.173034+00:00 | afk | afk-2026-11-06-scope | scope-aware account pick | succeeded |  | commit fdf998b |
