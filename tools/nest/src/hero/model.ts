import type { UiState } from '../data/state'
import { countByState, lastRun, openInbox, uiState } from '../data/state'
import type { InboxItem, Mission } from '../data/types'
import { title } from '../format'

// What the static hero shows, derived from the data layer only (item_155).

export type Mood = 'idle' | 'ask' | 'fail' | 'run'

const RANK: Record<UiState, number> = { needs: 0, failed: 1, working: 2, sleeping: 3, queued: 4, done: 5 }

export interface HeroMission { id: string; title: string; state: UiState; started?: string }
export interface HeroProject { name: string; state: UiState; live: boolean; missions: HeroMission[] }

export interface HeroModel {
  mood: Mood
  asks: number
  line: string
  projects: HeroProject[]
}

const WORDS = ['No', 'One', 'Two', 'Three', 'Four', 'Five', 'Six', 'Seven', 'Eight', 'Nine']
const word = (n: number) => WORDS[n] ?? String(n)

// The face follows the work, not the asks: the front planets and the chip carry the asks.
// Running work keeps the pass 7 grab pose; a failure braces only when nothing runs.
export function moodOf(c: Record<UiState, number>): Mood {
  if (c.working > 0) return 'run'
  if (c.failed > 0) return 'fail'
  return 'idle'
}

export function heroModel(missions: Mission[], inbox: InboxItem[]): HeroModel {
  const c = countByState(missions, inbox)
  const asks = openInbox(inbox).filter((r) => r.type !== 'fyi').length
  const mood = moodOf(c)

  const byProject = new Map<string, HeroMission[]>()
  for (const m of missions) {
    const run = lastRun(m)
    const hm: HeroMission = { id: m.id, title: title(m), state: uiState(m, inbox), started: run && !run.ended ? run.started : undefined }
    byProject.set(m.project, [...(byProject.get(m.project) ?? []), hm])
  }
  const projects = [...byProject].map(([name, ms]) => {
    ms.sort((a, b) => RANK[a.state] - RANK[b.state])
    const state = ms[0].state
    return { name, state, live: RANK[state] <= RANK.working, missions: ms.filter((m) => m.state !== 'done') }
  })
  projects.sort((a, b) => RANK[a.state] - RANK[b.state] || a.name.localeCompare(b.name))

  const parts: string[] = []
  if (asks) parts.push(`${word(asks)} ${asks === 1 ? 'ask waits' : 'asks wait'} on you, Sir.`)
  const failed = missions.filter((m) => uiState(m, inbox) === 'failed')
  if (failed.length === 1) parts.push(`${title(failed[0])} failed.`)
  else if (failed.length) parts.push(`${word(failed.length)} runs failed.`)
  if (!parts.length && c.working) parts.push(`${word(c.working)} ${c.working === 1 ? 'run' : 'runs'} going, Sir.`)
  if (!parts.length) parts.push('All quiet, Sir.')

  return { mood, asks, line: parts.join(' '), projects }
}
