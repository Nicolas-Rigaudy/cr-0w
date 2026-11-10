import { createContext, useContext, useEffect, useState } from 'react'
import type { DataSource } from './source'
import type { InboxItem, Mission } from './types'

export const DataContext = createContext<DataSource | null>(null)

export function useSource() {
  const src = useContext(DataContext)
  if (!src) throw new Error('DataContext missing')
  return src
}

export interface Snapshot {
  missions: Mission[]
  inbox: InboxItem[]
  now: Date
}

export function useSnapshot(): Snapshot | null {
  const src = useSource()
  const [snap, setSnap] = useState<Snapshot | null>(null)
  useEffect(() => {
    let live = true
    Promise.all([src.missions(), src.inbox()]).then(([missions, inbox]) => {
      if (live) setSnap({ missions, inbox, now: src.now() })
    })
    return () => { live = false }
  }, [src])
  return snap
}
