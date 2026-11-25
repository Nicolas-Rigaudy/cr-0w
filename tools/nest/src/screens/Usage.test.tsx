import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { usageFixture } from '../data/usageFixture'
import type { UsageData } from '../data/usage'
import { Router } from '../router'
import { pickAccount, UsageView } from './Usage'

const at = (w: number, search = '', pathname = '/usage/acct1') => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search, pathname })
}
const page = (w: number, opts: { account?: string; search?: string; data?: UsageData } = {}) => {
  const src = fixtureSource()
  const now = src.now()
  const data = opts.data ?? usageFixture(now)
  at(w, opts.search)
  return renderToStaticMarkup(<DataContext.Provider value={src}><Router><UsageView data={data} account={opts.account} now={now} /></Router></DataContext.Provider>)
}

afterEach(() => { vi.unstubAllGlobals() })

describe('Usage page, laptop (item_227, option C)', () => {
  it('shows the head with the account seg and the week summary', () => {
    const h = page(1440, { account: 'acct1' })
    expect(h).toContain('u-lp')
    expect(h).toContain('<h1>Usage</h1>')
    expect(h).toMatch(/class="seg on" aria-current="page"[^>]*>acct1</)
    expect(h).toMatch(/week <span class="c-n">\$\d+<\/span> · 2 accounts · 5h resets 18:00/)
  })

  it('lists the accounts with their room left and a sparkline, then the tiles', () => {
    const h = page(1440, { account: 'acct1' })
    expect(h).toContain('[1] Accounts')
    expect(h).toContain('5h 38% left · wk 62%')
    expect(h).toContain('5h 91% left · wk 30%')
    expect(h).toContain('class="usp"')
    expect(h).toContain('22%/h')
    for (const t of ['7 days', 'Cache hit', 'Cache saved', 'Week forecast', 'Hub start ctx']) expect(h).toContain(`<small>${t}</small>`)
  })

  it('opens the picked account window pane with facts, both charts and the legend', () => {
    const h = page(1440, { account: 'acct1' })
    expect(h).toContain('[2] acct1 · 5h window')
    for (const k of ['5h left', 'Burn', 'Empty at', 'Resets', 'Week left', 'Week resets']) expect(h).toContain(`<dt>${k}</dt>`)
    expect(h).toContain('aria-label="acct1 5h window: 62% used now')
    expect(h).toContain('aria-label="acct1 week: 38% used now')
    expect(h).toContain('floor 10% left')
    expect(h).toContain('% at reset')
    expect(h).toContain('y: % of the window used')
  })

  it('shows the analysis with the per day chart and this week per project', () => {
    const h = page(1440, { account: 'acct1' })
    expect(h).toContain('[3] Analysis')
    expect(h).toContain('one chart at a time')
    expect(h).toContain('aria-label="$ per day"')
    expect(h).toContain('<h3 class="ust">7 days</h3>')
    expect(h).toContain('>api<')
  })

  it('follows the route account and falls back to the active one', () => {
    const h = page(1440, { account: 'acct2' })
    expect(h).toContain('[2] acct2 · 5h window')
    const d = usageFixture(fixtureSource().now())
    expect(pickAccount(d.limits, 'nope')?.account).toBe(d.limits.find((l) => l.active)?.account ?? d.limits[0].account)
  })

  it('takes the floor from the limit row when it carries one', () => {
    const d = usageFixture(fixtureSource().now())
    d.limits = d.limits.map((l) => (l.account === 'acct1' ? { ...l, min_room: 20 } : l))
    expect(page(1440, { account: 'acct1', data: d })).toContain('floor 20% left')
  })

  it('says so when no limits are known', () => {
    const d = usageFixture(fixtureSource().now())
    d.limits = []
    expect(page(1440, { data: d })).toContain('No account limits yet.')
  })
})

describe('Usage page, phone', () => {
  it('stacks the picker, the 5h window, the week and the analysis', () => {
    const h = page(368, { account: 'acct1' })
    expect(h).toContain('u-ph')
    expect(h).not.toContain('[1] Accounts')
    expect(h).toContain('[2] 5h window')
    for (const k of ['Left', 'Burn', 'Empty']) expect(h).toContain(`<dt>${k}</dt>`)
    expect(h).toContain('[3] Week')
    expect(h).toContain('[4] Analysis')
    expect(h).toContain('Per project')
    expect(h.indexOf('[2] 5h window')).toBeLessThan(h.indexOf('[4] Analysis'))
  })
})
