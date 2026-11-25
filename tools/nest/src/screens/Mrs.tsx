import { useState, type ReactNode } from 'react'
import { countsOf, mrsOf, onlyProject, type Mr, type MrBranch, type MrJob, type MrsData } from '../data/mrs'
import type { Snapshot } from '../data/useData'
import { ago, time } from '../format'
import { useQuery } from '../router'
import { useMedia } from '../ui/useMedia'

// The MRs page (item_225, board round 2 option C): three deck panes (to review, mine,
// merged today) and a tabbed detail of the picked MR.

const LAPTOP = '(min-width: 1100px)'
export type Pane = 'review' | 'mine' | 'merged'
export type Tab = 'overview' | 'threads' | 'pipeline'
const PANES: Array<[Pane, string]> = [['review', 'To review'], ['mine', 'Mine'], ['merged', 'Merged today']]
const TABS: Tab[] = ['overview', 'threads', 'pipeline']
const keyOf = (m: Mr) => `${m.project}:${m.ref}`

export function useMrs(snap: Snapshot) {
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
  const open = (m: Mr) => { setPick(keyOf(m)); setTab('overview') }
  return { all, d, project, setProject, picked, pane, tab, setTab, open, now: snap.now }
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

function Head({ m }: { m: M }) {
  const c = countsOf(m.d)
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
export function Deck({ n, title, count, on, aside, children }: { n: number; title: string; count?: number; on?: boolean; aside?: ReactNode; children: ReactNode }) {
  return (
    <section className={`deck${on ? ' on' : ''}`} aria-label={title}>
      <div className="dt"><span>[{n}] {title}</span>{aside ?? (count !== undefined && <span className="num">{count}</span>)}</div>
      {children}
    </section>
  )
}

function Row({ m, r, pane }: { m: M; r: Mr; pane: Pane }) {
  const on = m.picked === r
  return (
    <button type="button" className={`mrow ${pane}${on ? ' on' : ''}`} aria-current={on || undefined} onClick={() => m.open(r)}>
      <span className="ref">{r.ref}</span>
      <span className="ti">{r.title}</span>
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

function BranchRow({ b }: { b: MrBranch }) {
  return (
    <div className="mrow br">
      <span className="ref">⑂</span>
      <span className="ti"><b>{b.branch}</b><span className="sub">{b.name || b.project} · {b.ahead} ahead · no MR</span></span>
      <span className="st"><button type="button" className="btn sm" disabled title="Opens with the row actions">Open MR</button></span>
    </div>
  )
}

function Panes({ m }: { m: M }) {
  const { d } = m
  const list = (p: Pane) => (p === 'review' ? d.review : p === 'mine' ? d.mine : d.merged)
  return (
    <div className="mpanes">
      {PANES.map(([p, title], i) => (
        <Deck key={p} n={i + 1} title={title} count={list(p).length} on={m.pane === p}>
          <div className="mlist">
            {list(p).map((r) => <Row key={keyOf(r)} m={m} r={r} pane={p} />)}
            {p === 'mine' && d.branches.map((b) => <BranchRow key={`${b.project}:${b.branch}`} b={b} />)}
            {!list(p).length && !(p === 'mine' && d.branches.length) && <div className="empty">{p === 'merged' ? 'Nothing merged today.' : 'Nothing here.'}</div>}
          </div>
        </Deck>
      ))}
    </div>
  )
}

function Facts({ r, now }: { r: Mr; now: Date }) {
  const f: Array<[string, ReactNode]> = [
    ['Pipeline', r.pipeline === 'failed' ? <span className="c-f">failed{r.jobs?.find((j) => j.state === 'failed') ? ` · ${r.jobs.find((j) => j.state === 'failed')!.name}` : ''}</span> : <span className={CI[r.pipeline][1]}>{r.pipeline}</span>],
    ['Threads', <><span className={r.threads_open ? 'c-n' : ''}>{r.threads_open} open</span> · {r.threads_resolved} resolved</>],
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
        </div>
        <div className="mbody">
          {m.tab === 'overview' && <><Facts r={r} now={m.now} /><Pipeline r={r} /><Threads r={r} now={m.now} /></>}
          {m.tab === 'threads' && <Threads r={r} now={m.now} all />}
          {m.tab === 'pipeline' && <Pipeline r={r} />}
        </div>
      </div>
    </Deck>
  )
}

function Off({ d }: { d: MrsData }) {
  return <div className="desk mrs"><div className="mhead"><h1>MRs</h1></div><div className="empty">{d.off ? 'The forge module is off.' : 'No MR open on your forges.'}</div></div>
}

export function Mrs({ snap }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const m = useMrs(snap)
  const laptop = useMedia(LAPTOP)
  if (m.all.off) return <Off d={m.all} />
  return (
    <div className={`desk mrs${laptop ? ' m-lp' : ' m-ph'}`}>
      <Head m={m} />
      {m.all.stale && <div className="empty c-f">MRs did not refresh: {m.all.error}</div>}
      <div className="mgrid"><Panes m={m} /><Detail m={m} /></div>
    </div>
  )
}
