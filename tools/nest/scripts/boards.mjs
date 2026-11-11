// Render design boards (*.dc.html) to PNG for side-by-side review with the app.
// Usage: node scripts/boards.mjs <boards_dir> <out_dir> <name>...   (name without .dc.html)
import { mkdirSync } from 'node:fs'
import { resolve } from 'node:path'
import { chromium } from 'playwright-core'

const [dir, out, ...names] = process.argv.slice(2)
if (!dir || !out || names.length === 0) {
  console.error('usage: node scripts/boards.mjs <boards_dir> <out_dir> <name>...')
  process.exit(2)
}
mkdirSync(out, { recursive: true })
const browser = await chromium.launch({ headless: true })
const page = await browser.newPage({ viewport: { width: 1500, height: 1000 } })
for (const n of names) {
  await page.goto(`file://${resolve(dir, `${n}.dc.html`)}`, { waitUntil: 'networkidle' }).catch(() => {})
  await page.waitForTimeout(500)
  await page.screenshot({ path: resolve(out, `${n}.png`), fullPage: true })
  console.log(resolve(out, `${n}.png`))
}
await browser.close()
