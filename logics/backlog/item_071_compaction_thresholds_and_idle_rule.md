## item_071_compaction_thresholds_and_idle_rule - Compaction thresholds and idle rule
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
> Complexity: Low
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.
> Indicators reviewed: 2026-11-03 17:33:16

# AI Context
- Summary: crow-ctx advises compaction only when it pays (hot from ~65k, cold only past ~125k) and nudges a /compact before the cache expires on an idle session (levers 9, 10).
- Keywords: compact, cache expiry, idle, threshold, crow-ctx, idleCompaction
- Use when: changing compaction advice or the long-session rule.
- Skip when: the 150k auto-compact cap itself (done).

# Problem
- The base load (system prompt, tools, CLAUDE.md) is 20-60k and reloads in full anyway. Compacting a small session costs more than it saves (colleague: 83k compacted cold = $0.85 vs $0.68 reloaded). An idle session past the cache TTL reloads at about 2x.

# Scope
- In:
  - `crow-ctx`: per live session, "compact now" only past 65k while the cache is warm; past 125k when cold. Thresholds in env like CTX_WARN.
  - Idle rule: a live session idle ~50 min with a 1h TTL gets a notification "`/compact <what to keep>` or close". Headless: `idleCompaction: false`.
  - runbooks/long-sessions.md updated with both rules.
- Out:
  - Auto-sending /compact into a session.

# Acceptance criteria
- AC1: `crow-ctx --json` gives `compact_advice` per session that matches the thresholds on three fixture sessions.
- AC2: An idle session at 50 min triggers one notification, not more.

# AC Traceability
- AC1 -> `crow-ctx` compact_advice (`--json` rows and `crow-ctx advice --json`). Proof: tests/test_crow_ctx_advice.py::test_three_fixture_sessions (warm 80k compact, cold 80k reload, cold 130k compact) and test_five_minute_cache_goes_cold_fast.
- AC2 -> `crow-ctx idle-hook` (profile Stop hook) and `idle-watch`. Proof: test_idle_nudge_fires_once_when_still_idle (one notify-send, wait 50 min for a 1h cache), test_idle_nudge_skips_new_activity_small_sessions_and_short_ttl; live: a second hook call for the same session replaced the first watcher.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: req_017_token_levers_quota_switch
- Primary task(s): (none yet)

# Priority
- Priority: Medium
- Rationale: Advice change, small code.

# Notes
- Check the hub's cache TTL (1h here) per account before fixing 50 min.
- Built 2026-11-03: TTL is read per session (1h when it wrote 1h cache, else 5 min; `CROW_CACHE_TTL_MIN` overrides), so the nudge fires at TTL minus `CROW_CTX_IDLE_LEAD_MIN` (10) and is skipped on 5 min caches. crow-usage sessions carry `cache_1h_tokens`. Thresholds `CROW_CTX_HOT` (65k), `CROW_CTX_COLD` (125k).
- Not done: `idleCompaction: false` for headless runs (not checked whether Claude Code has this setting; headless runs do not sit idle). The Stop hook moved from the hub settings to every profile's user settings on 2026-11-03 (owner OK); it skips headless runs (CLAUDE_CODE_ENTRYPOINT sdk-cli, checked on a real cdx run).
