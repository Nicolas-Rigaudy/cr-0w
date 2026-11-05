# Session retros

Sessions worth a look get a short retrospective, so lessons land in skills, runbooks and memories without the user repeating them.
One batch a day, over the sessions a no-model pre-filter keeps.

## Pipeline
1. **Timer**: `crow-retro.timer` (systemd user unit) runs `bin/crow-retro run` Mon–Fri at 19:00, `Persistent=true`. No weekend runs: Monday's window covers the weekend.
2. **Window**: from the end of the last successful batch (`local/state/retros/batch.json`, `last_end`; default 24 h ago) to 20 minutes ago, so live sessions settle first. A failed batch keeps the window for the next run.
3. **Pre-filter** (no model, deterministic): every Claude transcript in every cdx profile, modified in the window. Skipped with a reason:
   - headless runs (`entrypoint` `sdk-cli`: dispatch, retro, probes)
   - one-shot `/brief`, `/retro`, `/daily`, `/weekly`, `/eod`, `/wrap` sessions
   - under 3 typed prompts
   - no correction signal
   - already analysed (`done` ledger, 30 days)
   - over the cap of 8 kept sessions (lowest score first)

   Signals: correction words in a typed prompt after the first (the first is often a pasted hand-off), interrupts, rejected tool calls, the same tool error twice. Score = 2·corrections + 2·rejections + interrupts + repeated errors.
4. **Digest**: `bin/crow-digest <transcript.jsonl>` per kept session (~25 KB each).
5. **Analysis**: one `cdx run <account> --kind retro --failover --power low` (Sonnet) on `/retro batch=<manifest>`. The manifest is `local/state/retros/batches/YYYY-MM-DD.md`. Budget `min(1 + 0.75·n, 8)` USD, override with `CROW_RETRO_BUDGET`.
6. **Report**: `local/state/retros/YYYY-MM-DD-batch.md`, one section per session, then a "Skipped (pre-filter)" table appended by the script. The brief's 🔁 section lists it.

## Operate
- `bin/crow-retro scan [--since ISO]`: what the next run would keep and skip, no model, a few seconds.
- `bin/crow-retro run --dry-run`: scan, digests and manifest, no model call, no state change.
- `bin/crow-retro run`: refuses a second run the same day; `--force` overrides.
- `bin/crow-retro install-timer`: (re)writes and enables the units in the real home's `~/.config/systemd/user/`. Check with `systemctl --user list-timers crow-retro.timer`; logs with `journalctl --user -u crow-retro`.
- `bin/crow-retro status`: last batch, by-hand queue, worker log.
- `bin/crow-retro add <transcript.jsonl>`: retro one session by hand (queue + worker, one report `YYYY-MM-DD-<project>-<session8>.md`). Transcripts: `<profile>/.claude/projects/<cwd-slug>/<session>.jsonl`.
- By-hand failures move to `local/state/retros/failed/`; `bin/crow-retro retry` requeues usage-limit failures.
- Runs started by crow-retro carry `CROW_RETRO=1`.

## Known gaps
- A session resumed after its retro is skipped as already analysed. Use `add` for it.
- The correction regex (English and French) misses polite redirections; `scan` shows what it catches. Tune `CORRECTION` in `bin/crow-retro` and the fixtures in `tests/fixtures/retro/`.
- Sonnet tends to miss causes on the hub side (a vague handoff prompt). The report has a mandatory "Hub handoff" line per session; spot-check it.
- The old SessionEnd hook (one run per session end) was removed 2026-11-05. If cdx re-adds hooks to a profile's `.claude/settings.json`, check `jq .hooks.SessionEnd` there.
- Retrying `claude -p` unsandboxed to test auth is blocked by the auto-mode classifier. Test inside the sandbox or ask the user to run it with `!`.
