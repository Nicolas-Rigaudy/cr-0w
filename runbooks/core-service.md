# Core service (state API and Nest over the tailnet)

`crow-core serve` serves the read API under `/api` and the built Nest at `/` on one origin.

## Install and check
1. `cd tools/nest && npm run build` (core serves `tools/nest/dist`).
2. `bin/crow core install`: writes `~/.config/systemd/user/crow-core.service` under the real home (cdx sandboxes HOME), enables and restarts it, runs `tailscale serve --bg --https=443 http://127.0.0.1:8787`.
3. Linger keeps the unit up after logout: `loginctl enable-linger <user>` (install prints the hint when it is off).
4. `bin/crow core status`: unit state, listen address, socket, Nest dir, `tailscale serve status`.
5. After rebuilding the binary (any `bin/crow core|mission|inbox` call rebuilds when Go sources are newer): `bin/crow core restart`.

## Access rules
- TCP `127.0.0.1:8787` needs the `Tailscale-User-Login` header equal to the owner login. tailscale serve sets it for tailnet requests. The login comes from `serve.owner_login` in the config, else from `tailscale status --json`.
- Unix socket `local/state/core.sock` (mode 0600) is trusted: no identity header.
- Non-GET requests need `X-Crow: 1` and a same-origin `Origin`. No CORS headers, ever.
- Artifacts: only regular files under the mission's `outbox/`, within the kind's per-file cap; anything else is 404.

## Probe
- `curl -s --unix-socket local/state/core.sock http://x/api/state/inbox`
- `curl -s https://<host>.<tailnet>.ts.net/api/state/inbox` from any tailnet node of the owner.
- `curl -s -o /dev/null -w '%{http_code}' http://127.0.0.1:8787/api/state` answers 403 (no identity).
- Repeat with `If-None-Match: <etag>`: 304 until the store version moves.

## Config (`serve:` in the crow config)
- `listen` (default `127.0.0.1:8787`), `socket`, `nest` (relative to the crow home), `owner_login`.

## Nest against the live core
- `cd tools/nest && VITE_NEST_SOURCE=api npm run dev`: the dev server proxies `/api` to the unix socket. Plain `npm run dev` and tests use fixtures; production builds always read the API.
- After `npm run build` the core serves the new `dist/` at once; no restart.
- Test mission for a live check: `bin/crow mission new --kind research --project cr-0w --goal "..." --end-state "probe only"` lands `queued` (nothing ticks on its own). Then `bin/crow mission cancel <id> --reason ...`, remove `local/state/missions/<id>/` and run `tools/crow-core/crow-core rebuild` to drop it from the store.

## Fixtures
- `go test ./internal/api -update` in `tools/crow-core` rewrites `logics/specs/core_contract/fixtures/api/*.json`; `tests/test_core_contract.py` checks them against the record schemas.
