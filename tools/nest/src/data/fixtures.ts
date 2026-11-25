import { missions as coreMissions, inbox as coreInbox } from 'virtual:core-records'
import { heroes as brickHeroes, pages as brickPages, themes as nestThemes } from 'virtual:nest-modules'
import stateMissions from '../../fixtures/states/missions.json'
import stateInbox from '../../fixtures/states/inbox.json'
import stateAgenda from '../../fixtures/states/agenda.json'
import stateLive from '../../fixtures/states/live.json'
import stateWork from '../../fixtures/states/work.json'
import stateOps from '../../fixtures/states/ops.json'
import stateMrs from '../../fixtures/states/mrs.json'
import stateWeek from '../../fixtures/states/agenda-week.json'
import stateReports from '../../fixtures/states/reports.json'
import stateJobs from '../../fixtures/states/jobs.json'
import { ApiError } from './api'
import { fixtureButler, type ButlerMsg } from './butler'
import type { ChatTarget, ConversationPage, Entry, PageAsk } from './conversation'
import type { Actions, Dispatch, LiveReply } from './live'
import type { Modules } from './modules'
import { isLive, shortPath, type Changes, type Ops, type OpsQueued } from './ops'
import type { NestBuild, NestKnown, NestSave, NestSaved, NestSettings, SettingsApi } from './settings'
import type { DataSource } from './source'
import { CANCELLABLE } from './state'
import type { AgendaEvent, InboxItem, InboxVerb, Mission, SectionBody, SectionName } from './types'

// Shared core-contract records plus Nest-local records that cover every UI state.
const MISSIONS = [...coreMissions, ...stateMissions] as Mission[]
const INBOX = [...coreInbox, ...stateInbox] as InboxItem[]
const AGENDA = stateAgenda as AgendaEvent[]
const NOW = new Date('2026-11-09T16:10:00+01:00')
// Sections without a state fixture stay empty until their screens bring one (item_171/172).
const NO_SECTION: SectionBody = { data: null, stale: false, updated_at: null }
type LiveData = Partial<Record<SectionName, unknown>>
// The work feed (item_204) over the live fixture: one row per kind and state, a mission rollup.
// The MRs board (item_225) also carries the fields crow-mrs does not serve yet (src/data/mrs.ts).
// The Agenda (item_226) adds crow-today `week` and dated crow-checks rows.
// Reports (item_258): settled AFK reports, rollups and retros; jobs carries the retro routine.
const LIVE: LiveData = {
  ...stateLive, work: stateWork, mrs: stateMrs, reports: stateReports, jobs: stateJobs,
  today: { ...stateLive.today, week: stateWeek.week },
  checks: { ...stateLive.checks, rows: [...stateLive.checks.rows, ...stateWeek.checks] },
}

// In-memory live actions over the fixture sections: they check targets like core
// (404 unknown, 409 full slots) and change the section so a refetch shows the result.
type MissionActions = 'missionAnswer' | 'missionStop' | 'missionCancel' | 'inboxHide'
type QueueActions = 'missionMove' | 'missionNext' | 'missionStart'

