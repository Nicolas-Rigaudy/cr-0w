import { useEffect, useState, type ReactNode } from 'react'
import { useSource } from '../data/useData'
import type { UsageData, UsageLimit } from '../data/usage'
import { cut, floorOf, localDay, periodLabel, tiles, tok, usd, weekOf, window5h, type Period } from '../data/usageView'
import { time } from '../format'
import { Link, useQuery, usagePath } from '../router'
import { useMedia } from '../ui/useMedia'
import { Deck } from './Mrs'
import { Analysis } from './UsageAnalysis'
import { dayTicks, hourTicks, Legend, Spark, WindowChart } from './usageCharts'

// Usage page (item_227, board usage.html option C): the account list with burn sparklines
// and the tiles on the left; the picked account's 5h and week windows on the right, then
// the analysis pane. Phone: account picker, the two windows stacked, then the analysis.

const LAPTOP = '(min-width: 1100px)'

export function useUsage(): [UsageData | null, string] {
  const src = useSource()
  const [data, setData] = useState<UsageData | null>(null)
  const [err, setErr] = useState('')
  useEffect(() => {
    let live = true
    if (!src.usage) return
    const load = () => src.usage!().then((d) => { if (live) { setData(d); setErr('') } }, (e: Error) => { if (live) setErr(e.message) })
    void load()
    const t = src.poll ? setInterval(load, Math.max(src.poll, 30_000)) : undefined
    return () => { live = false; clearInterval(t) }
  }, [src])
  return [data, src.usage ? err : 'no usage data']
}

export function Usage({ account }: { account?: string }) {
  const [data, err] = useUsage()
  const now = useSource().now()
  if (!data) {
    return (
      <div className="desk usage">
        <div className="uhead"><h1>Usage{account ? ` · ${account}` : ''}</h1></div>
        <p className="empty">{err || 'Loading usage…'}</p>
      </div>
    )
  }
  return <UsageView data={data} account={account} now={now} />
}

const pct = (v: number) => `${Math.round(v)}%`

// The picked account: the route's, else the active one, else the first.
export function pickAccount(limits: UsageLimit[], account?: string): UsageLimit | undefined {
  return limits.find((l) => l.account === account) ?? limits.find((l) => l.active) ?? limits[0]
}

export function UsageView({ data, account, now }: { data: UsageData; account?: string; now: Date }) {
  const laptop = useMedia(LAPTOP)
  const q = useQuery()
  const [period, setPeriod] = useState<Period>((['D', 'W', 'M'] as const).find((p) => p === q.get('period')) ?? 'W')
  const lim = pickAccount(data.limits, account)
  const days = data.cost.data?.by_day ?? []
  const t = tiles(data, period, now)
  const p = { data, lim, now, period, setPeriod }
  const stale = data.cost.stale && data.cost.error && <div className="empty c-f">Cost audit: {data.cost.error}</div>
  const weekCost = cut(days, 'W').reduce((s, d) => s + d.cost_usd, 0)
  const sum = (
    <span className="sum">week <span className="c-n">{usd(weekCost)}</span> · {data.limits.length} account{data.limits.length === 1 ? '' : 's'}{lim?.reset_5h_at ? ` · 5h resets ${time(new Date(lim.reset_5h_at))}` : ''}</span>
  )
  const picker = data.limits.length > 0 && (
    <nav className="segs" aria-label="Accounts">
      {data.limits.map((l) => (
        <Link key={l.account} href={usagePath(l.account)} className={`seg${l === lim ? ' on' : ''}`} aria-current={l === lim ? 'page' : undefined}>{l.account}</Link>
      ))}
    </nav>
  )
  if (!laptop) {
    return (
      <div className="desk usage u-ph">
        <div className="uhead"><h1>Usage</h1><span className="sum">{periodLabel(period)} <span className="c-n">{usd(t.cost)}</span></span></div>
        {picker}
        {stale}
        {lim ? <><WindowPane lim={lim} now={now} days={days} phone /><WeekPane lim={lim} now={now} days={days} /></> : <div className="empty">No account limits yet.</div>}
        <Analysis {...p} phone />
        <Tiles t={t} period={period} phone />
      </div>
    )
  }
  return (
    <div className="desk usage u-lp">
      <div className="uhead"><h1>Usage</h1>{picker}{sum}</div>
      {stale}
      <div className="ugrid">
        <Deck n={1} title="Accounts">
          <div className="ulist">
            {data.limits.map((l) => <AccountRow key={l.account} l={l} on={l === lim} days={days} now={now} />)}
            {!data.limits.length && <div className="empty">No account limits yet.</div>}
          </div>
          <div className="ufill" />
          <Tiles t={t} period={period} />
        </Deck>
        <div className="uright">
          {lim ? <WindowPane lim={lim} now={now} days={days} /> : <Deck n={2} title="Window"><div className="empty">No account limits yet.</div></Deck>}
          <Analysis {...p} />
        </div>
      </div>
    </div>
  )
}

