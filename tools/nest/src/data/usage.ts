// Usage page (item_227) over GET /api/usage: the cost audit (crow-cost-audit --json), the
// account limits and the session list of crow-usage --json.

import type { Limit } from './home'

export interface TokenKinds {
  input_tokens: number
  output_tokens: number
  cache_read_tokens: number
  cache_write_tokens: number
}

export interface Group { cost_usd: number; tokens: number; turns: number }
export type ModelGroup = Group & TokenKinds

// One day's limits per account: weekly % used that day and the used % of each 5h window.
export interface DayLimits { week_used_pct: number; windows_5h: number[] }

export interface CostDay extends TokenKinds {
  date: string
  cost_usd: number
  turns: number
  tokens: number
  by_model: Record<string, ModelGroup>
  by_account: Record<string, Group>
  by_project: Record<string, Group>
  by_kind: Record<string, Group>
  // Null before the limit samples started.
  limits: Record<string, DayLimits> | null
}

export interface Share { cost_usd: number; turns: number; context_tokens: number; share_pct: number }

export interface CostAudit {
  days: number
  since: string
  turns: number
  cost_usd: number
  context_tokens: number
  by_model: Record<string, Share>
  by_project: Record<string, Share>
  by_account: Record<string, Share>
  buckets: Record<string, Share>
  start_context: { sessions: number; median: number | null }
  by_day: CostDay[]
}

export interface UsageLimit extends Limit {
  in_policy?: boolean
  burn_5h_pct_per_h?: number | null
  empty_at?: string | null
}

export interface UsageSession extends Partial<TokenKinds> {
  id: string
  account: string
  project?: string
  title?: string
  kind?: 'run' | 'interactive'
  models?: string[]
  start?: string
  end?: string
  active?: boolean
  tokens?: number
  cache_pct?: number
  cost_usd?: number
  cache_saved_usd?: number
  by_day?: Record<string, Partial<TokenKinds> & { cost_usd?: number; tokens?: number }>
}

export interface UsageData {
  version: string
  generated_at: string
  cost: { data: CostAudit | null; stale: boolean; error?: string; updated_at: string | null }
  limits: UsageLimit[]
  sessions: UsageSession[]
  stale?: boolean
  error?: string
}

const isAudit = (d: unknown): d is CostAudit => !!d && typeof d === 'object' && Array.isArray((d as CostAudit).by_day)

// The wire body with every list defaulted: a cold core answers an empty, stale audit.
export function usageBody(b: Partial<Omit<UsageData, 'cost'>> & { cost?: { data?: unknown; stale?: boolean; error?: string; updated_at?: string | null } }): UsageData {
  const c = b.cost ?? {}
  return {
    version: b.version ?? '',
    generated_at: b.generated_at ?? '',
    cost: { data: isAudit(c.data) ? c.data : null, stale: c.stale ?? true, error: c.error, updated_at: c.updated_at ?? null },
    limits: Array.isArray(b.limits) ? b.limits.filter((l) => l && typeof l.account === 'string') : [],
    sessions: Array.isArray(b.sessions) ? b.sessions.filter((s) => s && typeof s.id === 'string') : [],
    stale: b.stale,
    error: b.error,
  }
}
