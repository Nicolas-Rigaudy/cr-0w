import { renderToStaticMarkup } from 'react-dom/server'
import { describe, expect, it } from 'vitest'
import { fixtureSource } from '../data/fixtures'
import type { Dispatch } from '../data/live'
import { Queue } from './Ops'

const NOW = new Date('2026-11-09T14:00:00Z')
const ids = async (src = fixtureSource()) => (await src.ops('day=2026-11-09')).queue.filter((q) => q.kind === 'mission').map((q) => q.id)
const [A, B, C] = ['20261109-refresh-docs', '20261109-mail-team-bricks', '20261109-optional-tool-bricks']
const nop = () => {}

describe('Ops queue order (item_272)', () => {
  it('moves up, down, to a place and next; the ops answer follows', async () => {
    const src = fixtureSource()
    expect(await ids(src)).toEqual([A, B, C])
    await src.actions.missionMove(C, 'up')
    expect(await ids(src)).toEqual([A, C, B])
    await src.actions.missionMove(A, 'down')
    expect(await ids(src)).toEqual([C, A, B])
    await src.actions.missionMove(C, 9)
    expect(await ids(src)).toEqual([A, B, C])
    await src.actions.missionNext(B)
    expect(await ids(src)).toEqual([B, A, C])
    const q = (await src.ops('day=2026-11-09')).queue
    expect(q.map((x) => x.kind)).toEqual(['mission', 'mission', 'mission', 'routine', 'afk'])
  })

  it('refuses like core: unknown 404, not a mission 409, place 0 400', async () => {
    const src = fixtureSource()
    await expect(src.actions.missionMove('nope', 1)).rejects.toMatchObject({ status: 404 })
    await expect(src.actions.missionNext('crow-retro.timer')).rejects.toMatchObject({ status: 409 })
    await expect(src.actions.missionMove(A, 0)).rejects.toMatchObject({ status: 400 })
  })

  it('start now takes the mission off the queue and adds a run past the cap', async () => {
    const src = fixtureSource()
    const slots = async () => ((await src.section('dispatch')).data as Dispatch).slots
    const before = await slots()
    await src.actions.missionStart(B)
    expect(await ids(src)).toEqual([A, C])
    expect((await slots()).busy).toBe(before.busy + 1)
  })

  it('draws grip, place, arrows and RUN ▾ on mission rows; slots 3/2 shows over cap', async () => {
    const src = fixtureSource()
    const queue = (await src.ops('day=2026-11-09')).queue
    const html = renderToStaticMarkup(<Queue queue={queue} qsum="3 missions · 1 routine" now={NOW} slots={{ busy: 3, max: 2 }} src={src}
      reload={async () => {}} onCancel={nop} onRoutine={nop} onEdit={nop} onSkip={nop} />)
    expect(html.match(/class="ops-grip"/g)).toHaveLength(3)
    expect(html.match(/RUN ▾/g)).toHaveLength(3)
    expect(html).toContain('Drag ⠿ or use ▲ ▼ to reorder.')
    expect(html).toContain('class="ops-slots over">slots 3/2<span class="ops-busy"> over cap</span>')
    expect(html).toMatch(/aria-label="Move Refresh the docs index up" disabled=""/)
    expect(html).toMatch(/aria-label="Move Optional tool bricks down" disabled=""/)
    expect(html).not.toMatch(/aria-label="Move Mail, team and vault bricks (up|down)" disabled/)
  })
})
