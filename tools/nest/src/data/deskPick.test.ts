import { describe, expect, it } from 'vitest'
import { deskPick, PICK_KEY, readPick, savePick } from './deskPick'
import type { Agent } from './live'
import type { WorkItem } from './work'

const memory = () => {
  const m = new Map<string, string>()
  return { getItem: (k: string) => m.get(k) ?? null, setItem: (k: string, v: string) => void m.set(k, v) }
}
const blocked = {
  getItem: (): string | null => { throw new Error('SecurityError') },
  setItem: () => { throw new Error('SecurityError') },
}

const agent = (pid: number, kind: Agent['kind'], state: string, started = '2026-11-25T08:00:00Z'): Agent =>
  ({ kind, pid, account: 'a', project: 'p', cwd: '/p', state, started_at: started })
const item = (kind: WorkItem['kind'], ref: string): WorkItem =>
  ({ id: `${kind}:${ref}`, kind, project: 'p', title: ref, state: 'working', source: { section: kind === 'run' ? 'dispatch' : 'sessions', ref } })

const run = agent(1, 'run', 'working')
const idle = agent(2, 'session', 'idle', '2026-11-25T09:00:00Z')
const working = agent(3, 'session', 'working')
const waiting = agent(4, 'session', 'waiting')

describe('desk pick (item_210)', () => {
  it('keeps the remembered pid while it is live, a run included', () => {
    expect(deskPick([run, working], null, 1)).toBe(run)
    expect(deskPick([run, working], [item('session', '3')], 3)).toBe(working)
  })

  it('falls back when the remembered pid ended', () => {
    expect(deskPick([run, working], null, 99)).toBe(working)
  })

  it('takes the first feed session that maps to a live agent', () => {
    const work = [item('run', 'r1'), item('session', '77'), item('session', '2'), item('session', '4')]
    expect(deskPick([run, idle, waiting], work, null)).toBe(idle)
  })

  it('without the feed: needs before working before idle, newest first', () => {
    expect(deskPick([run, idle, working, waiting], null, null)).toBe(waiting)
    expect(deskPick([run, idle, working], null, null)).toBe(working)
    const newer = agent(5, 'session', 'working', '2026-11-25T10:00:00Z')
    expect(deskPick([working, newer], null, null)).toBe(newer)
    expect(deskPick([run, idle], null, null)).toBe(idle)
    expect(deskPick([working, agent(6, 'session', 'setup')], null, null)?.pid).toBe(6)
  })

  it('never picks a background run by default', () => {
    expect(deskPick([run], null, null)).toBeUndefined()
    expect(deskPick([run], [item('run', 'r1')], null)).toBeUndefined()
    expect(deskPick([], null, null)).toBeUndefined()
  })

  it('a feed with no live session falls back to the state order', () => {
    expect(deskPick([run, idle, working], [item('session', '77')], null)).toBe(working)
  })
})

describe('desk pick storage (item_210)', () => {
  it('reads back the saved pid, as after a reload', () => {
    const s = memory()
    expect(readPick(s)).toBeNull()
    savePick(42, s)
    expect(s.getItem(PICK_KEY)).toBe('42')
    expect(readPick(s)).toBe(42)
  })

  it('ignores an invalid saved value', () => {
    const s = memory()
    s.setItem(PICK_KEY, 'nope')
    expect(readPick(s)).toBeNull()
    s.setItem(PICK_KEY, '0')
    expect(readPick(s)).toBeNull()
  })

  it('survives blocked or missing storage', () => {
    expect(readPick(blocked)).toBeNull()
    expect(() => savePick(1, blocked)).not.toThrow()
    expect(readPick(undefined)).toBeNull()
  })
})
