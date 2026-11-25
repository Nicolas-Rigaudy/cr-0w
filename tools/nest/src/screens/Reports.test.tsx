import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { Router } from '../router'
import { ApiError } from '../data/api'
import { Reports, retroConfirm } from './Reports'

const snapshot = async (src: ReturnType<typeof fixtureSource>): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(src, (x) => { box.s = x })()
  return box.s!
}
const at = (w: number, search = '') => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search, pathname: '/reports' })
}
const refresh = async () => {}
const page = async (w: number, search = '', s?: Snapshot) => {
  const src = fixtureSource()
  const snap = s ?? await snapshot(src)
  at(w, search)
  return renderToStaticMarkup(<DataContext.Provider value={src}><Router><Reports snap={snap} refresh={refresh} /></Router></DataContext.Provider>)
}

afterEach(() => { vi.unstubAllGlobals() })

describe('Reports page, laptop (item_258, board B)', () => {
  it('shows the six ritual tiles with their counts, last output and foot', async () => {
    const h = await page(1440)
    expect(h).toContain('r-lp')
    expect(h).toContain('settled only · open decisions stay in Inbox')
    for (const n of ['AFK', 'Brief', 'Daily', 'Weekly', 'Retro', 'Eod']) expect(h).toContain(` ${n}</span>`)
    expect(h).toContain('from afk runs')
    expect(h).toContain('1 this week')
    // Brief, daily and eod save no output; weekly has output but no routine (item_234).
    expect(h.match(/no saved output/g)).toHaveLength(3)
    expect(h).toContain('no routine yet')
    // Retro is the only ritual with a Generate button.
    expect(h.match(/↻ Generate/g)).toHaveLength(1)
  })

  it('lists settled outputs newest first and opens the newest in the reader', async () => {
    const h = await page(1440)
    expect(h.indexOf('AFK night → 2026-11-09')).toBeLessThan(h.indexOf('Retro batch · 6 sessions'))
    expect(h.indexOf('Retro batch · 6 sessions')).toBeLessThan(h.indexOf('Week 45 rollup'))
    expect(h).toMatch(/class="rrow on"[^>]*><span class="rg c-m"[^>]*>☾<\/span><span class="rt">AFK night → 2026-11-09/)
    expect(h).toContain('5 points · 7 runs · no decision · Mon 09 07:55')
    expect(h).toContain('<h3>crow</h3>')
    expect(h).toContain('<b>Demo login</b>')
    expect(h).toContain('1 AFK report waits for a decision in Inbox')
  })

  it('filters by the ritual in the URL and shows the vault link of a note', async () => {
    const h = await page(1440, '?kind=retro')
    expect(h).toContain('class="rtile on"')
    expect(h).not.toContain('Week 45 rollup')
    expect(h).toContain('href="obsidian://open?vault=vault&amp;file=CR-0W%2FRetros%2F2026-11-08%20retro%20batch"')
    expect(h).toContain('<h4>Sessions</h4>')
    // The note title is the reader header, not repeated in the body.
    expect(h.match(/Retro batch · 6 sessions/g)!.length).toBe(2)
  })

  it('says when a ritual has no saved output', async () => {
    const h = await page(1440, '?kind=brief')
    expect(h).toContain('No saved brief output yet.')
    expect(h).toContain('Pick a report.')
  })

  it('hides Generate when the retro routine is not installed', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const snap = { ...s, sections: { ...s.sections, jobs: { data: [], stale: false, updated_at: '' } } }
    const h = await page(1440, '', snap)
    expect(h).not.toContain('↻ Generate')
  })
})

describe('Reports page, phone', () => {
  it('shows the 2×3 tiles and the newest five', async () => {
    const h = await page(368)
    expect(h).toContain('r-ph')
    expect(h).not.toContain('this week')
    expect(h).toContain('↻ Gen</button>')
    expect(h.match(/class="rrow/g)).toHaveLength(5)
    expect(h).not.toContain('class="rph"')
  })

  it('opens the report the URL names as its own page with a back button', async () => {
    const h = await page(368, '?id=weekly:2026-W45%20rollup')
    expect(h).toContain('◀ Reports')
    expect(h).toContain('Week 45 rollup')
    expect(h).toContain('Slides drafted, 3 demos · Fri 06 09:00')
    expect(h).not.toContain('class="rtiles"')
  })
})

describe('Retro Generate', () => {
  it('starts the retro routine through core and refreshes', async () => {
    const src = fixtureSource()
    const said: string[] = []
    const started = vi.fn()
    const refresh = vi.fn(async () => {})
    const routine = vi.spyOn(src.actions, 'routine')
    const c = retroConfirm(src.actions, refresh, (s) => said.push(s.text), started)
    expect(c.target).toContainEqual(['Routine', 'crow-retro.timer'])
    await c.run()
    expect(routine).toHaveBeenCalledWith('crow-retro.timer', 'run')
    expect(started).toHaveBeenCalledTimes(1)
    expect(refresh).toHaveBeenCalledTimes(1)
    expect(said).toHaveLength(1)
  })

  it('refreshes on a 409 and rethrows', async () => {
    const src = fixtureSource()
    const refresh = vi.fn(async () => {})
    const busy = { ...src.actions, routine: async () => { throw new ApiError(409, 'running') } }
    const c = retroConfirm(busy, refresh, () => {}, () => {})
    await expect(c.run()).rejects.toThrow('running')
    expect(refresh).toHaveBeenCalledTimes(1)
  })
})
