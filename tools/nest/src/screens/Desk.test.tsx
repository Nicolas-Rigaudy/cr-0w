import { renderToStaticMarkup } from 'react-dom/server'
import { afterEach, describe, expect, it, vi } from 'vitest'
import { feed, type Snapshot } from '../data/feed'
import { fixtureSource } from '../data/fixtures'
import { DataContext } from '../data/useData'
import { Desk } from './Desk'

const snapshot = async (src: ReturnType<typeof fixtureSource>): Promise<Snapshot> => {
  const box: { s?: Snapshot } = {}
  await feed(src, (x) => { box.s = x })()
  return box.s!
}
const at = (w: number, search = '') => {
  vi.stubGlobal('window', {
    matchMedia: (q: string) => ({ matches: Number(/min-width:\s*(\d+)/.exec(q)?.[1] ?? 0) <= w, addEventListener() {}, removeEventListener() {} }),
  })
  vi.stubGlobal('location', { search, pathname: '/desk' })
}
const refresh = async () => {}

afterEach(() => { vi.unstubAllGlobals() })

describe('Nest-started sessions (item_188)', () => {
  it('a start replies with its tmux session, which shows on the new agent row', async () => {
    const src = fixtureSource()
    const r = await src.actions.start({ project: 'demo', prompt: 'look at the build' })
    expect(r.tmux).toMatch(/^crow-demo-/)
    expect(r.message).toContain('running in Nest')
    const s = await snapshot(src)
    const a = (s.sections.dispatch?.data as { agents: Array<{ pid: number; tmux?: string }> }).agents.find((x) => x.tmux === r.tmux)
    expect(a).toBeDefined()
  })

  it('shows Terminal on every session, Log on background runs', async () => {
    const src = fixtureSource()
    const r = await src.actions.start({ project: 'demo' })
    const s = await snapshot(src)
    const agents = (s.sections.dispatch?.data as { agents: Array<{ pid: number; tmux?: string; kind: string }> }).agents
    const pid = agents.find((x) => x.tmux === r.tmux)!.pid
    const render = (q: string) => {
      at(1300, q)
      return renderToStaticMarkup(<DataContext.Provider value={src}><Desk snap={s} refresh={refresh} /></DataContext.Provider>)
    }
    expect(render(`?pid=${pid}&tab=transcript`)).toContain('>Terminal</button>')
    const kitty = agents.find((x) => x.kind === 'session' && !x.tmux)!.pid
    expect(render(`?pid=${kitty}&tab=transcript`)).toContain('>Terminal</button>')
    const run = agents.find((x) => x.kind === 'run')!.pid
    expect(render(`?pid=${run}&tab=transcript`)).not.toContain('>Terminal</button>')
    expect(render(`?pid=${run}&tab=transcript`)).toContain('>Log</button>')
  })
})

describe('Terminal tab and setup sessions (item_240)', () => {
  it('reads, types into and sizes a session; refuses a background run', async () => {
    const src = fixtureSource()
    const r = await src.actions.start({ project: 'demo' })
    const s = await snapshot(src)
    const agents = (s.sections.dispatch?.data as { agents: Array<{ pid: number; tmux?: string; kind: string }> }).agents
    const pid = agents.find((x) => x.tmux === r.tmux)!.pid
    const scr = await src.actions.screen(pid)
    expect(scr.lines.join('\n')).toContain('demo')
    expect(scr.cursor.shown).toBe(true)
    await expect(src.actions.keys(pid, [{ text: 'hi' }, { key: 'Enter' }])).resolves.toBeTruthy()
    await expect(src.actions.resize(pid, 100, 30)).resolves.toBeTruthy()
    const kitty = agents.find((x) => x.kind === 'session' && !x.tmux)!.pid
    await expect(src.actions.screen(kitty)).resolves.toBeTruthy()
    const run = agents.find((x) => x.kind === 'run')!.pid
    await expect(src.actions.screen(run)).rejects.toMatchObject({ status: 409 })
  })

  it('shows a setup session as waiting on setup, without the attach button', async () => {
    const src = fixtureSource()
    const r = await src.actions.start({ project: 'demo' })
    const s = await snapshot(src)
    const d = s.sections.dispatch!.data as { agents: Array<{ pid: number; tmux?: string; state?: string }> }
    const a = d.agents.find((x) => x.tmux === r.tmux)!
    a.state = 'setup'
    at(1300, `?pid=${a.pid}`)
    const html = renderToStaticMarkup(<DataContext.Provider value={src}><Desk snap={s} refresh={refresh} /></DataContext.Provider>)
    expect(html).toContain('waits on setup')
    expect(html).not.toContain('Attach a terminal')
  })

  it('a resume replies like a start when it runs headless', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const id = (s.sections.sessions?.data as { sessions: Array<{ conversation_id: string }> }).sessions[0].conversation_id
    const r = await src.actions.resume(id)
    expect(r.action).toBe('resume')
  })
})

describe('Desk layout C (item_295)', () => {
  const html = async (w = 1300, q = '') => {
    const src = fixtureSource()
    const s = await snapshot(src)
    at(w, q)
    return renderToStaticMarkup(<DataContext.Provider value={src}><Desk snap={s} refresh={refresh} /></DataContext.Provider>)
  }

  it('shows four panes with their counts in the titles and no count tiles', async () => {
    const h = await html()
    for (const t of ['[1] Needs you', '[2] Moving', '[3] Queued', '[4] Done · failed', '[5] ']) expect(h).toContain(t)
    expect(h).toMatch(/\[1\] Needs you<\/span><span class="num">\d+</)
    expect(h).toMatch(/\d+ · folded ▸/)
    expect(h).not.toContain('done today')
    expect(h).not.toContain('ok rate 7d')
    expect(h).toContain('◇ Butler')
    expect(h).not.toContain('aria-label="Ask CR-0W"')
  })

  it('details the picked unit in four tabs with its actions', async () => {
    const h = await html()
    for (const t of ['Overview', 'Transcript', 'Timeline', 'Files']) expect(h).toMatch(new RegExp(`role="tab"[^>]*>${t}</button>`))
    expect(h).toMatch(/aria-selected="true"[^>]*>Overview/)
    expect(h).toContain('Stop')
  })

  it('offers Now and Background, with AFK batch disabled until the core route', async () => {
    const h = await html()
    expect(h).toMatch(/aria-pressed="true"[^>]*>Now</)
    expect(h).toMatch(/>Background</)
    expect(h).toMatch(/disabled=""[^>]*>AFK batch</)
    expect(h).toContain('item_297')
  })

  it('shows one pane at a time on the phone, behind pane tabs', async () => {
    const h = await html(400)
    expect(h).toContain('aria-label="Panes"')
    expect(h).toMatch(/\[1\] Needs <span class="num">\d+/)
    expect(h).not.toContain('[2] Moving</span>')
  })

  it('opens the detail as its own page on the phone, with a back button', async () => {
    const src = fixtureSource()
    const s = await snapshot(src)
    const pid = (s.sections.dispatch!.data as { agents: Array<{ pid: number }> }).agents[0].pid
    at(400, `?pid=${pid}`)
    const h = renderToStaticMarkup(<DataContext.Provider value={src}><Desk snap={s} refresh={refresh} /></DataContext.Provider>)
    expect(h).toContain('class="mback"')
    expect(h).toContain('Overview')
  })
})
