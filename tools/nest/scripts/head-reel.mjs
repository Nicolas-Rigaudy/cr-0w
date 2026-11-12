// Head beat reel (item_198): drives HeadCtl live through a scripted timeline on the GPU head alone
// (no data needed for run/fail) and writes an animated WebP plus the source WebM (for stills). Usage: node scripts/head-reel.mjs <out.webp> [cell css px, default 2.75]
// Timeline: load wave, long hover (think), present left, drag spin, run (table taps), fail (shiver, brace, droop), back to idle.
import { spawn, execFileSync } from 'node:child_process'
import { mkdirSync, readdirSync, renameSync, rmSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { chromium } from 'playwright-core'
const out = process.argv[2] ?? 'head-reel.webp', cw = Number(process.argv[3] ?? 2.75)
const PORT = 4397
const server = spawn('npx', ['vite', '--port', String(PORT), '--strictPort'], { stdio: 'ignore', cwd: new URL('..', import.meta.url).pathname })
const base = `http://localhost:${PORT}`
for (let i = 0; i < 80; i++) { try { await fetch(base); break } catch { await new Promise((r) => setTimeout(r, 200)) } }
const browser = await chromium.launch({ headless: true, args: ['--use-angle=gl', '--enable-gpu', '--ignore-gpu-blocklist'] })
const tmp = join(dirname(out), '.reel'); rmSync(tmp, { recursive: true, force: true }); mkdirSync(tmp, { recursive: true })
try {
  { const p = await browser.newPage(); await p.goto(base + '/', { waitUntil: 'networkidle' }); await p.close() }
  const ctx = await browser.newContext({ viewport: { width: 900, height: 900 }, recordVideo: { dir: tmp, size: { width: 900, height: 900 } } })
  const page = await ctx.newPage()
  await page.goto(base + '/', { waitUntil: 'networkidle' })
  const box = await page.evaluate(async (cw) => {
    await document.fonts.load("500 12px 'Fira Code'")
    const H = await import('/src/hero/gl/head.ts'), C = await import('/src/hero/headCtl.ts')
    document.body.innerHTML = ''; document.body.style.cssText = 'margin:0;background:#05060d;color:#8af;font:14px monospace'
    const label = document.createElement('div'), cv = document.createElement('canvas')
    label.style.cssText = 'position:absolute;left:12px;top:8px'; document.body.append(label, cv)
    const h = H.createHead(cv), l = h.layout({ cw, ch: cw * 7 / 4, font: cw * 6.65 / 4, dpr: 1 })
    cv.style.cssText = `position:absolute;left:${(900 - l.w) / 2}px;top:40px;width:${l.w}px;height:${l.h}px`
    const ctl = new C.HeadCtl(C.CHARACTER.lively)
    // [start ms, label, input]
    const beats = [
      [0, 'load: wave', {}],
      [2600, 'long hover: think', { hover: true, look: { yaw: 14, pitch: -3 } }],
      [5600, 'select left: present', { look: { yaw: -16, pitch: 2 }, present: true }],
      [7800, 'drag spin', { drag: 0.25 }], [8600, 'drag spin back', { drag: -0.3 }], [9400, 'release', {}],
      [10600, 'working: tap the table', { face: 'run' }],
      [14000, 'fail: glitch, shiver, brace, droop', { face: 'fail' }],
      [17600, 'back to idle', {}], [19600, 'end', {}],
    ]
    const t0 = performance.now(); let last = t0, bi = -1
    await new Promise((done) => {
      const tick = (now) => {
        const t = now - t0
        while (bi + 1 < beats.length && t >= beats[bi + 1][0]) { bi++; label.textContent = beats[bi][1]; if (beats[bi][2].present) ctl.present(t) }
        if (bi === beats.length - 1) return done()
        const b = beats[bi][2], dt = Math.min(now - last, 100); last = now
        const { frame } = ctl.step({ t, dt, glances: [], face: b.face ?? 'idle', tip: [46, 242, 255], look: b.look, hover: b.hover, drag: b.drag,
          eye: b.face === 'fail' ? [255, 46, 99] : b.face === 'run' ? [46, 242, 255] : undefined })
        h.render(frame, t / 1000)
        requestAnimationFrame(tick)
      }
      requestAnimationFrame(tick)
    })
    return { w: Math.round(l.w), h: Math.round(l.h) + 40 }
  }, cw)
  await ctx.close()
  const webm = readdirSync(tmp).find((f) => f.endsWith('.webm'))
  const crop = `crop=${Math.min(900, box.w + 260)}:${Math.min(900, box.h + 10)}:${Math.max(0, Math.round((900 - box.w - 260) / 2))}:0`
  execFileSync('ffmpeg', ['-y', '-loglevel', 'error', '-ss', '0.6', '-i', join(tmp, webm), '-vf', `${crop},fps=20`, '-loop', '0', '-quality', '80', '-c:v', 'libwebp_anim', out])
  renameSync(join(tmp, webm), out.replace(/\.webp$/, '.webm'))
  console.log('wrote', out, box)
} finally { rmSync(tmp, { recursive: true, force: true }); await browser.close(); server.kill() }
