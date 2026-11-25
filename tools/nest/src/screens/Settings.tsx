import { useEffect, useRef, type PointerEvent, type KeyboardEvent, type ReactNode } from 'react'
import { useSource } from '../data/useData'
import {
  changedKeys, errorLines, HOME, movePage, needsBuild, pageMark, pendingCount, togglePage, type Draft, type SettingKey,
} from '../data/settings'
import { useSettings, type SettingsControl, type SettingsView } from '../data/useSettings'
import { ALL_TABS } from '../modules/nav'
import { Link, useQuery } from '../router'
import { PANE_ICONS, TabIcon } from '../ui/icons'
import { paintTheme, setTheme, unpaintTheme } from '../ui/theme'

// Nest settings (item_284, board D): panes in a sidebar on laptop, a pane list on phone.
// Edits stay in a draft until Save; the theme previews live, hero and pages need a build.

export type PaneKey = 'look' | 'pages' | 'modules' | 'general'

export const PANES: Array<{ key: PaneKey; label: string; hint: string }> = [
  { key: 'look', label: 'Look', hint: 'Theme and hero of the Nest' },
  { key: 'pages', label: 'Pages', hint: 'Which pages show and in what order' },
  { key: 'modules', label: 'Modules', hint: 'Bricks that run on this hub' },
  { key: 'general', label: 'General', hint: 'Language, time zone, refresh' },
]

const THEME_HINTS: Record<string, string> = {
  modern: 'Soft purple on dark, the default',
  cyber: 'Neon green and pink on black',
}
const PAGE_HINTS: Record<string, string> = {
  home: 'Always on',
  desk: 'Sessions, runs and picks',
  ops: 'Missions and the run timeline',
  inbox: 'Questions and asks waiting for you',
  mrs: 'Merge requests to review or fix',
  agenda: 'Calendar of the day',
  reports: 'Weekly and daily rollups',
}
const title = (k: string) => k.charAt(0).toUpperCase() + k.slice(1).replace(/-/g, ' ')
const pageLabel = (k: string) => ALL_TABS.find((t) => t.key === k)?.label ?? title(k)

export function paneOf(q: string | null, laptop: boolean): PaneKey | null {
  const p = PANES.find((x) => x.key === q)?.key
  return p ?? (laptop ? 'look' : null)
}

// Panes with pending edits get a dot in the nav.
export function paneDirty(base: Draft, draft: Draft, pane: PaneKey): boolean {
  if (pane === 'look') return base.theme !== draft.theme || base.hero !== draft.hero
  if (pane === 'pages') return draft.pages.some((p) => pageMark(base, draft, p.key))
  return false
}

const BUILD_WHAT: Record<string, string> = { hero: 'Hero', pages: 'Pages' }
const plural = (n: number, w: string) => `${n} ${w}${n === 1 ? '' : 's'}`

// The save bar's text: pending edits first, then the save or build that follows.
export function barText(v: SettingsView): { n: string; hint: string; tone: 'pending' | 'busy' | 'ok' | 'failed' } | null {
  const { base, draft, build, saved } = v
  if (!base || !draft) return null
  const n = pendingCount(base, draft)
  if (n) {
    const keys = changedKeys(base, draft).filter((k) => k !== 'theme')
    const hint = keys.length
      ? `${keys.map((k) => BUILD_WHAT[k]).join(' and ')} need${keys.length === 1 && keys[0] !== 'pages' ? 's' : ''} a rebuild, about 40 s`
      : 'The theme applies at once'
    return { n: plural(n, 'change'), hint: v.saveError ?? hint, tone: v.saveError ? 'failed' : 'pending' }
  }
  if (v.saving) return { n: 'Saving', hint: 'Writing crow.yaml', tone: 'busy' }
  if (build?.state === 'building')
    return { n: 'Building Nest', hint: `${build.elapsed_s} s${build.queued ? ', another build queued' : ''}`, tone: 'busy' }
  if (build?.state === 'failed')
    return { n: 'Build failed', hint: [build.brick, build.error].filter(Boolean).join(': ') || 'See the log below', tone: 'failed' }
  if (saved && needsBuild(saved) && build?.state === 'ok') return { n: 'New build ready', hint: 'Reload to use it', tone: 'ok' }
  if (saved?.length) return { n: 'Saved', hint: saved.includes('theme') ? 'Theme applied' : 'crow.yaml updated', tone: 'ok' }
  return null
}

