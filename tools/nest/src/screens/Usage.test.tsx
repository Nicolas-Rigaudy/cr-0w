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

  it('offers the six tabs, the unit toggle and the D/W/M period', () => {
    const h = page(1440, { account: 'acct1' })
    for (const t of ['Per day', 'Per project', 'Per model', 'Per kind', 'Per session', 'Top runs']) expect(h).toContain(`>${t}</button>`)
    expect(h).toMatch(/aria-label="Unit".*>\$<.*>tokens<.*>limits</)
    expect(h).toMatch(/aria-label="Period".*>D<.*aria-pressed="true"[^>]*>W<.*>M</)
  })

  it('puts the 5h cells under each day and the picked day summary line', () => {
    const h = page(1440, { account: 'acct1' })
    expect(h).toContain('class="ucell"')
    expect(h).toContain('5h cells: acct1')
    expect(h).toMatch(/class="uday"[^>]*><b>\w{3} \d\d\/\d\d<\/b> \$\d+ · [\d.]+M tok · \d+% cache · week −\d+% · 5h: /)
  })

  it('switches the per day unit to tokens and limits from the query', () => {
    expect(page(1440, { account: 'acct1', search: '?unit=tok' })).toContain('aria-label="tokens per day"')
    const h = page(1440, { account: 'acct1', search: '?unit=lim' })
    expect(h).toContain('aria-label="% of a weekly limit used per day"')
    expect(h).toContain('100% = one account-week')
  })

  it('shows limits as no data on days before the samples', () => {
    const d = usageFixture(fixtureSource().now())
    d.cost.data!.by_day = d.cost.data!.by_day.map((x, i, a) => (i < a.length - 2 ? { ...x, limits: null } : x))
    expect(page(1440, { account: 'acct1', data: d, search: '?unit=lim' })).toContain('>no data<')
  })

  it('splits per project, model and kind in $ and tokens', () => {
    expect(page(1440, { account: 'acct1', search: '?tab=project' })).toContain('$ per project · 7 days')
    const m = page(1440, { account: 'acct1', search: '?tab=model' })
    expect(m).toContain('tokens per model')
    expect(m).toContain('>opus<')
    expect(page(1440, { account: 'acct1', search: '?tab=kind' })).toContain('$ per kind')
  })

  it('lists the sessions of the period with full title, model, tokens and cost', () => {
    const h = page(1440, { account: 'acct1', search: '?tab=session' })
    for (const c of ['Session', 'Project', 'Model', 'Tokens', 'Cost']) expect(h).toContain(`>${c}</span>`)
    expect(h).toContain('Refactor the auth middleware')
    expect(h).toContain('opus + sonnet')
  })

  it('ranks the top runs by cost', () => {
    const h = page(1440, { account: 'acct1', search: '?tab=runs' })
    expect(h.indexOf('Refactor the auth middleware')).toBeLessThan(h.indexOf('Plan the plugin loader'))
    expect(h).toContain('$31')
  })

  it('cuts the period with D', () => {
    const h = page(1440, { account: 'acct1', search: '?tab=session&period=D' })
    expect(h).toContain('Review the billing endpoints')
    expect(h).not.toContain('Rewrite the completion engine')
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
    for (const t of ['Day', 'Project', 'Model', 'Kind', 'Session', 'Runs']) expect(h).toContain(`>${t}</button>`)
  })

  it('lists sessions as two-line rows on the phone', () => {
    const h = page(368, { account: 'acct1', search: '?tab=session' })
    expect(h).toContain('utbl uph')
    expect(h).not.toMatch(/class="([^"]* )?ph[ "]/)
    expect(h).toMatch(/class="sub">api · opus \+ sonnet · /)
  })

  it('marks the active session without the global live hatch class', () => {
    for (const tab of ['session', 'runs']) {
      const h = page(1280, { account: 'acct1', search: `?tab=${tab}` })
      expect(h).toContain('utr uact')
      expect(h).not.toMatch(/class="[^"]*\blive\b/)
    }
  })

  it('draws a cost bar on each session and run row, the top cost at full width', () => {
    for (const tab of ['session', 'runs']) {
      for (const w of [368, 1280]) {
        const h = page(w, { account: 'acct1', search: `?tab=${tab}` })
        const n = (h.match(/class="ucb"/g) ?? []).length
        expect(n).toBeGreaterThan(1)
        expect(h).toContain('width:100%')
      }
    }
  })
})
