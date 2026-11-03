# Personal service tokens (PATs)

API tokens for scripts and sessions (SonarCloud, Jira, Graph, any vendor PAT) live in the desktop keyring. One tool reads and writes them: `crow secret`. CI variables are a different job: see [ci-secrets](ci-secrets.md).

## Why the keyring
- The keyring is reached over D-Bus. Every cdx sandboxed HOME sees the same items; a file under `~` exists once per HOME and drifts.
- Never write a token to a file under `~`, the repo, `local/` or a prompt. Check what a CLI does: gh and glab keep their tokens in the keyring when `use_keyring` is on; a glab config written in a sandbox HOME holds the token in plaintext.
- Items use the `secret-tool` generic schema with two attributes: `service=<name>` and `user=<account>`.

## Registry
`local/tokens.yaml` lists every token: service, user, scope, created, expiry, where to rotate, who uses it. Format in `docs/tokens.example.yaml`. No values in it. `crow secret` takes the default user from it.

## Before creating a token
- Least scope: read-only, one project or org when the service allows it. When the service has no scopes (SonarCloud user tokens carry the account's rights), say so in the registry.
- Shortest expiry the job allows (30 or 90 days). Write `created` and `expiry` in the registry.

## Commands
| Step | Command | Who |
|---|---|---|
| Store or replace | `crow secret set <service> [--user U] [--label L]` | owner, real terminal |
| Read | `crow secret get <service> [--user U]` | agents, scripts |
| List | `crow secret list` | anyone: labels, users, dates, expiry, registry entries not stored; never values |
| Remove | `crow secret rm <service> [--user U]` | owner or agent on request |

- `set` prompts on the tty with echo off. It refuses argv values, pipes and the `!` shell (no tty there; an older secret-tool run stored an empty value silently). Paste the token in a real terminal.
- Rotate: create the new token on the service, `crow secret set <service>`, check a call works, revoke the old token on the service page from the registry, update `created` and `expiry`.
- Revoke: revoke on the service page first, then `crow secret rm <service>` and drop or mark the registry entry.

## Reading a token without printing it
Use it inside a command substitution, never bare (a bare `get` prints the value into the transcript):

```bash
curl -su "$(crow secret get sonarcloud):" "https://sonarcloud.io/api/issues/search?organization=<org>&componentKeys=<project>"
curl -sH "Authorization: Bearer $(crow secret get <service>)" https://api.example.com/...
t=$(crow secret get <service>); echo "${#t} chars"   # check it is stored
```

`get` exits 2 when nothing is stored. Python scripts import `secret_lookup` and `secret_store` from `bin/crow_lib.py` (`bin/crow-today` does).
