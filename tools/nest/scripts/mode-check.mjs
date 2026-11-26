// Headless check of the light/dark mode engine (item_303) on the fixture dev server.
// Usage: node scripts/mode-check.mjs [out_dir]
//   AC1: mode auto follows a system switch live, no reload.
//   AC2: a forced mode survives a reload and is on <html> before any script of the app runs.
//   Settings: the Mode control previews, Save keeps it per device, Discard drops it.
// out_dir gets a shot of the Look pane per size. Exits 1 on a failure. Never opens a window.
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (out) mkdirSync(out, { recursive: true })
const PORT = Number(process.env.SHOTS_PORT) || 4318
if (await fetch(`http://localhost:${PORT}`).then(() => true, () => false)) {
  console.error(`port ${PORT} is busy. Set SHOTS_PORT to a free port.`)
  process.exit(2)
}
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 100; i++) {
  try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) }
}

const failed = []
const check = (name, pass, got) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${name}${pass ? '' : ` (got ${got})`}`); if (!pass) failed.push(name) }
const browser = await chromium.launch()
try {
  const ctx = await browser.newContext({ colorScheme: 'dark', viewport: { width: 1440, height: 900 } })
  const page = await ctx.newPage()
  // The mode <html> carries when the first byte of body parses, before the app module runs.
  await page.addInitScript(() => {
    document.addEventListener('readystatechange', () => {
      if (document.readyState === 'interactive' && !window.__firstMode) window.__firstMode = document.documentElement.dataset.mode
    })
  })
  const mode = () => page.evaluate(() => document.documentElement.dataset.mode)
  const first = () => page.evaluate(() => window.__firstMode)

  // AC1: auto follows the system live.
  await page.goto(`${base}/`)
  await page.evaluate(() => localStorage.removeItem('nest.mode'))
  await page.reload()
  await page.waitForSelector('#root *')
  check('auto paints dark on a dark system', (await mode()) === 'dark', await mode())
  await page.emulateMedia({ colorScheme: 'light' })
  await page.waitForFunction(() => document.documentElement.dataset.mode === 'light', null, { timeout: 2000 }).catch(() => {})
  check('auto repaints light live on a system switch', (await mode()) === 'light', await mode())
  await page.emulateMedia({ colorScheme: 'dark' })
  await page.waitForFunction(() => document.documentElement.dataset.mode === 'dark', null, { timeout: 2000 }).catch(() => {})
  check('auto repaints dark live on the switch back', (await mode()) === 'dark', await mode())

  // AC2: a forced mode survives a reload with no flash of the other one.
  await page.evaluate(() => localStorage.setItem('nest.mode', 'light'))
  await page.reload()
  await page.waitForSelector('#root *')
  check('forced light is on <html> at first paint on a dark system', (await first()) === 'light', await first())
  check('forced light holds after the app starts', (await mode()) === 'light', await mode())
  await page.emulateMedia({ colorScheme: 'light' })
  await page.emulateMedia({ colorScheme: 'dark' })
  await page.waitForTimeout(200)
  check('forced light ignores a system switch', (await mode()) === 'light', await mode())
  const bar = await page.evaluate(() => document.querySelector('meta[name="theme-color"]')?.getAttribute('content'))
  const bg = await page.evaluate(() => getComputedStyle(document.documentElement).getPropertyValue('--bg').trim())
  check('browser bar takes the painted --bg', bar === bg, `${bar} vs ${bg}`)

  // Settings: preview, Discard, Save.
  await page.evaluate(() => localStorage.removeItem('nest.mode'))
  await page.goto(`${base}/settings?pane=look`)
  const seg = page.locator('[role="group"][aria-label="Mode"]')
  await seg.waitFor()
  check('Settings shows Auto picked', (await seg.locator('button.on').innerText()) === 'Auto', await seg.locator('button.on').innerText())
  await seg.getByRole('button', { name: 'Light' }).click()
  check('Light previews at once', (await mode()) === 'light', await mode())
  check('the preview is not kept yet', (await page.evaluate(() => localStorage.getItem('nest.mode'))) === null)
  await page.getByRole('button', { name: 'Discard' }).click()
  check('Discard paints the kept mode again', (await mode()) === 'dark', await mode())
  await seg.getByRole('button', { name: 'Light' }).click()
  if (out) await page.screenshot({ path: `${out}/settings-mode-laptop.png` })
  await page.getByRole('button', { name: 'Save' }).click()
  await page.waitForTimeout(300)
  check('Save keeps light on the device', (await page.evaluate(() => localStorage.getItem('nest.mode'))) === 'light')
  await page.goto(`${base}/`)
  await page.waitForSelector('#root *')
  check('leaving Settings keeps the saved light', (await mode()) === 'light', await mode())
  // Header theme popup: the quick mode toggle keeps the pick at once.
  await page.locator('.theme-menu > button').click()
  const quick = page.locator('.theme-pop [role="group"][aria-label="Mode"]')
  await quick.getByRole('button', { name: 'Dark' }).click()
  check('the popup toggle paints dark', (await mode()) === 'dark', await mode())
  check('the popup toggle keeps dark on the device', (await page.evaluate(() => localStorage.getItem('nest.mode'))) === 'dark')
  if (out) await page.screenshot({ path: `${out}/theme-pop-laptop.png`, clip: { x: 1040, y: 0, width: 400, height: 360 } })
  await page.evaluate(() => localStorage.removeItem('nest.mode'))

  if (out) {
    const phone = await browser.newContext({ colorScheme: 'dark', viewport: { width: 368, height: 812 }, deviceScaleFactor: 2 })
    const p = await phone.newPage()
    await p.goto(`${base}/settings?pane=look`)
    await p.locator('[role="group"][aria-label="Mode"]').waitFor()
    await p.screenshot({ path: `${out}/settings-mode-cover.png`, fullPage: true })
  }
} finally {
  await browser.close()
  try { process.kill(-server.pid) } catch { server.kill() }
}
if (failed.length) { console.error(`${failed.length} failed`); process.exit(1) }
console.log('all passed')
