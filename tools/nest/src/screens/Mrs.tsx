import { useState, type ReactNode } from 'react'
import { ApiError } from '../data/api'
import type { Actions, LiveReply } from '../data/live'
import { branchPrompt, countsOf, hitsOf, mrsOf, patternsFor, onlyProject, reviewLabel, reviewPrompt, threadsPrompt, type Mr, type MrBranch, type MrJob, type MrsData } from '../data/mrs'
import { useSource, type Snapshot } from '../data/useData'
import { ago, time } from '../format'
import { useQuery } from '../router'
import { Sheet, type Confirm } from '../ui/Sheet'
import { useMedia } from '../ui/useMedia'

// The MRs page (item_225, board round 2 option C): three deck panes (to review, mine,
// merged today) and a tabbed detail of the picked MR.

const LAPTOP = '(min-width: 1100px)'
export type Pane = 'review' | 'mine' | 'merged'
export type Tab = 'overview' | 'threads' | 'pipeline' | 'patterns'
const PANES: Array<[Pane, string]> = [['review', 'To review'], ['mine', 'Mine'], ['merged', 'Merged today']]
const TABS: Tab[] = ['overview', 'threads', 'pipeline', 'patterns']
// The patterns tab keeps to the newest rows; the rest stay in the ledger files.
const PATTERN_ROWS = 8
const keyOf = (m: Mr) => `${m.project}:${m.ref}`

// Row actions (item_225): every one goes through a confirm sheet and the existing live
// actions. Posting on the forge stays in the review and review-threads sessions.
type Said = { text: string; f?: boolean }
export function confirmsOf(act: Actions, refresh: () => Promise<void>, setSaid: (s: Said) => void) {
  const wrap = (call: () => Promise<LiveReply>) => async () => {
    try {
      const r = await call()
      setSaid({ text: r.message || `${r.action} ${r.target}`, f: !r.ok })
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) void refresh()
      throw e
    }
    void refresh()
  }
  const target = (m: Mr): Array<[string, string]> => [['MR', `${m.ref} ${m.title}`], ['Project', `${m.name} · ${m.forge}`]]
  return {
    review: (m: Mr): Confirm => ({
      title: `Review ${m.ref}`, lead: 'Starts a background review run on this MR.', quote: reviewPrompt(m),
      target: [...target(m), ['Run', reviewLabel(m)]], note: 'Comments wait for your OK before they post.', ok: 'Start review',
      run: wrap(() => act.dispatch({ project: m.project, label: reviewLabel(m), prompt: reviewPrompt(m) })),
    }),
    threads: (m: Mr): Confirm => ({
      title: `Answer threads on ${m.ref}`, lead: `Opens a session on the ${m.threads_open} open threads.`, quote: threadsPrompt(m),
      target: target(m), note: 'Replies post after your OK in the threads session.', ok: 'Open session',
      run: wrap(() => act.start({ project: m.project, prompt: threadsPrompt(m) })),
    }),
    session: (m: Mr): Confirm => ({
      title: `Session in ${m.name}`, lead: 'Opens a session in the repo of this MR.', target: [['Project', m.name], ['Branch', m.branch]],
      note: 'Nothing runs until you write in the session.', ok: 'Open session',
      run: wrap(() => act.start({ project: m.project })),
    }),
    branch: (b: MrBranch): Confirm => ({
      title: `Open an MR for ${b.branch}`, lead: 'Opens a session that prepares the MR.', quote: branchPrompt(b),
      target: [['Project', b.name || b.project], ['Branch', `${b.branch} · ${b.ahead} ahead`]], note: 'The MR opens after your OK in the session.', ok: 'Open session',
      run: wrap(() => act.start({ project: b.project, prompt: branchPrompt(b) })),
    }),
  }
}

export function useMrActions(refresh: () => Promise<void>) {
  const act = useSource().actions
  const [sheet, setSheet] = useState<Confirm | null>(null)
  const [said, setSaid] = useState<Said | null>(null)
  const c = confirmsOf(act, refresh, setSaid)
  return {
    sheet, setSheet, said,
    review: (m: Mr) => setSheet(c.review(m)),
    threads: (m: Mr) => setSheet(c.threads(m)),
    session: (m: Mr) => setSheet(c.session(m)),
    branch: (b: MrBranch) => setSheet(c.branch(b)),
  }
}
export type A = ReturnType<typeof useMrActions>

