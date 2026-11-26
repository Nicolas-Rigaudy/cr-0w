import { useEffect, useRef, useState, type ReactNode } from 'react'
import { ApiError } from '../data/api'
import { answerOf, designOf, doneToday, isDesign, isFyi, kindOf, kindWord, KINDS, type AnswerRow, type Design, type InboxKind } from '../data/inbox'
import type { Mail } from '../data/live'
import type { GateChoice, GateQuestion, InboxItem, InboxVerb, Mission } from '../data/types'
import { useSource, type Snapshot } from '../data/useData'
import type { WorkItem } from '../data/work'
import { time, title } from '../format'
import { Link } from '../router'
import { REJECT, useAnswer, verbLabel } from '../ui/answer'
import { Sheet, type Confirm } from '../ui/Sheet'
import { useMedia } from '../ui/useMedia'
import { TAG_ICON } from '../ui/WorkList'

// Inbox (item_223 and item_191, board 223-B approved 2026-11-25): everything that waits on
// the owner, by kind. Laptop: kind rail with the Reports log, list, reading pane. Fold and
// phone: kind chips over the list. On the phone a row opens a sheet; swiping it left or
// right moves to the next or previous row. FYI reports live in the read-only Reports log.
const LAPTOP = '(min-width: 1100px)'
const INNER = '(min-width: 600px)'
const short = (id: string) => id.replace(/^\d{8}-/, '')
const valid = (iso?: string | null) => !!iso && !isNaN(Date.parse(iso))
const LETTERS = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'

type Filter = 'all' | InboxKind

// ---- kind icons (board 223-B) ----
const ICON: Record<InboxKind | 'design' | 'doc' | 'done', ReactNode> = {
  decision: <><circle cx="11" cy="4.5" r="2" /><circle cx="5" cy="17.5" r="2" /><circle cx="17" cy="17.5" r="2" /><path d="M11 6.5v4l-6 5M11 10.5l6 5" /></>,
  review: <><rect x="3" y="4" width="7" height="14" rx="1.5" /><rect x="12" y="4" width="7" height="14" rx="1.5" /></>,
  design: <><rect x="3" y="4" width="7" height="14" rx="1.5" /><rect x="12" y="4" width="7" height="14" rx="1.5" /></>,
  doc: <><path d="M5.5 2.5h7l4 4v13h-11z" /><path d="M12.5 2.5v4h4M8 11h6M8 14h6M8 17h3.5" /></>,
  approval: <path d="m5 11.5 4 4 8-9" />,
  done: <path d="m5 11.5 4 4 8-9" />,
  ask: TAG_ICON.question,
  mail: TAG_ICON.mail,
  followup: <><circle cx="11" cy="11" r="8.5" /><path d="M11 6.5V11l3 2" /></>,
}
function Icon({ k }: { k: keyof typeof ICON }) {
  return (
    <svg aria-hidden="true" width="16" height="16" viewBox="0 0 22 22" fill="none" stroke="currentColor" strokeWidth="1.9" strokeLinejoin="round" strokeLinecap="round">{ICON[k]}</svg>
  )
}
const iconOf = (r: AnswerRow): keyof typeof ICON => {
  const k = kindOf(r)
  if (k === 'review') return kindWord(r) === 'design review' ? 'design' : 'doc'
  return k
}

// Where a waiting work item opens: its mission, the MRs, or the Desk for live work.
export function workHref(w: WorkItem) {
  if (w.kind === 'mission') return `/missions/${encodeURIComponent(w.source.ref)}`
  if (w.kind === 'mr') return '/mrs'
  if (w.mission) return `/missions/${encodeURIComponent(w.mission)}`
  return '/desk'
}

