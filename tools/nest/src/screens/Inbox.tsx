import { useState } from 'react'
import { ApiError } from '../data/api'
import type { InboxItem, Mission } from '../data/types'
import { useSource, type Snapshot } from '../data/useData'
import { time, title } from '../format'
import { Link } from '../router'
import { REJECT, useAnswer, verbLabel } from '../ui/answer'
import { Sheet, type Confirm } from '../ui/Sheet'
import { useMedia } from '../ui/useMedia'

// Inbox (item_138 boards fold-inbox-inner, fold-inbox-cover) on the mission inbox of
// item_168. The reading pane answers inline: the well holds the note, the bottom row
// sends one verb. Hide sits where the board has Discard.
const INNER = '(min-width: 600px)'
const LATER = 'Arrives with the live API'
const short = (id: string) => id.replace(/^\d{8}-/, '')
const day = new Intl.DateTimeFormat('en-GB', { weekday: 'short', day: '2-digit' })
const stamp = (iso: string) => `${day.format(new Date(iso))} · ${time(iso)}`
const newest = (a: InboxItem, b: InboxItem) => b.created.localeCompare(a.created)
const initials = (p: string) => p.split(/[^a-z0-9]+/i).filter(Boolean).slice(0, 2).map((w) => w[0]).join('').toUpperCase() || '?'

type Seg = 'all' | 'answer' | 'drafts' | 'fyi'
const ask = (r: InboxItem) => r.state === 'open' && r.type !== 'fyi'
const fyi = (r: InboxItem) => r.state === 'open' && r.type === 'fyi'

// ---- icons (board paths) ----
const Diamond = ({ s = 13 }: { s?: number }) => (
  <svg width={s} height={s} viewBox="0 0 16 16" style={{ flex: 'none' }} aria-hidden="true"><path d="M8 1 15 8 8 15 1 8Z" fill="var(--c-n)" /></svg>
)
const Pen = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="var(--accent)" strokeWidth="2" strokeLinecap="round" aria-hidden="true"><path d="M4 20h4L19 9l-4-4L4 16z" /></svg>
)
const Env = () => (
  <svg viewBox="0 0 24 24" fill="none" stroke="#8087a2" strokeWidth="2" strokeLinecap="round" aria-hidden="true"><rect x="3" y="5" width="18" height="14" rx="2" /><path d="m3 7 9 6 9-6" /></svg>
)
const Unfold = () => (
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
    <rect x="3" y="4" width="8" height="16" rx="2" /><path d="M15 4h4a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-4M14 12h4m-2-2 2 2-2 2" />
  </svg>
)

function useInbox(snap: Snapshot, refresh: () => Promise<void>) {
  const { inbox, missions } = snap
  const act = useSource().actions
  const [seg, setSeg] = useState<Seg>('all')
  const [sheet, setSheet] = useState<Confirm | null>(null)
  const asks = inbox.filter(ask).sort(newest)
  const fyis = inbox.filter(fyi).sort(newest)
  const [sel, setSel] = useState<string | undefined>(asks[0]?.id ?? fyis[0]?.id)
  const picked = inbox.find((r) => r.id === sel) ?? asks[0] ?? fyis[0]
  const a = useAnswer(picked, inbox, missions, refresh)
  const missionOf = (r: InboxItem) => missions.find((m) => m.id === r.mission)

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
  return { seg, setSeg, sel: picked?.id, setSel, picked, asks, fyis, a, missionOf, hide, sheet, setSheet }
}
type I = ReturnType<typeof useInbox>

// ---- shared pieces ----
function From({ r, m }: { r: InboxItem; m?: Mission }) {
  return <>from {m ? title(m) : short(r.mission)} · {r.project} · {r.source} v{r.version}</>
}

function Stale({ i }: { i: I }) {
  if (!i.a.newer) return i.a.error ? <div className="err" role="alert">{i.a.error}</div> : null
  return (
    <div className="stale" role="alert">
      <b>THE MISSION MOVED ON · NOW V{i.a.mission?.version}</b>
      <div className="sm">You answered v{i.picked?.version}. This is the newer question. Your text is kept.</div>
    </div>
  )
}

