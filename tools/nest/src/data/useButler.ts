import { createContext, createElement, useCallback, useContext, useEffect, useRef, useState, type ReactNode } from 'react'
import { mergeMsgs, type ButlerApi, type ButlerMsg, type Surface } from './butler'
import { useSource } from './useData'

export interface Butler {
  messages: ButlerMsg[]
  partial: string
  busy: boolean
  account: string
  error: string
  // account: undefined keeps the current one, '' means auto.
  say(text: string, surface: Surface, account?: string): Promise<void>
  // Start over: CR-0W forgets the conversation and the view empties.
  clear(): Promise<void>
  // Phone and Fold chat sheet: the header badge opens it, a phone message opens it too.
  sheet: boolean
  setSheet(open: boolean): void
}

const BUSY_MS = 1000
const IDLE_MS = 15000

// One conversation, polled by seq. Fast while a reply is on its way, slow when idle,
// stopped while the page is hidden.
function useButlerState(api: ButlerApi | undefined, conv: string): Butler | null {
  const [msgs, setMsgs] = useState<ButlerMsg[]>([])
  const [partial, setPartial] = useState('')
  const [busy, setBusy] = useState(false)
  const [account, setAccount] = useState('')
  const [error, setError] = useState('')
  const [sheet, setSheet] = useState(false)
  const last = useRef(0)
  const cleared = useRef(0)
  const kick = useRef<() => void>(() => {})

  useEffect(() => {
    if (!api) return
    let live = true
    let timer: ReturnType<typeof setTimeout> | undefined
    let fast = false
    const tick = async () => {
      clearTimeout(timer)
      try {
        const v = await api.get(conv, last.current)
        if (!live) return
        // Cleared elsewhere (the phone, another tab): drop what the clear hid.
        if ((v.cleared ?? 0) > cleared.current) {
          const c = (cleared.current = v.cleared ?? 0)
          setMsgs((h) => h.filter((m) => m.seq > c))
        }
        if (v.messages.length) {
          last.current = Math.max(last.current, ...v.messages.map((m) => m.seq))
          setMsgs((h) => mergeMsgs(h, v.messages))
        }
        setPartial(v.partial ?? '')
        setBusy(v.busy)
        setAccount(v.account ?? '')
        setError('')
        fast = v.busy
      } catch (e) {
        if (live) setError((e as Error).message)
      }
      if (live && document.visibilityState !== 'hidden') timer = setTimeout(tick, fast ? BUSY_MS : IDLE_MS)
    }
    kick.current = () => { fast = true; void tick() }
    const onVis = () => { if (document.visibilityState !== 'hidden') void tick(); else clearTimeout(timer) }
    void tick()
    document.addEventListener('visibilitychange', onVis)
    return () => {
      live = false
      clearTimeout(timer)
      document.removeEventListener('visibilitychange', onVis)
    }
  }, [api, conv])

  const say = useCallback(async (text: string, surface: Surface, acct?: string) => {
    if (!api) return
    if (surface === 'phone') setSheet(true)
    try {
      const m = await api.say(conv, text, surface, acct)
      setMsgs((h) => mergeMsgs(h, [m]))
      setBusy(true)
      setError('')
      kick.current()
    } catch (e) {
      setError((e as Error).message)
    }
  }, [api, conv])

  const clear = useCallback(async () => {
    if (!api) return
    try {
      const v = await api.clear(conv)
      cleared.current = v.cleared ?? v.last
      last.current = Math.max(last.current, v.last)
      setMsgs([])
      setPartial('')
      setBusy(false)
      setError('')
    } catch (e) {
      setError((e as Error).message)
    }
  }, [api, conv])

  if (!api) return null
  return { messages: msgs, partial, busy, account, error, say, clear, sheet, setSheet }
}

const ButlerContext = createContext<Butler | null>(null)

// The rail, the phone sheet and Desk share one poller on conversation "main".
export function ButlerProvider({ children }: { children: ReactNode }) {
  const b = useButlerState(useSource().butler, 'main')
  return createElement(ButlerContext.Provider, { value: b }, children)
}

export const useButler = () => useContext(ButlerContext)
