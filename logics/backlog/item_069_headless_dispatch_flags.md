## item_069_headless_dispatch_flags - Headless dispatch flags
> From version: 1.0.0
> Schema version: 1.0
> Status: Ready
> Understanding: 85%
> Confidence: 80%
> Progress: 0%
> Complexity: Low
> Theme: cost
> Reminder: Update status/understanding/confidence/progress and linked request/task references when you edit this doc.

# AI Context
- Summary: Every headless run (crow-runs, crow-review-run, crow-retro, titler) gets strict MCP config, a pinned model ID and effort, a budget cap, and the headless env fixes (levers 3, 4, 7, 16).
- Keywords: headless, strict-mcp-config, model id, effort, max-budget-usd, watchdog, DISABLE_BACKGROUND_TASKS
- Use when: changing how crow launches `claude -p` or `cdx run`.
- Skip when: interactive sessions.

# Problem
- Dispatch passes none of these. Each run loads every MCP connector (~44k tokens per call measured by the colleague), can drift in a loop, and sees phantom background results.

# Scope
- In:
  - Lever 3: `--strict-mcp-config` with an empty or per-kind MCP file.
  - Lever 4: exact model ID (`claude-opus-5-5`, `claude-sonnet-5-5`), effort fixed at launch, never changed during the run.
  - Lever 7: `--max-budget-usd` (or `cdx set --budget`) per run kind, plus a progress watchdog in `crow-runs` (no new tool call or file change in N minutes = kill, state `stuck`).
  - Lever 16: `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1`; prompts pass JSON through files, never inline in Bash.
  - One shared helper in `crow_lib.py` so every launcher uses it.
- Out:
  - Per-repo MCP needs beyond a per-kind allowlist.

# Acceptance criteria
- AC1: `crow-runs` dry-run prints the full command with all four flag groups for each run kind.
- AC2: A test run's transcript shows no MCP tool definitions and the pinned model ID on every turn.
- AC3: A run forced into a loop stops at the budget or the watchdog and shows `stuck` or `budget`.

# AC Traceability
- All ACs -> this item. Proof: to fill at build.

# Decision framing
- Product framing: Not needed
- Architecture framing: Not needed

# Links
- Request: req_017_token_levers_quota_switch
- Primary task(s): (none yet)

# Priority
- Priority: High
- Rationale: Cheapest change, applies to every run.

# Notes
- Check which flags `cdx run` forwards (`--extra-args`) before adding a raw `claude` path.
