import type { UiState } from '../data/state'
import type { HeroModel } from './model'
import type { AgendaEvent } from '../data/types'
import { elapsed } from '../format'

// Static ASCII scene, ported from the approved hero pass 7 renderer (item_139) with the
// option B head (item_126 shader port) pasted in from heads.json. Draws one frame per call.

export interface Heads { palette: string[]; moods: Record<string, { cols: number; glyphs: string[]; colours: string[] }> }
type RGB = [number, number, number]
export interface Palette { ring: RGB; dim: RGB; fg: RGB; accent: RGB; state: Record<UiState, RGB> }
export interface PlanetLabel { x: number; y: number; name: string; state: UiState; side: 'l' | 'r'; note?: string; target?: boolean; tag?: { x: number; y: number; text: string } }
export interface RimLabel { x: number; y: number; text: string; kind: 'hour' | 'event'; side?: 'l' | 'r' }
export interface HeadBox { x: number; y: number; top: number; left: number; right: number }
export interface SceneOut { planets: PlanetLabel[]; rim: RimLabel[]; head: HeadBox; lens: { x: number; y: number } }
export interface SceneOpts { now: Date; agenda: AgendaEvent[]; spin?: number }

const CW = 4
const CH = 7
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

interface Grid { cols: number; rows: number; z: Float32Array; ch: (string | null)[]; col: RGB[]; a: Float32Array }
function makeGrid(w: number, h: number): Grid {
  const cols = Math.ceil(w / CW), rows = Math.ceil(h / CH), n = cols * rows
  return { cols, rows, z: new Float32Array(n).fill(1e9), ch: new Array(n).fill(null), col: new Array(n), a: new Float32Array(n) }
}
function put(g: Grid, c: number, r: number, z: number, ch: string | null, col: RGB, a = 1) {
  if (c < 0 || r < 0 || c >= g.cols || r >= g.rows) return
  const i = r * g.cols + c
  if (z >= g.z[i]) return
  g.z[i] = z; g.ch[i] = ch; g.col[i] = col; g.a[i] = a
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
    const L = clamp(Math.max(0, dot(n, LIGHT)) * (1 - s.emit) + s.emit * 0.75 + 0.08, 0, 1) * s.dimf
    put(g, c, r, best * dot(d, cam.f), rampCh(L), mix(mix(s.col, [20, 20, 50], 0.5), mix(s.col, WHITE, s.emit * 0.4), L), 0.4 + 0.6 * L)
  }
}


// Holotable (item_162): a projector lens on a galaxy-map table, CR-0W floating above it.
// Planets = projects (size = open missions, colour = most urgent state). A project that waits
// on you rises on a stalk (height = wait). Rim = 24h ruler: day so far lit, now notch, agenda.
const RT = 4.6
const PY = -0.55
// Free placement on the table: [angle (90 = facing the camera), fraction of the table radius].
// Urgent planets sit in front, working ones at the sides, quiet ones at the back.
const SPOTS: Record<'urgent' | 'working' | 'quiet', [number, number][]> = {
  urgent: [[138, 0.7], [64, 0.88], [100, 0.4], [168, 0.44], [14, 0.44]],
  working: [[248, 0.5], [206, 0.78], [340, 0.4], [196, 0.4]],
  quiet: [[296, 0.5], [322, 0.66], [236, 0.86], [304, 0.88], [270, 0.8]],
}
const MAXP = 9
// Head scale on the table: the full option B grid (692 px) would cover the whole map.
const HS = 0.62
const hourAngle = (h: number) => 90 + h * 15
const hoursOf = (d: Date) => d.getHours() + d.getMinutes() / 60
export const fmtWait = (m: number) => (m < 60 ? `${m}m` : `${Math.floor(m / 60)}h${String(m % 60).padStart(2, '0')}`)
const two = (n: number) => String(n).padStart(2, '0')

