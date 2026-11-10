import type { DataSource } from './source'
import type { AgendaEvent, InboxItem, Mission } from './types'

// Live adapter over the core read API (item_144). Same origin, no CORS. Records in
// `data` already follow the core contract types, so the mapping is a type cast.

interface Section<T> {
  section: string
  version: string
  updated_at: string
  data: T[] | null
}

interface MissionDetail {
  mission: Mission
}

export class ApiError extends Error {
  constructor(readonly status: number, message: string) {
    super(message)
  }
}

const NO_AGENDA: AgendaEvent[] = []

type Fetch = (input: string, init?: RequestInit) => Promise<Response>

async function failure(res: Response): Promise<ApiError> {
  let msg = res.statusText
  try { msg = ((await res.json()) as { error?: string }).error ?? msg } catch { /* body not JSON */ }
  return new ApiError(res.status, msg)
}

// A section remembers its ETag and last body. A 304 returns the same array
// reference, so callers can skip work when nothing changed.
function section<T>(get: Fetch, url: string) {
  let etag = ''
  let last: T[] = []
  return async (): Promise<T[]> => {
    const res = await get(url, { headers: etag ? { 'If-None-Match': etag } : {}, cache: 'no-store' })
    if (res.status === 304) return last
    if (!res.ok) throw await failure(res)
    const body = (await res.json()) as Section<T>
    etag = res.headers.get('ETag') ?? ''
    last = body.data ?? []
    return last
  }
}

export function apiSource(base = '/api', get: Fetch = (i, init) => fetch(i, init)): DataSource {
  const missions = section<Mission>(get, `${base}/state/missions`)
  // all=1 keeps answered asks, which mission detail lists as history.
  const inbox = section<InboxItem>(get, `${base}/state/inbox?all=1`)
  return {
    missions,
    inbox,
    async mission(id) {
      const res = await get(`${base}/missions/${encodeURIComponent(id)}`, { cache: 'no-store' })
      if (res.status === 404) return undefined
      if (!res.ok) throw await failure(res)
      return ((await res.json()) as MissionDetail).mission
    },
    // No agenda route in the core yet. One shared array keeps snapshots unchanged.
    agenda: async () => NO_AGENDA,
    now: () => new Date(),
    poll: 5000,
  }
}
