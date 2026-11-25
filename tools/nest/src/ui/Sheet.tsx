import { useEffect, useId, useRef, useState } from 'react'
import { createPortal } from 'react-dom'

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
  // A second action shown before the main one (item_279: Restart run next to Stop).
  alt?: { label: string; run: () => Promise<void> }
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
  const go = (run: () => Promise<void>) => async () => {
    setBusy(true)
    setError('')
    try {
      await run()
      onClose()
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e))
      setBusy(false)
    }
  }
  const body = (
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
          {c.alt && <button type="button" className="btn" onClick={go(c.alt.run)} disabled={busy}>{c.alt.label}</button>}
          <button type="button" className={`btn ${c.danger ? 'xx' : 'p'}`} onClick={go(c.run)} disabled={busy} aria-busy={busy}>
            {busy ? 'Working…' : c.ok}
          </button>
        </div>
      </div>
    </>
  )
  // Mounted on the body (item_279): a sheet opened inside a clipped panel (.pn under
  // the cyber clip-path) was cut off with its action button.
  return typeof document === 'undefined' ? body : createPortal(body, document.body)
}
