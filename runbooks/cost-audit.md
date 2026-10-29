# Cost audit

Measure where hub tokens and money go, and re-measure after a change with the same method.

## Sources
- `crow usage --json --days 14`: per-session cost, kind, project, title.
- `cdx stats --by project`, `cdx runs --since 30d`: headless runs and their cost.
- Claude transcripts in every cdx profile: `<profiles>/*/claude-home/.claude/projects/*/*.jsonl`.

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