function fixtureActions(data: LiveData, bodies: Map<SectionName, SectionBody>, now: Date): Omit<Actions, MissionActions | QueueActions> {
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
    async answer(pid, prompt, option) {
      const a = agent(pid)
      const p = a.prompt
      if (!p || p.id !== prompt) throw new ApiError(409, `pid ${pid} no longer waits on prompt ${prompt}`)
      if (option === 'allow' && p.work_approval) throw new ApiError(409, `work approval on ${p.work_approval}: give it in approvals`)
      const d = dispatch()
      set('dispatch', { ...d, agents: d.agents.map((x) => (x === a ? { ...x, prompt: null, state: option === 'allow' ? 'working' : 'idle' } : x)) })
      return reply('answer', String(pid), { text: `${option}: ${p.tool} ${p.detail}`, account: a.account, project: a.project })
    },
    async ask(pid, prompt, answers) {
      const a = agent(pid)
      const p = a.prompt
      if (!p || p.id !== prompt) throw new ApiError(409, `pid ${pid} no longer waits on prompt ${prompt}`)
      if (!p.questions || answers.length !== p.questions.length) throw new ApiError(400, 'answers: one entry per question')
      const d = dispatch()
      set('dispatch', { ...d, agents: d.agents.map((x) => (x === a ? { ...x, prompt: null, state: 'working' } : x)) })
      return reply('answer', String(pid), { text: `answers: ${p.tool} ${p.detail}`, account: a.account, project: a.project })
    },
    async stop(pid) {
      const a = drop(pid)
      return reply('stop', String(pid), { account: a.account, project: a.project })
    },
    async attach(pid) {
      const a = dispatch().agents.find((x) => x.pid === pid)
      if (!a) throw new ApiError(404, `no agent with pid ${pid}`)
      return reply('attach', String(pid), { account: a.account, project: a.project, message: a.tmux ? `opened a terminal on tmux session ${a.tmux}` : 'focused its window' })
    },
    // A fixed claude screen with colors; keys and resize only check the target (sessions, not runs).
    async screen(pid) {
      const a = agent(pid)
      if (a.kind === 'run') throw new ApiError(409, `${pid} is a background run, it has no screen`)
      const lines = ['\x1b[38;5;174m╭───────────────────────────────╮\x1b[39m', `\x1b[38;5;174m│\x1b[39m \x1b[1m${a.project}\x1b[22m on ${a.account}`, '\x1b[38;5;174m╰───────────────────────────────╯\x1b[39m', '', '\x1b[2m> \x1b[22m']
      return { lines, cols: 80, rows: 24, cursor: { x: 2, y: 4, shown: true } }
    },
    async keys(pid) {
      if (agent(pid).kind === 'run') throw new ApiError(409, `${pid} is a background run, it has no screen`)
      return { ok: true }
    },
    async resize(pid) {
      if (agent(pid).kind === 'run') throw new ApiError(409, `${pid} is a background run, it has no screen`)
      return { ok: true }
    },
    async stopRun(run) {
      const a = dispatch().agents.find((x) => x.run_id === run)
      if (!a) throw new ApiError(404, `no running dispatch run ${run}`)
      drop(a.pid)
      return reply('run-stop', run, { account: a.account, project: a.project })
    },
    async retryRun(run) {
      const d = dispatch()
      const f = (d.finished ?? []).find((x) => x.run_id === run)
      if (!f) throw new ApiError(404, `no finished run "${run}"`)
      if (d.slots.busy >= d.slots.max) throw new ApiError(409, `too_many_runs: ${d.slots.busy} runs already running`)
      const id = `r${String(d.agents.length).padStart(7, '0')}`
      const account = f.account ?? 'work'
      set('dispatch', {
        ...d,
        slots: { ...d.slots, busy: d.slots.busy + 1 },
        finished: (d.finished ?? []).filter((x) => x !== f),
        agents: [...d.agents, { kind: 'run', pid: 9500 + d.agents.length, account, project: f.project, cwd: `/src/${f.project}.wt/${f.label ?? run}`, run_id: id, label: f.label ?? run, title: f.label ?? run, state: 'working', started_at: now.toISOString() }],
      })
      return reply('run-retry', run, { account, project: f.project, run_id: id, message: `retried ${run} as ${id}` })
    },
    async resume(id, account) {
      const rows = (data.sessions as { sessions: Array<{ conversation_id: string; account: string; project: string }> }).sessions
      const r = rows.find((x) => x.conversation_id === id)
      if (!r) throw new ApiError(404, `no session "${id}"`)
      return reply('resume', id, { account: account ?? r.account, project: r.project, message: `resumed ${r.project}` })
    },
    // A headless start (item_188): the session shows up at once with its tmux session.
    async start(s) {
      const d = dispatch()
      const account = s.account ?? 'main'
      const tmux = `crow-${s.project}-${String(d.agents.length).padStart(6, '0')}`
      set('dispatch', {
        ...d,
        agents: [...d.agents, { kind: 'session', pid: 8000 + d.agents.length, account, project: s.project, cwd: `/src/${s.project}`, state: 'working', started_at: now.toISOString(), last_prompt: s.prompt ?? null, tmux }],
      })
      return reply('start', s.project, { text: s.prompt, account, project: s.project, message: `started ${s.project} on ${account}, running in Nest`, tmux })
    },
    async dispatch(r) {
      const d = dispatch()
      if (d.slots.busy >= d.slots.max) throw new ApiError(409, `too_many_runs: ${d.slots.busy} runs already running`)
      const run = `f${String(d.agents.length).padStart(7, '0')}`
      const pid = 9000 + d.agents.length
      set('dispatch', {
        ...d,
        slots: { ...d.slots, busy: d.slots.busy + 1 },
        agents: [...d.agents, { kind: 'run', pid, account: r.account ?? 'work', project: r.project, cwd: `/src/${r.project}.wt/${r.label}`, run_id: run, label: r.label, title: r.prompt, state: 'working', started_at: now.toISOString() }],
      })
      return reply('dispatch', r.project, { text: r.prompt, account: r.account ?? 'work', project: r.project, run_id: run, message: `dispatched run ${run}` })
    },
    async mail(id, verb) {
      const t = data.today as { important: Array<{ id: string }>; flagged: Array<{ id: string }> }
      if (![...t.important, ...t.flagged].some((m) => m.id === id)) throw new ApiError(404, `no mail ${id}`)
      set('today', { ...t, important: t.important.filter((m) => m.id !== id), flagged: t.flagged.filter((m) => m.id !== id) })
      return reply(verb, id)
    },
    async routine(unit, verb) {
      if (!/^crow-[a-z0-9-]+\.timer$/.test(unit)) throw new ApiError(404, `no routine ${unit}`)
      return reply(`routine_${verb}`, unit)
    },
    async afkRead(id) {
      return { id, text: afkJobs.get(id) ?? `# ${id}\n\nRead the open items and fix the first one.\n` }
    },
    async afkSave(id, text) {
      afkJobs.set(id, text)
      return reply('afk_edit', id)
    },
    async afkSkip(id) {
      afkJobs.delete(id)
      return reply('afk_skip', id)
    },
    async module(brick, action, item) {
      return reply('module_action', `${brick}/${action}`, item ? { text: item } : {})
    },
  }
}

