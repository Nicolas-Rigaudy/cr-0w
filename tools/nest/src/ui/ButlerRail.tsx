import { useEffect, useMemo, useRef, useState } from 'react'
import { isCommand, type ButlerMsg } from '../data/butler'
import { accountsOf } from '../data/live'
import { useButler } from '../data/useButler'
import type { Snapshot } from '../data/useData'
import { time } from '../format'
import { heroLiveOf, heroModel, type HeroModel } from '../hero/model'
import { RailHead } from '../hero/RailHead'
import { Link, navigate } from '../router'
import { enterSubmits, rowsOf } from './keys'
import { HomeHud } from './HomeHud'
import { HOME_W, HOME_W_KEY, RAIL_W, RAIL_W_KEY, useRailWidth } from './railWidth'
import { stickBottom } from './stickBottom'
import { useDraft } from './draft'

const FOLD_KEY = 'nest.rail'
// Home keeps its own fold (item_259).
const HOME_FOLD_KEY = 'nest.rail.home'

type Fold = 'folded' | 'open' | null

// The owner's fold choice, remembered; null until they pick one.
function useFold(key: string, auto: boolean): [boolean, () => void] {
  const [pick, setPick] = useState<Fold>(() => {
    try { const v = localStorage.getItem(key); return v === 'folded' || v === 'open' ? v : null } catch { return null }
  })
  const folded = pick ? pick === 'folded' : auto
  const now = useRef(folded)
  now.current = folded
  const toggle = () => {
    const next: Fold = now.current ? 'open' : 'folded'
    try { localStorage.setItem(key, next) } catch { /* storage off: fold lasts the visit */ }
    setPick(next)
  }
  // Ctrl B is Nest's one shortcut: it folds and opens the butler.
  useEffect(() => {
    const key = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && !e.altKey && !e.shiftKey && e.code === 'KeyB') { e.preventDefault(); toggle() }
    }
    window.addEventListener('keydown', key)
    return () => window.removeEventListener('keydown', key)
  }, [])
  return [folded, toggle]
}

// Mood label subject: the failed project, the first one asking, or the live count.
export function railLabel(m: HeroModel): string {
  const by = (st: string) => m.projects.find((p) => p.state === st)
  if (m.mood === 'fail') return `fail · ${by('failed')?.name ?? 'run'}`
  if (m.mood === 'ask') return `ask · ${by('needs')?.name ?? `${m.asks} waiting`}`
  if (m.mood === 'run') return `run · ${m.projects.filter((p) => p.state === 'working').length} live`
  return 'idle'
}

const lower = (s: string) => s.charAt(0).toLowerCase() + s.slice(1)

// Note text with `code` spans: a login failure carries the command to copy.
function noteText(t: string) {
  return t.split('`').map((part, i) => (i % 2 ? <code key={i} className="rl-code">{part}</code> : part))
}

function Msg({ m }: { m: ButlerMsg }) {
  if (m.role === 'tool') return <li className="rl-msg tool">{m.text}</li>
  if (m.role === 'note') return <li className="rl-msg note">{noteText(m.text)}</li>
  const you = m.role === 'you'
  return (
    <li className={`rl-msg ${you ? 'you' : 'cr'}`}>
      <span className="rl-who">{you ? '' : <i className="dot" aria-hidden="true" />}{you ? 'you' : 'cr-0w'} · {time(m.ts)}</span>
      <p>{m.text}</p>
      {m.chips?.length ? (
        <span className="rl-chips">{m.chips.map((c) => <Link key={c.href + c.label} href={c.href} className="rl-chip">{c.label}</Link>)}</span>
      ) : null}
    </li>
  )
}

// Butler rail (item_176, board desk-rail-a): CR-0W on the right of every laptop screen.
// Hero on top, the hero line, the conversation, then the account and the input. A grip on the
// left edge sets the width (item_256). The fold button folds it to a strip, remembered.
// Home (item_259, board H3): the right third next to the hero, open and 480px wide by default,
// with its own fold and width; the HUD card (needs you, next, running) replaces head and line.
export function ButlerRail({ snap, home = false }: { snap: Snapshot; home?: boolean }) {
  const b = useButler()
  const [folded, fold] = useFold(home ? HOME_FOLD_KEY : FOLD_KEY, false)
  const m = useMemo(() => heroModel(snap.missions, snap.inbox, heroLiveOf(snap.sections), snap.now), [snap])
  const working = m.projects.some((p) => p.state === 'working')
  const grip = useRailWidth(home ? HOME_W_KEY : RAIL_W_KEY, home ? HOME_W : RAIL_W)

  if (!b) return null
  const label = railLabel(m)

  if (folded) {
    return (
      <aside className="b-rail folded" aria-label="CR-0W">
        <button type="button" className="rl-strip" onClick={fold} title="Open CR-0W" aria-expanded="false">
          <i className={`dia m-${m.mood}`} aria-hidden="true" />
          <span className="rl-vert">{label}</span>
          {b.busy && <i className="rl-busy" aria-label="CR-0W is answering" />}
        </button>
      </aside>
    )
  }

  return (
    <aside className={`b-rail${home ? ' home-rail' : ''}`} aria-label="CR-0W">
      <div className="rl-grip" {...grip} />
      <div className={home ? 'rl-hdr' : 'rl-hdr plain'}>
        <div className="rl-top">
          <span className={`rl-mood m-${m.mood}`}><i className="dia" aria-hidden="true" />{label}</span>
          <button type="button" className="rl-kbd" onClick={fold} title="Fold (Ctrl B)" aria-expanded="true">Fold</button>
        </div>
        {home && <HomeHud snap={snap} />}
      </div>
      {!home && <RailHead mood={m.mood} working={working} asks={m.asks} />}
      {!home && <p className="rl-line"><b>cr-0w</b> › {lower(m.line)}<span className="cur" aria-hidden="true">▌</span></p>}
      <ButlerChat snap={snap} fold={fold} framed={home} />
    </aside>
  )
}

