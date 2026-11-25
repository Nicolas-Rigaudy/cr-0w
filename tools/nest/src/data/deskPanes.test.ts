import { describe, expect, it } from 'vitest'
import type { Agent, FinishedRun } from './live'
import type { Mission } from './types'
import type { WorkItem } from './work'
import { agentOf, flatOf, itemRef, paneGroups, panesOf, unitOfAgent } from './deskPanes'

const w = (id: string, state: WorkItem['state'], extra: Partial<WorkItem> = {}): WorkItem => {
  const [kind, ref] = id.split(':') as [WorkItem['kind'], string]
  return { id, kind, project: 'p', title: id, state, source: { section: 'dispatch', ref }, ...extra }
}
const agent = (pid: number, kind: Agent['kind'], state: string, extra: Partial<Agent> = {}): Agent =>
  ({ kind, pid, account: 'main', project: 'p', cwd: '/', state, ...extra })
const mission = (id: string, state: Mission['state']) => ({ id, project: 'p', goal: `goal ${id}`, state, updated: '2026-11-09T10:00:00Z' }) as Mission
const fin = (run_id: string, status: string): FinishedRun => ({ run_id, project: 'p', label: run_id, status })

describe('desk panes (item_295)', () => {
  it('sorts the feed into the four panes, failures with what needs you', () => {
    const p = panesOf([w('session:1', 'needs'), w('run:a', 'failed'), w('run:b', 'working'), w('parked:c', 'sleeping'), w('mr:x', 'queued'), w('mission:m', 'done')], [], [])
    expect(p.needs.map((x) => x.id)).toEqual(['session:1', 'run:a'])
    expect(p.moving.map((x) => x.id)).toEqual(['run:b', 'parked:c'])
    expect(p.queued.map((x) => x.id)).toEqual(['mr:x'])
    expect(p.done.map((x) => x.id)).toEqual(['mission:m'])
  })

  it('adds queued missions and finished runs the feed does not list', () => {
    const p = panesOf([w('mission:m1', 'working')], [], [mission('m1', 'queued'), mission('m2', 'preflight')], [fin('r1', 'failed'), fin('r2', 'done')])
    expect(p.queued.map((x) => x.id)).toEqual(['mission:m2'])
    expect(p.queued[0].why).toBe('preflight')
    expect(p.done.map((x) => [x.id, x.state])).toEqual([['run:r1', 'failed'], ['run:r2', 'done']])
  })

  it('falls back to the live agents without the feed', () => {
    const p = panesOf(null, [agent(1, 'session', 'waiting'), agent(2, 'session', 'setup'), agent(3, 'run', 'working', { run_id: 'r3' })], [])
    expect(p.needs.map((x) => x.id)).toEqual(['session:1', 'session:2'])
    expect(p.needs[1].why).toBe('waits on setup')
    expect(p.moving.map((x) => x.id)).toEqual(['run:r3'])
  })

  it('filters on a project', () => {
    const p = panesOf([w('run:a', 'working'), w('run:b', 'working', { project: 'q' })], [], [], [], 'q')
    expect(p.moving.map((x) => x.id)).toEqual(['run:b'])
  })

  it('nests mission units under their goal and skips a folded queue', () => {
    const items = [w('mission:m', 'working'), w('run:a', 'working', { mission: 'm' }), w('run:b', 'working')]
    expect(paneGroups(items).map((g) => [g.item.id, g.children.map((c) => c.id)])).toEqual([['mission:m', ['run:a']], ['run:b', []]])
    const p = panesOf([...items, w('mr:x', 'queued')], [], [])
    expect(flatOf(p, true).map((x) => x.id)).toEqual(['mission:m', 'run:a', 'run:b'])
    expect(flatOf(p, false).map((x) => x.id)).toContain('mr:x')
  })

  it('links units and live agents both ways', () => {
    const run = agent(7, 'run', 'working', { run_id: 'a1' })
    const ses = agent(8, 'session', 'waiting')
    const p = panesOf([w('run:a1', 'working'), w('session:8', 'needs')], [run, ses], [])
    expect(agentOf(p.moving[0], [run, ses])).toBe(run)
    expect(unitOfAgent(p, ses)?.id).toBe('session:8')
    expect(itemRef('Implement item_295 end to end')).toBe('item_295')
    expect(itemRef('no ref')).toBeUndefined()
  })
})
