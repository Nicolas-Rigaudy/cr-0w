import { describe, expect, it } from 'vitest'
import { bakeLook, CHARACTER, faceOf, HeadCtl, lookAt, YAW_MAX } from './headCtl'

const tip: [number, number, number] = [46, 242, 255]
const run = (ctl: HeadCtl, from: number, to: number, look: (t: number) => { yaw: number; pitch: number } | undefined, face: 'idle' | 'run' | 'fail' = 'idle') => {
  const out = []
  for (let t = from; t <= to; t += 16) out.push({ t, ...ctl.step({ t, dt: 16, look: look(t), glances: [], face, tip }) })
  return out
}

describe('head character (items 196 + 198)', () => {
  it('a fast pointer sweep gives one smooth turn: monotonic, no overshoot, bounded speed', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    run(ctl, 0, 1000, () => ({ yaw: -YAW_MAX, pitch: 0 }))
    // Sweep across the hero in 120 ms, then hold.
    const f = run(ctl, 1016, 3000, (t) => ({ yaw: Math.min(YAW_MAX, -YAW_MAX + ((t - 1016) / 120) * 2 * YAW_MAX), pitch: 0 }))
    let maxV = 0
    for (let i = 1; i < f.length; i++) {
      const d = f[i].frame.yaw - f[i - 1].frame.yaw
      expect(d).toBeGreaterThanOrEqual(-1e-9)
      maxV = Math.max(maxV, d / 16)
    }
    expect(Math.max(...f.map((x) => x.frame.yaw))).toBeLessThanOrEqual(YAW_MAX + 1e-6)
    // Peak speed under 0.2 deg/ms (56 deg over a 220 ms half-life spring).
    expect(maxV).toBeLessThan(0.2)
    expect(f.at(-1)!.frame.yaw).toBeCloseTo(YAW_MAX, 1)
  })

  it('settles and stops asking for frames when nothing moves', () => {
    const ctl = new HeadCtl({ ...CHARACTER.calm, glance: null, scan: null, flicker: null, blink: [1e9, 1e9] }, 7, false)
    const f = run(ctl, 0, 4000, () => ({ yaw: 10, pitch: 2 }))
    expect(f.at(-1)!.busy).toBe(false)
  })

  it('blinks within the character range and the blink closes the eyes', () => {
    const ctl = new HeadCtl(CHARACTER.attentive)
    const f = run(ctl, 0, 8000, () => undefined)
    const at = f.find((x) => x.frame.blink > 0.9)
    expect(at).toBeDefined()
    expect(at!.t).toBeGreaterThan(CHARACTER.attentive.blink[0] - 1)
  })

  it('fail glitches first, then braces', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    run(ctl, 0, 600, () => undefined, 'run')
    const f = run(ctl, 616, 1600, () => undefined, 'fail')
    expect(f[0].frame.glitch).toBeGreaterThan(0.9)
    expect(f[0].frame.pose.to).toBe('grab')
    expect(f.at(-1)!.frame.glitch).toBe(0)
    expect(f.at(-1)!.frame.pose).toMatchObject({ to: 'brace', k: 1 })
  })

  it('nods on present and comes back level', () => {
    const ctl = new HeadCtl({ ...CHARACTER.attentive, breathe: 0 }, 7, false)
    run(ctl, 0, 2000, () => ({ yaw: 0, pitch: 0 }))
    ctl.present(2016)
    const f = run(ctl, 2016, 4500, () => ({ yaw: 0, pitch: 0 }))
    expect(Math.min(...f.map((x) => x.frame.pitch))).toBeLessThan(-2)
    expect(f.some((x) => x.frame.pose.to === 'present')).toBe(true)
    expect(Math.abs(f.at(-1)!.frame.pitch)).toBeLessThan(0.1)
    expect(f.at(-1)!.frame.pose.to).toBe('rest')
  })

  it('maps the screen to clamped head angles', () => {
    expect(lookAt(500, 300, 1000, 600, 300)).toEqual({ yaw: 0, pitch: 0 })
    expect(lookAt(5000, 9000, 1000, 600, 300).yaw).toBe(YAW_MAX)
  })

  it('bake fallback turns past 10 and comes back under 6, holding in between', () => {
    expect(bakeLook(12, 0)).toBe(1)
    expect(bakeLook(8, 1)).toBe(1)
    expect(bakeLook(8, 0)).toBe(0)
    expect(bakeLook(5, 1)).toBe(0)
    expect(bakeLook(-11, 0)).toBe(-1)
  })

  it('the face follows the work, never the ask face', () => {
    expect(faceOf('ask', true)).toBe('run')
    expect(faceOf('ask', false)).toBe('idle')
    expect(faceOf('fail', true)).toBe('fail')
    expect(faceOf('run', false)).toBe('run')
  })

  const step = (ctl: HeadCtl, from: number, to: number, x: Partial<Parameters<HeadCtl['step']>[0]> = {}) => {
    const out = []
    for (let t = from; t <= to; t += 16) out.push(ctl.step({ t, dt: 16, glances: [], face: 'idle', tip, ...x }))
    return out
  }

  it('waves once on load with happy eyes, then rests', () => {
    const ctl = new HeadCtl(CHARACTER.attentive)
    const f = step(ctl, 0, 4000)
    const mid = f[Math.round(1200 / 16)].frame
    expect(mid.pose.to).toBe('wave')
    expect(mid.eyes).toBe(2)
    expect(f.at(-1)!.frame.pose).toMatchObject({ to: 'rest', k: 1 })
    expect(f.at(-1)!.frame.eyes).toBeUndefined()
  })

  it('thinks only after a long hover, and stops when the hover ends', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    step(ctl, 0, 500)
    expect(step(ctl, 516, 1500, { hover: true }).at(-1)!.frame.pose.to).toBe('rest')
    expect(step(ctl, 1516, 2400, { hover: true }).at(-1)!.frame.pose.to).toBe('think')
    expect(step(ctl, 2416, 3500).at(-1)!.frame.pose).toMatchObject({ to: 'rest', k: 1 })
  })

  it('presents with the arm on the side of the target', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    step(ctl, 0, 500, { look: { yaw: -20, pitch: 0 } })
    ctl.present(516)
    expect(step(ctl, 516, 900, { look: { yaw: -20, pitch: 0 } }).at(-1)!.frame.pose).toMatchObject({ to: 'present', side: [1, -1] })
  })

  it('taps the table while working: claws alternate and sparks only on contact', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    const f = step(ctl, 0, 3000, { face: 'run' }).slice(100)
    const sides = new Set(f.filter((x) => x.frame.spark).map((x) => x.frame.spark!.side))
    expect(sides).toEqual(new Set([-1, 1]))
    expect(f.some((x) => !x.frame.spark)).toBe(true)
    expect(f.every((x) => x.busy)).toBe(true)
  })

  it('drag: eyes widen on the grab, arms trail the spin (bounded), then settle', () => {
    const ctl = new HeadCtl({ ...CHARACTER.calm, glance: null, scan: null, flicker: null, blink: [1e9, 1e9] }, 7, false)
    step(ctl, 0, 1000)
    const d = step(ctl, 1016, 1600, { drag: 0.3 })
    expect(Math.max(...d.map((x) => x.frame.eyeS!))).toBeGreaterThan(1.15)
    const dx = d.at(-1)!.frame.arms!.dw[1][0]
    expect(dx).toBeLessThan(-0.2)
    expect(dx).toBeGreaterThanOrEqual(-0.28 - 0.02)
    const after = step(ctl, 1616, 4000)
    expect(after.at(-1)!.busy).toBe(false)
    expect(Math.abs(after.at(-1)!.frame.arms!.dw[1][0])).toBeLessThan(0.01)
  })

  it('fail shivers, braces, then droops; the face swaps under a blink', () => {
    const ctl = new HeadCtl(CHARACTER.attentive, 7, false)
    step(ctl, 0, 600)
    const f = step(ctl, 616, 3000, { face: 'fail' })
    expect(f[0].frame.face).toBe('idle')
    expect(Math.max(...f.slice(0, 8).map((x) => x.frame.blink))).toBeGreaterThan(0.8)
    expect(f[10].frame.face).toBe('fail')
    expect(f.at(-1)!.frame.pose).toMatchObject({ to: 'droop', k: 1 })
    expect(f.at(-1)!.frame.mouth).toBeCloseTo(-0.6, 1)
  })
})
