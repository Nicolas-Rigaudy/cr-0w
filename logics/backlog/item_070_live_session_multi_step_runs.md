## item_070_live_session_multi_step_runs - One live session per multi-step run
> From version: 1.0.0
> Schema version: 1.0
> Status: Blocked
> Understanding: 75%
> Confidence: 65%
> Progress: 0%
> Complexity: High
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Multi-step headless work (crow-review-run first) runs in one `--input-format stream-json` process, sends the skill once, and feeds each step only the diff (levers 2, 5, 13).
- Keywords: stream-json, resume, cache rewrite, skill once, diff only, review run
- Use when: a headless job has several steps on the same context.
- Skip when: single-shot runs.

# Problem
- `claude -p --resume` per step rewrites the whole cache each time. The colleague measured -73% on one review ($6.39 to $1.56). Each `/skill` call re-injects about 10k tokens.

# Scope
- In:
  - A small driver in `bin/` that keeps one `claude --input-format stream-json --output-format stream-json` process and writes one message per step.
  - Skill text sent once at start; later steps are one line; re-sent only after a compaction event.
  - Steps that need new code state get the diff computed by script (`git diff` since the last step), not a re-read.
  - First user: `crow-review-run`. Then any other multi-step kind.
- Out:
  - Interactive sessions.

# Acceptance criteria
- AC1: The same review run before and after shows the cache_write share drop (item_072 method), numbers in the item.
- AC2: The skill text appears once in the transcript unless a compaction happened.

# AC Traceability
- All ACs -> this item. Proof: to fill at build.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: req_017_token_levers_quota_switch
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Biggest measured gain, biggest change. After the baseline (072) and flags (069).

# Notes
- Read how `crow-review-run` chains steps today before designing.
- Parked 2026-11-03 (owner decision). Measured before building: in 60 days of cdx runs (186), one run was a resume (`crow-review-run apply` on lambda MR 214, 17 min after its review). It wrote 12k cache tokens and read 421k (97% hits, $0.25). cdx already uses the 1h cache, so `--continue-run` does not rewrite the whole cache. The two review steps are split by the owner's decisions, which a live process cannot bridge past the cache TTL. `cdx run` has no stream-json passthrough, so the driver would run claude outside cdx.
- Reopen when resumed runs become frequent or their cache-write share passes 20% (item_072 method). The MR gate (code-review, review self, mr back to back) is the likely first user then.
