import { describe, expect, it } from 'vitest'
import { fixtureSource } from '../data/fixtures'
import { heroModel, moodOf } from './model'

const zero = { working: 0, needs: 0, failed: 0, sleeping: 0, queued: 0, done: 0 }

describe('hero model', () => {
  it('orders the mood fail > ask > run > idle', () => {
    expect(moodOf(zero)).toBe('idle')
    expect(moodOf({ ...zero, working: 1 })).toBe('run')
    expect(moodOf({ ...zero, working: 1 }, 1)).toBe('ask')
    expect(moodOf({ ...zero, working: 1, failed: 1 }, 1)).toBe('fail')
    expect(moodOf({ ...zero, failed: 1 })).toBe('fail')
    expect(moodOf({ ...zero, needs: 2 })).toBe('idle')
  })

  it('builds one node per project from the fixtures', async () => {
    const src = fixtureSource()
    const missions = await src.missions()
    const m = heroModel(missions, await src.inbox())
    expect(m.projects.length).toBe(new Set(missions.map((x) => x.project)).size)
    expect(m.projects[0].state).toBe('needs')
    expect(m.mood).toBe('fail') // the fixtures hold a failed run
    expect(m.line).toMatch(/Sir\./)
  })

  it('maps the holotable values from the fixtures', async () => {
    const src = fixtureSource()
    const m = heroModel(await src.missions(), await src.inbox(), src.now())
    const p = Object.fromEntries(m.projects.map((x) => [x.name, x]))
    expect(m.projects[0].name).toBe('demo') // longest wait first
    expect(p.demo.wait).toBe(290)
    expect(p['cr-0w'].wait).toBe(83)
    expect(p.nest).toMatchObject({ state: 'working', open: 2, runFrom: '2026-11-09T15:05:00+01:00' })
    expect(p.vault).toMatchObject({ state: 'done', open: 0 })
    expect(p.infra).toMatchObject({ state: 'failed', cause: 'Push token rotation' })
  })
})
