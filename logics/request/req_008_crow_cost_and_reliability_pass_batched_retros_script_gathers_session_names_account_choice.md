## req_008_crow_cost_and_reliability_pass_batched_retros_script_gathers_session_names_account_choice - crow cost and reliability pass: batched retros, script gathers, session names, account choice
> From version: 1.0.0
> Schema version: 1.0
> Status: Draft
> Understanding: 90%
> Confidence: 85%
> Complexity: Medium
> Theme: Hub cost
> Reminder: Update status/understanding/confidence and linked backlog/task references when you edit this doc.

# AI Context
- Summary: Cost and reliability pass from the 2026-10-29 audit: long Opus contexts, per-session retros, gather turns, model titles, busy-account launches, unprotected headless runs.
- Keywords: crow, cost, reliability, pass, batched, retros, script, gathers, session, names, account, choice
- Use when: planning or reviewing crow token spend and automatic routines.
- Skip when: UI work (req_007) or vault structure.

# Needs
- Over 14 days, half the spend came from Opus turns whose context was already above 150k tokens (median 242k). Long sessions should be compacted or restarted before that point, and the owner should see when a session crosses it.
- Every finished session starts its own retro run. Retros should run once a day as one batch, on the sessions that show a signal worth learning from, picked without a model.
- The brief, daily and weekly skills spend most of their tokens on tool turns that only gather data. A script should gather everything it can and hand the model one compact JSON.
- Sessions get titles from four sources (account name, a small-model titler, dispatch labels, /rename). One deterministic naming scheme should replace the model titler.
- crow open, term and handoff always launch on the default account, which prompts when a session already runs there; two sessions on one profile share OAuth and can log each other out. crow should pick a free allowed account.
- Headless runs should use the cdx features that cap and protect them: spend budget, failover on a usage limit, fallback model, run kind.
- Hub sessions start with 31-50k tokens of context; the parts the hub never uses should go.
- crow-usage attributes Opus and Sonnet tokens of multi-model sessions to Haiku, so the per-model numbers are wrong.

# Context
- Audit with the measured numbers and sources: local/state/research/crow-audit.md (private).
- cdx has its own retro (`cdx set --retro on`, `cdx retro`): per session, read-only, lessons as text. crow keeps its own skill (it applies lessons) and reuses `cdx run --kind retro --failover --budget --power low`.
- `cdx select` does not skip profiles with a live interactive session; `cdx sessions --json` lists them.
- Claude Code accepts a name flag at launch; cdx has no per-launch passthrough for it yet (cdx wishlist).
- cdx is not owned by the hub's owner: crow only calls its documented CLI.

# Acceptance criteria
- AC1: no SessionEnd retro hook remains; a daily batch analyses only pre-filtered sessions in one run and writes one dated report.
- AC2: brief, daily and weekly read one JSON from a gather script; each fork needs at most the tracker/mail calls plus its answer.
- AC3: every session shown by crow has a name `<project> · <kind> · <slug>` computed without a model; the model titler is gone.
- AC4: crow open, term and handoff launch on a free allowed account, or ask when none is free; they never start a second session on a busy profile unasked.
- AC5: dispatch, AFK and retro runs pass a budget and failover; the long-session warning and context cap are live; crow-usage per-model totals match the transcripts.

# Definition of Ready (DoR)
- [x] Problem statement is explicit and user impact is clear.
- [x] Scope boundaries (in/out) are explicit.
- [x] Acceptance criteria are testable.
- [x] Dependencies and known risks are listed.

# Companion docs
- Product brief(s): `prod_003_crow_cost_and_reliability_pass`
- Architecture decision(s): (none yet)

# References
- bin/crow
- bin/crow-retro
- bin/crow-digest
- bin/crow-usage
- bin/crow-runs
- bin/crow-accounts
- bin/crow-project
- .claude/skills/brief/SKILL.md
- .claude/skills/retro/SKILL.md
- runbooks/session-retros.md
- runbooks/dispatch-background-runs.md

# Backlog
- `item_023_context_cap_and_long_session_warning`
- `item_024_daily_retro_batch_with_a_no_model_pre_filter`
- `item_025_gather_script_for_brief_daily_and_weekly`
- `item_026_deterministic_session_names`
- `item_027_free_account_choice_for_crow_open_term_and_handoff`
- `item_028_budget_failover_and_kind_on_headless_runs`
- `item_029_trim_hub_start_context_and_fix_per_model_usage`
