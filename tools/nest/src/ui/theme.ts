import { useSyncExternalStore } from 'react'
import { config, themes } from 'virtual:nest-modules'
import type { Mode } from '../data/settings'

export type { Mode }

// The theme the owner picked (header switcher or Settings save), shared by every reader.
// index.html applies the stored one at first paint; crow.yaml modules.nest.theme is the default.
let current: string | null = null
const subs = new Set<() => void>()
const notify = () => { for (const f of subs) f() }
const subscribe = (f: () => void) => { subs.add(f); return () => { subs.delete(f) } }

function read(): string {
  if (current) return current
  const t = typeof document === 'undefined' ? undefined : document.documentElement?.dataset?.theme
  return (current = t && themes.includes(t) ? t : config.theme)
}

// The browser bar takes the painted background.
function syncBar() {
  const bg = getComputedStyle(document.documentElement).getPropertyValue('--bg').trim()
  if (bg) document.querySelector('meta[name="theme-color"]')?.setAttribute('content', bg)
}

// Paints a theme on the page without keeping it (Settings preview before save).
export function paintTheme(t: string) {
  document.documentElement.dataset.theme = t
  syncBar()
}

export function setTheme(t: string) {
  current = t
  paintTheme(t)
  try { localStorage.setItem('nest.theme', t) } catch { /* storage blocked: theme stays per visit */ }
  notify()
}

// Drops a preview: the kept theme paints again.
export const unpaintTheme = () => paintTheme(read())

export function useTheme(): string {
  return useSyncExternalStore(subscribe, read, read)
}

// Mode (item_303): auto follows the system, light and dark force it. Kept per device in
// localStorage, never in crow.yaml. `<html data-mode>` holds the resolved light or dark.
// A theme with no light block keeps one palette whatever the mode.
export type Painted = 'light' | 'dark'
export const MODES: readonly Mode[] = ['auto', 'light', 'dark']
export const MODE_KEY = 'nest.mode'
const DARK_QUERY = '(prefers-color-scheme: dark)'

export const isMode = (m: unknown): m is Mode => typeof m === 'string' && (MODES as readonly string[]).includes(m)
export const resolveMode = (m: Mode, systemDark: boolean): Painted => (m === 'auto' ? (systemDark ? 'dark' : 'light') : m)

let mode: Mode | null = null
// The mode a Settings preview paints until Save or Discard.
let preview: Mode | null = null

function storedMode(): Mode {
  try {
    const m = localStorage.getItem(MODE_KEY)
    return isMode(m) ? m : 'auto'
  } catch { return 'auto' }
}

const readMode = (): Mode => (mode ??= storedMode())
// The kept choice outside React (Settings loads its draft from it).
export const readKeptMode = readMode

// No matchMedia (tests, old browsers): dark, the look every theme has.
function systemDark(): boolean {
  return typeof matchMedia === 'function' ? matchMedia(DARK_QUERY).matches : true
}

function applyMode(m: Mode) {
  const p = resolveMode(m, systemDark())
  if (document.documentElement.dataset.mode === p) return
  document.documentElement.dataset.mode = p
  syncBar()
}

// Paints a mode without keeping it (Settings preview).
export function paintMode(m: Mode) {
  preview = m
  applyMode(m)
  notify()
}

export function setMode(m: Mode) {
  mode = m
  preview = null
  try { localStorage.setItem(MODE_KEY, m) } catch { /* storage blocked: mode stays per visit */ }
  applyMode(m)
  notify()
}

// Drops a preview: the kept mode paints again.
export function unpaintMode() {
  if (preview === null) return
  preview = null
  applyMode(readMode())
  notify()
}

// While the painted mode is auto, a system switch repaints at once.
export function watchSystemMode(): () => void {
  if (typeof matchMedia !== 'function') return () => {}
  const mq = matchMedia(DARK_QUERY)
  const on = () => {
    if ((preview ?? readMode()) !== 'auto') return
    applyMode('auto')
    notify()
  }
  mq.addEventListener('change', on)
  return () => mq.removeEventListener('change', on)
}

// The kept choice (auto, light or dark).
export function useMode(): Mode {
  return useSyncExternalStore(subscribe, readMode, readMode)
}

const readPainted = (): Painted => resolveMode(preview ?? readMode(), systemDark())

// What the page paints now: light or dark.
export function useResolvedMode(): Painted {
  return useSyncExternalStore(subscribe, readPainted, readPainted)
}

// Tests only: forget the cached choice so the next read takes storage again.
export function resetModeForTest() {
  mode = null
  preview = null
}
