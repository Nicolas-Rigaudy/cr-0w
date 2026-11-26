import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { isMode, paintMode, readKeptMode, resetModeForTest, resolveMode, setMode, unpaintMode, watchSystemMode } from './theme'

// A fake page: <html> dataset, theme-color meta, storage and a prefers-color-scheme query the test flips.
function page(stored: string | null, dark: boolean) {
  const root = { dataset: {} as Record<string, string> }
  const meta = { content: '', setAttribute(_: string, v: string) { this.content = v } }
  const store = new Map<string, string>(stored === null ? [] : [['nest.mode', stored]])
  const listeners = new Set<() => void>()
  const mq = {
    get matches() { return dark },
    addEventListener: (_: string, f: () => void) => { listeners.add(f) },
    removeEventListener: (_: string, f: () => void) => { listeners.delete(f) },
  }
  vi.stubGlobal('document', { documentElement: root, querySelector: () => meta })
  vi.stubGlobal('getComputedStyle', () => ({ getPropertyValue: () => (root.dataset.mode === 'light' ? '#f4f4f0' : '#05070c') }))
  vi.stubGlobal('localStorage', { getItem: (k: string) => store.get(k) ?? null, setItem: (k: string, v: string) => { store.set(k, v) } })
  vi.stubGlobal('matchMedia', () => mq)
  return {
    root, meta, store, listeners,
    flip(d: boolean) { dark = d; for (const f of listeners) f() },
  }
}

beforeEach(() => resetModeForTest())
afterEach(() => { vi.unstubAllGlobals() })

describe('resolveMode', () => {
  it('auto follows the system, light and dark force it', () => {
    expect(resolveMode('auto', true)).toBe('dark')
    expect(resolveMode('auto', false)).toBe('light')
    expect(resolveMode('light', true)).toBe('light')
    expect(resolveMode('dark', false)).toBe('dark')
  })

  it('knows the three choices only', () => {
    expect(['auto', 'light', 'dark'].every(isMode)).toBe(true)
    expect(isMode('dim')).toBe(false)
    expect(isMode(null)).toBe(false)
  })
})

describe('stored choice', () => {
  it('defaults to auto when nothing or junk is stored', () => {
    page(null, true)
    expect(readKeptMode()).toBe('auto')
    resetModeForTest()
    page('sepia', true)
    expect(readKeptMode()).toBe('auto')
  })

  it('reads the stored mode and setMode keeps and paints it with the bar colour', () => {
    const p = page('dark', false)
    expect(readKeptMode()).toBe('dark')
    setMode('light')
    expect(p.store.get('nest.mode')).toBe('light')
    expect(p.root.dataset.mode).toBe('light')
    expect(p.meta.content).toBe('#f4f4f0')
  })

  it('a preview paints without keeping, unpaint brings the kept mode back', () => {
    const p = page('dark', true)
    paintMode('light')
    expect(p.root.dataset.mode).toBe('light')
    expect(p.store.get('nest.mode')).toBe('dark')
    expect(readKeptMode()).toBe('dark')
    unpaintMode()
    expect(p.root.dataset.mode).toBe('dark')
  })
})

describe('system listener', () => {
  it('repaints live while auto', () => {
    const p = page(null, true)
    const stop = watchSystemMode()
    p.flip(false)
    expect(p.root.dataset.mode).toBe('light')
    p.flip(true)
    expect(p.root.dataset.mode).toBe('dark')
    stop()
    expect(p.listeners.size).toBe(0)
  })

  it('leaves a forced mode alone', () => {
    const p = page('dark', true)
    setMode('dark')
    const stop = watchSystemMode()
    p.flip(false)
    expect(p.root.dataset.mode).toBe('dark')
    stop()
  })

  it('follows the system during an auto preview over a forced mode', () => {
    const p = page('light', true)
    paintMode('auto')
    const stop = watchSystemMode()
    p.flip(false)
    expect(p.root.dataset.mode).toBe('light')
    p.flip(true)
    expect(p.root.dataset.mode).toBe('dark')
    stop()
  })
})
