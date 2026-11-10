import { describe, expect, it } from 'vitest'
import { fixtureSource } from './fixtures'
import { countByState, uiState } from './state'

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