// Row text: title and clock time over kind · source · status (board 223-B).
interface Line { title: string; sub: string[]; at?: string }
function lineOf(r: AnswerRow, missions: Mission[]): Line {
  const at = valid(r.at) ? time(r.at!) : undefined
  const k = kindWord(r)
  if (r.kind === 'prompt') {
    const w = r.item
    return { title: w.title, sub: [k, w.kind, w.project, ...(w.why ? [w.why] : [])], at }
  }
  if (r.kind === 'mail') {
    const m = r.item
    return { title: m.subject, sub: [k, m.from ?? 'Outlook', m.flag ?? 'flagged'] }
  }
  const q = r.item
  const m = missions.find((x) => x.id === q.mission)
  return { title: q.ask, sub: [k, m ? title(m) : short(q.mission), q.project], at }
}

// The options of a row: the question it asks, else the first question of its version
// that has choices.
function optionsOf(qs: GateQuestion[], r: InboxItem): { q?: GateQuestion; choices: GateChoice[] } {
  const qid = r.source === 'questions' ? r.id.slice(r.id.lastIndexOf('/') + 1) : undefined
  const q = qid ? qs.find((x) => x.id === qid) : qs.find((x) => x.choices?.length)
  return { q, choices: q?.choices ?? [] }
}
// A pick answers the row with its own pick verb, else the verb that resumes the asker.
const pickVerb = (verbs: InboxVerb[]): InboxVerb | undefined =>
  (['pick', 'answer', 'msg', 'go', 'ok', 'approve'] as InboxVerb[]).find((v) => verbs.includes(v))

function useQuestions(r: InboxItem | undefined) {
  const src = useSource()
  const [qs, setQs] = useState<{ key: string; list: GateQuestion[] }>({ key: '', list: [] })
  const key = r ? `${r.mission}/v${r.version}` : ''
  useEffect(() => {
    if (!r || !src.questions || key === qs.key) return
    let live = true
    src.questions(r.mission, r.version).then((list) => { if (live) setQs({ key, list }) }, () => { if (live) setQs({ key, list: [] }) })
    return () => { live = false }
  }, [key]) // eslint-disable-line react-hooks/exhaustive-deps
  return qs.key === key ? qs.list : []
}

function useInbox(snap: Snapshot, refresh: () => Promise<void>) {
  const { inbox, missions } = snap
  const act = useSource().actions
  const [filter, setFilter] = useState<Filter>('all')
  const [sheet, setSheet] = useState<Confirm | null>(null)
  const rows = answerOf(snap)
  // `?item=<inbox id>` opens that row (a link from a notification or a report).
  const [linked] = useState(() => {
    const id = new URLSearchParams(globalThis.location?.search ?? '').get('item')
    return id ? rows.find((r) => r.kind === 'ask' && r.item.id === id)?.key : undefined
  })
  const [open, setOpen] = useState(!!linked)
  const [choice, setChoice] = useState<string | undefined>()
  const shown = filter === 'all' ? rows : rows.filter((r) => kindOf(r) === filter)
  const done = filter === 'all' ? doneToday(inbox, snap.now) : []
  const fyis = inbox.filter(isFyi)
  const count = (k: InboxKind) => rows.filter((r) => kindOf(r) === k).length
  const [sel, setSelRaw] = useState<string | undefined>(linked ?? shown[0]?.key)
  const picked = shown.find((r) => r.key === sel) ?? shown[0]
  const q = picked?.kind === 'ask' ? picked.item : undefined
  const a = useAnswer(q, inbox, missions, refresh)
  const qs = useQuestions(q)
  // A design review offers its option letters (item_211); other rows their version questions.
  const design = q && isDesign(q) ? designOf(q) : undefined
  const opts: { q?: GateQuestion; choices: GateChoice[] } = design
    ? { choices: design.letters.map((l) => ({ id: l, label: design.label(l) ?? `Option ${l}` })) }
    : q ? optionsOf(qs, q) : { choices: [] }
  const chosen = opts.choices.find((c) => c.id === choice) ?? opts.choices.find((c) => c.id === (opts.q?.recommended ?? design?.recommend)) ?? opts.choices[0]
  const setSel = (k: string) => { setSelRaw(k); setChoice(undefined) }
  const step = (d: 1 | -1) => {
    const n = shown.findIndex((r) => r.key === picked?.key)
    const to = shown[n + d]
    if (to) setSel(to.key)
  }
  const line = (r: AnswerRow) => lineOf(r, missions)

  const hide = (r: InboxItem) => setSheet({
    title: 'Hide this item?', lead: 'It leaves the inbox without a reply.',
    target: [['item', short(r.id)], ['project', r.project]],
    note: 'Same as crow inbox hide. Logged with surface nest.', ok: 'Hide',
    run: async () => {
      try {
        await act.inboxHide(r.id)
      } catch (e) {
        if (e instanceof ApiError && e.status === 409) void refresh()
        throw e
      }
      await refresh()
    },
  })
  // Same sheet the Desk had for flagged mail: it changes the Outlook flag only.
  const flag = (m: Mail, verb: 'done' | 'snooze') => setSheet({
    title: verb === 'done' ? 'Mark done' : 'Snooze', lead: verb === 'done' ? 'Marks the Outlook flag as complete.' : 'Moves the Outlook flag to tomorrow.',
    target: [['mail', m.subject]], note: 'Changes the Outlook flag only; nothing is sent.', ok: verb === 'done' ? 'Done' : 'Snooze',
    run: async () => { await act.mail(m.id, verb); await refresh() },
  })
  return { filter, setFilter, rows, shown, done, fyis, count, sel: picked?.key, setSel, picked, q, a, design, opts, chosen, setChoice, open, setOpen, step, missions, line, hide, flag, sheet, setSheet }
}
type I = ReturnType<typeof useInbox>

