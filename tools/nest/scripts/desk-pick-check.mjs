// Headless check of the Desk pick (item_210): default, Ops and back, reload on a run.
// Usage: SHOTS_PORT=4391 node scripts/desk-pick-check.mjs   (runs vite dev on fixtures)
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'

const PORT = Number(process.env.SHOTS_PORT) || 4391
const base = `http://localhost:${PORT}`
if (await fetch(base).then(() => true, () => false)) { console.error(`port ${PORT} is busy`); process.exit(2) }
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
for (let i = 0; i < 100; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }

const failed = []
const check = (name, pass, got) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${name}${got ? ` (${got})` : ''}`); if (!pass) failed.push(name) }
const browser = await chromium.launch()
try {
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } })
  const picked = async () => { await page.waitForTimeout(800); return (await page.locator('button.row[aria-pressed="true"]').first().innerText().catch(() => '')).replace(/\s+/g, ' ').trim() }
  const runRow = page.locator('button.row.edge-w')
  await page.goto(base + '/desk', { waitUntil: 'networkidle' })
  await page.evaluate(() => localStorage.removeItem('nest.desk.pick'))
  await page.reload({ waitUntil: 'networkidle' })
  const first = await picked()
  const runs = await runRow.allInnerTexts().then((t) => t.map((x) => x.replace(/\s+/g, ' ').trim()))
  check('cleared storage picks something', !!first, first)
  check('cleared storage does not pick a run', !runs.includes(first), first)

  // Pick another session row (not a run), go to Ops and back.
  const rows = page.locator('button.row[aria-pressed]:not(.edge-w)')
  const n = await rows.count()
  let other = ''
  for (let i = 0; i < n; i++) {
    const t = (await rows.nth(i).innerText()).replace(/\s+/g, ' ').trim()
    if (t && t !== first) { await rows.nth(i).click(); other = t; break }
  }
  check('picked another session', !!other && (await picked()) === other, other)
  await page.goto(base + '/ops', { waitUntil: 'networkidle' })
  await page.goBack({ waitUntil: 'networkidle' })
  check('Ops and back keeps the pick', (await picked()) === other, await picked())

  // Pick a run, reload.
  if (await runRow.count()) {
    await runRow.first().click()
    const run = await picked()
    await page.reload({ waitUntil: 'networkidle' })
    check('reload on a run keeps the run', (await picked()) === run, await picked())
  } else check('fixtures have a run', false)
} finally {
  await browser.close()
  try { process.kill(-server.pid) } catch {}
}
process.exit(failed.length ? 1 : 0)
