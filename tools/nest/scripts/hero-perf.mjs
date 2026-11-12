// Headless hero perf: median drawScene time per frame on the laptop size at dpr 2.
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const PORT = 4391
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true })
try {
  for (const [w, h] of [[1440, 900], [707, 691], [368, 812]]) {
    const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2 })
    const page = await ctx.newPage()
    await page.goto(base + '/', { waitUntil: 'networkidle' })
    await page.waitForTimeout(800)
    const r = await page.evaluate(async () => {
      const sc = await import('/src/hero/scene.ts'); const md = await import('/src/hero/model.ts'); const fx = await import('/src/data/fixtures.ts')
      const heads = (await import('/src/hero/heads.json')).default
      const src = fx.fixtureSource(); const m = md.heroModel(await src.missions(), await src.inbox(), src.now())
      const box = document.querySelector('.hero-box'); const W = box.clientWidth, H = box.clientHeight
      const cv = document.createElement('canvas'); const pal = { ring: [79,195,255], dim: [95,120,168], fg: [216,230,255], accent: [192,92,255], state: { needs: [255,184,0], working: [46,242,255], failed: [255,46,99], sleeping: [111,134,184], queued: [138,160,200], done: [95,120,168] } }
      const t = []
      for (let i = 0; i < 12; i++) { const a = performance.now(); sc.drawScene(cv, W, H, 2, m, heads, pal, undefined, { now: src.now(), agenda: await src.agenda(), spin: i * 5 }); t.push(performance.now() - a) }
      t.sort((a, b) => a - b); return { W, H, med: t[6].toFixed(1), max: t[11].toFixed(1) }
    })
    console.log(w, h, JSON.stringify(r))
  }
} finally { await browser.close(); server.kill() }
