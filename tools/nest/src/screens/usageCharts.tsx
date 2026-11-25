import { useLayoutEffect, useRef, useState, type ReactNode } from 'react'
import { time } from '../format'
import { cellsLabel, dow, niceMax, type Bar } from '../data/usageView'

// Usage charts (item_227, board usage.html C): labelled axes, a legend under each chart,
// the floor line, the projection dashed, direct value labels. Colours are --c1..--c4.

export const COLORS = ['var(--c1)', 'var(--c2)', 'var(--c3)', 'var(--c4)']
const GRID = 'var(--u-grid)'

// With fit, the viewBox width follows the measured width, so chart text keeps its CSS size
// on any phone width; w is the first-render guess. Without fit the viewBox scales (laptop).
function useFit(w: number, fit?: boolean) {
  const ref = useRef<SVGSVGElement>(null)
  const [mw, setMw] = useState(0)
  useLayoutEffect(() => {
    const el = ref.current
    if (!fit || !el) return
    const read = () => { const x = Math.round(el.getBoundingClientRect().width); if (x > 0) setMw(x) }
    read()
    if (typeof ResizeObserver === 'undefined') return
    const ro = new ResizeObserver(read)
    ro.observe(el)
    return () => ro.disconnect()
  }, [fit])
  return [ref, fit && mw ? mw : w] as const
}

// % used over a window (5h or week): area to now, dashed line to the projection at reset.
export function WindowChart({ w: w0, h, from, to, now, points, used, proj, floorLeft, ticks, title, fit }: {
  w: number; h: number; from: Date; to: Date; now: Date; points: Array<[Date, number]>; used: number; proj: number | null
  floorLeft: number; ticks: Array<[Date, string]>; title: string; fit?: boolean
}) {
  const [ref, w] = useFit(w0, fit)
  const L = 40, R = 14, T = 20, B = 24, iw = w - L - R, ih = h - T - B
  const span = to.getTime() - from.getTime()
  const X = (d: Date) => L + Math.min(1, Math.max(0, (d.getTime() - from.getTime()) / span)) * iw
  const Y = (p: number) => T + ih * (1 - p / 100)
  const path = points.map(([d, p], i) => `${i ? 'L' : 'M'}${X(d).toFixed(1)},${Y(p).toFixed(1)}`).join('')
  const nx = X(now)
  const fl = 100 - floorLeft
  const nowLeft = nx - L > iw * 0.45
  return (
    <svg ref={ref} className="uch" viewBox={`0 0 ${w} ${h}`} height={fit ? h : undefined} role="img" aria-label={`${title}: ${used}% used now${proj !== null ? `, ${proj}% at reset` : ''}`}>
      {[0, 25, 50, 75, 100].map((p) => (
        <g key={p}><line x1={L} x2={w - R} y1={Y(p)} y2={Y(p)} stroke={GRID} /><text x={L - 6} y={Y(p) + 3.5} textAnchor="end">{p}%</text></g>
      ))}
      {ticks.map(([d, t], i) => (
        <text key={i} x={X(d)} y={h - 8} textAnchor={i === ticks.length - 1 ? 'end' : i ? 'middle' : 'start'}>{t}</text>
      ))}
      <line x1={L} x2={w - R} y1={Y(fl)} y2={Y(fl)} className="ufl" strokeDasharray="2 3" />
      <text x={L + 6} y={Y(fl) - 4} className="ufl-t">floor {floorLeft}% left</text>
      <path d={`${path}L${nx},${Y(0)}L${X(from)},${Y(0)}Z`} fill="var(--c1)" opacity=".18" />
      <path d={path} fill="none" stroke="var(--c1)" strokeWidth="2" />
      {proj !== null && (
        <>
          <line x1={nx} y1={Y(used)} x2={X(to)} y2={Y(proj)} stroke="var(--dim)" strokeWidth="2" strokeDasharray="5 4" />
          <text x={X(to)} y={Y(proj) - 7} textAnchor="end" className="vl">→ {proj}% at reset</text>
        </>
      )}
      <line x1={nx} x2={nx} y1={T} y2={T + ih} stroke="var(--accent)" />
      <circle cx={nx} cy={Y(used)} r="4" fill="var(--c1)" stroke="var(--bg)" strokeWidth="2" />
      <text x={nowLeft ? nx - 8 : nx + 8} y={Y(used) - 8} textAnchor={nowLeft ? 'end' : 'start'} className="vl">{used}% · now {time(now)}</text>
      <text x={L} y={T - 8} className="ax">% used</text>
    </svg>
  )
}

