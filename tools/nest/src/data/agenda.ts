import type { Snapshot } from './feed'
import { checksOf, type Check } from './home'
import { mrsOf, type Mr } from './mrs'
import type { SectionBody } from './types'
import { ago } from '../format'

// The Agenda page (item_226, board option B): seven days from today, the meetings of
// crow-today `week` and the dated items core already serves (crow-checks rows, MR waits).

export const DAYS = 7
const DAY = 86400000

export type DueKind = 'th' | 'mr' | 'dt'
export interface Meeting {
  key: string
  day: string
  title: string
  // Minutes since midnight, Paris wall clock (crow-today sends local times).
  from: number
  to: number
  start: string
  end: string
  state: '' | 'done' | 'now' | 'next'
}
export interface Due {
  key: string
  kind: DueKind
  day: string
  title: string
  sub: string
  check?: Check
  mr?: Mr
}
export interface AgendaDay {
  key: string
  date: Date
  today: boolean
  allDay: string[]
  meetings: Meeting[]
  due: Due[]
}
export interface AgendaData {
  off: boolean
  stale: boolean
  error?: string
  days: AgendaDay[]
  next: Meeting | null
  week: number
  now: number
}

interface WeekEvent { start?: unknown; end?: unknown; subject?: unknown; all_day?: unknown }

const str = (v: unknown) => (typeof v === 'string' ? v : '')
const obj = (b: SectionBody | undefined) => (b?.data && typeof b.data === 'object' ? (b.data as Record<string, unknown>) : {})
const list = <T>(v: unknown) => (Array.isArray(v) ? (v as T[]) : [])
const pad = (n: number) => String(n).padStart(2, '0')
export const dayKey = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
const minOf = (iso: string) => Number(iso.slice(11, 13)) * 60 + Number(iso.slice(14, 16))
export const hhmm = (m: number) => `${pad(Math.floor(m / 60))}:${pad(m % 60)}`

// ISO week number of a local date.
export function isoWeek(d: Date): number {
  const t = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()))
  t.setUTCDate(t.getUTCDate() + 4 - (t.getUTCDay() || 7))
  return Math.ceil(((t.getTime() - Date.UTC(t.getUTCFullYear(), 0, 1)) / DAY + 1) / 7)
}

// crow-today `week`; older payloads only carry today's events and tomorrow's.
function eventsOf(today: SectionBody | undefined): WeekEvent[] {
  const d = obj(today)
  return Array.isArray(d.week) ? list<WeekEvent>(d.week) : [...list<WeekEvent>(d.events), ...list<WeekEvent>(d.tomorrow)]
}

const LABEL: Record<string, string> = { mail: 'follow-up', sprint: 'sprint date', ritual: 'ritual', one_on_one: '1:1' }

function checkDue(c: Check, today: string): Due {
  const kind: DueKind = c.check === 'sprint' ? 'dt' : 'th'
  const day = c.due && c.due.slice(0, 10) > today ? c.due.slice(0, 10) : today
  return { key: `c:${c.id}`, kind, day, title: c.what, sub: [LABEL[c.check] ?? c.check, c.why].filter(Boolean).join(' · '), check: c }
}

function mrDue(m: Mr, today: string, now: Date): Due {
  return { key: `m:${m.project}:${m.ref}`, kind: 'mr', day: today, title: `${m.ref} review waits ${ago(m.created, now)}`, sub: `MR · ${m.name}`, mr: m }
}

export function agendaOf(snap: Snapshot): AgendaData {
  const { sections, now } = snap
  const body = sections.today
  const todayKey = dayKey(now)
  const nowMin = now.getHours() * 60 + now.getMinutes()
  const days: AgendaDay[] = Array.from({ length: DAYS }, (_, i) => {
    const date = new Date(now.getFullYear(), now.getMonth(), now.getDate() + i)
    return { key: dayKey(date), date, today: i === 0, allDay: [], meetings: [], due: [] }
  })
  const byKey = new Map(days.map((d) => [d.key, d]))

  for (const e of eventsOf(body)) {
    const start = str(e.start), end = str(e.end), title = str(e.subject) || '(no subject)'
    if (start.length < 16 || end.length < 16) continue
    if (e.all_day) {
      // An all-day event covers its start day up to the end day, excluded.
      for (const d of days) if (d.key >= start.slice(0, 10) && d.key < end.slice(0, 10) && !d.allDay.includes(title)) d.allDay.push(title)
      continue
    }
    const d = byKey.get(start.slice(0, 10))
    if (!d || d.meetings.some((m) => m.start === start && m.title === title)) continue
    const from = minOf(start)
    const to = end.slice(0, 10) === start.slice(0, 10) ? Math.max(minOf(end), from + 15) : 24 * 60
    const state = !d.today ? '' : to <= nowMin ? 'done' : from <= nowMin ? 'now' : ''
    d.meetings.push({ key: `${start} ${title}`, day: d.key, title, from, to, start, end, state })
  }
  for (const d of days) d.meetings.sort((a, b) => a.from - b.from || a.to - b.to)
  const next = days.flatMap((d) => d.meetings).find((m) => m.state !== 'done') ?? null
  if (next && next.state === '') next.state = 'next'

  // Due items: dated crow-checks rows (undated and overdue ones sit on today), then the
  // MRs that wait on the owner's review. Review checks repeat the MR rows, so they drop
  // when the forge serves MRs.
  const mrs = mrsOf(sections.mrs, now)
  const forge = !mrs.off && !!sections.mrs?.data
  const due = [
    ...checksOf(sections.checks).filter((c) => !(forge && c.check === 'review')).map((c) => checkDue(c, todayKey)),
    ...(forge ? mrs.review.map((m) => mrDue(m, todayKey, now)) : []),
  ]
  for (const x of due) byKey.get(x.day)?.due.push(x)

  return {
    off: !!body?.off, stale: !!body?.stale && !!body?.error, error: body?.error,
    days, next, week: isoWeek(now), now: nowMin,
  }
}

export interface AgendaCounts { todayMeetings: number; todayDue: number; meetings: number; due: number }
export function countsOf(a: AgendaData): AgendaCounts {
  const t = a.days[0]
  return {
    todayMeetings: t?.meetings.length ?? 0, todayDue: t?.due.length ?? 0,
    meetings: a.days.reduce((n, d) => n + d.meetings.length, 0), due: a.days.reduce((n, d) => n + d.due.length, 0),
  }
}

// Hour range of a grid: the board's 09:00 to 19:00, wider when a meeting falls outside.
export function hoursOf(days: AgendaDay[], lo: number, hi: number): [number, number] {
  const all = days.flatMap((d) => d.meetings)
  const h0 = Math.min(lo, ...all.map((m) => Math.floor(m.from / 60)))
  const h1 = Math.max(hi, ...all.map((m) => Math.ceil(m.to / 60)))
  return [h0, Math.min(h1, 24)]
}

// The session prompt of PREP: the meeting and where its context lives.
export function prepPrompt(m: Meeting, date: string): string {
  return `Prep the meeting "${m.title}" (${date}, ${hhmm(m.from)} to ${hhmm(m.to)}). Gather the context from the vault, open threads and recent mail, then draft the points to raise.`
}

// The session prompt of OPEN on a due check row.
export function duePrompt(x: Due): string {
  const c = x.check
  return c ? `${c.what}. ${c.why ? `${c.why}. ` : ''}Next step: ${c.action}.` : x.title
}
