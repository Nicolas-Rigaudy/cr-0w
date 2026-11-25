import { ApiError } from '../data/api'
import type { LiveReply } from '../data/live'
import type { Mission } from '../data/types'
import { title } from '../format'
import type { Confirm } from './Sheet'

// Mission sheets shared by Desk and the session chat (item_279).

type MissionActs = {
  missionStop(mission: string): Promise<unknown>
  missionCancel(mission: string, reason?: string): Promise<unknown>
}

const short = (id: string) => id.replace(/^\d{8}-/, '')

// A refused reply (ok false) throws, so the sheet that sent it shows the message.
export function refused(r: LiveReply): LiveReply {
  if (!r.ok) throw new Error(r.message || `${r.action} ${r.target} refused`)
  return r
}

// Runs one mission call; a 409 refetches first since the view was stale. The error goes
// back to the sheet.
async function call(f: () => Promise<unknown>, refresh: () => void) {
  try {
    await f()
  } catch (e) {
    if (e instanceof ApiError && e.status === 409) refresh()
    throw e
  }
}

// Stop on a run of a mission: the main action cancels the mission (its run stops and the
// tick leaves it). Restart run kills the run only; the tick resumes it.
export function stopMissionRun(o: {
  mission: Mission; run: string; name: string; act: MissionActs
  refresh: () => void; done: (what: 'cancelled' | 'restarted') => void
}): Confirm {
  const { mission: m, act, refresh } = o
  return {
    title: 'Stop mission run', lead: `Stop the run ${o.name} and cancel ${title(m)}.`,
    target: [['mission', short(m.id)], ['run', o.run], ['project', m.project]],
    note: 'Same as crow mission cancel, logged with surface nest. Its commits stay on the run branch. Restart run stops this run only; the mission starts a new one at the next tick.',
    ok: 'Stop and cancel the mission', danger: true,
    run: async () => { await call(() => act.missionCancel(m.id), refresh); o.done('cancelled'); refresh() },
    alt: { label: 'Restart run', run: async () => { await call(() => act.missionStop(m.id), refresh); o.done('restarted'); refresh() } },
  }
}

// Cancel a mission: a queued one leaves the queue, a running one has its run stopped.
export function cancelMission(o: { mission: Mission; act: MissionActs; refresh: () => void; done: () => void }): Confirm {
  const { mission: m, act, refresh } = o
  const running = m.state === 'running'
  return {
    title: 'Cancel the mission?',
    lead: running ? `${title(m)} stops its run and ends for good.` : `${title(m)} leaves the queue for good.`,
    target: [['mission', short(m.id)], ['project', m.project], ['state', m.state]],
    note: `Same as crow mission cancel. Logged with surface nest.${running ? ' Its commits stay on the run branch.' : ''}`,
    ok: 'Cancel mission', danger: true,
    run: async () => { await call(() => act.missionCancel(m.id), refresh); o.done(); refresh() },
  }
}
