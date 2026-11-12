import { Fragment, useEffect, useLayoutEffect, useMemo, useRef, useState, type CSSProperties, type KeyboardEvent, type PointerEvent } from 'react'
import type { UiState } from '../data/state'
import { lastRun, openInbox } from '../data/state'
import type { AgendaEvent, InboxItem, Mission } from '../data/types'
import { ago, title } from '../format'
import { navigate } from '../router'
import { Glyph } from '../ui/Glyph'
import { useMedia } from '../ui/useMedia'
import { heroModel } from './model'
import { declutterDom, type Box } from './declutter'
import { DEG_PER_WIDTH, detent, liftCurve, spinToFront, step, stepTo, type Spin } from './motion'
import type { Heads, Palette, PlanetLabel, SceneOut } from './scene'
import { drawScene } from './scene'

// Holotable hero (item_162, item_158): ASCII galaxy map of projects with the option B head above
// the lens. Drag spins the planets, hover lifts one and shows its moons, click focuses then opens.
// The frame loop only runs while something moves; otherwise it sleeps (0 fps).

type RGB = [number, number, number]
const VARS: Record<UiState, string> = { needs: '--c-n', working: '--c-w', failed: '--c-f', sleeping: '--c-s', queued: '--c-q', done: '--c-d' }
const PULSE_MS = 2400
const PULSE_FRAME_MS = 100
const LIFT_MS = 420
const GLITCH_MS = 180
const HOVER_MS = 200
const DRAG_PX = 6

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

const hit = (out: SceneOut | undefined, x: number, y: number): PlanetLabel | undefined => {
  if (!out) return undefined
  let best: PlanetLabel | undefined, d = Infinity
  for (const p of out.planets) {
    const dd = Math.hypot(x - p.cx, y - p.cy)
    if (dd < p.rad + 8 && dd < d) { d = dd; best = p }
  }
  return best
}
const onHead = (out: SceneOut | undefined, x: number, y: number) =>
  !!out && x > out.head.left + 30 && x < out.head.right - 30 && y > out.head.top && y < out.lens.y + 10

export const missionsHref = (project: string) => `/missions?project=${encodeURIComponent(project)}`

// Detail card for the focused planet: open missions, what waits on you, last run.
// Small heroes (phone, inner panel) get a compact card: one ask, two missions.
function Card({ project, missions, inbox, now, at, w }: { project: string; missions: Mission[]; inbox: InboxItem[]; now: Date; at: PlanetLabel; w: number }) {
  const nA = w < 760 ? 1 : 2, nM = w < 760 ? 2 : 4
  const ms = missions.filter((m) => m.project === project && m.state !== 'done' && m.state !== 'cancelled')
  const asks = openInbox(inbox).filter((r) => r.project === project && r.type !== 'fyi')
  const runs = missions.filter((m) => m.project === project).flatMap((m) => m.runs.map((r) => ({ r, m }))).sort((a, b) => b.r.started.localeCompare(a.r.started))
  const last = runs[0]
  const left = Math.max(8, Math.min(w - 248, at.cx - 120))
  return (
    <button type="button" className="hero-card" data-prio={1} data-up={2 * at.rad + 28} style={{ left, top: at.cy + at.rad + 14, '--cx': `${at.cx}px`, '--cy': `${at.cy}px`, '--r': `${at.rad}px` } as CSSProperties} onClick={() => navigate(missionsHref(project))} aria-label={`Open ${project} missions`}>
      <span className="hc-h"><Glyph state={at.state} /> {project.toUpperCase()}<span className="hc-go">open →</span></span>
      {asks.slice(0, nA).map((r) => <span key={r.id} className="hc-r c-n"><span className="hc-t">Your move · {r.ask}</span></span>)}
      {ms.slice(0, nM).map((m) => <span key={m.id} className="hc-r"><Glyph state={stateOf(m)} /> <span className="hc-t">{title(m)}</span></span>)}
      {(ms.length > nM || asks.length > nA) && <span className="hc-r dim">+{Math.max(0, ms.length - nM) + Math.max(0, asks.length - nA)} more</span>}
      {last && (
        <span className={`hc-r dim${at.state === 'failed' ? ' c-f' : ''}`}>
          Last run {last.r.ended ? ago(last.r.ended, now) : 'live'}{last.r.cause && last.r.cause !== 'done' ? ` · ${last.r.cause}${last.r.exit ? ` (exit ${last.r.exit})` : ''}` : ''}
        </span>
      )}
    </button>
  )
}
const stateOf = (m: Mission): UiState =>
  m.state === 'failed' ? 'failed' : m.state === 'waiting' || m.state === 'review' ? 'needs' : m.state === 'queued' ? (lastRun(m)?.cause === 'usage_limit' ? 'sleeping' : 'queued') : m.state === 'done' ? 'done' : 'working'

