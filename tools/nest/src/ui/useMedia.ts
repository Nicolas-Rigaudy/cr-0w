import { useEffect, useState } from 'react'

export const LAPTOP = '(min-width: 1200px)'

export function useMedia(query: string) {
  const [on, setOn] = useState(() => window.matchMedia(query).matches)
  useEffect(() => {
    const mq = window.matchMedia(query)
    const sync = () => setOn(mq.matches)
    sync()
    mq.addEventListener('change', sync)
    return () => mq.removeEventListener('change', sync)
  }, [query])
  return on
}