function SaveBar({ v, ctl }: { v: SettingsView; ctl: SettingsControl }) {
  const t = barText(v)
  if (!t) return null
  const pending = t.tone === 'pending' || (t.tone === 'failed' && !!v.saveError)
  const log = !pending && v.build?.state === 'failed' ? v.build.log : []
  return (
    <div className={`set-bar ${t.tone}`} role={t.tone === 'failed' ? 'alert' : 'status'}>
      <div className="set-bar-row">
        <div className="t"><b>{t.n}</b><span>{t.hint}</span></div>
        {pending && <>
          <button type="button" className="btn" onClick={ctl.discard} disabled={v.saving}>Discard</button>
          <button type="button" className="btn p" onClick={() => void ctl.save()} disabled={v.saving}>{v.saving ? 'Saving…' : 'Save'}</button>
        </>}
        {t.tone === 'ok' && t.n === 'New build ready' && (
          <button type="button" className="btn p" onClick={() => location.reload()}>Reload</button>
        )}
      </div>
      {log.length > 0 && <pre className="set-log">{log.slice(-12).join('\n')}</pre>}
    </div>
  )
}

function SettingRow({ label, hint, cls = '', children }: { label: string; hint: string; cls?: string; children: ReactNode }) {
  return (
    <div className={`set-r ${cls}`}>
      <div className="lab"><b>{label}</b><span>{hint}</span></div>
      {children}
    </div>
  )
}

function Look({ v, ctl }: { v: SettingsView; ctl: SettingsControl }) {
  const known = v.settings!.known
  const d = v.draft!
  return (
    <>
      <div className="set-sl">Theme</div>
      <div className="set-themes">
        {known.themes.map((t) => (
          <button key={t} type="button" className="set-tcard" aria-pressed={d.theme === t} onClick={() => ctl.edit((x) => ({ ...x, theme: t }))}>
            <span className="n"><b>{title(t)}</b>{d.theme === t && <span className="ck" aria-hidden="true">✓</span>}</span>
            <span className="set-mini" data-theme={t} aria-hidden="true"><i /><span><i /><i /><i /></span></span>
            <span className="h">{THEME_HINTS[t] ?? 'Theme from a brick'}</span>
          </button>
        ))}
      </div>
      <div className="set-sl">Hero</div>
      <div className="set-card">
        <SettingRow label="Hero" hint="What sits at the top of Home" cls={d.hero !== v.base!.hero ? 'chg' : ''}>
          <span className="seg set-seg" role="group" aria-label="Hero">
            {known.heroes.map((h) => (
              <button key={h} type="button" className={d.hero === h ? 'on' : ''} aria-pressed={d.hero === h} onClick={() => ctl.edit((x) => ({ ...x, hero: h }))}>
                {title(h)}
              </button>
            ))}
          </span>
        </SettingRow>
      </div>
    </>
  )
}