export function useMrs(snap: Snapshot, refresh: () => Promise<void>) {
  const a = useMrActions(refresh)
  const q = useQuery()
  const all = mrsOf(snap.sections.mrs, snap.now)
  const [project, setProject] = useState(q.get('project') ?? '')
  const d = onlyProject(all, all.projects.includes(project) ? project : '')
  const rows = [...d.review, ...d.mine, ...d.merged]
  // ?mr=239 or ?mr=!239 picks an MR (shots, links from other pages).
  const want = q.get('mr')
  const [pick, setPick] = useState<string | null>(want ? rows.find((m) => m.ref.replace(/^[!#]/, '') === want.replace(/^[!#]/, ''))?.ref ?? null : null)
  const picked = rows.find((m) => m.ref === pick || keyOf(m) === pick) ?? d.mine[0] ?? d.review[0] ?? d.merged[0] ?? null
  const pane: Pane | null = !picked ? null : d.review.includes(picked) ? 'review' : d.mine.includes(picked) ? 'mine' : 'merged'
  const tq = q.get('tab') as Tab | null
  const [tab, setTab] = useState<Tab>(tq && TABS.includes(tq) ? tq : 'overview')
  // Phone: one pane at a time (the picked MR's pane until a pane tab is tapped), the
  // detail as its own page once a row is tapped (or named in the URL).
  const [paneTab, setPaneTab] = useState<Pane | null>(null)
  const [shown, setShown] = useState(!!want && !!picked)
  const open = (m: Mr) => { setPick(keyOf(m)); setTab('overview'); setShown(true) }
  return {
    a, all, d, project, setProject, picked, pane, tab, setTab, open, now: snap.now,
    paneTab: paneTab ?? pane ?? 'mine', setPaneTab, shown: shown && !!picked, back: () => setShown(false),
  }
}
export type M = ReturnType<typeof useMrs>

const CI: Record<Mr['pipeline'], [string, string]> = {
  passed: ['●', 'c-d'], failed: ['✕', 'c-f'], running: ['◌', 'c-b'], none: ['', ''],
}
const JOB: Record<MrJob['state'], string> = { passed: '✓', failed: '✕', running: '◌', none: '', skipped: '', manual: '' }

function Ci({ m }: { m: Mr }) {
  const [g, c] = CI[m.pipeline]
  return g ? <span className={`ci ${c}`}>{g} CI</span> : null
}

function Head({ m, phone }: { m: M; phone?: boolean }) {
  const c = countsOf(m.d)
  if (phone) {
    return (
      <div className="mhead">
        <h1>MRs</h1>
        <div className="sum">{c.needs > 0 && <><span className="c-n">{c.needs} needs you</span> · </>}<span>{c.review} review</span></div>
      </div>
    )
  }
  const seg = (p: string, label: string) => (
    <button key={p || 'all'} type="button" className={`seg${m.project === p ? ' on' : ''}`} aria-pressed={m.project === p} onClick={() => m.setProject(p)}>{label}</button>
  )
  return (
    <div className="mhead">
      <h1>MRs</h1>
      {m.all.projects.length > 1 && (
        <div className="segs" role="group" aria-label="Projects">{seg('', 'All')}{m.all.projects.map((p) => seg(p, p))}</div>
      )}
      <div className="sum">
        <span>{c.review} to review</span> · <span className={c.needs ? 'c-n' : ''}>{c.needs} needs you</span> · <span>{c.waiting} waiting</span> · <span>merged today {c.merged}</span>
      </div>
    </div>
  )
}

// A deck pane: the title sits in the top border, the count at the right.
export function Deck({ n, title, count, on, aside, cls, children }: { n: number; title: string; count?: number; on?: boolean; aside?: ReactNode; cls?: string; children: ReactNode }) {
  return (
    <section className={`deck${cls ? ` ${cls}` : ''}${on ? ' on' : ''}`} aria-label={title}>
      <div className="dt"><span>[{n}] {title}</span>{aside ?? (count !== undefined && <span className="num">{count}</span>)}</div>
      {children}
    </section>
  )
}

function Row({ m, r, pane, phone }: { m: M; r: Mr; pane: Pane; phone?: boolean }) {
  const on = m.picked === r
  return (
    <button type="button" className={`mrow ${pane}${on ? ' on' : ''}`} aria-current={on || undefined} onClick={() => m.open(r)}>
      <span className="ref">{r.ref}</span>
      {phone ? <span className="ti"><b>{r.title}</b><span className="sub">{r.name} · {r.forge}</span></span> : <span className="ti">{r.title}</span>}
      <span className="st">
        {pane === 'merged' ? <span>merged {time(r.merged_at!)}</span> : (
          <>
            {r.draft && <span className="dim">draft</span>}
            <Ci m={r} />
            {r.threads_open > 0 && <span className="c-n">◆ {r.threads_open}</span>}
            <span className="dim">{ago(r.updated, m.now)}</span>
          </>
        )}
      </span>
    </button>
  )
}

function BranchRow({ b, a }: { b: MrBranch; a: A }) {
  return (
    <div className="mrow br">
      <span className="ref">⑂</span>
      <span className="ti"><b>{b.branch}</b><span className="sub">{b.name || b.project} · {b.ahead} ahead · no MR</span></span>
      <span className="st"><button type="button" className="btn sm" onClick={() => a.branch(b)}>Open MR</button></span>
    </div>
  )
}

// The detail's action row. Answer threads leads on your own MRs, Review on the others.
function Actions({ r, a, phone }: { r: Mr; a: A; phone?: boolean }) {
  const mine = r.role === 'AUTHOR'
  const merged = !!r.merged_at
  return (
    <div className="macts">
      <div className="btnrow">
        {!merged && mine && (
          <button type="button" className={`btn${r.threads_open ? ' p' : ''}`} disabled={!r.threads_open} onClick={() => a.threads(r)}>Answer threads</button>
        )}
        {!merged && !mine && <button type="button" className="btn p" onClick={() => a.review(r)}>Review</button>}
        <button type="button" className="btn" onClick={() => a.session(r)}>{phone ? 'Session' : 'Repo session'}</button>
        <a className="btn" href={r.url} target="_blank" rel="noreferrer">{phone ? 'Forge ↗' : 'Open on forge ↗'}</a>
      </div>
      {mine && !merged && !phone && <div className="note">Replies post after your OK in the threads session.</div>}
      {a.said && <div className={`note${a.said.f ? ' c-f' : ''}`} role="status">{a.said.text}</div>}
    </div>
  )
}

const listOf = (d: MrsData, p: Pane) => (p === 'review' ? d.review : p === 'mine' ? d.mine : d.merged)

function PaneBody({ m, p, phone }: { m: M; p: Pane; phone?: boolean }) {
  const list = listOf(m.d, p)
  const branches = p === 'mine' ? m.d.branches : []
  return (
    <div className="mlist">
      {list.map((r) => <Row key={keyOf(r)} m={m} r={r} pane={p} phone={phone} />)}
      {branches.map((b) => <BranchRow key={`${b.project}:${b.branch}`} b={b} a={m.a} />)}
      {!list.length && !branches.length && <div className="empty">{p === 'merged' ? 'Nothing merged today.' : 'Nothing here.'}</div>}
    </div>
  )
}

function Panes({ m }: { m: M }) {
  return (
    <div className="mpanes">
      {PANES.map(([p, title], i) => (
        <Deck key={p} n={i + 1} title={title} count={listOf(m.d, p).length} on={m.pane === p}><PaneBody m={m} p={p} /></Deck>
      ))}
    </div>
  )
}

// Phone: the pane tabs, then the one pane they pick.
const SHORT: Record<Pane, string> = { review: 'Review', mine: 'Mine', merged: 'Merged' }
function PhonePanes({ m }: { m: M }) {
  const i = PANES.findIndex(([p]) => p === m.paneTab)
  return (
    <>
      <div className="segs ptabs" role="tablist" aria-label="Panes">
        {PANES.map(([p], n) => {
          const c = listOf(m.d, p).length
          return <button key={p} type="button" role="tab" className={`seg${m.paneTab === p ? ' on' : ''}`} aria-selected={m.paneTab === p} onClick={() => m.setPaneTab(p)}>[{n + 1}] {SHORT[p]}{c && p !== 'merged' ? ` ${c}` : ''}</button>
        })}
      </div>
      <Deck n={i + 1} title={PANES[i][1]} on><PaneBody m={m} p={m.paneTab} phone /></Deck>
    </>
  )
}

function Facts({ r, now, phone }: { r: Mr; now: Date; phone?: boolean }) {
  const failed = r.jobs?.find((j) => j.state === 'failed')?.name
  const open = <span className={r.threads_open ? 'c-n' : ''}>{r.threads_open} open</span>
  if (phone) {
    const [g, c] = CI[r.pipeline]
    const f: Array<[string, ReactNode]> = [
      ['CI', <span className={c}>{g} {r.pipeline === 'failed' ? failed ?? 'failed' : r.pipeline}</span>],
      ['Threads', open],
    ]
    if (r.approvals) f.push(['Approvals', `${r.approvals.given}/${r.approvals.required}`])
    return <dl className="facts">{f.map(([k, v]) => <div key={k}><dt>{k}</dt><dd>{v}</dd></div>)}</dl>
  }
  const f: Array<[string, ReactNode]> = [
    ['Pipeline', r.pipeline === 'failed' ? <span className="c-f">failed{failed ? ` · ${failed}` : ''}</span> : <span className={CI[r.pipeline][1]}>{r.pipeline}</span>],
    ['Threads', <>{open} · {r.threads_resolved} resolved</>],
  ]
  if (r.approvals) f.push(['Approvals', `${r.approvals.given} / ${r.approvals.required}`])
  if (r.reviewers?.length) f.push([r.role === 'AUTHOR' ? 'Reviewer' : 'Author', r.role === 'AUTHOR' ? r.reviewers.join(', ') : r.author])
  else if (r.role !== 'AUTHOR') f.push(['Author', r.author])
  f.push(['Updated', `${ago(r.updated, now)} ago`])
  return (
    <dl className="facts">
      {f.map(([k, v]) => <div key={k}><dt>{k}</dt><dd>{v}</dd></div>)}
    </dl>
  )
}

export function Pipeline({ r }: { r: Mr }) {
  if (!r.jobs?.length && !r.failure) return <div className="empty">No pipeline details.</div>
  return (
    <div className="pipe">
      <div className="lbl"><span>Pipeline{r.pipeline_id ? ` #${r.pipeline_id}` : ''}</span></div>
      {!!r.jobs?.length && (
        <div className="jobs">
          {r.jobs.map((j) => <span key={j.name} className={`job ${j.state}`}>{j.name}{JOB[j.state] && ` ${JOB[j.state]}`}{j.duration && ` ${j.duration}`}</span>)}
        </div>
      )}
      {r.failure && <pre className="flog">{r.failure}</pre>}
    </div>
  )
}

export function Threads({ r, now, all }: { r: Mr; now: Date; all?: boolean }) {
  const open = (r.threads ?? []).filter((t) => all || !t.resolved)
  if (!r.threads) return <div className="empty">{r.threads_open} open · {r.threads_resolved} resolved. Thread texts open on the forge.</div>
  return (
    <div className="thr">
      <div className="lbl"><span>{all ? 'Threads' : 'Open threads'} · {all ? r.threads.length : r.threads_open}</span>{!all && <span>{r.threads_resolved} resolved</span>}</div>
      {open.map((t, i) => (
        <div key={i} className={`th${t.resolved ? ' done' : ''}`}>
          <div className="tm"><span className={t.resolved ? 'c-d' : 'c-n'}>{t.resolved ? 'resolved' : 'open'}</span>{t.path && <> · {t.path}{t.line ? `:${t.line}` : ''}</>} · {t.author} · {ago(t.at, now)}</div>
          <div className="tb">“{t.body}”</div>
        </div>
      ))}
      {!open.length && <div className="empty">No open thread.</div>}
    </div>
  )
}

// Review patterns from the ledgers (common + project), newest first, ◆ on the ones this
// MR hits. Only rendered when its tab is picked.
export function Patterns({ d, r }: { d: MrsData; r: Mr }) {
  const all = patternsFor(d, r)
  if (!all.length) return <div className="empty">No review pattern in the ledgers.</div>
  const rows = all.slice(0, PATTERN_ROWS)
  return (
    <div className="pats">
      <div className="lbl"><span>Review patterns · {r.name}</span><span>newest first</span></div>
      <table>
        <thead><tr><th>pattern</th><th>MRs</th><th>last</th></tr></thead>
        <tbody>
          {rows.map((p) => {
            const hit = !!p.hits?.includes(r.ref)
            return <tr key={p.pattern} className={hit ? 'hit' : undefined}><td>{hit && '◆ '}{p.pattern}</td><td>{p.count}</td><td className="c-w">{p.last ?? ''}</td></tr>
          })}
        </tbody>
      </table>
      <div className="note">◆ = this MR hits it · common + project ledger{all.length > rows.length ? ` · ${all.length - rows.length} more` : ''}</div>
    </div>
  )
}

function patTab(d: MrsData, r: Mr) {
  const n = hitsOf(patternsFor(d, r), r).length
  return <>Patterns{n ? ` ◆${n}` : ''}</>
}

function Detail({ m }: { m: M }) {
  const r = m.picked
  if (!r) return <Deck n={4} title="Detail"><div className="empty">No MR open on your forges.</div></Deck>
  const tab = (t: Tab, label: ReactNode) => (
    <button key={t} type="button" role="tab" className={`mtab${m.tab === t ? ' on' : ''}`} aria-selected={m.tab === t} onClick={() => m.setTab(t)}>{label}</button>
  )
  return (
    <Deck n={4} title={r.ref} on aside={<span className="dim">{r.name} · {r.forge}</span>}>
      <div className="mdet">
        <h2><span className="ref">{r.ref}</span> {r.title}</h2>
        <div className="meta">{r.branch} → {r.target} · opened {ago(r.created, m.now)} ago by {r.author} · updated {ago(r.updated, m.now)} ago</div>
        <div className="mtabs" role="tablist">
          {tab('overview', 'Overview')}
          {tab('threads', <>Threads {r.threads_open || ''}</>)}
          {tab('pipeline', <>Pipeline {CI[r.pipeline][0]}</>)}
          {tab('patterns', patTab(m.all, r))}
        </div>
        <div className="mbody">
          {m.tab === 'overview' && <><Facts r={r} now={m.now} /><Pipeline r={r} /><Threads r={r} now={m.now} /></>}
          {m.tab === 'threads' && <Threads r={r} now={m.now} all />}
          {m.tab === 'pipeline' && <Pipeline r={r} />}
          {m.tab === 'patterns' && <Patterns d={m.all} r={r} />}
        </div>
        <Actions r={r} a={m.a} />
      </div>
    </Deck>
  )
}

// Phone: the detail as its own page, back to the pane it came from.
function PhoneDetail({ m }: { m: M }) {
  const r = m.picked!
  const i = PANES.findIndex(([p]) => p === m.pane)
  const tab = (t: Tab, label: ReactNode) => (
    <button key={t} type="button" role="tab" className={`mtab${m.tab === t ? ' on' : ''}`} aria-selected={m.tab === t} onClick={() => m.setTab(t)}>{label}</button>
  )
  return (
    <div className="mdet ph">
      <button type="button" className="mback" onClick={m.back}>◀ [{i + 1}] {SHORT[m.pane!]}</button>
      <h2><span className="ref">{r.ref}</span> {r.title}</h2>
      <div className="meta">{r.name} · {r.branch} → {r.target} · {ago(r.updated, m.now)}</div>
      <Facts r={r} now={m.now} phone />
      <div className="mtabs" role="tablist">
        {tab('overview', 'Overview')}
        {tab('threads', <>Threads {r.threads_open || ''}</>)}
        {tab('patterns', patTab(m.all, r))}
      </div>
      <div className="mbody">
        {m.tab === 'threads' ? <Threads r={r} now={m.now} all /> : m.tab === 'patterns' ? <Patterns d={m.all} r={r} /> : <><Pipeline r={r} /><Threads r={r} now={m.now} /></>}
      </div>
      <Actions r={r} a={m.a} phone />
    </div>
  )
}

function Off({ d }: { d: MrsData }) {
  return <div className="desk mrs"><div className="mhead"><h1>MRs</h1></div><div className="empty">{d.off ? 'The forge module is off.' : 'No MR open on your forges.'}</div></div>
}

export function Mrs({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const m = useMrs(snap, refresh)
  const laptop = useMedia(LAPTOP)
  if (m.all.off) return <Off d={m.all} />
  const stale = m.all.stale && <div className="empty c-f">MRs did not refresh: {m.all.error}</div>
  const sheet = m.a.sheet && <Sheet c={m.a.sheet} onClose={() => m.a.setSheet(null)} />
  if (!laptop) {
    return (
      <div className="desk mrs m-ph">
        {m.shown ? <PhoneDetail m={m} /> : <><Head m={m} phone />{stale}<PhonePanes m={m} /></>}
        {sheet}
      </div>
    )
  }
  return (
    <div className="desk mrs m-lp">
      <Head m={m} />
      {stale}
      <div className="mgrid"><Panes m={m} /><Detail m={m} /></div>
      {sheet}
    </div>
  )
}
