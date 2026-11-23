import type { CSSProperties, ReactElement } from 'react'
import { homeModel, type HomeModel, type NeedKind, type NeedRow } from '../data/home'
import type { Snapshot } from '../data/useData'
import { time } from '../format'
import { Link } from '../router'

// Home rail header (item_259, board H3): Needs you as tiles, then Next and Running side by side.
// The Fold inner Home shows the compact form: the top ask and one summary line.
// Each part is its own framed panel (.hpn) on the dotted rail, like the Desk panels.

const ICON: Record<NeedKind, ReactElement> = {
  ask: <svg viewBox="0 0 16 16" className="c-n" aria-hidden="true"><path d="M8 1 15 8 8 15 1 8Z" fill="currentColor" /></svg>,
  alert: <svg viewBox="0 0 16 16" className="c-f" aria-hidden="true"><path d="M8 1.5 15 14H1Z" fill="currentColor" /></svg>,
  session: (
    <svg viewBox="0 0 16 16" className="c-n" fill="none" stroke="currentColor" strokeWidth="1.6" aria-hidden="true">
      <rect x="1.5" y="2.5" width="13" height="11" rx="1.5" /><path d="M4.5 6.5l2 1.5-2 1.5" />
    </svg>
  ),
}

const SHOWN = 4
// Minutes until the next event that fill the countdown bar.
const WINDOW_MIN = 60

function Need({ r, top }: { r: NeedRow; top: boolean }) {
  return (
    <div className={`hh-need${top ? ' top' : ''}`}>
      <span className={`hh-ico k-${r.kind}`}>{ICON[r.kind]}</span>
      <span className="t">{r.title}</span>
      <Link href={r.href} className={`hh-b${top ? ' p' : ''}`}>{r.action}</Link>
      <span className="m">{r.sub}</span>
    </div>
  )
}

const minsTo = (iso: string, now: Date) => Math.round((new Date(iso).getTime() - now.getTime()) / 60000)

function nextOf(m: HomeModel, now: Date) {
  const [ev, then] = m.agenda.events
  if (!ev) return null
  const to = minsTo(ev.start, now)
  const when = to <= 0 ? 'now' : to < 60 ? `in ${to} min` : m.agenda.label === 'Today' ? `in ${Math.floor(to / 60)} h ${String(to % 60).padStart(2, '0')}` : m.agenda.label.toLowerCase()
  const fill = to <= 0 ? 100 : Math.max(0, Math.round((1 - Math.min(to, WINDOW_MIN) / WINDOW_MIN) * 100))
  return { ev, then, when, fill }
}

export function HomeHud({ snap }: { snap: Snapshot }) {
  const m = homeModel(snap)
  const nx = nextOf(m, snap.now)
  const runs = m.running.filter((r) => r.kind === 'run').length
  const sessions = m.running.filter((r) => r.kind === 'session' && r.state === 'working').length
  const live = m.working === 1 && sessions === 1 ? 'live session' : m.working === sessions && sessions ? 'live sessions' : 'live'
  const pips = Math.max(m.slots.max, 1)
  return (
    <div className="hh">
      <section className="hpn hh-needs" aria-label="Needs you">
        <div className="hpn-h">
          Needs you
          <span className="r">
            {m.needs.length > 0 && <span className="c-n">{m.needs.length} open</span>}
            <Link href="/inbox">Inbox ›</Link>
          </span>
        </div>
        <div className="hpn-b">
          {m.needs.length ? m.needs.slice(0, SHOWN).map((r, i) => <Need key={r.key} r={r} top={i === 0} />) : <p className="hh-quiet">All quiet, Sir.</p>}
          {m.needs.length > SHOWN && <Link href="/inbox" className="hh-more">{m.needs.length - SHOWN} more in Inbox ›</Link>}
        </div>
      </section>
      <div className="hh-duo">
        {!snap.sections.today?.off && <div className="hpn hh-next">
          <span className="lb">Next <Link href="/agenda">Agenda ›</Link></span>
          {nx ? (
            <>
              <span className="big">{time(nx.ev.start)} <small>{nx.ev.title}</small></span>
              <span className="s">{nx.when}{nx.then ? ` · then ${time(nx.then.start)} ${nx.then.title}` : ''}</span>
              <span className="hh-bar"><i style={{ width: `${nx.fill}%` } as CSSProperties} /></span>
            </>
          ) : <span className="s">Nothing left today</span>}
        </div>}
        <div className="hpn hh-run">
          <span className="lb">Running <Link href="/ops">Ops ›</Link></span>
          <span className="big c-w">{m.working} <small>{live}</small></span>
          <span className="s">{runs} {runs === 1 ? 'run' : 'runs'} · slots {m.slots.busy}/{m.slots.max}</span>
          <span className="hh-slots" aria-hidden="true">{Array.from({ length: pips }, (_, i) => <i key={i} className={i < m.slots.busy ? 'on' : ''} />)}</span>
        </div>
      </div>
    </div>
  )
}

// Fold inner Home (board F1): the top ask, then "+N need you · next HH:MM title" and Inbox.
export function HomeHudCompact({ snap }: { snap: Snapshot }) {
  const m = homeModel(snap)
  const [first] = m.needs
  const nx = nextOf(m, snap.now)
  const bits = [
    m.needs.length > 1 ? `+${m.needs.length - 1} need you` : first ? '' : 'all quiet',
    nx ? `next ${time(nx.ev.start)} ${nx.ev.title}` : '',
  ].filter(Boolean)
  return (
    <section className="hpn hh compact" aria-label="Needs you">
      {first && <Need r={first} top />}
      <div className="hh-sum"><span>{bits.join(' · ')}</span><Link href="/inbox">Inbox ›</Link></div>
    </section>
  )
}
