// Live GL head probe (item_156/196/198): renders the head in key states on the GPU, times a frame,
// and writes a contact sheet. Usage: node scripts/head-probe.mjs <out.png> [cell css px, default 2.75] [states.json]
// states.json: [[name, frame overrides, optional js run with H (the head module) before the render]].
import { spawn } from 'node:child_process'
import { readFileSync } from 'node:fs'
import { chromium } from 'playwright-core'
const out = process.argv[2] ?? 'head-probe.png', cw = Number(process.argv[3] ?? 2.75)
const spec = process.argv[4] ? JSON.parse(readFileSync(process.argv[4], 'utf8')) : null
const PORT = 4393
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
try {
  const ctx = await browser.newContext({ viewport: { width: 1600, height: 900 }, deviceScaleFactor: 2 })
  const page = await ctx.newPage()
  page.on('console', (m) => { if (m.type() !== 'debug') console.log('[page]', m.text()) })
  await page.goto(base + '/', { waitUntil: 'networkidle' })
  const r = await page.evaluate(async ([cw, spec]) => {
    await document.fonts.load("500 12px 'Fira Code'")
    const H = await import('/src/hero/gl/head.ts')
    document.body.innerHTML = ''; document.body.style.cssText = 'margin:0;background:#05060d;display:grid;grid-template-columns:repeat(4,auto);gap:4px'
    const g0 = document.createElement('canvas').getContext('webgl2'), dbg = g0.getExtension('WEBGL_debug_renderer_info')
    const renderer = dbg ? g0.getParameter(dbg.UNMASKED_RENDERER_WEBGL) : '?'
    const base = { yaw: 0, pitch: 0, tilt: 0, blink: 0, gaze: [0, 0], face: 'idle', pose: { from: 'rest', to: 'rest', k: 0 }, tip: [46, 242, 255], hold: [46, 242, 255], glitch: 0, flicker: 1, shimmer: 0 }
    const states = spec ?? [
      ['front', {}], ['yaw -20', { yaw: -20, gaze: [-3, 0] }], ['yaw +20', { yaw: 20, gaze: [3, 0] }], ['blink', { blink: 0.9 }],
      ['run', { face: 'run', pose: { from: 'rest', to: 'grab', k: 1 }, eye: [46, 242, 255] }], ['fail', { face: 'fail', pose: { from: 'rest', to: 'brace', k: 1 }, eye: [255, 46, 99] }],
      ['present', { pose: { from: 'rest', to: 'present', k: 1 }, tilt: 7, yaw: 12 }], ['scan+shimmer', { scan: 0.45, shimmer: 0.02, glitch: 0 }],
    ]
    const times = []
    const keep = { pose: { ...H.POSE }, face: JSON.parse(JSON.stringify(H.FACE_STYLE)), st: { ...H.ST } }
    for (const [name, o, js] of states) {
      Object.assign(H.POSE, keep.pose); Object.assign(H.ST, keep.st); Object.assign(H.FACE_STYLE, JSON.parse(JSON.stringify(keep.face)))
      if (js) new Function('H', js)(H)
      const cv = document.createElement('canvas'), wrap = document.createElement('div')
      wrap.style.cssText = 'position:relative;color:#8af;font:11px monospace'; wrap.append(name, document.createElement('br'), cv); document.body.append(wrap)
      const h = H.createHead(cv); if (!h) return { err: 'no head' }
      const l = h.layout({ cw, ch: cw * 7 / 4, font: cw * 6.65 / 4, dpr: 2 }); cv.style.width = l.w + 'px'; cv.style.height = l.h + 'px'
      const f = { ...base, ...o }
      h.render(f, 1)
      const gl = cv.getContext('webgl2'); const px = new Uint8Array(4)
      for (let i = 0; i < 20; i++) { const a = performance.now(); h.render(f, 1 + i / 30); gl.readPixels(0, 0, 1, 1, gl.RGBA, gl.UNSIGNED_BYTE, px); times.push(performance.now() - a) }
    }
    times.sort((a, b) => a - b)
    return { renderer, med: times[times.length >> 1].toFixed(2), p95: times[Math.floor(times.length * 0.95)].toFixed(2) }
  }, [cw, spec])
  console.log(JSON.stringify(r))
  await page.waitForTimeout(200)
  await page.screenshot({ path: out, fullPage: true })
} finally { await browser.close(); server.kill() }
