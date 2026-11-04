# Cost audit

Measure where hub tokens and money go, and re-measure after a change with the same method.

## Sources
- `crow usage --json --days 14`: per-session cost, kind, project, title.
- `cdx stats --by project`, `cdx runs --since 30d`: headless runs and their cost.
- Claude transcripts in every cdx profile: `<profiles>/*/claude-home/.claude/projects/*/*.jsonl`.

## Cost split (quick)
- `bin/crow-usage --split`: one line, last 14 days, cost by part: reload (cache writes), reread (cache reads), output, input, plus the runs' reload share.
- `bin/crow-usage --json`: `cost_split.by_kind` (all, interactive, run) for the same window; `sessions[].cost_split` and each `by_day` entry carry the parts too.
- The parts sum to the session cost. They use the same price table as the totals.

## Script
- `bin/crow-cost-audit --days N [--json] [--hub PATH]` runs the method below: dedup, context buckets by model, per project, per skill class, and the first-turn context of interactive hub sessions (`start_context`). Use it for the baseline and the re-measure (same `--days`).
- `bin/crow-usage` `by_model` gives each model its own turns (a multi-model session no longer lands on its first model); `tests/test_crow_usage_by_model.py` recounts a transcript.

## Method
1. Read assistant events from the transcripts of the window. Deduplicate on `(message.id, requestId)`: one API call is written as several lines and counting every line inflates cost about 2x.
2. Context size of a turn = `input_tokens + cache_read_input_tokens + cache_creation_input_tokens`.
3. Price each turn with list prices per model (input, output, cache read at 0.1x input, 1h cache write at 2x input, 5m write at 1.25x). Take prices from cdx's table (`cdx_manager.run_usage.token_prices`), as crow-usage does.
4. Bucket by model and context size: `<50k`, `50-100k`, `100-150k`, `>150k`. Report cost, share, turn count, median and p90 context for the top bucket, and the top bucket by project.
5. Skills and hooks: classify sessions from crow usage by first prompt (`/brief`, `/daily`, `/weekly`, retro digest, probe, other headless, interactive) and sum cost per class.
6. Report only computed numbers. Note the window, the price table and the dedup key in the report.

## Gotchas
- crow-usage's per-model split of multi-model sessions has been wrong before (Haiku inflated). Cross-check `by_model` against a transcript recount.
- Subagent transcripts live under `<session>/subagents/`; include them, they share the parent's project.
- A re-measure must use the same window length and the same price table as the baseline.

## Hub start context (item_029)
Baseline 2026-11-04, interactive sessions started in the hub, first-turn context (input + cache reads + cache writes), last 14 days, `bin/crow-cost-audit`-style recount:
- account with 2 plugins (ui-ux-pro-max, claude-hud) and the Atlassian MCP server: 42 sessions, median 48.8k, max 50.4k.
- account with 1 plugin (humanizer) and no MCP server: 11 sessions, median 42.6k, min 40.6k, max 43.1k.
- Difference about 6k tokens per session start: the plugin skill listings and the MCP tool names. The rest (about 42k) is Claude's own system prompt and tools, the hub CLAUDE.md (7 KB), the memory index (5 KB), the repo skills and hooks output.
- Last 7 days overall: 5533 turns, $277.54; 1080 turns past 150k context.
After the owner trims the profile (plugins or MCP are profile settings outside the repo), re-run the same measure with the same window and record the "after" line here.