// Demo bricks for dev (?bricks): mail, forge and team off, one user brick with a list
// section and two actions on Desk.
export interface FixtureBricks { modules: Modules; sections: Record<string, unknown> }

export const DEMO_BRICKS: FixtureBricks = {
  modules: {
    contract_version: 1,
    version: 'demo',
    errors: [],
    modules: [
      { name: 'core-api', on: true, user: false, sections: ['sessions'], nest: {}, deck: { pages: [], boxes: [] } },
      { name: 'forge', on: false, user: false, sections: ['mrs'], nest: { sections: [{ id: 'mrs', page: 'mrs', source: 'mrs' }] }, deck: { pages: ['MRs'], boxes: ['mrs'] } },
      { name: 'mail', on: false, user: false, sections: ['today'], nest: { sections: [{ id: 'today', page: 'home', source: 'today' }] }, deck: { pages: [], boxes: ['today'] } },
      { name: 'team', on: false, user: false, sections: ['checks'], nest: {}, deck: { pages: [], boxes: ['checks'] } },
      {
        name: 'notes', on: true, user: true, sections: ['notes'], deck: { pages: [], boxes: [] },
        nest: { sections: [{ id: 'notes', page: 'desk', source: 'notes', title: 'Notes', view: 'list', actions: [
          { id: 'done', label: 'Done', confirm: true, row: true },
          { id: 'sweep', label: 'Sweep' },
        ] }] },
      },
    ],
  },
  sections: {
    notes: { items: [
      { id: 'n-1', title: 'Renew the domain', detail: 'circle registrar, auto-renew off', meta: 'due 28/11' },
      { id: 'n-2', title: 'Send the WSL notes to Paul', meta: 'today' },
      { id: 'n-3', title: 'Read the brick ADR again' },
    ] },
  },
}

const afkJobs = new Map<string, string>()

