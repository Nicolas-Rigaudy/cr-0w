# Headless WebGL render and frame-time bench

Render a WebGL/canvas page to PNG and measure its frame time without opening a window on the desktop.
The Playwright MCP is unreliable here; call the Playwright Chromium binary directly.

## Binary
`$CHROME` = the Playwright cache Chromium (`<cache>/playwright/chromium-*/chrome-linux64/chrome`).
Base flags: `--headless=new --no-sandbox --hide-scrollbars --window-size=W,H --force-device-scale-factor=2`.

## Render a frame
- Software GL (always works): add `--use-angle=swiftshader --enable-unsafe-swiftshader`.
- `--virtual-time-budget=4000 --screenshot=/abs/out.png "file:///abs/page.html?..."`.
- Make the page deterministic from URL params (`t=` fixed time, mood/state, `glitch=1`), and wait for `document.fonts.load` before building any glyph atlas.

## Bench
- Page side: render N frames, call `gl.readPixels(0,0,1,1)` after each to force the GPU to finish, then write `document.body.dataset.ms` plus the unmasked renderer (`WEBGL_debug_renderer_info`).
- Shell side: `--dump-dom ... | grep -o '<body[^>]*>'` reads the numbers.
- Real GPU in headless: `--use-angle=gl --ignore-gpu-blocklist` (or `--use-angle=vulkan`). Check `data-gpu` names the real GPU, not SwiftShader.
- SwiftShader timing is the software worst case, not the device number. Report both.

## Notes
- Software renderers (SwiftShader, llvmpipe) run a GPU ASCII shader 100x slower: detect them via the unmasked renderer string and fall back to static.
- Crop reference sheets with `magick in.png -crop WxH+X+Y +repage out.png` (device pixels).
- Show results with `bin/crow-compare` and open with `bin/crow-open`.
