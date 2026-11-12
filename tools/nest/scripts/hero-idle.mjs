// Headless idle check: counts hero canvas paints per second at rest (0 expected with reduced motion).
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const PORT = 4392
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true })
try {
  for (const rm of ['no-preference', 'reduce']) {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 }, reducedMotion: rm })
    await ctx.addInitScript(() => {
      window.__draws = 0
      const f = CanvasRenderingContext2D.prototype.drawImage
      CanvasRenderingContext2D.prototype.drawImage = function (...a) { window.__draws++; return f.apply(this, a) }
    })
    const page = await ctx.newPage()
    await page.goto(base + '/', { waitUntil: 'networkidle' })
    await page.waitForTimeout(3000)
    const a = await page.evaluate(() => window.__draws)
    await page.waitForTimeout(3000)
    const b = await page.evaluate(() => window.__draws)
    console.log(rm, 'paints in 3s idle:', (b - a) / 3)
  }
} finally { await browser.close(); server.kill() }
