import { Fragment, useEffect, useRef, useState } from 'react'
import { ApiError } from '../data/api'
import { deliveryOf, pendingOf, type Delivery, type Pending } from '../data/delivery'
import { useConversation } from '../data/useConversation'
import { DeliveryLine, useDeliveryClock } from '../ui/Delivery'
import { accountsOf, dispatchOf, projectsOf, resumableOf, type Agent, type AskAnswer, type LiveReply, type Resumable } from '../data/live'
import { useSource, type Snapshot } from '../data/useData'
import { ago, time, title } from '../format'
import type { Mission } from '../data/types'
import { CANCELLABLE, LABEL, missionOfRun } from '../data/state'
import { navigate, useQuery } from '../router'
import { Sheet, type Confirm } from '../ui/Sheet'
import { cancelMission, refused, stopMissionRun } from '../ui/missionSheets'
import { AskDock } from '../ui/AskDock'
import { promptView } from '../ui/prompt'
import { enterSubmits, rowsOf } from '../ui/keys'
import { Chat as AgentChat } from './Chat'
import { Deck } from './Mrs'
import { LAPTOP, useMedia } from '../ui/useMedia'
import { KIND_LABEL, workOf, workWhen, type WorkItem } from '../data/work'
import { KindIcon, type RowAct } from '../ui/WorkList'
import { stateClass } from '../ui/Glyph'
import { AccountField } from '../ui/AccountField'
import { isCommand } from '../data/butler'
import { useButler } from '../data/useButler'
import { deskPick, readPick, savePick } from '../data/deskPick'
import {
  PANE_IDS, PANE_TITLE, agentOf, agentUnit, flatOf, itemRef, paneGroups, panesOf, projectsOfPanes, summaryOf, unitOfAgent,
  type PaneId, type Panes,
} from '../data/deskPanes'

// Desk layout C (item_295): four deck panes (needs you, moving, queued, done) and a tabbed
// detail, on the live actions of item_170. Every action goes through the confirm sheet,
// then the feed refetches. The butler chat stays folded until asked for.
const MODELS = ['default', 'sonnet', 'opus', 'haiku']
const LATER = 'Arrives with the live API'
const CORE = 'item_297'
const nameOf = (a: Agent) => (a.kind === 'run' ? a.label || a.run_id || String(a.pid) : a.title || a.project)
// setup (item_240): a headless session stuck on a first-run dialog; its Transcript tab answers it.
const waits = (a: Agent) => a.state === 'waiting' || a.state === 'setup'
const said = (a: Agent) => (a.state === 'setup' ? 'waits on setup' : waits(a) ? 'waits' : a.state)

const Mic = ({ s = 18 }: { s?: number }) => (
  <svg width={s} height={s} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
    <rect x="9" y="3" width="6" height="11" rx="3" /><path d="M5 11a7 7 0 0 0 14 0M12 18v3" />
  </svg>
)

interface Line { t?: string; tool: string; text: string; says?: boolean }
// Timeline of an agent from what the runs section carries: start, last prompt, recent actions.
function linesOf(a: Agent, now: Date): Line[] {
  const out: Line[] = []
  if (a.started_at) out.push({ t: time(new Date(a.started_at)), tool: 'Start', text: `${a.project} · ${a.account} · ${ago(a.started_at, now)} ago` })
  if (a.last_prompt) out.push({ tool: 'You', text: a.last_prompt })
  const acts = a.recent?.length ? a.recent : a.last_action ? [{ action: a.last_action, ts: '' }] : []
  for (const r of acts) {
    const m = /^([\w-]+):\s*(.*)$/.exec(r.action)
    const tool = m ? m[1] : r.action
    const says = tool.toLowerCase() === 'says'
    out.push({ t: r.ts && !isNaN(Date.parse(r.ts)) ? time(new Date(r.ts)) : undefined, tool: says ? 'Says' : tool, text: m ? m[2] : '', says })
  }
  return out
}

// A Tell it message carries its pending send and delivery state.
interface Msg { me?: boolean; text: string; f?: boolean; p?: Pending; state?: Delivery }
type Form = 'run' | 'session' | null
export type Tab = 'overview' | 'transcript' | 'timeline' | 'files'
const TABS: [Tab, string][] = [['overview', 'Overview'], ['transcript', 'Transcript'], ['timeline', 'Timeline'], ['files', 'Files']]
type Mode = 'now' | 'background' | 'afk'
const MODES: [Mode, string][] = [['now', 'Now'], ['background', 'Background'], ['afk', 'AFK batch']]
type Focus = PaneId | 'detail'
// One action of the detail: its label, an optional key, and its look.
interface Act { label: string; run: () => void; key?: string; p?: boolean; x?: boolean }

const slug = (s: string) => s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '').split('-').slice(0, 4).join('-')

