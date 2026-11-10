import { createContext, useContext, useEffect, useState } from 'react'
import { feed, type Snapshot } from './feed'
import type { DataSource } from './source'

export type { Snapshot } from './feed'

export const DataContext = createContext<DataSource | null>(null)

export function useSource() {
  const src = useContext(DataContext)
  if (!src) throw new Error('DataContext missing')
  return src
}

// Loads every section, then polls while the page is visible when the source asks
// for it. Hidden pages stop polling and refetch everything on return.
export function useSnapshot(): Snapshot | null {
  const src = useSource()
  const [snap, setSnap] = useState<Snapshot | null>(null)
  useEffect(() => {
    let live = true
    const refresh = feed(src, (s) => { if (live) setSnap(s) })
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
  return snap
}
