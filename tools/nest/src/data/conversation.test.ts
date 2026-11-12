import { describe, expect, it } from 'vitest'
import { match } from '../router'
import { blocksOf, conversationPath, pageQuery, partsOf, spansOf, type Entry } from './conversation'
import { CHAT, fixtureSource } from './fixtures'

const e = (id: string, role: Entry['role'], text = '', tool_name?: string, tool_summary?: string): Entry =>
  ({ at: '2026-11-09T16:00:00+01:00', id, role, text, tool_name, tool_summary })

describe('blocksOf', () => {
  it('folds the board chat: two folds with call counts and tool names in first-use order', () => {
    const b = blocksOf(CHAT)
    expect(b.map((x) => x.kind)).toEqual(['me', 'tools', 'cr', 'me', 'tools', 'cr'])
    const folds = b.filter((x) => x.kind === 'tools')
    expect(folds.map((f) => f.kind === 'tools' && [f.calls, f.tools])).toEqual([[4, ['Read', 'Bash']], [9, ['Bash', 'Read', 'Edit']]])
  })

  it('shows the who label only when the side changes, folds aside', () => {
    const b = blocksOf([e('1', 'assistant', 'a'), e('2', 'tool', '', 'Read', 'Call recorded'), e('3', 'assistant', 'b'), e('4', 'user', 'c')])
    expect(b.flatMap((x) => (x.kind === 'tools' ? [] : [x.who]))).toEqual([true, false, true])
  })

  it('skips empty messages and counts a results-only fold as one call', () => {
    const b = blocksOf([e('1', 'tool', '', 'Bash', 'Result recorded'), e('2', 'assistant', '  ')])
    expect(b).toEqual([{ kind: 'tools', id: '1', calls: 1, tools: ['Bash'], entries: [expect.objectContaining({ id: '1' })] }])
  })
})

describe('partsOf / spansOf', () => {
  it('splits paragraphs, list lines, headings and code fences', () => {
    expect(partsOf('# Plan\nfirst line\nsame para\n\n- a\n* b\n2) c\n```\nx = 1\n```\nend')).toEqual([
      { t: 'h', text: 'Plan' },
      { t: 'p', text: 'first line\nsame para' },
      { t: 'li', mark: '•', text: 'a' },
      { t: 'li', mark: '•', text: 'b' },
      { t: 'li', mark: '2)', text: 'c' },
      { t: 'code', text: 'x = 1' },
      { t: 'p', text: 'end' },
    ])
  })

  it('marks bold and inline code, keeps the rest', () => {
    expect(spansOf('run **now** with `crow x` ok')).toEqual([
      { t: 'txt', text: 'run ' }, { t: 'b', text: 'now' }, { t: 'txt', text: ' with ' }, { t: 'code', text: 'crow x' }, { t: 'txt', text: ' ok' },
    ])
  })
})

describe('paths and routes', () => {
  it('builds the core path and the page query, cursor first', () => {
    expect(conversationPath({ pid: 42 })).toBe('/sessions/42/conversation')
    expect(conversationPath({ mission: 'a b' })).toBe('/missions/a%20b/conversation')
    expect(pageQuery({ tail: 50, cursor: 'c1', before: 'b1' })).toBe('?cursor=c1')
    expect(pageQuery({ tail: 50, before: 'b 1' })).toBe('?before=b%201')
    expect(pageQuery({ tail: 50 })).toBe('?tail=50')
    expect(pageQuery({})).toBe('')
  })

  it('routes /desk/<pid> and /desk/m/<mission> to the chat', () => {
    expect(match('/desk/4242')).toEqual({ tab: 'chat', target: { pid: 4242 } })
    expect(match('/desk/m/20261109-nest%20x')).toEqual({ tab: 'chat', target: { mission: '20261109-nest x' } })
    expect(match('/desk')).toEqual({ tab: 'desk' })
  })
})

describe('fixture conversation pages', () => {
  it('pages the tail, then the cursor, then a send shows after the cursor', async () => {
    const src = fixtureSource()
    const first = await src.conversation({ pid: 4242 }, { tail: 5 })
    expect(first.entries.map((x) => x.id)).toEqual(CHAT.slice(-5).map((x) => x.id))
    expect(first.has_more).toBe(true)
    const older = await src.conversation({ pid: 4242 }, { before: first.before_cursor })
    expect(older.entries.at(-1)?.id).toBe(CHAT.at(-6)?.id)
    expect((await src.conversation({ pid: 4242 }, { cursor: first.next_cursor })).entries).toEqual([])
    await src.actions.send(4242, 'hello')
    const after = await src.conversation({ pid: 4242 }, { cursor: first.next_cursor })
    expect(after.entries.map((x) => [x.role, x.text])).toEqual([['user', 'hello']])
  })

  it('answers 404 for an unknown pid or mission', async () => {
    const src = fixtureSource()
    await expect(src.conversation({ pid: 1 }, {})).rejects.toMatchObject({ status: 404 })
    await expect(src.conversation({ mission: 'nope' }, {})).rejects.toMatchObject({ status: 404 })
  })
})
