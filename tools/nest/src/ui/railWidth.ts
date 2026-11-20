import { useEffect, useRef, useState, type KeyboardEvent, type PointerEvent } from 'react'

// Butler rail width on laptop (item_256), remembered per device like nest.rail.
export const RAIL_W_KEY = 'nest.railw'
export const RAIL_W = 320
// Home keeps its own width (item_259): the chat sits wider next to the hero.
export const HOME_W_KEY = 'nest.railw.home'
export const HOME_W = 480
export const RAIL_MIN = 280
export const RAIL_MAX_RATIO = 0.6
export const RAIL_STEP = 24

type Store = Pick<Storage, 'getItem' | 'setItem'>
const local = (): Store | undefined => (typeof localStorage === 'undefined' ? undefined : localStorage)
const winW = () => (typeof window === 'undefined' ? Infinity : window.innerWidth)

// Width kept between 280px and 60% of the window.
export function clampRail(w: number, win: number): number {
  const max = Math.max(RAIL_MIN, Math.floor(win * RAIL_MAX_RATIO))
  return Math.round(Math.min(Math.max(w, RAIL_MIN), max))
}

// The saved width; the default when nothing valid is saved or storage is blocked.
export function readRailWidth(store: Store | undefined = local(), key = RAIL_W_KEY, def = RAIL_W): number {
  try {
    const v = Number(store?.getItem(key))
    return Number.isFinite(v) && v >= RAIL_MIN ? v : def
  } catch {
    return def
  }
}

export function saveRailWidth(w: number, store: Store | undefined = local(), key = RAIL_W_KEY) {
  try { store?.setItem(key, String(Math.round(w))) } catch { /* storage blocked: the width stays per visit */ }
}

// Set on the root while the grip is held: the hero waits for the release to refit (item_259).
const dragging = (on: boolean) => {
  if (on) document.documentElement.dataset.railDrag = ''
  else delete document.documentElement.dataset.railDrag
}

// Grip on the rail's left edge: drag, Left and Right by 24px, double-click back to the default.
// The width lands in --rail-w on the root, read by the rail and the app grid.
export function useRailWidth(key = RAIL_W_KEY, def = RAIL_W) {
  const [w, setW] = useState(() => readRailWidth(local(), key, def))
  const drag = useRef<{ x0: number; w0: number } | null>(null)
  useEffect(() => { document.documentElement.style.setProperty('--rail-w', `${w}px`) }, [w])
  const set = (next: number, save = true) => {
    const c = clampRail(next, winW())
    setW(c)
    if (save) saveRailWidth(c, local(), key)
  }
  const grip = {
    role: 'separator',
    'aria-orientation': 'vertical' as const,
    'aria-label': 'Resize CR-0W',
    'aria-valuenow': clampRail(w, winW()),
    'aria-valuemin': RAIL_MIN,
    'aria-valuemax': clampRail(Infinity, winW()),
    tabIndex: 0,
    title: `Drag to resize, double-click for ${def}px`,
    onPointerDown: (e: PointerEvent<HTMLElement>) => {
      e.preventDefault()
      e.currentTarget.setPointerCapture(e.pointerId)
      drag.current = { x0: e.clientX, w0: clampRail(w, winW()) }
      dragging(true)
    },
    onPointerMove: (e: PointerEvent<HTMLElement>) => {
      const d = drag.current
      if (d) set(d.w0 + d.x0 - e.clientX, false)
    },
    onPointerUp: (e: PointerEvent<HTMLElement>) => {
      const d = drag.current
      drag.current = null
      dragging(false)
      if (d) set(d.w0 + d.x0 - e.clientX)
    },
    onPointerCancel: () => { drag.current = null; dragging(false) },
    onDoubleClick: () => set(def),
    onKeyDown: (e: KeyboardEvent<HTMLElement>) => {
      const dir = e.key === 'ArrowLeft' ? 1 : e.key === 'ArrowRight' ? -1 : 0
      if (!dir) return
      e.preventDefault()
      set(clampRail(w, winW()) + dir * RAIL_STEP)
    },
  }
  return grip
}
