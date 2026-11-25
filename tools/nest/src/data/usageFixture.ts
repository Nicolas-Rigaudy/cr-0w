import type { CostDay, Group, ModelGroup, Share, UsageData, UsageSession } from './usage'

// Usage fixture (item_227): 35 days to the fixture clock for two accounts and two projects.
// Limit samples start five days before it, like the live data started on its first sample day.

const ACCTS = ['acct1', 'acct2']
const PROJECTS = ['api', 'shell']
const MODELS = ['opus-4-6', 'sonnet-4-5']
const KINDS = ['interactive', 'headless']
export const USAGE_LIMITS_FROM = '2026-11-04'

// Deterministic noise so tests and shots stay put.
function rng(seed: number) {
  let s = seed
  return () => ((s = (s * 1103515245 + 12345) % 2147483648) / 2147483648)
}

const day = (end: Date, back: number) => {
  const d = new Date(end.getFullYear(), end.getMonth(), end.getDate() - back)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const add = (into: Record<string, Group>, key: string, cost: number, tokens: number, turns: number) => {
  const g = (into[key] ??= { cost_usd: 0, tokens: 0, turns: 0 })
  g.cost_usd = Math.round((g.cost_usd + cost) * 100) / 100
  g.tokens += tokens
  g.turns += turns
}

export function usageFixture(now: Date): UsageData {
  const r = rng(227)
  const days: CostDay[] = []
  for (let back = 34; back >= 0; back--) {
    const date = day(now, back)
    const weekend = [0, 6].includes(new Date(`${date}T12:00:00`).getDay())
    const scale = (weekend ? 0.25 : 1) * (back === 0 ? 0.7 : 1) * (0.6 + r() * 0.8)
    const d: CostDay = { date, cost_usd: 0, turns: 0, tokens: 0, input_tokens: 0, output_tokens: 0, cache_read_tokens: 0, cache_write_tokens: 0,
      by_model: {}, by_account: {}, by_project: {}, by_kind: {}, limits: null }
    MODELS.forEach((m, mi) => {
      const read = Math.round((mi ? 40e6 : 120e6) * scale)
      const write = Math.round(read * (0.03 + r() * 0.02))
      const io = Math.round(read * 0.004)
      const cost = Math.round(((mi ? 0.3 : 1.5) * read / 1e6 + (mi ? 3.75 : 18.75) * write / 1e6 + (mi ? 9 : 45) * io / 1e6) * 100) / 100
      const turns = Math.round(read / 150e3)
      const g: ModelGroup = { cost_usd: cost, tokens: read + write + io, turns, input_tokens: Math.round(io * 0.2), output_tokens: io - Math.round(io * 0.2), cache_read_tokens: read, cache_write_tokens: write }
      d.by_model[m] = g
      d.cost_usd = Math.round((d.cost_usd + cost) * 100) / 100
      d.turns += turns
      d.tokens += g.tokens
      d.input_tokens += g.input_tokens
      d.output_tokens += g.output_tokens
      d.cache_read_tokens += read
      d.cache_write_tokens += write
    })
    const split = (into: Record<string, Group>, keys: string[], first: number) => {
      add(into, keys[0], d.cost_usd * first, Math.round(d.tokens * first), Math.round(d.turns * first))
      add(into, keys[1], d.cost_usd - into[keys[0]].cost_usd, d.tokens - into[keys[0]].tokens, d.turns - into[keys[0]].turns)
    }
    split(d.by_account, ACCTS, 0.55 + r() * 0.2)
    split(d.by_project, PROJECTS, 0.5 + r() * 0.3)
    split(d.by_kind, KINDS, 0.7 + r() * 0.15)
    if (date >= USAGE_LIMITS_FROM) {
      d.limits = Object.fromEntries(ACCTS.map((a, ai) => {
        const share = d.by_account[a].cost_usd / Math.max(1, d.cost_usd)
        const windows = weekend ? [Math.round(30 + r() * 40)] : [100, 100, Math.round(r() * 90) + 10].slice(ai)
        return [a, { week_used_pct: Math.round(share * d.cost_usd / 9), windows_5h: windows }]
      }))
    }
    days.push(d)
  }
  const total = days.reduce((s, d) => s + d.cost_usd, 0)
  const share = (pick: (d: CostDay) => Record<string, Group>): Record<string, Share> => {
    const out: Record<string, Share> = {}
    for (const d of days) {
      for (const [k, g] of Object.entries(pick(d))) {
        const s = (out[k] ??= { cost_usd: 0, turns: 0, context_tokens: 0, share_pct: 0 })
        s.cost_usd = Math.round((s.cost_usd + g.cost_usd) * 100) / 100
        s.turns += g.turns
        s.context_tokens += g.tokens
      }
    }
    for (const s of Object.values(out)) s.share_pct = Math.round((1000 * s.cost_usd) / total) / 10
    return out
  }
  const at = (back: number, hm: string) => new Date(`${day(now, back)}T${hm}:00+01:00`).toISOString()
  const sessions: UsageSession[] = [
    ['s01', 'acct1', 'api', 'Refactor the auth middleware', 'interactive', 0, '09:12', '12:40', 31.2, 182e6],
    ['s02', 'acct1', 'shell', 'Nightly lint and test run', 'run', 0, '03:00', '03:48', 6.4, 41e6],
    ['s03', 'acct2', 'api', 'Review the billing endpoints', 'interactive', 0, '13:05', '16:08', 18.9, 120e6],
    ['s04', 'acct2', 'shell', 'Theme tokens for the prompt', 'run', 1, '10:20', '11:02', 4.1, 22e6],
    ['s05', 'acct1', 'api', 'Fix flaky queue tests', 'run', 1, '14:00', '15:30', 12.7, 76e6],
    ['s06', 'acct1', 'shell', 'Plan the plugin loader', 'interactive', 2, '09:40', '12:10', 22.5, 140e6],
    ['s07', 'acct2', 'api', 'Migrate the session store', 'run', 3, '08:30', '10:55', 15.8, 98e6],
    ['s08', 'acct1', 'api', 'Draft the release notes', 'interactive', 4, '16:00', '16:40', 3.2, 18e6],
    ['s09', 'acct2', 'shell', 'Benchmark the startup path', 'run', 5, '11:10', '12:00', 7.9, 51e6],
    ['s10', 'acct1', 'shell', 'Rewrite the completion engine', 'interactive', 6, '09:00', '13:30', 27.4, 160e6],
  ].map(([id, account, project, title, kind, back, from, to, cost, tokens]) => {
    const tk = tokens as number
    const read = Math.round(tk * 0.94)
    return {
      id: id as string, account: account as string, project: project as string, title: title as string, kind: kind as 'run' | 'interactive',
      models: kind === 'run' ? ['sonnet-4-5'] : ['opus-4-6', 'sonnet-4-5'], start: at(back as number, from as string), end: at(back as number, to as string),
      active: back === 0 && id === 's03', tokens: tk, cache_read_tokens: read, cache_write_tokens: Math.round(tk * 0.05), input_tokens: Math.round(tk * 0.002), output_tokens: Math.round(tk * 0.008),
      cache_pct: 94, cost_usd: cost as number, cache_saved_usd: Math.round((cost as number) * 3.1 * 10) / 10,
      by_day: { [day(now, back as number)]: { cost_usd: cost as number, tokens: tk } },
    }
  })
  const iso = (hm: string, back = 0) => new Date(`${day(now, -back)}T${hm}:00+01:00`).toISOString()
  return {
    version: 'fixture',
    generated_at: now.toISOString(),
    cost: {
      stale: false, updated_at: now.toISOString(),
      data: {
        days: 35, since: new Date(now.getTime() - 35 * 864e5).toISOString(), turns: days.reduce((s, d) => s + d.turns, 0),
        cost_usd: Math.round(total * 100) / 100, context_tokens: days.reduce((s, d) => s + d.tokens, 0),
        by_model: share((d) => d.by_model), by_project: share((d) => d.by_project), by_account: share((d) => d.by_account), buckets: {},
        start_context: { sessions: 14, median: 38_400 }, by_day: days,
      },
    },
    limits: [
      { account: 'acct1', remaining_5h_pct: 38, reset_5h_at: iso('18:00'), remaining_week_pct: 62, reset_week_at: iso('09:00', 4), active: true, in_policy: true, burn_5h_pct_per_h: 22, empty_at: iso('17:53') },
      { account: 'acct2', remaining_5h_pct: 91, reset_5h_at: iso('19:30'), remaining_week_pct: 30, reset_week_at: iso('09:00', 2), active: false, in_policy: true, burn_5h_pct_per_h: 3, empty_at: null },
    ],
    sessions,
  }
}
