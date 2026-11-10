import type { AgendaEvent, InboxItem, Mission, SectionBody, SectionName } from './types'

// The only door between screens and data. Fixtures (fixtures.ts) and the core
// API (api.ts) implement it; main.tsx picks one.
export interface DataSource {
  missions(): Promise<Mission[]>
  mission(id: string): Promise<Mission | undefined>
  inbox(): Promise<InboxItem[]>
  agenda(): Promise<AgendaEvent[]>
  // One script section. An unchanged section returns the same object.
  section(name: SectionName): Promise<SectionBody>
  // Reference time for "today" and elapsed clocks. Fixtures pin it, the API uses the wall clock.
  now(): Date
  // Refresh period in ms while the page is visible. Absent: load once.
  poll?: number
}
