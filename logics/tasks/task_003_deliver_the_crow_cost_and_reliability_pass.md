## task_003_deliver_the_crow_cost_and_reliability_pass - Deliver the crow cost and reliability pass
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Implementation delivery
> Reminder: Update status/understanding/confidence/progress and linked request/backlog references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:12:33

# AI Context
- Summary: Delivery order for req_008: account choice, context cap, run failover, retro batch, gather script, names, start context, then re-measure.
- Keywords: deliver, crow, cost, reliability, pass
- Use when: implementing any item_023..item_029 slice.
- Skip when: work outside req_008.

# Context
- Orchestrate the scaffolded request chain and keep sibling implementation slices linked.

# Plan
- [ ] 1. 1. Free-account choice (unblocks parallel sessions).
- [ ] 2. 2. Context cap and long-session warning.
- [ ] 3. 3. Failover, fallback model and kind on headless runs (optional budget).
- [ ] 4. 4. Weekday retro batch, then remove the SessionEnd hook.
- [ ] 5. 5. Gather script for brief, daily and weekly.
- [ ] 6. 6. Deterministic session names.
- [ ] 7. 7. Trim start context and fix per-model usage.
- [ ] 8. 8. Re-measure 14-day spend with the audit method and record it.
- [ ] ADR 009 checkpoint: update affected Logics docs during each meaningful wave and leave the repo commit-ready.
- [ ] Keep commit creation under operator control; do not force one commit per micro-step.
- [ ] GATE: do not close until lint, audit, and scaffold validation pass.

# Backlog
- `item_023_context_cap_and_long_session_warning`
- `item_024_daily_retro_batch_with_a_no_model_pre_filter`
- `item_025_gather_script_for_brief_daily_and_weekly`
- `item_026_deterministic_session_names`
- `item_027_free_account_choice_for_crow_open_term_and_handoff`
- `item_028_budget_failover_and_kind_on_headless_runs`
- `item_029_trim_hub_start_context_and_fix_per_model_usage`

# Definition of Done (DoD)
- [ ] Generated request, product, backlog, and task docs are present.
- [ ] Context-pack handoff is available when requested.
- [ ] Validation passes.
- [ ] Meaningful waves followed ADR 009: affected docs updated and the repo left commit-ready without automatic commits.

# AC Traceability
- request-AC5 -> `item_023_context_cap_and_long_session_warning`. Proof deferred to slice closeout.
- request-AC1 -> `item_024_daily_retro_batch_with_a_no_model_pre_filter`. Proof deferred to slice closeout.
- request-AC2 -> `item_025_gather_script_for_brief_daily_and_weekly`. Proof deferred to slice closeout.
- request-AC3 -> `item_026_deterministic_session_names`. Proof deferred to slice closeout.
- request-AC4 -> `item_027_free_account_choice_for_crow_open_term_and_handoff`. Proof deferred to slice closeout.
- request-AC5 -> `item_028_budget_failover_and_kind_on_headless_runs`. Proof deferred to slice closeout.
- request-AC5 -> `item_029_trim_hub_start_context_and_fix_per_model_usage`. Proof deferred to slice closeout.

# Validation
- (no validation recorded yet)

# Report
- Not started.

# Links
- Request: `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)
