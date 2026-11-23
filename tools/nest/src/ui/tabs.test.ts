import { describe, expect, it } from 'vitest'
import { match } from '../router'
import { BAR_KEYS, MORE_KEYS, barTab, hiddenTabs, screenTab } from './tabs'
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

  it('splits every tab once between the bar and More', () => {
    expect([...BAR_KEYS, ...MORE_KEYS].sort()).toEqual(TABS.map((t) => t.key).sort())
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
