import { LABEL, uiState } from '../data/state'
import type { Run } from '../data/types'
import type { Snapshot } from '../data/useData'
import { duration, time, title } from '../format'
import { Link } from '../router'
import { Glyph, StateWord, stateClass } from '../ui/Glyph'
import type { UiState } from '../data/state'

function runState(r: Run): UiState {
  if (!r.ended) return 'working'
  if (r.cause === 'done') return 'done'
  if (r.cause === 'usage_limit') return 'sleeping'
  if (r.cause === 'stopped') return 'queued'
  return 'failed'
}

export function MissionDetail({ snap, id }: { snap: Snapshot; id: string }) {
  const { missions, inbox, now } = snap
  const m = missions.find((x) => x.id === id)
  if (!m) {
    return (
      <div className="screen">
        <section className="pn"><div className="ph"><span>Mission not found</span></div><Link href="/missions" className="btn">Back to Missions</Link></section>
      </div>
    )
  }
  const s = uiState(m, inbox)
  const rows = inbox.filter((r) => r.mission === m.id)
  const last = m.runs[m.runs.length - 1]
  const props: Array<[string, string]> = [
    ['project', m.project],
    ['kind', m.kind],
    ['scope', `${m.scope} · ${m.data_class}`],
    ['autonomy', m.autonomy],
    ['gates', m.gates.join(', ') || 'none'],
    ['account', last?.account ?? 'not picked'],
    ['version', String(m.version)],
    ['batch', m.batch ?? 'none'],
    ['branch', m.branch ?? 'none'],
    ['worktree', m.worktree ?? 'none'],
    ['files', m.file_scope.join(', ') || 'any'],
    ['budget', `${m.budget.time_min} min · $${m.budget.budget_usd} · ${m.budget.resumes} resumes`],
    ['MR', m.mr ?? 'none'],
    ['created', time(m.created)],
    ['updated', time(m.updated)],
  ]
  if (last && !last.ended) props.splice(6, 0, ['elapsed', duration(last.started, now)])

  return (
    <div className="screen detail">
      <section className="pn card span2" aria-labelledby="mt">
        <Link href="/missions" className="sub">‹ Missions</Link>
        <div className="wrap" style={{ gap: 10 }}>
          <Glyph state={s} size={22} />
          <h1 id="mt" style={{ font: '600 17px var(--font-head)', flex: 1, minWidth: 0 }}>{title(m)}</h1>
          <StateWord state={s} />
        </div>
        <div className="dim">{m.goal}</div>
        <div className="sub">end state: {m.end_state} · contract state <span className={stateClass(s)}>{m.state}</span></div>
      </section>

      <div className="stack">
        <section className="pn" aria-labelledby="rl">
          <div className="ph"><h2 id="rl" style={{ font: 'inherit' }}>Run ledger</h2><span className="sub num">{m.runs.length} {m.runs.length === 1 ? 'run' : 'runs'}</span></div>
          {m.runs.length === 0 && <div className="empty">No run yet. It starts when a slot and an account are free.</div>}
          {m.runs.map((r) => {
            const rs = runState(r)
            return (
              <div key={r.run} className={`row${rs === 'done' ? ' past' : ''}`}>
                <Glyph state={rs} />
                <span className="t">
                  <b>v{r.version} · {r.cold ? 'cold' : 'warm'} · {r.account}</b>
                  <span className="num">{time(r.started)} to {r.ended ? time(r.ended) : 'now'} · {duration(r.started, r.ended ?? now)}{r.exit !== null ? ` · exit ${r.exit}` : ''}</span>
                </span>
                <span className={`sub ${stateClass(rs)}`}>{r.ended ? (r.cause ?? 'ended').replace('_', ' ') : LABEL.working}</span>
              </div>
            )
          })}
        </section>
        <section className="pn" aria-labelledby="mi">
          <div className="ph"><h2 id="mi" style={{ font: 'inherit' }}>Inbox for this mission</h2><Link href="/inbox">Inbox ›</Link></div>
          {rows.length === 0 && <div className="empty">Nothing asked.</div>}
          {rows.map((r) => (
            <div key={r.id} className={`row${r.state !== 'open' ? ' past' : ''}`}>
              {r.state === 'open' && r.type !== 'fyi' ? <Glyph state="needs" /> : <Glyph state="done" />}
              <span className="t"><b>{r.ask}</b><span>{r.source} v{r.version} · {r.type} · {r.state}{r.answered ? ` ${time(r.answered)}` : ''}</span></span>
            </div>
          ))}
        </section>
      </div>

      <section className="pn" aria-labelledby="pr">
        <div className="ph"><h2 id="pr" style={{ font: 'inherit' }}>Mission</h2><span className="sub">properties</span></div>
        <dl className="props">
          {props.map(([k, v]) => (
            <div key={k} style={{ display: 'contents' }}><dt>{k}</dt><dd className="num">{v}</dd></div>
          ))}
        </dl>
      </section>
    </div>
  )
}
