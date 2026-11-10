import { countByState, openInbox, uiState } from '../data/state'
import type { Snapshot } from '../data/useData'
import { date, greeting, time, title } from '../format'
import { Link } from '../router'
import { Glyph } from '../ui/Glyph'
import { MissionRow } from '../ui/MissionRow'
import { byUrgency } from './order'

export function Home({ snap }: { snap: Snapshot }) {
  const { missions, inbox, now } = snap
  const c = countByState(missions, inbox)
  const asks = openInbox(inbox).filter((r) => r.type !== 'fyi')
  const fyi = openInbox(inbox).filter((r) => r.type === 'fyi')
  const first = asks[0]
  const firstMission = first && missions.find((m) => m.id === first.mission)
  const live = byUrgency(missions, inbox).filter((m) => uiState(m, inbox) !== 'done')
  const summary = c.needs
    ? `${c.needs} need you, ${c.working} working, ${c.failed} failed.`
    : `Nothing needs you. ${c.working} working.`

  return (
    <div className="screen home">
      <div className="stack ga">
        <section className="pn card" aria-labelledby="greet">
          <div className="hero tall" role="img" aria-label="CR-0W droid in the nest, idle" />
          <h1 id="greet" className="greet">{greeting(now)}</h1>
          <div className="sub num">{date(now)} · {time(now)}</div>
          <p className="dim" style={{ margin: 0 }}>{summary}</p>
          <div className="kpis only-cover" style={{ gridTemplateColumns: 'repeat(3, 1fr)' }}>
            <Kpi n={c.needs} label="need you" cls="c-n" />
            <Kpi n={c.working} label="working" cls="c-w" />
            <Kpi n={c.failed} label="failed" cls="c-f" />
          </div>
        </section>

        <section className="pn" aria-labelledby="needs">
          <div className="ph"><h2 id="needs" style={{ font: 'inherit' }}>Needs you first</h2><span className="sub num">{asks.length}</span></div>
          {first ? (
            <>
              <div className="well">
                <span className="lb c-m">{first.source} · {first.type}{first.blocking ? ' · blocking' : ''}</span>
                <span>{first.ask}</span>
                {firstMission && <span className="sub">{title(firstMission)} · {firstMission.project} · {time(first.created)}</span>}
              </div>
              <div className="wrap">
                <Link href="/inbox" className="btn p">Open in Inbox</Link>
                {firstMission && <Link href={`/missions/${firstMission.id}`} className="btn">Mission</Link>}
              </div>
              {asks.slice(1).map((r) => (
                <Link key={r.id} href="/inbox" className="row">
                  <Glyph state="needs" />
                  <span className="t"><b>{r.ask}</b><span>{r.project} · {r.source} · {time(r.created)}</span></span>
                  <span className="chev" aria-hidden="true">›</span>
                </Link>
              ))}
            </>
          ) : (
            <div className="empty">All quiet, Sir.</div>
          )}
        </section>
      </div>

      <section className="pn gb" aria-labelledby="ms">
        <div className="ph"><h2 id="ms" style={{ font: 'inherit' }}>Missions</h2><Link href="/missions">Open ›</Link></div>
        <div className="kpis only-inner">
          <Kpi n={c.done} label="done" cls="c-d" />
          <Kpi n={c.working} label="working" cls="c-w" />
          <Kpi n={c.needs} label="need you" cls="c-n" />
          <Kpi n={c.failed} label="failed" cls="c-f" />
          <Kpi n={c.sleeping + c.queued} label="queued" cls="c-q" />
        </div>
        <div className="gt"><span>Next</span><span className="num">{live.length}</span></div>
        {live.map((m) => <MissionRow key={m.id} m={m} inbox={inbox} now={now} big />)}
      </section>

      <div className="stack gc">
        <section className="pn" aria-labelledby="ib">
          <div className="ph"><h2 id="ib" style={{ font: 'inherit' }}>Inbox</h2><Link href="/inbox">Open ›</Link></div>
          <div className="wrap"><span className="kn c-n num">{asks.length}</span><span className="kl">to answer</span><span className="kn c-q num" style={{ marginLeft: 12 }}>{fyi.length}</span><span className="kl">fyi</span></div>
          {fyi.map((r) => (
            <Link key={r.id} href="/inbox" className="row">
              <span className="t"><b>{r.ask}</b><span>{r.project} · {r.source} · {time(r.created)}</span></span>
            </Link>
          ))}
        </section>
        {(['MRs', 'Reports', 'Agenda'] as const).map((name) => (
          <section key={name} className="pn">
            <div className="ph"><span>{name}</span><Link href={`/${name.toLowerCase()}`}>Open ›</Link></div>
            <div className="empty">Arrives with the live API.</div>
          </section>
        ))}
      </div>
    </div>
  )
}

export function Kpi({ n, label, cls }: { n: number; label: string; cls: string }) {
  return (
    <div className="kpi">
      <span className={`kn num ${cls}`}>{n}</span>
      <span className="kl">{label}</span>
    </div>
  )
}
