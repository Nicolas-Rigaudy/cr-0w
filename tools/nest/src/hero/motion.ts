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
