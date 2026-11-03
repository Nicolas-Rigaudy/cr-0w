## item_069_headless_dispatch_flags - Headless dispatch flags
> From version: 1.0.0
> Schema version: 1.0
> Status: Done
> Understanding: 95%
> Confidence: 90%
> Progress: 100%
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
- AC1 -> `crow-runs dispatch --dry-run`. Proof: dry runs for task, review, retro, oneshot and a `--continue-run` each print `--model claude-*-5-5 --power --timeout-seconds --budget` and the env pair; tests/test_crow_headless_flags.py::test_every_kind_pins_model_id_effort_budget_timeout.
- AC2 -> `crow_lib.HEADLESS_ENV`, `crow_lib.model_id`. Proof: live haiku runs on 2026-11-03 report model `claude-haiku-5-5` and no `mcp__` tool from the hub. Partial: in a work repo with a `.mcp.json`, its's `.mcp.json` servers (code graph, cloud servers) still load; see Notes.
- AC3 -> `cdx run --budget`, `crow-runs watchdog`, `run_state`. Proof: live run with a 0.004 USD budget stopped with `budget_exceeded` (exit 125, run dd25de73) and `run_state` read `budget`; test_watchdog_kills_idle_run_and_spares_young_one (idle 20 min = SIGTERM, state `stuck`).

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
- Built 2026-11-03: `crow_lib.run_flags(kind, provider, model, power, continue_run)` and `RUN_KINDS` (task 25 USD/120 min, review 10/60, retro 2/20, oneshot 3/20; env `CROW_<KIND>_BUDGET`, `_TIMEOUT_MIN`). Used by crow-runs dispatch (`--kind`, `--dry-run`), crow-review-run (kind review), crow-retro, `crow` one-shots. Titler pinned to `claude-haiku-5-5` (it already had `--strict-mcp-config`). Watchdog in the run notifier, `CROW_STUCK_MIN` (20), reads main and subagent transcripts. New states `budget` and `stuck` in crow-runs and crow-watch.
- Budget is Claude only: Codex budgets fail as `budget_unverifiable` without model evidence.
- Lever 3 gap: `cdx run` has no per-run extra args, so `--strict-mcp-config` cannot be passed. `ENABLE_CLAUDEAI_MCP_SERVERS=false` removes claude.ai connectors only; a repo's `.mcp.json` servers still load. `cdx set --extra-args` would hit interactive sessions too. Needs a cdx-manager change (per-run `--extra-args` or `--strict-mcp`).
