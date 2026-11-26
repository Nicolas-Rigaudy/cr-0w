import { useState, type CSSProperties } from 'react'
import { ApiError } from '../data/api'
import { agendaOf, countsOf, duePrompt, hhmm, hoursOf, prepPrompt, type AgendaData, type AgendaDay, type Due, type DueKind, type Meeting } from '../data/agenda'
import type { Actions, LiveReply } from '../data/live'
import { useSource, type Snapshot } from '../data/useData'
import { date, weekday } from '../format'
import { deskOpening, navigate, useQuery } from '../router'
import { Sheet, type Confirm } from '../ui/Sheet'
import { useMedia } from '../ui/useMedia'
import { confirmsOf as mrConfirms } from './Mrs'

// The Agenda page (item_226, board option B): the week as a grid with a due row on
// laptop, a day strip over one day's grid on phone, the picked meeting and the due list.

const LAPTOP = '(min-width: 1100px)'
// Meeting context lives in the hub (vault, threads, mail): PREP opens a hub session.
export const HUB = 'cr-0w'
const KIND: Record<DueKind, [string, string]> = { th: ['◆', 'c-n'], mr: ['⇄', 'c-t'], dt: ['⚑', 'c-o'] }
const dayLabel = (d: AgendaDay) => `${weekday(d.date)} ${String(d.date.getDate()).padStart(2, '0')}`
const span = (m: Meeting) => `${hhmm(m.from)} → ${hhmm(m.to)}`

type Said = { text: string; f?: boolean }
// A session that starts opens on the Desk (`open`).
export function confirmsOf(act: Actions, refresh: () => Promise<void>, setSaid: (s: Said) => void, open = (tmux: string) => navigate(deskOpening(tmux))) {
  const wrap = (call: () => Promise<LiveReply>) => async () => {
    try {
      const r = await call()
      setSaid({ text: r.message || `${r.action} ${r.target}`, f: !r.ok })
      if (r.ok && r.tmux) open(r.tmux)
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void refresh()
      throw e
    }
    void refresh()
  }
  const mr = mrConfirms(act, refresh, setSaid, open)
  return {
    prep: (m: Meeting, d: AgendaDay): Confirm => ({
      title: `Prep ${m.title}`, lead: 'Opens a hub session with this meeting as its context.', quote: prepPrompt(m, date(d.date)),
      target: [['Meeting', m.title], ['When', `${date(d.date)} · ${span(m)}`]], note: 'The session reads and drafts. Nothing is sent without your OK.', ok: 'Open session',
      run: wrap(() => act.start({ project: HUB, prompt: prepPrompt(m, date(d.date)) })),
    }),
    due: (x: Due, d: AgendaDay): Confirm => x.mr ? mr.review(x.mr) : ({
      title: x.title, lead: 'Opens a hub session on this item.', quote: duePrompt(x),
      target: [['Item', x.title], ['Due', `${date(d.date)} · ${x.sub}`]], note: 'Nothing is sent without your OK.', ok: 'Open session',
      run: wrap(() => act.start({ project: HUB, prompt: duePrompt(x) })),
    }),
  }
}

export function useAgenda(snap: Snapshot, refresh: () => Promise<void>) {
  const act = useSource().actions
  const q = useQuery()
  const a = agendaOf(snap)
  const [sheet, setSheet] = useState<Confirm | null>(null)
  const [said, setSaid] = useState<Said | null>(null)
  const c = confirmsOf(act, refresh, setSaid)
  const all = a.days.flatMap((d) => d.meetings)
  const [pick, setPick] = useState<string | null>(null)
  const picked = all.find((m) => m.key === pick) ?? a.next
  // Phone: ?day=0..6 picks the day of the strip (shots, links).
  const qd = Number(q.get('day') ?? 0)
  const [day, setDay] = useState(Number.isInteger(qd) && qd >= 0 && qd < a.days.length ? qd : 0)
  const dayOf = (k: string) => a.days.find((d) => d.key === k)!
  return {
    a, sheet, setSheet, said, picked, pick: (m: Meeting) => setPick(m.key), day, setDay,
    prep: (m: Meeting) => setSheet(c.prep(m, dayOf(m.day))),
    due: (x: Due) => setSheet(c.due(x, dayOf(x.day))),
  }
}
export type G = ReturnType<typeof useAgenda>

function Head({ g, phone }: { g: G; phone?: boolean }) {
  const c = countsOf(g.a)
  const today = g.a.days[0]
  if (phone) return <div className="ahead"><h1>Agenda</h1><span className="sum">W{g.a.week}</span></div>
  return (
    <div className="ahead">
      <h1>Agenda</h1>
      {today && <span className="sum0">{date(today.date)}</span>}
      <div className="sum">today {c.todayMeetings} meetings · {c.todayDue} due · week {c.meetings} meetings · {c.due} due</div>
    </div>
  )
}

