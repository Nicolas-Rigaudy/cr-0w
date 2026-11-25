import { useSyncExternalStore } from 'react'
import { config, themes } from 'virtual:nest-modules'

// The theme the owner picked (header switcher or Settings save), shared by every reader.
// index.html applies the stored one at first paint; crow.yaml modules.nest.theme is the default.
let current: string | null = null
const subs = new Set<() => void>()

function read(): string {
  if (current) return current
  const t = typeof document === 'undefined' ? undefined : document.documentElement?.dataset?.theme
  return (current = t && themes.includes(t) ? t : config.theme)
}

// Paints a theme on the page without keeping it (Settings preview before save).
export function paintTheme(t: string) {
  const root = document.documentElement
  root.dataset.theme = t
  const bg = getComputedStyle(root).getPropertyValue('--bg').trim()
  if (bg) document.querySelector('meta[name="theme-color"]')?.setAttribute('content', bg)
}

export function setTheme(t: string) {
  current = t
  paintTheme(t)
  try { localStorage.setItem('nest.theme', t) } catch { /* storage blocked: theme stays per visit */ }
  for (const f of subs) f()
}

// Drops a preview: the kept theme paints again.
export const unpaintTheme = () => paintTheme(read())

export function useTheme(): string {
  return useSyncExternalStore(
    (f) => { subs.add(f); return () => { subs.delete(f) } },
    read,
    read,
  )
}
