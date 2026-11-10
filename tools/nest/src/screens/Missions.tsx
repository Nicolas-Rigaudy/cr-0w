import { countByState, LABEL, UI_STATES, uiState, type UiState } from '../data/state'
import type { Snapshot } from '../data/useData'
import { sameDay } from '../format'
import { Glyph, stateClass } from '../ui/Glyph'
import { MissionRow } from '../ui/MissionRow'
import { RunsChart } from '../ui/RunsChart'
import { Kpi } from './Home'
import { byUrgency } from './order'

const VAR: Record<UiState, string> = { working: '--c-w', needs: '--c-n', failed: '--c-f', sleeping: '--c-s', queued: '--c-q', done: '--c-d' }

function Donut({ counts, total }: { counts: Record<UiState, number>; total: number }) {
  const r = 30, C = 2 * Math.PI * r
  let at = 0
  return (
    <svg width="84" height="84" viewBox="0 0 84 84" aria-hidden="true" style={{ flex: 'none' }}>
      <circle cx="42" cy="42" r={r} fill="none" stroke="var(--row)" strokeWidth="10" />
      {UI_STATES.filter((s) => counts[s]).map((s) => {
        const len = (counts[s] / total) * C
        const el = (
          <circle key={s} cx="42" cy="42" r={r} fill="none" stroke={`var(${VAR[s]})`} strokeWidth="10"
            strokeDasharray={`${Math.max(0, len - 2)} ${C}`} strokeDashoffset={-at} transform="rotate(-90 42 42)" />
        )
        at += len
        return el
      })}
      <text x="42" y="47" textAnchor="middle" style={{ fill: 'var(--fg)', font: '600 15px var(--font-display)' }}>{total}</text>
    </svg>
  )
}

export function Missions({ snap }: { snap: Snapshot }) {
  const { missions, inbox, now } = snap
  const c = countByState(missions, inbox)
  const sorted = byUrgency(missions, inbox)
  const doneToday = missions.filter((m) => uiState(m, inbox) === 'done' && sameDay(new Date(m.updated), now)).length
  const groups = UI_STATES.map((s) => ({ s, items: sorted.filter((m) => uiState(m, inbox) === s) })).filter((g) => g.items.length)

  return (
    <>
      <div className="backdrop" aria-hidden="true" />
      <div className="screen missions over">
        <section className="pn span2" aria-labelledby="mh">
          <div className="ph"><h1 id="mh" style={{ font: 'inherit' }}>Missions today</h1><span className="sub num">{missions.length} missions</span></div>
          <div style={{ display: 'flex', gap: 14, alignItems: 'center' }}>
            <span className="only-cover"><Donut counts={c} total={missions.length} /></span>
            <div className="kpis" style={{ flex: 1 }}>
              <Kpi n={c.needs} label="need you" cls="c-n" />
              <Kpi n={c.working} label="working" cls="c-w" />
              <Kpi n={c.failed} label="failed" cls="c-f" />
              <Kpi n={c.sleeping} label="sleeping" cls="c-s" />
              <Kpi n={c.queued} label="queued" cls="c-q" />
              <Kpi n={doneToday} label="done today" cls="c-d" />
            </div>
          </div>
        </section>

        <section className="pn" aria-labelledby="pl">
          <div className="ph"><h2 id="pl" style={{ font: 'inherit' }}>Pipeline</h2><span className="sub">by state</span></div>
          {groups.map(({ s, items }) => (
            <div key={s} className="stack8">
              <div className={`grp ${stateClass(s)}`}><span style={{ display: 'flex', gap: 6, alignItems: 'center' }}><Glyph state={s} size={12} />{LABEL[s]}</span><span className="num">{items.length}</span></div>
              {items.map((m) => <MissionRow key={m.id} m={m} inbox={inbox} now={now} />)}
            </div>
          ))}
        </section>

        <div className="stack">
          <section className="pn" aria-labelledby="rc">
            <div className="ph"><h2 id="rc" style={{ font: 'inherit' }}>Runs per day · 7 days</h2><span className="sub"><span className="c-w">■</span> ok <span className="c-f">■</span> error</span></div>
            <RunsChart missions={missions} now={now} />
          </section>
          <section className="pn" aria-labelledby="sl">
            <div className="ph"><h2 id="sl" style={{ font: 'inherit' }}>Waiting to run</h2><span className="sub num">{c.sleeping + c.queued}</span></div>
            {sorted.filter((m) => ['sleeping', 'queued'].includes(uiState(m, inbox))).map((m) => <MissionRow key={m.id} m={m} inbox={inbox} now={now} />)}
            <div className="sub">sleeping = stopped on a usage limit, resumes when the account resets</div>
          </section>
        </div>
      </div>
    </>
  )
}
