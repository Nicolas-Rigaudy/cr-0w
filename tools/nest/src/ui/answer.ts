import { useState } from 'react'
import { ApiError } from '../data/api'
import type { InboxItem, InboxVerb, Mission } from '../data/types'
import { useSource } from '../data/useData'

// Answering an inbox row inline (item_168). A 409 refetches; when the mission moved on,
// the newer question takes its place and the note stays.

export const REJECT: InboxVerb[] = ['no', 'abandon']
export const verbLabel = (v: InboxVerb) => v[0].toUpperCase() + v.slice(1)
// The question part of a row id (<mission>/v<N>/<question>), empty for the version row.
const question = (r: InboxItem) => r.id.slice(r.mission.length).replace(/^\/v\d+/, '')

// The open row to answer now: the one asked for, else the same question on a newer version.
export function current(inbox: InboxItem[], asked: InboxItem): InboxItem | undefined {
  const same = inbox.find((r) => r.id === asked.id && r.state === 'open')
  if (same) return same
  return inbox
    .filter((r) => r.mission === asked.mission && r.state === 'open' && r.type !== 'fyi' && question(r) === question(asked))
    .sort((a, b) => b.version - a.version)[0]
}

export function useAnswer(asked: InboxItem | undefined, inbox: InboxItem[], missions: Mission[], refresh: () => Promise<void>) {
  const act = useSource().actions
  const [text, setText] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [moved, setMoved] = useState(false)
  const [key, setKey] = useState(asked?.id)
  // A new selection starts clean.
  if (asked?.id !== key) {
    setKey(asked?.id); setText(''); setError(''); setMoved(false)
  }
  const item = asked && current(inbox, asked)
  const m = asked && missions.find((x) => x.id === asked.mission)
  const newer = moved && !!item && !!asked && item.id !== asked.id

  const send = async (verb: InboxVerb) => {
    if (!item || !m) return
    setBusy(true)
    setError('')
    try {
      await act.missionAnswer(m.id, item.id, verb, text, m.version)
      setText('')
      await refresh()
    } catch (e) {
      if (e instanceof ApiError && e.status === 409) {
        await refresh()
        setMoved(true)
      }
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setBusy(false)
    }
  }
  return { item, mission: m, text, setText, busy, error, newer, send }
}
