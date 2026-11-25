import type { Agent, FinishedRun } from './live'
import type { Mission } from './types'
import type { WorkGroup, WorkItem } from './work'

// Desk layout C (item_295): four deck panes over the work feed.
// Needs you holds what waits on the owner, failures included. Moving holds what runs or
// sleeps on a usage limit. Queued holds what waits for a slot or a review. Done is history.
export type PaneId = 'needs' | 'moving' | 'queued' | 'done'
export const PANE_IDS: PaneId[] = ['needs', 'moving', 'queued', 'done']
export const PANE_TITLE: Record<PaneId, string> = { needs: 'Needs you', moving: 'Moving', queued: 'Queued', done: 'Done · failed' }
export type Panes = Record<PaneId, WorkItem[]>

const PANE_OF: Record<WorkItem['state'], PaneId> = { needs: 'needs', failed: 'needs', working: 'moving', sleeping: 'moving', queued: 'queued', done: 'done' }

// A queued mission the feed does not list yet (it only knows started work).
const missionUnit = (m: Mission): WorkItem => ({
  id: `mission:${m.id}`, kind: 'mission', project: m.project, title: m.goal || m.id, state: 'queued',
  since: m.updated, why: m.state === 'preflight' ? 'preflight' : 'waits for a slot', mission: m.id,
  source: { section: 'missions', ref: m.id },
})

// Without the feed (old core) the live agents stand in: a waiting session needs you, the rest moves.
export const agentUnit = (a: Agent): WorkItem => {
  const run = a.kind === 'run'
  const waiting = a.state === 'waiting' || a.state === 'setup'
  return {
    id: run ? `run:${a.run_id ?? a.pid}` : `session:${a.pid}`, kind: run ? 'run' : 'session', project: a.project,
    title: (run ? a.label || a.run_id : a.title) || a.project, state: waiting ? 'needs' : 'working',
    since: a.started_at ?? undefined, why: a.state === 'setup' ? 'waits on setup' : undefined,
    source: { section: 'dispatch', ref: run ? String(a.run_id ?? a.pid) : String(a.pid) },
  }
}

// A finished run the feed no longer lists: history for the Done · failed pane.
const finishedUnit = (f: FinishedRun): WorkItem => ({
  id: `run:${f.run_id}`, kind: 'run', project: f.project, title: f.label || f.run_id, state: f.status === 'done' ? 'done' : 'failed',
  since: f.ended_at ?? undefined, why: f.status === 'done' ? undefined : f.status, source: { section: 'dispatch', ref: f.run_id },
})

export function panesOf(work: WorkItem[] | null, agents: Agent[], queue: Mission[], finished: FinishedRun[] = [], project = ''): Panes {
  const units = work ?? agents.map(agentUnit)
  const known = new Set(units.map((w) => w.id))
  const p: Panes = { needs: [], moving: [], queued: [], done: [] }
  const mine = (w: WorkItem) => !project || w.project === project
  for (const w of units) if (mine(w)) p[PANE_OF[w.state]].push(w)
  for (const w of queue.map(missionUnit)) if (!known.has(w.id) && mine(w)) p.queued.push(w)
  for (const w of finished.map(finishedUnit)) if (!known.has(w.id) && mine(w)) p.done.push(w)
  return p
}

export const projectsOfPanes = (work: WorkItem[] | null, agents: Agent[]) =>
  [...new Set((work ?? agents.map(agentUnit)).map((w) => w.project))].sort()

// A pane's rows: a mission's units sit under its goal when the goal is in the same pane.
export function paneGroups(items: WorkItem[]): WorkGroup[] {
  const goals = new Set(items.filter((w) => w.kind === 'mission').map((w) => w.source.ref))
  const under = (w: WorkItem) => w.kind !== 'mission' && !!w.mission && goals.has(w.mission)
  return items.filter((w) => !under(w)).map((w) => ({
    item: w, children: w.kind === 'mission' ? items.filter((c) => under(c) && c.mission === w.source.ref) : [],
  }))
}

// Row order across panes for j / k: the pane order, each pane in its own order.
export const flatOf = (p: Panes, folded: boolean): WorkItem[] =>
  PANE_IDS.flatMap((id) => (id === 'queued' && folded ? [] : paneGroups(p[id]).flatMap((g) => [g.item, ...g.children])))

// The live agent behind a unit, if any.
export function agentOf(w: WorkItem, agents: Agent[]): Agent | undefined {
  if (w.kind === 'run') return agents.find((a) => a.kind === 'run' && (a.run_id === w.source.ref || String(a.pid) === w.source.ref))
  if (w.kind === 'session') return agents.find((a) => a.kind === 'session' && String(a.pid) === w.source.ref)
  return undefined
}

// The unit of a live agent (a ?pid= link, the remembered pick).
export const unitOfAgent = (p: Panes, a: Agent): WorkItem | undefined =>
  PANE_IDS.flatMap((id) => p[id]).find((w) => agentOf(w, [a]) === a)

// The Logics ref a title or label names (item_226, task_033, req_070).
export const itemRef = (s: string | null | undefined) => /\b(item|task|req)_\d+/.exec(s ?? '')?.[0]

export const summaryOf = (p: Panes) => ({ needs: p.needs.length, moving: p.moving.length, queued: p.queued.length, done: p.done.length })
