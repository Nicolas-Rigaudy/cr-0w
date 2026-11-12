// Headless hero layout shots: rest, hover (moons), card, at phone / inner / laptop sizes.
// Each rest and hover shot also saves zoom crops (z-*.png): one per planet at rest, and the
// hovered planet with its moon titles. Review the crops, not only the full shots: small
// spheres and moons are a few pixels wide there.
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
    // Zoom crop: a 2x capture of a small CSS box is enough to read single glyphs.
    const zoom = async (file, box) => {
      const b = await hero.boundingBox(), pad = 14
      const x = Math.max(b.x, box.x0 - pad), y = Math.max(b.y, box.y0 - pad)
      const clip = { x, y, width: Math.min(b.x + b.width, box.x1 + pad) - x, height: Math.min(b.y + b.height, box.y1 + pad) - y }
      if (clip.width > 4 && clip.height > 4) await page.screenshot({ path: `${out}/${file}`, clip })
    }
    const planets = () => page.evaluate(() => [...document.querySelectorAll('.hero-pl[data-c]')].filter((e) => e.offsetParent).map((e) => {
      const [x, y] = e.dataset.c.split(',').map(Number), b = e.closest('.hero-box').getBoundingClientRect(), l = e.getBoundingClientRect()
      const cx = b.left + x, cy = b.top + y, r = 40
      return { name: e.textContent.trim().split(/\s/)[0].toLowerCase(), x0: Math.min(cx - r, l.left), x1: Math.max(cx + r, l.right), y0: Math.min(cy - r, l.top), y1: Math.max(cy + r, l.bottom) }
    }))
    await hero.screenshot({ path: `${out}/${size}-1-rest.png` })
    for (const z of await planets()) await zoom(`${size}-z-rest-${z.name}.png`, z)
    let p = await pos('demo'); await page.mouse.move(p.x, p.y); await page.waitForTimeout(600)
    await hero.screenshot({ path: `${out}/${size}-2-hover.png` })
    const mz = await page.evaluate(({ x, y }) => {
      const els = [...document.querySelectorAll('.hero-ml')].map((e) => e.getBoundingClientRect())
      const xs = [x - 60, x + 60, ...els.flatMap((r) => [r.left, r.right])], ys = [y - 60, y + 60, ...els.flatMap((r) => [r.top, r.bottom])]
      return { x0: Math.min(...xs), x1: Math.max(...xs), y0: Math.min(...ys), y1: Math.max(...ys), shown: els.filter((r) => r.width > 0).length, all: els.length }
    }, p)
    await zoom(`${size}-z-hover-moons.png`, mz)
    console.log(size, 'moon titles', [...await page.$$eval('.hero-ml', (es) => es.map((e) => getComputedStyle(e).visibility !== 'hidden' && getComputedStyle(e).opacity !== '0'))].filter(Boolean).length, '/', mz.all)
    await page.mouse.click(p.x, p.y); await page.waitForTimeout(1200); await page.mouse.move(2, 2); await page.waitForTimeout(500)
    await hero.screenshot({ path: `${out}/${size}-3-card.png` })
    await ctx.close()
  }
} finally { await browser.close(); server.kill() }
