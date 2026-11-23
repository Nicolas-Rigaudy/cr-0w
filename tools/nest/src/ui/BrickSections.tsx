import { useState } from 'react'
import { listRows, listSections, type ListSection, type ModuleAction } from '../data/modules'
import type { Snapshot } from '../data/useData'
import { useSource } from '../data/useData'

// Generic list sections of user bricks (item_237): drawn on the page their manifest names,
// after that page's own panels. Actions post to core; a confirm action asks once more.
export function BrickSections({ snap, page }: { snap: Snapshot; page: string }) {
  const secs = listSections(snap.modules ?? null, page)
  if (!secs.length) return null
  return (
    <div className="screen bricks">
      {secs.map((s) => <BrickList key={`${s.brick}/${s.id}`} sec={s} snap={snap} />)}
    </div>
  )
}

function BrickList({ sec, snap }: { sec: ListSection; snap: Snapshot }) {
  const body = snap.bricks?.[sec.source]
  const rows = listRows(body)
  const head = (sec.actions ?? []).filter((a) => !a.row)
  const row = (sec.actions ?? []).filter((a) => a.row)
  const [note, setNote] = useState('')
  const hid = `brick-${sec.brick}-${sec.id}`
  return (
    <section className="pn brick" aria-labelledby={hid}>
      <div className="ph">
        <h2 id={hid} style={{ font: 'inherit' }}>{sec.title ?? sec.id}</h2>
        <span className="brick-acts">{head.map((a) => <ActionButton key={a.id} brick={sec.brick} a={a} onDone={setNote} />)}</span>
      </div>
      {body?.stale && body.error && <div className="sub c-f">{body.error}</div>}
      {rows.length ? rows.map((r) => (
        <div key={r.id} className="brick-row">
          <span className="brick-t">
            {r.url ? <a href={r.url} target="_blank" rel="noreferrer">{r.title}</a> : r.title}
            {r.detail && <span className="sub"> {r.detail}</span>}
          </span>
          {r.meta && <span className="sub num">{r.meta}</span>}
          {row.map((a) => <ActionButton key={a.id} brick={sec.brick} a={a} item={r.id} label={r.title} onDone={setNote} />)}
        </div>
      )) : <div className="empty">{body && !body.updated_at && !body.stale ? 'Loading…' : 'Nothing here.'}</div>}
      {note && <div className="sub" role="status">{note}</div>}
    </section>
  )
}

function ActionButton({ brick, a, item, label, onDone }: { brick: string; a: ModuleAction; item?: string; label?: string; onDone: (msg: string) => void }) {
  const src = useSource()
  const [armed, setArmed] = useState(false)
  const [busy, setBusy] = useState(false)
  const run = async () => {
    if (a.confirm && !armed) return setArmed(true)
    setArmed(false)
    setBusy(true)
    try {
      const r = await src.actions.module(brick, a.id, item)
      onDone(r.message || `${a.label} done`)
    } catch (e) {
      onDone(e instanceof Error ? e.message : String(e))
    } finally {
      setBusy(false)
    }
  }
  const name = label ? `${a.label}: ${label}` : a.label
  return (
    <button type="button" className={`btn sm${armed ? ' p' : ''}`} disabled={busy} aria-label={armed ? `Confirm ${name}` : name}
      onClick={() => void run()} onBlur={() => setArmed(false)}>
      {armed ? `${a.label}?` : a.label}
    </button>
  )
}