// ---- list ----
function Row({ i, r }: { i: I; r: AnswerRow }) {
  const l = i.line(r)
  const on = r.key === i.sel
  return (
    <button type="button" className={`irow${on ? ' on' : ''}`} aria-pressed={on} onClick={() => { i.setSel(r.key); i.setOpen(true) }}
      aria-label={`${l.title}, ${l.sub.join(', ')}${l.at ? `, ${l.at}` : ''}`}>
      <span className={`ic k-${kindOf(r)}`}><Icon k={iconOf(r)} /></span>
      <span className="t"><b>{l.title}</b><span>{l.sub.join(' · ')}</span></span>
      {l.at && <span className="tm num">{l.at}</span>}
    </button>
  )
}

function DoneRow({ r }: { r: InboxItem }) {
  const at = time(r.answered!)
  return (
    <div className="irow done">
      <span className="ic k-done"><Icon k="done" /></span>
      <span className="t"><b>{r.ask}</b><span>{r.state === 'dropped' ? 'hidden' : 'answered'} {at} · {r.project}</span></span>
      <span className="tm num">{at}</span>
    </div>
  )
}

function List({ i }: { i: I }) {
  const show = (k: InboxKind) => i.filter === 'all' || i.filter === k
  return (
    <section className="pn ilist" aria-label="Inbox">
      <div className="grp c-n"><span>To answer</span><span className="num">{i.shown.length}</span></div>
      <div className="alist">
        {i.shown.map((r) => <Row key={r.key} i={i} r={r} />)}
        {i.shown.length === 0 && <div className="sm empty">Nothing to answer.</div>}
      </div>
      {show('followup') && (
        <>
          <div className="grp c-q"><span>Waiting on others</span><span className="num">0</span></div>
          <div className="sm empty">No follow-up tracked yet.</div>
        </>
      )}
      {i.done.length > 0 && (
        <>
          <div className="grp c-d"><span>Done today</span><span className="num">{i.done.length}</span></div>
          <div className="alist">{i.done.map((r) => <DoneRow key={r.id} r={r} />)}</div>
        </>
      )}
    </section>
  )
}

function Rail({ i }: { i: I }) {
  const item = (k: Filter, label: string, n: number) => (
    <button key={k} type="button" className={`kind${i.filter === k ? ' on' : ''}`} aria-pressed={i.filter === k} onClick={() => i.setFilter(k)}>
      <span>{label}</span><span className="num">{n}</span>
    </button>
  )
  return (
    <nav className="pn rail" aria-label="Kinds">
      <div className="grp"><span>Kinds</span></div>
      {item('all', 'All', i.rows.length)}
      {KINDS.map(([k, l]) => item(k, l, i.count(k)))}
      <div className="grp"><span>Log</span></div>
      <Link className="kind log" href="/reports"><span><b>Reports</b><span className="sm">read-only, {i.fyis.length}</span></span><span className="chev">›</span></Link>
    </nav>
  )
}

