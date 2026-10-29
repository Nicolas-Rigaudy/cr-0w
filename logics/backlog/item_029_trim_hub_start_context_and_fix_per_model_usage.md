## item_029_trim_hub_start_context_and_fix_per_model_usage - Trim hub start context and fix per-model usage
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 90%
> Confidence: 85%
> Progress: 0%
> Complexity: Low
> Theme: Hub cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-10-29 17:16:27

# AI Context
- Summary: Measure the 31-50k start context of hub sessions and drop unused parts; fix crow-usage's per-model split that shows Haiku spend absent from transcripts.
- Keywords: trim, hub, start, context, fix, per, model, usage
- Use when: trimming hub settings/plugins or touching the by_model code in crow-usage.
- Skip when: project CLAUDE.md content.

# Problem
- Hub sessions start at 31-50k tokens; crow-usage shows Haiku spend that the transcripts do not contain.

# Scope
- In:
  - `bin/crow-cost-audit --days N --json`: the runbooks/cost-audit.md method as a script (dedup, context buckets, per-project, per-skill), used for the baseline and the re-measure
  - measure the start context by part, drop plugins/skills/MCP servers the hub never uses
  - fix the per-model split in crow-usage with a test
- Out:
  - changes to project CLAUDE.md content

# Acceptance criteria
- AC1: start context measured before and after in the runbook.
- AC2: per-model totals match a transcript recount within 2%.

# AC Traceability
- request-AC5 -> This backlog slice. Proof: AC1: start context measured before and after in the runbook.
- request-AC3 -> This backlog slice. Proof: AC2: per-model totals match a transcript recount within 2%.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)
- Request: `req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice`
- Primary task(s): `task_003_deliver_the_crow_cost_and_reliability_pass`

# Priority
- Priority: Low
- Rationale: Set by scaffold input or defaulted for grooming.
