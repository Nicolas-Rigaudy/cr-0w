import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { Router } from '../router'
import { Mrs } from './Mrs'

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
    expect(h).toContain('[1] To review</span><span class="num">2</span>')
    expect(h).toContain('[2] Mine</span><span class="num">3</span>')
    expect(h).toContain('[3] Merged today</span><span class="num">1</span>')
    expect(h).toContain('1 to review')
    expect(h).toContain('<span class="c-n">1 needs you</span>')
    expect(h).toContain('merged today 1')
  })

  it('opens the first Mine MR by default with its facts, pipeline and open threads', async () => {
    const h = await page(1440)
    expect(h).toMatch(/class="mrow mine on"[^>]*><span class="ref">!239/)
    expect(h).toContain('[4] !239')
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
    expect(h).toContain('[4] !242')
    expect(h).toContain('aria-selected="true">Threads')
    const t = await page(1440, '?mr=239&tab=threads')
    expect(t).toContain('Typo in the docstring.')
  })

  it('filters by project from the segment', async () => {
    const h = await page(1440, '?project=cr-0w')
    expect(h).toContain('[4] #57')
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

describe('MRs page, phone (item_225, option C)', () => {
  it('shows the pane tabs and only the picked pane, rows with project and forge', async () => {
    const h = await page(368)
    expect(h).toContain('m-ph')
    expect(h).toContain('aria-selected="true">[2] Mine 3</button>')
    expect(h).toContain('[1] Review 2')
    expect(h).toContain('>[3] Merged</button>')
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
    expect(h).toContain('◀ [2] Mine')
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
