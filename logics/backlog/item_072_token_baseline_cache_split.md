## item_072_token_baseline_cache_split - Token baseline with cache split
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 90%
> Confidence: 85%
> Progress: 100%
> Complexity: Low
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Show per session and per run kind the cost split between cache reload (write), cache reread and output, and record a baseline before the other req_017 items (lever 1).
- Keywords: measure, cache_read, cache_write, output, baseline, crow usage
- Use when: judging whether a token change worked.
- Skip when: per-repo cost questions.

# Problem
- The colleague found reload 40%, reread 43%, output 12% on their side. Crow has the raw numbers in crow-usage but no split, so we would optimise blind.

# Scope
- In:
  - `crow-usage --json`: `cost_split {reload, reread, output, input}` per session and per kind.
  - `crow usage` text and the brief: one line with the 14-day split.
  - Baseline written to `local/state/research/token-baseline-2026-11.md` with the cost-audit method.
- Out:
  - New dashboards in the deck (later, if useful).

# Acceptance criteria
- AC1: The split for one session matches a hand recount from its transcript (dedup on message.id, requestId).
- AC2: The baseline file exists with window, price table and numbers.

# AC Traceability
- AC1 -> `bin/crow-usage` `cost_split`. Proof: session c4825f35 recounted from its raw transcript with dedup on (message.id, requestId), 273 calls: reload $4.8846, reread $12.6762, output $4.0749, input $0.0023, same as crow-usage. Unit test `tests/test_crow_usage_split.py`.
- AC2 -> `local/state/research/token-baseline-2026-11.md`. Proof: window 2026-10-21 to 2026-11-03, cdx price table, split by kind and account.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: req_017_token_levers_quota_switch
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: First: every other item is measured against it.

# Notes
- Method: runbooks/cost-audit.md.
- Built 2026-11-03. Baseline: $513 over 14 days, reload 32%, reread 46%, output 22%, input 0%. Runs: reload 44%. The brief shows `crow-usage --split`.
- The 14-day split always scans 14 days of transcripts; `sessions[]` in the JSON keeps the old 7-day (or --days) window.
