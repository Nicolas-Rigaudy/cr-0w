import { useCallback, useEffect, useId, useRef, useState, type FocusEvent, type ReactNode } from 'react'
import { countByState } from '../data/state'
import { homeCounts, quotaCls, quotaOf } from '../data/home'
import { inboxBadge } from '../data/inbox'
import type { Snapshot } from '../data/useData'
import { ago, headerClock, time } from '../format'
import { Link, usagePath } from '../router'
import { Glyph } from './Glyph'
import { MORE_ICON, SCOPE_ICONS, SETTINGS_ICON, TabIcon, type Tab, type TabKey } from './icons'
import { barTab, hiddenTabs, SETTINGS } from './tabs'
import { ALL_TABS, NAV, splitBar } from '../modules/nav'
import { themes } from 'virtual:nest-modules'
import { nextScope, scopeChip, type ScopeMode } from '../data/scope'
import { MODES, paintTheme, setMode, setTheme, useMode, useTheme, watchSystemMode } from './theme'

// Themes are token files (item_267): the built-in ones and those of on bricks, in that order.
// The switcher cycles them; crow.yaml modules.nest.theme is the first paint.
export function nextTheme(t: string, list: readonly string[] = themes): string {
  return list[(list.indexOf(t) + 1) % list.length]
}

function useThemeCycle(): string {
  const theme = useTheme()
  // First paint sets the browser bar colour too.
  useEffect(() => paintTheme(theme), [theme])
  // Mode auto follows the system live.
  useEffect(watchSystemMode, [])
  return theme
}

const title = (t: string) => t.charAt(0).toUpperCase() + t.slice(1)

// Header theme picker: every installed theme with its swatch, then the light/dark mode; picks apply at once.
function ThemeMenu({ theme }: { theme: string }) {
  const [open, setOpen] = useState(false)
  const mode = useMode()
  const box = useRef<HTMLSpanElement>(null)
  useEffect(() => {
    if (!open) return
    const off = (e: Event) => { if (!box.current?.contains(e.target as Node)) setOpen(false) }
    const esc = (e: KeyboardEvent) => { if (e.key === 'Escape') setOpen(false) }
    document.addEventListener('pointerdown', off)
    document.addEventListener('keydown', esc)
    return () => { document.removeEventListener('pointerdown', off); document.removeEventListener('keydown', esc) }
  }, [open])
  return (
    <span className="theme-menu" ref={box}>
      <button type="button" className="icon-btn" aria-haspopup="menu" aria-expanded={open} aria-label={`Theme: ${theme}`} title="Theme" onClick={() => setOpen(!open)}>
        <svg aria-hidden="true" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
          <path d="M12 3a9 9 0 1 0 0 18c1.2 0 2-.8 2-1.8 0-.5-.2-.9-.5-1.2-.3-.3-.5-.7-.5-1.2 0-1 .8-1.8 1.8-1.8H17a4 4 0 0 0 4-4C21 6.6 17 3 12 3Z" />
          <circle cx="7.5" cy="11" r="1.2" fill="currentColor" /><circle cx="10" cy="7" r="1.2" fill="currentColor" /><circle cx="15" cy="7.5" r="1.2" fill="currentColor" />
        </svg>
      </button>
      {open && (
        <span className="theme-pop" role="menu">
          {themes.map((t) => (
            <button key={t} type="button" role="menuitemradio" aria-checked={t === theme} className="theme-opt" onClick={() => { setTheme(t); setOpen(false) }}>
              <span className="set-mini" data-theme={t} aria-hidden="true"><i /><span><i /><i /><i /></span></span>
              <span className="n">{title(t)}{t === theme && <span className="ck" aria-hidden="true"> ✓</span>}</span>
            </button>
          ))}
          <span className="theme-mode">
            <span className="n">Mode</span>
            <span className="seg sm" role="group" aria-label="Mode">
              {MODES.map((m) => (
                <button key={m} type="button" className={m === mode ? 'on' : ''} aria-pressed={m === mode} onClick={() => setMode(m)}>{title(m)}</button>
              ))}
            </span>
          </span>
        </span>
      )}
    </span>
  )
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
  const theme = useThemeCycle()
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
          {pills.map((t) => (
            <Link key={t.key} href={t.path} className="pill" aria-current={t.key === active ? 'page' : undefined}>
              {t.label}
              {b[t.key]?.n ? <span className={`badge-t num${b[t.key]!.failed ? ' f' : ''}`}>{b[t.key]!.n}</span> : null}
            </Link>
          ))}
        </nav>
        <span className="grow" />
        {accts.length > 0 && (
          <span className="accts" role="group" aria-label="Quota left">
            {accts.map((q) => (
              <Link key={q.account} href={usagePath(q.account)} className="acct" title={`${q.account}: ${q.h}% of 5 hours, ${q.w}% of the week left`}>
                {q.account}
                <span className="gb" aria-hidden="true"><i className={quotaCls(q.h)} style={{ width: `${q.h}%` }} /></span>
                <span className={`pc num${quotaCls(q.w) === 'lo' ? ' c-f' : ''}`} aria-label={`${q.h}% of 5 hours, ${q.w}% of the week left`}>{q.h}%</span>
              </Link>
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
        <ThemeMenu theme={theme} />
        <Link href="/settings" className="icon-btn set-ic" aria-label="Settings" title="Settings" aria-current={active === SETTINGS ? 'page' : undefined}>
          <TabIcon>{SETTINGS_ICON}</TabIcon>
        </Link>
      </header>
      <main className="main" id="main">{children}</main>
      {rail}
      <TabBar active={active} b={b} hidden={hidden} />
    </div>
  )
}

type Badges = ReturnType<typeof badges>

const TAB_OF = new Map<string, Tab>(ALL_TABS.map((t) => [t.key, t]))
const SETTINGS_TAB: Tab = { key: SETTINGS, label: 'Settings', path: '/settings', icon: SETTINGS_ICON }

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
  // Settings closes the More sheet on every phone (item_284), so More always shows.
  const more = [...split.more.map((k) => TAB_OF.get(k)!), SETTINGS_TAB]
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
      <button
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
      </button>
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
