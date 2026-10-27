# Claude and Codex accounts

Crow drives any cdx account. The provider comes from `cdx status --cached --json` (`crow-accounts provider <name>`);
an account cdx does not describe counts as Claude. Policy (`use`, `ask-first`, `never`) in `local/accounts.yaml` is unchanged.

## What differs on Codex

| Area | Claude | Codex |
|---|---|---|
| Hub instructions | `CLAUDE.md` and linked memory, loaded by the client | `AGENTS.md` only; the prompt lists the files to read first (`CLAUDE.md`, `AGENTS.md`, `LOGICS.md`, owner prefs, context, memory index) |
| Hub skills (`brief`, `daily`, `weekly`, `retro`) | slash command | prompt names the skill file under `.claude/skills/`, `skills/` or `local/skills/`; a missing file stops the run |
| One-shot permission | `auto` | `default` (writes stay in the workspace, no network) |
| Background dispatch | `auto` | `default`, prompt starts with the read-first list |
| Retro model | `CROW_RETRO_MODEL` (default `sonnet`) | `CROW_RETRO_CODEX_MODEL` (default: the account's own model) |
| Model name | `claude-*`, opus, sonnet, haiku, fable | `gpt-*`, `o<N>`, `codex*`; a name from the other provider is refused |

## Guarded operations

The hub guard (`bin/crow-hub-guard`) is a Claude hook. Codex has no hook, so crow does not rely on it there.

| Operation | Claude | Codex |
|---|---|---|
| `crow` (interactive hub session) | hook asks before commit, push, rebase, MR | refused (`crow-accounts hub-session`) |
| `crow brief`, `daily`, `weekly`, retros | hook | sandboxed run (`default` permission) |
| `crow open/term/handoff <project>` | the repo's own rules | the repo's own rules (`AGENTS.md`) |
| `crow dispatch` | `auto` | `default` |

A brief on Codex has no network, so it cannot reach the forge. Run it on Claude for the full brief.

## Not covered yet

- Live activity of a Codex session (transcript reading): `crow watch` lists the process only.
- Session titles use the `claude` CLI directly (`CROW_TITLE_CMD`).
- Cross-provider relaunch, Codex usage and cost.

## Check

`crow-accounts provider <name>`, `crow-accounts permission <name>`, `crow-accounts prompt <name> brief`.
