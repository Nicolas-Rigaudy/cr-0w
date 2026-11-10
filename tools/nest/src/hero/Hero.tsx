import { useEffect, useLayoutEffect, useRef, useState } from 'react'
import type { UiState } from '../data/state'
import type { InboxItem, Mission } from '../data/types'
import { elapsed } from '../format'
import { heroModel } from './model'
import type { Heads, Palette, SceneOut } from './scene'
import { drawScene } from './scene'

// Static hero (item_155): pass 7 ASCII scene with the option B head. One frame per data,
// size or theme change. No animation loop.

type RGB = [number, number, number]
const VARS: Record<UiState, string> = { needs: '--c-n', working: '--c-w', failed: '--c-f', sleeping: '--c-s', queued: '--c-q', done: '--c-d' }

function readPalette(el: Element): Palette {
  const cs = getComputedStyle(el)
  const ctx = document.createElement('canvas').getContext('2d')
  const rgb = (name: string, fb: RGB): RGB => {
    const v = cs.getPropertyValue(name).trim()
    if (!v || !ctx) return fb
    ctx.fillStyle = '#000'; ctx.fillStyle = v
    const s = ctx.fillStyle
    if (s.startsWith('#')) { const n = parseInt(s.slice(1), 16); return [(n >> 16) & 255, (n >> 8) & 255, n & 255] }
    const m = s.match(/[\d.]+/g)
    return m ? [+m[0], +m[1], +m[2]] : fb
  }
  const state = {} as Record<UiState, RGB>
  const FB: Record<UiState, RGB> = { needs: [255, 184, 0], working: [46, 242, 255], failed: [255, 46, 99], sleeping: [111, 134, 184], queued: [138, 160, 200], done: [95, 120, 168] }
  for (const k of Object.keys(VARS) as UiState[]) state[k] = rgb(VARS[k], FB[k])
  return { ring: rgb('--c-w', [79, 195, 255]), dim: rgb('--dim', [95, 120, 168]), fg: rgb('--fg', [216, 230, 255]), accent: rgb('--accent', [192, 92, 255]), state }
}

export function Hero({ missions, inbox, now }: { missions: Mission[]; inbox: InboxItem[]; now: Date }) {
  const m = heroModel(missions, inbox)
  const box = useRef<HTMLDivElement>(null)
  const cv = useRef<HTMLCanvasElement>(null)
  const [heads, setHeads] = useState<Heads>()
  const [size, setSize] = useState<{ w: number; h: number }>()
  const [theme, setTheme] = useState(() => document.documentElement.dataset.theme)
  const [out, setOut] = useState<SceneOut>()
  const [turn, setTurn] = useState<{ l: Heads; r: Heads }>()

  useEffect(() => {
    let live = true
    import('./heads.json').then((d) => live && setHeads(d.default as unknown as Heads))
    // Same head turned 20 degrees, to look toward the most urgent project.
    Promise.all([import('./heads-l.json'), import('./heads-r.json')]).then(([l, r]) => live && setTurn({ l: l.default as unknown as Heads, r: r.default as unknown as Heads }))
    return () => { live = false }
  }, [])

  useEffect(() => {
    const el = box.current
    if (!el) return
    const ro = new ResizeObserver(([e]) => setSize({ w: Math.round(e.contentRect.width), h: Math.round(e.contentRect.height) }))
    ro.observe(el)
    const mo = new MutationObserver(() => setTheme(document.documentElement.dataset.theme))
    mo.observe(document.documentElement, { attributes: true, attributeFilter: ['data-theme'] })
    return () => { ro.disconnect(); mo.disconnect() }
  }, [])

  const key = JSON.stringify([m.mood, m.projects.map((p) => [p.name, p.state, p.missions.map((x) => [x.id, x.state])])])
  useLayoutEffect(() => {
    if (!cv.current || !box.current || !size || size.w < 50) return
    document.fonts.ready.then(() => {
      if (!cv.current || !box.current) return
      setOut(drawScene(cv.current, size.w, size.h, window.devicePixelRatio || 1, heroModel(missions, inbox), heads, readPalette(box.current), turn))
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key, size, theme, heads, turn])

  const started = new Map(m.projects.flatMap((p) => p.missions.map((x) => [x.title, x.started] as const)))

  return (
    <div ref={box} className="hero-box" role="img" aria-label={`CR-0W: ${m.line}`}>
      <canvas ref={cv} className="hero-cv" />
      <div className="hero-say"><b>CR-0W ›</b> {m.line}<i className="hero-cur" /></div>
      {out?.projects.map((p) => (
        <div key={p.text} className={`hero-lbl${p.live ? ' live' : ''}`} style={{ left: p.x, top: p.y }}>
          {p.text.toUpperCase()} <span style={{ color: `var(${VARS[p.state]})` }}>· {p.state.toUpperCase()}</span>
        </div>
      ))}
      {out?.missions.map((x, i) => {
        const st = started.get(x.text)
        return (
          <div key={i} className={`hero-ml${x.state === 'needs' ? ` up${x.left && !x.right ? ' l' : ''}` : x.left ? ' l' : ''}${x.right ? ' r' : ''}`} style={{ left: x.x, top: x.y, color: `var(${VARS[x.state]})` }}>
            {x.text}{x.state === 'working' && st ? ` · ${elapsed(st, now)}` : x.state === 'failed' ? ' · failed' : ''}
          </div>
        )
      })}
      {m.asks > 0 && out && (
        <div className="hero-chip" style={{ left: out.head.x + 230, top: out.head.y - 260 }}>
          {m.asks} {m.asks === 1 ? 'ASK' : 'ASKS'}
        </div>
      )}
    </div>
  )
}
