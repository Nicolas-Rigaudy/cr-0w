// Headless hero layout shots: rest, hover (moons), card, at phone / inner / laptop sizes.
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const out = process.argv[2]
const PORT = 4320, base = 'http://localhost:' + PORT
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
for (let i = 0; i < 50; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true })
const SIZES = { cover: [368, 812, 'modern'], 'cover-cyber': [368, 812, 'cyber'], inner: [707, 691, 'cyber'], laptop: [1440, 900, 'cyber'] }
const only = process.argv[3]
try {
  for (const [size, [w, h, theme]] of Object.entries(SIZES)) {
    if (only && !size.startsWith(only)) continue
    const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2, timezoneId: 'Europe/Paris' })
    await ctx.addInitScript((t) => localStorage.setItem('nest.theme', t), theme)
    const page = await ctx.newPage()
    await page.goto(base + '/', { waitUntil: 'networkidle' }); await page.waitForSelector('.hero-pl[data-c]', { timeout: 15000 }); await page.waitForTimeout(800)
    const hero = page.locator('.hero-box'); await hero.scrollIntoViewIfNeeded()
    const pos = async (name) => page.evaluate((n) => { const o = [...document.querySelectorAll('.hero-pl')].find((e) => e.textContent.toLowerCase().startsWith(n)) ?? document.querySelector('.hero-pl'); const [x, y] = o.dataset.c.split(',').map(Number); const b = o.closest('.hero-box').getBoundingClientRect(); return { x: b.left + x, y: b.top + y } }, name)
    await hero.screenshot({ path: `${out}/${size}-1-rest.png` })
    let p = await pos('demo'); await page.mouse.move(p.x, p.y); await page.waitForTimeout(600)
    await hero.screenshot({ path: `${out}/${size}-2-hover.png` })
    await page.mouse.click(p.x, p.y); await page.waitForTimeout(1200); await page.mouse.move(2, 2); await page.waitForTimeout(500)
    await hero.screenshot({ path: `${out}/${size}-3-card.png` })
    await ctx.close()
  }
} finally { await browser.close(); server.kill() }
