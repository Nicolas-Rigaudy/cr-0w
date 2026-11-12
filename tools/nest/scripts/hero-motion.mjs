// Headless motion mockup (item_198): records the hero through one scripted beat sequence per variant and
// writes a cropped animated WebP each (for crow-compare). Variants set the localStorage mockup switches.
// Usage: node scripts/hero-motion.mjs <out dir> <variants.json> [width height]
// variants.json: [{ "name": "calm", "ls": { "nest.char": "calm", "nest.fx": "{...}" } }, ...]
// Beats: idle, a fast pointer sweep, hover a planet, select it, idle again (about 12 s).
import { spawn, execFileSync } from 'node:child_process'
import { mkdirSync, readFileSync, readdirSync, renameSync, rmSync } from 'node:fs'
import { join } from 'node:path'
import { chromium } from 'playwright-core'
const [out, spec, W = '1440', H = '900'] = process.argv.slice(2)
const variants = JSON.parse(readFileSync(spec, 'utf8'))
mkdirSync(out, { recursive: true })
const PORT = 4396
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
const vw = Number(W), vh = Number(H)
try {
  // Warm vite once so the first variant does not record a dependency reload.
  { const p = await browser.newPage(); await p.goto(base + '/', { waitUntil: 'networkidle' }); await p.waitForTimeout(1500); await p.close() }
  for (const v of variants) {
    const tmp = join(out, `.rec-${v.name}`); rmSync(tmp, { recursive: true, force: true })
    const ctx = await browser.newContext({ viewport: { width: vw, height: vh }, recordVideo: { dir: tmp, size: { width: vw, height: vh } } })
    await ctx.addInitScript((ls) => { for (const [k, x] of Object.entries(ls)) localStorage.setItem(k, x) }, v.ls ?? {})
    const page = await ctx.newPage()
    const t0 = Date.now()
    await page.goto(base + '/', { waitUntil: 'networkidle' })
    await page.waitForTimeout(1200)
    const box = await page.locator('.hero-box').boundingBox()
    const cy = box.y + box.height * 0.45
    const lead = (Date.now() - t0) / 1000
    await page.mouse.move(box.x + 20, cy); await page.waitForTimeout(2500)
    // Fast sweep across the hero in about 150 ms.
    await page.mouse.move(box.x + box.width - 20, cy, { steps: 10 }); await page.waitForTimeout(1800)
    // The planet disc centre (data-c, hero px); the label may overlap the head, which opens Desk.
    const [px, py] = (await page.locator('.hero-pl').first().getAttribute('data-c')).split(',').map(Number)
    await page.mouse.move(box.x + px, box.y + py, { steps: 12 }); await page.waitForTimeout(1500)
    await page.mouse.click(box.x + px, box.y + py); await page.waitForTimeout(2200)
    if (!page.url().endsWith('/')) throw new Error(`${v.name}: the select click navigated to ${page.url()}`)
    await page.mouse.move(box.x + box.width / 2, box.y + 4, { steps: 8 })
    await page.mouse.move(2, 2, { steps: 4 }); await page.waitForTimeout(3000)
    await ctx.close()
    const webm = readdirSync(tmp).find((f) => f.endsWith('.webm'))
    const crop = `crop=${Math.round(box.width)}:${Math.round(box.height)}:${Math.round(box.x)}:${Math.round(box.y)}`
    const dst = join(out, `${v.name}.webp`)
    execFileSync('ffmpeg', ['-y', '-loglevel', 'error', '-ss', String(lead), '-i', join(tmp, webm), '-vf', `${crop},fps=15`, '-loop', '0', '-quality', '80', '-c:v', 'libwebp_anim', dst])
    renameSync(join(tmp, webm), join(out, `${v.name}.webm`)); rmSync(tmp, { recursive: true, force: true })
    console.log('wrote', dst)
  }
} finally { await browser.close(); server.kill() }
