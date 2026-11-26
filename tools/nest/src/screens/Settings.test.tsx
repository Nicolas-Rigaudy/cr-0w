import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSettings, fixtureSource } from '../data/fixtures'
import { movePage, togglePage, type NestKnown } from '../data/settings'
import { settingsControl, type SettingsControl, type SettingsView } from '../data/useSettings'
import { match } from '../router'
import { Shell } from '../ui/Shell'
import { barTab, screenTab, shownRoute } from '../ui/tabs'
import { barText, paneOf, SettingsPage, type PaneKey } from './Settings'

afterEach(() => { vi.unstubAllGlobals() })

const known: NestKnown = { themes: ['modern', 'cyber'], heroes: ['holotable', 'none'], pages: ['home', 'desk', 'ops', 'inbox', 'mrs', 'agenda', 'reports'] }

async function ready(): Promise<[() => SettingsView, SettingsControl]> {
  let v: SettingsView = { status: 'loading', saving: false }
  const ctl = settingsControl(fixtureSettings(new Date('2026-11-25T10:00:00Z'), known), (n) => (v = n))
  await ctl.load()
  return [() => v, ctl]
}

const render = (v: SettingsView, ctl: SettingsControl, pane: PaneKey | null, laptop: boolean) => {
  vi.stubGlobal('location', { search: '', pathname: '/settings' })
  return renderToStaticMarkup(<SettingsPage v={v} ctl={ctl} pane={pane} laptop={laptop} />)
}

describe('settings route (item_284)', () => {
  it('always shows, under More on the phone bar, even when crow.yaml lists few pages', () => {
    const r = match('/settings')
    expect(r).toEqual({ tab: 'settings' })
    expect(shownRoute(r, ['home', 'desk'])).toEqual(r)
    expect(barTab(screenTab(r))).toBe('more')
  })

  it('opens Look on laptop and the pane list on phone', () => {
    expect(paneOf(null, true)).toBe('look')
    expect(paneOf(null, false)).toBeNull()
    expect(paneOf('pages', false)).toBe('pages')
    expect(paneOf('nope', true)).toBe('look')
  })
})

describe('settings entries', () => {
  const snapshot = async (): Promise<Snapshot> => {
    const box: { s?: Snapshot } = {}
    await feed(fixtureSource(), (x) => { box.s = x })()
    return box.s!
  }
  it('puts Settings in the header and the More sheet button always shows', async () => {
    vi.stubGlobal('location', { search: '', pathname: '/settings' })
    vi.stubGlobal('document', { documentElement: { dataset: {} } })
    const html = renderToStaticMarkup(<Shell active="settings" snap={await snapshot()} scope="all" setScope={() => {}}>x</Shell>)
    expect(html).toMatch(/<a href="\/settings" class="icon-btn set-ic" aria-label="Settings" title="Settings" aria-current="page">/)
    expect(html).toMatch(/<button type="button" class="tab" aria-expanded="false"[^>]*aria-current="page">.*More<span class="sr-only">, Settings open<\/span><\/button>/)
  })
})

