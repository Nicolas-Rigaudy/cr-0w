import { Fragment, useEffect, useLayoutEffect, useRef, useState, type CSSProperties } from 'react'
import type { UiState } from '../data/state'
import type { AgendaEvent, InboxItem, Mission } from '../data/types'
import { heroModel } from './model'
import type { Heads, Palette, SceneOut } from './scene'
import { drawScene } from './scene'

// Holotable hero (item_162): ASCII galaxy map of projects with the option B head above the lens.
// One frame per data, minute, size or theme change. No animation loop yet (item_158).

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

export function Hero({ missions, inbox, agenda, now }: { missions: Mission[]; inbox: InboxItem[]; agenda: AgendaEvent[]; now: Date }) {
  const m = heroModel(missions, inbox, now)
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
    // Same head turned 20 degrees, to look toward the planet that waits longest.
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

  // Redraw per minute at most: wait and run clocks move the stalks and the now notch.
  const minute = Math.floor(now.getTime() / 60000)
  const key = JSON.stringify([m.mood, minute, agenda.map((e) => [e.id, e.start]), m.projects.map((p) => [p.name, p.state, p.open, p.wait, p.runFrom, p.cause])])
  useLayoutEffect(() => {
    if (!cv.current || !box.current || !size || size.w < 50) return
    document.fonts.ready.then(() => {
      if (!cv.current || !box.current) return
      setOut(drawScene(cv.current, size.w, size.h, window.devicePixelRatio || 1, heroModel(missions, inbox, now), heads, readPalette(box.current), turn, { now, agenda }))
    })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key, size, theme, heads, turn])

  return (
    <div ref={box} className="hero-box" role="img" aria-label={`CR-0W: ${m.line}`} style={out ? ({ '--lens-y': `${out.lens.y}px` } as CSSProperties) : undefined}>
      <canvas ref={cv} className="hero-cv" />
      <div className="hero-say"><b>CR-0W ›</b> {m.line}<i className="hero-cur" /></div>
      {out?.planets.map((p) => (
        <Fragment key={p.name}>
          {p.tag && <div className={`hero-tag${p.target ? ' on' : ''}`} style={{ left: p.tag.x, top: p.tag.y }}>{p.tag.text}</div>}
          <div className={`hero-pl ${p.side}${p.target ? ' on' : ''}`} style={{ left: p.x, top: p.y }}>
            {p.name.toUpperCase()}
            {p.note && <span style={{ color: `var(${VARS[p.state]})` }}> · {p.note}</span>}
          </div>
        </Fragment>
      ))}
      {out?.rim.map((r) => (
        <div key={r.text} className={`hero-rim ${r.kind}${r.side ? ` ${r.side}` : ''}`} style={{ left: r.x, top: r.y }}>{r.text}</div>
      ))}
      {m.asks > 0 && out && (
        <div className="hero-chip" style={{ left: Math.min(out.head.right - 40, (size?.w ?? 0) - 60), top: out.head.top + 40 }}>
          {m.asks} {m.asks === 1 ? 'ASK' : 'ASKS'}
        </div>
      )}
    </div>
  )
}
