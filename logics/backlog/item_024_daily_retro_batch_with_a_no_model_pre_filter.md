## item_024_daily_retro_batch_with_a_no_model_pre_filter - Daily retro batch with a no-model pre-filter
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Medium
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:12:33

# AI Context
- Summary: Replace the per-session SessionEnd retro with one daily run over sessions a no-model pre-filter keeps (correction words, interrupts, rejected tools, repeated errors).
- Keywords: daily, retro, batch, model, pre, filter
- Use when: changing bin/crow-retro, the retro skill or the retro schedule.
- Skip when: a by-hand retro of one session (crow-retro add stays).

# Problem
- One Sonnet retro per session (53 in 14 days); a usage limit drops the retro.

# Scope
- In:
  - crow-retro scan: day's transcripts, skip retros/probes/one-shots and sessions under 3 prompts, keep sessions with correction signals
  - one `cdx run --kind retro --failover --power low` over all kept digests
  - one dated report with a section per session
  - timer Monday to Friday evening (systemd user unit or cron), none on weekends, plus `crow-retro run` by hand; Monday's run also covers weekend sessions
  - remove the SessionEnd hook, update runbooks/session-retros.md and the brief
- Out:
  - cdx's built-in retro switch

# Acceptance criteria
- AC1: pre-filter is deterministic and tested on fixture transcripts.
- AC2: one run per day at most, report lists skipped sessions with the reason.
- AC3: profile settings have no SessionEnd retro hook.

# AC Traceability
- request-AC1 -> This backlog slice. Proof: AC1: pre-filter is deterministic and tested on fixture transcripts.

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
