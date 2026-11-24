// Ops timeline (item_190) over GET /api/state/ops and GET /api/runs/{id} (item_257).

export type SpanState = 'working' | 'idle' | 'waiting' | 'stuck' | 'failed' | 'done'
export interface OpsSpan { state: SpanState; start: string; end: string }
export interface OpsLane {
  id: string
  day: string
  project: string
  account: string
  kind: 'interactive' | 'run' | 'mission' | 'routine'
  title: string
  model?: string
  state: string
  spans: OpsSpan[]
  keys: string[]
  active_s: number
  wait_s: number
  cost_usd: number
  tokens: number
  commits: number
  runs: number
  // The lane's working directory when known (item_273).
  cwd?: string
}
export interface OpsQueued { kind: 'mission' | 'routine' | 'afk'; id: string; project?: string; title: string; next?: string; paused?: boolean; over_cap?: boolean }
export interface OpsTotal { by: 'project' | 'kind'; name: string; lanes: number; active_s: number; wait_s: number; stuck: number; failed: number; cost_usd: number; tokens: number }
export interface Ops { day?: string; week?: string; days: string[]; lanes: OpsLane[]; queue: OpsQueued[]; totals: OpsTotal[] }

export interface RunReport {
  run: string
  source: string
  project: string
  label: string
  account: string
  status: string
  started: string | null
  ended: string | null
  cost_usd: number | null
  commits: string[]
  branch: string | null
  worktree: string | null
  item: string | null
  mission: string | null
  report: string
  // Older cores omit it; null for live, mission and non-git runs.
  diffstat?: DiffStat | null
}

export interface FileStat {
  path: string
  insertions: number
  deletions: number
}

// git diff --numstat <start>...HEAD stored when a crow-runs run ends (item_269).
export interface DiffStat {
  base: string
  files: number
  insertions: number
  deletions: number
  per_file: FileStat[]
  truncated: boolean
}

// GET /api/changes?cwd= (item_273): the git state of a live lane's checkout.
export interface Commit { sha: string; subject: string; time: string }
export interface Changes {
  cwd: string
  branch: string
  base: string
  ahead: number
  behind: number
  uncommitted: number
  diffstat: DiffStat
  commits: Commit[]
  more: boolean
}

// The last two segments of a path: "cr-0w.wt/item_262".
export const shortPath = (p: string) => p.replace(/\/+$/, '').split('/').slice(-2).join('/')

// A run report commit line "sha subject" as a Commit without a time.
export function oneline(l: string): Commit {
  const i = l.indexOf(' ')
  return i < 0 ? { sha: l, subject: '', time: '' } : { sha: l.slice(0, i), subject: l.slice(i + 1), time: '' }
}

// Lane states that are live right now (lanes package glyph states).
const LIVE = new Set(['working', 'run', 'needs', 'idle'])
export const isLive = (l: OpsLane) => LIVE.has(l.state)

// The owner reads the axis in Paris time, wherever the page runs.
export const TZ = 'Europe/Paris'
const hm = new Intl.DateTimeFormat('en-GB', { timeZone: TZ, hour: '2-digit', minute: '2-digit', hourCycle: 'h23' })
export const clock = (t: Date | string) => hm.format(new Date(t))

// Hours since Paris midnight of the instant (DST-safe through the formatter).
export function parisHour(t: Date | string): number {
  const [h, m] = clock(t).split(':').map(Number)
  return h + m / 60
}

// Today in Paris as YYYY-MM-DD.
export const parisDay = (t: Date) => new Intl.DateTimeFormat('en-CA', { timeZone: TZ }).format(t)

export function shiftDay(day: string, n: number): string {
  const d = new Date(`${day}T12:00:00Z`)
  d.setUTCDate(d.getUTCDate() + n)
  return d.toISOString().slice(0, 10)
}

