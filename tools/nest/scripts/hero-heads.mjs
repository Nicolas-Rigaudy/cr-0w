// Render the hero head spike once per mood and save its glyph grid as text data.
// Usage: node scripts/hero-heads.mjs <path/to/holo-head.html> [out=src/hero/heads.json]
// Turned heads for the gaze: HEAD_EXTRA="&yaw=-20" (heads-l.json) and "&yaw=20" (heads-r.json).
// Head at the hero's 1.25x glyph density (scene.ts DENSITY): w=863 h=800 focal=993 dy=160.
// Headless only. Reads the cell pass (glyph index + colour per cell) back from WebGL.
import { writeFileSync } from 'node:fs'
import { chromium } from 'playwright-core'

const [src, out = 'src/hero/heads.json'] = process.argv.slice(2)
if (!src) { console.error('usage: hero-heads.mjs <holo-head.html> [out]'); process.exit(2) }

const MOODS = ['idle', 'ask', 'fail', 'run']
const SIZE = 'holo=0&w=863&h=800&focal=993&dy=160' + (process.env.HEAD_EXTRA ?? '') // holo=0: plain option B, no emitter disc; HEAD_EXTRA e.g. '&yaw=-30'
const STEPS = 12 // colour levels per channel

const b = await chromium.launch({ headless: true, args: ['--use-angle=swiftshader', '--enable-unsafe-swiftshader'] })
const p = await b.newPage({ viewport: { width: 980, height: 900 } })
const heads = {}
for (const mood of MOODS) {
  await p.goto(`file://${src}?mood=${mood}&t=1&dpr=1&${SIZE}`)
  await p.waitForSelector('body[data-done]', { timeout: 20000 })
  heads[mood] = await p.evaluate(() => {
    /* global gl, F, COLS, ROWS, RAMP */
    const read = (i) => {
      gl.bindFramebuffer(gl.READ_FRAMEBUFFER, F.cell)
      gl.readBuffer(gl.COLOR_ATTACHMENT0 + i)
      const px = new Uint8Array(COLS * ROWS * 4)
      gl.readPixels(0, 0, COLS, ROWS, gl.RGBA, gl.UNSIGNED_BYTE, px)
      return Array.from(px)
    }
    return { cols: COLS, rows: ROWS, ramp: RAMP, col: read(0), aux: read(1) }
  })
}
await b.close()

// Encode: per row a glyph string (space = empty) and a colour string (two key chars per cell, '  ' = empty).
const KEYS = [...Array(94)].map((_, i) => String.fromCharCode(33 + i)).filter((c) => !'"\\'.includes(c)).join('')
const palette = []
const key = (r, g, bl, a) => {
  const q = (v) => Math.round((v / 255) * (STEPS - 1))
  const hex = [r, g, bl].map((v) => Math.round((q(v) * 255) / (STEPS - 1)).toString(16).padStart(2, '0')).join('')
  const al = Math.max(1, Math.round(a / 255 * 16)) / 16 // finer steps keep the faint shading cells
  const id = `#${hex}/${al}`
  let i = palette.indexOf(id)
  if (i < 0) { palette.push(id); i = palette.length - 1 }
  if (i >= KEYS.length ** 2) throw new Error('palette too large')
  return KEYS[Math.floor(i / KEYS.length)] + KEYS[i % KEYS.length]
}
const data = { palette, moods: {} }
for (const [mood, h] of Object.entries(heads)) {
  const glyphs = [], colours = []
  for (let y = h.rows - 1; y >= 0; y--) { // GL rows are bottom-up
    let gs = '', cs = ''
    for (let x = 0; x < h.cols; x++) {
      const i = (y * h.cols + x) * 4
      const a = h.col[i + 3]
      const g = h.aux[i]
      if (!h.aux[i + 3] || !a) { gs += ' '; cs += '  '; continue }
      gs += h.ramp[g] ?? ' '
      cs += key(h.col[i], h.col[i + 1], h.col[i + 2], a)
    }
    glyphs.push(gs.trimEnd()); colours.push(cs.slice(0, gs.trimEnd().length * 2))
  }
  while (glyphs.length && !glyphs[0]) { glyphs.shift(); colours.shift() }
  while (glyphs.length && !glyphs.at(-1)) { glyphs.pop(); colours.pop() }
  data.moods[mood] = { cols: h.cols, glyphs, colours }
}
writeFileSync(out, JSON.stringify(data))
console.log(out, palette.length, 'colours', Object.entries(data.moods).map(([m, d]) => `${m} ${d.glyphs.length} rows`).join(', '))
