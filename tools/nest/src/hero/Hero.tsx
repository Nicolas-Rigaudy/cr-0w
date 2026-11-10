import { useEffect, useState } from 'react'
import type { UiState } from '../data/state'
import type { InboxItem, Mission } from '../data/types'
import { elapsed } from '../format'
import type { HeroProject, Mood } from './model'
import { heroModel } from './model'

// Static hero (item_155): orbit of projects around the ASCII head. No animation.

interface Heads { palette: string[]; moods: Record<Mood, { glyphs: string[]; colours: string[] }> }

const W = 1000
const H = 1200
const CX = 500
const CY = 790 // orbit centre
const IN = { rx: 340, ry: 125 } // projects
const OUT = { rx: 460, ry: 175 } // queue
const FRONT = [140, 40, 165, 15, 200, 340, 240, 300] // degrees, 90 = nearest the viewer
const BACK = [235, 305, 205, 335, 180, 0]
const FS = 15 // head font size
const LH = 15.6
const CW = FS * 0.6
const COLS = 80

const KEYS = [...Array(94)].map((_, i) => String.fromCharCode(33 + i)).filter((c) => !'"\\'.includes(c)).join('')

const COL: Record<UiState, string> = {
  needs: 'var(--c-n)', working: 'var(--c-w)', failed: 'var(--c-f)', sleeping: 'var(--c-s)', queued: 'var(--c-q)', done: 'var(--c-d)',
}

const at = (r: { rx: number; ry: number }, deg: number) => {
  const a = (deg * Math.PI) / 180
  return { x: CX + r.rx * Math.cos(a), y: CY + r.ry * Math.sin(a) }
}

function Mark({ s, x, y }: { s: UiState; x: number; y: number }) {
  const c = COL[s]
  switch (s) {
    case 'needs':
      return <rect x={x - 9} y={y - 9} width={18} height={18} transform={`rotate(45 ${x} ${y})`} style={{ fill: c, filter: 'drop-shadow(0 0 6px var(--c-n))' }} />
    case 'working':
      return <circle cx={x} cy={y} r={10} style={{ fill: c, filter: 'drop-shadow(0 0 8px var(--c-w))' }} />
    case 'failed':
      return <circle cx={x} cy={y} r={10} fill="none" strokeWidth={2.5} strokeDasharray="12 5" style={{ stroke: c }} />
    case 'sleeping':
      return <circle cx={x} cy={y} r={9} fill="none" strokeWidth={2} strokeDasharray="1 4" strokeLinecap="round" style={{ stroke: c }} />
    case 'queued':
      return <circle cx={x} cy={y} r={7} fill="none" strokeWidth={1.5} style={{ stroke: c }} />
    default:
      return <circle cx={x} cy={y} r={5} style={{ fill: c, opacity: 0.6 }} />
  }
}

function Head({ heads, mood }: { heads: Heads; mood: Mood }) {
  const h = heads.moods[mood]
  const top = CY - 40 - h.glyphs.length * LH
  const left = CX - (COLS / 2) * CW
  return (
    <text className="hero-head" fontSize={FS} xmlSpace="preserve">
      {h.glyphs.map((g, row) => {
        const cs = h.colours[row]
        const runs: { x: number; t: string; k: string }[] = []
        for (let i = 0; i < g.length; i++) {
          if (cs[i] === ' ') continue
          const last = runs.at(-1)
          if (last && last.k === cs[i] && last.x + last.t.length === i) last.t += g[i]
          else runs.push({ x: i, t: g[i], k: cs[i] })
        }
        return runs.map((r) => {
          const [fill, alpha] = heads.palette[KEYS.indexOf(r.k)].split('/')
          return <tspan key={`${row}-${r.x}`} x={left + r.x * CW} y={top + row * LH} fill={fill} fillOpacity={alpha}>{r.t}</tspan>
        })
      })}
    </text>
  )
}

const short = (t: string) => (t.length > 24 ? `${t.slice(0, 23)}…` : t)