// The board chat (chat-cover): a user turn, folded tools, a list reply, a question.
const at = (hm: string) => `2026-11-09T${hm}:00+01:00`
const tool = (id: string, hm: string, name: string): Entry[] => [
  { at: at(hm), id: `${id}c`, role: 'tool', text: '', tool_name: name, tool_summary: 'Call recorded' },
  { at: at(hm), id: `${id}r`, role: 'tool', text: '', tool_name: name, tool_summary: 'Result recorded' },
]
export const CHAT: Entry[] = [
  { at: at('15:58'), id: 'e1', role: 'user', text: 'basically i want to be able to chat with you easily from the phone, not the diffs, all the text' },
  ...tool('t1', '15:58', 'Read'), ...tool('t2', '15:58', 'Bash'), ...tool('t3', '15:59', 'Read'), ...tool('t4', '15:59', 'Bash'),
  { at: at('15:59'), id: 'e2', role: 'assistant', text: 'Clear. On the phone you get every word I write, no diffs or logs.\n\nWhat changes:\n- Tap a session: its chat opens\n- Replies show in full, line breaks kept\n- Tool calls fold to one line, tap to open\n- The box at the bottom talks to that session' },
  { at: at('16:01'), id: 'e3', role: 'user', text: 'show me how it looks in nest' },
  ...['Bash', 'Read', 'Edit', 'Bash', 'Read', 'Edit', 'Bash', 'Read', 'Edit'].flatMap((n, i) => tool(`u${i}`, '16:02', n)),
  { at: at('16:06'), id: 'e4', role: 'assistant', text: 'Agenda modern is updated: reminder card, +38 -6.\n\nPublish now or continue with Reports?' },
]

// Pages over one in-memory list per target: the cursor is the count already read.
function fixtureChats(known: (t: ChatTarget) => void) {
  const chats = new Map<string, Entry[]>()
  const list = (t: ChatTarget) => {
    known(t)
    const k = 'pid' in t ? `p${t.pid}` : `m${t.mission}`
    let l = chats.get(k)
    if (!l) chats.set(k, (l = [...CHAT]))
    return l
  }
  const page = async (t: ChatTarget, p: PageAsk): Promise<ConversationPage> => {
    const l = list(t)
    let from = Math.max(0, l.length - (p.tail ?? 50))
    let to = l.length
    if (p.cursor) from = Number(p.cursor)
    else if (p.before) [from, to] = [Math.max(0, Number(p.before) - 50), Number(p.before)]
    return {
      account: 'rk1', conversation_id: 'fixture', entries: l.slice(from, to),
      next_cursor: String(l.length), before_cursor: from > 0 ? String(from) : undefined, has_more: from > 0,
    }
  }
  const said = (pid: number, text: string, now: Date) => list({ pid }).push({ at: now.toISOString(), id: `s${Math.random()}`, role: 'user', text })
  return { page, said }
}

