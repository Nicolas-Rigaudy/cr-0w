import { useCallback, useEffect, useId, useRef, useState, type FocusEvent, type ReactNode } from 'react'
import { countByState } from '../data/state'
import { homeCounts, quotaCls, quotaOf } from '../data/home'
import { inboxBadge } from '../data/inbox'
import type { Snapshot } from '../data/useData'
import { ago, headerClock, time } from '../format'
import { Link } from '../router'
import { Glyph } from './Glyph'
import { MORE_ICON, SCOPE_ICONS, TabIcon, type Tab, type TabKey } from './icons'
import { barTab, hiddenTabs } from './tabs'
import { ALL_TABS, NAV, splitBar } from '../modules/nav'
import { config, themes } from 'virtual:nest-modules'
import { nextScope, scopeChip, type ScopeMode } from '../data/scope'

// Themes are token files (item_267): the built-in ones and those of on bricks, in that order.
// The switcher cycles them; crow.yaml modules.nest.theme is the first paint.
export function nextTheme(t: string, list: readonly string[] = themes): string {
  return list[(list.indexOf(t) + 1) % list.length]
}

function useTheme(): [string, () => void] {
  const [theme, setTheme] = useState(() => {
    const t = document.documentElement.dataset.theme
    return t && themes.includes(t) ? t : config.theme
  })
  useEffect(() => {
    document.documentElement.dataset.theme = theme
    const bg = getComputedStyle(document.documentElement).getPropertyValue('--bg').trim()
    if (bg) document.querySelector('meta[name="theme-color"]')?.setAttribute('content', bg)
    try { localStorage.setItem('nest.theme', theme) } catch { /* storage blocked: theme stays per visit */ }
  }, [theme])
  return [theme, () => setTheme((t) => nextTheme(t))]
}

function badges(snap: Snapshot | null): Partial<Record<TabKey, { n: number; failed?: boolean }>> {
  if (!snap) return {}
  const c = countByState(snap.missions, snap.inbox)
  return { ops: { n: c.failed, failed: true }, inbox: inboxBadge(snap) }
}

// Header A (item_206): one chip after the logo cycles the scope filter.
function ScopeChip({ scope, setScope }: { scope: ScopeMode; setScope: (m: ScopeMode) => void }) {
  const t = scopeChip(scope)
  return (
    <button type="button" className="scope-chip" aria-label={t.name} aria-pressed={t.pressed} title={t.name} onClick={() => setScope(nextScope(scope))}>
      <svg aria-hidden="true" width="14" height="14" viewBox="0 0 22 22" fill="none" stroke="currentColor" strokeWidth="1.9" strokeLinejoin="round" strokeLinecap="round">
        {SCOPE_ICONS[scope]}
      </svg>
      {t.label}
    </button>
  )
}

interface ShellProps { active: TabKey; snap: Snapshot | null; scope: ScopeMode; setScope: (m: ScopeMode) => void; rail?: ReactNode; badge?: ReactNode; children: ReactNode }

export function Shell({ active, snap, scope, setScope, rail, badge, children }: ShellProps) {
  const [theme, toggle] = useTheme()
  const b = badges(snap)
  const c = snap ? homeCounts(snap) : null
  // Account gauges (laptop), allowed accounts only: 5h left as bar and number, week in the
  // title; the number turns red when the week runs low.
  const accts = snap ? quotaOf(snap) : []
  const hidden = hiddenTabs(snap)
  const pills = NAV.map((k) => TAB_OF.get(k)!).filter((t) => !hidden.has(t.key))
  return (
    <div className={`app${rail ? ' has-rail' : ''}`}>
      <header className="top">
        {badge}
        <Link href="/" className="logo" aria-label="CR-0W home">CR-0W</Link>
        <ScopeChip scope={scope} setScope={setScope} />
        <span className="hdr-sep" aria-hidden="true" />
        <nav className="pills" aria-label="Tabs">
          {pills.map((t, i) => (
            <Link key={t.key} href={t.path} className="pill" aria-current={t.key === active ? 'page' : undefined}>
              <span className="tn">{i + 1}</span>{t.label}
              {b[t.key]?.n ? <span className={`badge-t num${b[t.key]!.failed ? ' f' : ''}`}>{b[t.key]!.n}</span> : null}
            </Link>
          ))}
        </nav>
        <span className="grow" />
        {accts.length > 0 && (
          <span className="accts" role="group" aria-label="Quota left">
            {accts.map((q) => (
              <span key={q.account} className="acct" title={`${q.account}: ${q.h}% of 5 hours, ${q.w}% of the week left`}>
                {q.account}
                <span className="gb" aria-hidden="true"><i className={quotaCls(q.h)} style={{ width: `${q.h}%` }} /></span>
                <span className={`pc num${quotaCls(q.w) === 'lo' ? ' c-f' : ''}`} aria-label={`${q.h}% of 5 hours, ${q.w}% of the week left`}>{q.h}%</span>
              </span>
            ))}
          </span>
        )}
        {c && (
          <span className="chips wrap" aria-label={`${c.needs} need you, ${c.working} working`}>
            <span className="chip c-n"><Glyph state="needs" size={12} />{c.needs}</span>
            <span className="chip c-w"><Glyph state="working" size={12} />{c.working}</span>
          </span>
        )}
        {snap?.down && (
          <span className="link-down" role="status">
            <span className="long">core </span>down{snap.synced ? <> · <span className="num">{ago(snap.synced.toISOString(), snap.now)}</span></> : null}
          </span>
        )}
        {snap?.newBuild && (
          <button type="button" className="chip build-chip" onClick={() => location.reload()}>new build ready · reload</button>
        )}
        {snap && <span className="clock"><span className="short">{time(snap.now)}</span><span className="long">{headerClock(snap.now, theme === 'cyber')}</span></span>}
        <button type="button" className="icon-btn" onClick={toggle} aria-label={`Theme: ${theme}. Switch to ${nextTheme(theme)}`} title="Switch theme">
          <svg aria-hidden="true" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
            <circle cx="12" cy="12" r="8" /><path d="M12 4a8 8 0 0 1 0 16Z" fill="currentColor" />
          </svg>
        </button>
      </header>
      <main className="main" id="main">{children}</main>
      {rail}
      <TabBar active={active} b={b} hidden={hidden} />
    </div>
  )
}

