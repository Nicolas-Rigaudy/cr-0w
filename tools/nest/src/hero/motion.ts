// Spin physics for the holotable (item_158). Angles in degrees; 90 faces the camera.

export const FRONT = 90
// Drag: degrees per pixel of horizontal travel, as a share of the hero width.
export const DEG_PER_WIDTH = 220
// Inertia half-life feel: velocity decays by exp(-dt / FRICTION_MS).
export const FRICTION_MS = 380
// Below this speed (deg/ms) inertia hands over to the detent.
export const SETTLE_SPEED = 0.025
// Detent ease time constant.
export const DETENT_MS = 140

export const wrap = (a: number) => ((((a + 180) % 360) + 360) % 360) - 180

// Spin that brings the table angle a0 to the front, nearest to the current spin.
export function spinToFront(a0: number, spin: number): number {
  return spin + wrap(FRONT - (a0 + spin))
}

// Spin of the detent: the planet nearest to the front comes to the front.
export function detent(a0s: number[], spin: number): number {
  if (!a0s.length) return spin
  let best = spin, d = Infinity
  for (const a of a0s) {
    const t = spinToFront(a, spin), dd = Math.abs(t - spin)
    if (dd < d) { d = dd; best = t }
  }
  return best
}

export interface Spin { spin: number; vel: number; to?: number }

// One physics step. Free spin with friction, then ease to the detent target. Returns true while moving.
export function step(s: Spin, dt: number, a0s: number[]): boolean {
  if (s.to === undefined) {
    s.spin += s.vel * dt
    s.vel *= Math.exp(-dt / FRICTION_MS)
    if (Math.abs(s.vel) > SETTLE_SPEED) return true
    s.vel = 0
    s.to = detent(a0s, s.spin)
  }
  const d = s.to - s.spin
  if (Math.abs(d) < 0.05) { s.spin = s.to; s.to = undefined; return false }
  s.spin += d * (1 - Math.exp(-dt / DETENT_MS))
  return true
}

// Planet names ordered left to right as seen from the camera, for Left/Right stepping.
export function byScreenX<T extends { name: string; a0: number }>(ps: T[], spin: number): string[] {
  return [...ps].sort((a, b) => Math.cos((b.a0 + spin) * Math.PI / 180) - Math.cos((a.a0 + spin) * Math.PI / 180)).map((p) => p.name)
}

// Next planet when stepping: from the focused one, or from the one at the front.
export function stepTo<T extends { name: string; a0: number }>(ps: T[], spin: number, current: string | undefined, dir: 1 | -1): string | undefined {
  if (!ps.length) return undefined
  const order = byScreenX(ps, spin)
  const front = ps.reduce((a, b) => (Math.abs(wrap(FRONT - (b.a0 + spin))) < Math.abs(wrap(FRONT - (a.a0 + spin))) ? b : a)).name
  const i = order.indexOf(current ?? front)
  if (current === undefined) return front
  return order[(i + dir + order.length) % order.length]
}

// needs-you lift: rise to 1 with one small overshoot tick (t in 0..1).
export function liftCurve(t: number): number {
  if (t >= 1) return 1
  if (t < 0.6) { const u = t / 0.6; return 1.25 * (1 - (1 - u) ** 3) }
  const u = (t - 0.6) / 0.4
  return 1.25 - 0.25 * u
}

// Critically damped spring (item_196), tuned by half-life: the time to close half the gap.
// Fastest settle with no overshoot from rest; peak speed is bounded by gap * 2 ln2 / (e * half-life).
export interface Spring { x: number; v: number }
export function spring(s: Spring, goal: number, halfMs: number, dt: number): void {
  const y = (2 * Math.LN2) / Math.max(halfMs, 1e-3)
  const j0 = s.x - goal, j1 = s.v + j0 * y, e = Math.exp(-y * dt)
  s.x = e * (j0 + j1 * dt) + goal
  s.v = e * (s.v - j1 * y * dt)
}
export const settled = (s: Spring, goal: number, eps = 0.05) => Math.abs(s.x - goal) < eps && Math.abs(s.v) < eps / 50

// Gaze (item_196): the eyes jump to the target first, the head follows on a slower spring.
// Inside the dead zone only the eyes move. Once the target leaves it the head tracks the target
// all the way, and goes back to eyes-only when it has arrived (hysteresis: no twitch at the edge).
export const GAZE = { eyeMs: 40, headMs: 220, dead: 6, eyeMax: 12 }
export interface Gaze { head: Spring; eye: Spring; goal: number; track?: boolean }
export const gaze = (yaw = 0): Gaze => ({ head: { x: yaw, v: 0 }, eye: { x: 0, v: 0 }, goal: yaw })
// Returns true while still moving.
export function stepGaze(g: Gaze, target: number, dt: number, o = GAZE): boolean {
  if (Math.abs(target - g.goal) > o.dead) g.track = true
  if (g.track) g.goal = target
  spring(g.head, g.goal, o.headMs, dt)
  if (g.track && settled(g.head, g.goal, 0.5)) g.track = false
  const off = Math.max(-o.eyeMax, Math.min(o.eyeMax, target - g.head.x))
  spring(g.eye, off, o.eyeMs, dt)
  return !settled(g.head, g.goal) || !settled(g.eye, off, 0.2)
}