// The exchange the desk-rail-a board draws (item_176).
const BUTLER_SEED: ButlerMsg[] = [
  { seq: 1, ts: '2026-11-09T14:20:00Z', role: 'you', text: 'Start a review of !224 lambda', surface: 'nest' },
  { seq: 2, ts: '2026-11-09T14:20:00Z', role: 'crow', text: 'Review mission started for !224, Sir. It runs on main.', account: 'main',
    chips: [{ label: 'Open mission', href: '/missions/20261109-build-mr-pipeline-badge' }] },
  { seq: 3, ts: '2026-11-09T14:31:00Z', role: 'you', text: "What's running?", surface: 'nest' },
  { seq: 4, ts: '2026-11-09T14:31:00Z', role: 'crow', text: 'Two sessions and one run, Sir. S1 Nest waits on you for a Bash approval. The infra push run failed: token rotation.', account: 'main',
    chips: [{ label: 'Open S1', href: '/desk' }, { label: 'Allow in Inbox', href: '/inbox' }, { label: 'Show failure', href: '/missions/20261109-fix-push-token-refresh' }] },
  { seq: 5, ts: '2026-11-09T14:33:00Z', role: 'you', text: 'Why did the token rotation fail?', surface: 'nest' },
  { seq: 6, ts: '2026-11-09T14:33:00Z', role: 'tool', text: 'read run log · infra push · 2 files' },
  { seq: 7, ts: '2026-11-09T14:34:00Z', role: 'crow', text: 'The refresh call got a 401, Sir. The deploy key expired yesterday at 18:00. A new key in the vault fixes it; the run can retry after that.', account: 'main',
    chips: [{ label: 'Show failure', href: '/missions/20261109-fix-push-token-refresh' }] },
  { seq: 8, ts: '2026-11-09T14:36:00Z', role: 'you', text: 'ok rotate it and retry', surface: 'nest' },
  { seq: 9, ts: '2026-11-09T14:36:00Z', role: 'crow', text: 'That needs your OK on infra, Sir. I put the request in your Inbox.', account: 'main',
    chips: [{ label: 'Allow in Inbox', href: '/inbox' }] },
  { seq: 10, ts: '2026-11-09T14:40:00Z', role: 'you', text: 'what is left for today?', surface: 'nest' },
  { seq: 11, ts: '2026-11-09T14:40:00Z', role: 'crow', text: 'Three things, Sir. The management weekly at 15:00, its prep is ready. Two reviews wait on you: !226 and !231. Sprint planning prep is due at 17:30.', account: 'main',
    chips: [{ label: 'Open prep', href: '/agenda' }, { label: 'Open reviews', href: '/mrs' }] },
]

// A drawn stand-in for a mockup image: one hue per option, the file name as caption.
const HUE: Record<string, number> = { A: 200, B: 150, C: 30 }
export function fixtureImage(path: string): string {
  const name = (path.split('/').pop() ?? path).replace(/\.[a-z]+$/i, '')
  const letter = /(?:^|[-_])([A-Z])(?:[-_]|$)/.exec(name)?.[1] ?? 'A'
  const h = HUE[letter] ?? 280
  const dark = /dark/.test(name)
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 640 400">`
    + `<rect width="640" height="400" fill="hsl(${h} 35% ${dark ? 10 : 22}%)"/>`
    + `<rect x="40" y="40" width="560" height="56" rx="8" fill="hsl(${h} 60% ${dark ? 30 : 45}%)"/>`
    + `<rect x="40" y="120" width="340" height="240" rx="8" fill="hsl(${h} 40% ${dark ? 16 : 30}%)"/>`
    + `<rect x="400" y="120" width="200" height="110" rx="8" fill="hsl(${h} 50% ${dark ? 22 : 38}%)"/>`
    + `<text x="400" y="290" font-family="sans-serif" font-size="22" fill="hsl(${h} 30% 85%)">${name}</text></svg>`
  return `data:image/svg+xml,${encodeURIComponent(svg)}`
}

