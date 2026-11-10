import type { SectionBody } from './types'

// Live sessions, dispatch runs and Today mail (item_170). Shapes are the parts of
// the script sections Nest reads (logics/specs/script_contract); actions are the
// core POST routes of tools/crow-core/internal/api/live.go.

// One agent of `crow-runs json`: an interactive session or a background run.
export interface Agent {
  kind: 'session' | 'run'
  pid: number
  account: string
  project: string
  cwd: string
  session_id?: string | null
  run_id?: string | null
  label?: string | null
  title?: string | null
  state?: string
  started_at?: string | null
  last_action?: string | null
  last_prompt?: string | null
}

export interface Dispatch {
  slots: { busy: number; max: number }
  agents: Agent[]
}

// A resumable conversation of `crow-sessions json`.
export interface Resumable {
  conversation_id: string
  account: string
  cwd: string
  project: string
  branch?: string | null
  limit?: string | null
  resets_at?: string | null
  last_prompt?: string | null
}

export interface Mail {
  id: string
  subject: string
  from?: string
}

export interface Account {
  name: string
  policy: string
}

// The reply of every live action: the event core wrote to actions.jsonl.
export interface LiveReply {
  ts: string
  actor: string
  surface: string
  action: string
  target: string
  text?: string
  account?: string
  project?: string
  ok: boolean
  code: number
  message: string
  run_id?: string
}

export interface NewSession { project: string; account?: string; prompt?: string }
export interface NewRun { project: string; label: string; prompt: string; model?: string }

// Each call resolves with the reply or rejects with an ApiError (4xx/5xx).
export interface Actions {
  send(pid: number, text: string): Promise<LiveReply>
  interrupt(pid: number): Promise<LiveReply>
  stop(pid: number): Promise<LiveReply>
  resume(id: string, account?: string): Promise<LiveReply>
  start(s: NewSession): Promise<LiveReply>
  dispatch(r: NewRun): Promise<LiveReply>
  stopRun(runId: string): Promise<LiveReply>
  mail(id: string, verb: 'done' | 'snooze'): Promise<LiveReply>
}

const obj = (b: SectionBody | undefined) => (b?.data && typeof b.data === 'object' ? (b.data as Record<string, unknown>) : {})
const list = <T>(v: unknown) => (Array.isArray(v) ? (v as T[]) : [])

export function dispatchOf(b: SectionBody | undefined): Dispatch {
  const d = obj(b)
  const slots = (d.slots as Dispatch['slots'] | undefined) ?? { busy: 0, max: 0 }
  return { slots, agents: list<Agent>(d.agents) }
}

export const resumableOf = (b: SectionBody | undefined) => list<Resumable>(obj(b).sessions)

// Important and flagged mail, flagged first, each id once.
export function mailOf(b: SectionBody | undefined): Mail[] {
  const d = obj(b)
  const seen = new Set<string>()
  return [...list<Mail>(d.flagged), ...list<Mail>(d.important)].filter((m) => !seen.has(m.id) && seen.add(m.id))
}

export const projectsOf = (b: SectionBody | undefined) => list<{ name: string }>(obj(b).projects).map((p) => p.name)

// Accounts a new session may use; "never" accounts are refused by core anyway.
export const accountsOf = (b: SectionBody | undefined) => list<Account>(obj(b).accounts).filter((a) => a.policy !== 'never')
