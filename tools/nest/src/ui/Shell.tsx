import { useEffect, useState, type ReactNode } from 'react'
import { countByState, openInbox } from '../data/state'
import type { Snapshot } from '../data/useData'
import { time } from '../format'
import { Link } from '../router'
import { Glyph } from './Glyph'
import { TABS, TabIcon, type TabKey } from './icons'

type Theme = 'modern' | 'cyber'

function useTheme(): [Theme, () => void] {
  const [theme, setTheme] = useState<Theme>(() => (document.documentElement.dataset.theme === 'cyber' ? 'cyber' : 'modern'))
  useEffect(() => {
    document.documentElement.dataset.theme = theme
    document.querySelector('meta[name="theme-color"]')?.setAttribute('content', theme === 'cyber' ? '#120b2c' : '#24273a')
    try { localStorage.setItem('nest.theme', theme) } catch { /* storage blocked: theme stays per visit */ }
  }, [theme])
  return [theme, () => setTheme((t) => (t === 'modern' ? 'cyber' : 'modern'))]
}

function badges(snap: Snapshot | null): Partial<Record<TabKey, { n: number; failed?: boolean }>> {
  if (!snap) return {}
  const c = countByState(snap.missions, snap.inbox)
  const asks = openInbox(snap.inbox).filter((r) => r.type !== 'fyi').length
  return { missions: { n: c.failed, failed: true }, inbox: { n: asks } }
}

export function Shell({ active, snap, children }: { active: TabKey; snap: Snapshot | null; children: ReactNode }) {
  const [theme, toggle] = useTheme()
  const b = badges(snap)
  const c = snap ? countByState(snap.missions, snap.inbox) : null
  return (
    <div className="app">
      <header className="top">
        <Link href="/" className="logo" aria-label="CR-0W home">CR-0W</Link>
        <nav className="pills" aria-label="Tabs">
          {TABS.map((t) => (
            <Link key={t.key} href={t.path} className="pill" aria-current={t.key === active ? 'page' : undefined}>
              {t.label}
              {b[t.key]?.n ? <span className={`badge-t num${b[t.key]!.failed ? ' f' : ''}`}>{b[t.key]!.n}</span> : null}
            </Link>
          ))}
        </nav>
        <span className="grow" />
        {c && (
          <span className="chips wrap" aria-label={`${c.needs} need you, ${c.working} working`}>
            <span className="chip c-n"><Glyph state="needs" size={12} />{c.needs}</span>
            <span className="chip c-w"><Glyph state="working" size={12} />{c.working}</span>
          </span>
        )}
        {snap && <span className="clock">{time(snap.now)}</span>}
        <button type="button" className="icon-btn" onClick={toggle} aria-label={`Theme: ${theme}. Switch to ${theme === 'modern' ? 'cyber' : 'modern'}`} title="Switch theme">
          <svg aria-hidden="true" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
            <circle cx="12" cy="12" r="8" /><path d="M12 4a8 8 0 0 1 0 16Z" fill="currentColor" />
          </svg>
        </button>
      </header>
      <main className="main" id="main">{children}</main>
      <nav className="tabbar" aria-label="Tabs">
        {TABS.map((t) => (
          <Link key={t.key} href={t.path} className="tab" aria-current={t.key === active ? 'page' : undefined}>
            <TabIcon>{t.icon}</TabIcon>
            {t.label}
            {b[t.key]?.n ? (
              <span className={`badge${b[t.key]!.failed ? ' f' : ''}`} aria-label={`${b[t.key]!.n} ${b[t.key]!.failed ? 'failed' : 'to answer'}`}>{b[t.key]!.n}</span>
            ) : null}
          </Link>
        ))}
      </nav>
    </div>
  )
}