// Clear starts a new conversation. The first click arms it for a few seconds, the second clears.
const ARM_MS = 3000

function ClearChat({ clear, empty }: { clear: () => Promise<void>; empty: boolean }) {
  const [armed, setArmed] = useState(false)
  useEffect(() => {
    if (!armed) return
    const t = setTimeout(() => setArmed(false), ARM_MS)
    return () => clearTimeout(t)
  }, [armed])
  return (
    <button type="button" className={`rl-clear${armed ? ' armed' : ''}`} disabled={empty}
      title="Start a new conversation: CR-0W forgets this one" onClick={() => { if (!armed) return setArmed(true); setArmed(false); void clear() }}>
      {armed ? 'Clear it?' : 'Clear'}
    </button>
  )
}

// The conversation and its input (the rail body). Framed (Home): a titled panel with Clear in
// its header; otherwise Clear sits in the account row.
export function ButlerChat({ snap, fold, close, framed = false }: { snap: Snapshot; fold?: () => void; close?: () => void; framed?: boolean }) {
  const b = useButler()
  const [text, setText] = useDraft('butler')
  const [acct, setAcct] = useState<string | null>(null)
  const log = useRef<HTMLOListElement>(null)
  const accounts = useMemo(() => accountsOf(snap.sections.accounts), [snap.sections.accounts])
  const stick = useRef<ReturnType<typeof stickBottom> | null>(null)
  const n = b?.messages.length ?? 0
  const has = !!b

  useEffect(() => {
    const el = log.current
    if (!el) return
    stick.current = stickBottom(el)
    return () => { stick.current?.stop(); stick.current = null }
  }, [has])
  useEffect(() => { stick.current?.follow() }, [n, b?.partial, b?.busy])

  if (!b) return null
  const send = () => {
    const t = text.trim()
    if (!t) return
    setText('')
    if (isCommand(t)) return navigate(`/desk?cmd=${encodeURIComponent(t)}`)
    void b.say(t, 'nest', acct ?? undefined)
  }

  const empty = !b.messages.length && !b.busy
  const body = (
    <>
      <ol className="rl-log" ref={log} aria-live="polite">
        {b.messages.map((x) => <Msg key={x.seq} m={x} />)}
        {b.busy && (
          <li className="rl-msg cr typing">
            <span className="rl-who"><i className="dot" aria-hidden="true" />cr-0w</span>
            <p>{b.partial || '…'}</p>
          </li>
        )}
      </ol>
      {b.error && <p className="rl-err c-f" role="alert">{b.error}</p>}
      <form className="rl-foot" onSubmit={(e) => { e.preventDefault(); send() }}>
        <div className="rl-meta">
          <label>account
            <select value={acct ?? b.account} onChange={(e) => setAcct(e.target.value)}>
              <option value="">auto</option>
              {accounts.map((a) => <option key={a.name} value={a.name}>{a.name}</option>)}
            </select>
          </label>
          {!framed && <ClearChat clear={b.clear} empty={empty} />}
          {fold && <button type="button" className="bare sm" onClick={fold}>Fold</button>}
          {close && <button type="button" className="rl-close" onClick={close} aria-label="Close the chat" title="Close"><i aria-hidden="true" /></button>}
        </div>
        <div className="rl-in">
          <textarea value={text} rows={rowsOf(text, 4)} placeholder="Talk to CR-0W, or /dispatch…" aria-label="Message CR-0W"
            onChange={(e) => setText(e.target.value)} onKeyDown={enterSubmits} />
          <button type="submit" disabled={!text.trim()}>Send</button>
        </div>
      </form>
    </>
  )
  if (!framed) return body
  return (
    <section className="hpn hh-chat" aria-label="Chat with CR-0W">
      <div className="hpn-h">Chat with CR-0W<span className="r"><ClearChat clear={b.clear} empty={empty} /></span></div>
      {body}
    </section>
  )
}

