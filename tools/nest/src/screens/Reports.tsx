import { Fragment, useState } from 'react'
import { ApiError } from '../data/api'
import type { Actions } from '../data/live'
import {
  blocksOf, KINDS, RETRO_UNIT, reportsOf, RITUAL, running, spansOf, stamp, tilesOf,
  type Kind, type Report, type Span, type Tile,
} from '../data/reports'
import { useSource, type Snapshot } from '../data/useData'
import { ago } from '../format'
import { useQuery } from '../router'
import { Sheet, type Confirm } from '../ui/Sheet'
import { useMedia } from '../ui/useMedia'

// Reports page (item_258, board B): one tile per ritual, the settled outputs newest
// first, a reader with the rendered Markdown. Only the retro routine can be started
// from here; the other rituals wait for item_234.
const LAPTOP = '(min-width: 1100px)'
const PHONE_ROWS = 5

type Said = { text: string; f?: boolean }

export function retroConfirm(act: Actions, refresh: () => Promise<void>, setSaid: (s: Said) => void, onStart: () => void): Confirm {
  return {
    title: 'Generate a retro', lead: 'Starts the retro routine now as a background run.',
    target: [['Routine', RETRO_UNIT], ['Output', 'Retros folder in the vault']],
    note: 'Changes to tracked files wait for your OK.', ok: 'Start retro',
    run: async () => {
      try {
        const r = await act.routine(RETRO_UNIT, 'run')
        setSaid({ text: r.message || `${r.action} ${r.target}`, f: !r.ok })
        if (r.ok) onStart()
      } catch (e) {
        if (e instanceof ApiError && e.status === 409) void refresh()
        throw e
      }
      void refresh()
    },
  }
}

export function useReports(snap: Snapshot, refresh: () => Promise<void>) {
  const act = useSource().actions
  const q = useQuery()
  const d = reportsOf(snap.sections.reports, snap.sections.jobs)
  const tiles = tilesOf(d.items, snap.now)
  const fq = q.get('kind') as Kind | null
  const [kind, setKind] = useState<Kind | null>(fq && KINDS.includes(fq) ? fq : null)
  const rows = kind ? d.items.filter((r) => r.kind === kind) : d.items
  // ?id=<report id> opens a report (shots, links from other pages).
  const want = q.get('id')
  const [pick, setPick] = useState<string | null>(want)
  const picked = rows.find((r) => r.id === pick) ?? rows[0] ?? null
  const [shown, setShown] = useState(!!want && !!picked)
  const [started, setStarted] = useState<string | null>(null)
  const [sheet, setSheet] = useState<Confirm | null>(null)
  const [said, setSaid] = useState<Said | null>(null)
  return {
    d, tiles, rows, kind, picked, now: snap.now, sheet, setSheet, said, setSaid,
    filter: (k: Kind) => { setKind(kind === k ? null : k); setPick(null) },
    open: (r: Report) => { setPick(r.id); setShown(true) },
    shown: shown && !!picked, back: () => setShown(false),
    started, busy: running(started, d.items, snap.now),
    generate: () => setSheet(retroConfirm(act, refresh, setSaid, () => setStarted(new Date().toISOString()))),
  }
}
export type R = ReturnType<typeof useReports>

function Glyph({ k }: { k: Kind }) {
  return <span className={`rg ${RITUAL[k].c}`} aria-hidden="true">{RITUAL[k].g}</span>
}

function Foot({ r, t, phone }: { r: R; t: Tile; phone?: boolean }) {
  if (t.kind === 'afk') return <span className="rfoot">from afk runs</span>
  if (t.kind === 'retro' && r.busy) return <span className="chip run"><i className="dot" />{phone ? '' : 'running '}{ago(r.started!, r.now)}</span>
  if (t.kind === 'retro' && r.d.retro) {
    return (
      <button type="button" className="btn p sm" onClick={(e) => { e.stopPropagation(); r.generate() }}>↻ {phone ? 'Gen' : 'Generate'}</button>
    )
  }
  if (t.last) return <span className="rfoot">no routine yet</span>
  return <span className="rfoot">no saved output</span>
}

function TileBox({ r, t, phone }: { r: R; t: Tile; phone?: boolean }) {
  const on = r.kind === t.kind
  const run = t.kind === 'retro' && r.busy
  return (
    <div className={`rtile${on ? ' on' : ''}${run ? ' run' : ''}`} onClick={() => r.filter(t.kind)}>
      <button type="button" className="rtn" aria-pressed={on} onClick={(e) => { e.stopPropagation(); r.filter(t.kind) }}>
        <span className={`rname ${RITUAL[t.kind].c}`}>{RITUAL[t.kind].g} {RITUAL[t.kind].name}</span>
        {!phone && <span className="rwk">{t.week} this week</span>}
      </button>
      <span className="rlast">{phone ? '' : 'last: '}{t.last ? stamp(t.last.at) : '-'}</span>
      <Foot r={r} t={t} phone={phone} />
    </div>
  )
}

function Tiles({ r, phone }: { r: R; phone?: boolean }) {
  return <div className="rtiles" role="group" aria-label="Rituals">{r.tiles.map((t) => <TileBox key={t.kind} r={r} t={t} phone={phone} />)}</div>
}

function Row({ r, x }: { r: R; x: Report }) {
  const on = r.picked?.id === x.id
  return (
    <button type="button" className={`rrow${on ? ' on' : ''}`} aria-current={on || undefined} onClick={() => r.open(x)}>
      <Glyph k={x.kind} />
      <span className="rt">{x.title}</span>
      <span className="rs">{x.summary ? `${x.summary} · ` : ''}{stamp(x.at)}</span>
    </button>
  )
}

