import { useEffect, useState } from 'react'
import { useSource } from '../data/useData'
import type { UsageData } from '../data/usage'

// Usage page (item_227): opened from a header gauge on /usage/<account>.

export function useUsage(): [UsageData | null, string] {
  const src = useSource()
  const [data, setData] = useState<UsageData | null>(null)
  const [err, setErr] = useState('')
  useEffect(() => {
    let live = true
    if (!src.usage) return
    const load = () => src.usage!().then((d) => { if (live) { setData(d); setErr('') } }, (e: Error) => { if (live) setErr(e.message) })
    void load()
    const t = src.poll ? setInterval(load, Math.max(src.poll, 30_000)) : undefined
    return () => { live = false; clearInterval(t) }
  }, [src])
  return [data, src.usage ? err : 'no usage data']
}

export function Usage({ account }: { account?: string }) {
  const [data, err] = useUsage()
  return (
    <div className="usage">
      <h1 className="sr-only">Usage{account ? ` · ${account}` : ''}</h1>
      {!data && <p className="empty">{err || 'Loading usage…'}</p>}
    </div>
  )
}
