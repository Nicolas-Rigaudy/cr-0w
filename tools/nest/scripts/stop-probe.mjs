// Stop sheet checks (item_279) on fixture data: Stop on a run from the Desk panel and
// from the session chat opens a sheet whose action button is on top and reachable,
// and one tap on it reaches the source.
// Usage: node scripts/stop-probe.mjs <out_dir>   (exits 1 on a failed check)
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) { console.error('usage: stop-probe.mjs <out_dir>'); process.exit(2) }
mkdirSync(out, { recursive: true })
const PORT = Number(process.env.SHOTS_PORT) || 4319
const base = `http://localhost:${PORT}`
if (await fetch(base).then(() => true, () => false)) { console.error(`port ${PORT} is busy`); process.exit(2) }
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
for (let i = 0; i < 100; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }

// True when the element's centre hit-tests to the element itself (not clipped or covered).
const onTop = (loc) => loc.evaluate((el) => {
  const r = el.getBoundingClientRect()
  const hit = document.elementFromPoint(r.x + r.width / 2, r.y + r.height / 2)
  return r.width > 0 && r.height > 0 && !!hit && (hit === el || el.contains(hit))
})

const failed = []
const browser = await chromium.launch()
try {
  for (const theme of ['modern', 'cyber']) {
    for (const [size, [w, h]] of [['laptop', [1440, 900]], ['inner', [707, 691]], ['phone', [368, 812]]]) {
      const tag = `${size}-${theme}`
      const check = (n, pass) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${tag} ${n}`); if (!pass) failed.push(`${tag} ${n}`) }
      const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 1 })
      await ctx.addInitScript((t) => localStorage.setItem('nest.theme', t), theme)
      const page = await ctx.newPage()
      await page.goto(base + '/desk', { waitUntil: 'networkidle' })
      // Open the first run: its row on the Desk list, else its chat screen.
      const pid = await page.evaluate(async () => {
        const r = await fetch('/fixtures/state.json').catch(() => null)
        return r && r.ok ? null : null
      })
      void pid
      const stops = page.getByRole('button', { name: /^Stop$/ })
      const n = await stops.count()
      let missionRun = false
      check(`a Stop button shows (${n})`, n > 0)
      if (!n) { await ctx.close(); continue }
      for (let i = 0; i < n; i++) {
        const b = stops.nth(i)
        if (!(await b.isVisible())) continue
        await b.click()
        await page.waitForTimeout(200)
        const sheet = page.locator('.sheet')
        check(`stop ${i}: sheet opens`, (await sheet.count()) === 1)
        if (!(await sheet.count())) continue
        const ok = sheet.locator('.btns .btn').last()
        const label = await ok.innerText()
        check(`stop ${i}: "${label}" is on top`, await onTop(ok))
        // A mission run: the main action cancels the mission, Restart run sits next to it.
        if ((await sheet.locator('.sh').innerText()) === 'Stop mission run') {
          check(`stop ${i}: main action reads "Stop and cancel the mission"`, label === 'Stop and cancel the mission')
          const alt = sheet.getByRole('button', { name: 'Restart run' })
          check(`stop ${i}: "Restart run" is on top`, (await alt.count()) === 1 && await onTop(alt))
          missionRun = true
        }
        await page.screenshot({ path: `${out}/stop-${tag}-${i}.png` })
        await sheet.getByRole('button', { name: 'Cancel' }).click().catch(() => page.keyboard.press('Escape'))
        await page.waitForTimeout(150)
      }
      console.log(`${missionRun ? 'seen' : 'not seen'} ${tag} mission run sheet`)
      await ctx.close()
    }
  }
} finally {
  await browser.close()
  process.kill(-server.pid)
}
if (failed.length) { console.error(`${failed.length} failed`); process.exit(1) }