function Chips({ i }: { i: I }) {
  const chip = (k: Filter, label: string, n: number) => (
    <button key={k} type="button" className={`ichip${i.filter === k ? ' on' : ''}`} aria-pressed={i.filter === k} onClick={() => i.setFilter(k)}>{label} {n}</button>
  )
  return (
    <div className="ichips" role="group" aria-label="Kinds">
      {chip('all', 'All', i.rows.length)}
      {KINDS.map(([k, l]) => chip(k, l, i.count(k)))}
    </div>
  )
}

// ---- reader ----
function Route({ i }: { i: I }) {
  const p = i.picked!
  if (p.kind === 'mail') return <div className="route">→ changes the Outlook flag only · nothing is sent</div>
  if (p.kind === 'prompt') return <div className="route">→ opens the waiting {p.item.kind}</div>
  const m = i.missions.find((x) => x.id === p.item.mission)
  return <div className="route">→ answer goes to mission {m ? title(m) : short(p.item.mission)} · it resumes</div>
}

function Answer({ i }: { i: I }) {
  const p = i.picked!
  if (p.kind === 'mail') {
    return (
      <div className="btnrow">
        <button type="button" className="btn p" onClick={() => i.flag(p.item, 'done')}>Done</button>
        <button type="button" className="btn" onClick={() => i.flag(p.item, 'snooze')}>Snooze</button>
      </div>
    )
  }
  if (p.kind === 'prompt') return <div className="btnrow"><Link className="btn p" href={workHref(p.item)}>Open</Link></div>
  const item = i.a.item
  const r = p.item
  if (!item || !i.a.mission) return null
  const pv = i.opts.choices.length ? pickVerb(item.verbs) : undefined
  // A design review answers with the letter and, from a compare spec, the option label.
  const lead = (c: GateChoice) => `${i.opts.q?.id ?? 'pick'}: ${c.id}${i.design?.label(c.id) ? ` (${c.label})` : ''}`
  const letter = (c: GateChoice) => LETTERS[i.opts.choices.indexOf(c)]
  const rest = item.verbs.filter((v) => v !== pv)
  return (
    <>
      <input className="inote" placeholder="Change or note (optional)" aria-label="Change or note"
        value={i.a.text} onChange={(e) => i.a.setText(e.target.value)} disabled={i.a.busy} />
      <div className="btnrow">
        {pv && i.chosen && (
          <button type="button" className="btn p" onClick={() => i.a.send(pv, lead(i.chosen!))} disabled={i.a.busy} aria-busy={i.a.busy}>
            {i.a.busy ? 'Sending…' : `Pick ${letter(i.chosen)}`}
          </button>
        )}
        {pv && !i.design && i.opts.choices.filter((c) => c !== i.chosen).map((c) => (
          <button key={c.id} type="button" className="btn" onClick={() => i.setChoice(c.id)} disabled={i.a.busy}>{letter(c)}</button>
        ))}
        {rest.map((v, n) => (
          <button key={v} type="button" className={`btn${!pv && n === 0 ? ' p' : REJECT.includes(v) ? ' x' : ''}`}
            onClick={() => i.a.send(v)} disabled={i.a.busy} aria-busy={i.a.busy}>{!pv && n === 0 && i.a.busy ? 'Sending…' : verbLabel(v)}</button>
        ))}
        {r.state === 'open' && <button type="button" className="btn x hide" onClick={() => i.hide(r)} disabled={i.a.busy}>Hide</button>}
      </div>
    </>
  )
}

