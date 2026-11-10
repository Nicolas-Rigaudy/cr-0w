import type { Mission } from '../data/types'
import { sameDay, weekday } from '../format'

// Runs started per day over the last 7 days, failed runs stacked on top.
export function RunsChart({ missions, now, height = 120 }: { missions: Mission[]; now: Date; height?: number }) {
  const days = Array.from({ length: 7 }, (_, i) => new Date(now.getFullYear(), now.getMonth(), now.getDate() - 6 + i))
  const runs = missions.flatMap((m) => m.runs)
  const data = days.map((d) => {
    const today = runs.filter((r) => sameDay(new Date(r.started), d))
    return { d, ok: today.filter((r) => r.cause !== 'error').length, failed: today.filter((r) => r.cause === 'error').length }
  })
  const max = Math.max(2, ...data.map((x) => x.ok + x.failed))
  const W = 300, H = height, L = 28, B = 30, T = 8
  const bw = (W - L - 8) / 7
  const y = (n: number) => T + (H - T - B) * (1 - n / max)
  const ticks = [0, Math.ceil(max / 2), max]
  return (
    <svg className="chart" viewBox={`0 0 ${W} ${H}`} width="100%" role="img"
      aria-label={`Runs per day, last 7 days: ${data.map((x) => `${weekday(x.d)} ${x.ok + x.failed}`).join(', ')}`}>
      {ticks.map((t) => (
        <g key={t}>
          <line className="axis" x1={L} x2={W - 4} y1={y(t)} y2={y(t)} strokeDasharray={t ? '2 3' : undefined} />
          <text x={L - 5} y={y(t) + 3} textAnchor="end">{t}</text>
        </g>
      ))}
      {data.map((x, i) => {
        const cx = L + 4 + i * bw
        return (
          <g key={i}>
            {x.ok > 0 && <rect className="bar" x={cx + bw * 0.2} width={bw * 0.6} y={y(x.ok)} height={y(0) - y(x.ok)} rx="2" />}
            {x.failed > 0 && <rect className="bar f" x={cx + bw * 0.2} width={bw * 0.6} y={y(x.ok + x.failed)} height={y(x.ok) - y(x.ok + x.failed)} rx="2" />}
            <text x={cx + bw / 2} y={H - B + 13} textAnchor="middle">{weekday(x.d)}</text>
          </g>
        )
      })}
      <text x={L + (W - L) / 2} y={H - 3} textAnchor="middle">day (runs started)</text>
      <text x={9} y={T + (H - T - B) / 2} textAnchor="middle" transform={`rotate(-90 9 ${T + (H - T - B) / 2})`}>runs</text>
    </svg>
  )
}
