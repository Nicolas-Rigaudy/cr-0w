// Render the Desk and Ops sketches (item_293) (scripts/mockups/desk-ops.html) to PNG.
// Usage: node tools/nest/scripts/desk-ops-mockups.mjs <out_dir> [scene,scene]
// Writes <desk|ops|desk-phone>-<A|B|C>.png
// and one contact sheet per scenario (options side by side, same framing).
import { mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) {
  console.error('usage: node tools/nest/scripts/desk-ops-mockups.mjs <out_dir>')
  process.exit(2)
}
mkdirSync(out, { recursive: true })
const src = resolve(dirname(fileURLToPath(import.meta.url)), 'mockups/desk-ops.html')
const VS = ['A', 'B', 'C']
const ONLY = process.argv[3] ? process.argv[3].split(',') : null
const SCENES = [
  { name: 'desk', q: 'p=desk&f=laptop', vp: { width: 1440, height: 900 } },
  { name: 'ops', q: 'p=ops&f=laptop', vp: { width: 1440, height: 900 }, vs: ['A', 'B', 'C', 'X', 'Y'] },
  { name: 'desk-phone', q: 'p=desk&f=phone', vp: { width: 820, height: 900 } },
]

const browser = await chromium.launch({ headless: true })
for (const sc of SCENES.filter((s) => !ONLY || ONLY.includes(s.name))) {
  const page = await browser.newPage({ viewport: sc.vp, deviceScaleFactor: 1 })
  const files = []
  for (const v of sc.vs || VS) {
    await page.goto(`file://${src}?v=${v}&${sc.q}`)
    await page.waitForTimeout(200)
    const file = resolve(out, `${sc.name}-${v}.png`)
    await page.screenshot({ path: file, fullPage: true })
    files.push(file)
    console.log(file)
  }
  await page.close()
  const phone = sc.name.endsWith('phone')
  const sheet = await browser.newPage({ viewport: { width: phone ? 2560 : 1440 * 3 + 80, height: 900 } })
  const imgs = files.map((f) => `<img src="file://${encodeURI(f)}">`).join('')
  await sheet.setContent(`<body style="margin:0;background:#181926;display:flex;gap:20px;padding:20px;align-items:flex-start;width:max-content">${imgs}</body>`)
  await sheet.waitForTimeout(300)
  const file = resolve(out, `sheet-${sc.name}.png`)
  await sheet.screenshot({ path: file, fullPage: true })
  console.log(file)
  await sheet.close()
}
await browser.close()