function Legend() {
  return (
    <div className="legend">
      <span><span className="c-w">◷</span> meeting</span>
      <span><span className="c-n">◆</span> thread / follow-up</span>
      <span><span className="c-t">⇄</span> MR wait</span>
      <span><span className="c-o">⚑</span> sprint / release date</span>
      <span>times: Paris</span>
    </div>
  )
}

function Block({ g, m, h0, rh }: { g: G; m: Meeting; h0: number; rh: number }) {
  const on = g.picked === m
  const style = { top: ((m.from - h0 * 60) / 60) * rh, height: Math.max(((m.to - m.from) / 60) * rh - 2, 16) }
  return (
    <button type="button" className={`ev ${m.state}${on ? ' on' : ''}`} style={style} aria-current={on || undefined} title={`${span(m)} ${m.title}`} onClick={() => g.pick(m)}>
      {hhmm(m.from)} {m.title}
    </button>
  )
}

function NowLine({ a, h0, h1, rh, label }: { a: AgendaData; h0: number; h1: number; rh: number; label: string }) {
  if (a.now < h0 * 60 || a.now > h1 * 60) return null
  return <div className="nowl" style={{ top: ((a.now - h0 * 60) / 60) * rh }}><b>{label}</b></div>
}

// Laptop: hours, then per day a header, the all-day chips, the due chips and the meetings.
function Week({ g }: { g: G }) {
  const { days } = g.a
  const [h0, h1] = hoursOf(days, 9, 19)
  const rh = 40
  const cols = days.map((d) => (d.meetings.length || d.due.length ? '1fr' : '.45fr')).join(' ')
  const style = { gridTemplateColumns: `40px ${cols}`, '--rh': `${rh}px` } as CSSProperties
  const hours = Array.from({ length: h1 - h0 }, (_, i) => h0 + i)
  const last = days[days.length - 1]
  return (
    <section className="apn wk" aria-label="Week">
      <div className="ph"><span>Week {g.a.week} · {date(days[0].date)} → {last && dayLabel(last)}</span><span className="c-s">click a block or a due chip</span></div>
      <div className="wg" style={style}>
        <span />{days.map((d) => <span key={d.key} className={`dh${d.today ? ' on' : ''}`}>{dayLabel(d)}</span>)}
        <span />{days.map((d) => <div key={d.key} className="adr">{d.allDay.map((t) => <span key={t} className="ad">{t}</span>)}</div>)}
        <span />{days.map((d) => (
          <div key={d.key} className="duer">
            {d.due.map((x) => <button key={x.key} type="button" className={`due ${x.kind}`} title={`${x.title} · ${x.sub}`} onClick={() => g.due(x)}>{KIND[x.kind][0]} {x.title}</button>)}
          </div>
        ))}
        <div className="hrs" style={{ height: (h1 - h0) * rh }}>{hours.map((h) => <span key={h}>{String(h).padStart(2, '0')}:00</span>)}</div>
        {days.map((d) => (
          <div key={d.key} className="col" style={{ height: (h1 - h0) * rh }}>
            {d.meetings.map((m) => <Block key={m.key} g={g} m={m} h0={h0} rh={rh} />)}
            {d.today && <NowLine a={g.a} h0={h0} h1={h1} rh={rh} label="now" />}
          </div>
        ))}
      </div>
      <Legend />
    </section>
  )
}

// Phone: one day as an hour grid, the hours in a left column.
function DayGrid({ g, d }: { g: G; d: AgendaDay }) {
  const [h0, h1] = hoursOf([d], 9, 18)
  const rh = 30
  const hours = Array.from({ length: h1 - h0 }, (_, i) => h0 + i)
  return (
    <div className="hg" style={{ '--rh': `${rh}px` } as CSSProperties}>
      {hours.map((h) => <span key={h} className="h">{String(h).padStart(2, '0')}:00</span>)}
      <div className="evs">
        {d.meetings.map((m) => <Block key={m.key} g={g} m={m} h0={h0} rh={rh} />)}
        {d.today && <NowLine a={g.a} h0={h0} h1={h1} rh={rh} label={`now ${hhmm(g.a.now)}`} />}
      </div>
    </div>
  )
}

