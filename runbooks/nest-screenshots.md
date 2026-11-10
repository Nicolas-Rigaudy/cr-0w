# Nest screenshots (headless)

Capture the Nest PWA screens for design review without opening a window.

1. `cd tools/nest && npm ci && npm run build`
2. Install the browser once, outside the sandboxed HOME so it survives sessions:
   `PLAYWRIGHT_BROWSERS_PATH=<real home>/.cache/ms-playwright npx playwright-core install chromium-headless-shell`
3. `PLAYWRIGHT_BROWSERS_PATH=<same path> npm run shots -- <out_dir>`
   writes `<screen>-<cover|inner|laptop>-<modern|cyber>.png` (24 files, 2x scale, Europe/Paris time).
4. Put review images in the vault crow folder (`Research/nest-design/<item>/`), never in `local/state/`.
   Compare with the boards using `bin/crow-compare`.

Gotchas:
- The script runs `vite preview` on port 4317; it needs `dist/`.
- Add a screen by editing `SCREENS` in `scripts/shots.mjs`.
- An empty `package.json` at the repo root breaks Vite resolution above `tools/nest`; the package reads shared fixtures through a virtual module for that reason.
