import type { UiState } from '../data/state'
import type { HeroModel } from './model'
import { elapsed } from '../format'

// Static ASCII scene, ported from the approved hero pass 7 renderer (item_139) with the
// option B head (item_126 shader port) pasted in from heads.json. Draws one frame per call.

export interface Heads { palette: string[]; moods: Record<string, { cols: number; glyphs: string[]; colours: string[] }> }
type RGB = [number, number, number]
export interface Palette { ring: RGB; dim: RGB; fg: RGB; accent: RGB; state: Record<UiState, RGB> }
// cx, cy, rad: projected planet disc for hit tests. a0: table angle without spin. back: behind the lens.
export interface PlanetLabel {
  x: number; y: number; name: string; state: UiState; side: 'l' | 'r'; note?: string; target?: boolean
  tag?: { x: number; y: number; text: string }; cx: number; cy: number; rad: number; a0: number; back: boolean; hidden?: boolean
}
export interface MoonOut { x: number; y: number; r: number; id: string; title: string; state: UiState }
export interface HeadBox { x: number; y: number; top: number; left: number; right: number }
export interface SceneOut { planets: PlanetLabel[]; head: HeadBox; lens: { x: number; y: number }; moons: MoonOut[]; more?: { x: number; y: number; n: number } }
// spin: degrees added to every planet angle (item_158). hover: lifted planet and its moons.
// look: head turn override (-1 left, 1 right, 0 front). lean: head tilt in degrees.
// pulse: 0..1 working pulse phase. lift: needs-you lift per project (0..1 of 4 px).
// glitch: failed glitch per project (0..1 left of the 180 ms).
export interface SceneOpts {
  now: Date; spin?: number; hover?: { name: string; k: number }
  look?: number; lean?: number; pulse?: number; sweep?: number; lift?: Record<string, number>; glitch?: Record<string, number>
}

// Cell size in px. The scene uses the head's own cell size, so the table, planets and moons
// keep the same ASCII definition as the head (item_158). Cells scale with the hero (item_194).
const CW0 = 4, CH0 = 7
let CW = CW0, CH = CH0, FONT = 6.65
const RAMP = '.,-~:;=!*#$@'
const KEYS = [...Array(94)].map((_, i) => String.fromCharCode(33 + i)).filter((c) => !'"\\'.includes(c)).join('')
const WHITE: RGB = [255, 240, 255]

interface V3 { x: number; y: number; z: number }
const V = (x: number, y: number, z: number): V3 => ({ x, y, z })
const add = (a: V3, b: V3) => V(a.x + b.x, a.y + b.y, a.z + b.z)
const sub = (a: V3, b: V3) => V(a.x - b.x, a.y - b.y, a.z - b.z)
const mul = (a: V3, s: number) => V(a.x * s, a.y * s, a.z * s)
const dot = (a: V3, b: V3) => a.x * b.x + a.y * b.y + a.z * b.z
const cross = (a: V3, b: V3) => V(a.y * b.z - a.z * b.y, a.z * b.x - a.x * b.z, a.x * b.y - a.y * b.x)
const norm = (a: V3) => mul(a, 1 / Math.sqrt(dot(a, a)))
const rotX = (p: V3, a: number) => { const c = Math.cos(a), s = Math.sin(a); return V(p.x, p.y * c - p.z * s, p.y * s + p.z * c) }
const clamp = (v: number, a: number, b: number) => Math.max(a, Math.min(b, v))
const lerp = (a: number, b: number, t: number) => a + (b - a) * t
const mix = (a: RGB, b: RGB, t: number): RGB => [lerp(a[0], b[0], t), lerp(a[1], b[1], t), lerp(a[2], b[2], t)]
const D2R = Math.PI / 180
const LIGHT = norm(V(-0.55, 0.75, -0.6))
const rampCh = (L: number) => RAMP[clamp(Math.floor(L * 12), 0, 11)]

