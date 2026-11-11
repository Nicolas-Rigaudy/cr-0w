# Nest screenshots (headless)

Capture the Nest PWA screens for design review without opening a window.

1. `cd tools/nest && npm ci && npm run build`
2. Install the browser once, outside the sandboxed HOME so it survives sessions:
   `PLAYWRIGHT_BROWSERS_PATH=<real home>/.cache/ms-playwright npx playwright-core install chromium-headless-shell`
3. `PLAYWRIGHT_BROWSERS_PATH=<same path> npm run shots -- <out_dir>`
   writes `<screen>-<cover|inner|laptop>-<modern|cyber>.png` (24 files, 2x scale, Europe/Paris time).
4. Put review images in the vault crow folder (`Research/nest-design/<item>/`), never in `local/state/`.
   Compare with the boards using `bin/crow-compare`.

## Board vs build review (one screen at a time)

1. Render the approved boards the item names (never the running app or older PNGs):
   `cd tools/nest && PLAYWRIGHT_BROWSERS_PATH=<same path> node scripts/boards.mjs <boards_dir> <out_dir> fold-inbox-inner fold-inbox-inner-cyber ...`
   Each board has a `-cyber` variant; render both themes.
2. Shoot the matching app screen with `npm run shots` (or a one-off Playwright script on the fixture build at cover 368x812 and inner 707x691, 2x).
3. Copy both sets into `Research/nest-design/<item>/<screen>-compare/` with a `spec.json`:
   one decision, `rows` = sizes/themes, two `options` (Board, Build) whose `images` line up with the rows.
   `context` lists every deliberate gap (data that does not exist yet). Run `bin/crow-compare <spec.json>`.
4. Open the page in the real browser (`HOME=<real home> xdg-open .../spec.html`), then ask for the OK.
5. When no board draws an element, ask where it goes. Never draw it from the running app.

Gotchas:
- The script runs `vite preview` on port 4317; it needs `dist/`.
- Add a screen by editing `SCREENS` in `scripts/shots.mjs`.
- An empty `package.json` at the repo root breaks Vite resolution above `tools/nest`; the package reads shared fixtures through a virtual module for that reason.
