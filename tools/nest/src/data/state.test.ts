import { describe, expect, it } from 'vitest'
import { fixtureSource } from './fixtures'
import type { Mission } from './types'
import { countByState, missionOfRun, uiState } from './state'

describe('fixture adapter', () => {
  it('loads core and Nest-local records', async () => {
    const src = fixtureSource()
    const missions = await src.missions()
    expect(missions.some((m) => m.id === '20261109-build-demo-login')).toBe(true)
    expect(missions.length).toBeGreaterThanOrEqual(9)
    expect(await src.mission('20261109-fix-push-token-refresh')).toBeDefined()
  })

  it('covers every UI state', async () => {
    const src = fixtureSource()
    const counts = countByState(await src.missions(), await src.inbox())
    for (const n of Object.values(counts)) expect(n).toBeGreaterThan(0)
  })
})

describe('uiState', () => {
  it('maps contract states', async () => {
    const src = fixtureSource()
    const inbox = await src.inbox()
    const of = async (id: string) => uiState((await src.mission(id))!, inbox)
    expect(await of('20261109-build-demo-login')).toBe('working')
    expect(await of('20261109-build-sse-endpoint')).toBe('working')
    expect(await of('20261109-fix-push-token-refresh')).toBe('failed')
    expect(await of('20261109-build-report-export')).toBe('sleeping')
    expect(await of('20261109-research-vault-search')).toBe('queued')
    expect(await of('20261109-build-inbox-filters')).toBe('needs')
    expect(await of('20261109-build-mr-pipeline-badge')).toBe('needs')
    expect(await of('20261109-research-android-push-options')).toBe('done')
    expect(await of('20261109-chore-prune-worktrees')).toBe('done')
  })
})

describe('missionOfRun', () => {
  const run = (id: string, ended: string | null) => ({ run: id, mission: 'm', account: 'a', version: 1, cold: false, started: '2026-11-11T10:00:00Z', ended, cause: null, exit: null })
  const m = (id: string, state: Mission['state'], runs: ReturnType<typeof run>[]) => ({ id, state, runs }) as unknown as Mission

  it('finds the running mission of a live run only', () => {
    const ms = [m('done1', 'done', [run('r1', null)]), m('run1', 'running', [run('r0', '2026-11-11T10:05:00Z'), run('r2', null)])]
    expect(missionOfRun(ms, 'r2')?.id).toBe('run1')
    expect(missionOfRun(ms, 'r0')).toBeUndefined()
    expect(missionOfRun(ms, 'r1')).toBeUndefined()
    expect(missionOfRun(ms, 'adhoc')).toBeUndefined()
  })
})
