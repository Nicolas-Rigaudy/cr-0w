import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { doneToday, kindOf, kindWord, type AnswerRow } from '../data/inbox'
import type { InboxItem } from '../data/types'
import { DataContext } from '../data/useData'
import { Desk } from './Desk'
import { Inbox } from './Inbox'

const snapshot = async (): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(fixtureSource(), (x) => { box.s = x })()
  return box.s!
}
// Width of the screen under test: phone below 600px, fold up to 1100px, laptop above.
const at = (w: number) => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search: '', pathname: '/' })
}
const render = (el: React.ReactElement) => renderToStaticMarkup(<DataContext.Provider value={fixtureSource()}>{el}</DataContext.Provider>)
const refresh = async () => {}

afterEach(() => { vi.unstubAllGlobals() })

describe('Inbox to answer (item_204)', () => {
  it('lists questions then mail in the board row, no failed work (owner 2026-11-16)', async () => {
    at(707)
    const html = render(<Inbox snap={await snapshot()} refresh={refresh} />)
    const list = html.slice(html.indexOf('class="alist"'), html.indexOf('Waiting on others'))
    const titles = [...list.matchAll(/class="irow[^"]*".*?<b>(.*?)<\/b>/g)].map((m) => m[1])
    expect(html).toContain('<span>Done today</span>')
    expect(titles.length).toBeGreaterThan(2)
    expect(titles).not.toContain('night job')
    expect(titles.some((t) => t.startsWith('!230'))).toBe(false)
    expect(titles.slice(-2)).toEqual(['Invoice to check', 'Contract renewal'])
    expect(html).toContain('<span>To answer</span><span class="num">' + titles.length + '</span>')
  })

  it('draws each row on two lines: kind icon, title, kind · source · status (board 223-B)', async () => {
    at(707)
    const html = render(<Inbox snap={await snapshot()} refresh={refresh} />)
    const list = html.slice(html.indexOf('class="alist"'), html.indexOf('Waiting on others'))
    const rows = list.split('<button').slice(1)
    expect(rows.length).toBeGreaterThan(0)
    for (const r of rows) {
      expect(r).toMatch(/<span class="ic k-[a-z]+"><svg/)
      expect(r).toMatch(/<span class="t"><b>[^<]+<\/b><span>/)
    }
    expect(list).toContain('<b>fix timeout</b><span>ask · session · lambda · asks Bash</span>')
    expect(list).toContain('session · website · asks AskUserQuestion')
    expect(list).toMatch(/<span>(decision|approval|doc review|design review) · /)
    expect(list).toContain('mail · Accounting · flagged')
    expect(list).toContain('Legal · important')
  })

  it('adds the why to a waiting run, only when it has one (owner 2026-11-16)', async () => {
    at(707)
    const s = await snapshot()
    const run = (id: string, title: string, why?: string) =>
      ({ id, kind: 'run', project: 'cr-0w', state: 'needs', title, since: '2026-11-09T09:00:00+01:00', why, source: { section: 'dispatch', ref: id } })
    const data = [...(s.sections.work!.data as unknown[]), run('run:r1', 'nightly', 'asks Edit'), run('run:r2', 'backup')]
    const html = render(<Inbox snap={{ ...s, sections: { ...s.sections, work: { ...s.sections.work!, data } } }} refresh={refresh} />)
    expect(html).toContain('<b>nightly</b><span>ask · run · cr-0w · asks Edit</span>')
    expect(html).toContain('<b>backup</b><span>ask · run · cr-0w</span>')
  })
})

describe('Inbox layout B (item_223, item_191)', () => {
  it('shows the kind rail and the Reports log on the laptop, three columns', async () => {
    at(1440)
    const html = render(<Inbox snap={await snapshot()} refresh={refresh} />)
    expect(html).toContain('class="desk ibx ib-lp"')
    expect(html).toContain('aria-label="Kinds"')
    for (const k of ['All', 'Decisions', 'Reviews', 'Approvals', 'Asks', 'Mail', 'Follow-ups']) expect(html).toContain(`<span>${k}</span>`)
    expect(html).toMatch(/href="\/reports"/)
    expect(html).toContain('read-only, ')
    expect(html).toContain('Waiting on others')
    expect(html).toContain('id="rt"')
    expect(html).toContain('→ ')
  })

  it('shows kind chips and no open sheet on the phone, never a failed item', async () => {
    at(368)
    const html = render(<Inbox snap={await snapshot()} refresh={refresh} />)
    expect(html).toContain('class="desk ibx ib-ph"')
    expect(html).toContain('class="ichip on"')
    expect(html).toContain('Mail 2')
    expect(html).not.toContain('role="dialog"')
    expect(html).not.toContain('night job')
  })

  it('opens the reading pane on the first row with its verbs and Hide on the fold', async () => {
    at(707)
    const html = render(<Inbox snap={await snapshot()} refresh={refresh} />)
    expect(html).toContain('class="desk ibx ib-in"')
    expect(html).toContain('placeholder="Change or note (optional)"')
    expect(html).toContain('>Hide</button>')
  })
})

describe('Inbox kinds and Done today', () => {
  const item = (over: Partial<InboxItem>): InboxItem => ({
    id: 'm/v1', mission: 'm', project: 'p', version: 1, file: 'f.md', ask: 'a', type: 'choice', source: 'questions',
    verbs: ['pick'], artifacts: [], state: 'open', created: '2026-11-25T08:00:00+01:00', answered: null, ...over,
  } as InboxItem)
  const ask = (r: InboxItem): AnswerRow => ({ key: r.id, kind: 'ask', at: r.created, item: r })

  it('maps each row to one kind', () => {
    expect(kindOf(ask(item({ type: 'choice' })))).toBe('decision')
    expect(kindOf(ask(item({ type: 'approval', source: 'plan' })))).toBe('approval')
    expect(kindWord(ask(item({ type: 'review', source: 'report' })))).toBe('doc review')
    expect(kindWord(ask(item({ type: 'choice', source: 'mockup', artifacts: ['a.png', 'b.png'] })))).toBe('design review')
    expect(kindWord(ask(item({ type: 'choice', artifacts: ['x-A.png', 'x-B.png'] })))).toBe('design review')
    expect(kindOf(ask(item({ type: 'needs-input' })))).toBe('ask')
  })

  it('keeps rows answered or hidden today, newest first', () => {
    const now = new Date('2026-11-25T18:00:00+01:00')
    const rows = [
      item({ id: 'a', state: 'answered', answered: '2026-11-25T09:00:00+01:00' }),
      item({ id: 'b', state: 'dropped', answered: '2026-11-25T11:00:00+01:00' }),
      item({ id: 'c', state: 'answered', answered: '2026-11-24T11:00:00+01:00' }),
      item({ id: 'd', state: 'open' }),
      item({ id: 'e', type: 'fyi', state: 'answered', answered: '2026-11-25T10:00:00+01:00' }),
    ]
    expect(doneToday(rows, now).map((r) => r.id)).toEqual(['b', 'a'])
  })
})

describe('Desk without flagged mail (item_204)', () => {
  it('keeps the failed run that left Inbox (owner 2026-11-16)', async () => {
    at(707)
    expect(render(<Desk snap={await snapshot()} refresh={refresh} />)).toContain('night job')
  })

  it('shows no mail group at any width', async () => {
    for (const w of [368, 707, 1440]) {
      at(w)
      const html = render(<Desk snap={await snapshot()} refresh={refresh} />)
      expect(html).not.toContain('Flagged mail')
      expect(html).not.toContain('Invoice to check')
    }
  })
})
