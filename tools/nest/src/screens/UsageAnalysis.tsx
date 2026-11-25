import { useState, type ReactNode } from 'react'
import type { CostDay, UsageData, UsageLimit, UsageSession } from '../data/usage'
import {
  bars, bySession, cut, dayLimits, dayShort, daySummary, limitAccounts, localDay, modelName, modelsOf, PERIODS,
  periodLabel, sessionsIn, sumGroups, TABS, tok, topRuns, UNITS, usd, type Period, type Tab, type Unit,
} from '../data/usageView'
import { useQuery } from '../router'
import { Deck } from './Mrs'
import { COLORS, DayChart, HBars, Legend } from './usageCharts'

// The ANALYSIS pane of the Usage page: the D/W/M period, one tab (one chart) at a time.

export interface AnalysisProps {
  data: UsageData
  lim?: UsageLimit
  now: Date
  period: Period
  setPeriod: (p: Period) => void
  phone?: boolean
}

const pctFmt = (v: number) => `${Math.round(v)}%`
const TOK_LEGEND: Array<[string, string]> = [['cache read', COLORS[0]], ['cache write', COLORS[1]], ['input + output', COLORS[2]]]

function Segs<T extends string>({ label, items, on, set }: { label: string; items: Array<[T, string]>; on: T; set: (v: T) => void }) {
  return (
    <nav className="segs" aria-label={label}>
      {items.map(([k, t]) => (
        <button key={k} type="button" className={`seg${k === on ? ' on' : ''}`} aria-pressed={k === on} onClick={() => set(k)}>{t}</button>
      ))}
    </nav>
  )
}

export function Analysis({ data, lim, now, period, setPeriod, phone }: AnalysisProps) {
  // ?tab= and ?unit= open a tab directly (links, shots).
  const q = useQuery()
  const [tab, setTab] = useState<Tab>(TABS.find(([k]) => k === q.get('tab'))?.[0] ?? 'day')
  const [unit, setUnit] = useState<Unit>(UNITS.find(([k]) => k === q.get('unit'))?.[0] ?? 'usd')
  const all = data.cost.data?.by_day ?? []
  const days = cut(all, period)
  const sess = sessionsIn(data.sessions, days)
  const tabs = TABS.map(([k, t, s]) => [k, phone ? s : t] as [Tab, string])
  const periods = PERIODS.map(([k]) => [k, k] as [Period, string])
  const top = (
    <div className="ubar-top">
      <Segs label="Analysis" items={tabs} on={tab} set={setTab} />
      {!phone && <span className="dim">one chart at a time</span>}
      <span className="uctl">
        {tab === 'day' && <Segs label="Unit" items={UNITS} on={unit} set={setUnit} />}
        <Segs label="Period" items={periods} on={period} set={setPeriod} />
      </span>
    </div>
  )
  let body: ReactNode
  if (!days.length && tab !== 'session' && tab !== 'runs') body = <div className="empty">No cost audit yet.</div>
  else if (tab === 'day') body = <PerDay days={days} lim={lim} limits={data.limits} unit={unit} period={period} now={now} phone={phone} />
  else if (tab === 'project') body = <Split rows={sumGroups(days, (d) => d.by_project)} period={period} what="project" phone={phone} />
  else if (tab === 'model') body = <Split rows={sumGroups(days, (d) => d.by_model)} period={period} what="model" phone={phone} name={modelName} colors />
  else if (tab === 'kind') body = <Split rows={sumGroups(days, (d) => d.by_kind)} period={period} what="kind" phone={phone} color="var(--c3)" />
  else if (tab === 'session') body = <Sessions rows={bySession(sess)} period={period} phone={phone} />
  else body = <Sessions rows={topRuns(sess)} period={period} phone={phone} runs />
  return (
    <Deck n={phone ? 4 : 3} title="Analysis">
      {top}
      {body}
    </Deck>
  )
}

function PerDay({ days, lim, limits, unit, period, now, phone }: {
  days: CostDay[]; lim?: UsageLimit; limits: UsageLimit[]; unit: Unit; period: Period; now: Date; phone?: boolean
}) {
  const [pickDate, setPick] = useState<string | null>(null)
  const accounts = limitAccounts(days, limits)
  const b = bars(days, unit, accounts)
  const i = Math.max(0, days.findIndex((d) => d.date === pickDate))
  const pick = pickDate && days.some((d) => d.date === pickDate) ? i : days.length - 1
  const d = days[pick]
  const legend = unit === 'usd'
    ? modelsOf(days).map((m, k) => [modelName(m), COLORS[k % COLORS.length]] as [string, string])
    : unit === 'tok' ? TOK_LEGEND : accounts.map((a, k) => [a, COLORS[k % COLORS.length]] as [string, string])
  const axis = unit === 'usd' ? '$ per day' : unit === 'tok' ? 'tokens per day' : '% of a weekly limit used per day'
  const note = unit === 'lim' ? '100% = one account-week' : undefined
  const cells = lim ? days.map((x) => { const l = dayLimits(x, lim.account); return l === null ? null : l?.windows_5h }) : undefined
  const chart = (
    <DayChart w={phone ? 330 : 620} h={phone ? 230 : 250} bars={b} fmt={unit === 'usd' ? usd : unit === 'tok' ? tok : pctFmt} axis={axis}
      colors={COLORS} pick={pick} onPick={(k) => setPick(days[k]?.date ?? null)} today={localDay(now)} cells={cells} fit={phone} />
  )
  const line = d && (
    <div className="uday" aria-live="polite"><b>{dayShort(d.date)}</b> {daySummary(d, lim?.account)}</div>
  )
  const legendEl = <Legend items={legend} note={note ?? (lim ? `5h cells: ${lim.account}` : undefined)} />
  const proj = sumGroups(days, (x) => x.by_project).map(([k, g]) => [k, g.cost_usd] as [string, number])
  if (phone) return <>{line}{chart}{legendEl}<h3 className="ust">Per project · {periodLabel(period)}</h3><HBars rows={proj} fmt={usd} label={60} /></>
  return (
    <div className="uan">
      <div className="umain">{chart}{line}</div>
      <div className="uside">{legendEl}<h3 className="ust">{periodLabel(period)}</h3><HBars rows={proj} fmt={usd} label={80} /></div>
    </div>
  )
}

