// Nest settings (item_284) over core GET/POST /api/settings/nest and GET /api/settings/nest/build
// (items 282 and 283). The page edits a draft; one save sends the keys that changed.

export interface NestKnown { themes: string[]; heroes: string[]; pages: string[] }

// GET /api/settings/nest: crow.yaml modules.nest as written, so a broken value still shows.
export interface NestSettings {
  file: string
  theme: unknown
  hero: unknown
  pages: unknown
  known: NestKnown
  errors: unknown[]
}

export type BuildState = 'idle' | 'building' | 'ok' | 'failed'

export interface NestBuild {
  state: BuildState
  started_at?: string
  elapsed_s: number
  queued: boolean
  brick?: string
  error?: string
  log: string[]
  served: string
}

export interface NestSaved extends NestSettings {
  changed: string[]
  needs_build: boolean
  build: NestBuild
}

// A missing key keeps its value; pages [] shows every page.
export interface NestSave { theme?: string; hero?: string; pages?: string[] }

export interface SettingsApi {
  get(): Promise<NestSettings>
  save(body: NestSave): Promise<NestSaved>
  build(): Promise<NestBuild>
}

// Engine defaults when crow.yaml leaves a key out (nest-modules.ts).
export const DEFAULT_THEME = 'modern'
export const DEFAULT_HERO = 'holotable'
export const HOME = 'home'

export interface PageRow { key: string; on: boolean }
// mode is per device (localStorage, item_303): it rides in the draft for preview and Save, never in crow.yaml.
export type Mode = 'auto' | 'light' | 'dark'
export interface Draft { theme: string; hero: string; pages: PageRow[]; mode: Mode }

const str = (v: unknown, known: string[], dflt: string) => (typeof v === 'string' && v ? v : known.includes(dflt) ? dflt : known[0] ?? dflt)

// The draft crow.yaml holds: listed pages first in their order and on, then the other known
// pages off. No list shows every known page. Home is always on.
export function draftOf(s: NestSettings, mode: Mode = 'auto'): Draft {
  const known = s.known.pages
  const listed = Array.isArray(s.pages) ? (s.pages as unknown[]).filter((p): p is string => typeof p === 'string' && known.includes(p)) : []
  const order = listed.length ? [...new Set(listed)] : [...known]
  if (!order.includes(HOME) && known.includes(HOME)) order.unshift(HOME)
  const pages = [
    ...order.map((key) => ({ key, on: true })),
    ...known.filter((k) => !order.includes(k)).map((key) => ({ key, on: false })),
  ]
  return { theme: str(s.theme, s.known.themes, DEFAULT_THEME), hero: str(s.hero, s.known.heroes, DEFAULT_HERO), pages, mode }
}

const shown = (d: Draft) => d.pages.filter((p) => p.on).map((p) => p.key)
const samePages = (a: Draft, b: Draft) => shown(a).join() === shown(b).join()

export type SettingKey = 'theme' | 'hero' | 'pages' | 'mode'

// The keys the draft changes, in save order.
export function changedKeys(base: Draft, draft: Draft): SettingKey[] {
  const out: SettingKey[] = []
  if (draft.theme !== base.theme) out.push('theme')
  if (draft.hero !== base.hero) out.push('hero')
  if (!samePages(base, draft)) out.push('pages')
  if (draft.mode !== base.mode) out.push('mode')
  return out
}

// The body of one save: only what changed, mode left out (it stays on the device). Every page shown in the known order saves as [].
export function saveBody(base: Draft, draft: Draft, known: NestKnown): NestSave {
  const keys = changedKeys(base, draft)
  const body: NestSave = {}
  if (keys.includes('theme')) body.theme = draft.theme
  if (keys.includes('hero')) body.hero = draft.hero
  if (keys.includes('pages')) {
    const list = shown(draft)
    body.pages = list.join() === known.pages.join() ? [] : list
  }
  return body
}

// Hero and pages are compiled into Nest; theme and mode switch live.
export const needsBuild = (keys: readonly SettingKey[]) => keys.includes('hero') || keys.includes('pages')

// Count of pending changes as the owner sees them: theme, mode, hero, and each page moved or switched.
export function pendingCount(base: Draft, draft: Draft): number {
  let n = (draft.theme !== base.theme ? 1 : 0) + (draft.mode !== base.mode ? 1 : 0) + (draft.hero !== base.hero ? 1 : 0)
  for (const p of draft.pages) if (pageMark(base, draft, p.key)) n++
  return n
}

export type PageMark = '' | 'shown' | 'hidden' | 'moved up' | 'moved down'

// What a page row changed against crow.yaml: switched first, then its place among shown pages.
export function pageMark(base: Draft, draft: Draft, key: string): PageMark {
  const b = base.pages.find((p) => p.key === key)
  const d = draft.pages.find((p) => p.key === key)
  if (!b || !d) return ''
  if (b.on !== d.on) return d.on ? 'shown' : 'hidden'
  if (!d.on) return ''
  // Order among the pages shown before and after, so switching another page moves nothing.
  const both = new Set(shown(base).filter((k) => shown(draft).includes(k)))
  const b0 = shown(base).filter((k) => both.has(k))
  const d0 = shown(draft).filter((k) => both.has(k))
  const bi = b0.indexOf(key)
  const di = d0.indexOf(key)
  if (bi === di) return ''
  // A page that only shifted because another one moved past it is not marked: the rows that
  // keep their order (longest common run of the full lists) stay, the others moved.
  if (kept(base.pages.map((p) => p.key), draft.pages.map((p) => p.key)).has(key)) return ''
  return di < bi ? 'moved up' : 'moved down'
}

// Longest common subsequence of two orders of the same keys.
function kept(a: string[], b: string[]): Set<string> {
  const n = a.length
  const m = b.length
  const t = Array.from({ length: n + 1 }, () => new Array<number>(m + 1).fill(0))
  for (let i = n - 1; i >= 0; i--)
    for (let j = m - 1; j >= 0; j--) t[i][j] = a[i] === b[j] ? t[i + 1][j + 1] + 1 : Math.max(t[i + 1][j], t[i][j + 1])
  const out = new Set<string>()
  for (let i = 0, j = 0; i < n && j < m; ) {
    if (a[i] === b[j]) { out.add(a[i]); i++; j++ } else if (t[i + 1][j] >= t[i][j + 1]) i++
    else j++
  }
  return out
}

// Move a page one step (or to an index). Home stays first.
export function movePage(d: Draft, key: string, to: number): Draft {
  const from = d.pages.findIndex((p) => p.key === key)
  if (from < 0 || key === HOME) return d
  const min = d.pages[0]?.key === HOME ? 1 : 0
  const at = Math.max(min, Math.min(d.pages.length - 1, to))
  if (at === from) return d
  const pages = [...d.pages]
  const [row] = pages.splice(from, 1)
  pages.splice(at, 0, row)
  return { ...d, pages }
}

export function togglePage(d: Draft, key: string): Draft {
  if (key === HOME) return d
  return { ...d, pages: d.pages.map((p) => (p.key === key ? { ...p, on: !p.on } : p)) }
}

// crow.yaml errors as one line each (crow-modules reports {key, error} or plain text).
export function errorLines(errors: unknown[]): string[] {
  return errors.map((e) => {
    if (typeof e === 'string') return e
    const o = (e ?? {}) as { key?: string; error?: string; message?: string }
    const msg = o.error ?? o.message ?? JSON.stringify(e)
    return o.key ? `${o.key}: ${msg}` : msg
  })
}
