import { describe, expect, it } from 'vitest'
import { fixtureSource } from '../data/fixtures'
import { heroModel, moodOf } from './model'

const zero = { working: 0, needs: 0, failed: 0, sleeping: 0, queued: 0, done: 0 }

describe('hero model', () => {
  it('picks the mood by urgency', () => {
    expect(moodOf(zero, 0)).toBe('idle')
    expect(moodOf({ ...zero, working: 1 }, 0)).toBe('run')
    expect(moodOf({ ...zero, working: 1, failed: 1 }, 0)).toBe('fail')
    expect(moodOf({ ...zero, failed: 1 }, 2)).toBe('ask')
  })

  it('builds one node per project from the fixtures', async () => {
    const src = fixtureSource()
    const missions = await src.missions()
    const m = heroModel(missions, await src.inbox())
    expect(m.projects.length).toBe(new Set(missions.map((x) => x.project)).size)
    expect(m.projects[0].state).toBe('needs')
    expect(m.mood).toBe('ask')
    expect(m.line).toMatch(/Sir\./)
  })
})