function Pages({ v, ctl }: { v: SettingsView; ctl: SettingsControl }) {
  const d = v.draft!
  const list = useRef<HTMLDivElement>(null)
  const drag = useRef<string | null>(null)
  // Pointer drag on the handle (mouse and touch): the row under the pointer takes the page's place.
  const down = (key: string) => (e: PointerEvent<HTMLButtonElement>) => {
    drag.current = key
    e.currentTarget.setPointerCapture(e.pointerId)
  }
  const move = (e: PointerEvent<HTMLButtonElement>) => {
    const key = drag.current
    if (!key) return
    const rows = [...(list.current?.querySelectorAll<HTMLElement>('[data-page]') ?? [])]
    const at = rows.findIndex((r) => { const b = r.getBoundingClientRect(); return e.clientY >= b.top && e.clientY < b.bottom })
    if (at >= 0 && rows[at].dataset.page !== key) ctl.edit((x) => movePage(x, key, at))
  }
  const up = () => { drag.current = null }
  const keys = (key: string, i: number) => (e: KeyboardEvent<HTMLButtonElement>) => {
    const to = e.key === 'ArrowUp' ? i - 1 : e.key === 'ArrowDown' ? i + 1 : null
    if (to === null) return
    e.preventDefault()
    ctl.edit((x) => movePage(x, key, to))
  }
  return (
    <>
      <div className="set-sl">Order and visibility</div>
      <div className="set-card" ref={list}>
        {d.pages.map((p, i) => {
          const mark = pageMark(v.base!, d, p.key)
          const home = p.key === HOME
          const label = pageLabel(p.key)
          return (
            <div key={p.key} className={`set-r${mark ? ' chg' : ''}${p.on ? '' : ' off'}`} data-page={p.key}>
              {home ? <span className="set-hdl" aria-hidden="true" /> : (
                <button type="button" className="set-hdl" aria-label={`Move ${label}, arrow keys or drag`}
                  onPointerDown={down(p.key)} onPointerMove={move} onPointerUp={up} onPointerCancel={up} onKeyDown={keys(p.key, i)}>⠿</button>
              )}
              <span className="set-num num">{i + 1}</span>
              <div className="lab"><b>{label}</b><span>{PAGE_HINTS[p.key] ?? 'Page from a brick'}</span></div>
              {mark && <span className="set-mark">{mark}</span>}
              <button type="button" role="switch" className="set-sw" aria-checked={p.on} aria-label={`Show ${label}`}
                disabled={home} title={home ? 'Home stays on' : undefined} onClick={() => ctl.edit((x) => togglePage(x, p.key))} />
            </div>
          )
        })}
      </div>
      <p className="set-note">Drag to reorder. Home stays on.</p>
    </>
  )
}

// Panes that land with later items: what they will hold, no controls yet.
const LATER: Record<'modules' | 'general', Array<[string, string]>> = {
  modules: [
    ['Brick switches', 'Turn a brick on or off from here'],
    ['What each brick adds', 'Its pages, sections, hero and theme'],
  ],
  general: [
    ['Language', 'Words and dates of the Nest'],
    ['Time zone', 'The clock and agenda times'],
    ['Refresh', 'How often the Nest reads core'],
  ],
}

function Later({ pane }: { pane: 'modules' | 'general' }) {
  return (
    <>
      <div className="set-sl">Lands here later</div>
      <div className="set-card">
        {LATER[pane].map(([l, h]) => <SettingRow key={l} label={l} hint={h} cls="later"><span className="tag">later</span></SettingRow>)}
      </div>
      <p className="set-note">{pane === 'modules' ? 'crow.yaml modules stays the place to switch bricks until then.' : 'crow.yaml stays the place for these until then.'}</p>
    </>
  )
}

function PaneBody({ pane, v, ctl }: { pane: PaneKey; v: SettingsView; ctl: SettingsControl }) {
  if (pane === 'modules' || pane === 'general') return <Later pane={pane} />
  if (v.status === 'off') return <p className="empty">Settings need the core API. Fixtures and the static build have none.</p>
  if (v.status === 'loading') return <p className="empty">Loading crow.yaml…</p>
  if (v.status === 'error') return <p className="empty c-f" role="alert">Could not read the settings: {v.error}</p>
  return pane === 'look' ? <Look v={v} ctl={ctl} /> : <Pages v={v} ctl={ctl} />
}

function Errors({ v }: { v: SettingsView }) {
  const lines = errorLines(v.settings?.errors ?? [])
  if (!lines.length) return null
  return (
    <div className="well set-errs" role="alert">
      <b>crow.yaml has values Nest cannot use</b>
      {lines.map((l) => <span key={l}>{l}</span>)}
    </div>
  )
}

