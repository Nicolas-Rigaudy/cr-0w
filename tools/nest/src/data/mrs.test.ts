import { describe, expect, it } from 'vitest'
import stateMrs from '../../fixtures/states/mrs.json'
import { countsOf, hitsOf, mrOf, mrsOf, needsYou, onlyProject, patternsFor } from './mrs'
import { scopeSnapshot } from './scope'
import type { Snapshot } from './feed'
import type { SectionBody } from './types'

const NOW = new Date('2026-11-09T16:10:00+01:00')
const body = (data: unknown, more: Partial<SectionBody> = {}): SectionBody => ({ data, stale: false, updated_at: null, ...more })

describe('MRs mapping (item_225)', () => {
  it('splits crow-mrs rows into to review, mine and merged today', () => {
    const d = mrsOf(body(stateMrs), NOW)
    expect(d.review.map((m) => m.ref)).toEqual(['!242', '!245'])
    expect(d.mine.map((m) => m.ref)).toEqual(['!239', '!241', '#57'])
    // !229 merged two days ago stays out of Merged today.
    expect(d.merged.map((m) => m.ref)).toEqual(['!235'])
    expect(d.branches).toEqual([expect.objectContaining({ branch: 'proxy-rate-limit', ahead: 2 })])
    expect(d.projects).toEqual(['api', 'cr-0w'])
  })

  it('puts drafts last in To review and what needs the owner first in Mine', () => {
    const rows = [
      { project: 'p', number: 1, role: 'REVIEWER', draft: true, updated: '2026-11-01T00:00:00Z' },
      { project: 'p', number: 2, role: 'ASSIGNEE', updated: '2026-11-08T00:00:00Z' },
      { project: 'p', number: 3, role: 'AUTHOR', pipeline: 'passed', updated: '2026-11-01T00:00:00Z' },
      { project: 'p', number: 4, role: 'AUTHOR', threads_open: 1, updated: '2026-11-08T00:00:00Z' },
    ]
    const d = mrsOf(body({ mrs: rows }), NOW)
    expect(d.review.map((m) => m.number)).toEqual([2, 1])
    expect(d.mine.map((m) => m.number)).toEqual([4, 3])
    expect(needsYou(d.mine[0])).toBe(true)
  })

  it('counts the header summary from the rows', () => {
    expect(countsOf(mrsOf(body(stateMrs), NOW))).toEqual({ review: 1, needs: 1, waiting: 2, merged: 1 })
  })

  it('fills the schema defaults and builds the ref from the forge', () => {
    const m = mrOf({ project: 'x', forge: 'github', number: 7, role: 'reviewer', pipeline: 'weird' })
    expect(m).toMatchObject({ ref: '#7', role: 'REVIEWER', pipeline: 'none', draft: false, threads_open: 0, name: 'x' })
    expect(mrOf({ forge: 'gitlab', number: 3 }).ref).toBe('!3')
  })

  it('reads off and stale from the section', () => {
    expect(mrsOf(body(null, { off: true }), NOW).off).toBe(true)
    expect(mrsOf(body({ off: true, mrs: [] }), NOW).off).toBe(true)
    const s = mrsOf(body(stateMrs, { stale: true, error: 'crow-mrs: timed out' }), NOW)
    expect(s.stale).toBe(true)
    expect(s.error).toBe('crow-mrs: timed out')
    expect(mrsOf(undefined, NOW)).toMatchObject({ off: false, review: [], mine: [], merged: [] })
  })

  it('filters every list by project', () => {
    const d = onlyProject(mrsOf(body(stateMrs), NOW), 'cr-0w')
    expect([...d.review, ...d.mine, ...d.merged].map((m) => m.ref)).toEqual(['#57'])
    expect(d.branches).toEqual([])
  })

  it('lists the project and common ledger patterns, newest hit first', () => {
    const d = mrsOf(body(stateMrs), NOW)
    const m = d.mine[0]
    const ps = patternsFor(d, m)
    expect(ps.map((p) => p.count)).toEqual([1, 7, 9, 10, 5, 4])
    expect(hitsOf(ps, m).map((p) => p.pattern)).toEqual(['Same loop copied instead of one helper', "Name doesn't follow its siblings"])
  })

  it('keeps the MRs of the scope the header shows', () => {
    const s = { sections: { mrs: body(stateMrs), projects: body({ projects: [{ name: 'cr-0w', scope: 'personal' }] }), dispatch: body({}), sessions: body({}), today: body({}), work: body([]), ready: body({}) }, missions: [], inbox: [] } as unknown as Snapshot
    const perso = mrsOf(scopeSnapshot(s, 'perso').sections.mrs, NOW)
    expect([...perso.review, ...perso.mine, ...perso.merged].map((m) => m.ref)).toEqual(['#57'])
    const work = mrsOf(scopeSnapshot(s, 'work').sections.mrs, NOW)
    expect(work.mine.map((m) => m.ref)).toEqual(['!239', '!241'])
    expect(work.branches).toHaveLength(1)
  })
})
