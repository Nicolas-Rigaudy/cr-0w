import type { SectionBody } from './types'

// Bricks as core serves them on GET /api/modules (item_237, adr_002). Nest builds its tabs
// and its generic sections from this list; nothing is loaded at runtime.

export interface ModuleAction {
  id: string
  label: string
  confirm?: boolean
  // A row action runs on one row and sends its id as item.
  row?: boolean
}

export interface ModuleSection {
  id: string
  page: string
  source: string
  title?: string
  view?: 'list'
  actions?: ModuleAction[]
}

export interface ModuleBrick {
  name: string
  on: boolean
  user: boolean
  sections: string[]
  nest: { sections?: ModuleSection[] }
  deck: { pages: string[]; boxes: string[] }
}

export interface Modules {
  contract_version: number
  version: string
  modules: ModuleBrick[]
  errors: Array<{ file: string; key?: string; error: string }>
}

// Pages that hold the core screens: they show whatever bricks are off.
export const CORE_PAGES: ReadonlySet<string> = new Set(['home', 'desk', 'ops', 'inbox', 'agenda', 'reports'])

// Core sections of off bricks. Nest does not ask core for them.
export function offSections(m: Modules | null): ReadonlySet<string> {
  const off = new Set<string>()
  for (const b of m?.modules ?? []) if (!b.on) for (const s of b.sections) off.add(s)
  return off
}

// A brick page shows when at least one of its sections is on.
export function hiddenPages(m: Modules | null): ReadonlySet<string> {
  const on = new Set<string>()
  const named = new Set<string>()
  for (const b of m?.modules ?? []) {
    for (const s of b.nest.sections ?? []) {
      named.add(s.page)
      if (b.on) on.add(s.page)
    }
  }
  return new Set([...named].filter((p) => !on.has(p) && !CORE_PAGES.has(p)))
}

export interface ListSection extends ModuleSection {
  brick: string
}

// Sections drawn with the generic list view, by page, from bricks that are on.
export function listSections(m: Modules | null, page?: string): ListSection[] {
  const out: ListSection[] = []
  for (const b of m?.modules ?? []) {
    if (!b.on) continue
    for (const s of b.nest.sections ?? []) {
      if (s.view === 'list' && (page === undefined || s.page === page)) out.push({ ...s, brick: b.name })
    }
  }
  return out
}

export interface ListRow {
  id: string
  title: string
  detail?: string
  meta?: string
  url?: string
}

const str = (v: unknown) => (typeof v === 'string' || typeof v === 'number' ? String(v) : undefined)

// List view rows: the section data is an array of rows or {items: [...]}. A row needs an id
// and shows title (else label, name, id), then detail and meta when present.
export function listRows(body: SectionBody | undefined): ListRow[] {
  const d = body?.data
  const items = Array.isArray(d) ? d : d && typeof d === 'object' && Array.isArray((d as { items?: unknown }).items) ? (d as { items: unknown[] }).items : []
  const rows: ListRow[] = []
  for (const it of items) {
    if (!it || typeof it !== 'object') continue
    const r = it as Record<string, unknown>
    const id = str(r.id)
    if (!id) continue
    const url = str(r.url)
    rows.push({
      id,
      title: str(r.title) ?? str(r.label) ?? str(r.name) ?? id,
      detail: str(r.detail),
      meta: str(r.meta),
      url: url && /^https?:\/\//.test(url) ? url : undefined,
    })
  }
  return rows
}
