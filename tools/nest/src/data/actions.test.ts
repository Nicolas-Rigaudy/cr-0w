import { describe, expect, it } from 'vitest'
import { apiActions, ApiError } from './api'

// A fake core that records each POST and answers with the given status and body.
function core(status = 200, reply: unknown = { ok: true }) {
  const sent: Array<{ url: string; body: unknown; csrf: string | null }> = []
  const get = async (url: string, init?: RequestInit) => {
    sent.push({ url, body: init?.body ? JSON.parse(init.body as string) : undefined, csrf: new Headers(init?.headers).get('X-Crow') })
    return new Response(JSON.stringify(reply), { status, headers: { 'Content-Type': 'application/json' } })
  }
  return { get, sent }
}

describe('mission actions', () => {
  it('posts answer, stop and cancel on the mission routes', async () => {
    const { get, sent } = core()
    const act = apiActions('/api', get)
    await act.missionAnswer('20261109-a b', '20261109-a b/v2/q1', 'approve', '  ', 2)
    await act.missionAnswer('m1', 'm1/v1', 'answer', 'use plan B', 1)
    await act.missionStop('m1')
    await act.missionCancel('m1')
    await act.missionCancel('m1', 'not needed')
    expect(sent).toEqual([
      { url: '/api/missions/20261109-a%20b/answer', body: { item: '20261109-a b/v2/q1', verb: 'approve', version: 2 }, csrf: '1' },
      { url: '/api/missions/m1/answer', body: { item: 'm1/v1', verb: 'answer', text: 'use plan B', version: 1 }, csrf: '1' },
      { url: '/api/missions/m1/stop', body: undefined, csrf: '1' },
      { url: '/api/missions/m1/cancel', body: undefined, csrf: '1' },
      { url: '/api/missions/m1/cancel', body: { reason: 'not needed' }, csrf: '1' },
    ])
  })

  it('keeps the slashes of an inbox row id when hiding it', async () => {
    const { get, sent } = core()
    await apiActions('/api', get).inboxHide('20261109-x y/v3/pick?')
    expect(sent[0].url).toBe('/api/inbox/20261109-x%20y/v3/pick%3F/hide')
  })

  it('throws a 409 with the core message on a stale version', async () => {
    const { get } = core(409, { error: 'version 1 is older than 2' })
    const err = await apiActions('/api', get).missionAnswer('m1', 'm1/v1', 'answer', '', 1).catch((e: unknown) => e)
    expect(err).toBeInstanceOf(ApiError)
    expect((err as ApiError).status).toBe(409)
    expect((err as ApiError).message).toBe('version 1 is older than 2')
  })
})
