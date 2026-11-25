import { describe, expect, it } from 'vitest'
import { match, usagePath } from '../router'
import { barTab, hiddenTabs, screenTab, shownRoute } from './tabs'
import { NAV, navKeys, splitBar } from '../modules/nav'
import type { Snapshot } from '../data/feed'
import { TABS } from './icons'
import { DEMO_BRICKS } from '../data/fixtures'

const bar = (path: string) => barTab(screenTab(match(path)))

describe('phone tab bar B', () => {
  it('highlights the matching tab for the four bar screens', () => {
    expect(bar('/')).toBe('home')
    expect(bar('/desk')).toBe('desk')
    expect(bar('/ops')).toBe('ops')
    expect(bar('/missions')).toBe('ops')
    expect(bar('/inbox')).toBe('inbox')
  })

  it('highlights More while MRs, Agenda or Reports is open', () => {
    expect(bar('/mrs')).toBe('more')
    expect(bar('/agenda')).toBe('more')
    expect(bar('/reports')).toBe('more')
  })

  it('keeps chats under Desk, old mission links under Ops, unknown pages under Home', () => {
    expect(bar('/desk/4242')).toBe('desk')
    expect(bar('/desk/m/20261109-build-inbox-filters')).toBe('desk')
    expect(bar('/missions/20261109-build-inbox-filters')).toBe('ops')
    expect(bar('/nope')).toBe('home')
  })

  it('opens Usage on /usage and /usage/<account>, under Home', () => {
    expect(match('/usage')).toEqual({ tab: 'usage' })
    expect(match('/usage/')).toEqual({ tab: 'usage' })
    expect(match('/usage/acct%201')).toEqual({ tab: 'usage', account: 'acct 1' })
    expect(match('/usage/a/b')).toEqual({ tab: 'missing' })
    expect(usagePath('acct 1')).toBe('/usage/acct%201')
    expect(bar('/usage/acct1')).toBe('home')
    expect(shownRoute(match('/usage/acct1'))).toEqual({ tab: 'usage', account: 'acct1' })
  })

  it('keeps today\'s bar and More with no nest config', () => {
    expect(splitBar(NAV)).toEqual({ bar: ['home', 'desk', 'ops', 'inbox'], more: ['mrs', 'agenda', 'reports'] })
    expect([...NAV].sort()).toEqual(TABS.map((t) => t.key).sort())
  })
})

describe('tabs of an off brick', () => {
  const snap = (forgeOn: boolean) => ({ ...DEMO_BRICKS, modules: { ...DEMO_BRICKS.modules, modules: DEMO_BRICKS.modules.modules.map((m) => (m.name === 'forge' ? { ...m, on: forgeOn } : m)) } }) as unknown as Snapshot
  it('hides MRs while the forge brick is off', () => {
    expect([...hiddenTabs(snap(false))]).toEqual(['mrs'])
  })
  it('keeps every tab with forge on, no module list or no snapshot', () => {
    expect(hiddenTabs(snap(true)).size).toBe(0)
    expect(hiddenTabs({ modules: null } as unknown as Snapshot).size).toBe(0)
    expect(hiddenTabs(null).size).toBe(0)
  })
})

describe('nest.pages (item_267)', () => {
  const all = ['home', 'desk', 'ops', 'inbox', 'mrs', 'agenda', 'reports', 'garden']
  it('keeps every page in default order when the list is empty', () => {
    expect(navKeys([], all)).toEqual(all)
  })
  it('shows only the listed pages, in list order, Home always kept', () => {
    expect(navKeys(['home', 'inbox'], all)).toEqual(['home', 'inbox'])
    expect(navKeys(['garden', 'desk'], all)).toEqual(['home', 'garden', 'desk'])
  })
  it('fills the bar with the first four shown pages and skips hidden ones', () => {
    const nav = ['home', 'garden', 'mrs', 'desk', 'inbox', 'ops']
    expect(splitBar(nav)).toEqual({ bar: ['home', 'garden', 'mrs', 'desk'], more: ['inbox', 'ops'] })
    expect(splitBar(nav, new Set(['mrs']))).toEqual({ bar: ['home', 'garden', 'desk', 'inbox'], more: ['ops'] })
    expect(splitBar(['home', 'inbox'])).toEqual({ bar: ['home', 'inbox'], more: [] })
  })
  it('sends an unlisted page to Home and keeps chats with Desk', () => {
    const nav = ['home', 'inbox', 'desk']
    expect(shownRoute(match('/ops'), nav)).toEqual({ tab: 'home' })
    expect(shownRoute(match('/inbox'), nav)).toEqual(match('/inbox'))
    expect(shownRoute(match('/desk/4242'), nav)).toEqual(match('/desk/4242'))
    expect(shownRoute(match('/nope'), nav)).toEqual(match('/nope'))
  })
  it('routes a brick page key', () => {
    expect(match('/garden', ['garden'])).toEqual({ tab: 'brick', key: 'garden' })
    expect(barTab(screenTab(match('/garden', ['garden'])), ['home', 'garden'])).toBe('garden')
  })
})
