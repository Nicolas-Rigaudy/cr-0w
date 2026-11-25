import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import type { LiveReply } from '../data/live'
import { cancelMission, refused, stopMissionRun } from './missionSheets'
import { Sheet } from './Sheet'

const snapshot = async (src: ReturnType<typeof fixtureSource>): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(src, (x) => { box.s = x })()
  return box.s!
}
const setup = async () => {
  const src = fixtureSource()
  const s = await snapshot(src)
  const m = s.missions.find((x) => x.state === 'running')!
  const state = async () => (await snapshot(src)).missions.find((x) => x.id === m.id)!.state
  return { src, m, state }
}

describe('Stop on a mission run (item_279)', () => {
  it('cancels the mission with one confirm', async () => {
    const { src, m, state } = await setup()
    const done = vi.fn()
    const c = stopMissionRun({ mission: m, run: 'r1', name: 'S1', act: src.actions, refresh: () => {}, done })
    expect(c.ok).toBe('Stop and cancel the mission')
    expect(c.danger).toBe(true)
    await c.run()
    expect(await state()).toBe('cancelled')
    expect(done).toHaveBeenCalledWith('cancelled')
  })

  it('Restart run stops the run only; the mission stays alive', async () => {
    const { src, m, state } = await setup()
    const done = vi.fn()
    const c = stopMissionRun({ mission: m, run: 'r1', name: 'S1', act: src.actions, refresh: () => {}, done })
    expect(c.alt?.label).toBe('Restart run')
    await c.alt!.run()
    expect(await state()).not.toBe('cancelled')
    expect(done).toHaveBeenCalledWith('restarted')
  })

  it('a refused cancel throws its message and refetches on 409', async () => {
    const { src, m } = await setup()
    await src.actions.missionCancel(m.id)
    const refresh = vi.fn()
    const done = vi.fn()
    const c = stopMissionRun({ mission: m, run: 'r1', name: 'S1', act: src.actions, refresh, done })
    await expect(c.run()).rejects.toThrow(/cannot be cancelled/)
    expect(refresh).toHaveBeenCalled()
    expect(done).not.toHaveBeenCalled()
  })

  it('cancel is offered on a running mission and stops its run', async () => {
    const { src, m, state } = await setup()
    const c = cancelMission({ mission: m, act: src.actions, refresh: () => {}, done: () => {} })
    expect(c.lead).toContain('stops its run')
    await c.run()
    expect(await state()).toBe('cancelled')
  })
})

describe('failed actions show in place (item_279)', () => {
  it('a refused reply throws its message', () => {
    const r = { ok: false, message: 'run r1 is gone', action: 'stop', target: 'r1' } as LiveReply
    expect(() => refused(r)).toThrow('run r1 is gone')
    const ok = { ...r, ok: true }
    expect(refused(ok)).toBe(ok)
  })

  it('the sheet shows Restart run before the main action', async () => {
    const { src, m } = await setup()
    const c = stopMissionRun({ mission: m, run: 'r1', name: 'S1', act: src.actions, refresh: () => {}, done: () => {} })
    const html = renderToStaticMarkup(<Sheet c={c} onClose={() => {}} />)
    const alt = html.indexOf('>Restart run</button>')
    const main = html.indexOf('>Stop and cancel the mission</button>')
    expect(alt).toBeGreaterThan(-1)
    expect(main).toBeGreaterThan(alt)
  })
})
