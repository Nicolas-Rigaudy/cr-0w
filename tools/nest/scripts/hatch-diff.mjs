// Ops live hatch (item_249): the transform hatch must look like the old background-position hatch.
// Shoots /ops with every animation frozen at the same times, once as built and once with the old
// rule injected, and counts the differing pixels. Usage: node scripts/hatch-diff.mjs <out dir>
import { spawn } from 'node:child_process'
import { writeFileSync, mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'
const OUT = process.argv[2]
const PORT = Number(process.env.PERF_PORT) || 4394
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const OLD = `.live::before { display: none !important; }
.live { background-image: repeating-linear-gradient(135deg, rgba(255, 255, 255, 0.28) 0 3px, transparent 3px 7px); animation: ops-hatch-old 1.2s linear infinite; background-size: 20px 20px; }
@keyframes ops-hatch-old { to { background-position: 20px 0; } }`
const browser = await chromium.launch({ headless: true })
let bad = 0
try {
  if (OUT) mkdirSync(OUT, { recursive: true })
  for (const [w, h] of [[1440, 900], [368, 812]]) {
    const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2, timezoneId: 'Europe/Paris' })
    const page = await ctx.newPage()
    await page.goto(base + '/ops', { waitUntil: 'networkidle' })
    await page.waitForTimeout(800)
    // The rail head and hero keep their own clocks: drop their canvases from both shots.
    await page.addStyleTag({ content: 'canvas { visibility: hidden !important; }' })
    const n = await page.locator('.live').count()
    const shoot = async (ms) => {
      await page.evaluate((ms) => { for (const a of document.getAnimations()) { a.pause(); a.currentTime = ms } }, ms)
      return page.screenshot()
    }
    for (const ms of [0, 300, 700]) {
      const now = await shoot(ms)
      const tag = await page.addStyleTag({ content: OLD })
      const old = await shoot(ms)
      await tag.evaluate((el) => el.remove())
      const diff = await page.evaluate(async ([a, b]) => {
        const img = async (s) => { const i = new Image(); i.src = 'data:image/png;base64,' + s; await i.decode(); return i }
        const [ia, ib] = [await img(a), await img(b)]
        const px = (i) => { const c = new OffscreenCanvas(i.width, i.height), x = c.getContext('2d'); x.drawImage(i, 0, 0); return x.getImageData(0, 0, i.width, i.height).data }
        const [da, db] = [px(ia), px(ib)]
        // Edge anti-aliasing of a composited layer moves a channel by a few levels; a stripe out of
        // place moves it by far more. Only the second counts as a visible change.
        let aa = 0, vis = 0, worst = 0
        for (let k = 0; k < da.length; k += 4) {
          const m = Math.max(Math.abs(da[k] - db[k]), Math.abs(da[k + 1] - db[k + 1]), Math.abs(da[k + 2] - db[k + 2]))
          worst = Math.max(worst, m)
          if (m > 24) vis++; else if (m > 2) aa++
        }
        return { aa, vis, worst }
      }, [now.toString('base64'), old.toString('base64')])
      bad += diff.vis
      console.log(`${w}x${h} t=${ms}ms live=${n} visible diff px: ${diff.vis}, anti-aliasing px: ${diff.aa}, worst channel delta: ${diff.worst}`)
      if (OUT && ms === 300) { writeFileSync(`${OUT}/ops-${w}-hatch-after.png`, now); writeFileSync(`${OUT}/ops-${w}-hatch-before.png`, old) }
    }
    await ctx.close()
  }
} finally { await browser.close(); server.kill() }
process.exitCode = bad ? 1 : 0
