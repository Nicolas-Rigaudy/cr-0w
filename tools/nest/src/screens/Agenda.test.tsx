import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { Router } from '../router'
import { ApiError } from '../data/api'
import { agendaOf } from '../data/agenda'
import type { SectionBody } from '../data/types'
import { Agenda, confirmsOf, HUB } from './Agenda'

const snapshot = async (src: ReturnType<typeof fixtureSource>): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(src, (x) => { box.s = x })()
  return box.s!
}
const at = (w: number, search = '') => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search, pathname: '/agenda' })
}
const refresh = async () => {}
const page = async (w: number, search = '', s?: Snapshot) => {
  const src = fixtureSource()
  const snap = s ?? await snapshot(src)
  at(w, search)
  return renderToStaticMarkup(<DataContext.Provider value={src}><Router><Agenda snap={snap} refresh={refresh} /></Router></DataContext.Provider>)
}
const withToday = (s: Snapshot, today: SectionBody): Snapshot => ({ ...s, sections: { ...s.sections, today } })

afterEach(() => { vi.unstubAllGlobals() })

describe('Agenda page, laptop (item_226, option B)', () => {
  it('shows the header summary and the week grid from today', async () => {
    const h = await page(1440)
    expect(h).toContain('desk agenda m-lp')
    expect(h).toContain('today 3 meetings · 1 due · week 11 meetings · 5 due')
    expect(h).toContain('Week 46')
    expect(h).toContain('class="dh on">Mon 09</span>')
    expect(h).toContain('class="ad">Remote</span>')
    expect(h).toContain('times: Paris')
  })

  it('places meetings as blocks with their state and the now line on today', async () => {
    const h = await page(1440)
    expect(h).toMatch(/class="ev done"[^>]*>11:30 Daily Meeting/)
    expect(h).toMatch(/class="ev next on"[^>]*>17:00 Sprint review/)
    expect(h).toContain('class="nowl"')
    expect(h).not.toContain('Management weekly')
  })

  it('puts due chips in the due row by kind', async () => {
    const h = await page(1440)
    expect(h).toContain('class="due mr"')
    expect(h).toContain('⇄ !242 review waits 4d</button>')
    expect(h).toContain('⚑ Sprint 34 ends Fri 13 Nov</button>')
    expect(h).toContain('class="due th"')
  })

  it('opens the next meeting card and the due list', async () => {
    const h = await page(1440)
    expect(h).toContain('Next · in 50 min')
    expect(h).toContain('class="ttl">Sprint review</div>')
    expect(h).toContain('class="btn p">Prep session</button>')
    expect(h).toContain('Due this week</span><span class="c-s">5</span>')
    expect(h).toContain('class="btn sm">Review</button>')
    expect(h).toContain('class="btn sm">Open</button>')
  })
})

describe('Agenda page, phone (item_226, option B)', () => {
  it('shows the day strip on today with its due chip and hour grid', async () => {
    const h = await page(368)
    expect(h).toContain('desk agenda m-ph')
    expect(h).toContain('W46')
    expect(h).toMatch(/class="dtab on" aria-selected="true"><span>Mon<\/span><b>09<\/b>/)
    expect(h).toContain('<span class="act">REVIEW</span>')
    expect(h).toContain('class="hg"')
    expect(h).toContain('now 16:10')
    expect(h).not.toContain('class="wg"')
  })

  it('picks another day from ?day=', async () => {
    const h = await page(368, '?day=4')
    expect(h).toMatch(/class="dtab on" aria-selected="true"><span>Fri<\/span><b>13<\/b>/)
    expect(h).toContain('⚑ Sprint 34 ends Fri 13 Nov')
    expect(h).not.toContain('now 16:10')
  })
})

describe('Agenda states (item_226)', () => {
  it('says when the calendar is off or did not refresh', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const off = await page(1440, '', withToday(s, { data: null, stale: false, updated_at: null, off: true }))
    expect(off).toContain('The mail module is off.')
    const stale = await page(1440, '', withToday(s, { ...s.sections.today, stale: true, error: 'not signed in' }))
    expect(stale).toContain('The calendar did not refresh: not signed in')
  })
})

describe('Agenda actions (item_226)', () => {
  const ctx = async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    return { src, a: agendaOf(s) }
  }

  it('preps a meeting and opens a due check in a hub session', async () => {
    const { src, a } = await ctx()
    const said: string[] = []
    const r = vi.fn(async () => {})
    const opened: string[] = []
    const c = confirmsOf(src.actions, r, (s) => said.push(s.text), (t) => opened.push(t))
    const start = vi.spyOn(src.actions, 'start')
    await c.prep(a.next!, a.days[0]).run()
    expect(start).toHaveBeenCalledWith(expect.objectContaining({ project: HUB, prompt: expect.stringContaining('Sprint review') }))
    const x = a.days[1].due[0]
    await c.due(x, a.days[1]).run()
    expect(start).toHaveBeenCalledWith({ project: HUB, prompt: 'Legal: contract questions. unanswered for 3d. Next step: draft a reply.' })
    expect(r).toHaveBeenCalledTimes(2)
    expect(said).toHaveLength(2)
    expect(opened).toHaveLength(2)
  })

  it('sends an MR wait to the MRs review run', async () => {
    const { src, a } = await ctx()
    const dispatch = vi.spyOn(src.actions, 'dispatch')
    const open = vi.fn()
    const c = confirmsOf(src.actions, refresh, () => {}, open)
    await c.due(a.days[0].due[0], a.days[0]).run()
    expect(open).not.toHaveBeenCalled()
    expect(dispatch).toHaveBeenCalledWith(expect.objectContaining({ project: 'api', label: 'review-api-242' }))
  })

  it('refreshes and rethrows when the session slots are full', async () => {
    const { src, a } = await ctx()
    const r = vi.fn(async () => {})
    const full = { ...src.actions, start: async () => { throw new ApiError(409, 'too_many_sessions') } }
    const c = confirmsOf(full, r, () => {})
    await expect(c.prep(a.next!, a.days[0]).run()).rejects.toThrow('too_many_sessions')
    expect(r).toHaveBeenCalledTimes(1)
  })
})
