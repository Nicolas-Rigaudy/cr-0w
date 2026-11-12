// Headless drag check: table paints/s, GL head frames/s and frame gaps (p50/p95/max ms) while the
// pointer drags the table for 2 s, then during the 1.5 s release coast. Smooth means about 60 fps.
// Usage: node scripts/hero-drag.mjs [width height]
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const [W = '1440', H = '900'] = process.argv.slice(2)
const PORT = 4393
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
try {
  const ctx = await browser.newContext({ viewport: { width: Number(W), height: Number(H) }, deviceScaleFactor: 2 })
  await ctx.addInitScript(() => {
    window.__t = []; window.__g = []
    const f = CanvasRenderingContext2D.prototype.drawImage
    // One table paint ends with the sharp layer: an unfiltered draw onto the visible canvas.
    CanvasRenderingContext2D.prototype.drawImage = function (...a) { if (this.canvas.isConnected && this.filter === 'none') window.__t.push(performance.now()); return f.apply(this, a) }
    const d = WebGL2RenderingContext.prototype.drawArrays
    WebGL2RenderingContext.prototype.drawArrays = function (...a) { if (this.getParameter(this.FRAMEBUFFER_BINDING) === null) window.__g.push(performance.now()); return d.apply(this, a) }
  })
  const page = await ctx.newPage()
  await page.goto(base + '/', { waitUntil: 'networkidle' })
  await page.waitForTimeout(2500)
  const box = await page.locator('.hero-box').boundingBox()
  const y = box.y + box.height * 0.85, x0 = box.x + box.width * 0.2
  const stats = (ts, a, b) => {
    const v = ts.filter((t) => t >= a && t <= b), g = v.slice(1).map((t, i) => t - v[i]).sort((p, q) => p - q)
    const q = (k) => (g.length ? g[Math.min(g.length - 1, Math.floor(k * g.length))].toFixed(1) : '-')
    return `${(v.length / ((b - a) / 1000)).toFixed(0)} fps, gap p50 ${q(0.5)} p95 ${q(0.95)} max ${g.length ? g.at(-1).toFixed(0) : '-'}`
  }
  await page.mouse.move(x0, y); await page.mouse.down()
  const a = await page.evaluate(() => performance.now())
  for (let i = 1; i <= 120; i++) { await page.mouse.move(x0 + (box.width * 0.6 * i) / 120, y); await page.waitForTimeout(16) }
  const b = await page.evaluate(() => performance.now())
  await page.mouse.up(); await page.waitForTimeout(1500)
  const c = await page.evaluate(() => performance.now())
  const [t, g] = await page.evaluate(() => [window.__t, window.__g])
  console.log(`${W}x${H} drag   table ${stats(t, a, b)} | head ${stats(g, a, b)}`)
  console.log(`${W}x${H} coast  table ${stats(t, b, c)} | head ${stats(g, b, c)}`)
} finally { await browser.close(); server.kill() }