function AccountRow({ l, on, days, now }: { l: UsageLimit; on: boolean; days: Days; now: Date }) {
  const wk = weekOf(l, days, now)
  const vals = wk ? wk.points.map(([, v]) => v) : []
  return (
    <Link href={usagePath(l.account)} className={`urow${on ? ' on' : ''}`} aria-current={on ? 'page' : undefined}>
      <span className="ti"><b>{l.account}</b><span className="sub">5h {pct(l.remaining_5h_pct)} left · wk {pct(l.remaining_week_pct)}</span></span>
      <span className="usw">
        <Spark values={vals} floorLeft={floorOf(l)} label={`${l.account}: week used ${vals.map((v) => pct(v ?? 0)).join(', ')}`} />
        <span className="dim">{typeof l.burn_5h_pct_per_h === 'number' && l.burn_5h_pct_per_h > 0 ? `${l.burn_5h_pct_per_h}%/h` : 'idle'}</span>
      </span>
    </Link>
  )
}

function Facts({ lim, phone }: { lim: UsageLimit; phone?: boolean }) {
  const burn = typeof lim.burn_5h_pct_per_h === 'number' ? `${lim.burn_5h_pct_per_h}%/h` : '–'
  const empty = lim.empty_at ? <span className="c-n">{time(new Date(lim.empty_at))}</span> : <span className="dim">–</span>
  const f: Array<[string, ReactNode]> = phone
    ? [['Left', pct(lim.remaining_5h_pct)], ['Burn', burn], ['Empty', empty]]
    : [
        ['5h left', pct(lim.remaining_5h_pct)], ['Burn', burn], ['Empty at', empty],
        ['Resets', lim.reset_5h_at ? time(new Date(lim.reset_5h_at)) : '–'], ['Week left', pct(lim.remaining_week_pct)],
        ['Week resets', lim.reset_week_at ? weekReset(new Date(lim.reset_week_at)) : '–'],
      ]
  return <dl className="facts">{f.map(([k, v]) => <div key={k}><dt>{k}</dt><dd>{v}</dd></div>)}</dl>
}

const weekReset = (d: Date) => `${d.toLocaleDateString('en-GB', { weekday: 'short' })} ${time(d)}`

type Days = NonNullable<UsageData['cost']['data']>['by_day']

function FiveHour({ lim, now, w, h, short }: { lim: UsageLimit; now: Date; w: number; h: number; short?: boolean }) {
  const win = window5h(lim, now)
  if (!win) return <div className="empty">No 5h window yet.</div>
  const ticks = hourTicks(win.from, win.to).map(([d, t]) => [d, short ? `${t.slice(0, 2)}h` : t] as [Date, string])
  return <WindowChart w={w} h={h} title={`${lim.account} 5h window`} from={win.from} to={win.to} now={now} points={[[win.from, 0], [now, win.used]]} used={win.used} proj={win.proj} floorLeft={floorOf(lim)} ticks={ticks} fit={short} />
}

function WeekChart({ lim, now, days, w, h, short }: { lim: UsageLimit; now: Date; days: Days; w: number; h: number; short?: boolean }) {
  const wk = weekOf(lim, days, now)
  if (!wk) return <div className="empty">No week window yet.</div>
  return <WindowChart w={w} h={h} title={`${lim.account} week`} from={wk.from} to={wk.to} now={now} points={wk.points} used={wk.used} proj={wk.proj} floorLeft={floorOf(lim)} ticks={dayTicks(wk.from, wk.to, short)} fit={short} />
}

const WIN_LEGEND: Array<[string, string, ('box' | 'dash' | 'floor')?]> = [['used', 'var(--c1)'], ['projected', '', 'dash'], ['floor', '', 'floor']]

function WindowPane({ lim, now, days, phone }: { lim: UsageLimit; now: Date; days: Days; phone?: boolean }) {
  if (phone) {
    return (
      <Deck n={2} title="5h window" on>
        <Facts lim={lim} phone />
        <FiveHour lim={lim} now={now} w={330} h={170} short />
        <Legend items={WIN_LEGEND} />
      </Deck>
    )
  }
  return (
    <Deck n={2} title={`${lim.account} · 5h window`} on aside={<span className="dim">used to now · projection at the current burn</span>}>
      <Facts lim={lim} />
      <div className="uwin">
        <FiveHour lim={lim} now={now} w={530} h={190} />
        <WeekChart lim={lim} now={now} days={days} w={530} h={190} />
      </div>
      <Legend items={WIN_LEGEND} note={`x: local time · y: % of the window used · ${localDay(now) === localDay(new Date()) ? 'live' : 'as of'} ${time(now)}`} />
    </Deck>
  )
}

function WeekPane({ lim, now, days }: { lim: UsageLimit; now: Date; days: Days }) {
  return (
    <Deck n={3} title="Week">
      <WeekChart lim={lim} now={now} days={days} w={330} h={130} short />
    </Deck>
  )
}

function Tiles({ t, period, phone }: { t: ReturnType<typeof tiles>; period: Period; phone?: boolean }) {
  const tile = (k: string, v: string, s: string) => <div key={k} className="utile"><small>{k}</small><b>{v}</b><span>{s}</span></div>
  return (
    <div className={`utiles${phone ? ' uph' : ''}`}>
      {tile(periodLabel(period), usd(t.cost), `${t.turns.toLocaleString('en-GB')} turns`)}
      {tile('Cache hit', t.cache === null ? '–' : `${t.cache}%`, 'of context read')}
      {tile('Cache saved', usd(t.saved), `${t.savedFrom} session${t.savedFrom === 1 ? '' : 's'}`)}
      {tile('Week forecast', usd(t.forecast), `by Sunday · ${usd(t.weekSoFar)} so far`)}
      {tile('Hub start ctx', t.startCtx === null ? '–' : tok(t.startCtx), `median · ${t.startSessions} sessions`)}
    </div>
  )
}