export function fixtureSource(missions: Mission[] = MISSIONS, inbox: InboxItem[] = INBOX, now: Date = NOW, agenda: AgendaEvent[] = AGENDA, live: LiveData = LIVE, bricks?: FixtureBricks): DataSource {
  const data: LiveData = structuredClone(live)
  // One body object per section version, so unchanged sections keep their reference.
  const bodies = new Map<SectionName, SectionBody>()
  const section = async (name: SectionName) => {
    if (!(name in data)) return NO_SECTION
    let b = bodies.get(name)
    if (!b) bodies.set(name, (b = { data: data[name], stale: false, updated_at: now.toISOString() }))
    return b
  }
  const brickBodies = new Map(Object.entries(bricks?.sections ?? {}).map(([k, v]): [string, SectionBody] => [k, { data: v, stale: false, updated_at: now.toISOString() }]))
  const rec = { missions, inbox }
  const ops = { queue: structuredClone((stateOps as unknown as Ops).queue) }
  const busy = (n: number) => {
    const d = data.dispatch as Dispatch
    data.dispatch = { ...d, slots: { ...d.slots, busy: d.slots.busy + n } }
    bodies.delete('dispatch')
  }
  const actions = { ...fixtureActions(data, bodies, now), ...missionActions(rec, now), ...queueActions(ops, busy) }
  const chats = fixtureChats((t) => {
    if ('pid' in t && !(data.dispatch as Dispatch | undefined)?.agents.some((a) => a.pid === t.pid)) throw new ApiError(404, `no live session with pid ${t.pid}`)
    if ('mission' in t && !rec.missions.some((m) => m.id === t.mission)) throw new ApiError(404, 'no mission')
  })
  const send = actions.send
  actions.send = async (pid, text) => {
    const r = await send(pid, text)
    chats.said(pid, text, now)
    return r
  }
  return {
    missions: async () => rec.missions,
    mission: async (id) => rec.missions.find((m) => m.id === id),
    // The filters plan asks one choice; other missions carry no options.
    questions: async (id, version) => id === '20261109-build-inbox-filters' && version === 1 ? [{
      id: 'q1', text: 'Keep the filter in the URL?', blocking: true, recommended: 'url',
      choices: [{ id: 'url', label: 'In the URL. A shared link opens the same filter.' }, { id: 'local', label: 'In local storage. URLs stay short.' }],
    }] : [],
    artifact: (_id, path) => fixtureImage(path),
    inbox: async () => rec.inbox,
    agenda: async () => agenda,
    section,
    conversation: chats.page,
    ops: async (query) => fixtureOps(query, ops.queue),
    run: async (id) => fixtureRun(id),
    changes: async (cwd) => fixtureChanges(cwd),
    actions,
    now: () => now,
    butler: fixtureButler(now, BUTLER_SEED),
    settings: fixtureSettings(now),
    ...(bricks && {
      modules: async () => bricks.modules,
      brickSection: async (source: string) => brickBodies.get(source) ?? NO_SECTION,
    }),
  }
}


// In-memory mission actions with the checks of core (internal/mission): exact ids,
// open rows, the version the owner saw, verbs of the row. Each change swaps the array
// so a refetch sees new references.
function missionActions(rec: { missions: Mission[]; inbox: InboxItem[] }, now: Date): Pick<Actions, MissionActions> {
  const ts = now.toISOString()
  const find = (id: string) => {
    const m = rec.missions.find((x) => x.id === id)
    if (!m) throw new ApiError(404, `no mission "${id}"`)
    return m
  }
  const put = (m: Mission) => { rec.missions = rec.missions.map((x) => (x.id === m.id ? m : x)) }
  const row = (r: InboxItem) => { rec.inbox = rec.inbox.map((x) => (x.id === r.id ? r : x)) }
  return {
    async missionAnswer(mission: string, item: string, verb: InboxVerb, text: string, version: number) {
      const m = find(mission)
      const id = item.startsWith(`${mission}/`) ? item : `${mission}/${item}`
      const r = rec.inbox.find((x) => x.id === id)
      if (!r) throw new ApiError(404, `no inbox item "${id}"`)
      if (version !== m.version) throw new ApiError(409, `${id}: answered from v${version}, the mission is at v${m.version}`)
      if (r.state !== 'open') throw new ApiError(409, `${id} is ${r.state}`)
      if (!r.verbs.includes(verb)) throw new ApiError(400, `verb "${verb}", want one of ${r.verbs.join(', ')}`)
      row({ ...r, state: 'answered', answered: ts })
      if (!rec.inbox.some((x) => x.mission === mission && x.state === 'open')) put({ ...m, state: 'queued', updated: ts })
      return { item: id, verb, text, state: 'queued' }
    },
    async missionStop(mission: string) {
      const m = find(mission)
      if (m.state !== 'running') throw new ApiError(409, `${mission} is ${m.state}; only a running mission stops`)
      put({ ...m, state: 'waiting', updated: ts, runs: m.runs.map((r) => (r.ended ? r : { ...r, ended: ts, cause: 'stopped' as const })) })
      return find(mission)
    },
    async missionCancel(mission: string) {
      const m = find(mission)
      if (!CANCELLABLE.includes(m.state)) throw new ApiError(409, `${mission} is ${m.state} and cannot be cancelled`)
      put({ ...m, state: 'cancelled', updated: ts, runs: m.runs.map((r) => (r.ended ? r : { ...r, ended: ts, cause: 'stopped' as const })) })
      rec.inbox = rec.inbox.map((x) => (x.mission === mission && x.state === 'open' ? { ...x, state: 'dropped' as const } : x))
      return find(mission)
    },
    async inboxHide(item: string) {
      const r = rec.inbox.find((x) => x.id === item)
      if (!r) throw new ApiError(404, `no inbox item "${item}"`)
      if (r.state !== 'open') throw new ApiError(409, `${item} is ${r.state}`)
      row({ ...r, state: 'dropped' })
      return rec.inbox.find((x) => x.id === item)
    },
  }
}