function Note({ i, big }: { i: I; big?: boolean }) {
  const r = i.picked!
  return (
    <div className={`well${big ? ' big' : ''}`}>
      <div className="wh"><span className="lb">Note to the mission</span><span className="sm">{r.type} · v{i.a.mission?.version ?? r.version}</span></div>
      <textarea className="note" placeholder="Optional. Sent with the reply." aria-label="Note to the mission"
        value={i.a.text} onChange={(e) => i.a.setText(e.target.value)} disabled={i.a.busy} rows={3} />
      {!big && i.picked!.artifacts.map((f) => (
        <div className="row" key={f}><span className="f1">{f.split('/').pop()}</span><span className="sm">artifact</span></div>
      ))}
    </div>
  )
}

function Verbs({ i, hide, cls }: { i: I; hide?: boolean; cls: string }) {
  const item = i.a.item
  const r = i.picked!
  return (
    <div className={`btnrow ${cls}`}>
      {item && i.a.mission && item.verbs.map((v, n) => (
        <button key={v} type="button" className={`btn${n === 0 ? ' p' : REJECT.includes(v) ? ' x' : ''}`}
          onClick={() => i.a.send(v)} disabled={i.a.busy} aria-busy={i.a.busy}>{n === 0 && i.a.busy ? 'Sending…' : verbLabel(v)}</button>
      ))}
      {hide && r.state === 'open' && <button type="button" className="btn x hide" onClick={() => i.hide(r)} disabled={i.a.busy}>Hide</button>}
    </div>
  )
}

// ---- fold inner and laptop (board fold-inbox-inner) ----
function Row({ i, r }: { i: I; r: InboxItem }) {
  const m = i.missionOf(r)
  if (r.type === 'fyi') {
    return (
      <button type="button" className={`row${r.id === i.sel ? ' on' : ''}`} aria-pressed={r.id === i.sel} onClick={() => i.setSel(r.id)}>
        <span className="sm num">{time(r.created)}</span><span className="f1">{r.project} · {r.ask}</span>
      </button>
    )
  }
  return (
    <button type="button" className={`row two${r.id === i.sel ? ' on' : ''}`} aria-pressed={r.id === i.sel} onClick={() => i.setSel(r.id)}>
      <span className="l1"><span className="f1">{m ? title(m) : short(r.mission)}</span><span className="sm num">{time(r.created)}</span></span>
      <span className="sm f1">{r.ask}</span>
    </button>
  )
}

function Group({ i, label, cls, rows, later }: { i: I; label: string; cls: string; rows: InboxItem[]; later?: boolean }) {
  return (
    <>
      <div className={`grp ${cls}`}><span>{label}</span><span className="num">{later ? '' : rows.length}</span></div>
      {later ? <div className="sm later">{LATER}</div> : rows.map((r) => <Row key={r.id} i={i} r={r} />)}
    </>
  )
}

function InboxInner({ i }: { i: I }) {
  const r = i.picked
  const m = r && i.missionOf(r)
  const all = i.asks.length + i.fyis.length
  const segs: Array<[Seg, string]> = [['all', `All ${all}`], ['answer', 'Answer'], ['drafts', 'Drafts'], ['fyi', 'FYI']]
  const show = (s: Seg) => i.seg === 'all' || i.seg === s
  return (
    <div className="desk desk-in ibx">
      <div className="col">
        <div className="segs" role="group" aria-label="Filter">
          {segs.map(([k, l]) => (
            <button key={k} type="button" className={`seg${i.seg === k ? ' on' : ''}`} aria-pressed={i.seg === k} onClick={() => i.setSeg(k)}>{l}</button>
          ))}
        </div>
        <section className="pn grow list-pn" aria-label="Inbox">
          {show('answer') && <Group i={i} label="To answer" cls="c-n" rows={i.asks} />}
          {show('drafts') && <Group i={i} label="Drafts ready" cls="c-m" rows={[]} later />}
          {show('fyi') && <Group i={i} label="FYI" cls="c-q" rows={i.fyis} />}
        </section>
        <section className="pn none" aria-labelledby="tr">
          <h2 className="ph" id="tr"><span>Triage</span></h2>
          <div className="row"><Diamond /><span className="f1">Rules</span><Link className="lk" href="/reports">Reports ›</Link></div>
          <div className="sm">{LATER} · crow never sends</div>
        </section>
      </div>
      {r ? (
        <section className="pn grow read" aria-labelledby="rt" aria-live="polite">
          <h2 className="ph" id="rt"><span className="rt">{m ? title(m) : short(r.mission)}</span><span className="sm num">{stamp(r.created)}</span></h2>
          <div className="sm"><From r={r} m={m} /> · <Link className="lk" href={`/missions/${r.mission}`}>Mission ›</Link></div>
          <div className="body">{(i.a.item ?? r).ask}</div>
          <Stale i={i} />
          {r.state === 'open' && r.type !== 'fyi' && <Note i={i} />}
          <div className="sm">source: {r.file}</div>
          {r.state === 'open' && <Verbs i={i} hide cls="end act" />}
        </section>
      ) : (
        <section className="pn grow read"><div className="sm">Inbox is empty.</div></section>
      )}
    </div>
  )
}

