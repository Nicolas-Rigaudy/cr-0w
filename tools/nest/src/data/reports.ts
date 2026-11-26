import type { SectionBody } from './types'

// Reports page (item_258): what is settled. Core serves the `reports` section (AFK
// reports with no open decision, weekly rollups, retros); the `jobs` section says
// whether the retro routine exists. Brief, daily and eod save no output yet.
export type Kind = 'afk' | 'brief' | 'daily' | 'weekly' | 'retro' | 'eod'
export const KINDS: Kind[] = ['afk', 'brief', 'daily', 'weekly', 'retro', 'eod']

// Glyph, colour class and tile name, as board B.
export const RITUAL: Record<Kind, { g: string; c: string; name: string }> = {
  afk: { g: '☾', c: 'c-m', name: 'AFK' },
  brief: { g: '☀', c: 'c-n', name: 'Brief' },
  daily: { g: '◷', c: 'c-w', name: 'Daily' },
  weekly: { g: '▦', c: 'c-t', name: 'Weekly' },
  retro: { g: '↺', c: 'c-f', name: 'Retro' },
  eod: { g: '☽', c: 'c-s', name: 'Eod' },
}

// The only ritual with a routine core can start (POST /api/routines/<unit>/run).
// Brief, daily, weekly and eod wait for item_234.
export const RETRO_UNIT = 'crow-retro.timer'

export interface Report {
  id: string
  kind: Kind
  title: string
  at: string
  summary: string
  path: string
  link: string
  body: string
}

export interface ReportsData {
  off: boolean
  stale: boolean
  error: string
  items: Report[]
  held: number
  // The retro routine is installed (jobs carries crow-retro.timer).
  retro: boolean
}

const obj = (v: unknown) => (v && typeof v === 'object' && !Array.isArray(v) ? (v as Record<string, unknown>) : {})
const arr = <T>(v: unknown) => (Array.isArray(v) ? (v as T[]) : [])
const num = (v: unknown) => (typeof v === 'number' ? v : Number(v) || 0)
const str = (v: unknown) => (typeof v === 'string' ? v : v == null ? '' : String(v))

function reportOf(v: unknown): Report | null {
  const r = obj(v)
  const kind = str(r.kind) as Kind
  if (!KINDS.includes(kind)) return null
  return {
    id: str(r.id), kind, title: str(r.title) || str(r.id), at: str(r.at),
    summary: str(r.summary), path: str(r.path), link: str(r.link), body: str(r.body),
  }
}

export function reportsOf(b: SectionBody | undefined, jobs: SectionBody | undefined): ReportsData {
  const d = obj(b?.data)
  const items = arr<unknown>(d.items).map(reportOf).filter((r): r is Report => !!r)
  items.sort((x, y) => Date.parse(y.at) - Date.parse(x.at))
  return {
    off: !!b?.off,
    stale: !!b?.stale && !!b?.error,
    error: str(b?.error),
    items,
    held: num(d.held),
    retro: arr<unknown>(jobs?.data).some((j) => str(obj(j).unit) === RETRO_UNIT),
  }
}

// Monday 00:00 of the week of now, local time.
export function weekStart(now: Date) {
  const d = new Date(now)
  d.setHours(0, 0, 0, 0)
  d.setDate(d.getDate() - ((d.getDay() + 6) % 7))
  return d
}

export interface Tile {
  kind: Kind
  week: number
  last: Report | null
}

export function tilesOf(items: Report[], now: Date): Tile[] {
  const from = weekStart(now).getTime()
  return KINDS.map((kind) => {
    const mine = items.filter((r) => r.kind === kind)
    return { kind, week: mine.filter((r) => Date.parse(r.at) >= from).length, last: mine[0] ?? null }
  })
}

// "Wed 25 08:05", the board's time stamp.
const wd = new Intl.DateTimeFormat('en-GB', { weekday: 'short' })
const hm = new Intl.DateTimeFormat('en-GB', { hour: '2-digit', minute: '2-digit' })
export function stamp(iso: string) {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return `${wd.format(d)} ${String(d.getDate()).padStart(2, '0')} ${hm.format(d)}`
}

