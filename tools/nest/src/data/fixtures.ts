import { missions as coreMissions, inbox as coreInbox } from 'virtual:core-records'
import stateMissions from '../../fixtures/states/missions.json'
import stateInbox from '../../fixtures/states/inbox.json'
import stateAgenda from '../../fixtures/states/agenda.json'
import stateLive from '../../fixtures/states/live.json'
import { ApiError } from './api'
import type { Actions, Dispatch, LiveReply } from './live'
import type { DataSource } from './source'
import type { AgendaEvent, InboxItem, Mission, SectionBody, SectionName } from './types'

// Shared core-contract records plus Nest-local records that cover every UI state.
const MISSIONS = [...coreMissions, ...stateMissions] as Mission[]
const INBOX = [...coreInbox, ...stateInbox] as InboxItem[]
const AGENDA = stateAgenda as AgendaEvent[]
const NOW = new Date('2026-11-09T16:10:00+01:00')
// Sections without a state fixture stay empty until their screens bring one (item_171/172).
const NO_SECTION: SectionBody = { data: null, stale: false, updated_at: null }
type LiveData = Partial<Record<SectionName, unknown>>

// In-memory live actions over the fixture sections: they check targets like core
// (404 unknown, 409 full slots) and change the section so a refetch shows the result.
function fixtureActions(data: LiveData, bodies: Map<SectionName, SectionBody>, now: Date): Actions {
  const set = (name: SectionName, v: unknown) => {
    data[name] = v
    bodies.delete(name)
  }
  const dispatch = () => data.dispatch as Dispatch
  const reply = (action: string, target: string, extra: Partial<LiveReply> = {}): LiveReply =>
    ({ ts: now.toISOString(), actor: 'owner', surface: 'nest', action, target, ok: true, code: 0, message: 'done', ...extra })
  const agent = (pid: number) => {
    const a = dispatch().agents.find((x) => x.pid === pid)
    if (!a) throw new ApiError(404, `no live session with pid ${pid}`)
    return a
  }
  const drop = (pid: number) => {
    const d = dispatch()
    const a = agent(pid)
    set('dispatch', { ...d, slots: { ...d.slots, busy: d.slots.busy - (a.kind === 'run' ? 1 : 0) }, agents: d.agents.filter((x) => x !== a) })
    return a
  }
  return {
    async send(pid, text) {
      const a = agent(pid)
      if (!text.trim()) throw new ApiError(400, 'text is empty')
      return reply('send', String(pid), { text, account: a.account, project: a.project })
    },
    async interrupt(pid) {
      const a = agent(pid)
      return reply('interrupt', String(pid), { account: a.account, project: a.project })
    },
    async stop(pid) {
      const a = drop(pid)
      return reply('stop', String(pid), { account: a.account, project: a.project })
    },
    async stopRun(run) {
      const a = dispatch().agents.find((x) => x.run_id === run)
      if (!a) throw new ApiError(404, `no running dispatch run ${run}`)
      drop(a.pid)
      return reply('run-stop', run, { account: a.account, project: a.project })
    },
    async resume(id, account) {
      const rows = (data.sessions as { sessions: Array<{ conversation_id: string; account: string; project: string }> }).sessions
      const r = rows.find((x) => x.conversation_id === id)
      if (!r) throw new ApiError(404, `no session "${id}"`)
      return reply('resume', id, { account: account ?? r.account, project: r.project, message: `resumed ${r.project}` })
    },
    async start(s) {
      return reply('start', s.project, { text: s.prompt, account: s.account ?? 'main', project: s.project, message: `opened a window on ${s.project}` })
    },
    async dispatch(r) {
      const d = dispatch()
      if (d.slots.busy >= d.slots.max) throw new ApiError(409, `too_many_runs: ${d.slots.busy} runs already running`)
      const run = `f${String(d.agents.length).padStart(7, '0')}`
      const pid = 9000 + d.agents.length
      set('dispatch', {
        ...d,
        slots: { ...d.slots, busy: d.slots.busy + 1 },
        agents: [...d.agents, { kind: 'run', pid, account: 'work', project: r.project, cwd: `/src/${r.project}.wt/${r.label}`, run_id: run, label: r.label, title: r.prompt, state: 'working', started_at: now.toISOString() }],
      })
      return reply('dispatch', r.project, { text: r.prompt, account: 'work', project: r.project, run_id: run, message: `dispatched run ${run}` })
    },
    async mail(id, verb) {
      const t = data.today as { important: Array<{ id: string }>; flagged: Array<{ id: string }> }
      if (![...t.important, ...t.flagged].some((m) => m.id === id)) throw new ApiError(404, `no mail ${id}`)
      set('today', { ...t, important: t.important.filter((m) => m.id !== id), flagged: t.flagged.filter((m) => m.id !== id) })
      return reply(verb, id)
    },
  }
}

export function fixtureSource(missions: Mission[] = MISSIONS, inbox: InboxItem[] = INBOX, now: Date = NOW, agenda: AgendaEvent[] = AGENDA, live: LiveData = stateLive): DataSource {
  const data: LiveData = structuredClone(live)
  // One body object per section version, so unchanged sections keep their reference.
  const bodies = new Map<SectionName, SectionBody>()
  const section = async (name: SectionName) => {
    if (!(name in data)) return NO_SECTION
    let b = bodies.get(name)
    if (!b) bodies.set(name, (b = { data: data[name], stale: false, updated_at: now.toISOString() }))
    return b
  }
  return {
    missions: async () => missions,
    mission: async (id) => missions.find((m) => m.id === id),
    inbox: async () => inbox,
    agenda: async () => agenda,
    section,
    actions: fixtureActions(data, bodies, now),
    now: () => now,
  }
}