// The ticks of a 5h window (hours) and of a week (days).
export function hourTicks(from: Date, to: Date): Array<[Date, string]> {
  const out: Array<[Date, string]> = []
  for (let t = from.getTime(); t <= to.getTime() + 1; t += 3600e3) out.push([new Date(t), time(new Date(t))])
  return out
}
export function dayTicks(from: Date, to: Date, short = false): Array<[Date, string]> {
  const out: Array<[Date, string]> = []
  for (let t = from.getTime(); t <= to.getTime() + 1; t += 864e5) {
    const d = new Date(t)
    const name = dow(`${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`)
    out.push([d, short ? name[0] : name])
  }
  return out
}

export interface DayChartProps {
  w: number; h: number; bars: Bar[]; fmt: (v: number) => string; axis: string; colors: string[]
  pick: number; onPick: (i: number) => void; today: string
  // The 5h windows of one account under each day: used % per window, null = no data.
  cells?: Array<number[] | null | undefined>
  fit?: boolean
}

// Stacked bars per day; the dearest day and the picked day get a direct value label.
export function DayChart({ w: w0, h, bars, fmt, axis, colors, pick, onPick, today, cells, fit }: DayChartProps) {
  const [ref, w] = useFit(w0, fit)
  const withCells = !!cells
  const C = withCells ? 30 : 0
  const L = 44, R = 8, T = 14, B = 24 + C, iw = w - L - R, ih = h - T - B
  // Headroom so the peak value label stays clear of the axis title.
  const mx = niceMax(1.15 * Math.max(0, ...bars.map((b) => b.total)))
  const bw = iw / Math.max(1, bars.length)
  const Y = (v: number) => T + ih * (1 - v / mx)
  const top = bars.reduce((m, b, i) => (b.total > (bars[m]?.total ?? -1) ? i : m), 0)
  const step = bars.length > 14 ? Math.ceil(bars.length / 10) : 1
  const labelCells = bw >= 64
  return (
    <svg ref={ref} className="uch" viewBox={`0 0 ${w} ${h}`} height={fit ? h : undefined} role="img" aria-label={axis}>
      {[0, 0.25, 0.5, 0.75, 1].map((f) => (
        <g key={f}><line x1={L} x2={w - R} y1={Y(mx * f)} y2={Y(mx * f)} stroke={GRID} /><text x={L - 6} y={Y(mx * f) + 3.5} textAnchor="end">{fmt(mx * f)}</text></g>
      ))}
      {bars.map((b, i) => {
        const x = L + i * bw + bw * 0.2
        const bwi = bw * 0.6
        let acc = 0
        const isToday = b.date === today
        const parts: ReactNode[] = b.stacks.map((s, k) => {
          const y0 = Y(acc)
          acc += s.value
          const y1 = Y(acc)
          const hh = Math.max(0, y0 - y1 - (k ? 1.5 : 0))
          return hh > 0 ? <rect key={s.key} x={x} y={y1} width={bwi} height={hh} fill={colors[k % colors.length]} rx={k === b.stacks.length - 1 ? 2.5 : 1} /> : null
        })
        const cy = T + ih + 22
        const cl = cells?.[i]
        return (
          <g key={b.date} className={`ubar${i === pick ? ' on' : ''}`} onMouseEnter={() => onPick(i)} onClick={() => onPick(i)}>
            <rect x={L + i * bw} y={T} width={bw} height={ih + 18 + C} className="uhit" />
            {parts}
            {b.none && <text x={x + bwi / 2} y={Y(0) - 6} textAnchor="middle" className="nd">no data</text>}
            {(i === top || i === pick) && b.total > 0 && <text x={x + bwi / 2} y={Y(b.total) - 5} textAnchor="middle" className="vl">{fmt(b.total)}</text>}
            {i % step === 0 || i === bars.length - 1 ? (
              <text x={x + bwi / 2} y={T + ih + 14} textAnchor="middle" className={isToday || i === pick ? 'vl' : undefined}>{isToday ? 'today' : bars.length > 10 ? b.date.slice(8, 10) : dow(b.date)}</text>
            ) : null}
            {withCells && <Cells x={L + i * bw + 3} y={cy} w={bw - 6} cells={cl} label={labelCells} />}
          </g>
        )
      })}
      <text x={L} y={T - 4} className="ax">{axis}</text>
      {withCells && <text x={4} y={T + ih + 32} className="ax">5h</text>}
    </svg>
  )
}

