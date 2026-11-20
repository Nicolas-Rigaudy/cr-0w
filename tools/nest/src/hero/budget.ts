// Frame budget (item_220 AC5): the hero times each table frame and falls back to the plain table
// when the median of a window of frames is over one 60 fps frame. Decided once per page load.
// localStorage nest.table = 'lit' | 'plain' forces a table and skips the watch.

export const BUDGET_MS = 16
// Frames skipped first (fonts, caches, GL setup), then frames in the window.
export const WARM = 5, WINDOW = 30

export type TablePref = 'lit' | 'plain' | undefined

export function tablePref(): TablePref {
  try {
    const v = localStorage.getItem('nest.table')
    return v === 'lit' || v === 'plain' ? v : undefined
  } catch { return undefined }
}

// Feed one frame time per call. Returns 'miss' or 'ok' once, on the frame that fills the window, then undefined.
export function budgetWatch(ms = BUDGET_MS, warm = WARM, n = WINDOW) {
  let seen = 0
  const xs: number[] = []
  return (dt: number): 'miss' | 'ok' | undefined => {
    if (++seen <= warm || xs.length >= n) return undefined
    xs.push(dt)
    if (xs.length < n) return undefined
    const s = [...xs].sort((a, b) => a - b)
    return (s[(n - 1) >> 1] + s[n >> 1]) / 2 > ms ? 'miss' : 'ok'
  }
}