// ---- fold cover (board fold-inbox-cover) ----
function InboxCover({ i }: { i: I }) {
  const r = i.picked
  const m = r && i.missionOf(r)
  const rest = i.asks.filter((x) => x.id !== r?.id)
  return (
    <div className="desk desk-cv ibx">
      {r ? (
        <section className="gcard dhero ask" aria-labelledby="ct">
          <div className="gt"><span className={r.type === 'fyi' ? 'c-q' : 'c-n'}>{r.type === 'fyi' ? 'FYI' : `To answer · ${i.asks.length}`}</span><span className="num">{time(r.created)}</span></div>
          <div className="who">
            <span className="av" aria-hidden="true">{initials(r.project)}</span>
            <div className="hm"><h2 className="gh" id="ct">{m ? title(m) : short(r.mission)}</h2><span className="gs"><From r={r} m={m} /></span></div>
          </div>
          <div className="body">{(i.a.item ?? r).ask}</div>
          <Stale i={i} />
          {r.state === 'open' && r.type !== 'fyi' && <Note i={i} big />}
        </section>
      ) : (
        <section className="gcard dhero"><div className="gt"><span>Inbox</span></div><div className="gs">Nothing to answer.</div></section>
      )}
      <div className="gl">
        <div className="gls"><span>Also today</span><span className="num">{rest.length + i.fyis.length}</span></div>
        {rest.slice(0, 2).map((x) => {
          const xm = i.missionOf(x)
          return (
            <button key={x.id} type="button" className="row r2" onClick={() => i.setSel(x.id)}>
              <span className="ic pen"><Pen /></span>
              <span className="t"><b>{xm ? title(xm) : short(x.mission)}</b><span>{x.ask}</span></span>
              <span className="sm num">{time(x.created)}</span><span className="chev">›</span>
            </button>
          )
        })}
        {i.fyis.length > 0 && (
          <button type="button" className="row r2" onClick={() => i.setSel(i.fyis[0].id)}>
            <span className="ic"><Env /></span>
            <span className="t"><b>{i.fyis.length} FYI</b><span>{[...new Set(i.fyis.map((x) => x.project))].join(', ')}</span></span>
            <span className="chev">›</span>
          </button>
        )}
        <Link className="row r2" href="/reports">
          <span className="ic tri"><Diamond /></span>
          <span className="t"><b>Triage rules</b><span>{LATER}</span></span>
          <span className="chev">›</span>
        </Link>
      </div>
      <div className="uf"><Unfold />unfold to read and hide</div>
      {r && r.state === 'open' && r.type !== 'fyi' && <div className="dock"><Verbs i={i} cls="gbtns" /></div>}
    </div>
  )
}

export function Inbox({ snap, refresh }: { snap: Snapshot; refresh: () => Promise<void> }) {
  const i = useInbox(snap, refresh)
  const inner = useMedia(INNER)
  return (
    <>
      {inner ? <InboxInner i={i} /> : <InboxCover i={i} />}
      {i.sheet && <Sheet c={i.sheet} onClose={() => i.setSheet(null)} />}
    </>
  )
}