// ISO week (YYYY-Www) holding day.
export function isoWeek(day: string): string {
  const d = new Date(`${day}T12:00:00Z`)
  const wd = (d.getUTCDay() + 6) % 7
  d.setUTCDate(d.getUTCDate() - wd + 3)
  const y = d.getUTCFullYear()
  const first = new Date(Date.UTC(y, 0, 4))
  const n = 1 + Math.round((d.getTime() - first.getTime()) / 864e5 / 7 + ((first.getUTCDay() + 6) % 7 - 3) / 7)
  return `${y}-W${String(n).padStart(2, '0')}`
}

// The axis window in hours. 'full' is the whole day. 'work' is a fixed 08:00-20:00 so days
// compare at a glance, widened to whole hours when a span or the now line falls outside.
export type Scale = 'work' | 'full'
export function axis(ls: OpsLane[], today: boolean, nowH: number, scale: Scale = 'work'): [number, number] {
  if (scale === 'full') return [0, 24]
  let lo = 8, hi = 20
  for (const l of ls) for (const s of l.spans) {
    lo = Math.min(lo, parisHour(s.start))
    hi = Math.max(hi, parisHour(s.end) || 24)
  }
  if (today) {
    lo = Math.min(lo, nowH)
    hi = Math.max(hi, nowH + 0.5)
  }
  return [Math.floor(lo), Math.min(24, Math.ceil(hi))]
}

// Lanes by project, the project with the most active time first; inside, live lanes first.
export function groups(ls: OpsLane[]): Array<{ project: string; lanes: OpsLane[] }> {
  const by = new Map<string, OpsLane[]>()
  for (const l of ls) by.set(l.project, [...(by.get(l.project) ?? []), l])
  const sum = (x: OpsLane[]) => x.reduce((a, l) => a + l.active_s, 0)
  return [...by.entries()]
    .map(([project, lanes]) => ({ project, lanes: [...lanes].sort((a, b) => Number(isLive(b)) - Number(isLive(a)) || first(a) - first(b)) }))
    .sort((a, b) => sum(b.lanes) - sum(a.lanes))
}

const first = (l: OpsLane) => (l.spans[0] ? Date.parse(l.spans[0].start) : Infinity)

// The run id of a lane: a crow-runs key (32 hex, no dashes).
export const runId = (l: OpsLane) => l.keys.find((k) => /^[0-9a-f]{32}$/.test(k))

export interface DayStats { lanes: number; active_s: number; wait_s: number; cost: number; failed: number; sessions: number; runs: number; done: number }
export function stats(ls: OpsLane[]): DayStats {
  const s: DayStats = { lanes: ls.length, active_s: 0, wait_s: 0, cost: 0, failed: 0, sessions: 0, runs: 0, done: 0 }
  for (const l of ls) {
    s.active_s += l.active_s
    s.wait_s += l.wait_s
    s.cost += l.cost_usd
    if (l.state === 'failed' || l.spans.some((x) => x.state === 'failed')) s.failed++
    if (l.kind === 'interactive') s.sessions++
    else s.runs++
    if (l.spans.at(-1)?.state === 'done' && l.kind !== 'interactive') s.done++
  }
  return s
}

// Busy stretches of a day: the union of working spans, in Paris hours.
export function busy(ls: OpsLane[]): Array<[number, number]> {
  const xs = ls.flatMap((l) => l.spans.filter((s) => s.state !== 'idle').map((s) => [parisHour(s.start), parisHour(s.end) || 24] as [number, number]))
  xs.sort((a, b) => a[0] - b[0])
  const out: Array<[number, number]> = []
  for (const x of xs) {
    const last = out.at(-1)
    if (last && x[0] <= last[1]) last[1] = Math.max(last[1], x[1])
    else out.push([...x])
  }
  return out
}

export const hours = (s: number) => (s / 3600).toFixed(1)
export const usd = (n: number) => `$${n < 10 ? n.toFixed(2) : n.toFixed(0)}`
