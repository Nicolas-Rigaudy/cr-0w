import type { Agent } from './live'
import type { WorkItem } from './work'

// The Desk pick is remembered per device in localStorage (item_210), like nest.scope.
export const PICK_KEY = 'nest.desk.pick'

type Store = Pick<Storage, 'getItem' | 'setItem'>
const local = (): Store | undefined => (typeof localStorage === 'undefined' ? undefined : localStorage)

// The saved pid; null when nothing valid is saved or storage is blocked.
export function readPick(store: Store | undefined = local()): number | null {
  try {
    const n = Number(store?.getItem(PICK_KEY))
    return Number.isInteger(n) && n > 0 ? n : null
  } catch {
    return null
  }
}

export function savePick(pid: number, store: Store | undefined = local()) {
  try { store?.setItem(PICK_KEY, String(pid)) } catch { /* storage blocked: the pick stays per visit */ }
}

// Session order without the feed: waiting (or setup), then working, then idle, then the rest.
const RANK: Record<string, number> = { waiting: 0, setup: 0, working: 1, idle: 2 }
const rank = (a: Agent) => RANK[a.state ?? ''] ?? 3
const start = (a: Agent) => (a.started_at ? Date.parse(a.started_at) || 0 : 0)

// The agent Desk opens: the remembered pid if live (a run included), else the first
// session of the work feed with a live agent, else the session that needs the owner
// most, newest first. Never a background run by default; undefined opens nothing.
export function deskPick(agents: Agent[], work: WorkItem[] | null, remembered: number | null): Agent | undefined {
  const kept = remembered == null ? undefined : agents.find((a) => a.pid === remembered)
  if (kept) return kept
  const sessions = agents.filter((a) => a.kind === 'session')
  if (work) {
    for (const w of work) {
      if (w.kind !== 'session') continue
      const a = sessions.find((s) => String(s.pid) === w.source.ref)
      if (a) return a
    }
  }
  return [...sessions].sort((x, y) => rank(x) - rank(y) || start(y) - start(x))[0]
}
