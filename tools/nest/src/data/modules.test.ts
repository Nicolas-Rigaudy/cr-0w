import { describe, expect, it, vi } from 'vitest'
import api from 'virtual:core-api'
import { apiActions, apiSource } from './api'
import { feed, type Snapshot } from './feed'
import { DEMO_BRICKS, fixtureSource } from './fixtures'
import { hiddenPages, listRows, listSections, offSections, type Modules } from './modules'

const contract = JSON.parse(api['modules.json']) as Modules

describe('bricks from /api/modules (item_237)', () => {
  it('reads the contract fixture: mail off, its today section off, MRs page kept', () => {
    expect([...offSections(contract)]).toEqual(['today'])
    expect(hiddenPages(contract).size).toBe(0)
    expect(listSections(contract).map((s) => `${s.brick}/${s.id}@${s.page}`)).toEqual(['notes/notes@desk'])
  })

  it('hides a brick page only when every section on it is off, never a core page', () => {
    const m = DEMO_BRICKS.modules
    expect([...hiddenPages(m)]).toEqual(['mrs'])
    expect(listSections(m, 'desk')).toHaveLength(1)
    expect(listSections(m, 'home')).toHaveLength(0)
  })

  it('takes rows from an array or {items}, skips rows without an id, drops non-http links', () => {
    const body = (data: unknown) => ({ data, stale: false, updated_at: null })
    expect(listRows(body([{ id: 1, name: 'a' }, { title: 'no id' }]))).toEqual([{ id: '1', title: 'a', detail: undefined, meta: undefined, url: undefined }])
    expect(listRows(body({ items: [{ id: 'x', title: 'X', url: 'javascript:alert(1)' }] }))[0].url).toBeUndefined()
    expect(listRows(body({ items: [{ id: 'x', url: 'https://e.x' }] }))[0]).toMatchObject({ title: 'x', url: 'https://e.x' })
    expect(listRows(undefined)).toEqual([])
  })

  it('never fetches off sections and reads the bricks of on modules', async () => {
    const src = fixtureSource(undefined, undefined, undefined, undefined, undefined, DEMO_BRICKS)
    const section = vi.spyOn(src, 'section')
    let snap: Snapshot | null = null
    await feed(src, (s) => { snap = s })()
    const s = snap! as Snapshot
    expect(section.mock.calls.map((c) => c[0])).not.toContain('today')
    expect(section.mock.calls.map((c) => c[0])).not.toContain('mrs')
    expect(s.sections.today.off && s.sections.mrs.off && s.sections.checks.off).toBe(true)
    expect(s.agenda).toEqual([])
    expect(listRows(s.bricks?.notes)).toHaveLength(3)
    expect(s.down).toBe(false)
  })

  it('flags a new build once core serves another Nest build, and keeps the data live (item_283)', async () => {
    const src = fixtureSource(undefined, undefined, undefined, undefined, undefined, DEMO_BRICKS)
    let build: string | undefined
    src.modules = async () => ({ ...DEMO_BRICKS.modules, nest_build: build })
    const snaps: Snapshot[] = []
    const refresh = feed(src, (s) => snaps.push(s))
    await refresh() // a core without a served dist: nothing to compare
    build = 'aaa'
    await refresh()
    await refresh()
    expect(snaps.some((s) => s.newBuild)).toBe(false)
    build = 'bbb'
    await refresh()
    expect(snaps.at(-1)?.newBuild).toBe(true)
    expect(snaps.at(-1)?.modules?.nest_build).toBe('bbb')
    build = 'aaa' // the flag stays: this page still runs the first build
    await refresh()
    expect(snaps.at(-1)?.newBuild).toBe(true)
  })

  it('posts a module action with the row id and keeps the agenda when today fails', async () => {
    const calls: Array<[string, RequestInit | undefined]> = []
    const ok = (body: unknown) => new Response(JSON.stringify(body), { status: 200 })
    const acts = apiActions('/api', async (u, init) => { calls.push([u, init]); return ok({ message: 'closed' }) })
    await acts.module('notes', 'done', 'n-1')
    expect(calls[0][0]).toBe('/api/modules/notes/actions/done')
    expect(calls[0][1]?.body).toBe('{"item":"n-1"}')
    expect((calls[0][1]?.headers as Record<string, string>)['X-Crow']).toBe('1')
    const src = apiSource('/api', async (u) => (u.endsWith('/state/today')
      ? new Response('{"error":"module off"}', { status: 404 })
      : ok({ section: 'x', version: '1', updated_at: null, data: [] })))
    await expect(src.agenda()).resolves.toEqual([])
  })
})
