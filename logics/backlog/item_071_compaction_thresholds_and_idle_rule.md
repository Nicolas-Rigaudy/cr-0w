## item_071_compaction_thresholds_and_idle_rule - Compaction thresholds and idle rule
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 75%
> Progress: 0%
> Complexity: Low
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

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
- All ACs -> this item. Proof: to fill at build.

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
