import { describe, expect, it } from 'vitest'
import { fixtureSource } from '../data/fixtures'
import { heroModel, moodOf } from './model'

const zero = { working: 0, needs: 0, failed: 0, sleeping: 0, queued: 0, done: 0 }

describe('hero model', () => {
  it('picks the face from the work, never from asks', () => {
    expect(moodOf(zero)).toBe('idle')
    expect(moodOf({ ...zero, working: 1 })).toBe('run')
    expect(moodOf({ ...zero, working: 1, failed: 1 })).toBe('run')
    expect(moodOf({ ...zero, failed: 1 })).toBe('fail')
    expect(moodOf({ ...zero, needs: 2 })).toBe('idle')
  })

  it('builds one node per project from the fixtures', async () => {
    const src = fixtureSource()
    const missions = await src.missions()
    const m = heroModel(missions, await src.inbox())
    expect(m.projects.length).toBe(new Set(missions.map((x) => x.project)).size)
    expect(m.projects[0].state).toBe('needs')
    expect(m.mood).toBe('run')
    expect(m.line).toMatch(/Sir\./)
  })
})