// A retro run started here is shown running until a newer retro appears (30 min cap).
export const RUN_CAP_MS = 30 * 60 * 1000
export function running(started: string | null, items: Report[], now: Date) {
  if (!started) return false
  const t = Date.parse(started)
  if (now.getTime() - t > RUN_CAP_MS) return false
  return !items.some((r) => r.kind === 'retro' && Date.parse(r.at) >= t)
}

// Minimal Markdown for the reader: headings, lists, paragraphs, tables, code fences,
// **bold**, `code`, [text](url). Front matter is dropped. Raw HTML stays text.
export type Block =
  | { t: 'h'; level: 1 | 2 | 3; text: string }
  | { t: 'ul'; items: string[] }
  | { t: 'p'; text: string }
  | { t: 'table'; head: string[]; rows: string[][] }
  | { t: 'pre'; text: string }

const cells = (l: string) => l.replace(/^\|/, '').replace(/\|$/, '').split(/(?<!\\)\|/).map((c) => c.trim().replace(/\\\|/g, '|'))
const isSep = (l: string) => /^\|?\s*:?-{2,}:?\s*(\|\s*:?-{2,}:?\s*)*\|?$/.test(l)

export function blocksOf(md: string): Block[] {
  let lines = md.replace(/\r\n/g, '\n').split('\n')
  if (lines[0]?.trim() === '---') {
    const end = lines.findIndex((l, i) => i > 0 && l.trim() === '---')
    if (end > 0) lines = lines.slice(end + 1)
  }
  const out: Block[] = []
  let para: string[] = []
  let list: string[] | null = null
  const flush = () => {
    if (para.length) out.push({ t: 'p', text: para.join(' ') })
    para = []
    if (list) out.push({ t: 'ul', items: list })
    list = null
  }
  let table: { head: string[]; rows: string[][] } | null = null
  let fence: string[] | null = null
  for (let i = 0; i < lines.length; i++) {
    const raw = lines[i]
    const l = raw.trim()
    if (fence) {
      if (l.startsWith('```')) { out.push({ t: 'pre', text: fence.join('\n') }); fence = null } else fence.push(raw)
      continue
    }
    if (l.startsWith('```')) { flush(); fence = []; continue }
    if (table) {
      if (l.startsWith('|')) { table.rows.push(cells(l)); continue }
      out.push({ t: 'table', ...table }); table = null
    }
    if (l.startsWith('|') && isSep(lines[i + 1]?.trim() ?? '')) { flush(); table = { head: cells(l), rows: [] }; i++; continue }
    const h = /^(#{1,6})\s+(.*)$/.exec(l)
    const li = /^[-*+]\s+(.*)$/.exec(l) ?? /^\d+[.)]\s+(.*)$/.exec(l)
    if (!l) flush()
    else if (h) { flush(); out.push({ t: 'h', level: Math.min(3, h[1].length) as 1 | 2 | 3, text: h[2] }) }
    else if (li) { if (para.length) { out.push({ t: 'p', text: para.join(' ') }); para = [] } (list ??= []).push(li[1].replace(/^\[[ xX]\]\s+/, '')) }
    else if (list && /^\s/.test(raw)) list[list.length - 1] += ` ${l}`
    else { if (list) { out.push({ t: 'ul', items: list }); list = null } para.push(l) }
  }
  if (table) out.push({ t: 'table', ...table })
  if (fence) out.push({ t: 'pre', text: fence.join('\n') })
  flush()
  return out
}

export type Span = { t: 'text' | 'b' | 'code'; s: string } | { t: 'a'; s: string; href: string }

export function spansOf(text: string): Span[] {
  const out: Span[] = []
  const re = /\*\*([^*]+)\*\*|`([^`]+)`|\[([^\]]+)\]\(([^)\s]+)\)|\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/g
  let i = 0
  for (let m = re.exec(text); m; m = re.exec(text)) {
    if (m.index > i) out.push({ t: 'text', s: text.slice(i, m.index) })
    if (m[1] !== undefined) out.push({ t: 'b', s: m[1] })
    else if (m[2] !== undefined) out.push({ t: 'code', s: m[2] })
    else if (m[3] !== undefined) out.push(/^(https?|obsidian):/.test(m[4]) ? { t: 'a', s: m[3], href: m[4] } : { t: 'text', s: m[3] })
    else out.push({ t: 'text', s: m[6] ?? m[5] })
    i = m.index + m[0].length
  }
  if (i < text.length) out.push({ t: 'text', s: text.slice(i) })
  return out
}
