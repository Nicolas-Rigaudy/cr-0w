import { describe, expect, it } from 'vitest'
import { blocksOf, reportsOf, running, spansOf, tilesOf, weekStart } from './reports'

const NOW = new Date('2026-11-09T16:10:00+01:00')
const rep = (kind: string, at: string, id = `${kind}:${at}`) => ({ id, kind, title: id, at, path: '/x', body: '' })

describe('reportsOf (item_258)', () => {
  it('reads items newest first, held and the retro routine from jobs', () => {
    const d = reportsOf(
      { data: { items: [rep('weekly', '2026-11-06T09:00:00+01:00'), rep('afk', '2026-11-09T07:55:00+01:00'), { kind: 'bogus' }], held: 2 }, stale: false, updated_at: '' },
      { data: [{ unit: 'crow-retro.timer' }], stale: false, updated_at: '' },
    )
    expect(d.items.map((r) => r.kind)).toEqual(['afk', 'weekly'])
    expect(d.held).toBe(2)
    expect(d.retro).toBe(true)
  })

  it('has no retro routine when jobs lacks the timer, and marks a failed refresh stale', () => {
    const d = reportsOf({ data: null, stale: true, error: 'boom', updated_at: '' }, { data: [{ unit: 'crow-brief.timer' }], stale: false, updated_at: '' })
    expect(d.retro).toBe(false)
    expect(d.stale).toBe(true)
    expect(d.items).toEqual([])
  })
})

describe('tilesOf', () => {
  it('counts this week from Monday and keeps the newest as last', () => {
    expect(weekStart(NOW).getDay()).toBe(1)
    const items = reportsOf({ data: { items: [rep('retro', '2026-11-09T09:00:00+01:00'), rep('retro', '2026-11-08T16:01:00+01:00')] }, stale: false, updated_at: '' }, undefined).items
    const t = tilesOf(items, NOW)
    expect(t.map((x) => x.kind)).toEqual(['afk', 'brief', 'daily', 'weekly', 'retro', 'eod'])
    const retro = t.find((x) => x.kind === 'retro')!
    expect(retro.week).toBe(1)
    expect(retro.last?.at).toBe('2026-11-09T09:00:00+01:00')
    expect(t.find((x) => x.kind === 'brief')!.last).toBeNull()
  })
})

describe('running', () => {
  it('runs until a newer retro lands or 30 minutes pass', () => {
    const started = '2026-11-09T15:00:00Z'
    const now = new Date('2026-11-09T15:10:00Z')
    expect(running(started, [], now)).toBe(true)
    const items = reportsOf({ data: { items: [rep('retro', '2026-11-09T15:05:00Z')] }, stale: false, updated_at: '' }, undefined).items
    expect(running(started, items, now)).toBe(false)
    expect(running(started, [], new Date('2026-11-09T15:40:00Z'))).toBe(false)
    expect(running(null, [], now)).toBe(false)
  })
})

describe('Markdown blocks', () => {
  it('reads headings, lists and paragraphs and drops front matter', () => {
    const b = blocksOf('---\ntags: x\n---\n# Title\n\nOne line\nsame para\n\n- a\n- [x] b\n  more\n\n### Sub\n1. n')
    expect(b).toEqual([
      { t: 'h', level: 1, text: 'Title' },
      { t: 'p', text: 'One line same para' },
      { t: 'ul', items: ['a', 'b more'] },
      { t: 'h', level: 3, text: 'Sub' },
      { t: 'ul', items: ['n'] },
    ])
  })

  it('reads bold, code, safe links and wiki links; other links stay text', () => {
    expect(spansOf('**a** `b` [c](https://x.y) [[Note|d]] [e](javascript:void)')).toEqual([
      { t: 'b', s: 'a' }, { t: 'text', s: ' ' }, { t: 'code', s: 'b' }, { t: 'text', s: ' ' },
      { t: 'a', s: 'c', href: 'https://x.y' }, { t: 'text', s: ' ' }, { t: 'text', s: 'd' }, { t: 'text', s: ' ' },
      { t: 'text', s: 'e' },
    ])
  })
})