// Queue order over the ops fixture (item_272), with the checks of core (internal/mission/queue.go):
// exact ids, missions only, place 1 or more. Start now runs the mission at once (busy + 1).
function queueActions(ops: { queue: OpsQueued[] }, busy: (n: number) => void): Pick<Actions, QueueActions> {
  const missions = () => ops.queue.filter((q) => q.kind === 'mission')
  const at = (id: string) => {
    const i = missions().findIndex((q) => q.id === id)
    if (i < 0) throw new ApiError(ops.queue.some((q) => q.id === id) ? 409 : 404, `${id} is not a queued mission`)
    return i
  }
  const place = (id: string, to: number) => {
    const ms = missions()
    const i = at(id)
    if (to < 1) throw new ApiError(400, `place ${to}: want 1 or more`)
    const [m] = ms.splice(i, 1)
    ms.splice(Math.min(to, ms.length + 1) - 1, 0, m)
    ops.queue = [...ms, ...ops.queue.filter((q) => q.kind !== 'mission')]
    return ms
  }
  return {
    async missionMove(mission, to) {
      return place(mission, typeof to === 'number' ? to : Math.max(1, at(mission) + 1 + (to === 'up' ? -1 : 1)))
    },
    async missionNext(mission) { return place(mission, 1) },
    async missionStart(mission) {
      at(mission)
      ops.queue = ops.queue.filter((q) => q.id !== mission)
      busy(1)
      return { id: mission, state: 'running' }
    },
  }
}

// Ops fixture (item_190): one day. A week answer repeats its lanes on each past day of the week.
function fixtureOps(query: string, queue: OpsQueued[]): Ops {
  const day = { ...(stateOps as unknown as Ops), queue }
  const week = /^week=(.+)$/.exec(query)?.[1]
  if (week) {
    const days = ['2026-11-09', '2026-11-10', '2026-11-11', '2026-11-12', '2026-11-13', '2026-11-14', '2026-11-15']
    return { week, days, lanes: day.lanes, queue: day.queue, totals: day.totals }
  }
  const d = /^day=(.+)$/.exec(query)?.[1] ?? day.day ?? ''
  return d === day.day ? day : { day: d ?? '', days: [d ?? ''], lanes: [], queue: day.queue, totals: [] }
}

function fixtureRun(id: string) {
  const l = (stateOps as unknown as Ops).lanes.find((x) => x.keys.includes(id))
  if (!l) return undefined
  return {
    run: id, source: 'dispatch', project: l.project, label: l.title, account: l.account, status: l.state,
    started: l.spans[0]?.start ?? null, ended: l.spans.at(-1)?.end ?? null, cost_usd: l.cost_usd, commits: [],
    branch: l.title, worktree: `/work/${l.project}.wt/${l.title}`, item: null, mission: null,
    report: `## ${l.title}\n\nFixture report: what the run did and what is left.`,
    ...(isLive(l) ? {} : {
      commits: ['4e1a9c2 chg[MD]: ops docs name the lane sheet tabs', '91b0d3e new[MD]: runbook for the ops page'],
      diffstat: { base: '0f3c2a1', files: 2, insertions: 64, deletions: 9, truncated: false, per_file: [
        { path: 'runbooks/ops-page.md', insertions: 58, deletions: 0 }, { path: 'docs/ops.md', insertions: 6, deletions: 9 },
      ] },
    }),
  }
}

