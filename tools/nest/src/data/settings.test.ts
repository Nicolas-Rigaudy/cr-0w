import { describe, expect, it } from 'vitest'
import { apiSettings, ApiError } from './api'
import { fixtureSettings } from './fixtures'
import {
  changedKeys, draftOf, errorLines, movePage, needsBuild, pageMark, pendingCount, saveBody, togglePage,
  type NestKnown, type NestSettings,
} from './settings'
import { settingsControl, type SettingsView } from './useSettings'

const known: NestKnown = { themes: ['modern', 'cyber'], heroes: ['holotable', 'none'], pages: ['home', 'desk', 'ops', 'inbox'] }
const nest = (p: Partial<NestSettings> = {}): NestSettings => ({ file: 'crow.yaml', theme: null, hero: null, pages: null, known, errors: [], ...p })
const keys = (d: { pages: { key: string; on: boolean }[] }) => d.pages.map((p) => `${p.key}${p.on ? '' : '-'}`).join(' ')

describe('draftOf', () => {
  it('falls back to the engine defaults and shows every page', () => {
    const d = draftOf(nest())
    expect(d.theme).toBe('modern')
    expect(d.hero).toBe('holotable')
    expect(keys(d)).toBe('home desk ops inbox')
  })

  it('keeps the listed order, puts Home first and the rest off', () => {
    const d = draftOf(nest({ theme: 'cyber', pages: ['ops', 'desk', 'gone', 'ops'] }))
    expect(d.theme).toBe('cyber')
    expect(keys(d)).toBe('home ops desk inbox-')
  })
})

describe('draft edits', () => {
  const base = draftOf(nest())

  it('moves pages but never Home, and toggles all but Home', () => {
    expect(movePage(base, 'home', 3)).toBe(base)
    expect(keys(movePage(base, 'ops', 0))).toBe('home ops desk inbox')
    expect(togglePage(base, 'home')).toBe(base)
    expect(keys(togglePage(base, 'desk'))).toBe('home desk- ops inbox')
  })

  it('marks only the page that moved or switched', () => {
    const d = togglePage(movePage(base, 'inbox', 1), 'ops')
    expect(pageMark(base, d, 'inbox')).toBe('moved up')
    expect(pageMark(base, d, 'desk')).toBe('')
    expect(pageMark(base, d, 'ops')).toBe('hidden')
    expect(pendingCount(base, { ...d, theme: 'cyber' })).toBe(3)
  })

  it('saves only the changed keys; every page in known order saves as []', () => {
    expect(saveBody(base, base, known)).toEqual({})
    const theme = { ...base, theme: 'cyber' }
    expect(changedKeys(base, theme)).toEqual(['theme'])
    expect(needsBuild(['theme'])).toBe(false)
    expect(saveBody(base, theme, known)).toEqual({ theme: 'cyber' })
    const moved = movePage(base, 'ops', 1)
    expect(saveBody(base, moved, known)).toEqual({ pages: ['home', 'ops', 'desk', 'inbox'] })
    expect(needsBuild(changedKeys(base, moved))).toBe(true)
    const listed = draftOf(nest({ pages: ['home', 'desk'] }))
    expect(saveBody(listed, togglePage(togglePage(listed, 'ops'), 'inbox'), known)).toEqual({ pages: [] })
  })

  it('formats crow.yaml errors one per line', () => {
    expect(errorLines([{ key: 'hero', error: 'unknown' }, 'plain'])).toEqual(['hero: unknown', 'plain'])
  })
})

describe('apiSettings', () => {
  it('reads with GET and saves with the CSRF header', async () => {
    const seen: Array<{ url: string; method: string; crow: string | null; body?: string }> = []
    const api = apiSettings('/api', async (url, init) => {
      seen.push({ url, method: init?.method ?? 'GET', crow: new Headers(init?.headers).get('X-Crow'), body: init?.body as string | undefined })
      if (init?.method === 'POST') return new Response(JSON.stringify({ error: 'crow-modules: unknown theme "x"' }), { status: 400 })
      return new Response(JSON.stringify(url.endsWith('/build') ? { state: 'idle' } : nest()), { status: 200 })
    })
    expect((await api.get()).file).toBe('crow.yaml')
    expect((await api.build()).state).toBe('idle')
    const err = await api.save({ theme: 'x' }).catch((e: unknown) => e)
    expect(err).toBeInstanceOf(ApiError)
    expect((err as ApiError).status).toBe(400)
    expect((err as ApiError).message).toContain('unknown theme')
    expect(seen.map((s) => `${s.method} ${s.url} ${s.crow}`)).toEqual([
      'GET /api/settings/nest null',
      'GET /api/settings/nest/build null',
      'POST /api/settings/nest 1',
    ])
    expect(seen[2].body).toBe('{"theme":"x"}')
  })
})

describe('fixtureSettings', () => {
  it('refuses unknown names and builds after a pages save', async () => {
    const api = fixtureSettings(new Date('2026-11-25T10:00:00Z'), known)
    await expect(api.save({ hero: 'nope' })).rejects.toBeInstanceOf(ApiError)
    const theme = await api.save({ theme: 'cyber' })
    expect(theme).toMatchObject({ changed: ['theme'], needs_build: false, theme: 'cyber' })
    const r = await api.save({ pages: ['home', 'ops'] })
    expect(r).toMatchObject({ changed: ['pages'], needs_build: true, build: { state: 'building' } })
    expect((await api.build()).state).toBe('building')
    const done = await api.build()
    expect(done.state).toBe('ok')
    expect(done.served).not.toBe(r.build.served)
  })
})

describe('settingsControl', () => {
  it('loads, edits, saves and rebases the draft on the reply', async () => {
    let v: SettingsView = { status: 'loading', saving: false }
    const ctl = settingsControl(fixtureSettings(new Date('2026-11-25T10:00:00Z'), known), (n) => (v = n))
    await ctl.load()
    expect(v.status).toBe('ready')
    expect(v.build?.state).toBe('idle')
    ctl.edit((d) => togglePage(d, 'inbox'))
    expect(changedKeys(v.base!, v.draft!)).toEqual(['pages'])
    ctl.discard()
    expect(v.draft).toBe(v.base)
    ctl.edit((d) => ({ ...togglePage(d, 'inbox'), hero: 'none' }))
    await ctl.save()
    expect(v.saved).toEqual(['hero', 'pages'])
    expect(v.build?.state).toBe('building')
    expect(changedKeys(v.base!, v.draft!)).toEqual([])
    expect(keys(v.draft!)).toBe('home desk ops inbox-')
    await ctl.tick()
    await ctl.tick()
    expect(v.build?.state).toBe('ok')
  })

  it('keeps the draft and shows the refusal when a save fails', async () => {
    let v: SettingsView = { status: 'loading', saving: false }
    const ctl = settingsControl(fixtureSettings(new Date(), known), (n) => (v = n))
    await ctl.load()
    ctl.edit((d) => ({ ...d, theme: 'neon' }))
    await ctl.save()
    expect(v.saveError).toContain('unknown theme')
    expect(v.draft?.theme).toBe('neon')
  })

  it('says off without an api', async () => {
    let v: SettingsView | undefined
    await settingsControl(undefined, (n) => (v = n)).load()
    expect(v?.status).toBe('off')
  })
})
