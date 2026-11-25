import { sameDay } from '../format'
import type { SectionBody } from './types'

// The MRs page (item_225) over the crow-mrs section (logics/specs/script_contract/schemas/mrs.schema.json).
// Fields marked "not in crow-mrs yet" render when present and stay hidden otherwise.

export type Pipeline = 'passed' | 'failed' | 'running' | 'none'
export type MrRole = 'AUTHOR' | 'REVIEWER' | 'ASSIGNEE' | ''

export interface MrJob { name: string; state: Pipeline | 'skipped' | 'manual'; duration?: string }
export interface MrThread { path?: string; line?: number; author: string; at: string; body: string; resolved: boolean }

export interface Mr {
  project: string
  name: string
  scope: string
  forge: string
  repo: string
  number: number
  ref: string
  title: string
  author: string
  role: MrRole
  draft: boolean
  // The owner already approved it (reviewer rows only).
  approved: boolean
  pipeline: Pipeline
  updated: string
  created: string
  branch: string
  target: string
  url: string
  threads_open: number
  threads_resolved: number
  // Not in crow-mrs yet.
  approvals?: { given: number; required: number }
  reviewers?: string[]
  pipeline_id?: number
  jobs?: MrJob[]
  failure?: string
  threads?: MrThread[]
  merged_at?: string
}

export interface MrBranch { project: string; name?: string; scope: string; branch: string; ahead: number }

// One review ledger row (local/reviews/common.md and the project ledger). Not in crow-mrs yet.
export interface Pattern {
  pattern: string
  count: number
  // The last MR that hit it and when.
  last?: string
  at?: string
  // MR refs a review run matched against it.
  hits?: string[]
}

export interface MrsData {
  off: boolean
  stale: boolean
  error?: string
  review: Mr[]
  mine: Mr[]
  merged: Mr[]
  branches: MrBranch[]
  patterns: Record<string, Pattern[]>
  projects: string[]
}

const obj = (v: unknown) => (v && typeof v === 'object' && !Array.isArray(v) ? (v as Record<string, unknown>) : {})
const arr = <T>(v: unknown) => (Array.isArray(v) ? (v as T[]) : [])
const num = (v: unknown) => (typeof v === 'number' ? v : Number(v) || 0)
const str = (v: unknown) => (typeof v === 'string' ? v : v == null ? '' : String(v))
const PIPES: Pipeline[] = ['passed', 'failed', 'running', 'none']
const ROLES: MrRole[] = ['AUTHOR', 'REVIEWER', 'ASSIGNEE']

export function mrOf(v: unknown): Mr {
  const r = obj(v)
  const pipe = str(r.pipeline) as Pipeline
  const role = str(r.role).toUpperCase() as MrRole
  const forge = str(r.forge)
  const number = num(r.number)
  return {
    ...(r as Partial<Mr>),
    project: str(r.project), name: str(r.name) || str(r.project), scope: str(r.scope), forge, repo: str(r.repo), number,
    ref: str(r.ref) || `${forge === 'github' ? '#' : '!'}${number}`, title: str(r.title), author: str(r.author),
    role: ROLES.includes(role) ? role : '', draft: !!r.draft, approved: !!r.approved, pipeline: PIPES.includes(pipe) ? pipe : 'none',
    updated: str(r.updated), created: str(r.created), branch: str(r.branch), target: str(r.target), url: str(r.url),
    threads_open: num(r.threads_open), threads_resolved: num(r.threads_resolved),
  }
}

// Mine needs the owner when its pipeline failed or a thread is open.
export const needsYou = (m: Mr) => m.role === 'AUTHOR' && (m.pipeline === 'failed' || m.threads_open > 0)
const t = (iso: string) => Date.parse(iso) || 0

export function mrsOf(b: SectionBody | undefined, now: Date): MrsData {
  const d = obj(b?.data)
  const rows = arr<unknown>(d.mrs).map(mrOf)
  // To review: ready MRs the owner has not approved yet, the oldest wait first (unblock devs first).
  const review = rows.filter((m) => m.role !== 'AUTHOR' && !m.draft && !m.approved).sort((a, c) => t(a.updated) - t(c.updated))
  // Mine: what needs the owner first, then the oldest.
  const mine = rows.filter((m) => m.role === 'AUTHOR').sort((a, c) => Number(needsYou(c)) - Number(needsYou(a)) || t(a.updated) - t(c.updated))
  const merged = arr<unknown>(d.merged).map(mrOf).filter((m) => m.merged_at && sameDay(new Date(m.merged_at), now))
    .sort((a, c) => t(c.merged_at!) - t(a.merged_at!))
  const branches = arr<Record<string, unknown>>(d.branches).map((x) => ({ ...x, project: str(x.project), scope: str(x.scope), branch: str(x.branch), ahead: num(x.ahead) }) as MrBranch)
  const patterns: Record<string, Pattern[]> = {}
  for (const [k, v] of Object.entries(obj(d.patterns))) patterns[k] = arr<Pattern>(v)
  const projects = [...new Set([...rows, ...merged].map((m) => m.project))].sort()
  return { off: !!b?.off || !!d.off, stale: !!b?.stale && !!b?.error, error: b?.error, review, mine, merged, branches, patterns, projects }
}

// The ledger rows for one MR's project: the common ledger plus the project one, newest hit first.
export function patternsFor(d: MrsData, m: Mr): Pattern[] {
  const rows = [...(d.patterns[m.project] ?? []), ...(d.patterns.common ?? [])]
  return rows.map((p, i) => ({ p, i })).sort((a, c) => t(c.p.at ?? '') - t(a.p.at ?? '') || a.i - c.i).map((x) => x.p)
}
export const hitsOf = (ps: Pattern[], m: Mr) => ps.filter((p) => p.hits?.includes(m.ref))

export interface MrsCounts { review: number; needs: number; waiting: number; merged: number }
export function countsOf(d: MrsData): MrsCounts {
  const needs = d.mine.filter(needsYou).length
  return { review: d.review.length, needs, waiting: d.mine.length - needs, merged: d.merged.length }
}

// A project filter keeps the rows, branches and counts of that project.
export function onlyProject(d: MrsData, p: string): MrsData {
  if (!p) return d
  const f = (m: { project: string }) => m.project === p
  return { ...d, review: d.review.filter(f), mine: d.mine.filter(f), merged: d.merged.filter(f), branches: d.branches.filter(f) }
}

// Prompts of the row actions. The review and review-threads skills read the MR from its URL.
export const reviewPrompt = (m: Mr) => `/review ${m.url || `${m.project} ${m.ref}`}`
export const threadsPrompt = (m: Mr) => `/review-threads ${m.url || `${m.project} ${m.ref}`}`
export const reviewLabel = (m: Mr) => `review-${m.project}-${m.number}`
export const branchPrompt = (b: MrBranch) => `/mr ${b.branch}`