function Split({ rows, period, what, phone, name = (k) => k, colors, color }: {
  rows: Array<[string, { cost_usd: number; tokens: number; turns: number }]>; period: Period; what: string; phone?: boolean
  name?: (k: string) => string; colors?: boolean; color?: string
}) {
  if (!rows.length) return <div className="empty">No {what} data in {periodLabel(period)}.</div>
  const total = rows.reduce((s, [, g]) => s + g.cost_usd, 0) || 1
  const cost = rows.map(([k, g], i) => [name(k), g.cost_usd, colors ? COLORS[i % COLORS.length] : undefined] as [string, number, string?])
  const toks = rows.map(([k, g], i) => [name(k), g.tokens, colors ? COLORS[i % COLORS.length] : undefined] as [string, number, string?])
  const label = what === 'kind' ? (phone ? 90 : 150) : phone ? 70 : 100
  return (
    <div className={phone ? 'usplit ph' : 'usplit'}>
      <div><h3 className="ust">$ per {what} · {periodLabel(period)}</h3><HBars rows={cost} fmt={usd} label={label} color={color} /></div>
      <div><h3 className="ust">tokens per {what}</h3><HBars rows={toks} fmt={tok} label={label} color={color} /></div>
      <div className="uleg"><span className="dim">{rows.map(([k, g]) => `${name(k)} ${Math.round((100 * g.cost_usd) / total)}% · ${g.turns.toLocaleString('en-GB')} turns`).join(' · ')}</span></div>
    </div>
  )
}

const models = (s: UsageSession) => (s.models ?? []).map(modelName).filter((m, i, a) => a.indexOf(m) === i).join(' + ') || '–'
const when = (s: UsageSession) => (s.start ? dayShort(localDay(new Date(s.start))) : '–')

// Cost cell: a thin bar scaled to the largest cost in the list, then the value.
function Cost({ v, mx }: { v: number; mx: number }) {
  return (
    <>
      <span className="ucb" aria-hidden="true"><i style={{ width: `${(v / mx) * 100}%` }} /></span>
      <span className="ucv">{usd(v)}</span>
    </>
  )
}

function Sessions({ rows, period, phone, runs }: { rows: UsageSession[]; period: Period; phone?: boolean; runs?: boolean }) {
  if (!rows.length) return <div className="empty">No sessions in {periodLabel(period)}.</div>
  const head = runs ? 'Top runs' : 'Session'
  const mx = Math.max(0, ...rows.map((s) => s.cost_usd ?? 0)) || 1
  if (phone) {
    return (
      <div className="utbl ph" role="table" aria-label={head}>
        {rows.map((s) => (
          <div key={s.id} className="utr" role="row">
            <span className="ti" role="cell">{s.title ?? s.id}<span className="sub">{s.project ?? '–'} · {runs ? s.kind ?? '–' : models(s)} · {tok(s.tokens ?? 0)} · {when(s)}</span></span>
            <span className="num ucost" role="cell"><Cost v={s.cost_usd ?? 0} mx={mx} /></span>
          </div>
        ))}
      </div>
    )
  }
  const cols = runs ? ['Top runs', 'Project', 'Kind', 'Account', 'Tokens', 'Cost'] : ['Session', 'Project', 'Model', 'Day', 'Tokens', 'Cost']
  return (
    <div className={`utbl${runs ? ' runs' : ''}`} role="table" aria-label={head}>
      <div className="utr uth" role="row">{cols.map((c, i) => <span key={c} role="columnheader" className={i >= 4 ? 'rt' : undefined}>{c}</span>)}</div>
      {rows.map((s) => (
        <div key={s.id} className={`utr${s.active ? ' uact' : ''}`} role="row">
          <span role="cell" className="ti" title={s.title}>{s.title ?? s.id}</span>
          <span role="cell" className="dim">{s.project ?? '–'}</span>
          <span role="cell" className="dim">{runs ? s.kind ?? '–' : models(s)}</span>
          <span role="cell" className="dim">{runs ? s.account : when(s)}</span>
          <span role="cell" className="rt num">{tok(s.tokens ?? 0)}</span>
          <span role="cell" className="rt num ucost"><Cost v={s.cost_usd ?? 0} mx={mx} /></span>
        </div>
      ))}
    </div>
  )
}
