import { useEffect, useId, useRef, useState } from 'react'

// The confirm sheet (item_170, variant A): what happens, the exact text, the
// target, then Cancel and the action. A bottom sheet on the cover screen, a side
// panel from 600px. item_168's answer sheet reuses it.
export interface Confirm {
  title: string
  lead: string
  // The exact text that will be sent, shown as written.
  quote?: string
  target: Array<[string, string]>
  note: string
  ok: string
  // Heavy action: red button. The tap that opened the sheet was the first one.
  danger?: boolean
  run: () => Promise<void>
}

export function Sheet({ c, onClose }: { c: Confirm; onClose: () => void }) {
  const id = useId()
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const first = useRef<HTMLButtonElement>(null)
  useEffect(() => {
    first.current?.focus()
    const key = (e: KeyboardEvent) => { if (e.key === 'Escape' && !busy) onClose() }
    addEventListener('keydown', key)
    return () => removeEventListener('keydown', key)
  }, [busy, onClose])
  const go = async () => {
    setBusy(true)
    setError('')
    try {
      await c.run()
      onClose()
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
      setBusy(false)
    }
  }
  return (
    <>
      <div className="scrim" onClick={busy ? undefined : onClose} aria-hidden="true" />
      <div className="sheet" role="dialog" aria-modal="true" aria-labelledby={id}>
        <div className="grab" aria-hidden="true" />
        <h2 id={id} className="sh">{c.title}</h2>
        <div>{c.lead}</div>
        {c.quote !== undefined && <div className="quote">{c.quote}</div>}
        <dl className="props">
          {c.target.map(([k, v]) => [<dt key={`${k}-t`}>{k}</dt>, <dd key={`${k}-d`}>{v}</dd>])}
        </dl>
        <div className="sub">{c.note}</div>
        {error && <div className="err" role="alert">{error}</div>}
        <div className="btns">
          <button ref={first} type="button" className="btn" onClick={onClose} disabled={busy}>Cancel</button>
          <button type="button" className={`btn ${c.danger ? 'xx' : 'p'}`} onClick={go} disabled={busy} aria-busy={busy}>
            {busy ? 'Working…' : c.ok}
          </button>
        </div>
      </div>
    </>
  )
}
