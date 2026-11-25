import type { CostDay, DayLimits, Group, UsageData, UsageLimit, UsageSession } from './usage'

// Usage page view model (item_227): the period cut, the per-day series in three units,
// the 5h window cells, the tiles and the session tables. Pure, so the page stays thin.

export type Period = 'D' | 'W' | 'M'
export type Unit = 'usd' | 'tok' | 'lim'
export type Tab = 'day' | 'project' | 'model' | 'kind' | 'session' | 'runs'
export const PERIODS: Array<[Period, number, string]> = [['D', 1, 'today'], ['W', 7, '7 days'], ['M', 30, '30 days']]
export const TABS: Array<[Tab, string, string]> = [
  ['day', 'Per day', 'Day'], ['project', 'Per project', 'Project'], ['model', 'Per model', 'Model'],
  ['kind', 'Per kind', 'Kind'], ['session', 'Per session', 'Session'], ['runs', 'Top runs', 'Runs'],
]
export const UNITS: Array<[Unit, string]> = [['usd', '$'], ['tok', 'tokens'], ['lim', 'limits']]
// Dispatch skips an account under this much room left (accounts.yaml min_room). A limit row
// that carries its own min_room wins.
export const FLOOR_LEFT = 10
export const floorOf = (l?: UsageLimit) => l?.min_room ?? FLOOR_LEFT

export const daysOf = (p: Period) => PERIODS.find((x) => x[0] === p)![1]
export const periodLabel = (p: Period) => PERIODS.find((x) => x[0] === p)![2]

// The last n days of the series, oldest first.
export function cut(days: CostDay[], p: Period): CostDay[] {
  return days.slice(-daysOf(p))
}

export const usd = (v: number) => (v === 0 ? '$0' : v >= 100 ? `$${Math.round(v)}` : v >= 10 ? `$${v.toFixed(0)}` : `$${v.toFixed(v >= 1 ? 1 : 2)}`)
export function tok(v: number): string {
  if (v >= 1e9) return `${(v / 1e9).toFixed(1)}B`
  if (v >= 1e6) return `${(v / 1e6).toFixed(v >= 1e8 ? 0 : 1)}M`
  if (v >= 1e3) return `${Math.round(v / 1e3)}k`
  return String(Math.round(v))
}

// Cache hit: the share of the context read from the cache.
export function cachePct(t: { cache_read_tokens?: number; cache_write_tokens?: number; input_tokens?: number }): number | null {
  const all = (t.cache_read_tokens ?? 0) + (t.cache_write_tokens ?? 0) + (t.input_tokens ?? 0)
  return all ? Math.round((100 * (t.cache_read_tokens ?? 0)) / all) : null
}

// "3 full + 10%": the 5h windows of a day, full ones counted, the rest summed.
export function cellsLabel(w: number[]): string {
  if (!w.length) return 'none'
  const full = w.filter((x) => x >= 100).length
  const rest = Math.round(w.filter((x) => x < 100).reduce((s, x) => s + x, 0))
  if (!full) return `${rest}%`
  return rest ? `${full} full + ${rest}%` : `${full} full`
}

export const dayLimits = (d: CostDay, account?: string): DayLimits | null | undefined =>
  d.limits === null ? null : account ? d.limits[account] : undefined

const DOW = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
export const dow = (date: string) => DOW[new Date(`${date}T12:00:00`).getDay()]
export const dayShort = (date: string) => `${dow(date)} ${date.slice(8, 10)}/${date.slice(5, 7)}`

// The hover/tap line of one day: "$84 · 232M tok · 98% cache · week −20% · 5h: 3 full + 10%".
export function daySummary(d: CostDay, account?: string): string {
  const parts = [usd(d.cost_usd), `${tok(d.tokens)} tok`]
  const c = cachePct(d)
  if (c !== null) parts.push(`${c}% cache`)
  const l = dayLimits(d, account)
  if (l === null) parts.push('limits: no data')
  else if (l) parts.push(`week −${Math.round(l.week_used_pct)}%`, `5h: ${cellsLabel(l.windows_5h)}`)
  return parts.join(' · ')
}

// Accounts seen in the limits of the period, in a stable order.
export function limitAccounts(days: CostDay[], limits: UsageLimit[]): string[] {
  const seen = new Set(limits.map((l) => l.account))
  for (const d of days) for (const a of Object.keys(d.limits ?? {})) seen.add(a)
  return [...seen]
}

// Model order: biggest cost first over the period.
export function modelsOf(days: CostDay[]): string[] {
  return sumGroups(days, (d) => d.by_model).map(([k]) => k)
}

// A model key without the version suffix ("opus-4-6" → "opus").
export const modelName = (m: string) => m.replace(/^claude-/, '').replace(/-\d.*$/, '')

export interface Stack { key: string; label: string; value: number }
export interface Bar { date: string; total: number; stacks: Stack[]; none?: boolean }

// The PER DAY bars in one unit.
export function bars(days: CostDay[], unit: Unit, accounts: string[]): Bar[] {
  const models = modelsOf(days)
  return days.map((d) => {
    if (unit === 'usd') {
      const stacks = models.map((m) => ({ key: m, label: modelName(m), value: d.by_model[m]?.cost_usd ?? 0 }))
      return { date: d.date, total: d.cost_usd, stacks }
    }
    if (unit === 'tok') {
      const stacks = [
        { key: 'read', label: 'cache read', value: d.cache_read_tokens },
        { key: 'write', label: 'cache write', value: d.cache_write_tokens },
        { key: 'io', label: 'input + output', value: d.input_tokens + d.output_tokens },
      ]
      return { date: d.date, total: stacks.reduce((s, x) => s + x.value, 0), stacks }
    }
    if (!d.limits) return { date: d.date, total: 0, stacks: [], none: true }
    const stacks = accounts.map((a) => ({ key: a, label: a, value: d.limits![a]?.week_used_pct ?? 0 }))
    return { date: d.date, total: stacks.reduce((s, x) => s + x.value, 0), stacks }
  })
}

