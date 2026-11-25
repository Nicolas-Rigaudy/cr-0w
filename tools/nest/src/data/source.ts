import type { ButlerApi } from './butler'
import type { ChatTarget, ConversationPage, PageAsk } from './conversation'
import type { Actions } from './live'
import type { Modules } from './modules'
import type { Changes, Ops, RunReport } from './ops'
import type { SettingsApi } from './settings'
import type { AgendaEvent, GateQuestion, InboxItem, Mission, SectionBody, SectionName } from './types'
import type { UsageData } from './usage'

// The only door between screens and data. Fixtures (fixtures.ts) and the core
// API (api.ts) implement it; main.tsx picks one.
export interface DataSource {
  missions(): Promise<Mission[]>
  mission(id: string): Promise<Mission | undefined>
  // The questions and options of one mission version (Inbox decisions, item_223).
  questions?(mission: string, version: number): Promise<GateQuestion[]>
  // The URL of one mission artifact, an image to compare in a design review (item_211).
  artifact?(mission: string, path: string): string
  inbox(): Promise<InboxItem[]>
  agenda(): Promise<AgendaEvent[]>
  // One script section. An unchanged section returns the same object.
  section(name: SectionName): Promise<SectionBody>
  // Bricks and switches (item_237). Absent: every section is on, no generic sections.
  modules?(): Promise<Modules>
  // A generic brick section by source name (/api/state/<source>).
  brickSection?(source: string): Promise<SectionBody>
  // A page of a session's or a mission's conversation (item_177).
  conversation(target: ChatTarget, page: PageAsk): Promise<ConversationPage>
  // Ops timeline for `day=YYYY-MM-DD` or `week=YYYY-Www` (item_190), and one run's report.
  ops(query: string): Promise<Ops>
  run(id: string): Promise<RunReport | undefined>
  // The git changes of a live lane's checkout (item_273). Undefined outside a registered repo.
  changes(cwd: string): Promise<Changes | undefined>
  // The Usage page data (item_227). Absent: the page says it has no data.
  usage?(): Promise<UsageData>
  // Live actions (item_170). Callers refetch after each one; nothing is optimistic.
  actions: Actions
  // Reference time for "today" and elapsed clocks. Fixtures pin it, the API uses the wall clock.
  now(): Date
  // Butler chat (item_176). Absent: no rail, no sheet.
  butler?: ButlerApi
  // Nest settings and the rebuild after a save (item_284). Absent: the Settings page says so.
  settings?: SettingsApi
  // Refresh period in ms while the page is visible. Absent: load once.
  poll?: number
}
