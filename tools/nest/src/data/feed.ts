import { listSections, offSections, type Modules } from './modules'
import type { DataSource } from './source'
import { SCRIPT_SECTIONS, type AgendaEvent, type InboxItem, type Mission, type SectionBody, type SectionName } from './types'
import { withLinkedStates, workOf } from './work'

export interface Snapshot {
  missions: Mission[]
  inbox: InboxItem[]
  agenda: AgendaEvent[]
  // Script sections by name, each with its own stale flag.
  sections: Record<SectionName, SectionBody>
  // Bricks and switches (item_237), null when the source has none. A section of an off
  // brick reads { off: true } and is never fetched.
  modules?: Modules | null
  // Generic list sections of bricks that are on, by source name.
  bricks?: Record<string, SectionBody>
  now: Date
  // Last good refresh. down is true while the latest refresh failed; the data
  // then stays as it was at synced (null when nothing ever loaded).
  synced: Date | null
  down: boolean
  // Core serves another Nest build than this page (item_283). The page keeps running:
  // the header offers a reload and the next route change loads the new build.
  newBuild?: boolean
}

const minute = (d: Date) => Math.floor(d.getTime() / 60000)

const NO_AGENDA: AgendaEvent[] = []
const OFF: SectionBody = { data: null, stale: false, off: true, updated_at: null }
const NO_BRICKS: Record<string, SectionBody> = {}

const NO_SECTIONS = Object.fromEntries(
  SCRIPT_SECTIONS.map((n) => [n, { data: null, stale: true, error: 'not loaded yet', updated_at: null }]),
) as Record<SectionName, SectionBody>

// One failing section (an older core without it, a script error) shows that panel empty and
// stale; only missions and inbox decide that the core is down. A repeat failure reuses the
// stale body so the snapshot does not churn.
async function readSection(src: DataSource, n: SectionName, prev: SectionBody | undefined): Promise<SectionBody> {
  try {
    return await src.section(n)
  } catch (e) {
    if (prev?.stale && prev.error) return prev
    return { ...(prev ?? NO_SECTIONS[n]), stale: true, error: e instanceof Error ? e.message : String(e) }
  }
}

async function readSections(src: DataSource, prev: Record<SectionName, SectionBody> | undefined, off: ReadonlySet<string>) {
  const bodies = await Promise.all(SCRIPT_SECTIONS.map((n) => (off.has(n) ? OFF : readSection(src, n, prev?.[n]))))
  if (prev && bodies.every((b, i) => prev[SCRIPT_SECTIONS[i]] === b)) return prev
  return Object.fromEntries(SCRIPT_SECTIONS.map((n, i) => [n, bodies[i]])) as Record<SectionName, SectionBody>
}

// Generic list sections, read like script sections: one failing shows stale.
async function readBricks(src: DataSource, mods: Modules | null, prev: Record<string, SectionBody> | undefined) {
  if (!src.brickSection) return NO_BRICKS
  const names = [...new Set(listSections(mods).map((s) => s.source))]
  if (!names.length) return NO_BRICKS
  const bodies = await Promise.all(names.map(async (n): Promise<SectionBody> => {
    try {
      return await src.brickSection!(n)
    } catch (e) {
      const p = prev?.[n]
      if (p?.stale && p.error) return p
      return { ...(p ?? { data: null, updated_at: null }), stale: true, error: e instanceof Error ? e.message : String(e) }
    }
  }))
  if (prev && Object.keys(prev).length === names.length && names.every((n, i) => prev[n] === bodies[i])) return prev
  return Object.fromEntries(names.map((n, i) => [n, bodies[i]]))
}

// The module list; a failure (an older core) keeps the last one.
async function readModules(src: DataSource, prev: Modules | null) {
  if (!src.modules) return null
  try {
    return await src.modules()
  } catch {
    return prev
  }
}

// A feed refreshes every section of a source and emits a new snapshot only when
// a section body, the link state or the clock minute changed. A call during a
// refresh waits for it and runs once more, so a refetch after an action never
// returns data read before the action. The first nest_build core reports is the one
// this page runs; a later different one sets newBuild for good (item_283).
export function feed(src: DataSource, emit: (s: Snapshot) => void) {
  let snap: Snapshot | null = null
  let served: string | undefined
  let newBuild = false
  let lastRaw: Mission[] | null = null
  let ok: Date | null = null
  let running: Promise<void> | null = null
  let again: Promise<void> | null = null
  return function refresh(): Promise<void> {
    if (running) return (again ??= running.then(() => { again = null; return refresh() }))
    return (running = once().finally(() => { running = null }))
  }
  async function once() {
    const now = src.now()
    const same = snap !== null && minute(snap.now) === minute(now)
    try {
      const modules = src.modules ? await readModules(src, snap?.modules ?? null) : null
      const build = modules?.nest_build
      served ??= build
      const fresh = !newBuild && !!build && build !== served
      if (fresh) newBuild = true
      const off = offSections(modules)
      const [raw, inbox, agendaRaw, sections, bricks] = await Promise.all([
        src.missions(), src.inbox(), src.agenda(), readSections(src, snap?.sections, off), readBricks(src, modules, snap?.bricks),
      ])
      const agenda = off.has('today') ? NO_AGENDA : agendaRaw
      ok = now
      const missions = raw === lastRaw && sections === snap?.sections ? snap!.missions : withLinkedStates(raw, workOf(sections.work))
      lastRaw = raw
      if (same && !fresh && !snap!.down && snap!.missions === missions && snap!.inbox === inbox && snap!.agenda === agenda && snap!.sections === sections && snap!.modules === modules && snap!.bricks === bricks) return
      snap = { missions, inbox, agenda, sections, modules, bricks, now, synced: now, down: false, ...(newBuild && { newBuild }) }
    } catch {
      if (same && snap!.down) return
      snap = { ...(snap ?? { missions: [], inbox: [], agenda: [], sections: NO_SECTIONS, modules: null, bricks: NO_BRICKS }), now, synced: ok, down: true }
    }
    emit(snap)
  }
}
