// Headless frame budget check (item_220 AC5): the hero keeps the lit table on a normal run and falls
// back to the plain table when the CPU is throttled (CDP) so frames miss the budget. The
// nest.table = lit | plain override wins over the watch either way. Prints one line per case and ok.
//   PLAYWRIGHT_BROWSERS_PATH=<playwright cache> node scripts/hero-budget.mjs [throttle=20]
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const PORT = 4393, RATE = +(process.argv[2] ?? 20)
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
const cases = [
  { name: 'normal', rate: 1, want: 'lit' },
  { name: 'slow', rate: RATE, want: 'plain' },
  { name: 'slow, pref lit', rate: RATE, pref: 'lit', want: 'lit' },
  { name: 'normal, pref plain', rate: 1, pref: 'plain', want: 'plain' },
]
let ok = true
try {
  for (const c of cases) {
    const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 } })
    if (c.pref) await ctx.addInitScript((v) => localStorage.setItem('nest.table', v), c.pref)
    const page = await ctx.newPage()
    const cdp = await ctx.newCDPSession(page)
    await cdp.send('Emulation.setCPUThrottlingRate', { rate: c.rate })
    await page.goto(base + '/', { waitUntil: 'load' })
    await page.waitForSelector('.hero-box', { timeout: 60000 })
    // Idle ticks run at 15 fps; warm-up plus window is 35 frames, slower when throttled.
    let got = ''
    for (let i = 0; i < 40; i++) {
      await page.waitForTimeout(500)
      got = await page.$eval('.hero-box', (e) => e.dataset.table ?? '')
      if (got !== 'lit' || i >= (c.rate > 1 ? 39 : 8)) break
    }
    const pass = got === c.want
    ok &&= pass
    console.log(`${c.name}: table=${got} want=${c.want} ${pass ? 'ok' : 'FAIL'}`)
    await ctx.close()
  }
} finally { await browser.close(); server.kill() }
console.log(ok ? 'ok' : 'FAIL')
process.exit(ok ? 0 : 1)
