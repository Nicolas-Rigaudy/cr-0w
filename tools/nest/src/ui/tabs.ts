import type { Snapshot } from '../data/feed'
import { hiddenPages } from '../data/modules'
import { NAV, splitBar } from '../modules/nav'
import type { Route } from '../router'
import type { TabKey } from './icons'

export type BarKey = TabKey | 'more'

// The screen a route belongs to: a chat sits under Desk, Usage and an unknown page under Home.
export function screenTab(route: Route): TabKey {
  if (route.tab === 'missing' || route.tab === 'usage') return 'home'
  if (route.tab === 'chat') return 'desk'
  if (route.tab === 'brick') return route.key
  return route.tab
}

// The highlighted bar tab (phone bar B, item_206): More while one of its screens is open.
export function barTab(active: TabKey, nav: readonly string[] = NAV, hidden?: ReadonlySet<string>): BarKey {
  return splitBar(nav, hidden).more.includes(active) ? 'more' : active
}

// Screens of bricks that are off (adr_002, item_237): a brick page shows when at least one
// of its sections is on, so MRs leaves the pills and the More sheet while forge is off.
export function hiddenTabs(snap: Snapshot | null): ReadonlySet<TabKey> {
  return hiddenPages(snap?.modules ?? null)
}

// The screen a route shows (item_267): a page left out of crow.yaml modules.nest.pages,
// or a chat while Desk is left out, lands on Home.
export function shownRoute(route: Route, nav: readonly string[] = NAV): Route {
  return nav.includes(screenTab(route)) ? route : { tab: 'home' }
}