function Node({ p, deg, now }: { p: HeroProject; deg: number; now: Date }) {
  const { x, y } = at(IN, deg)
  // Missions fan outward from the core around their project node.
  const ux = x - CX, uy = (y - CY) * 2.6
  const base = Math.atan2(uy, ux)
  const sats = p.missions.filter((m) => m.state !== 'queued')
  const step = 0.62
  const side = x < CX ? 'end' : 'start'
  return (
    <g>
      <line x1={CX} y1={CY} x2={x} y2={y} strokeWidth={p.live ? 1.6 : 1} strokeDasharray={p.live ? undefined : '3 6'}
        style={{ stroke: p.live ? COL[p.state] : 'var(--line)', opacity: p.live ? 0.75 : 0.6 }} />
      {sats.map((m, i) => {
        const a = base + (i - (sats.length - 1) / 2) * step
        const sx = x + Math.cos(a) * 92, sy = y + Math.sin(a) * 78
        const tx = sx + (side === 'end' ? -18 : 18)
        return (
          <g key={m.id}>
            <line x1={x} y1={y} x2={sx} y2={sy} style={{ stroke: COL[m.state], opacity: 0.45 }} />
            <Mark s={m.state} x={sx} y={sy} />
            <text x={tx} y={sy + 5} textAnchor={side} className="hero-sat" style={{ fill: COL[m.state] }}>
              {m.state === 'working' && m.started ? elapsed(m.started, now) : short(m.title)}
            </text>
          </g>
        )
      })}
      <circle cx={x} cy={y} r={11} className="hero-hub" style={{ stroke: COL[p.state] }} />
      <circle cx={x} cy={y} r={4} style={{ fill: COL[p.state] }} />
      <text x={x} y={y - 24} textAnchor="middle" className="hero-lbl" style={{ fill: p.live ? 'var(--fg)' : 'var(--dim)' }}>
        {p.name.toUpperCase()} <tspan style={{ fill: COL[p.state] }}>· {p.state.toUpperCase()}</tspan>
      </text>
    </g>
  )
}

export function Hero({ missions, inbox, now }: { missions: Mission[]; inbox: InboxItem[]; now: Date }) {
  const m = heroModel(missions, inbox)
  const [heads, setHeads] = useState<Heads>()
  useEffect(() => {
    let live = true
    import('./heads.json').then((d) => live && setHeads(d.default as Heads))
    return () => { live = false }
  }, [])

  const queued = m.projects.flatMap((p) => p.missions.filter((x) => x.state === 'queued').map((x) => ({ ...x, project: p.name })))

  return (
    <svg className="hero-svg" viewBox={`0 0 ${W} ${H}`} preserveAspectRatio="xMidYMid meet" role="img" aria-label={`CR-0W: ${m.line}`}>
      <defs>
        <radialGradient id="hero-glow" cx="50%" cy="66%" r="55%">
          <stop offset="0" stopColor="var(--accent)" stopOpacity="0.22" />
          <stop offset="0.55" stopColor="var(--accent)" stopOpacity="0.06" />
          <stop offset="1" stopColor="var(--accent)" stopOpacity="0" />
        </radialGradient>
      </defs>
      <rect width={W} height={H} fill="url(#hero-glow)" />

      <text x={40} y={70} className="hero-line">
        <tspan className="hero-who">CR-0W ›</tspan> {m.line}
      </text>

      <ellipse cx={CX} cy={CY} rx={OUT.rx} ry={OUT.ry} fill="none" strokeDasharray="4 8" style={{ stroke: 'var(--line)' }} />
      <ellipse cx={CX} cy={CY} rx={IN.rx} ry={IN.ry} fill="none" strokeWidth={1.4} style={{ stroke: 'var(--accent)', opacity: 0.55 }} />
      <ellipse cx={CX} cy={CY} rx={90} ry={32} fill="none" style={{ stroke: 'var(--accent)', opacity: 0.4 }} />

      {queued.slice(0, BACK.length).map((q, i) => {
        const { x, y } = at(OUT, BACK[i])
        return (
          <g key={q.id}>
            <Mark s="queued" x={x} y={y} />
            <text x={x} y={y - 16} textAnchor="middle" className="hero-q">Q{i + 1} · {q.title}</text>
          </g>
        )
      })}

      {m.projects.slice(0, FRONT.length).map((p, i) => <Node key={p.name} p={p} deg={FRONT[i]} now={now} />)}

      <ellipse cx={CX} cy={CY} rx={58} ry={18} style={{ fill: 'var(--accent)', opacity: 0.18 }} />
      <ellipse cx={CX} cy={CY} rx={34} ry={10} style={{ fill: 'var(--accent)', opacity: 0.45 }} />

      {heads && <Head heads={heads} mood={m.mood} />}

      {m.asks > 0 && (
        <g>
          <rect x={CX - 52} y={CY + 30} width={104} height={28} rx={14} className="hero-chip" />
          <text x={CX} y={CY + 49} textAnchor="middle" className="hero-chip-t">{m.asks} {m.asks === 1 ? 'ASK' : 'ASKS'}</text>
        </g>
      )}
    </svg>
  )
}