describe('Settings page', () => {
  it('shows the laptop window with panes, theme cards with their own theme and the hero choice', async () => {
    const [v, ctl] = await ready()
    const html = render(v(), ctl, 'look', true)
    expect(html).toContain('class="set-win"')
    expect([...html.matchAll(/class="set-nv"[^>]*>.*?<\/svg><\/span>(\w+)/g)].map((m) => m[1])).toEqual(['Look', 'Pages', 'Modules', 'General'])
    expect(html).toContain('<span class="set-mini" data-theme="cyber" aria-hidden="true">')
    expect(html).toContain('aria-pressed="true"><span class="n"><b>Modern</b>')
    expect(html).toContain('Neon green and pink on black')
    expect(html).toMatch(/<button type="button" class="on" aria-pressed="true">Holotable<\/button>/)
    expect(html).not.toContain('set-bar')
  })

  it('marks pending pages, locks Home and offers Discard and Save with the rebuild hint', async () => {
    const [v, ctl] = await ready()
    ctl.edit((d) => togglePage(movePage(d, 'agenda', 3), 'reports'))
    const html = render(v(), ctl, 'pages', true)
    expect(html).toMatch(/data-page="agenda">.*?<span class="set-mark">moved up<\/span>/)
    expect(html).toMatch(/data-page="reports">.*?<span class="set-mark">hidden<\/span>/)
    expect(html).toMatch(/role="switch" class="set-sw" aria-checked="true" aria-label="Show Home" disabled="" title="Home stays on"/)
    expect(html).toContain('<b>2 changes</b><span>Pages need a rebuild, about 40 s</span>')
    expect(html).toContain('>Discard</button>')
    expect(html).toContain('>Save</button>')
    // The Pages nav item carries the pending dot, Look does not.
    expect(html).toMatch(/aria-current="page">.*?Pages<span class="set-dot"/)
    expect(html).not.toMatch(/Look<span class="set-dot"/)
  })

  it('lists panes on phone with a pending dot, then one pane with a back link', async () => {
    const [v, ctl] = await ready()
    ctl.edit((d) => ({ ...d, theme: 'cyber' }))
    const list = render(v(), ctl, null, false)
    expect(list).toContain('<h2>Settings</h2><p>Changes apply after Save</p>')
    expect(list).toMatch(/href="\/settings\?pane=look" class="set-r">.*?<b>Look<\/b>.*?<span class="set-dot"/)
    expect(list).toContain('<b>1 change</b><span>The theme applies at once</span>')
    const pane = render(v(), ctl, 'modules', false)
    expect(pane).toContain('<a href="/settings" class="set-back">‹ Settings</a>')
    expect(pane).toContain('Lands here later')
    // No dead controls in a later pane: the only buttons are the Save bar's.
    expect(pane.split('class="set-bar')[0]).not.toMatch(/<button/)
  })
})

describe('save bar', () => {
  it('follows the save into the build, then offers Reload', async () => {
    const [v, ctl] = await ready()
    ctl.edit((d) => ({ ...d, hero: 'none' }))
    expect(barText(v())?.hint).toBe('Hero needs a rebuild, about 40 s')
    await ctl.save()
    expect(barText(v())).toMatchObject({ n: 'Building Nest', tone: 'busy' })
    await ctl.tick()
    await ctl.tick()
    expect(barText(v())).toMatchObject({ n: 'New build ready', tone: 'ok' })
    expect(render(v(), ctl, 'look', true)).toContain('>Reload</button>')
  })

  it('shows a failed build with its brick, error and log', () => {
    const v: SettingsView = {
      status: 'ready', saving: false,
      base: { theme: 'modern', hero: 'holotable', pages: [] }, draft: { theme: 'modern', hero: 'holotable', pages: [] },
      build: { state: 'failed', elapsed_s: 9, queued: false, brick: 'forge', error: 'tsc failed', log: ['src/x.ts: error'], served: 'a' },
    }
    expect(barText(v)).toEqual({ n: 'Build failed', hint: 'forge: tsc failed', tone: 'failed' })
    const html = renderToStaticMarkup(<SettingsPage v={{ ...v, settings: { file: '', theme: null, hero: null, pages: null, known, errors: [] } }} ctl={{} as SettingsControl} pane="look" laptop />)
    expect(html).toContain('<pre class="set-log">src/x.ts: error</pre>')
  })

  it('keeps the draft and shows the refusal of a save', async () => {
    const [v, ctl] = await ready()
    ctl.edit((d) => ({ ...d, theme: 'neon' }))
    await ctl.save()
    expect(barText(v())).toMatchObject({ n: '1 change', tone: 'failed' })
    expect(barText(v())?.hint).toContain('unknown theme')
  })
})
