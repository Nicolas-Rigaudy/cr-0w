// Ops queue checks (item_272) on fixture data: reorder by arrows, by drag (mouse on laptop,
// long-press on phone), the RUN ▾ menu, then start now past the cap (slots 3/2).
// Usage: node scripts/queue-shots.mjs <out_dir>   (exits 1 on a failed check)
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) { console.error('usage: queue-shots.mjs <out_dir>'); process.exit(2) }
mkdirSync(out, { recursive: true })
const PORT = Number(process.env.SHOTS_PORT) || 4318
const base = `http://localhost:${PORT}`
if (await fetch(base).then(() => true, () => false)) { console.error(`port ${PORT} is busy`); process.exit(2) }
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
for (let i = 0; i < 100; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }

const failed = []
const browser = await chromium.launch()
try {
  for (const theme of ['modern', 'cyber']) {
    for (const [size, [w, h], touch] of [['laptop', [1440, 900], false], ['phone', [368, 812], true]]) {
      const tag = `${size}-${theme}`
      const check = (n, pass) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${tag} ${n}`); if (!pass) failed.push(`${tag} ${n}`) }
      const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2, hasTouch: touch })
      await ctx.addInitScript((t) => localStorage.setItem('nest.theme', t), theme)
      const page = await ctx.newPage()
      await page.goto(base + '/ops', { waitUntil: 'networkidle' })
      await page.evaluate(() => document.fonts.ready)
      const rows = page.locator('.ops-q:has(.ops-grip)')
      const titles = () => rows.locator('.ops-qt').allInnerTexts()
      const queue = page.locator('.ops-q').first().locator('xpath=..')
      const shot = async (n) => {
        await page.waitForTimeout(250)
        const box = await queue.boundingBox()
        const file = `${out}/built-${n}-${tag}.png`
        await page.screenshot({ path: file, clip: { x: 0, y: Math.max(0, box.y - 60), width: w, height: Math.min(box.height + 260, h) } })
        console.log(file)
      }
      await rows.nth(1).evaluate((e) => e.scrollIntoView({ block: 'center' }))
      const start = await titles()
      check('3 mission rows', start.length === 3)
      await shot('queue')
      // ▼ on the first, then ▲ back.
      await rows.nth(0).getByRole('button', { name: /down$/ }).click()
      await page.waitForTimeout(200)
      check('down', (await titles())[1] === start[0])
      await rows.nth(1).getByRole('button', { name: /up$/ }).click()
      await page.waitForTimeout(200)
      check('up', (await titles())[0] === start[0])
      // Drag the last row to the top.
      const grip = rows.nth(2).locator('.ops-grip')
      const g = await grip.boundingBox()
      const top = await rows.nth(0).boundingBox()
      const [x, y] = [g.x + g.width / 2, g.y + g.height / 2]
      if (touch) {
        const cdp = await ctx.newCDPSession(page)
        const tp = (type, yy) => cdp.send('Input.dispatchTouchEvent', { type, touchPoints: type === 'touchEnd' ? [] : [{ x, y: yy }] })
        await tp('touchStart', y)
        await page.waitForTimeout(450)
        for (let yy = y; yy > top.y + 2; yy -= 8) await tp('touchMove', yy)
        await tp('touchMove', top.y + 2)
        await shot('drag')
        await tp('touchEnd')
      } else {
        await page.mouse.move(x, y)
        await page.mouse.down()
        for (let yy = y; yy > top.y + 2; yy -= 8) await page.mouse.move(x, yy)
        await page.mouse.move(x, top.y + 2)
        await shot('drag')
        await page.mouse.up()
      }
      await page.waitForTimeout(300)
      { const t = await titles(); check(`drag to top (${t.join(' / ')})`, t[0] === start[2]) }
      // RUN ▾ menu on the second row.
      await rows.nth(1).locator('.ops-run').click()
      const menu = page.locator('.ops-menu')
      check('menu open', (await menu.count()) === 1)
      check('Next in line first', (await menu.locator('button').first().innerText()).startsWith('Next in line'))
      await shot('menu')
      await page.keyboard.press('Escape')
      check('Escape closes', (await menu.count()) === 0)
      // Start now twice: the second goes past the cap.
      for (let i = 0; i < 2; i++) {
        await rows.nth(0).locator('.ops-run').click()
        await menu.locator('button').nth(1).click()
        await page.waitForTimeout(300)
      }
      const slots = await page.locator('.ops-slots').innerText()
      check(`slots over cap (${slots.replace(/\s+/g, ' ')})`, /slots 3\/2/i.test(slots) && /over cap/i.test(slots))
      check('1 row left', (await rows.count()) === 1)
      await shot('overcap')
      await ctx.close()
    }
  }
} finally {
  await browser.close()
  process.kill(-server.pid)
}
if (failed.length) { console.error(`${failed.length} failed`); process.exit(1) }
