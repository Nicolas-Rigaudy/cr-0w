// Headless main-thread cost per screen (item_249): 10 s idle and 10 s scrolling on the laptop
// layout, with the rail head on and stubbed (webgl2 refused). Reports the main-thread busy share
// (CDP TaskDuration over wall time), the script and layout/style split, frame intervals while
// scrolling, long tasks, and GL draw calls per second.
// Usage: node scripts/nest-perf.mjs [seconds=10] [screens=/,/desk,/ops,/inbox]
import { spawn } from 'node:child_process'
import { chromium } from 'playwright-core'
const SECS = Number(process.argv[2]) || 10
const SCREENS = (process.argv[3] || '/,/desk,/ops,/inbox').split(',')
const PORT = Number(process.env.PERF_PORT) || 4393
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
const bcdp = await browser.newBrowserCDPSession()

// Runs before the app: counts GL draws, long tasks and frame intervals; ?nohead refuses webgl2.
const probe = () => {
  const w = window
  w.__perf = { draws: 0, long: [], frames: [] }
  const proto = WebGL2RenderingContext.prototype
  for (const k of ['drawArrays', 'drawElements']) {
    const f = proto[k]
    proto[k] = function (...a) { w.__perf.draws++; return f.apply(this, a) }
  }
  if (location.search.includes('nohead')) {
    const g = HTMLCanvasElement.prototype.getContext
    HTMLCanvasElement.prototype.getContext = function (t, o) { return t === 'webgl2' ? null : g.call(this, t, o) }
  }
  new PerformanceObserver((l) => { for (const e of l.getEntries()) w.__perf.long.push(e.duration) }).observe({ type: 'longtask', buffered: true })
  // Frame intervals only while scrolling: an idle rAF probe would itself keep frames coming.
  let last = 0
  const tick = (t) => { if (last) w.__perf.frames.push(t - last); last = t; requestAnimationFrame(tick) }
  w.__perf.watch = () => requestAnimationFrame(tick)
}
// CPU seconds of every browser process (renderer, GPU, browser): GL work lands in the GPU process.
const cpu = async (bcdp) => (await bcdp.send('SystemInfo.getProcessInfo')).processInfo.reduce((s, p) => s + p.cpuTime, 0)

const metric = async (cdp) => Object.fromEntries((await cdp.send('Performance.getMetrics')).metrics.map((m) => [m.name, m.value]))
const pct = (v, s) => `${((100 * v) / s).toFixed(1)}%`
const q = (a, p) => { if (!a.length) return 0; const s = [...a].sort((x, y) => x - y); return s[Math.min(s.length - 1, Math.floor(p * s.length))] }

async function run(path, mode, nohead) {
  const ctx = await browser.newContext({ viewport: { width: 1440, height: 900 }, deviceScaleFactor: 2, timezoneId: 'Europe/Paris' })
  const page = await ctx.newPage()
  await page.addInitScript(probe)
  const cdp = await ctx.newCDPSession(page)
  await cdp.send('Performance.enable')
  // A focused window, as on the owner's laptop: the head and hero loops pause on blur.
  await cdp.send('Emulation.setFocusEmulationEnabled', { enabled: true })
  await page.goto(base + path + (nohead ? '?nohead' : ''), { waitUntil: 'networkidle' })
  await page.bringToFront()
  await page.waitForTimeout(1500)
  await page.evaluate(() => { const p = window.__perf; p.draws = 0; p.long = []; p.frames = [] })
  const c0 = await cpu(bcdp), a = await metric(cdp), t0 = Date.now()
  if (mode === 'scroll') {
    await page.evaluate(() => window.__perf.watch())
    await page.mouse.move(800, 500)
    const end = t0 + SECS * 1000
    for (let d = 1; Date.now() < end; d = -d) for (let i = 0; i < 12 && Date.now() < end; i++) { await page.mouse.wheel(0, 120 * d); await page.waitForTimeout(40) }
  } else await page.waitForTimeout(SECS * 1000)
  const b = await metric(cdp), c1 = await cpu(bcdp), wall = (Date.now() - t0) / 1000
  const p = await page.evaluate(() => ({ ...window.__perf, focus: document.hasFocus(), head: !!document.querySelector('.rail-head.on') }))
  await ctx.close()
  const d = (k) => b[k] - a[k]
  return {
    screen: path, mode, head: nohead ? 'stub' : p.head ? 'on' : 'absent', focus: p.focus,
    busy: pct(d('TaskDuration'), wall), script: pct(d('ScriptDuration'), wall), layout: pct(d('LayoutDuration') + d('RecalcStyleDuration'), wall),
    cpu: pct(c1 - c0, wall), p50: q(p.frames, 0.5).toFixed(1), p95: q(p.frames, 0.95).toFixed(1), long: p.long.length, longMax: Math.round(Math.max(0, ...p.long)),
    glDraws: (p.draws / wall).toFixed(0) + '/s',
  }
}

try {
  const gpu = (await bcdp.send('SystemInfo.getInfo')).gpu.devices.map((d) => d.deviceString).join(', ')
  console.log(`GPU: ${gpu || 'unknown'}, ${SECS} s per run, 1440x900 dpr 2`)
  const rows = []
  for (const s of SCREENS) for (const mode of ['idle', 'scroll']) for (const nohead of [false, true]) {
    const r = await run(s, mode, nohead)
    rows.push(r)
    console.log(JSON.stringify(r))
  }
  console.log('\n| Screen | Mode | Head | Main busy | Script | Layout+style | All-process CPU | Frame p50/p95 ms | Long tasks (max ms) | GL draws |')
  console.log('| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |')
  for (const r of rows) console.log(`| ${r.screen} | ${r.mode} | ${r.head} | ${r.busy} | ${r.script} | ${r.layout} | ${r.cpu} | ${r.mode === 'scroll' ? `${r.p50}/${r.p95}` : '-'} | ${r.long} (${r.longMax}) | ${r.glDraws} |`)
} finally { await browser.close(); server.kill() }
