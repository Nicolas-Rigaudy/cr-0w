import { describe, expect, it } from 'vitest'
import stateMrs from '../../fixtures/states/mrs.json'
import stateWeek from '../../fixtures/states/agenda-week.json'
import { agendaOf, countsOf, duePrompt, hoursOf, isoWeek } from './agenda'
import type { Snapshot } from './feed'
import type { SectionBody } from './types'

const NOW = new Date('2026-11-09T16:10:00+01:00')
const body = (data: unknown, more: Partial<SectionBody> = {}): SectionBody => ({ data, stale: false, updated_at: null, ...more })
const snap = (sections: Record<string, SectionBody>, now = NOW) => ({ sections, now }) as unknown as Snapshot
const REVIEW = { id: 'mr-wait', check: 'review', what: 'Review waits on demo !12', why: 'open 2 days', action: 'review it', due: null }
const full = () => snap({
  today: body({ week: stateWeek.week }),
  checks: body({ rows: [REVIEW, ...stateWeek.checks] }),
  mrs: body(stateMrs),
})

describe('Agenda mapping (item_226)', () => {
  it('lays seven days from today with all-day chips and timed meetings in order', () => {
    const a = agendaOf(full())
    expect(a.days.map((d) => d.key)).toEqual(['2026-11-09', '2026-11-10', '2026-11-11', '2026-11-12', '2026-11-13', '2026-11-14', '2026-11-15'])
    expect(a.days[0].today).toBe(true)
    expect(a.days.map((d) => d.allDay.join())).toEqual(['Office', 'Office', 'Office', 'Remote', 'Remote', '', ''])
    expect(a.days[3].meetings.map((m) => m.title)).toEqual(['Weekly prep', 'Daily Meeting', 'HR point'])
    // The eighth day of crow-today `week` stays out of the grid.
    expect(a.days.flatMap((d) => d.meetings).some((m) => m.title === 'Management weekly')).toBe(false)
    expect(a.week).toBe(46)
  })

  it('marks what is over today and picks the next meeting', () => {
    const a = agendaOf(full())
    expect(a.days[0].meetings.map((m) => [m.title, m.state])).toEqual([['Daily Meeting', 'done'], ['Sprint review', 'next'], ['Call with the accountant', '']])
    expect(a.next).toMatchObject({ title: 'Sprint review', from: 17 * 60, to: 17 * 60 + 45 })
  })

  it('puts dated checks on their day, undated ones and MR waits on today', () => {
    const a = agendaOf(full())
    expect(a.days[0].due.map((x) => [x.kind, x.title])).toEqual([['mr', '!242 review waits 4d']])
    expect(a.days[1].due.map((x) => [x.kind, x.sub])).toEqual([['th', 'follow-up · unanswered for 3d']])
    expect(a.days[4].due.map((x) => [x.kind, x.title])).toEqual([['dt', 'Sprint 34 ends Fri 13 Nov'], ['th', 'Weekly rollup']])
    expect(countsOf(a)).toEqual({ todayMeetings: 3, todayDue: 1, meetings: 11, due: 5 })
  })

  it('keeps the review check when the forge serves no MR, overdue checks on today', () => {
    const a = agendaOf(snap({ today: body({ week: [] }), checks: body({ rows: [REVIEW, { ...REVIEW, id: 'old', check: 'ritual', due: '2026-11-02' }] }) }))
    expect(a.days[0].due.map((x) => x.key)).toEqual(['c:mr-wait', 'c:old'])
  })

  it('falls back to events and tomorrow when crow-today has no week', () => {
    const a = agendaOf(snap({ today: body({ events: [{ start: '2026-11-09T18:00', end: '2026-11-09T18:30', subject: 'A', all_day: false }], tomorrow: [{ start: '2026-11-10T09:00', end: '2026-11-10T09:30', subject: 'B', all_day: false }] }) }))
    expect(a.days[0].meetings.map((m) => m.title)).toEqual(['A'])
    expect(a.days[1].meetings.map((m) => m.title)).toEqual(['B'])
  })

  it('reports an off or stale calendar', () => {
    expect(agendaOf(snap({ today: body(null, { off: true }) })).off).toBe(true)
    expect(agendaOf(snap({ today: body({ week: [] }, { stale: true, error: 'not signed in' }) }))).toMatchObject({ stale: true, error: 'not signed in' })
  })

  it('widens the hour range past 09:00 to 19:00 only when a meeting needs it', () => {
    const a = agendaOf(full())
    expect(hoursOf(a.days, 9, 19)).toEqual([9, 20])
    expect(hoursOf(a.days.slice(1), 9, 19)).toEqual([9, 19])
  })

  it('writes the session prompt of a due check', () => {
    const x = agendaOf(full()).days[1].due[0]
    expect(duePrompt(x)).toBe('Legal: contract questions. unanswered for 3d. Next step: draft a reply.')
  })

  it('numbers ISO weeks', () => {
    expect(isoWeek(new Date(2026, 10, 24))).toBe(48)
    expect(isoWeek(new Date(2027, 0, 1))).toBe(53)
  })
})
