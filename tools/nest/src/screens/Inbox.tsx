import { useState } from 'react'
import type { InboxItem } from '../data/types'
import type { Snapshot } from '../data/useData'
import { time, title } from '../format'
import { Link } from '../router'
import { Glyph } from '../ui/Glyph'

const GROUPS: Array<{ key: string; label: string; cls: string; test: (r: InboxItem) => boolean }> = [
  { key: 'ask', label: 'To answer', cls: 'c-n', test: (r) => r.state === 'open' && r.type !== 'fyi' },
  { key: 'fyi', label: 'FYI', cls: 'c-q', test: (r) => r.state === 'open' && r.type === 'fyi' },
  { key: 'past', label: 'Answered', cls: 'c-d', test: (r) => r.state !== 'open' },
]

export function Inbox({ snap }: { snap: Snapshot }) {
  const { inbox, missions } = snap
  const ordered = GROUPS.flatMap((g) => inbox.filter(g.test).sort((a, b) => b.created.localeCompare(a.created)))
  const [sel, setSel] = useState(ordered[0]?.id)
  const item = inbox.find((r) => r.id === sel)
  const mission = item && missions.find((m) => m.id === item.mission)

  return (
    <>
      <div className="backdrop" aria-hidden="true" />
      <div className="screen inbox over">
        <section className="pn" aria-labelledby="il">
          <div className="ph"><h1 id="il" style={{ font: 'inherit' }}>Inbox</h1><span className="sub">read only</span></div>
          {GROUPS.map((g) => {
            const items = inbox.filter(g.test).sort((a, b) => b.created.localeCompare(a.created))
            if (!items.length) return null
            return (
              <div key={g.key} className="stack8">
                <div className={`grp ${g.cls}`}><span>{g.label}</span><span className="num">{items.length}</span></div>
                {items.map((r) => (
                  <button key={r.id} type="button" className={`row${r.id === sel ? ' on' : ''}${g.key === 'past' ? ' past' : ''}`}
                    aria-pressed={r.id === sel} onClick={() => setSel(r.id)}>
                    {g.key === 'ask' ? <Glyph state="needs" /> : g.key === 'past' ? <Glyph state="done" /> : null}
                    <span className="t"><b>{r.ask}</b><span>{r.project} · {r.source} · {r.type}</span></span>
                    <span className="sub num">{time(r.created)}</span>
                  </button>
                ))}
              </div>
            )
          })}
        </section>

        {item ? (
          <section className="pn card" aria-labelledby="id" aria-live="polite">
            <div className="ph">
              <h2 id="id" style={{ font: '600 14px var(--font-head)', color: 'var(--fg)' }}>{item.ask}</h2>
              <span className="sub num">{time(item.created)}</span>
            </div>
            <div className="wrap">
              <span className="tag">{item.type}</span>
              <span className="tag">{item.source} v{item.version}</span>
              {item.blocking && <span className="tag c-f" style={{ borderColor: 'currentColor' }}>blocking</span>}
              <span className={`tag ${item.state === 'open' ? 'c-n' : 'c-d'}`} style={{ borderColor: 'currentColor' }}>{item.state}</span>
            </div>
            {mission && (
              <Link href={`/missions/${mission.id}`} className="row">
                <span className="t"><b>{title(mission)}</b><span>{mission.project} · {mission.kind} · {mission.goal}</span></span>
                <span className="chev" aria-hidden="true">›</span>
              </Link>
            )}
            <dl className="props">
              <dt>file</dt><dd>{item.file}</dd>
              <dt>artifacts</dt><dd>{item.artifacts.join(', ') || 'none'}</dd>
              <dt>replies</dt><dd>{item.verbs.join(' · ')}</dd>
              {item.answered && <><dt>answered</dt><dd className="num">{time(item.answered)}</dd></>}
            </dl>
            <div className="sub">Replies arrive with the live API. This slice is read only.</div>
          </section>
        ) : (
          <section className="pn"><div className="empty">Inbox is empty.</div></section>
        )}
      </div>
    </>
  )
}
