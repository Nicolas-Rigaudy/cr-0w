// Render round 2 of the Nest Settings mockups (scripts/mockups/settings-v2.html) to PNG.
// Usage: node tools/nest/scripts/settings-v2-mockups.mjs <out_dir>
import { mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright-core'

const [out] = process.argv.slice(2)
if (!out) {
  console.error('usage: node tools/nest/scripts/settings-v2-mockups.mjs <out_dir>')
  process.exit(2)
}
const dir = resolve(out)
const src = resolve(dirname(fileURLToPath(import.meta.url)), 'mockups/settings-v2.html')
const SHOTS = [
  ['D', 'laptop', 'look'], ['D', 'laptop', 'pages'], ['E', 'laptop', 'look'],
  ['D', 'phone', 'list'], ['D', 'phone', 'look'], ['D', 'phone', 'pages'], ['E', 'phone', 'look'],
]
const SIZE = { laptop: { width: 1440, height: 900 }, phone: { width: 380, height: 900 } }

mkdirSync(dir, { recursive: true })
const browser = await chromium.launch({ headless: true })
for (const [v, f, p] of SHOTS) {
  const page = await browser.newPage({ viewport: SIZE[f] })
  await page.goto(`file://${src}?v=${v}&f=${f}&p=${p}`)
  await page.waitForTimeout(150)
  const file = resolve(dir, `settings-${v}-${f}${v === 'E' ? '' : `-${p}`}.png`)
  await page.screenshot({ path: file, fullPage: true })
  console.log(file)
  await page.close()
}
await browser.close()