export function drawScene(cv: HTMLCanvasElement, w: number, h: number, dpr: number, m: HeroModel, heads: Heads | undefined, pal: Palette, turn: { l: Heads; r: Heads } | undefined, opts: SceneOpts): SceneOut {
  cv.width = Math.round(w * dpr); cv.height = Math.round(h * dpr)
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

  // Rim = 24h ruler, midnight at the front, running clockwise seen from above. Fixed: time, not space.
  const nowH = hoursOf(opts.now)
  for (let hh = 0; hh < 24; hh += 0.04) {
    const lit = hh <= nowH, front = Math.max(0, Math.sin(hourAngle(hh) * D2R))
    splat(g, cam, at(hourAngle(hh), 1), lit ? '=' : '-', lit ? SC.working : pal.ring, lit ? 0.95 : 0.18 + 0.14 * front)
  }
  for (let i = 0; i < 96; i++) {
    const big = i % 4 === 0
    for (let f = big ? 0.93 : 0.965; f < 0.995; f += 0.012) splat(g, cam, at(hourAngle(i / 4), f), big ? '|' : '.', pal.ring, big ? 0.55 : 0.3)
  }
  for (let f = 0.86; f < 1.05; f += 0.008) splat(g, cam, at(hourAngle(nowH), f), '#', WHITE, 0.95)
  const today = opts.agenda.filter((e) => new Date(e.start).toDateString() === opts.now.toDateString())
  for (const e of today) for (let f = 1.01; f < 1.08; f += 0.01) splat(g, cam, at(hourAngle(hoursOf(new Date(e.start))), f), '^', SC.needs, 0.9)

  // Lens: the projector core, and a sparse cone that carries the head up.
  for (let r = 0; r < 0.34; r += 0.03) for (let t = 0; t < 360; t += 4) splat(g, cam, add(LENS, V(r * Math.cos(t * D2R), 0.02, r * Math.sin(t * D2R))), r < 0.18 ? '@' : '*', mix(WHITE, SC.working, r * 2.5), 0.9)
  ringOnPlane(g, cam, add(LENS, V(0, 0.02, 0)), 0.6, '-', SC.working, 0.6, 0.02)
  for (let i = 0; i < 700; i++) {
    const y = Math.pow(rnd(), 1.4) * 3.2, rad = (0.2 + y * 0.42) * Math.sqrt(rnd()), t = rnd() * Math.PI * 2
    splat(g, cam, add(LENS, V(rad * Math.cos(t), y + 0.05, rad * Math.sin(t))), rnd() < 0.2 ? ':' : '.', mix(SC.working, pal.accent, rnd() * 0.3), 0.08 + 0.2 * (1 - y / 3.2))
  }

  // Planets.
  const used = { urgent: 0, working: 0, quiet: 0 }
  const ps = m.projects.slice(0, MAXP).map((p) => {
    const k = p.state === 'needs' || p.state === 'failed' ? 'urgent' : p.state === 'working' ? 'working' : 'quiet'
    const [a, f] = SPOTS[k][used[k]++ % SPOTS[k].length]
    const base = at(a + spin, f)
    const r = 0.16 + 0.07 * Math.min(p.open, 5)
    const stalk = p.state === 'needs' ? 0.45 + 1.75 * Math.min((p.wait ?? 0) / 300, 1) : p.state === 'working' ? 0.24 : 0
    return { ...p, base, r, stalk, c: add(base, V(0, stalk + r, 0)) }
  })
  // Bot intent: longest wait first (model order), then a failure, then live work.
  const target = ps.find((p) => p.state === 'needs') ?? ps.find((p) => p.state === 'failed') ?? ps.find((p) => p.state === 'working')
  const spheres: Sphere[] = []
  for (const p of ps) {
    const col = SC[p.state], flat = add(p.base, V(0, 0.02, 0))
    if (p.state === 'done') { ringOnPlane(g, cam, flat, 0.34, 'o', col, 0.5, 0.05); ringOnPlane(g, cam, flat, 0.18, '.', col, 0.35, 0.09); continue }
    if (p.state === 'failed') { brokenRing(g, cam, add(p.base, V(0, 0.08, 0)), col, 0.12); continue }
    ringOnPlane(g, cam, flat, p.r * 1.2, '.', mix(col, [20, 20, 50], 0.35), 0.5, 0.06)
    for (let y = 0.05; y < p.stalk; y += 0.07) splat(g, cam, add(p.base, V(0, y, 0)), ':', col, 0.55)
    const quiet = p.state === 'sleeping' || p.state === 'queued'
    spheres.push({ c: p.c, r: p.r, col: quiet ? pal.ring : col, emit: p === target ? 0.55 : p.state === 'working' ? 0.3 : 0.12, dimf: quiet ? 0.45 : p === target ? 1 : 0.85 })
    if (p === target && p.state === 'needs') octa(g, cam, add(p.c, V(0, p.r + 0.34, 0)), 0.12, col, 1.4)
  }
  raySpheres(g, cam, spheres)

  // Head: option B glyph grid on its own layer at HS scale, above the lens, turned toward the
  // target. Its eyes (the violet cells) take the target's state colour.
  const lens = cam.proj(LENS)
  let head: HeadBox = { x: w / 2, y: lens.y, top: lens.y, left: w / 2, right: w / 2 }
  let layer: HeadLayer | undefined
  const tx = target ? cam.proj(target.c).x : w / 2
  const H = turn && target && Math.abs(tx - w / 2) > 40 ? (tx < w / 2 ? turn.l : turn.r) : heads
  if (H) {
    const hd = H.moods[m.mood], rows = hd.glyphs.length
    const hg: Grid = { cols: hd.cols, rows, z: new Float32Array(hd.cols * rows).fill(1e9), ch: new Array(hd.cols * rows).fill(null), col: new Array(hd.cols * rows), a: new Float32Array(hd.cols * rows) }
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
    const hw = hd.cols * CW * HS, hh = rows * CH * HS
    const top = Math.max(60, (70 + cam.proj(at(270, 1)).y - 50 - hh) / 2)
    layer = { g: hg, scale: HS, x: w / 2 - hw / 2, y: top, zSplit: lens.z }
    head = { x: w / 2, y: top + hh / 2, top, left: w / 2 - hw / 2, right: w / 2 + hw / 2 }
  }

  paint(cv, g, dpr, layer)

  const gaps = new Map<PlanetLabel, number>()
  const out: SceneOut = { planets: [], rim: [], head, lens: { x: lens.x, y: lens.y } }
  for (const p of ps) {
    const flat = p.state === 'done' || p.state === 'failed'
    const s = cam.proj(flat ? p.base : p.c)
    const rad = ((flat ? 0.34 : p.r) * cam.focal) / s.z
    const side = s.x < w / 2 ? 'l' : 'r'
    const note = p.state === 'working' && p.runFrom ? elapsed(p.runFrom, opts.now)
      : p.state === 'failed' ? `${p.cause ?? 'run'} failed`
      : p.state === 'done' || p.state === 'sleeping' || p.state === 'queued' ? p.state === 'sleeping' ? 'sleep' : p.state : undefined
    const label: PlanetLabel = { x: s.x + (side === 'r' ? rad + 8 : -rad - 8), y: s.y, name: p.name, state: p.state, side, note, target: p === target }
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
    if (p.state === 'done' || p.state === 'failed') continue
    const c = cam.proj(p.c), rr = (p.r * cam.focal) / c.z
    boxes.push({ x0: c.x - rr, x1: c.x + rr, y0: c.y - rr, y1: c.y + rr, own: p.name })
  }
  const next = today.filter((e) => hoursOf(new Date(e.start)) >= nowH).sort((a, b) => a.start.localeCompare(b.start))[0]
  let ev: RimLabel | undefined
  if (next) {
    const d = new Date(next.start), s = cam.proj(at(hourAngle(hoursOf(d)), 1.12))
    ev = { x: s.x, y: s.y, text: `${two(d.getHours())}:${two(d.getMinutes())} ${next.title}`, kind: 'event', side: s.x < w / 2 ? 'l' : 'r' }
    // Matches .hero-rim.event: right side ends at the notch, left side starts there; 6px below.
    const tw = ev.text.length * 6.4, x0 = ev.side === 'r' ? ev.x - tw + 8 : ev.x - 8, x1 = x0 + tw
    // Keep it off planet spheres: slide down past any disc it crosses.
    for (const p of ps) {
      if (p.state === 'done' || p.state === 'failed') continue
      const c = cam.proj(p.c), rr = (p.r * cam.focal) / c.z
      if (c.x + rr > x0 && c.x - rr < x1 && c.y + rr > ev.y && c.y - rr < ev.y + 20) ev.y = c.y + rr + 2
    }
    boxes.push(line(x0, x1, ev.y + 12))
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
  // Rim hour labels (hidden where the head covers them) and the next agenda event.
  const inHead = (x: number, y: number) => x > head.left + 20 && x < head.right - 20 && y > head.top && y < 2 * head.y - head.top
  for (const hh of [0, 6, 12, 18]) {
    const s = cam.proj(at(hourAngle(hh), 1.09))
    if (!inHead(s.x, s.y)) out.rim.push({ x: s.x, y: s.y, text: two(hh), kind: 'hour' })
  }
  if (ev) out.rim.push(ev)
  return out
}

interface HeadLayer { g: Grid; scale: number; x: number; y: number; zSplit: number }
function paint(cv: HTMLCanvasElement, g: Grid, dpr: number, head?: HeadLayer) {
  const ctx = cv.getContext('2d')
  if (!ctx) return
  const L = document.createElement('canvas')
  L.width = cv.width; L.height = cv.height
  const l = L.getContext('2d')!
  l.scale(dpr, dpr)
  l.textBaseline = 'middle'; l.textAlign = 'center'
  const cells = (gr: Grid, keep: (z: number) => boolean) => {
    for (let r = 0; r < gr.rows; r++) for (let c = 0; c < gr.cols; c++) {
      const i = r * gr.cols + c, ch = gr.ch[i]
      if (!ch || !keep(gr.z[i])) continue
      if (ch === ' ') { l.clearRect(c * CW - 0.5, r * CH - 0.5, CW + 1, CH + 1); continue }
      const col = gr.col[i]
      l.fillStyle = `rgba(${col[0] | 0},${col[1] | 0},${col[2] | 0},${gr.a[i].toFixed(3)})`
      l.fillText(ch, c * CW + CW / 2, r * CH + CH / 2 + 0.5)
    }
  }
  l.font = "500 6.65px 'Fira Code', ui-monospace, monospace"
  // Behind the lens, then the head (its silhouette hides the table back), then the front.
  cells(g, (z) => !head || z > head.zSplit)
  if (head) {
    l.save(); l.translate(head.x, head.y); l.scale(head.scale, head.scale)
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
