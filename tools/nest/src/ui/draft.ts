// Typed text survives a tab change: screens remount on every route, so each input keeps its
// draft per key in memory and in sessionStorage (a reload of the page keeps it too).
import { useCallback, useState } from 'react'

const mem = new Map<string, string>()
const SK = 'nest.draft.'

export function readDraft(key: string): string {
  if (mem.has(key)) return mem.get(key)!
  try { return sessionStorage.getItem(SK + key) ?? '' } catch { return '' }
}

export function writeDraft(key: string, text: string) {
  if (text) mem.set(key, text); else mem.delete(key)
  try { if (text) sessionStorage.setItem(SK + key, text); else sessionStorage.removeItem(SK + key) } catch { /* storage off */ }
}

// useState('') for an input, kept per key. A new key shows that key's own draft.
export function useDraft(key: string): [string, (t: string) => void] {
  const [s, setS] = useState(() => ({ key, text: readDraft(key) }))
  const cur = s.key === key ? s : { key, text: readDraft(key) }
  if (cur !== s) setS(cur)
  const set = useCallback((t: string) => { writeDraft(key, t); setS({ key, text: t }) }, [key])
  return [cur.text, set]
}
