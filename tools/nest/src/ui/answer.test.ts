import { describe, expect, it } from 'vitest'
import type { InboxItem } from '../data/types'
import { current } from './answer'

const row = (id: string, version: number, state: InboxItem['state'] = 'open', type: InboxItem['type'] = 'choice'): InboxItem => ({
  id, mission: 'm1', kind: 'build', project: 'p', type, source: 'questions', version, file: '', ask: id,
  verbs: ['answer'], blocking: true, artifacts: [], state, created: '', answered: null,
})

describe('current', () => {
  it('returns the asked row while it is open', () => {
    const asked = row('m1/v1/q1', 1)
    expect(current([asked], asked)).toBe(asked)
  })

  it('moves to the same question on the newest version', () => {
    const asked = row('m1/v1/q1', 1)
    const v2 = row('m1/v2/q1', 2)
    const v3 = row('m1/v3/q1', 3)
    const other = row('m1/v3/q2', 3)
    expect(current([row('m1/v1/q1', 1, 'answered'), v2, other, v3], asked)).toBe(v3)
  })

  it('matches the version row by its empty question and skips fyi rows', () => {
    const asked = row('m1/v1', 1)
    expect(current([row('m1/v2', 2, 'open', 'fyi')], asked)).toBeUndefined()
    const v2 = row('m1/v2', 2)
    expect(current([v2, row('m1/v2/q1', 2)], asked)).toBe(v2)
  })

  it('returns nothing once the question is closed everywhere', () => {
    const asked = row('m1/v1/q1', 1)
    expect(current([row('m1/v1/q1', 1, 'answered'), row('m1/v2/q1', 2, 'dropped')], asked)).toBeUndefined()
  })
})
