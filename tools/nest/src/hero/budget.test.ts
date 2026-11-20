import { describe, expect, it } from 'vitest'
import { budgetWatch } from './budget'

const feed = (w: ReturnType<typeof budgetWatch>, xs: number[]) => xs.map(w)

describe('frame budget watch', () => {
  it('skips the warm-up frames and decides once on the frame that fills the window', () => {
    const r = feed(budgetWatch(16, 2, 4), [100, 100, 5, 5, 5, 5, 50, 50])
    expect(r).toEqual([undefined, undefined, undefined, undefined, undefined, 'ok', undefined, undefined])
  })

  it('misses when the median is over the budget', () => {
    expect(feed(budgetWatch(16, 0, 5), [20, 3, 30, 18, 25]).at(-1)).toBe('miss')
  })

  it('ignores a few slow spikes', () => {
    expect(feed(budgetWatch(16, 0, 5), [90, 4, 80, 5, 6]).at(-1)).toBe('ok')
  })

  it('takes the mean of the two middle frames on an even window', () => {
    expect(feed(budgetWatch(16, 0, 4), [10, 17, 17, 40]).at(-1)).toBe('miss')
    expect(feed(budgetWatch(16, 0, 4), [10, 15, 17, 40]).at(-1)).toBe('ok')
  })
})