// One day's 5h windows as cells, each filled to its used %; the label under them.
function Cells({ x, y, w, cells, label }: { x: number; y: number; w: number; cells?: number[] | null; label: boolean }) {
  if (cells === null) return <text x={x + w / 2} y={y + 8} textAnchor="middle" className="nd">–</text>
  if (!cells?.length) return null
  const n = cells.length
  const cw = Math.min(10, (w - (n - 1) * 2) / n)
  const total = n * cw + (n - 1) * 2
  const x0 = x + (w - total) / 2
  return (
    <g aria-label={cellsLabel(cells)}>
      {cells.map((c, i) => {
        const hh = 9 * Math.min(100, c) / 100
        return (
          <g key={i}>
            <rect x={x0 + i * (cw + 2)} y={y} width={cw} height={9} className="ucell" />
            <rect x={x0 + i * (cw + 2)} y={y + 9 - hh} width={cw} height={hh} fill={c >= 100 ? 'var(--c2)' : 'var(--c1)'} />
          </g>
        )
      })}
      {label && <text x={x + w / 2} y={y + 20} textAnchor="middle" className="cl">{cellsLabel(cells)}</text>}
    </g>
  )
}

// Horizontal bars with direct value labels.
export function HBars({ rows, color = 'var(--c1)', label = 90, fmt }: { rows: Array<[string, number, string?]>; color?: string; label?: number; fmt: (v: number) => string }) {
  const mx = Math.max(0, ...rows.map((r) => r[1])) || 1
  return (
    <div className="uhb" style={{ gridTemplateColumns: `${label}px 1fr 64px` }}>
      {rows.map(([k, v, c]) => (
        <div key={k} className="uhr">
          <span className="k" title={k}>{k}</span>
          <span className="tr"><i style={{ width: `${(v / mx) * 100}%`, background: c ?? color }} /></span>
          <span className="v num">{fmt(v)}</span>
        </div>
      ))}
    </div>
  )
}

// Burn sparkline of an account: used % of the week per day, with the floor line and a y label.
export function Spark({ values, floorLeft, label }: { values: Array<number | null>; floorLeft: number; label: string }) {
  const w = 90, h = 24
  const mx = Math.max(100 - floorLeft, ...values.map((v) => v ?? 0)) * 1.05
  const pts = values.map((v, i) => [i * (w / Math.max(1, values.length - 1)), v === null ? null : h - 2 - (v / mx) * (h - 4)] as const)
  const line = pts.filter((p) => p[1] !== null).map(([x, y]) => `${x.toFixed(1)},${y!.toFixed(1)}`).join(' ')
  const fy = h - 2 - ((100 - floorLeft) / mx) * (h - 4)
  return (
    <svg className="usp" viewBox={`0 0 ${w} ${h}`} width={w} height={h} role="img" aria-label={label}>
      <line x1="0" x2={w} y1={fy} y2={fy} className="ufl" strokeDasharray="2 3" />
      {line && <polyline points={line} fill="none" stroke="var(--c1)" strokeWidth="2" />}
    </svg>
  )
}

export function Legend({ items, note }: { items: Array<[string, string, ('box' | 'dash' | 'floor')?]>; note?: string }) {
  return (
    <div className="uleg">
      {items.map(([t, c, k]) => (
        <span key={t} className={k === 'floor' ? 'c-f' : undefined}>{k === 'dash' ? <i className="ln" /> : k === 'floor' ? '┄ ' : <i style={{ background: c }} />}{t}</span>
      ))}
      {note && <span className="dim">{note}</span>}
    </div>
  )
}
