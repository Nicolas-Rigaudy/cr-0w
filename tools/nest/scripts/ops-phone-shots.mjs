// Ops phone checks (item_274) on fixture data: at 368x812 no card or the page overflows
// horizontally, day and week, both themes; the Now card is first with a working OPEN and a free slot row.
// Also shoots the laptop. Usage: node scripts/ops-phone-shots.mjs <out_dir>   (exits 1 on a failed check)
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) { console.error('usage: ops-phone-shots.mjs <out_dir>'); process.exit(2) }
mkdirSync(out, { recursive: true })
const PORT = Number(process.env.SHOTS_PORT) || 4319
const base = `http://localhost:${PORT}`
if (await fetch(base).then(() => true, () => false)) { console.error(`port ${PORT} is busy`); process.exit(2) }
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
for (let i = 0; i < 100; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }

const failed = []
const browser = await chromium.launch()
try {
  for (const theme of ['modern', 'cyber']) {
    for (const [size, [w, h], touch] of [['phone', [368, 812], true], ['laptop', [1440, 900], false]]) {
      const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2, hasTouch: touch })
      await ctx.addInitScript((t) => localStorage.setItem('nest.theme', t), theme)
      const page = await ctx.newPage()
      await page.goto(base + '/ops', { waitUntil: 'networkidle' })
      await page.evaluate(() => document.fonts.ready)
      for (const mode of ['day', 'week']) {
        const tag = `${size}-${theme}-${mode}`
        const check = (n, pass) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${tag} ${n}`); if (!pass) failed.push(`${tag} ${n}`) }
        if (mode === 'week') await page.getByRole('tab', { name: /^week$/i }).first().click()
        await page.waitForTimeout(400)
        if (size === 'phone') {
          const over = await page.evaluate(() => {
            const bad = []
            for (const e of document.querySelectorAll('.ops .pn, .ops')) if (e.scrollWidth > e.clientWidth) bad.push(`${e.className.toString().slice(0, 30)} ${e.scrollWidth}>${e.clientWidth}`)
            const d = document.documentElement
            if (d.scrollWidth > d.clientWidth) bad.push(`page ${d.scrollWidth}>${d.clientWidth}`)
            return bad
          })
          check(`no horizontal overflow ${over.join('; ')}`, over.length === 0)
          if (mode === 'day') {
            const first = await page.locator('.ops .pn').first().getAttribute('class')
            check(`Now card first (${first})`, /ops-nowc/.test(first ?? ''))
            const opens = await page.locator('.ops-nowc a:has-text("OPEN")').count()
            const lives = await page.locator('.ops-nowc .ops-nl').count()
            check(`OPEN per live row (${opens}/${lives})`, lives > 0 && opens > 0)
            check('free slot row', (await page.locator('.ops-nowc .ops-free').count()) > 0)
          }
        }
        await page.screenshot({ path: `${out}/${tag}.png`, fullPage: true })
        console.log(`${out}/${tag}.png`)
      }
      await ctx.close()
    }
  }
} finally {
  await browser.close()
  process.kill(-server.pid)
}
if (failed.length) { console.error(`${failed.length} failed`); process.exit(1) }
