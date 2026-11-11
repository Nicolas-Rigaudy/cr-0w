import type { UiState } from '../data/state'
import { countByState, lastRun, openInbox, uiState } from '../data/state'
import type { InboxItem, Mission } from '../data/types'
import { title } from '../format'

// What the hero shows, derived from the data layer only (item_155, holotable item_162).

export type Mood = 'idle' | 'ask' | 'fail' | 'run'

const RANK: Record<UiState, number> = { needs: 0, failed: 1, working: 2, sleeping: 3, queued: 4, done: 5 }

export interface HeroMission { id: string; title: string; state: UiState; started?: string }
// open = missions not done (planet size). wait = minutes the oldest ask has waited (stalk height).
// runFrom = start of the oldest live run (working label). cause = what failed (failed label).
export interface HeroProject {
  name: string; state: UiState; live: boolean; missions: HeroMission[]
  open: number; wait?: number; runFrom?: string; cause?: string
}

export interface HeroModel {
  mood: Mood
  asks: number
  line: string
  projects: HeroProject[]
}

const WORDS = ['No', 'One', 'Two', 'Three', 'Four', 'Five', 'Six', 'Seven', 'Eight', 'Nine']
const word = (n: number) => WORDS[n] ?? String(n)

// Core mood, most urgent first (item_158, owner 2026-11-11): a failed run, then an open
// question, then live work, otherwise idle.
export function moodOf(c: Record<UiState, number>, asks = 0): Mood {
  if (c.failed > 0) return 'fail'
  if (asks > 0) return 'ask'
  if (c.working > 0) return 'run'
  return 'idle'
}

export function heroModel(missions: Mission[], inbox: InboxItem[], now: Date = new Date()): HeroModel {
  const c = countByState(missions, inbox)
  const asks = openInbox(inbox).filter((r) => r.type !== 'fyi').length
  const mood = moodOf(c, asks)

  const byProject = new Map<string, HeroMission[]>()
  for (const m of missions) {
    const run = lastRun(m)
    const hm: HeroMission = { id: m.id, title: title(m), state: uiState(m, inbox), started: run && !run.ended ? run.started : undefined }
    byProject.set(m.project, [...(byProject.get(m.project) ?? []), hm])
  }
  const asked = (p: string) => [
    ...openInbox(inbox).filter((r) => r.project === p && r.type !== 'fyi').map((r) => r.created),
    ...missions.filter((m) => m.project === p && (m.state === 'waiting' || m.state === 'review')).map((m) => m.updated),
  ]
  const mins = (iso: string) => Math.max(0, Math.round((now.getTime() - Date.parse(iso)) / 60000))
  const projects: HeroProject[] = [...byProject].map(([name, ms]) => {
    ms.sort((a, b) => RANK[a.state] - RANK[b.state])
    const state = ms[0].state
    const open = ms.filter((m) => m.state !== 'done')
    const since = asked(name).sort()[0]
    const runFrom = open.map((m) => m.started).filter((x): x is string => !!x).sort()[0]
    const failed = open.find((m) => m.state === 'failed')
    return {
      name, state, live: RANK[state] <= RANK.working, missions: open, open: open.length,
      wait: state === 'needs' && since ? mins(since) : undefined, runFrom, cause: failed?.title,
    }
  })
  // Longest wait first: the head looks at projects[0] when it waits on you.
  projects.sort((a, b) => RANK[a.state] - RANK[b.state] || (b.wait ?? 0) - (a.wait ?? 0) || a.name.localeCompare(b.name))

  const parts: string[] = []
  if (asks) parts.push(`${word(asks)} ${asks === 1 ? 'ask waits' : 'asks wait'} on you, Sir.`)
  const failed = missions.filter((m) => uiState(m, inbox) === 'failed')
  if (failed.length === 1) parts.push(`${title(failed[0])} failed.`)
  else if (failed.length) parts.push(`${word(failed.length)} runs failed.`)
  if (!parts.length && c.working) parts.push(`${word(c.working)} ${c.working === 1 ? 'run' : 'runs'} going, Sir.`)
  if (!parts.length) parts.push('All quiet, Sir.')

  return { mood, asks, line: parts.join(' '), projects }
}
