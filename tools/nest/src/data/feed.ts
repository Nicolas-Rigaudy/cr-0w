import type { DataSource } from './source'
import type { AgendaEvent, InboxItem, Mission } from './types'

export interface Snapshot {
  missions: Mission[]
  inbox: InboxItem[]
  agenda: AgendaEvent[]
  now: Date
  // Last good refresh. down is true while the latest refresh failed; the data
  // then stays as it was at synced (null when nothing ever loaded).
  synced: Date | null
  down: boolean
}

const minute = (d: Date) => Math.floor(d.getTime() / 60000)

// A feed refreshes every section of a source and emits a new snapshot only when
// a section body, the link state or the clock minute changed.
export function feed(src: DataSource, emit: (s: Snapshot) => void) {
  let snap: Snapshot | null = null
  let ok: Date | null = null
  let busy = false
  return async function refresh() {
    if (busy) return
    busy = true
    try {
      const now = src.now()
      const same = snap !== null && minute(snap.now) === minute(now)
      try {
        const [missions, inbox, agenda] = await Promise.all([src.missions(), src.inbox(), src.agenda()])
        ok = now
        if (same && !snap!.down && snap!.missions === missions && snap!.inbox === inbox && snap!.agenda === agenda) return
        snap = { missions, inbox, agenda, now, synced: now, down: false }
      } catch {
        if (same && snap!.down) return
        snap = { ...(snap ?? { missions: [], inbox: [], agenda: [] }), now, synced: ok, down: true }
      }
      emit(snap)
    } finally {
      busy = false
    }
  }
}
