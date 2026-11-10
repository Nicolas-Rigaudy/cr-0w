import { createContext, useCallback, useContext, useEffect, useRef, useState } from 'react'
import { feed, type Snapshot } from './feed'
import type { DataSource } from './source'

export type { Snapshot } from './feed'

export const DataContext = createContext<DataSource | null>(null)

export function useSource() {
  const src = useContext(DataContext)
  if (!src) throw new Error('DataContext missing')
  return src
}

export const useSnapshot = (): Snapshot | null => useFeed()[0]

// Loads every section, then polls while the page is visible when the source asks
// for it. Hidden pages stop polling and refetch everything on return. The second
// value refetches at once (after a live action).
export function useFeed(): [Snapshot | null, () => Promise<void>] {
  const src = useSource()
  const [snap, setSnap] = useState<Snapshot | null>(null)
  const now = useRef<() => Promise<void>>(async () => {})
  useEffect(() => {
    let live = true
    const refresh = feed(src, (s) => { if (live) setSnap(s) })
    now.current = refresh
    void refresh()
    if (!src.poll) return () => { live = false }
    let timer: ReturnType<typeof setInterval> | undefined
    const start = () => { timer ??= setInterval(refresh, src.poll) }
    const stop = () => { clearInterval(timer); timer = undefined }
    const onVis = () => {
      if (document.visibilityState === 'hidden') return stop()
      void refresh()
      start()
    }
    if (document.visibilityState !== 'hidden') start()
    document.addEventListener('visibilitychange', onVis)
    return () => {
      live = false
      stop()
      document.removeEventListener('visibilitychange', onVis)
    }
  }, [src])
  return [snap, useCallback(() => now.current(), [])]
}