function Reader({ i, back }: { i: I; back?: () => void }) {
  const p = i.picked
  if (!p) return <section className="pn read"><div className="sm">Inbox is empty.</div></section>
  const l = i.line(p)
  const d = i.design
  const meta = d ? [...l.sub, `${d.letters.length} options × ${d.scenarios.length} scenario${d.scenarios.length > 1 ? 's' : ''}`]
    : p.kind === 'ask' ? [...l.sub, `v${p.item.version}`, p.item.file] : l.sub
  return (
    <section className="pn read" aria-labelledby="rt" aria-live="polite">
      {back && <button type="button" className="btn back" onClick={back}>← Inbox</button>}
      <h2 id="rt" className="rt">{l.title}</h2>
      <div className="meta">{meta.join(' · ')}</div>
      <Body i={i} />
      <div className="foot">
        {i.a.error && <div className="err" role="alert">{i.a.error}</div>}
        {i.a.newer && <div className="stale" role="alert"><b>THE MISSION MOVED ON · NOW V{i.a.mission?.version}</b><span className="sm">Your text is kept.</span></div>}
        <Answer i={i} />
        <Route i={i} />
      </div>
    </section>
  )
}

function Body({ i }: { i: I }) {
  const p = i.picked!
  if (p.kind === 'mail') return <div className="body">from {p.item.from ?? 'unknown'}{p.item.scope && <> · {p.item.scope} mailbox</>}. Crow does not read the body yet. Open it in Outlook.</div>
  if (p.kind === 'prompt') return <div className="body">{p.item.why ?? 'Waiting on you.'}</div>
  const r = i.a.item ?? p.item
  const { q, choices } = i.opts
  if (i.design) return <><Stage i={i} r={r} d={i.design} /><div className="body">{r.ask}</div></>
  return (
    <>
      <h3 className="sh">Context</h3>
      <div className="body">{q && q.text !== r.ask ? `${r.ask}\n\n${q.text}` : r.ask}</div>
      {choices.length > 0 && (
        <>
          <h3 className="sh">Options</h3>
          <div className="opts" role="radiogroup" aria-label="Options">
            {choices.map((c, n) => (
              <button key={c.id} type="button" role="radio" aria-checked={c === i.chosen}
                className={`opt${c === i.chosen ? ' on' : ''}${c.id === q?.recommended ? ' rec' : ''}`} onClick={() => i.setChoice(c.id)}>
                <b>{LETTERS[n]}</b><span className="f1">{c.label}{c.note && <span className="sm"> {c.note}</span>}</span>
                {c.id === q?.recommended && <span className="tag">recommended</span>}
              </button>
            ))}
          </div>
        </>
      )}
    </>
  )
}

