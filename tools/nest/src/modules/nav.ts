// Nav pages (item_267): engine tabs plus brick pages, in crow.yaml modules.nest.pages order.
import { createElement } from 'react'
import { config, pages as brickPages } from 'virtual:nest-modules'
import { TABS, type Tab } from '../ui/icons'

const BAR_SIZE = 4

const title = (k: string) => k.charAt(0).toUpperCase() + k.slice(1).replace(/-/g, ' ')
const PAGE_ICON = createElement('rect', { x: 3.5, y: 3.5, width: 15, height: 15, rx: 2.5 })

export const BRICK_TABS: Tab[] = brickPages.map((p) => ({ key: p.key, label: p.mod.label ?? title(p.key), path: `/${p.key}`, icon: p.mod.icon ?? PAGE_ICON }))
export const ALL_TABS: Tab[] = [...TABS, ...BRICK_TABS]
export const BRICK_PAGE_KEYS: readonly string[] = brickPages.map((p) => p.key)

// Pages in nav order. An empty list keeps every page in the default order. Home always
// stays: it is where unlisted routes land.
export function navKeys(listed: readonly string[], all: readonly string[] = ALL_TABS.map((t) => t.key)): string[] {
  if (!listed.length) return [...all]
  const keys = listed.filter((k) => all.includes(k))
  return keys.includes('home') ? keys : ['home', ...keys]
}

export const NAV: readonly string[] = navKeys(config.pages)

// Phone tab bar: the first four shown pages, the rest in the More sheet.
export function splitBar(nav: readonly string[], hidden: ReadonlySet<string> = new Set()): { bar: string[]; more: string[] } {
  const shown = nav.filter((k) => !hidden.has(k))
  return { bar: shown.slice(0, BAR_SIZE), more: shown.slice(BAR_SIZE) }
}
