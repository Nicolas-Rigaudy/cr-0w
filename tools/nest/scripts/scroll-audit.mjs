// Scroll audit: every Nest screen at phone and laptop size, with long content forced in.
// Fails when content is cut off with no way to scroll to it.
// Usage: node scripts/scroll-audit.mjs [--only=reports,desk]
// Runs vite dev on fixtures (SHOTS_PORT, default 4318). Never opens a visible window.
//
// Each text block of the page (#main; the rail has its own scrolling log) is repeated 30 times, so every list, reader and pane overflows.
// Then each element that clips (overflow hidden or clip) while its content is taller than it,
// and is not a one-line ellipsis, is reported, as is a page taller than the window with no scroll.
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'

const only = process.argv.find((a) => a.startsWith('--only='))?.slice(7).split(',').filter(Boolean)
const SCREENS = {
  home: '/', missions: '/missions', detail: '/missions/20261109-build-inbox-filters', inbox: '/inbox',
  review: '/inbox?item=20261109-fix-push-token-refresh/v1', desk: '/desk', mrs: '/mrs', agenda: '/agenda',
  'mrs-review': '/mrs?mr=242', reports: '/reports', settings: '/settings?pane=look', usage: '/usage/acct1',
}
const SIZES = { phone: [368, 812], inner: [707, 691], laptop: [1440, 900] }
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

// Runs in the page: grow the content, then list what is cut off.
function audit() {
  const blocks = [...document.querySelectorAll('#main :is(p, li, .rrow, tr, .md > *, .tl)')]
  for (const b of blocks.slice(0, 400)) {
    const n = b.parentElement
    if (!n) continue
    for (let i = 0; i < 30; i++) n.appendChild(b.cloneNode(true))
  }
  const out = []
  const scrolls = (el) => /(auto|scroll)/.test(getComputedStyle(el).overflowY)
  const scrollableAncestor = (el) => {
    for (let p = el.parentElement; p; p = p.parentElement) if (scrolls(p) && p.scrollHeight > p.clientHeight + 2) return true
    return document.scrollingElement.scrollHeight > innerHeight + 2 && getComputedStyle(document.documentElement).overflowY !== 'hidden'
      && getComputedStyle(document.body).overflowY !== 'hidden'
  }
  for (const el of document.querySelectorAll('#main, #main *')) {
    const s = getComputedStyle(el)
    if (!/(hidden|clip)/.test(s.overflowY) || el.clientHeight < 40) continue
    if (s.whiteSpace === 'nowrap' || s.textOverflow === 'ellipsis' || s.webkitLineClamp !== 'none') continue
    if (el.closest('canvas, svg, .hero, [aria-hidden="true"]')) continue
    const lost = el.scrollHeight - el.clientHeight
    if (lost > 24 && !el.querySelector(':scope [style*="overflow"]')) {
      const cls = typeof el.className === 'string' ? el.className.trim().split(/\s+/).slice(0, 3).join('.') : ''
      out.push(`${el.tagName.toLowerCase()}${cls ? '.' + cls : ''} clips ${lost}px`)
    }
  }
  // Content below the window with neither the page nor any ancestor scrolling.
  const last = [...document.querySelectorAll('#main :is(p, li, .rrow, tr, .md > *)')].filter((e) => e.getClientRects().length && !e.closest('.vh, [aria-hidden="true"]')).at(-1)
  if (last) {
    const r = last.getBoundingClientRect()
    if (r.top > innerHeight && !scrollableAncestor(last)) out.push(`last block at y=${Math.round(r.top)} unreachable`)
  }
  return out
}

let bad = 0
const browser = await chromium.launch({ headless: true })
try {
  for (const [size, [w, h]] of Object.entries(SIZES)) {
    const ctx = await browser.newContext({ viewport: { width: w, height: h }, timezoneId: 'Europe/Paris' })
    await ctx.addInitScript(() => { localStorage.setItem('nest.rail', 'folded'); localStorage.setItem('nest.rail.home', 'folded') })
    const page = await ctx.newPage()
    for (const [name, path] of Object.entries(SCREENS)) {
      if (only && !only.includes(name)) continue
      await page.goto(base + path, { waitUntil: 'networkidle' })
      await page.waitForSelector('html:not([data-boot])', { timeout: 15000 }).catch(() => {})
      await page.waitForTimeout(500)
      if (name === 'reports') await page.locator('.rrow').first().click().catch(() => {})
      const found = await page.evaluate(audit)
      for (const f of found) console.log(`${size} ${name}: ${f}`)
      bad += found.length
    }
    await ctx.close()
  }
} finally {
  await browser.close()
  process.kill(-server.pid)
}
console.log(bad ? `${bad} clipped spots` : 'no clipped content')
process.exit(bad ? 1 : 0)
