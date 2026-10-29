# Long sessions

Cache reads on a long context are most of the spend: Opus turns above 150k tokens of context cost about half of all spend in the 14-day audit (see `cost-audit.md`).

## Limits
- Cap: auto-compact window at 150k tokens in every allowed cdx profile (`env.CLAUDE_CODE_AUTO_COMPACT_WINDOW` in the profile's `settings.json`).
- Warning: 120k tokens. Override with `CROW_CTX_WARN` and `CROW_CTX_CAP` (read by `crow-usage` and `crow-ctx`).
- Context size is the last turn's input + cache read + cache write tokens of the session itself (subagents excluded).

## Set and check the cap
- `bin/crow-ctx cap check`: one line per allowed profile, exit 1 when one lacks the setting.
- `bin/crow-ctx cap apply`: writes the setting, keeps every other key, skips accounts with policy `never`. Idempotent. New sessions pick it up; running ones keep their window.
- Run `cap apply` again after adding an account.

## Where the warning shows
- `bin/crow-ctx` (text), `--json`, `--brief` (one line per live session, empty when none).
- `crow-usage --json`: `context_tokens` and `context_level` (`ok` | `warn` | `over`) per session.
- `crow-watch` Usage view: a live session today past the warning gets a yellow `▲` (red past the cap) and `NNNk ctx` before its title.
- The brief lists them under "Needs you".

## When to wrap
- A session is listed, or a topic just closed.
- Run `/wrap`, then `/clear` or start fresh. Do not push on to a second topic in the same context.

## Scope the rest
- One session, one end state: a result you can name in one line and check (an MR opened, a test green, a doc written).
- List what is left, cut it into pieces that each fit that rule, and start each with `crow handoff <project> --task "<one piece>" --end-state "<its result>"`.
- Pieces that touch different repos or different files go in separate handoffs.
