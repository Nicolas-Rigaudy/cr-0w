import { useCallback, useEffect, useRef, useState } from 'react'
import type { ChatTarget, Entry } from './conversation'
import { useSource } from './useData'

// The chat of one target (item_177): the last 50 entries, then a read every few seconds
// from the cursor that appends what came after. A reset reloads the tail. "Earlier"
// pages back with the before cursor. Hidden pages stop reading.
export const CHAT_POLL = 3000

export interface ConversationView {
  entries: Entry[]
  account: string
  error: string
  loading: boolean
  more: boolean
  earlier: () => Promise<void>
  refresh: () => Promise<void>
}

export function useConversation(target: ChatTarget): ConversationView {
  const src = useSource()
  const key = 'pid' in target ? `p${target.pid}` : `m${target.mission}`
  const [entries, setEntries] = useState<Entry[]>([])
  const [account, setAccount] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [more, setMore] = useState(false)
  const st = useRef({ next: '', before: '', key: '', busy: false })
  const tgt = useRef(target)
  tgt.current = target

  const read = useCallback(async () => {
    const s = st.current
    if (s.busy) return
    s.busy = true
    const k = s.key
    try {
      const p = await src.conversation(tgt.current, s.next ? { cursor: s.next } : { tail: 50 })
      if (st.current.key !== k) return
      if (p.reset_required) {
        s.next = ''
        return
      }
      const first = !s.next
      if (first) {
        // cdx has_more on a tail page speaks of later turns; older ones exist when a
        // before cursor comes back.
        s.before = p.before_cursor ?? ''
        setMore(!!s.before)
        setEntries(p.entries)
      } else if (p.entries.length) {
        setEntries((l) => {
          const seen = new Set(l.map((e) => e.id))
          return [...l, ...p.entries.filter((e) => !seen.has(e.id))]
        })
      }
      if (p.next_cursor) s.next = p.next_cursor
      setAccount(p.account)
      setError('')
    } catch (e) {
      if (st.current.key === k) setError(e instanceof Error ? e.message : String(e))
    } finally {
      s.busy = false
      if (st.current.key === k) setLoading(false)
    }
  }, [src])

  const earlier = useCallback(async () => {
    const s = st.current
    if (!s.before) return
    const k = s.key
    try {
      const p = await src.conversation(tgt.current, { before: s.before })
      if (st.current.key !== k) return
      s.before = p.before_cursor ?? ''
      setMore(!!s.before && p.entries.length > 0)
      setEntries((l) => {
        const seen = new Set(l.map((e) => e.id))
        return [...p.entries.filter((e) => !seen.has(e.id)), ...l]
      })
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
    }
  }, [src])

  useEffect(() => {
    st.current = { next: '', before: '', key, busy: false }
    setEntries([]); setError(''); setLoading(true); setMore(false)
    void read()
    let timer: ReturnType<typeof setInterval> | undefined
    const start = () => { timer ??= setInterval(read, CHAT_POLL) }
    const stop = () => { clearInterval(timer); timer = undefined }
    const onVis = () => {
      if (document.visibilityState === 'hidden') return stop()
      void read()
      start()
    }
    if (src.poll && document.visibilityState !== 'hidden') start()
    if (src.poll) document.addEventListener('visibilitychange', onVis)
    return () => {
      st.current.key = ''
      stop()
      document.removeEventListener('visibilitychange', onVis)
    }
  }, [key, read, src])

  return { entries, account, error, loading, more, earlier, refresh: read }
}