// Groups summed over the period, biggest cost first.
export function sumGroups(days: CostDay[], pick: (d: CostDay) => Record<string, Group>): Array<[string, Group]> {
  const out: Record<string, Group> = {}
  for (const d of days) {
    for (const [k, g] of Object.entries(pick(d))) {
      const s = (out[k] ??= { cost_usd: 0, tokens: 0, turns: 0 })
      s.cost_usd += g.cost_usd
      s.tokens += g.tokens
      s.turns += g.turns
    }
  }
  return Object.entries(out).sort((a, b) => b[1].cost_usd - a[1].cost_usd)
}

// Sessions that ran in the period (a day of their by_day, else their start day).
export function sessionsIn(sessions: UsageSession[], days: CostDay[]): UsageSession[] {
  const keep = new Set(days.map((d) => d.date))
  return sessions.filter((s) => (s.by_day ? Object.keys(s.by_day).some((k) => keep.has(k)) : !!s.start && keep.has(localDay(new Date(s.start)))))
}

export const localDay = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`

export interface Tiles { cost: number; turns: number; cache: number | null; saved: number; savedFrom: number; forecast: number; weekSoFar: number; startCtx: number | null; startSessions: number }

// The left tiles: period cost, cache hit, cache savings, cost forecast to Sunday, hub start context.
export function tiles(u: UsageData, p: Period, now: Date): Tiles {
  const all = u.cost.data?.by_day ?? []
  const days = cut(all, p)
  const sum = (f: (d: CostDay) => number) => days.reduce((s, d) => s + f(d), 0)
  const sess = sessionsIn(u.sessions, days)
  // Week from Monday: spent so far, then the 7-day average per day for the days left.
  const back = (now.getDay() + 6) % 7
  const monday = localDay(new Date(now.getFullYear(), now.getMonth(), now.getDate() - back))
  const weekSoFar = all.filter((d) => d.date >= monday).reduce((s, d) => s + d.cost_usd, 0)
  const last7 = all.slice(-7)
  const avg = last7.length ? last7.reduce((s, d) => s + d.cost_usd, 0) / last7.length : 0
  return {
    cost: sum((d) => d.cost_usd),
    turns: sum((d) => d.turns),
    cache: cachePct({ cache_read_tokens: sum((d) => d.cache_read_tokens), cache_write_tokens: sum((d) => d.cache_write_tokens), input_tokens: sum((d) => d.input_tokens) }),
    saved: sess.reduce((s, x) => s + (x.cache_saved_usd ?? 0), 0),
    savedFrom: sess.length,
    forecast: weekSoFar + avg * (6 - back),
    weekSoFar,
    startCtx: u.cost.data?.start_context.median ?? null,
    startSessions: u.cost.data?.start_context.sessions ?? 0,
  }
}

// The 5h window of an account: start, reset, used now and used at reset at the current burn.
export interface Window { from: Date; to: Date; used: number; proj: number | null }
export function window5h(l: UsageLimit, now: Date): Window | null {
  if (!l.reset_5h_at) return null
  const to = new Date(l.reset_5h_at)
  const from = new Date(to.getTime() - 5 * 3600e3)
  const used = 100 - l.remaining_5h_pct
  const burn = l.burn_5h_pct_per_h
  const left = Math.max(0, (to.getTime() - now.getTime()) / 3600e3)
  return { from, to, used, proj: typeof burn === 'number' ? Math.min(100, Math.round(used + burn * left)) : null }
}

// The week of an account: the used % at the end of each sampled day since the week start,
// then now; the projection keeps the average pace of the week so far.
export interface Week { from: Date; to: Date; used: number; points: Array<[Date, number]>; proj: number | null }
export function weekOf(l: UsageLimit, days: CostDay[], now: Date): Week | null {
  if (!l.reset_week_at) return null
  const to = new Date(l.reset_week_at)
  const from = new Date(to.getTime() - 7 * 864e5)
  const used = 100 - l.remaining_week_pct
  const points: Array<[Date, number]> = [[from, 0]]
  let acc = 0
  const today = localDay(now)
  for (const d of days) {
    const end = new Date(`${d.date}T23:59:59`)
    const lim = d.limits?.[l.account]
    if (!lim || d.date === today || end <= from || end >= now) continue
    acc += lim.week_used_pct
    points.push([end, Math.min(acc, used)])
  }
  points.push([now, used])
  const spent = (now.getTime() - from.getTime()) / (to.getTime() - from.getTime())
  const proj = spent > 0.02 ? Math.min(100, Math.round(used / spent)) : null
  return { from, to, used, points, proj }
}

// Sessions for the PER SESSION table (newest first) and TOP RUNS (dearest first).
export const bySession = (s: UsageSession[]) => [...s].sort((a, b) => (b.start ?? '').localeCompare(a.start ?? ''))
export const topRuns = (s: UsageSession[], n = 10) => [...s].sort((a, b) => (b.cost_usd ?? 0) - (a.cost_usd ?? 0)).slice(0, n)

// Axis ticks: 5 lines over a round maximum that splits into four even steps (12/24/36/48).
export function niceMax(v: number): number {
  if (v <= 0) return 1
  const p = 10 ** Math.floor(Math.log10(v))
  for (const m of [1, 1.2, 1.6, 2, 2.4, 3.2, 4, 4.8, 6, 8, 10]) if (v <= m * p) return m * p
  return 10 * p
}
