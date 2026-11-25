// Render the Agenda / MRs / Reports page mockups (scripts/mockups/pages.html) to PNG.
// Usage: node tools/nest/scripts/pages-mockups.mjs <out_dir> [page...]   (pages: agenda mrs reports)
// Writes <out_dir>/<page>/<page>-<A|B|C>-<laptop|phone>.png and one contact sheet per page and form factor.
import { mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright-core'

const [out, ...want] = process.argv.slice(2)
if (!out) {
  console.error('usage: node tools/nest/scripts/pages-mockups.mjs <out_dir> [agenda|mrs|reports...]')
  process.exit(2)
}
const pages = want.length ? want : ['agenda', 'mrs', 'reports']
const src = resolve(dirname(fileURLToPath(import.meta.url)), 'mockups/pages.html')
const VS = ['A', 'B', 'C']
const SIZE = { laptop: { width: 1440, height: 900 }, phone: { width: 380, height: 900 } }

const browser = await chromium.launch({ headless: true })
for (const p of pages) {
  const dir = resolve(out, p)
  mkdirSync(dir, { recursive: true })
  for (const f of Object.keys(SIZE)) {
    const page = await browser.newPage({ viewport: SIZE[f] })
    for (const v of VS) {
      await page.goto(`file://${src}?p=${p}&v=${v}&f=${f}`)
      await page.waitForTimeout(150)
      const file = resolve(dir, `${p}-${v}-${f}.png`)
      await page.screenshot({ path: file, fullPage: true })
      console.log(file)
    }
    await page.close()
    // Contact sheet: the three options side by side, same framing.
    const sheet = await browser.newPage({ viewport: { width: f === 'phone' ? 1220 : 1440 * 3 + 80, height: 900 } })
    const imgs = VS.map((v) => `<img src="file://${resolve(dir, `${p}-${v}-${f}.png`)}" style="${f === 'phone' ? '' : 'width:1440px'}">`).join('')
    await sheet.setContent(`<body style="margin:0;background:#181926;display:flex;gap:20px;padding:20px;align-items:flex-start">${imgs}</body>`)
    await sheet.waitForTimeout(200)
    const file = resolve(dir, `${p}-sheet-${f}.png`)
    await sheet.screenshot({ path: file, fullPage: true })
    console.log(file)
    await sheet.close()
  }
}
await browser.close()
