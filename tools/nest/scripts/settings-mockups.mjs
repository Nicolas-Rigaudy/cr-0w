// Render the Nest Settings page mockups (scripts/mockups/settings.html) to PNG.
// Usage: node tools/nest/scripts/settings-mockups.mjs <out_dir>
// Writes <out_dir>/settings-<A|B|C>-<laptop|phone>.png and one contact sheet per form factor.
import { mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright-core'

const [out] = process.argv.slice(2)
if (!out) {
  console.error('usage: node tools/nest/scripts/settings-mockups.mjs <out_dir>')
  process.exit(2)
}
const dir = resolve(out)
const src = resolve(dirname(fileURLToPath(import.meta.url)), 'mockups/settings.html')
const VS = ['A', 'B', 'C']
const SIZE = { laptop: { width: 1440, height: 900 }, phone: { width: 380, height: 900 } }

mkdirSync(dir, { recursive: true })
const browser = await chromium.launch({ headless: true })
for (const f of Object.keys(SIZE)) {
  const page = await browser.newPage({ viewport: SIZE[f] })
  for (const v of VS) {
    await page.goto(`file://${src}?v=${v}&f=${f}`)
    await page.waitForTimeout(150)
    const file = resolve(dir, `settings-${v}-${f}.png`)
    await page.screenshot({ path: file, fullPage: true })
    console.log(file)
  }
  await page.close()
  // Contact sheet: the three options side by side, same framing.
  const sheet = await browser.newPage({ viewport: { width: f === 'phone' ? 1220 : 1440 * 3 + 80, height: 900 } })
  const imgs = VS.map((v) => `<img src="file://${resolve(dir, `settings-${v}-${f}.png`)}" style="${f === 'phone' ? '' : 'width:1440px'}">`).join('')
  await sheet.setContent(`<body style="margin:0;background:#181926;display:flex;gap:20px;padding:20px;align-items:flex-start">${imgs}</body>`)
  await sheet.waitForTimeout(200)
  const file = resolve(dir, `settings-sheet-${f}.png`)
  await sheet.screenshot({ path: file, fullPage: true })
  console.log(file)
  await sheet.close()
}
await browser.close()
