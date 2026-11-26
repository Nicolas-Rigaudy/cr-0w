import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { Router } from '../router'
import { ApiError } from '../data/api'
import { mrsOf } from '../data/mrs'
import { confirmsOf, Mrs } from './Mrs'

const snapshot = async (src: ReturnType<typeof fixtureSource>): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(src, (x) => { box.s = x })()
  return box.s!
}
const at = (w: number, search = '') => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search, pathname: '/mrs' })
}
const refresh = async () => {}
const page = async (w: number, search = '', s?: Snapshot) => {
  const src = fixtureSource()
  const snap = s ?? await snapshot(src)
  at(w, search)
  return renderToStaticMarkup(<DataContext.Provider value={src}><Router><Mrs snap={snap} refresh={refresh} /></Router></DataContext.Provider>)
}

afterEach(() => { vi.unstubAllGlobals() })

describe('MRs page, laptop (item_225, option C)', () => {
  it('shows the three deck panes with their counts and the header summary', async () => {
    const h = await page(1440)
    expect(h).toContain('m-lp')
    expect(h).toContain('To review</span><span class="num">1</span>')
    expect(h).toContain('Mine</span><span class="num">3</span>')
    expect(h).toContain('Merged today</span><span class="num">1</span>')
    expect(h).toContain('1 to review')
    expect(h).toContain('<span class="c-n">1 needs you</span>')
    expect(h).toContain('merged today 1')
  })

  it('opens the first Mine MR by default with its facts, pipeline and open threads', async () => {
    const h = await page(1440)
    expect(h).toMatch(/class="mrow mine on"[^>]*><span class="ref">!239/)
    expect(h).toContain('!239')
    expect(h).toContain('sync-loops-helper → main')
    expect(h).toContain('Pipeline #8812')
    expect(h).toContain('class="job failed">test ✕ 3:12')
    expect(h).toContain('test_sync.py::test_reset_drops_roles FAILED')
    expect(h).toContain('Same loop as in sync_grants. One helper for both?')
    // Resolved threads stay out of the overview.
    expect(h).not.toContain('Typo in the docstring.')
    expect(h).toContain('0 / 1')
  })

  it('lists the branch without an MR in Mine and merged rows with their time', async () => {
    const h = await page(1440)
    expect(h).toContain('proxy-rate-limit')
    expect(h).toContain('2 ahead · no MR')
    expect(h).toMatch(/merged \d\d:42/)
  })

  it('opens the MR and tab the URL names', async () => {
    const h = await page(1440, '?mr=242&tab=threads')
    expect(h).toContain('!242')
    expect(h).toContain('aria-selected="true">Threads')
    const t = await page(1440, '?mr=239&tab=threads')
    expect(t).toContain('Typo in the docstring.')
  })

  it('filters by project from the segment', async () => {
    const h = await page(1440, '?project=cr-0w')
    expect(h).toContain('#57')
    expect(h).not.toContain('!242')
  })

  it('shows the forge-off state with no rows', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const off = { ...s, sections: { ...s.sections, mrs: { data: null, stale: false, off: true, updated_at: null } } }
    const h = await page(1440, '', off)
    expect(h).toContain('The forge module is off.')
    expect(h).not.toContain('mrow')
  })

  it('says when the section did not refresh on both layouts', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const stale = { ...s, sections: { ...s.sections, mrs: { ...s.sections.mrs, stale: true, error: 'crow-mrs: timed out' } } }
    expect(await page(368, '', stale)).toContain('MRs did not refresh')
  })

  it('says when the section did not refresh', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const stale = { ...s, sections: { ...s.sections, mrs: { ...s.sections.mrs, stale: true, error: 'crow-mrs: timed out' } } }
    expect(await page(1440, '', stale)).toContain('MRs did not refresh: crow-mrs: timed out')
  })
})

