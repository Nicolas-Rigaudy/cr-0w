// Render the MRs page round 2 mockups (scripts/mockups/mrs-v2.html) to PNG.
// Usage: node tools/nest/scripts/mrs-v2-mockups.mjs <out_dir>
// Writes mrs-<A|B|C>-laptop.png, mrs-<A|B|C>-laptop-patterns.png (patterns opened), mrs-<A|B|C>-phone.png
// and one contact sheet per scenario (options side by side, same framing).
import { mkdirSync } from 'node:fs'
import { dirname, resolve } from 'node:path'
import { fileURLToPath } from 'node:url'
import { chromium } from 'playwright-core'

const out = process.argv[2]
if (!out) {
  console.error('usage: node tools/nest/scripts/mrs-v2-mockups.mjs <out_dir>')
  process.exit(2)
}
mkdirSync(out, { recursive: true })
const src = resolve(dirname(fileURLToPath(import.meta.url)), 'mockups/mrs-v2.html')
const VS = ['A', 'B', 'C']
const SCENES = [
  { name: 'laptop', q: 'f=laptop', vp: { width: 1440, height: 900 } },
  { name: 'laptop-patterns', q: 'f=laptop&s=pat', vp: { width: 1440, height: 900 } },
  { name: 'phone', q: 'f=phone', vp: { width: 820, height: 900 } },
]

const browser = await chromium.launch({ headless: true })
for (const sc of SCENES) {
  const page = await browser.newPage({ viewport: sc.vp, deviceScaleFactor: 1 })
  const files = []
  for (const v of VS) {
    await page.goto(`file://${src}?v=${v}&${sc.q}`)
    await page.waitForTimeout(200)
    const file = resolve(out, `mrs-${v}-${sc.name}.png`)
    await page.screenshot({ path: file, fullPage: true })
    files.push(file)
    console.log(file)
  }
  await page.close()
  const phone = sc.name === 'phone'
  const sheet = await browser.newPage({ viewport: { width: phone ? 2560 : 1440 * 3 + 80, height: 900 } })
  const imgs = files.map((f) => `<img src="file://${f}">`).join('')
  await sheet.setContent(`<body style="margin:0;background:#181926;display:flex;gap:20px;padding:20px;align-items:flex-start;width:max-content">${imgs}</body>`)
  await sheet.waitForTimeout(300)
  const file = resolve(out, `mrs-sheet-${sc.name}.png`)
  await sheet.screenshot({ path: file, fullPage: true })
  console.log(file)
  await sheet.close()
}
await browser.close()
