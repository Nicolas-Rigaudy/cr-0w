# Nest

Mobile-first PWA for CR-0W: Home, Missions, Mission detail and Inbox. This slice reads core contract fixtures only. It is read only: no replies, no state changes.

## Run

```bash
npm ci
npm run dev          # vite on all interfaces, port 5173
npm run typecheck
npm test             # vitest: data layer and state mapping
npm run build        # tsc + vite build into dist/
npm run preview      # serve dist/ (service worker active)
```

Fixture validity against the core schemas is checked from the repo root:
`python3 -m pytest -q tests/test_nest_fixtures.py`.

## Install on a phone

A PWA needs HTTPS to install. On the tailnet:

```bash
npm run build && npm run preview -- --port 4173
tailscale serve --bg 4173
```

Open the `https://<host>.<tailnet>.ts.net` URL on the phone, then "Add to Home screen".

## Data

- `src/data/source.ts`: the `DataSource` interface. Screens only see this and `useSnapshot()`.
- `src/data/fixtures.ts`: fixture adapter. It merges the shared records in
  `logics/specs/core_contract/fixtures/records/` with `fixtures/states/`, which adds one
  mission per UI state.
- `src/data/api.ts`: live adapter over the core read API (`/api/state/*`, same origin).
  `src/data/feed.ts` polls it every 5 s while the page is visible, with ETags, and keeps
  the last data while the core is down (header shows `core down · <age>`).
- `src/main.tsx` picks the adapter: production builds read the API; `npm run dev` and
  tests use fixtures unless `VITE_NEST_SOURCE=api`, which proxies `/api` to the core socket.
- The shared records and API response fixtures come through the `virtual:core-records`
  and `virtual:core-api` Vite modules (read with fs), so Vite never resolves files outside
  this package.

UI states map from contract states:

| UI | contract |
|---|---|
| working | running, preflight, publishing |
| needs you | waiting, review, or an open non-fyi inbox row |
| failed | failed |
| sleeping | queued, last run ended on usage_limit |
| queued | queued |
| done | done, cancelled |

## Screens and layout

- Cover (< 600px): bottom tab bar, state chips in the header.
- Inner (>= 600px): header pills, two columns. Laptop (>= 1200px): Home in three columns.
- Fold spanning (`horizontal-viewport-segments: 2`) splits on the hinge.
- Themes: `modern` and `cyber` on `<html data-theme>`, toggled in the header, kept in localStorage.
- Desk, MRs, Agenda and Reports are placeholders until the live API.

## Budget

JS stays under 90 kB gzip (React is about 60 kB of it). Last build (2026-11-11): JS 96.9 kB gzip, over budget; item_185 brings it back.

## Screenshots

```bash
npm run build
PLAYWRIGHT_BROWSERS_PATH=<cache> npx playwright-core install chromium-headless-shell
PLAYWRIGHT_BROWSERS_PATH=<cache> npm run shots -- <out_dir>
```

Headless only. Captures Home, Missions, Mission detail and Inbox at 368x812, 707x691 and
1440x900 in both themes.

`npm run build` gives the live-API build, so the shots show whatever core returns. For fixture
data, build with `NODE_ENV=development npx vite build --mode development`, take the shots, then
run `npm run build` again so `dist` serves the phone the live build.
