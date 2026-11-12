// Conversation pages of a session or a mission's last run (item_177 chat). The core
// routes GET /api/sessions/{pid}/conversation and /api/missions/{id}/conversation pass
// `cdx conversation` pages through: entries in order, a cursor for what came after,
// a before cursor for older turns.

export interface Entry {
  at: string
  id: string
  role: 'user' | 'assistant' | 'tool'
  text: string
  tool_name?: string | null
  tool_summary?: string | null
}

export interface ConversationPage {
  account: string
  conversation_id: string
  entries: Entry[]
  next_cursor?: string
  before_cursor?: string
  has_more: boolean
  reset_required?: boolean
}

export type ChatTarget = { pid: number } | { mission: string }
export interface PageAsk { tail?: number; cursor?: string; before?: string }

// What the chat shows: a message, or one fold for a run of tool entries.
export type Block =
  | { kind: 'me' | 'cr'; id: string; at: string; text: string; who: boolean }
  | { kind: 'tools'; id: string; calls: number; tools: string[]; entries: Entry[] }

// Consecutive tool entries fold into one line: calls counted on "Call recorded", tool
// names listed once in first-use order. A message shows its who label unless the
// message just before it (folds aside) came from the same side.
export function blocksOf(entries: Entry[]): Block[] {
  const out: Block[] = []
  let last: 'me' | 'cr' | null = null
  for (const e of entries) {
    if (e.role === 'tool') {
      let f = out[out.length - 1]
      if (f?.kind !== 'tools') out.push((f = { kind: 'tools', id: e.id, calls: 0, tools: [], entries: [] }))
      f.entries.push(e)
      if ((e.tool_summary ?? '').startsWith('Call')) f.calls++
      if (e.tool_name && !f.tools.includes(e.tool_name)) f.tools.push(e.tool_name)
      continue
    }
    if (!e.text.trim()) continue
    const kind = e.role === 'user' ? 'me' : 'cr'
    out.push({ kind, id: e.id, at: e.at, text: e.text, who: last !== kind })
    last = kind
  }
  // A fold of results alone (call on an earlier page) still counts as one call.
  for (const b of out) if (b.kind === 'tools' && b.calls === 0) b.calls = 1
  return out
}

// Assistant text as display parts: paragraphs, list lines, code blocks, headings.
export type Part =
  | { t: 'p' | 'h'; text: string }
  | { t: 'li'; mark: string; text: string }
  | { t: 'code'; text: string }

export function partsOf(text: string): Part[] {
  const out: Part[] = []
  const lines = text.replace(/\r/g, '').split('\n')
  let para: string[] = []
  const flush = () => { if (para.length) out.push({ t: 'p', text: para.join('\n') }); para = [] }
  for (let i = 0; i < lines.length; i++) {
    const l = lines[i]
    if (/^\s*```/.test(l)) {
      flush()
      const code: string[] = []
      while (++i < lines.length && !/^\s*```/.test(lines[i])) code.push(lines[i])
      out.push({ t: 'code', text: code.join('\n') })
      continue
    }
    const h = /^#{1,6}\s+(.*)$/.exec(l)
    const li = /^(\s*)([-*•]|\d+[.)])\s+(.*)$/.exec(l)
    if (!l.trim()) flush()
    else if (h) { flush(); out.push({ t: 'h', text: h[1] }) }
    else if (li) { flush(); out.push({ t: 'li', mark: /\d/.test(li[2]) ? li[2] : '•', text: li[3] }) }
    else para.push(l)
  }
  flush()
  return out
}

// Inline marks: **bold** and `code`; everything else stays as written.
export type Span = { t: 'txt' | 'b' | 'code'; text: string }
export function spansOf(s: string): Span[] {
  const out: Span[] = []
  const re = /\*\*([^*]+)\*\*|`([^`]+)`/g
  let at = 0
  for (let m; (m = re.exec(s)); at = re.lastIndex) {
    if (m.index > at) out.push({ t: 'txt', text: s.slice(at, m.index) })
    out.push(m[1] != null ? { t: 'b', text: m[1] } : { t: 'code', text: m[2] })
  }
  if (at < s.length) out.push({ t: 'txt', text: s.slice(at) })
  return out
}

export const conversationPath = (t: ChatTarget) =>
  'pid' in t ? `/sessions/${t.pid}/conversation` : `/missions/${encodeURIComponent(t.mission)}/conversation`

export function pageQuery(p: PageAsk) {
  if (p.cursor) return `?cursor=${encodeURIComponent(p.cursor)}`
  if (p.before) return `?before=${encodeURIComponent(p.before)}`
  return p.tail ? `?tail=${p.tail}` : ''
}
