import { describe, expect, it } from 'vitest'
import { detent, FRONT, liftCurve, spinToFront, step, stepTo, wrap, type Spin } from './motion'

const ps = [{ name: 'a', a0: 138 }, { name: 'b', a0: 248 }, { name: 'c', a0: 64 }]

describe('hero motion', () => {
  it('wraps angles to -180..180', () => {
    expect(wrap(190)).toBe(-170)
    expect(wrap(-190)).toBe(170)
    expect(wrap(360)).toBe(0)
  })

  it('brings a planet to the front by the shortest turn', () => {
    expect(wrap(138 + spinToFront(138, 0) - FRONT)).toBe(0)
    expect(spinToFront(138, 0)).toBe(-48)
    expect(spinToFront(138, 350)).toBe(312)
  })

  it('snaps to the planet nearest the front', () => {
    expect(detent(ps.map((p) => p.a0), 0)).toBe(26)
    expect(detent([], 12)).toBe(12)
  })

  it('slows down, then settles exactly on a detent', () => {
    const s: Spin = { spin: 0, vel: 0.5 }
    let n = 0
    while (step(s, 16, ps.map((p) => p.a0)) && n < 2000) n++
    expect(n).toBeLessThan(2000)
    expect(s.vel).toBe(0)
    expect(ps.some((p) => Math.abs(wrap(p.a0 + s.spin - FRONT)) < 1e-6)).toBe(true)
  })

  it('steps through planets left to right and wraps', () => {
    expect(stepTo(ps, 0, undefined, 1)).toBe('c')
    const seen = new Set<string>()
    let cur: string | undefined = 'a'
    for (let i = 0; i < 3; i++) { cur = stepTo(ps, 0, cur, 1); seen.add(cur!) }
    expect(seen.size).toBe(3)
    expect(stepTo(ps, 0, stepTo(ps, 0, 'a', 1), -1)).toBe('a')
  })

  it('lifts with one tick and lands on 1', () => {
    expect(liftCurve(0)).toBe(0)
    expect(liftCurve(0.6)).toBeCloseTo(1.25)
    expect(liftCurve(1)).toBe(1)
  })
})