// say replaces the model's line (phone Home counts sessions and checks too).
export function Hero({ missions, inbox, agenda, now, say }: { missions: Mission[]; inbox: InboxItem[]; agenda: AgendaEvent[]; now: Date; say?: string }) {
  const m = useMemo(() => heroModel(missions, inbox, now), [missions, inbox, now])
  const reduced = useMedia('(prefers-reduced-motion: reduce)')
  const box = useRef<HTMLDivElement>(null)
  const cv = useRef<HTMLCanvasElement>(null)
  const [heads, setHeads] = useState<Heads>()
  const [size, setSize] = useState<{ w: number; h: number }>()
  const [theme, setTheme] = useState(() => document.documentElement.dataset.theme)
  const [out, setOut] = useState<SceneOut>()
  const [turn, setTurn] = useState<{ l: Heads; r: Heads }>()
  const [focus, setFocus] = useState<string>()

  // Mutable frame state: read by the loop without re-rendering.
  const st = useRef({
    spin: { spin: 0, vel: 0 } as Spin, drag: undefined as undefined | { id: number; x0: number; x: number; t: number; spin0: number; moved: boolean },
    hover: undefined as undefined | { name: string; k: number; dir: 1 | -1 }, look: undefined as number | undefined, lean: 0,
    lift: {} as Record<string, number>, glitch: {} as Record<string, number>, prev: undefined as Map<string, UiState> | undefined,
    raf: 0, timer: 0, last: 0, pal: undefined as Palette | undefined, out: undefined as SceneOut | undefined,
  })

  useEffect(() => {
    let live = true
    import('./heads.json').then((d) => live && setHeads(d.default as unknown as Heads))
    // Same head turned 20 degrees, to look toward a planet.
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
    return () => { ro.disconnect(); mo.disconnect(); cancelAnimationFrame(st.current.raf); clearTimeout(st.current.timer) }
  }, [])

  const live = m.projects.some((p) => p.state === 'working')
  const a0s = () => (st.current.out?.planets ?? []).map((p) => p.a0)

  // One frame. Returns true when another frame is needed now (not the slow pulse).
  const frame = (t: number) => {
    const s = st.current
    if (!cv.current || !box.current || !size || size.w < 50) return false
    const dt = s.last ? Math.min(t - s.last, 64) : 16
    s.last = t
    let busy = false
    if (!s.drag && (s.spin.vel !== 0 || s.spin.to !== undefined)) {
      if (reduced) { s.spin.vel = 0; if (s.spin.to !== undefined) { s.spin.spin = s.spin.to; s.spin.to = undefined } }
      else busy = step(s.spin, dt, a0s()) || busy
    }
    if (s.drag) busy = true
    if (s.hover) {
      const k = Math.max(0, Math.min(1, s.hover.k + (reduced ? s.hover.dir : (s.hover.dir * dt) / HOVER_MS)))
      s.hover.k = k
      if (k === 0 && s.hover.dir < 0) s.hover = undefined
      else if (k > 0 && k < 1) busy = true
    }
    for (const [map, ms] of [[s.lift, LIFT_MS], [s.glitch, GLITCH_MS]] as const) {
      for (const n of Object.keys(map)) {
        map[n] += dt / ms
        if (map[n] >= 1) delete map[n]
        else busy = true
      }
    }
    const lean = s.drag ? Math.max(-4, Math.min(4, s.spin.vel * 30)) : 0
    s.lean += (lean - s.lean) * Math.min(1, dt / 120)
    if (Math.abs(s.lean - lean) > 0.05) busy = true
    s.pal ??= readPalette(box.current)
    const lift: Record<string, number> = {}
    for (const n of Object.keys(s.lift)) lift[n] = liftCurve(s.lift[n])
    const glitch: Record<string, number> = {}
    for (const n of Object.keys(s.glitch)) glitch[n] = 1 - s.glitch[n]
    const o = drawScene(cv.current, size.w, size.h, window.devicePixelRatio || 1, m, heads, s.pal, turn, {
      now, agenda, spin: s.spin.spin, hover: s.hover && { name: s.hover.name, k: s.hover.k },
      look: s.look, lean: s.lean, pulse: live && !reduced ? (t % PULSE_MS) / PULSE_MS : undefined, lift, glitch,
    })
    s.out = o
    setOut(o)
    return busy
  }

  // Wake the loop. Fast frames while busy; a slow tick for the working pulse; else sleep.
  const kick = () => {
    const s = st.current
    clearTimeout(s.timer)
    if (s.raf) return
    const loop = (t: number) => {
      s.raf = 0
      if (frame(t)) { s.raf = requestAnimationFrame(loop); return }
      s.last = 0
      if (live && !reduced && document.visibilityState === 'visible') s.timer = window.setTimeout(() => kick(), PULSE_FRAME_MS)
    }
    s.raf = requestAnimationFrame(loop)
  }

  // Data changes: start a lift (entered needs you) or a glitch (entered failed), then redraw.
  const minute = Math.floor(now.getTime() / 60000)
  const key = JSON.stringify([m.mood, minute, agenda.map((e) => [e.id, e.start]), m.projects.map((p) => [p.name, p.state, p.open, p.wait, p.runFrom, p.cause])])
  useLayoutEffect(() => {
    const s = st.current
    const cur = new Map(m.projects.map((p) => [p.name, p.state]))
    if (s.prev && !reduced) for (const [n, state] of cur) {
      if (s.prev.get(n) === state) continue
      if (state === 'needs') s.lift[n] = 0
      if (state === 'failed') s.glitch[n] = 0
    }
    s.prev = cur
    s.pal = undefined
    document.fonts.ready.then(kick)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key, size, theme, heads, turn, live, reduced])

  // Turn to a planet: ease (or step) the spin so it sits at the front.
  const bring = (name: string) => {
    const s = st.current, p = s.out?.planets.find((q) => q.name === name)
    if (!p) return
    s.spin.vel = 0
    s.spin.to = spinToFront(p.a0, s.spin.spin)
    if (reduced) { s.spin.spin = s.spin.to; s.spin.to = undefined }
    kick()
  }
  const setHover = (name: string | undefined) => {
    const s = st.current
    if (s.hover && name === s.hover.name && s.hover.dir > 0) return
    if (name) s.hover = { name, k: s.hover?.name === name ? s.hover.k : 0, dir: 1 }
    else if (s.hover) s.hover.dir = -1
    // The head looks at the hovered planet; on leave it turns back to the longest wait.
    const p = name ? s.out?.planets.find((q) => q.name === name) : undefined
    s.look = p && size ? (Math.abs(p.cx - size.w / 2) > 40 ? Math.sign(p.cx - size.w / 2) : 0) : undefined
    kick()
  }
  const select = (name: string) => {
    if (focus === name) { navigate(missionsHref(name)); return }
    setFocus(name); setHover(name); bring(name)
  }
  const clear = () => { setFocus(undefined); setHover(undefined) }

  const local = (e: PointerEvent) => { const r = box.current!.getBoundingClientRect(); return { x: e.clientX - r.left, y: e.clientY - r.top } }
  const down = (e: PointerEvent<HTMLDivElement>) => {
    if ((e.target as HTMLElement).closest('.hero-card, .hero-chip')) return
    const s = st.current
    s.drag = { id: e.pointerId, x0: e.clientX, x: e.clientX, t: e.timeStamp, spin0: s.spin.spin, moved: false }
    s.spin.vel = 0; s.spin.to = undefined
  }
  const move = (e: PointerEvent<HTMLDivElement>) => {
    const s = st.current, d = s.drag
    if (d && d.id === e.pointerId) {
      if (!d.moved && Math.abs(e.clientX - d.x0) > DRAG_PX) { d.moved = true; box.current?.setPointerCapture(e.pointerId) }
      if (d.moved && size) {
        const dt = Math.max(1, e.timeStamp - d.t)
        const v = (((e.clientX - d.x) / size.w) * DEG_PER_WIDTH) / dt
        s.spin.vel = 0.6 * v + 0.4 * s.spin.vel
        d.x = e.clientX; d.t = e.timeStamp
        s.spin.spin = d.spin0 + ((e.clientX - d.x0) / size.w) * DEG_PER_WIDTH
        // The head follows the drag direction.
        s.look = Math.abs(s.spin.vel) > 0.01 ? Math.sign(s.spin.vel) : s.look
        kick()
        return
      }
    }
    if (e.pointerType !== 'mouse' || d?.moved) return
    const { x, y } = local(e)
    // Over a moon of the hovered planet: keep the hover.
    if (st.current.out?.moons.some((o) => Math.hypot(x - o.x, y - o.y) < o.r + 4)) return
    const p = hit(st.current.out, x, y)
    if (p?.name !== st.current.hover?.name || !p) { if (!focus || p) setHover(p?.name ?? focus) }
    if (box.current) box.current.style.cursor = p || onHead(st.current.out, x, y) ? 'pointer' : 'grab'
  }
  const up = (e: PointerEvent<HTMLDivElement>) => {
    const s = st.current, d = s.drag
    if (!d || d.id !== e.pointerId) return
    s.drag = undefined
    if (d.moved) {
      // Release: inertia then detent; the head turns back to the longest wait.
      if (e.timeStamp - d.t > 80) s.spin.vel = 0
      s.look = undefined
      if (reduced) s.spin.vel = 0
      if (s.spin.vel === 0) s.spin.to = detent(a0s(), s.spin.spin)
      kick()
      return
    }
    const { x, y } = local(e)
    const p = hit(s.out, x, y)
    if (p) select(p.name)
    else if (onHead(s.out, x, y)) navigate('/desk')
    else clear()
  }
  const cancel = () => { const s = st.current; if (s.drag?.moved) { s.spin.vel = 0; s.spin.to = detent(a0s(), s.spin.spin); kick() } s.drag = undefined }
  const leave = (e: PointerEvent<HTMLDivElement>) => { if (e.pointerType === 'mouse' && !st.current.drag) setHover(focus) }

  const key2 = (e: KeyboardEvent<HTMLDivElement>) => {
    const ps = st.current.out?.planets ?? []
    if (e.key === 'ArrowLeft' || e.key === 'ArrowRight') {
      e.preventDefault()
      const n = stepTo(ps, st.current.spin.spin, focus, e.key === 'ArrowRight' ? 1 : -1)
      if (n) { setFocus(n); setHover(n); bring(n) }
    } else if (e.key === 'Enter' && focus) { e.preventDefault(); navigate(missionsHref(focus)) }
    else if (e.key === 'Escape' && focus) { e.preventDefault(); clear() }
  }

  // Measured overlay layout: nothing over the head, the planets, the moons or another label.
  useLayoutEffect(() => {
    if (!box.current || !out) return
    const disc = (x: number, y: number, r: number): Box => ({ x: x - r, y: y - r, w: 2 * r, h: 2 * r })
    const head: Box = { x: out.head.left + 30, y: out.head.top, w: Math.max(0, out.head.right - out.head.left - 60), h: 2 * (out.head.y - out.head.top) }
    const root = box.current, obs = [head, ...out.planets.map((p) => p.hidden ? disc(-99, -99, 0) : disc(p.cx, p.cy, p.rad)), ...out.moons.map((o) => disc(o.x, o.y, o.r + 6))]
    declutterDom(root, obs)
    // Web fonts change label widths: measure again once they are in.
    let on = true
    document.fonts?.ready.then(() => { if (on) declutterDom(root, obs) })
    return () => { on = false }
  })

  const fp = focus ? out?.planets.find((p) => p.name === focus) : undefined
  const hname = st.current.hover?.name
  const hp = hname ? out?.planets.find((p) => p.name === hname) : undefined
  return (
    <div ref={box} className="hero-box" tabIndex={0} role="group" aria-roledescription="holotable"
      aria-label={`CR-0W: ${say ?? m.line}. Left and Right pick a project, Enter opens its missions.`}
      style={out ? ({ '--lens-y': `${out.lens.y}px` } as CSSProperties) : undefined}
      onPointerDown={down} onPointerMove={move} onPointerUp={up} onPointerCancel={cancel} onPointerLeave={leave} onKeyDown={key2}>
      <canvas ref={cv} className="hero-cv" aria-hidden="true" />
      <div className="hero-say" data-prio={0} data-fixed aria-hidden="true"><b>CR-0W ›</b> {say ?? m.line}<i className="hero-cur" /></div>
      {out?.planets.map((p, i) => (
        <Fragment key={p.name}>
          {p.tag && <div className={`hero-tag${p.target ? ' on' : ''}`} style={{ left: p.tag.x, top: p.tag.y }} data-prio={p.target ? 2 : 3} data-own={1 + i} aria-hidden="true">{p.tag.text}</div>}
          {!p.hidden && (
            <div className={`hero-pl ${p.side}${p.target || p.name === focus ? ' on' : ''}`} style={{ left: p.x, top: p.y }} data-c={`${Math.round(p.cx)},${Math.round(p.cy)}`} data-prio={p.name === focus ? 2 : 4} data-own={1 + i} aria-hidden="true">
              {p.name.toUpperCase()}
              {p.note && <span style={{ color: `var(${VARS[p.state]})` }}> · {p.note}</span>}
            </div>
          )}
        </Fragment>
      ))}
      {out?.more && <div className="hero-rim" style={{ left: out.more.x, top: out.more.y }} data-prio={7} data-hide aria-hidden="true">+{out.more.n}</div>}
      {hp && out?.moons.map((o, j) => {
        const l = o.x < hp.cx
        return (
          <div key={o.id} className={`hero-ml${l ? ' l' : ''}`} style={{ left: o.x + (l ? -1 : 1) * (o.r + 12), top: o.y }} data-prio={5} data-hide data-own={1 + out.planets.length + j} aria-hidden="true">
            <Glyph state={o.state} size={11} /> <span>{o.title}</span>
          </div>
        )
      })}
      {out?.rim.map((r) => (
        <div key={r.text} className={`hero-rim ${r.kind}${r.side ? ` ${r.side}` : ''}`} style={{ left: r.x, top: r.y }} data-prio={r.kind === 'event' ? 6 : 8} data-hide aria-hidden="true">{r.text}</div>
      ))}
      {m.asks > 0 && out && (
        <div className="hero-chip" data-prio={0} data-fixed style={{ left: Math.min(out.head.right - 40, (size?.w ?? 0) - 60), top: out.head.top + 40 }} aria-hidden="true">
          {m.asks} {m.asks === 1 ? 'ASK' : 'ASKS'}
        </div>
      )}
      {fp && size && <Card project={fp.name} missions={missions} inbox={inbox} now={now} at={fp} w={size.w} />}
      {/* Screen readers and keyboard users get the planets as a plain list. */}
      <ul className="vh">
        {m.projects.map((p) => (
          <li key={p.name}><a href={missionsHref(p.name)} onClick={(e) => { e.preventDefault(); navigate(missionsHref(p.name)) }}>
            {p.name}: {p.state === 'needs' ? 'needs you' : p.state}, {p.open} open
          </a></li>
        ))}
      </ul>
    </div>
  )
}
