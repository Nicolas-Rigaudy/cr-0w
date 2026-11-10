import type { UiState } from '../data/state'
import { uiState } from '../data/state'
import type { InboxItem, Mission } from '../data/types'

const RANK: Record<UiState, number> = { needs: 0, failed: 1, working: 2, sleeping: 3, queued: 4, done: 5 }

// Owner first: what needs you, then what broke, then what runs.
export function byUrgency(missions: Mission[], inbox: InboxItem[]) {
  return [...missions].sort((a, b) => RANK[uiState(a, inbox)] - RANK[uiState(b, inbox)] || b.updated.localeCompare(a.updated))
}
