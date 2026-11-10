// Headless screenshots of the Nest screens at the three board sizes, both themes.
// Usage: node scripts/shots.mjs <out_dir>   (run `npm run build` first)
// Never opens a visible window: uses chromium-headless-shell from playwright-core.
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) { console.error('usage: shots.mjs <out_dir>'); process.exit(2) }
mkdirSync(out, { recursive: true })

const PORT = 4317
const SIZES = { cover: [368, 812], inner: [707, 691], laptop: [1440, 900] }
const SCREENS = {
  home: '/',
  missions: '/missions',
  detail: '/missions/20261109-build-inbox-filters',
  inbox: '/inbox',
  desk: '/desk',
}

const server = spawn('npx', ['vite', 'preview', '--port', String(PORT), '--strictPort'], { stdio: 'ignore' })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 50; i++) {
  try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) }
}

const browser = await chromium.launch({ headless: true })
try {
  for (const theme of ['modern', 'cyber']) {
    for (const [size, [w, h]] of Object.entries(SIZES)) {
      const ctx = await browser.newContext({ viewport: { width: w, height: h }, deviceScaleFactor: 2, timezoneId: 'Europe/Paris' })
      await ctx.addInitScript((t) => localStorage.setItem('nest.theme', t), theme)
      const page = await ctx.newPage()
      for (const [name, path] of Object.entries(SCREENS)) {
        await page.goto(base + path, { waitUntil: 'networkidle' })
        await page.evaluate(() => document.fonts.ready)
        const file = `${out}/${name}-${size}-${theme}.png`
        await page.screenshot({ path: file })
        console.log(file)
      }
      await ctx.close()
    }
  }
} finally {
  await browser.close()
  server.kill()
}