function Card({ g }: { g: G }) {
  const m = g.picked
  if (!m) return <section className="apn card"><div className="empty">No meeting in the next 7 days.</div></section>
  const d = g.a.days.find((x) => x.key === m.day)!
  const head = m.state === 'now' ? <span className="c-m">Now · ends in {m.to - g.a.now} min</span>
    : m === g.a.next && d.today ? <span className="c-m">Next · in {m.from - g.a.now} min</span>
      : <span className={m === g.a.next ? 'c-m' : ''}>{m === g.a.next ? 'Next · ' : ''}{date(d.date)}</span>
  return (
    <section className={`apn card${m === g.a.next ? ' on' : ''}`} aria-label="Meeting">
      <div className="ph">{head}<span className="c-s">{span(m)}</span></div>
      <div className="ttl">{m.title}</div>
      <div className="kv">
        <span>day</span><span>{date(d.date)}{d.allDay.length ? ` · ${d.allDay.join(', ')}` : ''}</span>
        <span>length</span><span>{m.to - m.from} min</span>
      </div>
      <div className="btnrow">
        <button type="button" className="btn p" disabled={m.state === 'done'} onClick={() => g.prep(m)}>Prep session</button>
      </div>
    </section>
  )
}

function DueRow({ g, x, d }: { g: G; x: Due; d: AgendaDay }) {
  const [icon, c] = KIND[x.kind]
  return (
    <div className="arow">
      <span className={`k ${c}`}>{icon}</span>
      <span className="t">{x.title}</span>
      <span className="s">{d.today ? 'Today' : dayLabel(d)} · {x.sub}</span>
      <button type="button" className="btn sm" onClick={() => g.due(x)}>{x.kind === 'mr' ? 'Review' : 'Open'}</button>
    </div>
  )
}

function DueList({ g }: { g: G }) {
  const rows = g.a.days.flatMap((d) => d.due.map((x) => ({ x, d })))
  return (
    <section className="apn dues" aria-label="Due this week">
      <div className="ph"><span>Due this week</span><span className="c-s">{rows.length}</span></div>
      <div className="alist">
        {rows.map(({ x, d }) => <DueRow key={x.key} g={g} x={x} d={d} />)}
        {!rows.length && <div className="empty">Nothing due this week.</div>}
      </div>
      {g.said && <div className={`note${g.said.f ? ' c-f' : ''}`} role="status">{g.said.text}</div>}
    </section>
  )
}

// Phone: the day strip (meeting dots, ◆ when something is due), then the picked day.
function Strip({ g }: { g: G }) {
  return (
    <div className="strip" role="tablist" aria-label="Days">
      {g.a.days.map((d, i) => (
        <button key={d.key} type="button" role="tab" className={`dtab${g.day === i ? ' on' : ''}`} aria-selected={g.day === i} onClick={() => g.setDay(i)}>
          <span>{weekday(d.date)}</span>
          <b>{String(d.date.getDate()).padStart(2, '0')}</b>
          <span className="dots">{'•'.repeat(Math.min(d.meetings.length, 3))}{d.due.length ? '◆' : ''}</span>
        </button>
      ))}
    </div>
  )
}

function PhoneDay({ g }: { g: G }) {
  const d = g.a.days[g.day]
  return (
    <section className="apn day" aria-label={dayLabel(d)}>
      <div className="ph"><span>{dayLabel(d)}{d.allDay.length ? ` · ${d.allDay.join(', ')}` : ''}</span><span className="c-s">{d.meetings.length}{d.due.length ? ` · ${d.due.length} due` : ''}</span></div>
      {d.due.map((x) => (
        <button key={x.key} type="button" className={`due ${x.kind}`} onClick={() => g.due(x)}>
          <span className="t">{KIND[x.kind][0]} {x.title}</span><span className="act">{x.kind === 'mr' ? 'REVIEW' : 'OPEN'}</span>
        </button>
      ))}
      <DayGrid g={g} d={d} />
      {g.said && <div className={`note${g.said.f ? ' c-f' : ''}`} role="status">{g.said.text}</div>}
    </section>
  )
}

export function Agenda({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const g = useAgenda(snap, refresh)
  const laptop = useMedia(LAPTOP)
  if (g.a.off) return <div className="desk agenda"><div className="ahead"><h1>Agenda</h1></div><div className="empty">The mail module is off.</div></div>
  const stale = g.a.stale && <div className="empty c-f">The calendar did not refresh: {g.a.error}</div>
  const sheet = g.sheet && <Sheet c={g.sheet} onClose={() => g.setSheet(null)} />
  if (!laptop) {
    return (
      <div className="desk agenda m-ph">
        <Head g={g} phone />
        {stale}
        <Strip g={g} />
        <PhoneDay g={g} />
        <Card g={g} />
        {sheet}
      </div>
    )
  }
  return (
    <div className="desk agenda m-lp">
      <Head g={g} />
      {stale}
      <Week g={g} />
      <div className="abot"><Card g={g} /><DueList g={g} /></div>
      {sheet}
    </div>
  )
}
