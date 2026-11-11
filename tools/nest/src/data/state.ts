import type { InboxItem, Mission } from './types'

// UI state language (item_138). Seven contract states fold into six UI states.
export type UiState = 'working' | 'needs' | 'failed' | 'sleeping' | 'queued' | 'done'

export const UI_STATES: UiState[] = ['needs', 'working', 'failed', 'sleeping', 'queued', 'done']

export const LABEL: Record<UiState, string> = {
  working: 'working',
  needs: 'needs you',
  failed: 'failed',
  sleeping: 'sleeping',
  queued: 'queued',
  done: 'done',
}

// Contract states core may move to cancelled (contract Moves, adr_001).
export const CANCELLABLE: Mission['state'][] = ['preflight', 'queued', 'running', 'waiting', 'review', 'failed']

export function lastRun(m: Mission) {
  return m.runs.length ? m.runs[m.runs.length - 1] : undefined
}

// sleeping = queued whose last run stopped on a usage limit: it waits for the account to reset.
// needs = the mission waits on the owner (waiting, review) or has an open inbox row.
export function uiState(m: Mission, inbox: InboxItem[] = []): UiState {
  if (m.state === 'failed') return 'failed'
  if (m.state === 'done' || m.state === 'cancelled') return 'done'
  if (m.state === 'waiting' || m.state === 'review') return 'needs'
  if (inbox.some((r) => r.mission === m.id && r.state === 'open' && r.type !== 'fyi')) return 'needs'
  if (m.state === 'queued') return lastRun(m)?.cause === 'usage_limit' ? 'sleeping' : 'queued'
  return 'working'
}

export function countByState(missions: Mission[], inbox: InboxItem[]) {
  const counts: Record<UiState, number> = { working: 0, needs: 0, failed: 0, sleeping: 0, queued: 0, done: 0 }
  for (const m of missions) counts[uiState(m, inbox)]++
  return counts
}

export function openInbox(inbox: InboxItem[]) {
  return inbox.filter((r) => r.state === 'open')
}

// The running mission a live run belongs to. Stopping such a run goes through the mission
// (crow mission stop) so core logs the owner's action event (item_168 AC4).
export function missionOfRun(missions: Mission[], runId: string): Mission | undefined {
  return missions.find((m) => m.state === 'running' && m.runs.some((r) => r.run === runId && !r.ended))
}