function useDesk(snap: Snapshot, refresh: () => Promise<void>) {
  const src = useSource()
  const act = src.actions
  const { slots, agents, finished = [] } = dispatchOf(snap.sections.dispatch)
  const resumable = resumableOf(snap.sections.sessions)
  const projects = projectsOf(snap.sections.projects)
  const accounts = accountsOf(snap.sections.accounts)
  const q = new URLSearchParams(location.search)

  const [sheet, setSheet] = useState<Confirm | null>(null)
  const [log, setLog] = useState<Msg[]>([])
  // Home links open a unit directly: /desk?pid=N (&tell=1 opens the reply box, &tab=transcript a tab).
  // Otherwise the last pick of this device (item_210): Desk remounts on every tab change.
  const [sel, setSel] = useState<number | null>(() => Number(q.get('pid')) || readPick())
  useEffect(() => { if (sel != null) savePick(sel) }, [sel])
  const [uid, setUid] = useState<string | null>(null)
  const [tell, setTell] = useState(() => q.get('tell') === '1')
  const [tab, setTab] = useState<Tab | null>(() => (TABS.some(([t]) => t === q.get('tab')) ? (q.get('tab') as Tab) : null))
  // Phone: the detail is its own page, open on a ?pid= link or a tap.
  const [page, setPage] = useState(() => !!q.get('pid'))
  const [text, setText] = useState('')
  const [ask, setAsk] = useState('')
  const [chat, setChat] = useState(false)
  const butler = useButler()
  const [form, setForm] = useState<Form>(null)
  const [mode, setMode] = useState<Mode>('now')
  const [task, setTask] = useState('')
  const [run, setRun] = useState({ project: '', label: '', model: 'default', account: '', prompt: '' })
  const [start, setStart] = useState({ project: '', account: '', prompt: '' })
  const [proj, setProj] = useState('')
  const [focusSet, setFocus] = useState<Focus | null>(null)
  const [folded, setFolded] = useState(true)
  // tmux session of a session just started here: its detail opens once a poll shows it (item_188).
  const [opening, setOpening] = useState<string | null>(null)

  // The work feed (item_202) fills the panes; null on a core without it (live agents stand in).
  const work = workOf(snap.sections.work)
  const queue = snap.missions.filter((m) => m.state === 'queued' || m.state === 'preflight')
  const panes = panesOf(work, agents, queue, finished, proj)
  const units = PANE_IDS.flatMap((id) => panes[id])
  const paneOf = (w: WorkItem) => PANE_IDS.find((id) => panes[id].includes(w))
  const picked = deskPick(agents, work, sel)
  const pickedUnit = picked && (!proj || picked.project === proj) ? unitOfAgent(panes, picked) ?? agentUnit(picked) : undefined
  const unit: WorkItem | undefined = units.find((w) => w.id === uid) ?? pickedUnit ?? panes.needs[0] ?? panes.moving[0] ?? units[0]
  const agent = unit && agentOf(unit, agents)
  const focus: Focus = focusSet ?? (unit && paneOf(unit)) ?? 'needs'

  // Success closes the sheet and answers in the chat; a refusal stays in the sheet
  // (thrown) and refetches when the view was stale (409).
  const wrap = (said: string, call: () => Promise<LiveReply>, done?: () => void) => async () => {
    try {
      const r = refused(await call())
      setLog((l) => [...l, { me: true, text: said }, { text: r.message || `${r.action} ${r.target}`, f: !r.ok }])
      done?.()
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void refresh()
      throw e
    }
    void refresh()
  }
  const where = (a: Agent) => `${a.project} (${a.account})`
  const tgt = (a: Agent): Array<[string, string]> => [['session', nameOf(a)], ['account', a.account], ['pid', String(a.pid)]]

  // A message typed into a session shows "sending…" in the chat until its transcript
  // carries the turn (item_208). The tail read before the send is the baseline.
  const tellIt = async (a: Agent, t: string) => {
    const before = await src.conversation({ pid: a.pid }, { tail: 50 }).then((r) => r.entries, () => [])
    try {
      await act.send(a.pid, t)
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void refresh()
      throw e
    }
    const p = pendingOf(a.pid, t, before, Date.now())
    setLog((l) => [...l.filter((m) => m.p?.pid !== a.pid || m.p.text !== t || m.state !== 'failed'), { me: true, text: `Tell ${nameOf(a)}: ${t}`, p, state: 'sending' }])
    void refresh()
  }
  const settle = (p: Pending, state: Delivery) => setLog((l) => l.map((m) => (m.p === p ? { ...m, state } : m)))
  const retell = (p: Pending) => {
    const a = agents.find((x) => x.pid === p.pid)
    if (!a) return setLog((l) => [...l, { text: 'That session is gone.', f: true }])
    tellIt(a, p.text).catch((e: Error) => setLog((l) => [...l, { text: e.message, f: true }]))
  }
  const send = (a: Agent) => setSheet({
    title: 'Send to session', lead: `Type this answer into ${where(a)} and press enter:`, quote: text, target: tgt(a),
    note: 'Lands in the open window, like typing it yourself.', ok: 'Send',
    run: async () => { await tellIt(a, text); setText(''); setTell(false) },
  })
  const pause = (a: Agent) => setSheet({
    title: 'Pause session', lead: `Interrupt the current turn of ${where(a)} (Escape in its window).`, target: tgt(a),
    note: 'The session stays open and waits for you.', ok: 'Pause', run: wrap(`Pause ${nameOf(a)}`, () => act.interrupt(a.pid)),
  })
  const stop = (a: Agent) => {
    const m = a.kind === 'run' && a.run_id ? missionOfRun(snap.missions, a.run_id) : undefined
    if (m) {
      setSheet(stopMissionRun({
        mission: m, run: a.run_id!, name: nameOf(a), act, refresh: () => void refresh(),
        done: (what) => {
          setLog((l) => [...l, what === 'cancelled'
            ? { me: true, text: `Stop ${nameOf(a)} and cancel the mission` }
            : { me: true, text: `Restart ${nameOf(a)}` }, { text: what === 'cancelled' ? `${title(m)} cancelled.` : `${title(m)}: run stopped, the next tick starts a new one.` }])
          setSel(null)
        },
      }))
      return
    }
    stopPlain(a)
  }
  const stopPlain = (a: Agent) => setSheet(a.kind === 'run' && a.run_id
    ? {
      title: 'Stop run', lead: `Stop the run ${nameOf(a)} of ${a.project}.`, target: [['run', a.run_id], ['project', a.project], ['account', a.account]],
      note: 'Tap again to confirm. Its commits stay on the run branch.', ok: 'Tap again to stop', danger: true,
      run: wrap(`Stop ${nameOf(a)}`, () => act.stopRun(a.run_id!), () => setSel(null)),
    }
    : {
      title: 'Stop session', lead: `End the claude session of ${a.project} (SIGTERM).`, target: tgt(a),
      note: 'Tap again to confirm. The window closes; the conversation can be resumed.', ok: 'Tap again to stop', danger: true,
      run: wrap(`Stop ${nameOf(a)}`, () => act.stop(a.pid), () => setSel(null)),
    })
  const dispatch = () => setSheet({
    title: 'Start a run', lead: 'Opens a new run on a free allowed account. Nothing is pushed.', quote: run.prompt,
    target: [['project', run.project], ['label', run.label], ['model', run.model], ['account', run.account || 'auto'], ['slots', `${slots.busy} of ${slots.max} busy`]],
    note: 'The prompt above is sent as written.', ok: 'Start run',
    run: wrap(`/dispatch ${run.project} ${run.label}`, () => act.dispatch({ project: run.project, label: run.label, prompt: run.prompt, model: run.model === 'default' ? undefined : run.model, account: run.account || undefined }),
      () => { setRun({ project: '', label: '', model: 'default', account: '', prompt: '' }); setForm(null) }),
  })
  const open = () => setSheet({
    title: 'New session', lead: `Starts a claude session on ${start.project} inside Nest, no desktop terminal. Its detail opens once it shows up.`, quote: start.prompt || undefined,
    target: [['project', start.project], ['account', start.account || 'auto']],
    note: 'Fails when that account is busy; nothing else is started. Its Transcript tab shows the session screen.', ok: 'Start session',
    run: wrap(`/session ${start.project}`, async () => {
      const r = await act.start({ project: start.project, account: start.account || undefined, prompt: start.prompt || undefined })
      if (r.tmux) setOpening(r.tmux)
      return r
    }, () => { setStart({ project: '', account: '', prompt: '' }); setForm(null) }),
  })
  const resume = (r: Resumable) => setSheet({
    title: 'Resume session', lead: `Reopens the conversation of ${r.project} inside Nest, no desktop terminal. Its detail opens once it shows up.`,
    target: [['project', r.project], ['account', r.account], ['id', r.conversation_id.slice(0, 8)]],
    note: 'Same conversation, same account.', ok: 'Resume', run: wrap(`Resume ${r.project}`, async () => {
      const res = await act.resume(r.conversation_id, r.account)
      if (res.tmux) setOpening(res.tmux)
      return res
    }),
  })
  const cancel = (m: Mission) => setSheet(cancelMission({
    mission: m, act, refresh: () => void refresh(),
    done: () => setLog((l) => [...l, { me: true, text: `Cancel ${title(m)}` }, { text: `${title(m)} cancelled.` }]),
  }))

  // Picking a unit opens its detail (a page on the phone): the chat for a live agent, else Overview.
  const choose = (w: WorkItem) => {
    setUid(w.id)
    const a = agentOf(w, agents)
    if (a) setSel(a.pid)
    setTell(false)
    setTab(null)
    setPage(true)
  }
  const pick = (a: Agent, talk = false) => {
    setUid(agentUnit(a).id)
    setSel(a.pid)
    setTell(talk)
    setTab(talk ? 'overview' : null)
    setPage(true)
  }
  const fresh = opening ? agents.find((a) => a.tmux === opening) : undefined
  useEffect(() => {
    if (!fresh) return
    setOpening(null)
    pick(fresh)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fresh?.pid])

  // A prompt answer is one tap: the card already shows the exact command (item_174).
  const answer = (a: Agent, key: 'allow' | 'deny') => {
    const p = a.prompt
    if (!p) return
    wrap(`${key === 'allow' ? 'Allow once' : 'Deny'} ${nameOf(a)}: ${p.tool} ${p.detail}`, () => act.answer(a.pid, p.id, key))()
      .catch((e: Error) => setLog((l) => [...l, { text: e.message, f: true }]))
  }
  // An AskUserQuestion answer sends every step at once (item_193).
  const [asking, setAsking] = useState(false)
  const answerAsk = (a: Agent, answers: AskAnswer[]) => {
    const p = a.prompt
    if (!p) return
    setAsking(true)
    wrap(`Answer ${nameOf(a)}`, () => act.ask(a.pid, p.id, answers))()
      .catch((e: Error) => setLog((l) => [...l, { text: e.message, f: true }]))
      .finally(() => setAsking(false))
  }
  const retry = (w: WorkItem) => {
    const f = finished.find((x) => x.run_id === w.source.ref)
    setSheet({
      title: 'Retry run', lead: `Start ${w.title} of ${w.project} again on the same prompt.`,
      target: [['run', w.source.ref], ['project', w.project], ['account', f?.account || 'same'], ['slots', `${slots.busy} of ${slots.max} busy`]],
      note: 'A new run with its own worktree. Nothing is pushed.', ok: 'Retry',
      run: wrap(`Retry ${w.title}`, () => act.retryRun(w.source.ref)),
    })
  }
  const allow = (a: Agent) => {
    const p = a.prompt!
    setSheet({
      title: 'Allow once', lead: `${nameOf(a)} asks to run ${p.tool}${p.description ? `: ${p.description}` : ''}.`, quote: p.detail,
      target: [['session', nameOf(a)], ['project', a.project], ['account', a.account]],
      note: 'Same as Allow once in its window.', ok: 'Allow once',
      run: wrap(`Allow once ${nameOf(a)}: ${p.tool} ${p.detail}`, () => act.answer(a.pid, p.id, 'allow')),
    })
  }
  // The main action of a unit (item_246): Retry a failed run, open an MR, allow or answer a session.
  const rowAct = (w: WorkItem): RowAct | undefined => {
    if (w.kind === 'run' && w.state === 'failed' && finished.some((x) => x.run_id === w.source.ref)) return { label: 'Retry', run: () => retry(w) }
    if (w.kind === 'mr') return { label: 'Open MR', run: () => window.open(w.source.ref, '_blank', 'noopener') }
    if (w.kind === 'session' && w.state === 'needs') {
      const a = agents.find((x) => String(x.pid) === w.source.ref)
      if (!a) return undefined
      if (a.prompt && promptView(a.prompt) === 'options' && a.prompt.options.some((o) => o.key === 'allow')) return { label: 'Allow', run: () => allow(a) }
      return { label: 'Answer', run: () => pick(a, !a.prompt) }
    }
    if (w.kind === 'mission' && w.state === 'needs') return { label: 'Answer', run: () => navigate(`/ops?mission=${encodeURIComponent(w.source.ref)}`) }
    return undefined
  }
  const missionOf = (w: WorkItem) => snap.missions.find((m) => m.id === (w.kind === 'mission' ? w.source.ref : w.mission))
  // Every action of a unit, main one first. Keys: a answer, r retry or resume, x stop or cancel, o open session.
  const actsOf = (w: WorkItem): Act[] => {
    const out: Act[] = []
    const main = rowAct(w)
    if (main) out.push({ ...main, p: true, key: main.label === 'Retry' ? 'r' : main.label === 'Open MR' ? undefined : 'a' })
    const a = agentOf(w, agents)
    if (a) {
      out.push({ label: 'Open session', key: 'o', run: () => navigate(`/desk/${a.pid}`) })
      if (a.kind === 'session') {
        out.push({ label: 'Message', run: () => { setTab('overview'); setTell(!tell) } })
        out.push({ label: 'Pause', run: () => pause(a) })
      }
      out.push({ label: 'Stop', key: 'x', x: true, run: () => stop(a) })
    }
    if (w.kind === 'parked') {
      const r = resumable.find((x) => x.conversation_id === w.source.ref)
      if (r) out.push({ label: 'Resume', key: 'r', p: !main, run: () => resume(r) })
    }
    const m = missionOf(w)
    if (m) {
      out.push({ label: 'Open mission', run: () => navigate(`/ops?mission=${encodeURIComponent(m.id)}`) })
      if (!a && CANCELLABLE.includes(m.state)) out.push({ label: 'Cancel', key: 'x', x: true, run: () => cancel(m) })
    }
    return out.filter((x, i) => out.findIndex((y) => y.label === x.label) === i)
  }

  // Launcher: "project, prompt" (or the project as first word) opens the form of the mode.
  const launch = () => {
    const t = task.trim()
    const comma = t.indexOf(',')
    const head = comma >= 0 ? t.slice(0, comma).trim() : t.split(/\s+/)[0]
    const p = projects.find((x) => x === head)
    const prompt = p ? (comma >= 0 ? t.slice(comma + 1) : t.slice(head.length)).trim() : t
    if (mode === 'now') {
      setStart((s) => ({ ...s, project: p ?? s.project, prompt: prompt || s.prompt }))
      setForm('session')
    } else if (mode === 'background') {
      setRun((r) => ({ ...r, project: p ?? r.project, prompt: prompt || r.prompt, label: r.label || itemRef(prompt) || slug(prompt) }))
      setForm('run')
    }
    setTask('')
  }

  // /dispatch and /session open their form, prefilled with the project and prompt typed after.
  const command = (t: string) => {
    const [cmd, ...rest] = t.trim().split(/\s+/)
    const words = rest.join(' ')
    const p = projects.find((x) => x === rest[0])
    if (cmd === '/dispatch') setRun((r) => ({ ...r, project: p ?? r.project, prompt: (p ? rest.slice(1).join(' ') : words) || r.prompt }))
    else setStart((s) => ({ ...s, project: p ?? s.project, prompt: (p ? rest.slice(1).join(' ') : words) || s.prompt }))
    setMode(cmd === '/dispatch' ? 'background' : 'now')
    setForm(cmd === '/dispatch' ? 'run' : 'session')
  }

  // The ask box: commands open their form; anything else goes to the butler (item_176).
  const submitAsk = () => {
    const t = ask.trim()
    if (!t) return
    if (isCommand(t)) command(t)
    else void butler?.say(t, matchMedia(LAPTOP).matches ? 'nest' : 'phone')
    setAsk('')
  }

  // The butler rail sends commands here as /desk?cmd=...: open the form, then drop the param.
  const cmd = useQuery().get('cmd')
  useEffect(() => {
    if (!cmd || !isCommand(cmd)) return
    command(cmd)
    const q = new URLSearchParams(location.search)
    q.delete('cmd')
    history.replaceState(null, '', location.pathname + (q.size ? `?${q}` : ''))
    dispatchEvent(new PopStateEvent('popstate'))
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [cmd])

  // Keys (laptop board): 1-4 a pane (3 unfolds Queued), 5 the detail, j k move, n new task,
  // c chat, a r x o the unit's actions. Off while typing or while a sheet is open.
  const key = (k: string): boolean => {
    if (k >= '1' && k <= '4') {
      const id = PANE_IDS[Number(k) - 1]
      const again = id === 'queued' && focus === 'queued'
      if (id === 'queued') setFolded((f) => (again ? !f : false))
      setFocus(id)
      const first = panes[id][0]
      if (first && !again) choose(first)
      return true
    }
    if (k === '5') { setFocus('detail'); return true }
    if (k === 'j' || k === 'k') {
      const flat = flatOf(panes, folded)
      if (!flat.length) return false
      const i = unit ? flat.findIndex((w) => w.id === unit.id) : -1
      const next = flat[Math.min(flat.length - 1, Math.max(0, i + (k === 'j' ? 1 : -1)))]
      choose(next)
      setFocus(paneOf(next) ?? 'needs')
      return true
    }
    if (k === 'n') { document.getElementById('dk-new')?.focus(); return true }
    if (k === 'c') { setChat((c) => !c); return true }
    const hit = unit && actsOf(unit).find((x) => x.key === k)
    if (hit) { hit.run(); return true }
    return false
  }

  return {
    snap, refresh, slots, agents, finished, resumable, projects, accounts,
    sheet, setSheet, log, tell, setTell, text, setText, ask, setAsk, submitAsk, chat, setChat,
    form, setForm, run, setRun, start, setStart, mode, setMode, task, setTask, launch, dispatch, open,
    send, settle, retell, answer, answerAsk, asking,
    work, panes, unit, agent, focus, setFocus, folded, setFolded, tab: tab ?? (agent && agent.state !== 'setup' ? 'transcript' : 'overview') as Tab, setTab, page, setPage, proj, setProj,
    allProjects: projectsOfPanes(work, agents), choose, actsOf, missionOf, paneOf, key,
  }
}
type D = ReturnType<typeof useDesk>

// ---- shared pieces ----
function TellBox({ d, a }: { d: D; a: Agent }) {
  return (
    <form className="in tell" onSubmit={(e) => { e.preventDefault(); if (d.text.trim()) d.send(a) }}>
      <textarea autoFocus className="f1" rows={rowsOf(d.text)} value={d.text} onKeyDown={enterSubmits} onChange={(e) => d.setText(e.target.value)}
        placeholder={`Tell ${nameOf(a)}…`} aria-label={`Text for ${nameOf(a)}`} />
      <button type="submit" className="btn p" disabled={!d.text.trim()}>Send</button>
    </form>
  )
}

// The permission prompt of a waiting session: the exact command and one button per option.
// A work push or MR shows only "Open in approvals" (item_175).
function Ask({ d, a }: { d: D; a: Agent }) {
  const p = a.prompt
  if (!p) return null
  const view = promptView(p)
  if (view === 'questions') {
    return <AskDock key={p.id} who={nameOf(a)} questions={p.questions!} busy={d.asking}
      onSend={(x) => d.answerAsk(a, x)} onCancel={() => d.answer(a, 'deny')} />
  }
  return (
    <div className="pq" role="group" aria-label={`Permission asked by ${nameOf(a)}`}>
      <div className="pqh"><span className="tool c-n">{p.tool}</span>{p.description && <span className="sm f1">{p.description}</span>}</div>
      <code className="pqc">{p.detail}</code>
      <div className="btnrow">
        {view === 'approval'
          ? <button type="button" className="btn p" disabled title="Approvals arrive with the approvals screen">Open in approvals</button>
          : p.options.map((o) => (
            <button type="button" key={o.key} className={`btn${o.key === 'allow' ? ' p' : ' x'}`} onClick={() => d.answer(a, o.key)}>{o.label}</button>
          ))}
      </div>
      {p.work_approval && <div className="sm">Work repo {p.work_approval}: give it on the exact text in approvals.</div>}
    </div>
  )
}

function Timeline({ d, a, term }: { d: D; a: Agent; term?: boolean }) {
  const lines = linesOf(a, d.snap.now)
  if (term) {
    return (
      <div className="term">
        {lines.map((l, i) => (
          <div key={i}><span className={i === lines.length - 1 && a.state === 'working' ? 'c-w' : 'c-d'}>{i === lines.length - 1 && a.state === 'working' ? '›' : '✓'}</span>{' '}
            <span className={l.says ? 'c-m' : 'c-b'}>{l.tool}</span> {l.text}</div>
        ))}
      </div>
    )
  }
  return (
    <div className="tls">
      {lines.map((l, i) => (
        <div key={i} className="tl"><b>{l.t ?? ''}</b><span className={`tool ${l.says ? 'c-m' : 'c-b'}`}>{l.tool}</span><span className="f1">{l.text}</span></div>
      ))}
    </div>
  )
}

function Forms({ d }: { d: D }) {
  if (d.form === 'run') {
    const r = d.run
    return (
      <form className="form msg cr" onSubmit={(e) => { e.preventDefault(); d.dispatch() }} aria-label="New run">
        <div className="fh"><b>Background · /dispatch</b><button type="button" className="bare sm" onClick={() => d.setForm(null)}>Close</button></div>
        <label>Project<select className="field" value={r.project} onChange={(e) => d.setRun({ ...r, project: e.target.value })}>
          <option value="">Pick a project</option>{d.projects.map((p) => <option key={p}>{p}</option>)}
        </select></label>
        <label>Label<input className="field" value={r.label} onChange={(e) => d.setRun({ ...r, label: e.target.value })} placeholder="short-name" /></label>
        <label>Model<select className="field" value={r.model} onChange={(e) => d.setRun({ ...r, model: e.target.value })}>
          {MODELS.map((m) => <option key={m}>{m}</option>)}
        </select></label>
        <AccountField value={r.account} accounts={d.accounts} set={(account) => d.setRun({ ...r, account })} />
        <label>Prompt<textarea className="field" value={r.prompt} onChange={(e) => d.setRun({ ...r, prompt: e.target.value })} /></label>
        <button type="submit" className="btn p" disabled={!r.project || !r.label.trim() || !r.prompt.trim()}>Review run</button>
      </form>
    )
  }
  if (d.form === 'session') {
    const s = d.start
    return (
      <form className="form msg cr" onSubmit={(e) => { e.preventDefault(); d.open() }} aria-label="New session">
        <div className="fh"><b>Now · /session</b><button type="button" className="bare sm" onClick={() => d.setForm(null)}>Close</button></div>
        <label>Project<select className="field" value={s.project} onChange={(e) => d.setStart({ ...s, project: e.target.value })}>
          <option value="">Pick a project</option>{d.projects.map((p) => <option key={p}>{p}</option>)}
        </select></label>
        <AccountField value={s.account} accounts={d.accounts} set={(account) => d.setStart({ ...s, account })} />
        <label>First prompt<input className="field" value={s.prompt} onChange={(e) => d.setStart({ ...s, prompt: e.target.value })} placeholder="optional" /></label>
        <button type="submit" className="btn p" disabled={!s.project}>Review session</button>
      </form>
    )
  }
  return null
}

// Reads the session transcript while a Tell it send is open, then reports delivered or
// not delivered and unmounts, which ends the reads.
function Watch({ p, onDone }: { p: Pending; onDone: (s: Delivery) => void }) {
  const conv = useConversation({ pid: p.pid })
  const s = deliveryOf(p, conv.entries, useDeliveryClock([p]))
  const done = useRef(onDone)
  done.current = onDone
  useEffect(() => { if (s !== 'sending') done.current(s) }, [s])
  return null
}

function Msgs({ d }: { d: D }) {
  return (
    <div className="msgs" role="log" aria-live="polite">
      {d.log.length === 0 && <div className="msg cr">Ask CR-0W, or type /dispatch or /session. Each action asks you to confirm first.</div>}
      {d.log.map((m, i) => (
        <Fragment key={i}>
          <div className={`msg ${m.me ? 'me' : 'cr'}${m.f ? ' f' : ''}`}>{m.text}</div>
          {m.p && m.state === 'sending' && <Watch p={m.p} onDone={(s) => d.settle(m.p!, s)} />}
          {m.p && m.state && <DeliveryLine state={m.state} onRetry={() => d.retell(m.p!)} />}
        </Fragment>
      ))}
    </div>
  )
}

// The butler chat, folded to one line by default; the line shows the last answer.
function Butler({ d }: { d: D }) {
  const last = d.log[d.log.length - 1]
  if (!d.chat) {
    return (
      <button type="button" className="bfold" aria-expanded={false} onClick={() => d.setChat(true)}>
        <span className="c-m">◇ Butler</span><span className="dim">· c</span>{last && <span className={`f1 sub${last.f ? ' c-f' : ''}`}>{last.text}</span>}
      </button>
    )
  }
  return (
    <section className="bopen" aria-label="Butler">
      <div className="lbl"><span>◇ Butler</span><button type="button" className="bare sm" aria-expanded onClick={() => d.setChat(false)}>Fold · c</button></div>
      <Msgs d={d} />
      <form className="in ask" onSubmit={(e) => { e.preventDefault(); d.submitAsk() }}>
        <input className="f1" value={d.ask} onChange={(e) => d.setAsk(e.target.value)} placeholder="Ask CR-0W: run, queue, explain…" aria-label="Ask CR-0W" />
        <button type="button" className="mic" disabled title={`Voice: ${LATER.toLowerCase()}`} aria-label="Voice"><Mic /></button>
        <button type="submit" className="btn p">Send</button>
      </form>
    </section>
  )
}

// NOW opens a session, BACKGROUND a run; AFK BATCH waits on the core route.
function Launcher({ d }: { d: D }) {
  return (
    <div className="launch">
      <form className="lrow" onSubmit={(e) => { e.preventDefault(); d.launch() }} aria-label="Launch">
        <input id="dk-new" className="field f1" value={d.task} onChange={(e) => d.setTask(e.target.value)} placeholder="new task… project, prompt" aria-label="New task" />
        <div className="segs" role="group" aria-label="Launch mode">
          {MODES.map(([m, label]) => (
            <button type="button" key={m} className={`seg${d.mode === m ? ' on' : ''}`} aria-pressed={d.mode === m} disabled={m === 'afk'}
              title={m === 'afk' ? `Arrives with the core AFK batch route (${CORE})` : undefined} onClick={() => d.setMode(m)}>{label}</button>
          ))}
        </div>
        <button type="submit" className="btn p sm">Launch</button>
      </form>
      <div className="lnote dim">AFK batch arrives with the core AFK batch route ({CORE}).</div>
      <Forms d={d} />
    </div>
  )
}

function Summary({ p }: { p: Panes }) {
  const s = summaryOf(p)
  return <span className="sum"><span className={s.needs ? 'c-n' : undefined}>{s.needs} need you</span> · {s.moving} moving · {s.queued} queued</span>
}

// ---- panes ----
const EMPTY: Record<PaneId, string> = {
  needs: 'Nothing waits on you, Sir.', moving: 'Nothing runs.', queued: 'Nothing queued.', done: 'Nothing finished yet.',
}

function Unit({ d, w, child }: { d: D; w: WorkItem; child?: boolean }) {
  const on = d.unit?.id === w.id
  return (
    <button type="button" className={`mrow unit ${w.state}${child ? ' child' : ''}${on ? ' on' : ''}`} aria-current={on || undefined} onClick={() => d.choose(w)}>
      <span className={`ki ${stateClass(w.state)}`}><KindIcon kind={w.kind} /></span>
      <span className="ti"><b>{w.title}</b><span className="sub">{w.project}{w.why ? ` · ${w.why}` : ''}</span></span>
      <span className="st">{workWhen(w, d.snap.now)}</span>
    </button>
  )
}

function Pane({ d, id, phone }: { d: D; id: PaneId; phone?: boolean }) {
  const items = d.panes[id]
  const n = PANE_IDS.indexOf(id) + 1
  const on = d.focus === id
  if (id === 'queued' && d.folded && !phone) {
    return (
      <Deck n={n} title={PANE_TITLE[id]} cls={`p-${id} folded`} on={on} aside={<span className="num">{items.length} · folded ▸</span>}>
        <button type="button" className="mrow qline" aria-expanded={false} onClick={() => { d.setFolded(false); d.setFocus('queued') }}>
          <span className="ti sub">{items.length ? items.map((w) => w.title).join(' · ') : EMPTY.queued}</span>
        </button>
      </Deck>
    )
  }
  return (
    <Deck n={n} title={PANE_TITLE[id]} cls={`p-${id}`} on={on} count={items.length}>
      <div className="mlist">
        {items.length === 0 && <div className="empty dim">{EMPTY[id]}</div>}
        {paneGroups(items).map((g) => (
          <Fragment key={g.item.id}>
            <Unit d={d} w={g.item} />
            {g.children.map((c) => <Unit key={c.id} d={d} w={c} child />)}
          </Fragment>
        ))}
      </div>
    </Deck>
  )
}

// ---- detail ----
function Facts({ d, w, a }: { d: D; w: WorkItem; a?: Agent }) {
  const f = d.finished.find((x) => x.run_id === w.source.ref)
  const m = d.missionOf(w)
  const item = itemRef(w.title) ?? itemRef(a?.label) ?? itemRef(a?.last_prompt) ?? itemRef(m?.goal)
  return (
    <dl className="facts">
      <div><dt>State</dt><dd className={stateClass(w.state)}>{a?.state === 'setup' ? said(a) : LABEL[w.state]}</dd></div>
      <div><dt>When</dt><dd>{workWhen(w, d.snap.now) || '—'}</dd></div>
      <div><dt>Account</dt><dd>{a?.account ?? f?.account ?? '—'}</dd></div>
      <div><dt>Last action</dt><dd>{a?.last_action ?? '—'}</dd></div>
      <div><dt>Item</dt><dd>{item ?? '—'}</dd></div>
    </dl>
  )
}

function Overview({ d, w, a }: { d: D; w: WorkItem; a?: Agent }) {
  return (
    <>
      <Facts d={d} w={w} a={a} />
      {w.state === 'failed' && <div className="why c-f">{w.why ? `Failed: ${w.why}.` : 'Failed. The cause is not recorded.'}</div>}
      {w.state !== 'failed' && w.why && <div className="why">{w.why}</div>}
      {a?.prompt && <Ask d={d} a={a} />}
      {d.tell && a?.kind === 'session' && <TellBox d={d} a={a} />}
      {a?.last_prompt && <div><div className="lbl"><span>Last prompt</span></div><div className="quote">{a.last_prompt}</div></div>}
      <div className="note dim">Context, commits and files arrive with the core fields ({CORE}).</div>
    </>
  )
}

const NO_TRANSCRIPT: Record<WorkItem['kind'], string> = {
  session: 'This session is gone. Its conversation can be resumed from Parked.',
  run: 'This run has ended. Its commits stay on the run branch.',
  parked: 'Parked on a usage limit. Resume it to read the transcript.',
  mission: 'A mission has no transcript of its own. Open one of its runs.',
  mr: 'An MR has no transcript. Open it on the forge.',
}

function MissionRuns({ m }: { m: Mission }) {
  if (!m.runs.length) return <div className="empty dim">No run yet.</div>
  return (
    <div className="tls">
      {m.runs.map((r) => (
        <div key={r.run} className="tl"><b>{time(new Date(r.started))}</b><span className="tool c-b">{r.run.slice(0, 8)}</span>
          <span className="f1">{r.account}{r.ended ? ` · ended ${time(new Date(r.ended))}` : ' · running'}{r.cause ? ` · ${r.cause}` : ''}</span></div>
      ))}
    </div>
  )
}

function Body({ d, w, a }: { d: D; w: WorkItem; a?: Agent }) {
  if (d.tab === 'transcript') {
    return a
      ? <AgentChat key={a.pid} snap={d.snap} refresh={d.refresh} target={{ pid: a.pid }} log={<Timeline d={d} a={a} term />} />
      : <div className="empty dim">{NO_TRANSCRIPT[w.kind]}</div>
  }
  if (d.tab === 'timeline') {
    const m = d.missionOf(w)
    if (a) return <Timeline d={d} a={a} />
    if (m) return <MissionRuns m={m} />
    return <div className="empty dim">{w.since ? `${LABEL[w.state]} since ${time(new Date(w.since))}.` : 'No timeline for this unit.'}</div>
  }
  if (d.tab === 'files') return <div className="empty dim">Files arrive with the core fields ({CORE}).</div>
  return <Overview d={d} w={w} a={a} />
}

function Detail({ d, w, phone }: { d: D; w: WorkItem; phone?: boolean }) {
  const a = d.agent
  // The session chat carries its own reply box, Pause and Stop; the footer keeps the rest.
  const inChat = d.tab === 'transcript' && a?.kind === 'session'
  const acts = d.actsOf(w).filter((x) => !inChat || !CHAT_ACTS.includes(x.label))
  const pane = d.paneOf(w)
  const meta = [KIND_LABEL[w.kind].toLowerCase(), w.project, a?.account, a?.cwd, a?.started_at && `started ${time(new Date(a.started_at))}`]
    .filter(Boolean).join(' · ')
  const inner = (
    <div className={`mdet${phone ? ' ph' : ''}`}>
      {phone && pane && <button type="button" className="mback" onClick={() => d.setPage(false)}>◀ [{PANE_IDS.indexOf(pane) + 1}] {PANE_TITLE[pane]}</button>}
      <h2><span className={`ki ${stateClass(w.state)}`}><KindIcon kind={w.kind} size={15} /></span> {w.title}</h2>
      <div className="meta">{meta}</div>
      <div className="mtabs" role="tablist" aria-label="Detail">
        {TABS.map(([t, label]) => (
          <button type="button" role="tab" key={t} className={`mtab${d.tab === t ? ' on' : ''}`} aria-selected={d.tab === t} onClick={() => d.setTab(t)}>{label}</button>
        ))}
      </div>
      <div className={`mbody tab-${d.tab}`}><Body d={d} w={w} a={a} /></div>
      {acts.length > 0 && (
        <div className="macts">
          <div className="btnrow">
            {acts.map((x) => (
              <button type="button" key={x.label} className={`btn sm${x.p ? ' p' : ''}${x.x ? ' x' : ''}`} onClick={x.run}>
                {x.label}{x.key && !phone && <span className="kk">{x.key}</span>}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
  if (phone) return inner
  return <Deck n={5} title={w.title} cls="p-detail" on={d.focus === 'detail'} aside={<span>{w.project} · {KIND_LABEL[w.kind].toLowerCase()}</span>}>{inner}</Deck>
}

const KEYS: [string, string][] = [['1-5', 'pane'], ['j k', 'move'], ['n', 'new task'], ['a', 'answer'], ['r', 'retry'], ['x', 'cancel'], ['o', 'open session'], ['c', 'chat']]

// ---- laptop ----
function DeskLaptop({ d }: { d: D }) {
  return (
    <div className="desk dk m-lp">
      <div className="mhead">
        <h1>Desk</h1>
        <div className="segs" role="group" aria-label="Project">
          {['', ...d.allProjects].map((p) => (
            <button type="button" key={p || 'all'} className={`seg${d.proj === p ? ' on' : ''}`} aria-pressed={d.proj === p} onClick={() => d.setProj(p)}>{p || 'All'}</button>
          ))}
        </div>
        <Summary p={d.panes} />
      </div>
      <Launcher d={d} />
      <div className="mgrid">
        <div className="mpanes">{PANE_IDS.map((id) => <Pane key={id} d={d} id={id} />)}</div>
        {d.unit
          ? <Detail d={d} w={d.unit} />
          : <Deck n={5} title="Detail" cls="p-detail"><div className="empty dim">Nothing on the desk, Sir. Launch a task above.</div></Deck>}
      </div>
      <div className="dfoot">
        <div className="keys dim" aria-label="Keys">{KEYS.map(([k, what]) => <span key={k}><b>{k}</b> {what}</span>)}</div>
        <Butler d={d} />
      </div>
    </div>
  )
}

const CHAT_ACTS = ['Answer', 'Allow', 'Message', 'Pause', 'Stop']

// ---- phone: pane tabs, one pane, the detail as its own page ----
const SHORT: Record<PaneId, string> = { needs: 'Needs', moving: 'Moving', queued: 'Queued', done: 'Done' }
function DeskPhone({ d }: { d: D }) {
  if (d.page && d.unit) return <div className="desk dk m-ph"><Detail d={d} w={d.unit} phone /></div>
  const tab: PaneId = d.focus === 'detail' ? 'needs' : d.focus
  return (
    <div className="desk dk m-ph">
      <div className="mhead"><h1>Desk</h1><Summary p={d.panes} /></div>
      <div className="segs ptabs" role="tablist" aria-label="Panes">
        {PANE_IDS.map((id, i) => (
          <button type="button" role="tab" key={id} className={`seg${tab === id ? ' on' : ''}`} aria-selected={tab === id} onClick={() => d.setFocus(id)}>
            <span className="pk">[{i + 1}] </span>{SHORT[id]} <span className="num">{d.panes[id].length}</span>
          </button>
        ))}
      </div>
      <Pane d={d} id={tab} phone />
      <Launcher d={d} />
      <Butler d={d} />
    </div>
  )
}

export function Desk({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const d = useDesk(snap, refresh)
  const lp = useMedia(LAPTOP)
  const keyRef = useRef(d.key)
  keyRef.current = d.key
  const sheetOpen = !!d.sheet
  useEffect(() => {
    if (!lp || sheetOpen) return
    const on = (e: KeyboardEvent) => {
      const t = e.target as HTMLElement | null
      if (e.metaKey || e.ctrlKey || e.altKey || t?.closest?.('input, textarea, select, [contenteditable="true"]')) return
      if (keyRef.current(e.key)) e.preventDefault()
    }
    window.addEventListener('keydown', on)
    return () => window.removeEventListener('keydown', on)
  }, [lp, sheetOpen])
  return (
    <>
      {lp ? <DeskLaptop d={d} /> : <DeskPhone d={d} />}
      {d.sheet && <Sheet c={d.sheet} onClose={() => d.setSheet(null)} />}
    </>
  )
}