function Head({ pane, icon }: { pane: PaneKey; icon: boolean }) {
  const p = PANES.find((x) => x.key === pane)!
  return (
    <div className="set-ph">
      {icon && <span className="ic"><TabIcon>{PANE_ICONS[pane]}</TabIcon></span>}
      <div><h2>{p.label}</h2><p>{p.hint}</p></div>
    </div>
  )
}

const href = (p: PaneKey) => `/settings?pane=${p}`

export function SettingsPage({ v, ctl, pane, laptop }: { v: SettingsView; ctl: SettingsControl; pane: PaneKey | null; laptop: boolean }) {
  const dirty = (p: PaneKey) => !!v.base && !!v.draft && paneDirty(v.base, v.draft, p)
  if (laptop) {
    const cur = pane ?? 'look'
    return (
      <div className="screen set-stage">
        <div className="set-win">
          <nav className="set-side" aria-label="Settings panes">
            {PANES.map((p) => (
              <Link key={p.key} href={href(p.key)} className="set-nv" aria-current={p.key === cur ? 'page' : undefined}>
                <span className="ic"><TabIcon>{PANE_ICONS[p.key]}</TabIcon></span>{p.label}
                {dirty(p.key) && <span className="set-dot" aria-label="unsaved changes" />}
              </Link>
            ))}
          </nav>
          <section className="set-pane" aria-label={`${PANES.find((p) => p.key === cur)!.label} settings`}>
            <Head pane={cur} icon />
            <div className="set-body"><Errors v={v} /><PaneBody pane={cur} v={v} ctl={ctl} /></div>
            <SaveBar v={v} ctl={ctl} />
          </section>
        </div>
      </div>
    )
  }
  if (!pane)
    return (
      <div className="screen set-phone">
        <div className="set-ph"><div><h2>Settings</h2><p>Changes apply after Save</p></div></div>
        <Errors v={v} />
        <nav className="set-card set-list" aria-label="Settings panes">
          {PANES.map((p) => (
            <Link key={p.key} href={href(p.key)} className="set-r">
              <span className="ic"><TabIcon>{PANE_ICONS[p.key]}</TabIcon></span>
              <div className="lab"><b>{p.label}</b><span>{p.hint}</span></div>
              {dirty(p.key) && <span className="set-dot" aria-label="unsaved changes" />}
              <span className="chev" aria-hidden="true">›</span>
            </Link>
          ))}
        </nav>
        <SaveBar v={v} ctl={ctl} />
      </div>
    )
  return (
    <div className="screen set-phone">
      <Link href="/settings" className="set-back">‹ Settings</Link>
      <Head pane={pane} icon={false} />
      <Errors v={v} />
      <PaneBody pane={pane} v={v} ctl={ctl} />
      <SaveBar v={v} ctl={ctl} />
    </div>
  )
}

export function Settings({ laptop, refresh }: { laptop: boolean; refresh: () => Promise<void> }) {
  const [v, ctl] = useSettings(useSource().settings)
  const pane = paneOf(useQuery().get('pane'), laptop)
  // Theme preview: the draft theme paints the page until Save keeps it or Discard drops it.
  const draftTheme = v.draft?.theme
  const baseTheme = v.base?.theme
  useEffect(() => {
    if (draftTheme && draftTheme !== baseTheme) paintTheme(draftTheme)
    else unpaintTheme()
  }, [draftTheme, baseTheme])
  useEffect(() => unpaintTheme, [])
  // A saved theme becomes the kept one, the header switcher included.
  const saved = v.saved
  useEffect(() => {
    if (saved?.includes('theme') && baseTheme) setTheme(baseTheme)
  }, [saved, baseTheme])
  // A passing build: the feed sees the new build and the header shows its reload chip.
  const built = !!saved && needsBuild(saved as SettingKey[]) && v.build?.state === 'ok'
  useEffect(() => { if (built) void refresh() }, [built, refresh])
  return <SettingsPage v={v} ctl={ctl} pane={pane} laptop={laptop} />
}
