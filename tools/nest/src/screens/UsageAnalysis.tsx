import { useState } from 'react'
import type { UsageData, UsageLimit } from '../data/usage'
import { bars, cut, localDay, modelName, modelsOf, periodLabel, sumGroups, usd, type Period } from '../data/usageView'
import { Deck } from './Mrs'
import { COLORS, DayChart, HBars, Legend } from './usageCharts'

// The ANALYSIS pane of the Usage page: one chart at a time.

export interface AnalysisProps {
  data: UsageData
  lim?: UsageLimit
  now: Date
  period: Period
  setPeriod: (p: Period) => void
  phone?: boolean
}

export function Analysis({ data, now, period, phone }: AnalysisProps) {
  const days = cut(data.cost.data?.by_day ?? [], period)
  const [pick, setPick] = useState(days.length - 1)
  const models = modelsOf(days)
  const legend = <Legend items={models.map((m, i) => [modelName(m), COLORS[i % COLORS.length]])} />
  const chart = <DayChart w={phone ? 330 : 620} h={phone ? 200 : 220} bars={bars(days, 'usd', [])} fmt={usd} axis="$ per day" colors={COLORS} pick={pick} onPick={setPick} today={localDay(now)} />
  const proj = sumGroups(days, (d) => d.by_project).map(([k, g]) => [k, g.cost_usd] as [string, number])
  return (
    <Deck n={phone ? 4 : 3} title="Analysis">
      <div className="ubar-top"><nav className="segs" aria-label="Analysis"><span className="seg on">{phone ? 'Day' : 'Per day'}</span></nav><span className="dim">one chart at a time</span></div>
      {!days.length ? <div className="empty">No cost audit yet.</div> : phone ? (
        <>{chart}{legend}<h3 className="ust">Per project</h3><HBars rows={proj} fmt={usd} label={60} /></>
      ) : (
        <div className="uan">
          {chart}
          <div className="uside">{legend}<h3 className="ust">{periodLabel(period)}</h3><HBars rows={proj} fmt={usd} label={80} /></div>
        </div>
      )}
    </Deck>
  )
}