function rng(seed: number) {
  let a = seed
  return () => {
    a |= 0; a = (a + 0x6d2b79f5) | 0
    let t = Math.imul(a ^ (a >>> 15), 1 | a)
    t = (t + Math.imul(t ^ (t >>> 7), 61 | t)) ^ t
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

type Cam = ReturnType<typeof camera>
function camera(pos: V3, tgt: V3, focal: number, cx: number, cy: number) {
  const f = norm(sub(tgt, pos)), r = norm(cross(V(0, 1, 0), f)), u = cross(f, r)
  return {
    pos, f, focal,
    proj(p: V3) { const v = sub(p, pos), z = dot(v, f); return { x: cx + (focal * dot(v, r)) / z, y: cy - (focal * dot(v, u)) / z, z } },
    ray(sx: number, sy: number) { return norm(add(f, add(mul(r, (sx - cx) / focal), mul(u, -(sy - cy) / focal)))) },
  }
}

// bg: optional cell backdrop alpha (sphere cells), so a sphere reads as one solid disc.
// k: cell scale vs the scene grid (head bakes at another density).
interface Grid { cols: number; rows: number; z: Float32Array; ch: (string | null)[]; col: RGB[]; a: Float32Array; bg?: Float32Array; k?: number }
function makeGrid(w: number, h: number): Grid {
  const cols = Math.ceil(w / CW), rows = Math.ceil(h / CH), n = cols * rows
  return { cols, rows, z: new Float32Array(n).fill(1e9), ch: new Array(n).fill(null), col: new Array(n), a: new Float32Array(n), bg: new Float32Array(n) }
}
function put(g: Grid, c: number, r: number, z: number, ch: string | null, col: RGB, a = 1, bg = 0) {
  if (c < 0 || r < 0 || c >= g.cols || r >= g.rows) return
  const i = r * g.cols + c
  if (z >= g.z[i]) return
  g.z[i] = z; g.ch[i] = ch; g.col[i] = col; g.a[i] = a
  if (g.bg) g.bg[i] = bg
}
function splat(g: Grid, cam: Cam, p: V3, ch: string, col: RGB, a = 1, zb = 0) {
  const s = cam.proj(p)
  if (s.z <= 0.1) return
  put(g, Math.floor(s.x / CW), Math.floor(s.y / CH), s.z - zb, ch, col, a)
}
function splatLit(g: Grid, cam: Cam, p: V3, n: V3, col: RGB, gain = 1, floor = 0.08, a = 1) {
  const L = clamp(Math.max(0, dot(n, LIGHT)) * gain + floor, 0, 1)
  splat(g, cam, p, rampCh(L), mix(mix(col, [20, 20, 50], 0.55), col, L), a * (0.45 + 0.55 * L))
}
function ringOnPlane(g: Grid, cam: Cam, c: V3, R: number, ch: string, col: RGB, a = 1, step = 0.012, dash = 0) {
  for (let t = 0; t < Math.PI * 2; t += step) {
    if (dash && Math.floor(t / dash) % 2) continue
    splat(g, cam, add(c, V(R * Math.cos(t), 0, R * Math.sin(t))), ch, col, a)
  }
}
function octa(g: Grid, cam: Cam, c: V3, s: number, col: RGB, gain = 1) {
  const rnd = rng(7)
  for (let f = 0; f < 8; f++) {
    const sx = f & 1 ? 1 : -1, sy = f & 2 ? 1 : -1, sz = f & 4 ? 1 : -1, n = norm(V(sx, sy * 0.8, sz))
    for (let k = 0; k < 420; k++) {
      let a = rnd(), b = rnd()
      if (a + b > 1) { a = 1 - a; b = 1 - b }
      splatLit(g, cam, add(c, V(sx * s * a, sy * s * 1.35 * b, sz * s * (1 - a - b))), n, col, gain, 0.25)
    }
  }
}
// tilt 0 lies flat on the table, PI/2 - 0.25 stands nearly upright.
function brokenRing(g: Grid, cam: Cam, c: V3, col: RGB, tilt = Math.PI / 2 - 0.25) {
  const R = 0.26, r = 0.06
  for (let i = 0; i < 220; i++) {
    const th = (i / 220) * Math.PI * 2
    if ((th > 0.5 && th < 1.0) || (th > 3.3 && th < 3.6)) continue
    for (let j = 0; j < 8; j++) {
      const ph = (j / 8) * Math.PI * 2
      const p = rotX(V((R + r * Math.cos(ph)) * Math.cos(th), r * Math.sin(ph), (R + r * Math.cos(ph)) * Math.sin(th)), tilt)
      const n = rotX(V(Math.cos(ph) * Math.cos(th), Math.sin(ph), Math.cos(ph) * Math.sin(th)), tilt)
      splatLit(g, cam, add(c, p), n, col, 1.4, 0.3)
    }
  }
}
interface Sphere { c: V3; r: number; col: RGB; emit: number; dimf: number }
function raySpheres(g: Grid, cam: Cam, spheres: Sphere[]) {
  const boxes = spheres.map((s) => { const p = cam.proj(s.c); return { p, rad: (s.r * cam.focal) / p.z } })
  for (let r = 0; r < g.rows; r++) for (let c = 0; c < g.cols; c++) {
    const sx = c * CW + CW / 2, sy = r * CH + CH / 2
    let d: V3 | null = null, best = 1e9, hs: Sphere | null = null
    spheres.forEach((s, i) => {
      const { p, rad } = boxes[i]
      if (Math.abs(sx - p.x) > rad + CW || Math.abs(sy - p.y) > rad + CH) return
      d ??= cam.ray(sx, sy)
      const oc = sub(cam.pos, s.c), b = dot(oc, d), h = b * b - (dot(oc, oc) - s.r * s.r)
      if (h < 0) return
      const t = -b - Math.sqrt(h)
      if (t > 0 && t < best) { best = t; hs = s }
    })
    if (!hs || !d) continue
    const s: Sphere = hs
    const n = norm(sub(add(cam.pos, mul(d, best)), s.c))
    const L = clamp(Math.max(0, dot(n, LIGHT)) * (1 - s.emit) + s.emit * 0.75 + 0.22, 0, 1) * s.dimf
    const z = best * dot(d, cam.f)
    put(g, c, r, z, rampCh(L), mix(mix(s.col, [20, 20, 50], 0.25), mix(s.col, WHITE, s.emit * 0.4), L), 0.7 + 0.3 * L, (0.14 + 0.16 * L) * s.dimf)
  }
}


// Holotable (item_162): a projector lens on a galaxy-map table, CR-0W floating above it.
// Planets = projects, colour = most urgent state.
// Sky ring (item_197): planets sit evenly on a ring parallel to the table, under the head. Height
// above the table = priority (failed, needs you, working, quiet, done just above the table). The highest
// sits at the front, the next ones alternate left and right; a spin shows the rest.
const RT = 4.6
const PY = -0.55
const RING = 0.8
const RANK: Record<UiState, number> = { failed: 0, needs: 1, working: 2, sleeping: 3, queued: 3, done: 4 }
const ALT: Record<UiState, number> = { failed: 0.95, needs: 0.62, working: 0.38, sleeping: 0.2, queued: 0.2, done: 0.08 }
// Size by priority, so the same state always reads the same (not by open missions).
const SIZE: Record<UiState, number> = { failed: 0.34, needs: 0.34, working: 0.3, sleeping: 0.26, queued: 0.26, done: 0.24 }
const MAXP = 9
// Cell scale (item_194): one rule for every size, calibrated on the phone cover hero (344x400 at
// 0.434). The smaller of the width and height fits wins, so head and scene keep their share of the hero.
// DENSITY: glyphs per head-scale cell (item_195: 1.25x so spheres and head read rounder). The
// head is baked at the same density (scripts/hero-heads.mjs); HEAD_COLS is its width at 1x, so
// a bake at another density keeps the head the same size.
const HS = 0.434, REF_W = 344, REF_H = 400, HEAD_COLS = 173, DENSITY = 1.25
export const fmtWait = (m: number) => (m < 60 ? `${m}m` : `${Math.floor(m / 60)}h${String(m % 60).padStart(2, '0')}`)

export function drawScene(cv: HTMLCanvasElement, w: number, h: number, dpr: number, m: HeroModel, heads: Heads | undefined, pal: Palette, turn: { l: Heads; r: Heads } | undefined, opts: SceneOpts): SceneOut {
  const cw = Math.round(w * dpr), ch = Math.round(h * dpr)
  if (cv.width !== cw || cv.height !== ch) { cv.width = cw; cv.height = ch }
  const cs = HS * Math.min(w / REF_W, h / REF_H)
  CW = (CW0 * cs) / DENSITY; CH = (CH0 * cs) / DENSITY; FONT = (6.65 * cs) / DENSITY
  const g = makeGrid(w, h), rnd = rng(42), SC = pal.state
  const spin = opts.spin ?? 0
  const POS = V(0, 4.4, -14.5), TGT = V(0, 0.2, 0), LENS = V(0, PY, 0)
  const focal = ((w / 2 - 28) * 14.5) / RT
  let cam = camera(POS, TGT, focal, w / 2, h / 2)
  cam = camera(POS, TGT, focal, w / 2, h / 2 + (h * 0.8 - cam.proj(LENS).y))
  const at = (a: number, f: number, y = PY) => V(RT * f * Math.cos(a * D2R), y, -RT * f * Math.sin(a * D2R))

  // Table: faint rings, plus spokes that turn with the planets so a spin reads (item_158).
  for (const f of [0.25, 0.5, 0.75]) for (let t = 0; t < 360; t += 2.4 / f) splat(g, cam, at(t, f), '·', pal.ring, 0.2)
  for (let k = 0; k < 8; k++) for (let f = 0.2; f < 0.92; f += 0.05) splat(g, cam, at(k * 45 + 22.5 + spin, f), '·', pal.ring, 0.1)

  // Rim (item_197, the 24h clock is gone): a bezel that turns with the planets, minor ticks every
  // 10 deg, long accent ticks every 45. Front is lit more, the back fades.
  const TC = mix(pal.ring, pal.fg, 0.55)
  const rimA = (t: number) => 0.4 + 0.45 * Math.max(0, Math.sin((((t % 360) + 360) % 360) * D2R))
  for (let t = 0; t < 360; t += 0.8) splat(g, cam, at(t, 1), '-', TC, rimA(t) * 0.8)
  for (let k = 0; k < 36; k++) {
    const t = k * 10 + spin, major = k % 9 === 0
    for (let f = 1.015; f < (major ? 1.11 : 1.045); f += 0.012) splat(g, cam, at(t, f), major ? '#' : "'", major ? pal.accent : TC, rimA(t) * (major ? 2.2 : 1.6))
  }
  // Scan sweep (idle): a bright head with a fading tail running round the rim.
  if (opts.sweep !== undefined) {
    const s0 = opts.sweep * 360
    for (let d = 0; d < 60; d += 0.5) for (const f of [0.98, 1, 1.02]) splat(g, cam, at(s0 - d, f), d < 3 ? '#' : d < 12 ? '=' : '-', mix(SC.working, WHITE, d < 3 ? 0.5 : 0), 1.6 * (1 - d / 60) ** 1.3)
  }

  // Lens: the projector core, and a sparse cone that carries the head up.
  for (let r = 0; r < 0.34; r += 0.03) for (let t = 0; t < 360; t += 4) splat(g, cam, add(LENS, V(r * Math.cos(t * D2R), 0.02, r * Math.sin(t * D2R))), r < 0.18 ? '@' : '*', mix(WHITE, SC.working, r * 2.5), 0.9)
  ringOnPlane(g, cam, add(LENS, V(0, 0.02, 0)), 0.6, '-', SC.working, 0.6, 0.02)
  for (let i = 0; i < 700; i++) {
    const y = Math.pow(rnd(), 1.4) * 3.2, rad = (0.2 + y * 0.42) * Math.sqrt(rnd()), t = rnd() * Math.PI * 2
    splat(g, cam, add(LENS, V(rad * Math.cos(t), y + 0.05, rad * Math.sin(t))), rnd() < 0.2 ? ':' : '.', mix(SC.working, pal.accent, rnd() * 0.3), 0.08 + 0.2 * (1 - y / 3.2))
  }

  // Planets.
  const hov = opts.hover
  const ranked = m.projects.slice(0, MAXP).map((p, i) => ({ p, i })).sort((a, b) => RANK[a.p.state] - RANK[b.p.state] || a.i - b.i).map((o) => o.p)
  const stepA = 360 / Math.max(ranked.length, 1)
  const ps = ranked.map((p, i) => {
    const a = 90 + (i % 2 ? 1 : -1) * Math.ceil(i / 2) * stepA
    let base = at(a + spin, RING)
    // Back half of the table (behind the lens): smaller and dimmer, labels hidden.
    const depth = clamp(base.z / (RT * 0.5), 0, 1)
    const back = base.z > RT * 0.15
    // Failed: one short sideways glitch.
    const gl = opts.glitch?.[p.name] ?? 0
    if (gl > 0) base = add(base, V((Math.random() - 0.5) * 0.5 * gl, 0, 0))
    const r = SIZE[p.state] * (1 - 0.25 * depth)
    const z0 = cam.proj(base).z, px = z0 / cam.focal
    const hk = hov?.name === p.name ? hov.k : 0
    const stalk = ALT[p.state] + (p.state === 'needs' ? 0.18 * Math.min((p.wait ?? 0) / 300, 1) : 0)
      + (opts.lift?.[p.name] ?? (p.state === 'needs' ? 1 : 0)) * 4 * px + hk * 0.3
    return { ...p, a0: a, base, r, stalk, depth, back, hk, c: add(base, V(0, stalk + r, 0)) }
  })
  // Rim footprints: each planet marks the bezel with an arc in its colour, brighter when hovered
  // and when it turns into the front, so a spin shows which one has the focus.
  for (const p of ps) {
    const t0 = p.a0 + spin, off = Math.abs(((((t0 - 90) % 360) + 540) % 360) - 180)
    const front = Math.max(0, 1 - off / 25), k = 1 + 0.6 * (1 - p.depth) + p.hk + 0.8 * front
    for (let d = -6; d <= 6; d += 0.5) for (const f of [0.985, 1.015]) splat(g, cam, at(t0 + d, f), '=', mix(SC[p.state], WHITE, 0.25 * front), k * (1 - Math.abs(d) / 8))
  }
  // Bot intent: longest wait first (model order), then a failure, then live work.
  const target = ps.find((p) => p.state === 'needs') ?? ps.find((p) => p.state === 'failed') ?? ps.find((p) => p.state === 'working')
  const spheres: Sphere[] = []
  const dimOthers = hov ? 1 - 0.3 * hov.k : 1
  const pulse = opts.pulse === undefined ? 0 : 0.5 - 0.5 * Math.cos(opts.pulse * Math.PI * 2)
  for (const p of ps) {
    const col = SC[p.state], flat = add(p.base, V(0, 0.02, 0))
    const fade = (p.hk > 0 ? 1 : dimOthers) * (1 - 0.55 * p.depth)
    ringOnPlane(g, cam, flat, p.r * 1.2, '.', mix(col, [20, 20, 50], 0.35), 0.5, 0.06)
    for (let y = 0.05; y < p.stalk; y += 0.07) splat(g, cam, add(p.base, V(0, y, 0)), ':', col, 0.55)
    const quiet = p.state === 'sleeping' || p.state === 'queued' || p.state === 'done'
    // One light for every live planet; the target stands out by its marker and the head's gaze.
    const emit = (quiet ? 0.12 : 0.45) + (p.state === 'working' ? 0.25 * pulse : 0) + 0.25 * p.hk
    spheres.push({ c: p.c, r: p.r, col: quiet && p.state !== 'done' ? pal.ring : p.state === 'failed' ? mix(col, WHITE, 0.2) : col, emit, dimf: (quiet ? 0.5 : 1) * fade })
    if (p.state === 'failed') brokenRing(g, cam, p.c, col)
    if (p === target && p.state === 'needs') octa(g, cam, add(p.c, V(0, p.r + 0.34, 0)), 0.12, col, 1.4)
  }
  // Hovered planet: its open missions orbit it as moons, max 6 then +n.
  const moons: { c: V3; r: number; id: string; title: string; state: UiState }[] = []
  let more: { c: V3; n: number } | undefined
  const hp = ps.find((p) => p.hk > 0)
  if (hp) {
    // Moons fly out as the hover eases in, leaving room for their titles. The orbit is
    // tilted toward the camera so the titles spread vertically instead of stacking.
    const list = hp.missions.slice(0, 6), R = (Math.max(hp.r, 0.2) * 1.9 + 0.12) * (1 + 0.5 * hp.hk)
    list.forEach((mi, i) => {
      const t = (i / Math.max(list.length, 1)) * Math.PI * 2 + spin * D2R * 0.5 + 0.4
      const c = add(hp.c, V(R * Math.cos(t), R * 0.55 * Math.sin(t), R * 0.8 * Math.sin(t)))
      const r = 0.21 * hp.hk
      moons.push({ c, r, id: mi.id, title: mi.title, state: mi.state })
      spheres.push({ c, r, col: SC[mi.state], emit: 0.3, dimf: hp.hk })
    })
    if (hp.missions.length > 6) more = { c: add(hp.c, V(0, hp.r + 0.25, 0)), n: hp.missions.length - 6 }
  }
  raySpheres(g, cam, spheres)

  // Head: option B glyph grid on its own layer at HS scale, above the lens, turned toward the
  // target. Its eyes (the violet cells) take the target's state colour.
  const lens = cam.proj(LENS)
  let head: HeadBox = { x: w / 2, y: lens.y, top: lens.y, left: w / 2, right: w / 2 }
  let layer: HeadLayer | undefined
  const tx = target ? cam.proj(target.c).x : w / 2
  const dir = opts.look ?? (target && Math.abs(tx - w / 2) > 40 ? Math.sign(tx - w / 2) : 0)
  const H = turn && dir ? (dir < 0 ? turn.l : turn.r) : heads
  if (H) {
    const hd = H.moods[m.mood], rows = hd.glyphs.length
    const hg: Grid = { cols: hd.cols, rows, z: new Float32Array(hd.cols * rows).fill(1e9), ch: new Array(hd.cols * rows).fill(null), col: new Array(hd.cols * rows), a: new Float32Array(hd.cols * rows), k: (HEAD_COLS * DENSITY) / hd.cols }
    const eye = target ? SC[target.state] : undefined
    hd.glyphs.forEach((row, r) => {
      const first = row.search(/\S/), last = row.trimEnd().length
      if (first < 0) return
      const eyes = eye && r > rows * 0.36 && r < rows * 0.54
      for (let c = first; c < last; c++) {
        const k = hd.colours[r].slice(2 * c, 2 * c + 2)
        if (row[c] === ' ' || k.length < 2 || k === '  ') { put(hg, c, r, 0, ' ', pal.dim, 0); continue }
        const [hex, alpha] = H.palette[KEYS.indexOf(k[0]) * KEYS.length + KEYS.indexOf(k[1])].split('/')
        const n = parseInt(hex.slice(1), 16)
        let col: RGB = [(n >> 16) & 255, (n >> 8) & 255, n & 255]
        if (eyes && col[0] > col[1] + 40 && col[2] > col[1]) col = mix(col, eye, 0.9)
        put(hg, c, r, 0, row[c], col, Math.min(1, Number(alpha) * 1.5))
      }
    })
    const hw = hd.cols * CW * hg.k!, hh = rows * CH * hg.k!
    const top = Math.max(60, (70 + cam.proj(at(270, 1)).y - 50 - hh) / 2)
    layer = { g: hg, scale: 1, x: w / 2 - hw / 2, y: top, zSplit: lens.z, lean: (opts.lean ?? 0) * D2R }
    head = { x: w / 2, y: top + hh / 2, top, left: w / 2 - hw / 2, right: w / 2 + hw / 2 }
  }

  paint(cv, g, dpr, layer)

  const gaps = new Map<PlanetLabel, number>()
  const out: SceneOut = {
    planets: [], head, lens: { x: lens.x, y: lens.y },
    moons: moons.map((o) => { const s = cam.proj(o.c); return { x: s.x, y: s.y, r: Math.max(6, (o.r * cam.focal) / s.z), id: o.id, title: o.title, state: o.state } }),
    more: more && { ...cam.proj(more.c), n: more.n },
  }
  for (const p of ps) {
    const s = cam.proj(p.c)
    const rad = (p.r * cam.focal) / s.z
    const side = s.x < w / 2 ? 'l' : 'r'
    const note = p.state === 'working' && p.runFrom ? elapsed(p.runFrom, opts.now)
      : p.state === 'failed' ? `${p.cause ?? 'run'} failed`
      : p.state === 'done' || p.state === 'sleeping' || p.state === 'queued' ? p.state === 'sleeping' ? 'sleep' : p.state : undefined
    const label: PlanetLabel = {
      x: s.x + (side === 'r' ? rad + 8 : -rad - 8), y: s.y, name: p.name, state: p.state, side, note, target: p === target,
      cx: s.x, cy: s.y, rad: Math.max(rad, 10), a0: p.a0, back: p.back && p.hk === 0,
    }
    if (p.state === 'needs' && p.wait !== undefined) {
      const t = cam.proj(add(p.c, V(0, p.r + (p === target ? 0.62 : 0.2), 0)))
      label.tag = { x: t.x, y: t.y, text: `YOUR MOVE · ${fmtWait(p.wait)}` }
    }
    out.planets.push(label)
    gaps.set(label, rad + 8)
  }
  // Nudge planet labels apart (approximate text boxes), keep tags clear of the head.
  type Box = { x0: number; x1: number; y0: number; y1: number; own?: string }
  const boxes: Box[] = []
  const line = (x0: number, x1: number, y: number): Box => ({ x0, x1, y0: y - 8, y1: y + 8 })
  for (const p of ps) {
    const c = cam.proj(p.c), rr = (p.r * cam.focal) / c.z
    boxes.push({ x0: c.x - rr, x1: c.x + rr, y0: c.y - rr, y1: c.y + rr, own: p.name })
  }
  for (const l of out.planets) if (l.tag) { const tw = l.tag.text.length * 8.6; boxes.push(line(l.tag.x - tw / 2, l.tag.x + tw / 2, l.tag.y)) }
  for (const l of [...out.planets].sort((a, b) => a.y - b.y)) {
    const tw = (l.name.length + (l.note ? l.note.length + 3 : 0)) * 7.2
    const span = () => (l.side === 'r' ? [l.x, l.x + tw] : [l.x - tw, l.x])
    const hit = (y: number) => { const [x0, x1] = span(); return boxes.some((b) => b.own !== l.name && b.y0 < y + 8 && b.y1 > y - 8 && b.x0 < x1 && x0 < b.x1) }
    // Blocked where it is: try the other side of the planet, then the nearest free row up or down.
    if (hit(l.y)) {
      const g = 2 * (gaps.get(l) ?? 0), was = { x: l.x, side: l.side }
      l.x += l.side === 'r' ? -g : g; l.side = l.side === 'r' ? 'l' : 'r'
      if (hit(l.y)) Object.assign(l, was)
    }
    const [x0, x1] = span()
    let dy = 0
    while (hit(l.y + dy) && hit(l.y - dy) && dy < 200) dy += 4
    if (hit(l.y + dy)) dy = -dy
    l.y += dy
    boxes.push(line(x0, x1, l.y))
  }
  // Labels hidden where the head covers them.
  const inHead = (x: number, y: number) => x > head.left + 20 && x < head.right - 20 && y > head.top && y < 2 * head.y - head.top
  // Readability (item_158): back labels hide (YOUR MOVE tags stay); nothing sits under the head.
  for (const l of out.planets) {
    const tw = (l.name.length + (l.note ? l.note.length + 3 : 0)) * 7.2, x0 = l.side === 'r' ? l.x : l.x - tw
    l.hidden = l.back || inHead(x0, l.y) || inHead(x0 + tw, l.y) || inHead(x0 + tw / 2, l.y)
    if (l.tag && inHead(l.tag.x, l.tag.y)) l.tag = undefined
  }
  return out
}

interface HeadLayer { g: Grid; scale: number; x: number; y: number; zSplit: number; lean: number }
let offscreen: HTMLCanvasElement | undefined
function paint(cv: HTMLCanvasElement, g: Grid, dpr: number, head?: HeadLayer) {
  const ctx = cv.getContext('2d')
  if (!ctx) return
  const L = (offscreen ??= document.createElement('canvas'))
  if (L.width !== cv.width || L.height !== cv.height) { L.width = cv.width; L.height = cv.height }
  const l = L.getContext('2d')!
  l.setTransform(1, 0, 0, 1, 0, 0); l.clearRect(0, 0, L.width, L.height)
  l.scale(dpr, dpr)
  l.textBaseline = 'middle'; l.textAlign = 'center'
  const CWg = CW, CHg = CH
  const cells = (gr: Grid, keep: (z: number) => boolean) => {
    const k = gr.k ?? 1, CW = CWg * k, CH = CHg * k
    l.font = `500 ${(FONT * k).toFixed(2)}px 'Fira Code', ui-monospace, monospace`
    for (let r = 0; r < gr.rows; r++) for (let c = 0; c < gr.cols; c++) {
      const i = r * gr.cols + c, ch = gr.ch[i]
      if (!ch || !keep(gr.z[i])) continue
      if (ch === ' ') { l.clearRect(c * CW - 0.5, r * CH - 0.5, CW + 1, CH + 1); continue }
      const col = gr.col[i], bg = gr.bg?.[i] ?? 0
      if (bg > 0) { l.fillStyle = `rgba(${col[0] | 0},${col[1] | 0},${col[2] | 0},${bg.toFixed(3)})`; l.fillRect(c * CW, r * CH, CW, CH) }
      l.fillStyle = `rgba(${col[0] | 0},${col[1] | 0},${col[2] | 0},${gr.a[i].toFixed(3)})`
      l.fillText(ch, c * CW + CW / 2, r * CH + CH / 2 + 0.5)
    }
  }
  // Behind the lens, then the head (its silhouette hides the table back), then the front.
  cells(g, (z) => !head || z > head.zSplit)
  if (head) {
    const hk = head.g.k ?? 1, hw = head.g.cols * CW * hk * head.scale, hh = head.g.rows * CH * hk * head.scale
    l.save(); l.translate(head.x + hw / 2, head.y + hh); l.rotate(head.lean); l.translate(-hw / 2, -hh); l.scale(head.scale, head.scale)
    cells(head.g, () => true)
    l.restore()
    cells(g, (z) => z <= head.zSplit)
  }
  ctx.clearRect(0, 0, cv.width, cv.height)
  ctx.save()
  ctx.globalCompositeOperation = 'lighter'
  ctx.filter = `blur(${14 * dpr}px)`; ctx.globalAlpha = 0.55; ctx.drawImage(L, 0, 0)
  ctx.filter = `blur(${3 * dpr}px)`; ctx.globalAlpha = 0.6; ctx.drawImage(L, 0, 0)
  ctx.restore()
  ctx.drawImage(L, 0, 0)
}