// Changes fixture (item_273): the live lanes of the ops fixture have a cwd.
function fixtureChanges(cwd: string): Changes | undefined {
  if (!(stateOps as unknown as Ops).lanes.some((l) => l.cwd === cwd)) return undefined
  const t = (m: number) => new Date(Date.parse('2026-11-12T11:00:00Z') - m * 60e3).toISOString()
  return {
    cwd, branch: shortPath(cwd).split('/').pop() === 'cr-0w' ? 'main' : 'item-90-ops', base: 'main', ahead: 3, behind: 0, uncommitted: 2,
    diffstat: { base: '0f3c2a1', files: 5, insertions: 233, deletions: 38, truncated: false, per_file: [
      { path: 'tools/nest/src/screens/Ops.tsx', insertions: 142, deletions: 0 },
      { path: 'tools/nest/src/theme.css', insertions: 61, deletions: 12 },
      { path: 'tools/crow-core/internal/api/changes.go', insertions: 18, deletions: 4 },
      { path: 'tools/nest/src/data/ops.ts', insertions: 9, deletions: 2 },
      { path: 'logics/backlog/item_90_ops_page.md', insertions: 3, deletions: 20 },
    ] },
    commits: [
      { sha: 'a41c9e0', subject: 'new[TS/CSS]: lane sheet tabs', time: t(12) },
      { sha: '7d02b5f', subject: 'chg[GO]: core serves lane changes', time: t(48) },
      { sha: 'c9e13aa', subject: 'chg[MD]: ops page item scope', time: t(95) },
    ],
    more: false,
  }
}

// Names the fixture Nest can build: the engine ones plus the compiled-in bricks (item_284).
export const FIXTURE_KNOWN: NestKnown = {
  themes: [...nestThemes],
  heroes: ['holotable', 'none', ...Object.keys(brickHeroes)],
  pages: ['home', 'desk', 'ops', 'inbox', 'mrs', 'agenda', 'reports', ...brickPages.map((p) => p.key)],
}

type NestYaml = Record<'theme' | 'hero' | 'pages', unknown>

// In-memory crow.yaml modules.nest with the checks of crow-modules set nest (an unknown name
// answers 400 and writes nothing). A hero or pages save starts a build that passes on the second poll.
export function fixtureSettings(now: Date, known: NestKnown = FIXTURE_KNOWN, conf: Partial<NestYaml> = {}): SettingsApi {
  const yaml: NestYaml = { theme: conf.theme ?? null, hero: conf.hero ?? null, pages: conf.pages ?? null }
  let served = 1
  let polls = 0
  let build: NestBuild = { state: 'idle', elapsed_s: 0, queued: false, log: [], served: `fixture-${served}` }
  const get = (): NestSettings => ({ file: '/fixture/crow.yaml', ...structuredClone(yaml), known, errors: [] })
  const check = (k: string, v: string, list: string[]) => {
    if (!list.includes(v)) throw new ApiError(400, `crow-modules: unknown ${k} "${v}" (known: ${list.join(', ')})`)
  }
  return {
    get: async () => get(),
    async build() {
      if (build.state === 'building' && ++polls >= 2) {
        served++
        build = { ...build, state: 'ok', elapsed_s: 38, log: ['nest build: ok'], served: `fixture-${served}` }
      } else if (build.state === 'building') build = { ...build, elapsed_s: build.elapsed_s + 12 }
      return { ...build }
    },
    async save(body: NestSave): Promise<NestSaved> {
      if (body.theme !== undefined) check('theme', body.theme, known.themes)
      if (body.hero !== undefined) check('hero', body.hero, known.heroes)
      for (const p of body.pages ?? []) check('page', p, known.pages)
      const changed: string[] = []
      for (const k of ['theme', 'hero', 'pages'] as const) {
        const v = body[k]
        if (v === undefined || JSON.stringify(v) === JSON.stringify(yaml[k])) continue
        yaml[k] = v
        changed.push(k)
      }
      const needs = changed.includes('hero') || changed.includes('pages')
      if (needs) {
        polls = 0
        build = { state: 'building', started_at: now.toISOString(), elapsed_s: 0, queued: false, log: [], served: `fixture-${served}` }
      }
      return { ...get(), changed, needs_build: needs, build: { ...build } }
    },
  }
}