type Badges = ReturnType<typeof badges>

const TAB_OF = new Map<string, Tab>(ALL_TABS.map((t) => [t.key, t]))

// Phone tab bar B (item_206): Home, Desk, Missions, Inbox, More. More opens a sheet
// with MRs, Agenda and Reports and stays highlighted while one of them is open.
function TabBar({ active, b, hidden }: { active: TabKey; b: Badges; hidden: ReadonlySet<TabKey> }) {
  const [open, setOpen] = useState(false)
  const btn = useRef<HTMLButtonElement>(null)
  const sheet = useRef<HTMLUListElement>(null)
  const id = useId()
  const cur = barTab(active, NAV, hidden)
  const close = useCallback((refocus: boolean) => {
    setOpen(false)
    if (refocus) btn.current?.focus()
  }, [])
  useEffect(() => {
    if (!open) return
    const items = sheet.current
    ;(items?.querySelector<HTMLElement>('[aria-current="page"]') ?? items?.querySelector<HTMLElement>('a'))?.focus()
    const key = (e: KeyboardEvent) => { if (e.key === 'Escape') { e.preventDefault(); close(true) } }
    const down = (e: PointerEvent) => {
      const t = e.target as Node
      if (!items?.contains(t) && !btn.current?.contains(t)) close(false)
    }
    addEventListener('keydown', key)
    document.addEventListener('pointerdown', down)
    return () => {
      removeEventListener('keydown', key)
      document.removeEventListener('pointerdown', down)
    }
  }, [open, close])
  const blur = (e: FocusEvent<HTMLUListElement>) => {
    const to = e.relatedTarget as Node | null
    if (to && !sheet.current?.contains(to) && to !== btn.current) close(false)
  }
  const split = splitBar(NAV, hidden)
  const tabs = split.bar.map((k) => TAB_OF.get(k)!)
  const more = split.more.map((k) => TAB_OF.get(k)!)
  const openMore = more.find((t) => t.key === active)
  return (
    <nav className="tabbar" aria-label="Tabs">
      {tabs.map((t) => (
        <Link key={t.key} href={t.path} className="tab" aria-current={t.key === cur ? 'page' : undefined}>
          <TabIcon>{t.icon}</TabIcon>
          {t.label}
          {b[t.key]?.n ? (
            <span className={`badge${b[t.key]!.failed ? ' f' : ''}`} aria-label={`${b[t.key]!.n} ${b[t.key]!.failed ? 'failed' : 'to answer'}`}>{b[t.key]!.n}</span>
          ) : null}
        </Link>
      ))}
      {more.length > 0 && <button
        ref={btn}
        type="button"
        className="tab"
        aria-expanded={open}
        aria-controls={id}
        aria-current={cur === 'more' ? 'page' : undefined}
        onClick={() => (open ? close(false) : setOpen(true))}
      >
        <TabIcon>{MORE_ICON}</TabIcon>
        More{openMore && <span className="sr-only">, {openMore.label} open</span>}
      </button>}
      {open && (
        <ul ref={sheet} id={id} className="more-sheet" aria-label="More screens" onBlur={blur}>
          {more.map((t) => (
            <li key={t.key}>
              <Link href={t.path} className="more-item" aria-current={t.key === active ? 'page' : undefined} onClick={() => close(true)}>
                <TabIcon>{t.icon}</TabIcon>
                {t.label}
              </Link>
            </li>
          ))}
        </ul>
      )}
    </nav>
  )
}
