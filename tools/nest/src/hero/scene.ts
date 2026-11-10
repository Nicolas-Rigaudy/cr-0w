import type { UiState } from '../data/state'
import type { HeroModel } from './model'

// Static ASCII scene, ported from the approved hero pass 7 renderer (item_139) with the
// option B head (item_126 shader port) pasted in from heads.json. Draws one frame per call.

export interface Heads { palette: string[]; moods: Record<string, { cols: number; glyphs: string[]; colours: string[] }> }
type RGB = [number, number, number]
export interface Palette { ring: RGB; dim: RGB; fg: RGB; accent: RGB; state: Record<UiState, RGB> }
export interface Label { x: number; y: number; text: string; state: UiState; live?: boolean; right?: boolean; left?: boolean }
export interface SceneOut { projects: Label[]; missions: Label[]; head: { x: number; y: number } }

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
function bez(a: V3, b: V3, lift: number) {
  const m = add(mul(add(a, b), 0.5), V(0, lift, 0))
  return (t: number) => add(add(mul(a, (1 - t) * (1 - t)), mul(m, 2 * t * (1 - t))), mul(b, t * t))
}
function link(g: Grid, cam: Cam, f: (t: number) => V3, col: RGB, a: number) {
  for (let t = 0; t <= 1; t += 0.004) splat(g, cam, f(t), '.', col, a)
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
function brokenRing(g: Grid, cam: Cam, c: V3, col: RGB) {
  const R = 0.26, r = 0.06, tilt = Math.PI / 2 - 0.25
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

// Urgency is closeness: needs/failed planets come to the front (90 = facing the camera),
// working ones sit at the sides, quiet ones shrink to the back behind the head.
const HOT: UiState[] = ['working', 'needs', 'failed']
const URGENT: UiState[] = ['needs', 'failed']
const SLOTS = { urgent: [30, 150, 58, 122], working: [8, 172, 350, 190], quiet: [270, 235, 305, 215, 325] }
const MAX = 8
function angles(states: UiState[]) {
  const used = { urgent: 0, working: 0, quiet: 0 }
  return states.map((st) => {
    const k = URGENT.includes(st) ? 'urgent' : st === 'working' ? 'working' : 'quiet'
    return SLOTS[k][used[k]++ % SLOTS[k].length]
  })
}

export function drawScene(cv: HTMLCanvasElement, w: number, h: number, dpr: number, m: HeroModel, heads: Heads | undefined, pal: Palette, turn?: { l: Heads; r: Heads }): SceneOut {
  cv.width = Math.round(w * dpr); cv.height = Math.round(h * dpr)
  const g = makeGrid(w, h), rnd = rng(42)
  const ORB = 4.0, OUT = 4.7, PY = -0.55
  const focal = ((w / 2 - 70) * 16.4) / OUT // side nodes and their missions stay in view; the outer ring may brush the panels
  const cam = camera(V(0, 10.0, -13.0), V(0, 0.35, 0), focal, w / 2, h * 0.63)
  const onRing = (a: number, R = ORB) => V(R * Math.cos(a * D2R), PY, -R * Math.sin(a * D2R))
  const HEAD = V(0, 1.15, 0)
  const SC = pal.state

  // Background: faint vault stars and their nearest links.
  const stars: [number, number, number][] = []
  for (let i = 0; i < Math.round((w * h) / 4300); i++) stars.push([rnd() * w, rnd() * h, rnd()])
  for (const [x, y, b] of stars) put(g, Math.floor(x / CW), Math.floor(y / CH), 1e8, b > 0.93 ? '+' : b > 0.6 ? '.' : '·', b > 0.93 ? pal.ring : pal.dim, 0.25 + b * 0.35)
  for (const [x, y] of stars) {
    let bd = 1e9, bx = 0, by = 0
    for (const [x2, y2] of stars) { const d = Math.hypot(x2 - x, y2 - y); if (d > 0 && d < bd) { bd = d; bx = x2; by = y2 } }
    if (bd > 140) continue
    for (let t = 0.15; t < 0.85; t += 10 / bd) put(g, Math.floor(lerp(x, bx, t) / CW), Math.floor(lerp(y, by, t) / CH), 1.1e8, '·', pal.dim, 0.14)
  }

  // One orbit ring for the projects (queued missions sit just outside it, no ring of their own).
  for (let t = 0; t < Math.PI * 2; t += 0.0035) {
    const front = -Math.sin(t)
    splat(g, cam, V(ORB * Math.cos(t), PY, ORB * Math.sin(t)), front > 0.5 ? '-' : front > -0.2 ? '~' : '.', pal.ring, 0.26 + 0.3 * Math.max(0, front))
  }

  // Positions.
  const ang = angles(m.projects.slice(0, MAX).map((p) => p.state))
  // Quiet planets drop their missions, urgent ones keep only what needs you; the head looks at the first urgent one.
  const projects = m.projects.slice(0, MAX).map((p, i) => ({ ...p, a: ang[i], at: onRing(ang[i]), missions: !HOT.includes(p.state) ? [] : URGENT.includes(p.state) ? p.missions.filter((x) => URGENT.includes(x.state)) : p.missions }))
  const target = projects.find((p) => URGENT.includes(p.state))
  const mpos = new Map<string, V3>()
  for (const p of projects) {
    const rad = V(Math.cos(p.a * D2R), 0, -Math.sin(p.a * D2R)), tan = V(-Math.sin(p.a * D2R), 0, -Math.cos(p.a * D2R))
    const list = p.missions.filter((x) => x.state !== 'queued')
    list.forEach((x, i) => {
      const k = i - (list.length - 1) / 2
      mpos.set(x.id, add(p.at, add(mul(rad, 0.85 + Math.abs(k) * 0.06), mul(tan, k * 1.3))))
    })
    p.missions.filter((x) => x.state === 'queued').forEach((x, i) => mpos.set(x.id, onRing(p.a - 6 - i * 7, OUT)))
  }
  const lifted = m.asks > 0 ? projects.flatMap((p) => p.missions).find((x) => x.state === 'needs')?.id : undefined

  // Short tethers project → missions (no links to the core).
  for (const p of projects) {
    for (const x of p.missions) {
      const mp = mpos.get(x.id)!, fm = bez(p.at, mp, 0.12)
      link(g, cam, fm, SC[x.state], x.state === 'working' || x.state === 'needs' ? 0.5 : 0.18)
    }
  }

  // Nodes as lit spheres.
  const spheres: Sphere[] = projects.map((p) => {
    const quiet = !HOT.includes(p.state)
    const r = URGENT.includes(p.state) ? 0.44 : quiet ? 0.17 : 0.28
    return { c: p.at, r, col: quiet ? pal.ring : SC[p.state], emit: URGENT.includes(p.state) ? 0.45 : p.state === 'working' ? 0.5 : 0.1, dimf: quiet ? 0.45 : 1 }
  })
  for (const p of projects) for (const x of p.missions) {
    const c = mpos.get(x.id)!
    if (x.state === 'working') spheres.push({ c, r: 0.15, col: SC.working, emit: 0.7, dimf: 1 })
    if (x.state === 'sleeping') spheres.push({ c, r: 0.12, col: SC.sleeping, emit: 0, dimf: 0.5 })
  }
  raySpheres(g, cam, spheres)

  // Dust under the floating head.
  for (let i = 0; i < 40; i++) {
    const r = rnd() * 0.8, t = rnd() * Math.PI * 2, y = PY + 0.1 + rnd() * 0.9
    splat(g, cam, V(r * Math.cos(t), y, r * Math.sin(t)), rnd() < 0.15 ? ':' : '·', mix(pal.accent, SC.working, rnd()), 0.2 + rnd() * 0.3)
  }

  // State marks.
  for (const p of projects) for (const x of p.missions) {
    const mp = mpos.get(x.id)!
    if (x.state === 'working') ringOnPlane(g, cam, mp, 0.34, ':', SC.working, 0.6, 0.03)
    if (x.state === 'sleeping') { splat(g, cam, add(mp, V(0.18, 0.32, 0)), 'z', SC.sleeping, 0.6); splat(g, cam, add(mp, V(0.3, 0.48, 0)), 'z', SC.sleeping, 0.4) }
    if (x.state === 'queued') ringOnPlane(g, cam, mp, 0.13, 'o', SC.queued, 0.7, 0.5)
    if (x.state === 'needs') {
      const up = x.id === lifted ? 0.95 : 0.55
      octa(g, cam, add(mp, V(0, up, 0)), x.id === lifted ? 0.17 : 0.12, SC.needs, x.id === lifted ? 1.4 : 1)
      for (let y = 0.05; y < up - 0.2; y += 0.07) splat(g, cam, add(mp, V(0, y, 0)), ':', SC.needs, 0.4)
      ringOnPlane(g, cam, mp, 0.2, '-', mix(SC.needs, [20, 20, 50], 0.55), 0.6, 0.08)
      if (x.id === lifted) ringOnPlane(g, cam, mp, 0.45, '~', SC.needs, 0.6, 0.012)
    }
    if (x.state === 'failed') brokenRing(g, cam, add(mp, V(0, 0.3, 0)), SC.failed)
  }

  // Head: option B glyph grid, occluding what lies behind it.
  const hp = cam.proj(HEAD)
  const tx = target ? cam.proj(target.at).x : hp.x
  const H = turn && target && Math.abs(tx - hp.x) > 40 ? (tx < hp.x ? turn.l : turn.r) : heads
  if (H) {
    const heads = H
    const hd = heads.moods[m.mood]
    const c0 = Math.round(hp.x / CW - hd.cols / 2)
    const r0 = Math.max(Math.ceil(52 / CH), Math.round(cam.proj(V(0, PY, 0)).y / CH - hd.glyphs.length * 1.02)) // clear the say line
    hd.glyphs.forEach((row, r) => {
      const first = row.search(/\S/)
      for (let c = Math.max(0, first); c < row.length; c++) {
        const k = hd.colours[r].slice(2 * c, 2 * c + 2)
        if (row[c] === ' ' || k.length < 2 || k === '  ') { put(g, c0 + c, r0 + r, hp.z, null, pal.dim, 0); continue }
        const [hex, alpha] = heads.palette[KEYS.indexOf(k[0]) * KEYS.length + KEYS.indexOf(k[1])].split('/')
        const n = parseInt(hex.slice(1), 16)
        put(g, c0 + c, r0 + r, hp.z - 1, row[c], [(n >> 16) & 255, (n >> 8) & 255, n & 255], Number(alpha))
      }
    })
  }

  paint(cv, g, dpr)

  const out: SceneOut = { projects: [], missions: [], head: { x: hp.x, y: hp.y } }
  for (const p of projects) {
    const s = cam.proj(p.at)
    out.projects.push({ x: s.x, y: s.y - 30, text: p.name, state: p.state, live: HOT.includes(p.state) })
    for (const x of p.missions) {
      if (!URGENT.includes(x.state)) continue // working runs stay a lit sphere, no text
      if (!URGENT.includes(p.state)) continue // only the front planets carry mission labels
      const lift = x.state === 'needs' ? (x.id === lifted ? 0.95 : 0.55) : 0
      const q = cam.proj(add(mpos.get(x.id)!, V(0, lift, 0)))
      out.missions.push({ x: q.x, y: q.y, text: x.title, state: x.state, right: q.x > w - 170 || (x.state === 'needs' && q.x < w / 2 && q.x > x.title.length * 7 + 24), left: q.x < 170 }) // raised labels face outward on the left half when they fit
    }
  }
  // Nudge mission labels down until they no longer overlap (approximate text boxes).
  const placed: { cx: number; cy: number; w: number }[] = []
  for (const l of out.missions.sort((a, b) => a.y - b.y)) {
    const tw = l.text.length * 6.6 + 50, up = l.state === 'needs', side = l.right ? -1 : l.left || up ? 1 : 0
    const cx = l.x + side * (tw / 2 + (up ? 14 : 0))
    let cy = up ? l.y : l.y + 20
    while (placed.some((b) => Math.abs(b.cy - cy) < 18 && Math.abs(b.cx - cx) < (b.w + tw) / 2)) { cy += 18; l.y += 18 }
    placed.push({ cx, cy, w: tw })
  }
  return out
}

function paint(cv: HTMLCanvasElement, g: Grid, dpr: number) {
  const ctx = cv.getContext('2d')
  if (!ctx) return
  const L = document.createElement('canvas')
  L.width = cv.width; L.height = cv.height
  const l = L.getContext('2d')!
  l.scale(dpr, dpr)
  l.font = "500 6.65px 'Fira Code', ui-monospace, monospace"
  l.textBaseline = 'middle'; l.textAlign = 'center'
  for (let r = 0; r < g.rows; r++) for (let c = 0; c < g.cols; c++) {
    const i = r * g.cols + c, ch = g.ch[i]
    if (!ch) continue
    const col = g.col[i]
    l.fillStyle = `rgba(${col[0] | 0},${col[1] | 0},${col[2] | 0},${g.a[i].toFixed(3)})`
    l.fillText(ch, c * CW + CW / 2, r * CH + CH / 2 + 0.5)
  }
  ctx.clearRect(0, 0, cv.width, cv.height)
  ctx.save()
  ctx.globalCompositeOperation = 'lighter'
  ctx.filter = `blur(${14 * dpr}px)`; ctx.globalAlpha = 0.55; ctx.drawImage(L, 0, 0)
  ctx.filter = `blur(${3 * dpr}px)`; ctx.globalAlpha = 0.6; ctx.drawImage(L, 0, 0)
  ctx.restore()
  ctx.drawImage(L, 0, 0)
}