describe('MRs row actions (item_225)', () => {
  const mrs = async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    return { src, d: mrsOf(s.sections.mrs, s.now) }
  }

  it('shows Answer threads, Repo session and Open on forge on your MR, Review on the others', async () => {
    const h = await page(1440)
    expect(h).toContain('class="btn p">Answer threads</button>')
    expect(h).toContain('>Repo session</button>')
    expect(h).toContain('href="https://gitlab.example.com/team/api/-/merge_requests/239" target="_blank" rel="noreferrer">Open on forge ↗</a>')
    expect(h).toContain('Replies post after your OK in the threads session.')
    const r = await page(1440, '?mr=242')
    expect(r).toContain('class="btn p">Review</button>')
    expect(r).not.toContain('Answer threads')
    const p = await page(368, '?mr=239')
    expect(p).toContain('>Session</button>')
    expect(p).toContain('>Forge ↗</a>')
  })

  it('starts a review run, a threads session, a repo session and an MR session through the live actions', async () => {
    const { src, d } = await mrs()
    const said: string[] = []
    const refresh = vi.fn(async () => {})
    const c = confirmsOf(src.actions, refresh, (s) => said.push(s.text))
    const dispatch = vi.spyOn(src.actions, 'dispatch')
    const start = vi.spyOn(src.actions, 'start')
    const review = d.review.find((m) => m.ref === '!242')!
    const mine = d.mine.find((m) => m.ref === '!239')!
    await c.review(review).run()
    expect(dispatch).toHaveBeenCalledWith({ project: 'api', label: 'review-api-242', prompt: '/review https://gitlab.example.com/team/api/-/merge_requests/242' })
    await c.threads(mine).run()
    expect(start).toHaveBeenCalledWith({ project: 'api', prompt: '/review-threads https://gitlab.example.com/team/api/-/merge_requests/239' })
    await c.session(mine).run()
    expect(start).toHaveBeenCalledWith({ project: 'api' })
    await c.branch(d.branches[0]).run()
    expect(start).toHaveBeenCalledWith({ project: 'api', prompt: '/mr proxy-rate-limit' })
    expect(said[0]).toMatch(/^dispatched run /)
    expect(refresh).toHaveBeenCalledTimes(4)
  })

  it('refreshes and rethrows when the run slots are full', async () => {
    const { src, d } = await mrs()
    const refresh = vi.fn(async () => {})
    const full = { ...src.actions, dispatch: async () => { throw new ApiError(409, 'too_many_runs') } }
    const c = confirmsOf(full, refresh, () => {})
    await expect(c.review(d.review[0]).run()).rejects.toThrow('too_many_runs')
    expect(refresh).toHaveBeenCalledTimes(1)
  })
})

describe('MRs review patterns (item_225)', () => {
  it('stays a tab with the hit count until picked', async () => {
    const h = await page(1440)
    expect(h).toContain('>Patterns ◆2</button>')
    expect(h).not.toContain('class="pats"')
  })

  it('lists the common and project ledger newest first with the hits marked', async () => {
    const h = await page(1440, '?mr=239&tab=patterns')
    expect(h).toContain('Review patterns · api')
    const order = [...h.matchAll(/<tr[^>]*><td>(?:◆ )?([^<]+)<\/td><td>(\d+)<\/td>/g)].map((x) => Number(x[2]))
    expect(order).toEqual([1, 7, 9, 10, 5, 4])
    expect(h).toContain('<tr class="hit"><td>◆ Same loop copied instead of one helper</td>')
    expect(h).toContain('common + project ledger')
  })

  it('shows no hit mark on an MR the ledgers do not name, and the tab on phone', async () => {
    const h = await page(1440, '?mr=242&tab=patterns')
    expect(h).toContain('>Patterns</button>')
    expect(h).not.toContain('class="hit"')
    const p = await page(368, '?mr=239&tab=patterns')
    expect(p).toContain('Review patterns · api')
  })
})

describe('MRs page, phone (item_225, option C)', () => {
  it('shows the pane tabs and only the picked pane, rows with project and forge', async () => {
    const h = await page(368)
    expect(h).toContain('m-ph')
    expect(h).toContain('aria-selected="true">Mine 3</button>')
    expect(h).toContain('Review 1')
    expect(h).toContain('>Merged</button>')
    expect(h).toContain('<span class="c-n">1 needs you</span> · <span>1 review</span>')
    expect(h).toContain('<span class="sub">api · gitlab</span>')
    expect(h).toContain('<span class="sub">cr-0w · github</span>')
    expect(h).toContain('2 ahead · no MR')
    // One pane at a time and no detail until a row is tapped.
    expect(h).not.toContain('raise alerts from detection events')
    expect(h).not.toContain('class="mdet')
  })

  it('opens the MR the URL names as its own page with a back link', async () => {
    const h = await page(368, '?mr=239')
    expect(h).toContain('◀ Mine')
    expect(h).toContain('api · sync-loops-helper → main · 2h')
    expect(h).toContain('<dt>CI</dt><dd><span class="c-f">✕ test</span></dd>')
    expect(h).toContain('<dt>Approvals</dt><dd>0/1</dd>')
    expect(h).toContain('class="job failed">test ✕ 3:12')
    expect(h).toContain('Same loop as in sync_grants. One helper for both?')
    expect(h).not.toContain('class="deck')
  })

  it('uses the phone layout below the laptop width', async () => {
    expect(await page(707)).toContain('m-ph')
  })
})
