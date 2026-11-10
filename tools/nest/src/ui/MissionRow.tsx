import type { InboxItem, Mission } from '../data/types'
import { LABEL, lastRun, uiState } from '../data/state'
import { duration, time, title } from '../format'
import { Link } from '../router'
import { Glyph, stateClass } from './Glyph'

export function subline(m: Mission, now: Date) {
  const run = lastRun(m)
  if (!run) return `${m.project} · ${m.kind} · not started`
  if (!run.ended) return `${m.project} · ${m.kind} · running ${duration(run.started, now)}`
  if (run.cause === 'usage_limit') return `${m.project} · ${m.kind} · usage limit at ${time(run.ended)}`
  if (run.cause === 'error') return `${m.project} · ${m.kind} · error, exit ${run.exit}`
  return `${m.project} · ${m.kind} · ${run.cause ?? 'ended'} at ${time(run.ended)}`
}

export function MissionRow({ m, inbox, now, big }: { m: Mission; inbox: InboxItem[]; now: Date; big?: boolean }) {
  const s = uiState(m, inbox)
  return (
    <Link href={`/missions/${m.id}`} className={`row${big ? ' r2' : ''}${s === 'done' ? ' past' : ''}`}>
      <Glyph state={s} size={big ? 18 : 14} />
      <span className="t">
        <b>{title(m)}</b>
        <span>{subline(m, now)}</span>
      </span>
      <span className={`sub ${stateClass(s)}`}>{LABEL[s]}</span>
      <span className="chev" aria-hidden="true">›</span>
    </Link>
  )
}
