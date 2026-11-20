import { describe, expect, it } from 'vitest'
import { apiButler, fixtureButler } from './butler'

describe('butler clear (item_259 follow-up)', () => {
  it('posts to /clear with the CSRF header and returns the view', async () => {
    const calls: [string, RequestInit | undefined][] = []
    const api = apiButler('/api', async (i, init) => {
      calls.push([i, init])
      return new Response(JSON.stringify({ conv: 'main', messages: [], busy: false, last: 9, cleared: 9 }), { status: 200 })
    })
    const v = await api.clear('main')
    expect(calls[0][0]).toBe('/api/butler/main/clear')
    expect(calls[0][1]?.method).toBe('POST')
    expect((calls[0][1]?.headers as Record<string, string>)['X-Crow']).toBe('1')
    expect(v.cleared).toBe(9)
  })

  it('fixture: clear empties the view and seqs keep growing', async () => {
    const now = new Date('2026-11-20T10:00:00Z')
    const api = fixtureButler(now, [{ seq: 1, ts: now.toISOString(), role: 'you', text: 'hi' }])
    const v = await api.clear('main')
    expect(v.messages).toEqual([])
    expect(v.cleared).toBe(2)
    const m = await api.say('main', 'again', 'nest')
    expect(m.seq).toBe(3)
    expect((await api.get('main', 0)).messages.map((x) => x.seq)).toEqual([3])
  })
})