function List({ r, phone }: { r: R; phone?: boolean }) {
  const rows = phone ? r.rows.slice(0, PHONE_ROWS) : r.rows
  const empty = r.kind ? `No saved ${RITUAL[r.kind].name.toLowerCase()} output yet.` : 'Nothing settled yet.'
  return (
    <section className="rpane rlist" aria-label="Reports">
      {rows.length ? rows.map((x) => <Row key={x.id} r={r} x={x} />) : <div className="empty">{empty}</div>}
      {r.d.held > 0 && (!r.kind || r.kind === 'afk') && (
        <div className="rheld">{r.d.held} AFK report{r.d.held > 1 ? 's wait' : ' waits'} for a decision in Inbox</div>
      )}
    </section>
  )
}

function Inline({ s }: { s: Span[] }) {
  return <>{s.map((x, i) => x.t === 'b' ? <b key={i}>{x.s}</b> : x.t === 'code' ? <code key={i}>{x.s}</code> : x.t === 'a' ? <a key={i} href={x.href} target="_blank" rel="noreferrer">{x.s}</a> : <Fragment key={i}>{x.s}</Fragment>)}</>
}

export function Markdown({ text }: { text: string }) {
  return (
    <div className="md">
      {blocksOf(text).map((b, i) => {
        if (b.t === 'h') return b.level === 3 ? <h4 key={i}><Inline s={spansOf(b.text)} /></h4> : <h3 key={i}><Inline s={spansOf(b.text)} /></h3>
        if (b.t === 'table') return (
          <div key={i} className="tw"><table>
            <thead><tr>{b.head.map((c, j) => <th key={j}><Inline s={spansOf(c)} /></th>)}</tr></thead>
            <tbody>{b.rows.map((r, j) => <tr key={j}>{r.map((c, k) => <td key={k}><Inline s={spansOf(c)} /></td>)}</tr>)}</tbody>
          </table></div>
        )
        if (b.t === 'pre') return <pre key={i}>{b.text}</pre>
        if (b.t === 'ul') return <ul key={i}>{b.items.map((t, j) => <li key={j}><Inline s={spansOf(t)} /></li>)}</ul>
        return <p key={i}><Inline s={spansOf(b.text)} /></p>
      })}
    </div>
  )
}

// Drops the note's own title line: the reader header carries it.
function bodyOf(x: Report) {
  const l = x.body.replace(/^\s+/, '')
  const m = /^#\s+(.*)\n?/.exec(l)
  return m && m[1].trim() === x.title ? l.slice(m[0].length) : l
}

function ReaderHead({ r, x }: { r: R; x: Report }) {
  const copy = () => {
    void navigator.clipboard?.writeText(x.body).then(() => r.setSaid({ text: 'Copied' }), () => r.setSaid({ text: 'Copy failed', f: true }))
  }
  return (
    <div className="rph">
      <span className={`rtitle ${RITUAL[x.kind].c}`}>{RITUAL[x.kind].g} {x.title}</span>
      <span className="racts">
        {x.link && <a className="btn sm" href={x.link}>Vault ↗</a>}
        <button type="button" className="btn sm" onClick={copy}>Copy</button>
      </span>
    </div>
  )
}

function Reader({ r }: { r: R }) {
  const x = r.picked
  if (!x) return <section className="rpane rread" aria-label="Reader"><div className="empty">Pick a report.</div></section>
  return (
    <section className="rpane rread" aria-label="Reader">
      <ReaderHead r={r} x={x} />
      <div className="rbody"><Markdown text={bodyOf(x)} /></div>
    </section>
  )
}

function PhoneReader({ r }: { r: R }) {
  const x = r.picked!
  return (
    <div className="rread ph">
      <button type="button" className="mback" onClick={r.back}>◀ Reports</button>
      <ReaderHead r={r} x={x} />
      <div className="rmeta">{x.summary ? `${x.summary} · ` : ''}{stamp(x.at)}</div>
      <div className="rbody"><Markdown text={bodyOf(x)} /></div>
    </div>
  )
}

function SaidLine({ r }: { r: R }) {
  return r.said ? <div className={`rsaid${r.said.f ? ' c-f' : ''}`} role="status">{r.said.text}</div> : null
}

export function Reports({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const r = useReports(snap, refresh)
  const laptop = useMedia(LAPTOP)
  const head = (
    <div className="rhead">
      <h1>Reports</h1>
      {laptop && <span className="sum">settled only · open decisions stay in Inbox</span>}
    </div>
  )
  if (r.d.off) return <div className="desk reports">{head}<div className="empty">The reports section is off.</div></div>
  const stale = r.d.stale && <div className="empty c-f">Reports did not refresh: {r.d.error}</div>
  const sheet = r.sheet && <Sheet c={r.sheet} onClose={() => r.setSheet(null)} />
  if (!laptop) {
    return (
      <div className="desk reports r-ph">
        {r.shown ? <PhoneReader r={r} /> : <>{head}{stale}<Tiles r={r} phone /><SaidLine r={r} /><List r={r} phone /></>}
        {sheet}
      </div>
    )
  }
  return (
    <div className="desk reports r-lp">
      {head}
      {stale}
      <Tiles r={r} />
      <SaidLine r={r} />
      <div className="rgrid"><List r={r} /><Reader r={r} /></div>
      {sheet}
    </div>
  )
}
