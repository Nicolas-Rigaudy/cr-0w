// Headless check of the Desk chat: the picked session's chat scrolls inside its pane,
// its question dock and reply box stay on screen. Usage: node scripts/desk-chat-check.mjs [out.png]
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'

const PORT = Number(process.env.SHOTS_PORT) || 4392
const base = `http://localhost:${PORT}`
if (await fetch(base).then(() => true, () => false)) { console.error(`port ${PORT} is busy`); process.exit(2) }
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname, detached: true })
for (let i = 0; i < 100; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }

const failed = []
const check = (name, pass, got) => { console.log(`${pass ? 'PASS' : 'FAIL'} ${name}${got != null ? ` (${got})` : ''}`); if (!pass) failed.push(name) }
const browser = await chromium.launch()
try {
  for (const [w, h] of [[1440, 900], [1920, 1080]]) {
    const page = await browser.newPage({ viewport: { width: w, height: h } })
    await page.goto(base + '/desk', { waitUntil: 'networkidle' })
    await page.locator('button.mrow.unit', { hasText: 'footer links' }).first().click()
    await page.waitForTimeout(1200)
    const m = await page.evaluate(() => {
      const r = (s) => document.querySelector(s)?.getBoundingClientRect()
      const chat = document.querySelector('.dk .chat')
      return {
        page: document.scrollingElement.scrollHeight - innerHeight,
        tab: document.querySelector('.mtab.on')?.textContent,
        ask: r('.dk .ask-q')?.bottom, input: r('.dk .chat-cv form.in')?.bottom, vh: innerHeight,
        scrolls: chat ? chat.scrollHeight > chat.clientHeight : null,
        opt: getComputedStyle(document.querySelector('.dk .ask-q .oc') ?? document.body).display,
      }
    })
    const tag = `${w}x${h}`
    check(`${tag} opens on the chat`, m.tab === 'Transcript', m.tab)
    check(`${tag} question dock on screen`, m.ask != null && m.ask <= m.vh, `${Math.round(m.ask)} / ${m.vh}`)
    check(`${tag} reply box on screen`, m.input != null && m.input <= m.vh, `${Math.round(m.input)} / ${m.vh}`)
    check(`${tag} options styled as cards`, m.opt === 'flex', m.opt)
    check(`${tag} chat scrolls in its pane`, m.scrolls === true, m.scrolls)
    if (process.argv[2] && w === 1440) await page.screenshot({ path: process.argv[2] })
    await page.close()
  }
} finally {
  await browser.close()
  process.kill(-server.pid)
}
process.exit(failed.length ? 1 : 0)