// ---- flip stage (item_211, board 211-B): one option at a time, flip through the letters ----
function Stage({ i, r, d }: { i: I; r: InboxItem; d: Design }) {
  const src = useSource()
  const [scenario, setScenario] = useState(d.scenarios[0] ?? '')
  const [zoom, setZoom] = useState(false)
  const on = i.chosen?.id ?? d.letters[0]
  const flip = () => i.setChoice(d.letters[(d.letters.indexOf(on) + 1) % d.letters.length])
  const url = (l: string) => {
    const f = d.file(l, scenario)
    return f ? src.artifact?.(r.mission, f) : undefined
  }
  const img = url(on)
  return (
    <div className={`stage${zoom ? ' zoom' : ''}`}>
      <button type="button" className="shot" onClick={flip} aria-label={`Option ${on}${scenario ? `, ${scenario}` : ''}. Flip to the next option`}>
        {img ? <img src={img} alt="" /> : <span className="sm">No image for this option.</span>}
        <b className="lt">{on}</b>
        <span className="hint lp">Click flips</span>
      </button>
      {d.label(on) && (
        <div className="cap">
          <b>{on} · {d.label(on)}</b>{on === d.recommend && <span className="tag">recommended</span>}
          {d.note(on) && <span className="sm">{d.note(on)}</span>}
        </div>
      )}
      <div className="side">
        <button type="button" className="btn zoomb" onClick={() => setZoom((z) => !z)} aria-pressed={zoom}>{zoom ? 'Fit image' : 'Zoom image'}</button>
        <div className="thumbs" role="radiogroup" aria-label="Options">
          {d.letters.map((l) => (
            <button key={l} type="button" role="radio" aria-checked={l === on} className={`thumb${l === on ? ' on' : ''}`} onClick={() => i.setChoice(l)}>
              {url(l) && <img src={url(l)} alt="" />}<b>{l}</b>
            </button>
          ))}
          <span className="hint ph">tap image to flip</span>
        </div>
        {d.scenarios.length > 1 && (
          <>
            <h3 className="sh">Scenario</h3>
            <div className="scen">
              {d.scenarios.map((s, n) => (
                <button key={s} type="button" aria-pressed={s === scenario} className={`kind${s === scenario ? ' on' : ''}`} onClick={() => setScenario(s)}>
                  <span className="num">{n + 1}</span><span>{s || 'default'}</span>
                </button>
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  )
}

// ---- phone sheet: swipe left for the next row, right for the previous ----
function PhoneSheet({ i }: { i: I }) {
  const start = useRef<{ x: number; y: number } | null>(null)
  useEffect(() => {
    const key = (e: KeyboardEvent) => { if (e.key === 'Escape') i.setOpen(false) }
    window.addEventListener('keydown', key)
    return () => window.removeEventListener('keydown', key)
  })
  if (!i.open || !i.picked) return null
  const n = i.shown.findIndex((r) => r.key === i.picked!.key)
  return (
    <div className="iscrim" onClick={(e) => { if (e.target === e.currentTarget) i.setOpen(false) }}>
      <div className="isheet" role="dialog" aria-modal="true" aria-labelledby="rt"
        onTouchStart={(e) => { start.current = { x: e.touches[0].clientX, y: e.touches[0].clientY } }}
        onTouchEnd={(e) => {
          const s = start.current
          start.current = null
          if (!s) return
          const dx = e.changedTouches[0].clientX - s.x
          const dy = e.changedTouches[0].clientY - s.y
          if (Math.abs(dx) > 60 && Math.abs(dx) > 1.5 * Math.abs(dy)) i.step(dx < 0 ? 1 : -1)
        }}>
        <button type="button" className="handle" aria-label="Close" onClick={() => i.setOpen(false)}><span /></button>
        <div className="pos sm num" aria-label={`Item ${n + 1} of ${i.shown.length}, swipe for the next`}>{n + 1} / {i.shown.length}</div>
        <Reader i={i} />
      </div>
    </div>
  )
}

export function Inbox({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const i = useInbox(snap, refresh)
  const laptop = useMedia(LAPTOP)
  const inner = useMedia(INNER)
  // A design review takes the page on laptop and tablet (board 211-B), next to a slim list of
  // titles; "← Inbox" or Escape brings the full list back until another row is picked.
  const [listFor, setListFor] = useState<string | null>(null)
  const key = i.picked?.key ?? null
  const wide = (laptop || inner) && !!i.design && listFor !== key
  useEffect(() => {
    if (!wide) return
    const esc = (e: KeyboardEvent) => { if (e.key === 'Escape') setListFor(key) }
    window.addEventListener('keydown', esc)
    return () => window.removeEventListener('keydown', esc)
  }, [wide, key])
  return (
    <>
      {wide ? (
        <div className="desk ibx ib-wide"><List i={i} /><Reader i={i} back={() => setListFor(key)} /></div>
      ) : laptop ? (
        <div className="desk ibx ib-lp"><Rail i={i} /><List i={i} /><Reader i={i} /></div>
      ) : inner ? (
        <div className="desk ibx ib-in"><div className="col"><Chips i={i} /><List i={i} /></div><Reader i={i} /></div>
      ) : (
        <div className="desk ibx ib-ph"><Chips i={i} /><List i={i} /><PhoneSheet i={i} /></div>
      )}
      {i.sheet && <Sheet c={i.sheet} onClose={() => i.setSheet(null)} />}
    </>
  )
}
